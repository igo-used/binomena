package test

import (
	"strings"
	"testing"

	"github.com/binomena/wallet"
)

func TestWalletCreation(t *testing.T) {
	// Create a new wallet
	w, err := wallet.NewWallet()
	if err != nil {
		t.Fatalf("Failed to create wallet: %v", err)
	}
	
	// Check that the address has the correct prefix
	if !strings.HasPrefix(w.Address, "AdNe") {
		t.Errorf("Expected wallet address to start with 'AdNe', got %s", w.Address)
	}
	
	// Check that the private key is not empty
	privateKey := w.ExportPrivateKey()
	if privateKey == "" {
		t.Errorf("Expected non-empty private key")
	}
	
	// Test signing
	message := []byte("test message")
	signature, err := w.Sign(message)
	if err != nil {
		t.Fatalf("Failed to sign message: %v", err)
	}
	
	// Verify signature
	valid := wallet.VerifySignature(w.PublicKey, message, signature)
	if !valid {
		t.Errorf("Signature verification failed")
	}
}

func TestWalletImport(t *testing.T) {
	// Create a new wallet
	originalWallet, err := wallet.NewWallet()
	if err != nil {
		t.Fatalf("Failed to create wallet: %v", err)
	}
	
	// Export private key
	privateKey := originalWallet.ExportPrivateKey()
	
	// Import wallet from private key
	importedWallet, err := wallet.ImportPrivateKey(privateKey)
	if err != nil {
		t.Fatalf("Failed to import wallet: %v", err)
	}
	
	// Check that the addresses match
	if importedWallet.Address != originalWallet.Address {
		t.Errorf("Expected imported wallet address to be %s, got %s", originalWallet.Address, importedWallet.Address)
	}
	
	// Test signing with imported wallet
	message := []byte("test message")
	signature, err := importedWallet.Sign(message)
	if err != nil {
		t.Fatalf("Failed to sign message with imported wallet: %v", err)
	}
	
	// Verify signature
	valid := wallet.VerifySignature(importedWallet.PublicKey, message, signature)
	if !valid {
		t.Errorf("Signature verification failed for imported wallet")
	}
}
