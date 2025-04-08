package token

import (
	"fmt"
	"sync"
)

// BinomToken represents the native Binom (BNM) token
type BinomToken struct {
	maxSupply         float64
	circulatingSupply float64
	balances          map[string]float64
	mu                sync.RWMutex
}

// NewBinomToken creates a new Binom token with 1 billion max supply
func NewBinomToken() *BinomToken {
	maxSupply := 1000000000.0 // 1 billion
	
	token := &BinomToken{
		maxSupply:         maxSupply,
		circulatingSupply: maxSupply,
		balances:          make(map[string]float64),
	}
	
	// Allocate initial supply to treasury
	token.balances["treasury"] = maxSupply
	
	return token
}

// Transfer transfers tokens from one address to another
func (bt *BinomToken) Transfer(from, to string, amount float64) error {
	bt.mu.Lock()
	defer bt.mu.Unlock()
	
	// Check if sender has enough balance
	if bt.balances[from] < amount {
		return fmt.Errorf("insufficient balance")
	}
	
	// Transfer tokens
	bt.balances[from] -= amount
	bt.balances[to] += amount
	
	return nil
}

// GetBalance returns the balance of an address
func (bt *BinomToken) GetBalance(address string) float64 {
	bt.mu.RLock()
	defer bt.mu.RUnlock()
	
	return bt.balances[address]
}

// GetCirculatingSupply returns the circulating supply of tokens
func (bt *BinomToken) GetCirculatingSupply() float64 {
	bt.mu.RLock()
	defer bt.mu.RUnlock()
	
	return bt.circulatingSupply
}

// Burn burns tokens, reducing the circulating supply
func (bt *BinomToken) Burn(amount float64) {
	bt.mu.Lock()
	defer bt.mu.Unlock()
	
	bt.circulatingSupply -= amount
	
	fmt.Printf("Burned %.2f BNM tokens. New circulating supply: %.2f\n", amount, bt.circulatingSupply)
}

// Mint mints new tokens, increasing the circulating supply
// This is restricted to not exceed the max supply
func (bt *BinomToken) Mint(to string, amount float64) error {
	bt.mu.Lock()
	defer bt.mu.Unlock()
	
	// Check if minting would exceed max supply
	if bt.circulatingSupply+amount > bt.maxSupply {
		return fmt.Errorf("minting would exceed max supply")
	}
	
	// Mint tokens
	bt.balances[to] += amount
	bt.circulatingSupply += amount
	
	return nil
}
