package audit

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/binomena/core"
)

// SecurityLevel represents the security level of an audit
type SecurityLevel int

const (
	// InfoLevel represents informational audits
	InfoLevel SecurityLevel = iota
	// WarningLevel represents warning audits
	WarningLevel
	// ErrorLevel represents error audits
	ErrorLevel
	// CriticalLevel represents critical audits
	CriticalLevel
)

// AuditEvent represents a security audit event
type AuditEvent struct {
	ID        string        `json:"id"`
	Timestamp int64         `json:"timestamp"`
	Level     SecurityLevel `json:"level"`
	Type      string        `json:"type"`
	Message   string        `json:"message"`
	Data      interface{}   `json:"data,omitempty"`
}

// AuditService provides blockchain security auditing
type AuditService struct {
	events     []AuditEvent
	blockchain *core.Blockchain
	mu         sync.RWMutex
}

// NewAuditService creates a new audit service
func NewAuditService(blockchain *core.Blockchain) *AuditService {
	service := &AuditService{
		events:     make([]AuditEvent, 0),
		blockchain: blockchain,
	}

	// Start background audit tasks
	go service.runPeriodicAudits()

	return service
}

// LogEvent logs an audit event
func (a *AuditService) LogEvent(level SecurityLevel, eventType, message string, data interface{}) {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Create event ID with "AdNe" prefix
	idHash := sha256.Sum256([]byte(fmt.Sprintf("%d%s%s%v", time.Now().UnixNano(), eventType, message, data)))
	id := "AdNe" + hex.EncodeToString(idHash[:])[:60]

	// Create audit event
	event := AuditEvent{
		ID:        id,
		Timestamp: time.Now().Unix(),
		Level:     level,
		Type:      eventType,
		Message:   message,
		Data:      data,
	}

	// Add event to log
	a.events = append(a.events, event)

	// Log critical events immediately
	if level == CriticalLevel {
		log.Printf("[CRITICAL AUDIT] %s: %s", eventType, message)
	}
}

// GetEvents returns all audit events
func (a *AuditService) GetEvents() []AuditEvent {
	a.mu.RLock()
	defer a.mu.RUnlock()

	// Return a copy of the events
	eventsCopy := make([]AuditEvent, len(a.events))
	copy(eventsCopy, a.events)

	return eventsCopy
}

// GetEventsByLevel returns audit events filtered by security level
func (a *AuditService) GetEventsByLevel(level SecurityLevel) []AuditEvent {
	a.mu.RLock()
	defer a.mu.RUnlock()

	// Filter events by level
	filtered := make([]AuditEvent, 0)
	for _, event := range a.events {
		if event.Level == level {
			filtered = append(filtered, event)
		}
	}

	return filtered
}

// AuditBlockchain performs a full audit of the blockchain
func (a *AuditService) AuditBlockchain() []AuditEvent {
	// Get a copy of the blockchain
	chain := a.blockchain.GetChain()

	auditEvents := make([]AuditEvent, 0)

	// Verify each block
	for i := 1; i < len(chain); i++ {
		block := chain[i]
		prevBlock := chain[i-1]

		// Verify block index
		if block.Index != prevBlock.Index+1 {
			a.LogEvent(
				CriticalLevel,
				"BlockIndexMismatch",
				fmt.Sprintf("Block %d has invalid index", block.Index),
				block,
			)
		}

		// Verify previous hash
		if block.PreviousHash != prevBlock.Hash {
			a.LogEvent(
				CriticalLevel,
				"PreviousHashMismatch",
				fmt.Sprintf("Block %d has invalid previous hash", block.Index),
				block,
			)
		}

		// Verify block hash
		calculatedHash := core.CalculateHash(block)
		if calculatedHash != block.Hash {
			a.LogEvent(
				CriticalLevel,
				"HashMismatch",
				fmt.Sprintf("Block %d has invalid hash", block.Index),
				block,
			)
		}

		// Verify transactions
		for _, tx := range block.Data {
			// Verify transaction ID prefix
			if tx.ID[:4] != "AdNe" {
				a.LogEvent(
					ErrorLevel,
					"InvalidTransactionPrefix",
					fmt.Sprintf("Transaction %s has invalid prefix", tx.ID),
					tx,
				)
			}

			// Note: Full transaction signature verification would require public keys
			// which we don't store in this simplified implementation
		}
	}

	return auditEvents
}

// runPeriodicAudits runs periodic security audits
func (a *AuditService) runPeriodicAudits() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			a.AuditBlockchain()
			a.LogEvent(
				InfoLevel,
				"PeriodicAudit",
				"Completed periodic blockchain audit",
				nil,
			)
		}
	}
}
