package test

import (
	"strings"
	"testing"

	"github.com/binomena/core"
	"github.com/binomena/wallet"
)

func TestTransactionCreation(t *testing.T) {
	// Create sender wallet
	senderWallet, err := wallet.NewWallet()
	if err != nil {
		t.Fatalf("Failed to create sender wallet: %v", err)
	}
	
	// Create receiver wallet
	receiverWallet, err := wallet.NewWallet()
	if err != nil {
		t.Fatalf("Failed to create receiver wallet: %v", err)
	}
	
	// Create transaction
	tx, err := core.NewTransaction(senderWallet.Address, receiverWallet.Address, 100.0, senderWallet)
	if err != nil {
		t.Fatalf("Failed to create transaction: %v", err)
	}
	
	// Check transaction ID prefix
	if !strings.HasPrefix(tx.ID, "AdNe") {
		t.Errorf("Expected transaction ID to start with 'AdNe', got %s", tx.ID)
	}
	
	// Check transaction fields
	if tx.From != senderWallet.Address {
		t.Errorf("Expected sender address to be %s, got %s", senderWallet.Address, tx.From)
	}
	
	if tx.To != receiverWallet.Address {
		t.Errorf("Expected receiver address to be %s, got %s", receiverWallet.Address, tx.To)
	}
	
	if tx.Amount != 100.0 {
		t.Errorf("Expected amount to be 100.0, got %f", tx.Amount)
	}
	
	// Check transaction fee
	fee := tx.CalculateFee()
	expectedFee := 0.1 // 0.1% of 100.0
	if fee != expectedFee {
		t.Errorf("Expected fee to be %f, got %f", expectedFee, fee)
	}
}

func TestTransactionVerification(t *testing.T) {
	// Create sender wallet
	senderWallet, err := wallet.NewWallet()
	if err != nil {
		t.Fatalf("Failed to create sender wallet: %v", err)
	}
	
	// Create receiver wallet
	receiverWallet, err := wallet.NewWallet()
	if err != nil {
		t.Fatalf("Failed to create receiver wallet: %v", err)
	}
	
	// Create transaction
	tx, err := core.NewTransaction(senderWallet.Address, receiverWallet.Address, 100.0, senderWallet)
	if err != nil {
		t.Fatalf("Failed to create transaction: %v", err)
	}
	
	// Verify transaction
	valid := core.VerifyTransaction(tx, senderWallet.PublicKey)
	if !valid {
		t.Errorf("Transaction verification failed")
	}
	
	// Verify with wrong public key (should fail)
	valid = core.VerifyTransaction(tx, receiverWallet.PublicKey)
	if valid {
		t.Errorf("Transaction verification should have failed with wrong public key")
	}
}
