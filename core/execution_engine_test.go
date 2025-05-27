package core

import (
	"fmt"
	"testing"
	"time"
)

// MockDPoSConsensus implements DelegateCounter for testing
type MockDPoSConsensus struct {
	activeDelegateCount int
}

func (m *MockDPoSConsensus) GetActiveDelegateCount() int {
	return m.activeDelegateCount
}

func (m *MockDPoSConsensus) SetActiveDelegateCount(count int) {
	m.activeDelegateCount = count
}

// MockTokenSystem implements basic token transfer for testing
type MockTokenSystem struct {
	balances map[string]float64
}

func NewMockTokenSystem() *MockTokenSystem {
	return &MockTokenSystem{
		balances: make(map[string]float64),
	}
}

func (m *MockTokenSystem) Transfer(from, to string, amount float64) error {
	if m.balances[from] < amount {
		return fmt.Errorf("insufficient balance")
	}
	m.balances[from] -= amount
	m.balances[to] += amount
	return nil
}

func (m *MockTokenSystem) SetBalance(address string, balance float64) {
	m.balances[address] = balance
}

func TestExecutionEngine_SingleThreadedMode(t *testing.T) {
	// Create test blockchain
	blockchain := NewBlockchain()

	// Create test token system
	tokenSystem := NewMockTokenSystem()
	tokenSystem.SetBalance("AdNetest1234567890abcdef1234567890abcdef12345678", 1000.0)

	// Create execution engine with default config
	engine := NewExecutionEngine(nil)

	// Create test transactions
	transactions := []Transaction{
		{
			ID:        "AdNetest1234567890abcdef1234567890abcdef12345678901234567890",
			From:      "AdNetest1234567890abcdef1234567890abcdef12345678",
			To:        "AdNetest9876543210fedcba9876543210fedcba98765432",
			Amount:    100.0,
			Timestamp: time.Now().Unix(),
		},
		{
			ID:        "AdNetest9876543210fedcba9876543210fedcba98765432109876543210",
			From:      "AdNetest1234567890abcdef1234567890abcdef12345678",
			To:        "AdNetest1111111111111111111111111111111111111111",
			Amount:    50.0,
			Timestamp: time.Now().Unix(),
		},
	}

	// Update mode with low delegate count (should be single-threaded)
	engine.UpdateMode(5)

	if engine.GetMode() != SingleThreaded {
		t.Errorf("Expected SingleThreaded mode, got %v", engine.GetMode())
	}

	// Execute transactions
	results, err := engine.ExecuteTransactions(transactions, blockchain, tokenSystem)
	if err != nil {
		t.Fatalf("Failed to execute transactions: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}

	// Check that all transactions were processed
	for i, result := range results {
		if !result.Success {
			t.Errorf("Transaction %d failed: %v", i, result.Error)
		}
	}
}

func TestExecutionEngine_MultiThreadedMode(t *testing.T) {
	// Create test blockchain
	blockchain := NewBlockchain()

	// Create test token system
	tokenSystem := NewMockTokenSystem()
	tokenSystem.SetBalance("AdNetest1234567890abcdef1234567890abcdef12345678", 1000.0)

	// Create execution engine with custom config for testing
	config := &ExecutionConfig{
		DelegateThreshold:     5, // Lower threshold for testing
		MaxWorkers:            2, // Limit workers for testing
		BatchSize:             10,
		Timeout:               10 * time.Second,
		EnableIntegrityChecks: true,
	}
	engine := NewExecutionEngine(config)

	// Create test transactions
	transactions := []Transaction{
		{
			ID:        "AdNetest1234567890abcdef1234567890abcdef12345678901234567890",
			From:      "AdNetest1234567890abcdef1234567890abcdef12345678",
			To:        "AdNetest9876543210fedcba9876543210fedcba98765432",
			Amount:    100.0,
			Timestamp: time.Now().Unix(),
		},
		{
			ID:        "AdNetest9876543210fedcba9876543210fedcba98765432109876543210",
			From:      "AdNetest1234567890abcdef1234567890abcdef12345678",
			To:        "AdNetest1111111111111111111111111111111111111111",
			Amount:    50.0,
			Timestamp: time.Now().Unix(),
		},
	}

	// Update mode with high delegate count (should be multi-threaded)
	engine.UpdateMode(12)

	if engine.GetMode() != MultiThreaded {
		t.Errorf("Expected MultiThreaded mode, got %v", engine.GetMode())
	}

	// Execute transactions
	results, err := engine.ExecuteTransactions(transactions, blockchain, tokenSystem)
	if err != nil {
		t.Fatalf("Failed to execute transactions: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}

	// Check that all transactions were processed
	for i, result := range results {
		if !result.Success {
			t.Errorf("Transaction %d failed: %v", i, result.Error)
		}
	}
}

func TestExecutionEngine_ModeSwitch(t *testing.T) {
	engine := NewExecutionEngine(nil)

	// Start in single-threaded mode
	engine.UpdateMode(5)
	if engine.GetMode() != SingleThreaded {
		t.Errorf("Expected SingleThreaded mode, got %v", engine.GetMode())
	}

	// Switch to multi-threaded mode
	engine.UpdateMode(15)
	if engine.GetMode() != MultiThreaded {
		t.Errorf("Expected MultiThreaded mode, got %v", engine.GetMode())
	}

	// Switch back to single-threaded mode
	engine.UpdateMode(8)
	if engine.GetMode() != SingleThreaded {
		t.Errorf("Expected SingleThreaded mode, got %v", engine.GetMode())
	}
}

func TestProtocol_Integration(t *testing.T) {
	// Create test components
	blockchain := NewBlockchain()
	consensus := &MockDPoSConsensus{activeDelegateCount: 5}
	tokenSystem := NewMockTokenSystem()
	tokenSystem.SetBalance("AdNetest1234567890abcdef1234567890abcdef12345678", 1000.0)

	// Create protocol with custom config
	config := &ProtocolConfig{
		ExecutionConfig: &ExecutionConfig{
			DelegateThreshold:     10,
			MaxWorkers:            2,
			BatchSize:             5,
			Timeout:               5 * time.Second,
			EnableIntegrityChecks: true,
		},
		DelegateCheckInterval: 1 * time.Second,
		EnableAutoMode:        true,
	}

	protocol := NewProtocol(blockchain, consensus, tokenSystem, config)

	// Start protocol
	err := protocol.Start()
	if err != nil {
		t.Fatalf("Failed to start protocol: %v", err)
	}
	defer protocol.Stop()

	// Initially should be single-threaded (5 delegates < 10 threshold)
	if protocol.IsMultiThreaded() {
		t.Error("Expected single-threaded mode initially")
	}

	// Create test transactions
	transactions := []Transaction{
		{
			ID:        "AdNetest1234567890abcdef1234567890abcdef12345678901234567890",
			From:      "AdNetest1234567890abcdef1234567890abcdef12345678",
			To:        "AdNetest9876543210fedcba9876543210fedcba98765432",
			Amount:    100.0,
			Timestamp: time.Now().Unix(),
		},
	}

	// Process transactions in single-threaded mode
	results, err := protocol.ProcessTransactions(transactions)
	if err != nil {
		t.Fatalf("Failed to process transactions: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
	}

	// Increase delegate count to trigger multi-threaded mode
	consensus.SetActiveDelegateCount(15)

	// Give time for delegate monitor to detect change
	time.Sleep(2 * time.Second)

	// Now should be multi-threaded (15 delegates > 10 threshold)
	if !protocol.IsMultiThreaded() {
		t.Error("Expected multi-threaded mode after delegate increase")
	}

	// Process more transactions in multi-threaded mode
	moreTransactions := []Transaction{
		{
			ID:        "AdNetest9999999999999999999999999999999999999999999999999999",
			From:      "AdNetest1234567890abcdef1234567890abcdef12345678",
			To:        "AdNetest2222222222222222222222222222222222222222",
			Amount:    25.0,
			Timestamp: time.Now().Unix(),
		},
	}

	results, err = protocol.ProcessTransactions(moreTransactions)
	if err != nil {
		t.Fatalf("Failed to process transactions in multi-threaded mode: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
	}

	// Validate configuration
	err = protocol.ValidateConfiguration()
	if err != nil {
		t.Errorf("Configuration validation failed: %v", err)
	}
}

func TestExecutionEngine_Stats(t *testing.T) {
	config := &ExecutionConfig{
		DelegateThreshold:     11,
		MaxWorkers:            4,
		BatchSize:             50,
		Timeout:               30 * time.Second,
		EnableIntegrityChecks: true,
	}

	engine := NewExecutionEngine(config)

	stats := engine.GetStats()

	// Verify stats contain expected keys
	expectedKeys := []string{"mode", "max_workers", "batch_size", "delegate_threshold", "timeout", "integrity_checks"}
	for _, key := range expectedKeys {
		if _, exists := stats[key]; !exists {
			t.Errorf("Stats missing key: %s", key)
		}
	}

	// Verify specific values
	if stats["max_workers"] != 4 {
		t.Errorf("Expected max_workers to be 4, got %v", stats["max_workers"])
	}

	if stats["delegate_threshold"] != 11 {
		t.Errorf("Expected delegate_threshold to be 11, got %v", stats["delegate_threshold"])
	}
}

func BenchmarkExecutionEngine_Sequential(b *testing.B) {
	blockchain := NewBlockchain()
	tokenSystem := NewMockTokenSystem()
	tokenSystem.SetBalance("AdNetest1234567890abcdef1234567890abcdef12345678", 100000.0)

	engine := NewExecutionEngine(nil)
	engine.UpdateMode(5) // Force single-threaded mode

	// Create a batch of transactions
	transactions := make([]Transaction, 100)
	for i := 0; i < 100; i++ {
		transactions[i] = Transaction{
			ID:        fmt.Sprintf("AdNe%058d", i),
			From:      "AdNetest1234567890abcdef1234567890abcdef12345678",
			To:        fmt.Sprintf("AdNe%058d", i+1000),
			Amount:    1.0,
			Timestamp: time.Now().Unix(),
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := engine.ExecuteTransactions(transactions, blockchain, tokenSystem)
		if err != nil {
			b.Fatalf("Execution failed: %v", err)
		}
	}
}

func BenchmarkExecutionEngine_Parallel(b *testing.B) {
	blockchain := NewBlockchain()
	tokenSystem := NewMockTokenSystem()
	tokenSystem.SetBalance("AdNetest1234567890abcdef1234567890abcdef12345678", 100000.0)

	config := &ExecutionConfig{
		DelegateThreshold:     5,
		MaxWorkers:            4,
		BatchSize:             25,
		Timeout:               30 * time.Second,
		EnableIntegrityChecks: false, // Disable for performance
	}

	engine := NewExecutionEngine(config)
	engine.UpdateMode(15) // Force multi-threaded mode

	// Create a batch of transactions
	transactions := make([]Transaction, 100)
	for i := 0; i < 100; i++ {
		transactions[i] = Transaction{
			ID:        fmt.Sprintf("AdNe%058d", i),
			From:      "AdNetest1234567890abcdef1234567890abcdef12345678",
			To:        fmt.Sprintf("AdNe%058d", i+1000),
			Amount:    1.0,
			Timestamp: time.Now().Unix(),
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := engine.ExecuteTransactions(transactions, blockchain, tokenSystem)
		if err != nil {
			b.Fatalf("Execution failed: %v", err)
		}
	}
}
