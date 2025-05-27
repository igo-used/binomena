package core

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"
)

// ExecutionMode represents the transaction execution mode
type ExecutionMode int

const (
	// SingleThreaded execution processes transactions sequentially
	SingleThreaded ExecutionMode = iota
	// MultiThreaded execution processes transactions in parallel
	MultiThreaded
)

// ExecutionConfig holds configuration for the transaction execution engine
type ExecutionConfig struct {
	// DelegateThreshold is the minimum number of delegates required to enable multithreading
	DelegateThreshold int
	// MaxWorkers is the maximum number of worker goroutines for parallel execution
	MaxWorkers int
	// BatchSize is the number of transactions to process in each batch
	BatchSize int
	// Timeout is the maximum time to wait for transaction execution
	Timeout time.Duration
	// EnableIntegrityChecks enables additional state integrity checks
	EnableIntegrityChecks bool
}

// DefaultExecutionConfig returns the default execution configuration
func DefaultExecutionConfig() *ExecutionConfig {
	return &ExecutionConfig{
		DelegateThreshold:     11,               // Enable multithreading when > 11 delegates
		MaxWorkers:            runtime.NumCPU(), // Use all available CPU cores
		BatchSize:             100,              // Process 100 transactions per batch
		Timeout:               30 * time.Second, // 30 second timeout
		EnableIntegrityChecks: true,             // Enable integrity checks by default
	}
}

// TransactionResult represents the result of transaction execution
type TransactionResult struct {
	Transaction *Transaction
	Success     bool
	Error       error
	StateHash   string // Hash of state after transaction execution
}

// ExecutionEngine manages transaction execution with support for parallel processing
type ExecutionEngine struct {
	config      *ExecutionConfig
	mode        ExecutionMode
	workerPool  chan struct{}
	resultsChan chan TransactionResult
	mu          sync.RWMutex
	isRunning   bool
	ctx         context.Context
	cancel      context.CancelFunc
	stateLocker sync.RWMutex // Protects global state during execution
}

// NewExecutionEngine creates a new transaction execution engine
func NewExecutionEngine(config *ExecutionConfig) *ExecutionEngine {
	if config == nil {
		config = DefaultExecutionConfig()
	}

	ctx, cancel := context.WithCancel(context.Background())

	engine := &ExecutionEngine{
		config:      config,
		mode:        SingleThreaded, // Start in single-threaded mode
		workerPool:  make(chan struct{}, config.MaxWorkers),
		resultsChan: make(chan TransactionResult, config.BatchSize*2),
		ctx:         ctx,
		cancel:      cancel,
	}

	log.Printf("Transaction execution engine initialized - Mode: %s, Max Workers: %d, Delegate Threshold: %d",
		engine.getModeString(), config.MaxWorkers, config.DelegateThreshold)

	return engine
}

// UpdateMode updates the execution mode based on the number of active delegates
func (e *ExecutionEngine) UpdateMode(activeDelegateCount int) {
	e.mu.Lock()
	defer e.mu.Unlock()

	oldMode := e.mode

	if activeDelegateCount > e.config.DelegateThreshold {
		e.mode = MultiThreaded
	} else {
		e.mode = SingleThreaded
	}

	if oldMode != e.mode {
		log.Printf("Execution mode changed: %s -> %s (Active Delegates: %d, Threshold: %d)",
			e.getModeString(oldMode), e.getModeString(e.mode), activeDelegateCount, e.config.DelegateThreshold)
	}
}

// GetMode returns the current execution mode
func (e *ExecutionEngine) GetMode() ExecutionMode {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.mode
}

// ExecuteTransactions executes a batch of transactions based on the current mode
func (e *ExecutionEngine) ExecuteTransactions(transactions []Transaction, blockchain BlockchainInterface, tokenSystem interface{}) ([]TransactionResult, error) {
	if len(transactions) == 0 {
		return []TransactionResult{}, nil
	}

	e.mu.RLock()
	currentMode := e.mode
	e.mu.RUnlock()

	log.Printf("Executing %d transactions in %s mode", len(transactions), e.getModeString(currentMode))

	start := time.Now()
	var results []TransactionResult
	var err error

	if currentMode == MultiThreaded {
		results, err = e.executeParallel(transactions, blockchain, tokenSystem)
	} else {
		results, err = e.executeSequential(transactions, blockchain, tokenSystem)
	}

	duration := time.Since(start)
	log.Printf("Transaction execution completed in %v - Processed: %d, Successful: %d, Failed: %d",
		duration, len(results), e.countSuccessful(results), e.countFailed(results))

	return results, err
}

