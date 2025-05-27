package main

import (
	"fmt"
	"log"
	"time"

	"github.com/igo-used/binomena/consensus"
	"github.com/igo-used/binomena/core"
	"github.com/igo-used/binomena/token"
)

func main() {
	log.Println("=== Threaded Transaction Execution Engine Demo ===")

	// Initialize blockchain
	blockchain := core.NewBlockchain()
	log.Println("✓ Blockchain initialized")

	// Initialize token system
	tokenSystem := token.NewBinomToken()
	log.Println("✓ Token system initialized")

	// Initialize DPoS consensus
	founderAddress := "AdNe6c3ce54e4371d056c7c566675ba16909eb2e9534"
	communityAddress := "AdNebaefd75d426056bffbc622bd9f334ed89450efae"
	dposConsensus := consensus.NewDPoSConsensus(founderAddress, communityAddress)
	log.Println("✓ DPoS consensus initialized")

	// Register founder as the first delegate
	if err := dposConsensus.RegisterDelegate(founderAddress, 400000000.0); err != nil {
		log.Printf("Warning: Failed to register founder as delegate: %v", err)
	} else {
		log.Println("✓ Founder registered as first delegate")
	}

	// Create protocol with custom configuration
	protocolConfig := &core.ProtocolConfig{
		ExecutionConfig: &core.ExecutionConfig{
			DelegateThreshold:     11,               // Enable multithreading when > 11 delegates
			MaxWorkers:            4,                // Use 4 worker goroutines
			BatchSize:             50,               // Process 50 transactions per batch
			Timeout:               30 * time.Second, // 30 second timeout
			EnableIntegrityChecks: true,             // Enable integrity checks
		},
		DelegateCheckInterval: 5 * time.Second, // Check delegate count every 5 seconds
		EnableAutoMode:        true,            // Enable automatic mode switching
	}

	protocol := core.NewProtocol(blockchain, dposConsensus, tokenSystem, protocolConfig)
	log.Println("✓ Protocol layer initialized")

	// Start the protocol layer
	if err := protocol.Start(); err != nil {
		log.Fatalf("Failed to start protocol: %v", err)
	}
	defer protocol.Stop()
	log.Println("✓ Protocol layer started")

	// Validate and display initial configuration
	log.Println("\n=== Initial Configuration ===")
	protocol.ValidateConfiguration()

	// Demonstrate execution in single-threaded mode (< 11 delegates)
	log.Println("\n=== Phase 1: Single-Threaded Mode (Low Delegate Count) ===")

	activeDelegates := dposConsensus.GetActiveDelegateCount()
	log.Printf("Current active delegates: %d", activeDelegates)
	log.Printf("Current execution mode: %s", protocol.GetCurrentMode())

	// Create sample transactions
	sampleTransactions := createSampleTransactions(5)

	// Process transactions in single-threaded mode
	log.Println("Processing 5 transactions...")
	results, err := protocol.ProcessTransactions(sampleTransactions)
	if err != nil {
		log.Printf("Error processing transactions: %v", err)
	} else {
		logResults("Single-threaded", results)
	}

	// Simulate network growth by registering more delegates
	log.Println("\n=== Phase 2: Network Growth Simulation ===")
	log.Println("Simulating network growth by registering more delegates...")

	// Register additional delegates to trigger multi-threaded mode
	additionalDelegates := []string{
		"AdNetest1111111111111111111111111111111111111111",
		"AdNetest2222222222222222222222222222222222222222",
		"AdNetest3333333333333333333333333333333333333333",
		"AdNetest4444444444444444444444444444444444444444",
		"AdNetest5555555555555555555555555555555555555555",
		"AdNetest6666666666666666666666666666666666666666",
		"AdNetest7777777777777777777777777777777777777777",
		"AdNetest8888888888888888888888888888888888888888",
		"AdNetest9999999999999999999999999999999999999999",
		"AdNetestaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		"AdNetestbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
		"AdNetestcccccccccccccccccccccccccccccccccccccccc",
	}

	for i, delegateAddr := range additionalDelegates {
		stake := 10000.0 + float64(i)*1000.0 // Varying stake amounts
		if err := dposConsensus.RegisterDelegate(delegateAddr, stake); err != nil {
			log.Printf("Failed to register delegate %s: %v", delegateAddr, err)
		} else {
			log.Printf("✓ Registered delegate %d: %s (stake: %.0f)", i+2, delegateAddr, stake)
		}
	}

	// Wait for delegate monitor to detect the change
	log.Println("Waiting for automatic mode switch detection...")
	time.Sleep(6 * time.Second)

	// Check new delegate count and execution mode
	log.Println("\n=== Phase 3: Multi-Threaded Mode (High Delegate Count) ===")

	newActiveDelegates := dposConsensus.GetActiveDelegateCount()
	log.Printf("New active delegates count: %d", newActiveDelegates)
	log.Printf("New execution mode: %s", protocol.GetCurrentMode())

	if protocol.IsMultiThreaded() {
		log.Println("🚀 Multi-threaded execution is now ACTIVE!")
	} else {
		log.Println("⚠️  Still in single-threaded mode")
	}

	// Process more transactions in multi-threaded mode
	log.Println("Processing 20 transactions in parallel...")
	largerBatch := createSampleTransactions(20)

	results, err = protocol.ProcessTransactions(largerBatch)
	if err != nil {
		log.Printf("Error processing transactions: %v", err)
	} else {
		logResults("Multi-threaded", results)
	}

	// Demonstrate block creation with processed transactions
	log.Println("\n=== Phase 4: Block Creation ===")

	// Add some pending transactions
	for _, tx := range createSampleTransactions(3) {
		blockchain.AddTransaction(tx)
	}

	// Create a block using the protocol
	block, err := protocol.CreateBlock(founderAddress, "demo_signature")
	if err != nil {
		log.Printf("Error creating block: %v", err)
	} else {
		log.Printf("✓ Block created successfully:")
		log.Printf("  - Index: %d", block.Index)
		log.Printf("  - Transactions: %d", len(block.Data))
		log.Printf("  - Hash: %s", block.Hash[:16]+"...")
		log.Printf("  - Validator: %s", block.Validator)
	}

	// Display final statistics
	log.Println("\n=== Final Statistics ===")
	stats := protocol.GetExecutionStats()
	for key, value := range stats {
		log.Printf("  - %s: %v", key, value)
	}

	log.Println("\n=== Demo Completed Successfully ===")
}

