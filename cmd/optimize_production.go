package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/igo-used/binomena/consensus"
	"github.com/igo-used/binomena/core"
	"github.com/igo-used/binomena/token"
)

// ProductionOptimizer handles safe performance optimization for live servers
type ProductionOptimizer struct {
	protocol     *core.Protocol
	initialTPS   float64
	targetTPS    float64
	rollbackTime time.Time
	monitoring   bool
}

func main() {
	fmt.Println("🚀 Binomena Production Performance Optimizer")
	fmt.Println("===========================================")

	// Initialize system components
	optimizer, err := initializeOptimizer()
	if err != nil {
		log.Fatalf("Failed to initialize optimizer: %v", err)
	}

	// Setup graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\n🛑 Graceful shutdown initiated...")
		optimizer.emergencyRollback()
		os.Exit(0)
	}()

	// Start optimization process
	optimizer.runOptimizationProcess()
}

func initializeOptimizer() (*ProductionOptimizer, error) {
	// Initialize blockchain components
	blockchain := core.NewBlockchain()
	tokenSystem := token.NewBinomToken()

	// Initialize consensus
	founderAddr := "AdNe6c3ce54e4371d056c7c566675ba16909eb2e9534"
	communityAddr := "AdNebaefd75d426056bffbc622bd9f334ed89450efae"
	consensus := consensus.NewDPoSConsensus(founderAddr, communityAddr)

	// Create protocol with conservative settings
	protocol := core.NewProtocol(blockchain, consensus, tokenSystem, nil)

	// Start protocol
	if err := protocol.Start(); err != nil {
		return nil, fmt.Errorf("failed to start protocol: %v", err)
	}

	return &ProductionOptimizer{
		protocol:   protocol,
		monitoring: true,
	}, nil
}

func (po *ProductionOptimizer) runOptimizationProcess() {
	fmt.Println("📊 Starting performance optimization process...")

	// Step 1: Baseline measurement
	po.measureBaseline()

	// Step 2: Apply safe optimization
	po.applySafeOptimization()

	// Step 3: Monitor for 10 minutes
	po.monitorPerformance(10 * time.Minute)

	// Step 4: Optionally apply balanced optimization
	po.applyBalancedOptimization()

	// Step 5: Long-term monitoring
	po.longTermMonitoring()
}

func (po *ProductionOptimizer) measureBaseline() {
	fmt.Println("\n📏 Measuring baseline performance...")

	// Create test transactions
	testTxs := po.createTestTransactions(100)

	start := time.Now()
	results, err := po.protocol.ProcessTransactions(testTxs)
	duration := time.Since(start)

	if err != nil {
		log.Printf("Error during baseline: %v", err)
		return
	}

	successful := 0
	for _, result := range results {
		if result.Success {
			successful++
		}
	}

	po.initialTPS = float64(successful) / duration.Seconds()

	fmt.Printf("✅ Baseline TPS: %.2f\n", po.initialTPS)
	fmt.Printf("✅ Success rate: %.1f%%\n", float64(successful)/float64(len(results))*100)
}

func (po *ProductionOptimizer) applySafeOptimization() {
	fmt.Println("\n🔧 Applying SAFE optimization level...")

	expectedTPS, warnings, err := po.protocol.ApplyProductionOptimization("safe")
	if err != nil {
		log.Printf("Failed to apply optimization: %v", err)
		return
	}

	po.targetTPS = float64(expectedTPS)
	po.rollbackTime = time.Now().Add(30 * time.Minute) // Auto-rollback in 30 min if issues

	fmt.Printf("🎯 Target TPS: %d\n", expectedTPS)
	fmt.Println("⚠️ Warnings:")
	for _, warning := range warnings {
		fmt.Printf("   - %s\n", warning)
	}

	// Test immediately after optimization
	po.testOptimizedPerformance()
}

