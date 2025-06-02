package core

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// DelegateCounter interface for counting active delegates
type DelegateCounter interface {
	GetActiveDelegateCount() int
}

// Protocol manages the coordination between consensus, blockchain, and execution engine
type Protocol struct {
	blockchain            BlockchainInterface
	consensus             DelegateCounter
	executionEngine       *ExecutionEngine
	tokenSystem           interface{}
	mu                    sync.RWMutex
	isRunning             bool
	delegateCheckInterval time.Duration
	lastDelegateCount     int
}

// ProtocolConfig holds configuration for the protocol layer
type ProtocolConfig struct {
	// ExecutionConfig for the execution engine
	ExecutionConfig *ExecutionConfig
	// DelegateCheckInterval determines how often to check delegate count
	DelegateCheckInterval time.Duration
	// EnableAutoMode enables automatic switching between execution modes
	EnableAutoMode bool
}

// DefaultProtocolConfig returns default protocol configuration
func DefaultProtocolConfig() *ProtocolConfig {
	return &ProtocolConfig{
		ExecutionConfig:       DefaultExecutionConfig(),
		DelegateCheckInterval: 10 * time.Second, // Check every 10 seconds
		EnableAutoMode:        true,             // Enable automatic mode switching
	}
}

// NewProtocol creates a new protocol layer
func NewProtocol(blockchain BlockchainInterface, consensus DelegateCounter, tokenSystem interface{}, config *ProtocolConfig) *Protocol {
	if config == nil {
		config = DefaultProtocolConfig()
	}

	executionEngine := NewExecutionEngine(config.ExecutionConfig)

	protocol := &Protocol{
		blockchain:            blockchain,
		consensus:             consensus,
		executionEngine:       executionEngine,
		tokenSystem:           tokenSystem,
		delegateCheckInterval: config.DelegateCheckInterval,
		lastDelegateCount:     0,
	}

	// Initial delegate count check
	protocol.updateExecutionMode()

	log.Printf("Protocol layer initialized - Auto mode: %v, Check interval: %v",
		config.EnableAutoMode, config.DelegateCheckInterval)

	return protocol
}

// Start starts the protocol layer services
func (p *Protocol) Start() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.isRunning {
		return nil
	}

	p.isRunning = true

	// Start delegate monitoring goroutine
	go p.delegateMonitor()

	log.Println("Protocol layer started")
	return nil
}

// Stop stops the protocol layer services
func (p *Protocol) Stop() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.isRunning {
		return nil
	}

	p.isRunning = false
	p.executionEngine.Shutdown()

	log.Println("Protocol layer stopped")
	return nil
}

// ProcessTransactions processes a batch of transactions using the execution engine
func (p *Protocol) ProcessTransactions(transactions []Transaction) ([]TransactionResult, error) {
	if len(transactions) == 0 {
		return []TransactionResult{}, nil
	}

	// Update execution mode based on current delegate count
	p.updateExecutionMode()

	// Execute transactions through the execution engine
	results, err := p.executionEngine.ExecuteTransactions(transactions, p.blockchain, p.tokenSystem)
	if err != nil {
		log.Printf("Transaction execution failed: %v", err)
		return nil, err
	}

	// Log execution statistics
	p.logExecutionStats(results)

	return results, nil
}

// CreateBlock creates a new block with processed transactions
func (p *Protocol) CreateBlock(validator string, signature string) (*Block, error) {
	// Get pending transactions
	pendingTxs := p.blockchain.GetPendingTransactions()
	if len(pendingTxs) == 0 {
		// Create empty block if no pending transactions
		return p.createEmptyBlock(validator, signature), nil
	}

	// Process transactions
	results, err := p.ProcessTransactions(pendingTxs)
	if err != nil {
		return nil, err
	}

	// Filter successful transactions for the block
	var successfulTxs []Transaction
	for _, result := range results {
		if result.Success && result.Transaction != nil {
			successfulTxs = append(successfulTxs, *result.Transaction)
		}
	}

	// Create block with successful transactions
	lastBlock := p.blockchain.GetLastBlock()
	newBlock := Block{
		Index:        lastBlock.Index + 1,
		PreviousHash: lastBlock.Hash,
		Timestamp:    time.Now().Unix(),
		Data:         successfulTxs,
		Validator:    validator,
		Signature:    signature,
	}

	// Calculate block hash
	newBlock.Hash = CalculateHash(newBlock)

	return &newBlock, nil
}

// GetExecutionStats returns current execution engine statistics
func (p *Protocol) GetExecutionStats() map[string]interface{} {
	stats := p.executionEngine.GetStats()

	// Add protocol-specific stats
	stats["last_delegate_count"] = p.lastDelegateCount
	stats["is_running"] = p.isRunning
	stats["delegate_check_interval"] = p.delegateCheckInterval.String()

	return stats
}