// executeSequential processes transactions sequentially (single-threaded)
func (e *ExecutionEngine) executeSequential(transactions []Transaction, blockchain BlockchainInterface, tokenSystem interface{}) ([]TransactionResult, error) {
	results := make([]TransactionResult, 0, len(transactions))

	for i, tx := range transactions {
		result := e.executeTransaction(&tx, blockchain, tokenSystem, fmt.Sprintf("seq-%d", i))
		results = append(results, result)

		// Perform integrity check if enabled
		if e.config.EnableIntegrityChecks {
			if err := e.performIntegrityCheck(blockchain, tokenSystem); err != nil {
				log.Printf("Integrity check failed after transaction %s: %v", tx.ID, err)
				result.Success = false
				result.Error = fmt.Errorf("integrity check failed: %v", err)
			}
		}
	}

	return results, nil
}

// executeParallel processes transactions in parallel (multi-threaded)
func (e *ExecutionEngine) executeParallel(transactions []Transaction, blockchain BlockchainInterface, tokenSystem interface{}) ([]TransactionResult, error) {
	numTransactions := len(transactions)
	results := make([]TransactionResult, numTransactions)
	var wg sync.WaitGroup

	// Process transactions in batches to maintain some ordering and reduce contention
	batchSize := e.config.BatchSize
	if batchSize > numTransactions {
		batchSize = numTransactions
	}

	for i := 0; i < numTransactions; i += batchSize {
		end := i + batchSize
		if end > numTransactions {
			end = numTransactions
		}

		batch := transactions[i:end]
		batchResults := make([]TransactionResult, len(batch))

		// Process batch in parallel
		for j, tx := range batch {
			wg.Add(1)

			go func(index int, transaction *Transaction, batchIndex int) {
				defer wg.Done()

				// Acquire worker slot
				select {
				case e.workerPool <- struct{}{}:
					defer func() { <-e.workerPool }()
				case <-e.ctx.Done():
					batchResults[batchIndex] = TransactionResult{
						Transaction: transaction,
						Success:     false,
						Error:       fmt.Errorf("execution cancelled"),
					}
					return
				}

				// Execute transaction with state locking for critical sections
				result := e.executeTransactionWithLocking(transaction, blockchain, tokenSystem, fmt.Sprintf("par-%d-%d", index, batchIndex))
				batchResults[batchIndex] = result
			}(i+j, &tx, j)
		}

		wg.Wait()

		// Copy batch results to main results array
		copy(results[i:end], batchResults)

		// Perform integrity check after each batch if enabled
		if e.config.EnableIntegrityChecks {
			if err := e.performIntegrityCheck(blockchain, tokenSystem); err != nil {
				log.Printf("Integrity check failed after batch %d-%d: %v", i, end-1, err)
				// Mark remaining transactions as failed
				for k := i; k < end; k++ {
					if !results[k].Success {
						results[k].Error = fmt.Errorf("batch integrity check failed: %v", err)
					}
				}
			}
		}
	}

	return results, nil
}

// executeTransaction executes a single transaction
func (e *ExecutionEngine) executeTransaction(tx *Transaction, blockchain BlockchainInterface, tokenSystem interface{}, executionID string) TransactionResult {
	log.Printf("[%s] Executing transaction %s: %s -> %s (%.6f)", executionID, tx.ID, tx.From, tx.To, tx.Amount)

	result := TransactionResult{
		Transaction: tx,
		Success:     false,
	}

	// Validate transaction
	if err := e.validateTransaction(tx); err != nil {
		result.Error = fmt.Errorf("validation failed: %v", err)
		return result
	}

	// Execute transaction (placeholder - implement actual business logic)
	if err := e.applyTransaction(tx, blockchain, tokenSystem); err != nil {
		result.Error = fmt.Errorf("execution failed: %v", err)
		return result
	}

	// Calculate state hash for integrity checking
	result.StateHash = e.calculateStateHash(blockchain, tokenSystem)
	result.Success = true

	log.Printf("[%s] Transaction %s executed successfully", executionID, tx.ID)
	return result
}

// executeTransactionWithLocking executes a transaction with proper locking for parallel execution
func (e *ExecutionEngine) executeTransactionWithLocking(tx *Transaction, blockchain BlockchainInterface, tokenSystem interface{}, executionID string) TransactionResult {
	// For parallel execution, we need to be more careful about state modifications
	// This is a simplified version - in production, you'd want more sophisticated locking

	e.stateLocker.Lock()
	defer e.stateLocker.Unlock()

	return e.executeTransaction(tx, blockchain, tokenSystem, executionID)
}

