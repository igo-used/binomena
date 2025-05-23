package audit

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/igo-used/binomena/core"
	"github.com/igo-used/binomena/database"
)

// AuditServiceDB provides database-backed blockchain security auditing
type AuditServiceDB struct {
	blockchain interface{} // Accept any blockchain implementation
	mu         sync.RWMutex
}

// NewAuditServiceWithDB creates a new database-backed audit service
func NewAuditServiceWithDB(blockchain interface{}) *AuditServiceDB {
	service := &AuditServiceDB{
		blockchain: blockchain,
	}

	// Start background audit tasks
	go service.runPeriodicAudits()

	return service
}

// LogEvent logs an audit event to the database
func (a *AuditServiceDB) LogEvent(level SecurityLevel, eventType, message string, data interface{}) {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Create event ID with "AdNe" prefix
	idHash := sha256.Sum256([]byte(fmt.Sprintf("%d%s%s%v", time.Now().UnixNano(), eventType, message, data)))
	id := "AdNe" + hex.EncodeToString(idHash[:])[:60]

	// Serialize data to JSON
	var dataJSON string
	if data != nil {
		dataBytes, err := json.Marshal(data)
		if err != nil {
			log.Printf("Failed to serialize audit data: %v", err)
			dataJSON = "{\"error\":\"serialization_failed\"}"
		} else {
			dataJSON = string(dataBytes)
		}
	} else {
		dataJSON = "null"
	}

	// Create audit event for database
	dbEvent := database.AuditEvent{
		EventID:   id,
		Timestamp: time.Now().Unix(),
		Level:     int(level),
		Type:      eventType,
		Message:   message,
		Data:      dataJSON,
	}

	// Save to database
	if err := database.DB.Create(&dbEvent).Error; err != nil {
		log.Printf("Failed to save audit event: %v", err)
		return
	}

	// Log critical events immediately
	if level == CriticalLevel {
		log.Printf("[CRITICAL AUDIT] %s: %s", eventType, message)
	}
}

// GetEvents returns all audit events from the database
func (a *AuditServiceDB) GetEvents() []AuditEvent {
	a.mu.RLock()
	defer a.mu.RUnlock()

	var dbEvents []database.AuditEvent
	if err := database.DB.Order("timestamp desc").Find(&dbEvents).Error; err != nil {
		log.Printf("Failed to get audit events: %v", err)
		return []AuditEvent{}
	}

	events := make([]AuditEvent, len(dbEvents))
	for i, dbEvent := range dbEvents {
		// Parse data JSON back to interface{}
		var data interface{}
		if dbEvent.Data != "" {
			if err := json.Unmarshal([]byte(dbEvent.Data), &data); err != nil {
				log.Printf("Failed to parse audit data: %v", err)
				data = dbEvent.Data // Use raw string if JSON parsing fails
			}
		}

		events[i] = AuditEvent{
			ID:        dbEvent.EventID,
			Timestamp: dbEvent.Timestamp,
			Level:     SecurityLevel(dbEvent.Level),
			Type:      dbEvent.Type,
			Message:   dbEvent.Message,
			Data:      data,
		}
	}

	return events
}

// GetEventsByLevel returns audit events filtered by security level
func (a *AuditServiceDB) GetEventsByLevel(level SecurityLevel) []AuditEvent {
	a.mu.RLock()
	defer a.mu.RUnlock()

	var dbEvents []database.AuditEvent
	if err := database.DB.Where("level = ?", int(level)).Order("timestamp desc").Find(&dbEvents).Error; err != nil {
		log.Printf("Failed to get audit events by level: %v", err)
		return []AuditEvent{}
	}

	events := make([]AuditEvent, len(dbEvents))
	for i, dbEvent := range dbEvents {
		// Parse data JSON back to interface{}
		var data interface{}
		if dbEvent.Data != "" {
			if err := json.Unmarshal([]byte(dbEvent.Data), &data); err != nil {
				log.Printf("Failed to parse audit data: %v", err)
				data = dbEvent.Data // Use raw string if JSON parsing fails
			}
		}

		events[i] = AuditEvent{
			ID:        dbEvent.EventID,
			Timestamp: dbEvent.Timestamp,
			Level:     SecurityLevel(dbEvent.Level),
			Type:      dbEvent.Type,
			Message:   dbEvent.Message,
			Data:      data,
		}
	}

	return events
}

// AuditBlockchain performs a comprehensive audit of the blockchain
func (a *AuditServiceDB) AuditBlockchain() {
	a.LogEvent(InfoLevel, "BlockchainAudit", "Starting comprehensive blockchain audit", nil)

	// Try to cast blockchain to the interface we need
	switch bc := a.blockchain.(type) {
	case interface{ GetChain() []core.Block }:
		a.auditBlocks(bc)
	default:
		a.LogEvent(WarningLevel, "BlockchainAudit", "Blockchain type not supported for detailed auditing", nil)
	}

	a.LogEvent(InfoLevel, "BlockchainAudit", "Blockchain audit completed", nil)
}

// auditBlocks audits individual blocks
func (a *AuditServiceDB) auditBlocks(blockchain interface{ GetChain() []core.Block }) {
	chain := blockchain.GetChain()

	for i, block := range chain {
		// Verify block hash
		calculatedHash := core.CalculateHash(block)
		if calculatedHash != block.Hash {
			a.LogEvent(CriticalLevel, "BlockHashMismatch",
				fmt.Sprintf("Block %d has invalid hash", block.Index),
				map[string]interface{}{
					"blockIndex":   block.Index,
					"expectedHash": calculatedHash,
					"actualHash":   block.Hash,
				})
		}

		// Verify block sequence
		if i > 0 {
			prevBlock := chain[i-1]
			if block.PreviousHash != prevBlock.Hash {
				a.LogEvent(CriticalLevel, "BlockSequenceError",
					fmt.Sprintf("Block %d has invalid previous hash", block.Index),
					map[string]interface{}{
						"blockIndex":       block.Index,
						"expectedPrevHash": prevBlock.Hash,
						"actualPrevHash":   block.PreviousHash,
					})
			}

			if block.Index != prevBlock.Index+1 {
				a.LogEvent(CriticalLevel, "BlockIndexError",
					fmt.Sprintf("Block %d has invalid index sequence", block.Index),
					map[string]interface{}{
						"blockIndex":    block.Index,
						"expectedIndex": prevBlock.Index + 1,
					})
			}
		}

		// Audit transactions in the block
		for _, tx := range block.Data {
			if len(tx.ID) < 4 || tx.ID[:4] != "AdNe" {
				a.LogEvent(ErrorLevel, "InvalidTransactionID",
					fmt.Sprintf("Transaction %s has invalid ID format", tx.ID),
					map[string]interface{}{
						"transactionId": tx.ID,
						"blockIndex":    block.Index,
					})
			}
		}
	}
}

// runPeriodicAudits runs periodic security audits
func (a *AuditServiceDB) runPeriodicAudits() {
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			a.LogEvent(InfoLevel, "PeriodicAudit", "Running scheduled security audit", nil)
			a.AuditBlockchain()
		}
	}
}
