package test

import (
	"testing"
	"time"

	"github.com/igo-used/binomena/core"
)

func TestBlockchain(t *testing.T) {
	// Create a new blockchain
	blockchain := core.NewBlockchain()

	// Check that the blockchain has a genesis block
	if blockchain.GetBlockCount() != 1 {
		t.Errorf("Expected blockchain to have 1 block, got %d", blockchain.GetBlockCount())
	}

	// Create a transaction
	tx := core.Transaction{
		ID:        "tx1",
		From:      "alice",
		To:        "bob",
		Amount:    100.0,
		Timestamp: time.Now().Unix(),
		Signature: "signature",
	}

	// Add transaction to blockchain
	blockchain.AddTransaction(tx)

	// Check that the transaction was added
	pendingTxs := blockchain.GetPendingTransactions()
	if len(pendingTxs) != 1 {
		t.Errorf("Expected 1 pending transaction, got %d", len(pendingTxs))
	}

	// Create a new block
	lastBlock := blockchain.GetLastBlock()
	newBlock := core.Block{
		Index:        lastBlock.Index + 1,
		PreviousHash: lastBlock.Hash,
		Timestamp:    time.Now().Unix(),
		Data:         pendingTxs,
		Validator:    "validator",
		Signature:    "signature",
	}

	// Calculate hash for the new block
	newBlock.Hash = calculateTestHash(newBlock)

	// Add block to blockchain
	err := blockchain.AddBlock(newBlock)
	if err != nil {
		t.Errorf("Failed to add block: %v", err)
	}

	// Check that the block was added
	if blockchain.GetBlockCount() != 2 {
		t.Errorf("Expected blockchain to have 2 blocks, got %d", blockchain.GetBlockCount())
	}

	// Check that pending transactions were cleared
	pendingTxs = blockchain.GetPendingTransactions()
	if len(pendingTxs) != 0 {
		t.Errorf("Expected 0 pending transactions, got %d", len(pendingTxs))
	}
}

// Helper function to calculate hash for testing
func calculateTestHash(_ core.Block) string {
	// This is a simplified version for testing
	return "testhash"
}