// validateTransaction performs basic transaction validation
func (e *ExecutionEngine) validateTransaction(tx *Transaction) error {
	if tx == nil {
		return fmt.Errorf("transaction is nil")
	}

	if tx.ID == "" {
		return fmt.Errorf("transaction ID is empty")
	}

	if tx.From == "" || tx.To == "" {
		return fmt.Errorf("from or to address is empty")
	}

	if tx.Amount <= 0 {
		return fmt.Errorf("invalid transaction amount: %f", tx.Amount)
	}

	if tx.From[:4] != "AdNe" || tx.To[:4] != "AdNe" {
		return fmt.Errorf("addresses must start with 'AdNe'")
	}

	return nil
}

// applyTransaction applies a transaction to the blockchain and token system
func (e *ExecutionEngine) applyTransaction(tx *Transaction, blockchain BlockchainInterface, tokenSystem interface{}) error {
	// Add transaction to blockchain
	if err := blockchain.AddTransaction(*tx); err != nil {
		return fmt.Errorf("failed to add transaction to blockchain: %v", err)
	}

	// Apply token transfer if token system supports it
	if transferer, ok := tokenSystem.(interface {
		Transfer(string, string, float64) error
	}); ok {
		if err := transferer.Transfer(tx.From, tx.To, tx.Amount); err != nil {
			return fmt.Errorf("failed to transfer tokens: %v", err)
		}
	}

	return nil
}

// calculateStateHash calculates a hash of the current state for integrity checking
func (e *ExecutionEngine) calculateStateHash(blockchain BlockchainInterface, tokenSystem interface{}) string {
	// This is a simplified version - in production, you'd want a more comprehensive state hash
	lastBlock := blockchain.GetLastBlock()
	pendingCount := len(blockchain.GetPendingTransactions())

	return fmt.Sprintf("%s-%d-%d", lastBlock.Hash, lastBlock.Index, pendingCount)
}

// performIntegrityCheck performs state integrity validation
func (e *ExecutionEngine) performIntegrityCheck(blockchain BlockchainInterface, tokenSystem interface{}) error {
	// Verify blockchain integrity
	chain := blockchain.GetChain()
	for i := 1; i < len(chain); i++ {
		if chain[i].PreviousHash != chain[i-1].Hash {
			return fmt.Errorf("blockchain integrity violation at block %d", i)
		}

		if chain[i].Index != chain[i-1].Index+1 {
			return fmt.Errorf("blockchain index integrity violation at block %d", i)
		}
	}

	// Additional integrity checks can be added here
	// e.g., token balance consistency, transaction signature verification, etc.

	return nil
}

// countSuccessful counts successful transaction results
func (e *ExecutionEngine) countSuccessful(results []TransactionResult) int {
	count := 0
	for _, result := range results {
		if result.Success {
			count++
		}
	}
	return count
}

// countFailed counts failed transaction results
func (e *ExecutionEngine) countFailed(results []TransactionResult) int {
	count := 0
	for _, result := range results {
		if !result.Success {
			count++
		}
	}
	return count
}

// getModeString returns a string representation of the execution mode
func (e *ExecutionEngine) getModeString(mode ...ExecutionMode) string {
	var m ExecutionMode
	if len(mode) > 0 {
		m = mode[0]
	} else {
		m = e.mode
	}

	switch m {
	case SingleThreaded:
		return "Single-Threaded"
	case MultiThreaded:
		return "Multi-Threaded"
	default:
		return "Unknown"
	}
}

// Shutdown gracefully shuts down the execution engine
func (e *ExecutionEngine) Shutdown() {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.isRunning {
		e.cancel()
		close(e.resultsChan)
		e.isRunning = false
		log.Println("Transaction execution engine shut down")
	}
}

// IsMultiThreaded returns true if the engine is currently in multi-threaded mode
func (e *ExecutionEngine) IsMultiThreaded() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.mode == MultiThreaded
}

// GetStats returns execution engine statistics
func (e *ExecutionEngine) GetStats() map[string]interface{} {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return map[string]interface{}{
		"mode":               e.getModeString(),
		"max_workers":        e.config.MaxWorkers,
		"batch_size":         e.config.BatchSize,
		"delegate_threshold": e.config.DelegateThreshold,
		"timeout":            e.config.Timeout.String(),
		"integrity_checks":   e.config.EnableIntegrityChecks,
		"is_running":         e.isRunning,
	}
}
