package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"testing"
)

func TestWalletCreation(t *testing.T) {
	// Create a new wallet
	wallet, err := NewWallet()
	if err != nil {
		t.Fatalf("Failed to create wallet: %v", err)
	}

	// Check if address is properly formatted
	if len(wallet.Address) < 5 || wallet.Address[:4] != "AdNe" {
		t.Errorf("Invalid wallet address format: %s", wallet.Address)
	}

	// Check if private key is generated
	if wallet.PrivateKey == nil {
		t.Error("Private key not generated")
	}

	// Check if public key is derived
	if wallet.PublicKey == nil {
		t.Error("Public key not derived")
	}
}

func TestWalletSigning(t *testing.T) {
	// Create a wallet
	wallet, _ := NewWallet()

	// Test message
	message := []byte("Test message for signing")

	// Sign the message
	signature, err := wallet.Sign(message)
	if err != nil {
		t.Fatalf("Failed to sign message: %v", err)
	}

	// Verify the signature
	valid := VerifySignature(wallet.PublicKey, message, signature)
	if !valid {
		t.Error("Signature verification failed")
	}

	// Test with invalid signature
	invalidSignature := make([]byte, len(signature))
	copy(invalidSignature, signature)
	invalidSignature[0] = ^invalidSignature[0] // Flip bits to invalidate

	valid = VerifySignature(wallet.PublicKey, message, invalidSignature)
	if valid {
		t.Error("Invalid signature was verified as valid")
	}
}

// WalletFromPrivateKey creates a wallet from an existing private key
func WalletFromPrivateKey(privateKey *ecdsa.PrivateKey) (*Wallet, error) {
	if privateKey == nil {
		return nil, errors.New("private key cannot be nil")
	}

	// Derive public key
	publicKey := &privateKey.PublicKey

	// Generate address from public key
	address, err := generateAddressFromPublicKey(publicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to generate address: %v", err)
	}

	return &Wallet{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Address:    address,
	}, nil
}

// Helper function to generate address from public key
func generateAddressFromPublicKey(publicKey *ecdsa.PublicKey) (string, error) {
	// Convert public key to bytes
	pubBytes := elliptic.Marshal(publicKey.Curve, publicKey.X, publicKey.Y)

	// Hash the public key
	hash := sha256.Sum256(pubBytes)

	// Create address with "AdNe" prefix
	address := "AdNe" + hex.EncodeToString(hash[:])[:60]

	return address, nil
}