func (po *ProductionOptimizer) testOptimizedPerformance() {
	fmt.Println("\n🧪 Testing optimized performance...")

	testTxs := po.createTestTransactions(200) // Larger test batch

	start := time.Now()
	results, err := po.protocol.ProcessTransactions(testTxs)
	duration := time.Since(start)

	if err != nil {
		log.Printf("❌ Optimization test failed: %v", err)
		po.emergencyRollback()
		return
	}

	successful := 0
	for _, result := range results {
		if result.Success {
			successful++
		}
	}

	currentTPS := float64(successful) / duration.Seconds()
	improvement := ((currentTPS - po.initialTPS) / po.initialTPS) * 100

	fmt.Printf("✅ Optimized TPS: %.2f (%.1f%% improvement)\n", currentTPS, improvement)

	if improvement < 20 {
		fmt.Printf("⚠️ Improvement below expected (%.1f%%), investigating...\n", improvement)
	} else {
		fmt.Printf("🚀 Optimization successful! TPS improved by %.1f%%\n", improvement)
	}
}

func (po *ProductionOptimizer) monitorPerformance(duration time.Duration) {
	fmt.Printf("\n👁️ Monitoring performance for %v...\n", duration)

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	endTime := time.Now().Add(duration)

	for time.Now().Before(endTime) {
		select {
		case <-ticker.C:
			stats := po.protocol.GetExecutionStats()
			fmt.Printf("📊 [%s] Mode: %s, Active: %v\n",
				time.Now().Format("15:04:05"),
				stats["mode"],
				stats["is_running"])

			// Check if rollback is needed
			if time.Now().After(po.rollbackTime) {
				fmt.Println("⏰ Auto-rollback time reached")
				po.emergencyRollback()
				return
			}
		}
	}

	fmt.Println("✅ Monitoring period completed successfully")
}

func (po *ProductionOptimizer) applyBalancedOptimization() {
	fmt.Println("\n🚀 Applying BALANCED optimization level...")

	expectedTPS, warnings, err := po.protocol.ApplyProductionOptimization("balanced")
	if err != nil {
		log.Printf("Failed to apply balanced optimization: %v", err)
		return
	}

	fmt.Printf("🎯 New target TPS: %d\n", expectedTPS)
	fmt.Println("⚠️ Additional warnings:")
	for _, warning := range warnings {
		fmt.Printf("   - %s\n", warning)
	}

	// Test the balanced optimization
	po.testOptimizedPerformance()
}

func (po *ProductionOptimizer) longTermMonitoring() {
	fmt.Println("\n🔄 Starting long-term monitoring (Ctrl+C to stop)...")

	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Run periodic health check
			po.performHealthCheck()
		}
	}
}

func (po *ProductionOptimizer) performHealthCheck() {
	fmt.Printf("🏥 [%s] Health check...\n", time.Now().Format("15:04:05"))

	// Quick performance test
	testTxs := po.createTestTransactions(50)
	start := time.Now()
	results, err := po.protocol.ProcessTransactions(testTxs)
	duration := time.Since(start)

	if err != nil {
		fmt.Printf("❌ Health check failed: %v\n", err)
		po.emergencyRollback()
		return
	}

	successful := 0
	for _, result := range results {
		if result.Success {
			successful++
		}
	}

	currentTPS := float64(successful) / duration.Seconds()

	if currentTPS < po.initialTPS*0.8 { // Performance dropped below 80% of baseline
		fmt.Printf("⚠️ Performance degradation detected (%.2f TPS), rolling back...\n", currentTPS)
		po.emergencyRollback()
	} else {
		fmt.Printf("✅ System healthy - Current TPS: %.2f\n", currentTPS)
	}
}

func (po *ProductionOptimizer) emergencyRollback() {
	fmt.Println("\n🔙 EMERGENCY ROLLBACK INITIATED")

	err := po.protocol.RollbackOptimization()
	if err != nil {
		log.Printf("❌ Rollback failed: %v", err)
		fmt.Println("🚨 MANUAL INTERVENTION REQUIRED")
		return
	}

	fmt.Println("✅ System rolled back to conservative settings")
	fmt.Println("📧 Consider notifying operations team")
}

func (po *ProductionOptimizer) createTestTransactions(count int) []core.Transaction {
	transactions := make([]core.Transaction, count)

	for i := 0; i < count; i++ {
		transactions[i] = core.Transaction{
			ID:        fmt.Sprintf("AdNeoptimize%054d", i),
			From:      "AdNetest1234567890abcdef1234567890abcdef12345678",
			To:        fmt.Sprintf("AdNetest%058d", i+10000),
			Amount:    float64(1 + i%100),
			Timestamp: time.Now().Unix(),
			Signature: fmt.Sprintf("optimize_sig_%d", i),
		}
	}

	return transactions
}
