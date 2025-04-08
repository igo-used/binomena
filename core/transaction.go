package core

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/igo-used/binomena/wallet"
)

// Transaction represents a transaction in the blockchain
type Transaction struct {
	ID        string  `json:"id"`
	From      string  `json:"from"`
	To        string  `json:"to"`
	Amount    float64 `json:"amount"`
	Timestamp int64   `json:"timestamp"`
	Signature string  `json:"signature"`
}

// NewTransaction creates a new transaction
func NewTransaction(from, to string, amount float64, senderWallet *wallet.Wallet) (*Transaction, error) {
	// Validate addresses
	if from[:4] != "AdNe" || to[:4] != "AdNe" {
		return nil, fmt.Errorf("addresses must start with 'AdNe'")
	}

	// Create transaction
	tx := &Transaction{
		From:      from,
		To:        to,
		Amount:    amount,
		Timestamp: time.Now().Unix(),
	}

	// Generate transaction ID with "AdNe" prefix
	txHash := sha256.Sum256([]byte(fmt.Sprintf("%s%s%f%d", from, to, amount, tx.Timestamp)))
	tx.ID = "AdNe" + hex.EncodeToString(txHash[:])[:60]

	// Sign the transaction
	signature, err := senderWallet.Sign([]byte(tx.ID))
	if err != nil {
		return nil, err
	}
	tx.Signature = hex.EncodeToString(signature)

	return tx, nil
}

// VerifyTransaction verifies the transaction signature
func VerifyTransaction(tx *Transaction, publicKey *ecdsa.PublicKey) bool {
	// Decode signature
	signature, err := hex.DecodeString(tx.Signature)
	if err != nil {
		return false
	}

	// Verify signature
	return wallet.VerifySignature(publicKey, []byte(tx.ID), signature)
}

// CalculateFee calculates the transaction fee (0.1% of the amount)
func (tx *Transaction) CalculateFee() float64 {
	return tx.Amount * 0.001
}
