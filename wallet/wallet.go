package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
)

// Wallet represents a cryptocurrency wallet
type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
	Address    string
}

// NewWallet creates a new wallet with a key pair
func NewWallet() (*Wallet, error) {
	// Generate ECDSA private key
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	// Derive public key from private key
	publicKey := &privateKey.PublicKey

	// Generate wallet address with "AdNe" prefix
	address, err := generateAddress(publicKey)
	if err != nil {
		return nil, err
	}

	return &Wallet{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Address:    address,
	}, nil
}

// generateAddress generates a wallet address from a public key with "AdNe" prefix
func generateAddress(publicKey *ecdsa.PublicKey) (string, error) {
	// Convert public key to bytes
	pubKeyBytes := elliptic.Marshal(publicKey.Curve, publicKey.X, publicKey.Y)

	// Hash the public key using SHA-256
	sha := sha256.New()
	_, err := sha.Write(pubKeyBytes)
	if err != nil {
		return "", err
	}
	hash := sha.Sum(nil)

	// Convert hash to hex string and add "AdNe" prefix
	address := "AdNe" + hex.EncodeToString(hash)[:40]
	return address, nil
}

// Sign signs data with the wallet's private key
func (w *Wallet) Sign(data []byte) ([]byte, error) {
	// Hash the data
	hash := sha256.Sum256(data)

	// Sign the hash with the private key
	r, s, err := ecdsa.Sign(rand.Reader, w.PrivateKey, hash[:])
	if err != nil {
		return nil, err
	}

	// Combine r and s into a single signature
	signature := append(r.Bytes(), s.Bytes()...)
	return signature, nil
}

// VerifySignature verifies a signature against a public key and data
func VerifySignature(publicKey *ecdsa.PublicKey, data, signature []byte) bool {
	// Hash the data
	hash := sha256.Sum256(data)

	// Split signature into r and s
	sigLen := len(signature)
	if sigLen != 64 {
		return false
	}

	r := new(big.Int).SetBytes(signature[:32])
	s := new(big.Int).SetBytes(signature[32:])

	// Verify the signature
	return ecdsa.Verify(publicKey, hash[:], r, s)
}

// ExportPrivateKey exports the private key as a hex string
func (w *Wallet) ExportPrivateKey() string {
	return fmt.Sprintf("%x", w.PrivateKey.D.Bytes())
}

// ImportPrivateKey imports a wallet from a private key hex string
func ImportPrivateKey(privateKeyHex string) (*Wallet, error) {
	// Decode hex string to bytes
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return nil, err
	}

	// Create big.Int from bytes
	privateKeyInt := new(big.Int).SetBytes(privateKeyBytes)

	// Create private key
	privateKey := new(ecdsa.PrivateKey)
	privateKey.Curve = elliptic.P256()
	privateKey.D = privateKeyInt
	privateKey.PublicKey.Curve = elliptic.P256()
	privateKey.PublicKey.X, privateKey.PublicKey.Y = privateKey.Curve.ScalarBaseMult(privateKeyInt.Bytes())

	// Generate address
	address, err := generateAddress(&privateKey.PublicKey)
	if err != nil {
		return nil, err
	}

	return &Wallet{
		PrivateKey: privateKey,
		PublicKey:  &privateKey.PublicKey,
		Address:    address,
	}, nil
}