// createSampleTransactions creates sample transactions for testing
func createSampleTransactions(count int) []core.Transaction {
	transactions := make([]core.Transaction, count)

	fromAddress := "AdNetest1234567890abcdef1234567890abcdef12345678"

	for i := 0; i < count; i++ {
		toAddress := fmt.Sprintf("AdNetest%058d", i+1000)
		txID := fmt.Sprintf("AdNetest%058d", i+2000)

		transactions[i] = core.Transaction{
			ID:        txID,
			From:      fromAddress,
			To:        toAddress,
			Amount:    float64(10 + i*5), // Varying amounts
			Timestamp: time.Now().Unix(),
			Signature: fmt.Sprintf("signature_%d", i),
		}
	}

	return transactions
}

// logResults logs the results of transaction execution
func logResults(mode string, results []core.TransactionResult) {
	successful := 0
	failed := 0

	for _, result := range results {
		if result.Success {
			successful++
		} else {
			failed++
		}
	}

	log.Printf("📊 %s execution results:", mode)
	log.Printf("  - Total: %d", len(results))
	log.Printf("  - Successful: %d", successful)
	log.Printf("  - Failed: %d", failed)

	if failed > 0 {
		log.Printf("  - Failed transactions:")
		for _, result := range results {
			if !result.Success && result.Error != nil {
				log.Printf("    • %s: %v", result.Transaction.ID[:16]+"...", result.Error)
			}
		}
	}
}
