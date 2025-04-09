package core

import (
	"strings"
	"testing"

	"github.com/igo-used/binomena/wallet"
)

func TestTransactionCreation(t *testing.T) {
	// Create wallets for testing
	senderWallet, _ := wallet.NewWallet()
	receiverWallet, _ := wallet.NewWallet()

	// Create transaction
	tx, err := NewTransaction(
		senderWallet.Address,
		receiverWallet.Address,
		100.0,
		senderWallet,
	)

	if err != nil {
		t.Fatalf("Failed to create transaction: %v", err)
	}

	// Check transaction properties
	if !strings.HasPrefix(tx.ID, "AdNe") {
		t.Errorf("Transaction ID doesn't have correct prefix: %s", tx.ID)
	}

	if tx.From != senderWallet.Address {
		t.Errorf("Sender address mismatch: expected=%s, got=%s",
			senderWallet.Address, tx.From)
	}

	if tx.To != receiverWallet.Address {
		t.Errorf("Receiver address mismatch: expected=%s, got=%s",
			receiverWallet.Address, tx.To)
	}

	if tx.Amount != 100.0 {
		t.Errorf("Amount mismatch: expected=%.2f, got=%.2f", 100.0, tx.Amount)
	}

	// Verify signature
	isValid := VerifyTransaction(tx, senderWallet.PublicKey)
	if !isValid {
		t.Error("Transaction signature verification failed")
	}
}

func TestTransactionFee(t *testing.T) {
	// Create wallets for testing
	senderWallet, _ := wallet.NewWallet()
	receiverWallet, _ := wallet.NewWallet()

	// Create transaction
	tx, _ := NewTransaction(
		senderWallet.Address,
		receiverWallet.Address,
		1000.0,
		senderWallet,
	)

	// Calculate fee
	fee := tx.CalculateFee()
	expectedFee := 1000.0 * 0.001 // 0.1% fee

	if fee != expectedFee {
		t.Errorf("Fee calculation incorrect: expected=%.4f, got=%.4f",
			expectedFee, fee)
	}
}