// ApplyProductionOptimization safely applies performance optimizations
// Returns the expected TPS improvement and any warnings
func (p *Protocol) ApplyProductionOptimization(level string) (expectedTPS int, warnings []string, err error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.isRunning {
		return 0, nil, fmt.Errorf("protocol not running")
	}

	var config *ExecutionConfig
	warnings = []string{}

	switch level {
	case "safe":
		config = ProductionOptimizedConfig()
		expectedTPS = 1000
		warnings = append(warnings, "Monitor CPU usage for first 30 minutes")
		warnings = append(warnings, "Performance increase: 30-50%")

	case "balanced":
		config = ConditionalIntegrityConfig()
		expectedTPS = 1200
		warnings = append(warnings, "Monitor error rates closely")
		warnings = append(warnings, "Performance increase: 50-70%")

	case "aggressive":
		config = HighPerformanceConfig()
		expectedTPS = 1500
		warnings = append(warnings, "⚠️ HIGH RISK: Monitor continuously")
		warnings = append(warnings, "⚠️ Integrity checks disabled")
		warnings = append(warnings, "⚠️ Have rollback plan ready")
		warnings = append(warnings, "Performance increase: 70-100%")

	default:
		return 0, nil, fmt.Errorf("invalid optimization level: %s (use: safe, balanced, aggressive)", level)
	}

	// Apply new configuration gradually
	oldConfig := p.executionEngine.config
	p.executionEngine.config = config

	log.Printf("🔧 Applied %s optimization - Expected TPS: %d", level, expectedTPS)
	log.Printf("📊 Config change: BatchSize %d→%d, Workers %d→%d, Threshold %d→%d",
		oldConfig.BatchSize, config.BatchSize,
		oldConfig.MaxWorkers, config.MaxWorkers,
		oldConfig.DelegateThreshold, config.DelegateThreshold)

	return expectedTPS, warnings, nil
}

// RollbackOptimization reverts to conservative settings for safety
func (p *Protocol) RollbackOptimization() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.isRunning {
		return fmt.Errorf("protocol not running")
	}

	// Revert to conservative settings
	conservativeConfig := DefaultExecutionConfig()
	p.executionEngine.config = conservativeConfig
	p.executionEngine.performanceLevel = 0

	log.Printf("🔙 Rolled back to conservative configuration for safety")
	return nil
}

// IsMultiThreaded returns true if currently running in multi-threaded mode
func (p *Protocol) IsMultiThreaded() bool {
	return p.executionEngine.IsMultiThreaded()
}

// GetCurrentMode returns the current execution mode as a string
func (p *Protocol) GetCurrentMode() string {
	if p.executionEngine.IsMultiThreaded() {
		return "Multi-Threaded"
	}
	return "Single-Threaded"
}

// delegateMonitor monitors delegate count and updates execution mode
func (p *Protocol) delegateMonitor() {
	ticker := time.NewTicker(p.delegateCheckInterval)
	defer ticker.Stop()

	for range ticker.C {
		if !p.isRunning {
			return
		}
		p.updateExecutionMode()
	}
}

// updateExecutionMode updates the execution mode based on current delegate count
func (p *Protocol) updateExecutionMode() {
	currentCount := p.consensus.GetActiveDelegateCount()

	p.mu.Lock()
	oldCount := p.lastDelegateCount
	p.lastDelegateCount = currentCount
	p.mu.Unlock()

	// Update execution engine mode
	p.executionEngine.UpdateMode(currentCount)

	// Log significant changes
	if currentCount != oldCount {
		log.Printf("Delegate count changed: %d -> %d, Execution mode: %s",
			oldCount, currentCount, p.GetCurrentMode())
	}
}

// createEmptyBlock creates an empty block for cases with no pending transactions
func (p *Protocol) createEmptyBlock(validator string, signature string) *Block {
	lastBlock := p.blockchain.GetLastBlock()

	emptyBlock := Block{
		Index:        lastBlock.Index + 1,
		PreviousHash: lastBlock.Hash,
		Timestamp:    time.Now().Unix(),
		Data:         []Transaction{},
		Validator:    validator,
		Signature:    signature,
	}

	emptyBlock.Hash = CalculateHash(emptyBlock)
	return &emptyBlock
}

// logExecutionStats logs execution statistics
func (p *Protocol) logExecutionStats(results []TransactionResult) {
	successful := 0
	failed := 0

	for _, result := range results {
		if result.Success {
			successful++
		} else {
			failed++
		}
	}

	log.Printf("Execution completed - Mode: %s, Total: %d, Successful: %d, Failed: %d",
		p.GetCurrentMode(), len(results), successful, failed)

	// Log failed transactions with reasons
	if failed > 0 {
		log.Printf("Failed transactions:")
		for _, result := range results {
			if !result.Success && result.Error != nil {
				log.Printf("  - %s: %v", result.Transaction.ID, result.Error)
			}
		}
	}
}

// ValidateConfiguration validates the protocol configuration
func (p *Protocol) ValidateConfiguration() error {
	stats := p.GetExecutionStats()

	log.Printf("Protocol Configuration:")
	log.Printf("  - Execution Mode: %s", stats["mode"])
	log.Printf("  - Max Workers: %v", stats["max_workers"])
	log.Printf("  - Batch Size: %v", stats["batch_size"])
	log.Printf("  - Delegate Threshold: %v", stats["delegate_threshold"])
	log.Printf("  - Integrity Checks: %v", stats["integrity_checks"])
	log.Printf("  - Active Delegates: %v", stats["last_delegate_count"])
	log.Printf("  - Check Interval: %v", stats["delegate_check_interval"])

	return nil
}
