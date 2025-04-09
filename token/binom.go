package token

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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

// SaveBalances saves token balances to a JSON file
func (bt *BinomToken) SaveBalances(dataDir string) error {
	bt.mu.RLock()
	defer bt.mu.RUnlock()

	// Create balances directory if it doesn't exist
	balancesDir := filepath.Join(dataDir, "balances")
	if err := os.MkdirAll(balancesDir, 0755); err != nil {
		return fmt.Errorf("failed to create balances directory: %v", err)
	}

	// Marshal balances to JSON
	data, err := json.MarshalIndent(bt.balances, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal balances: %v", err)
	}

	// Write to file
	balancesFile := filepath.Join(balancesDir, "token_balances.json")
	if err := os.WriteFile(balancesFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write balances file: %v", err)
	}

	// Also save circulating supply
	supplyData := fmt.Sprintf("%.8f", bt.circulatingSupply)
	supplyFile := filepath.Join(balancesDir, "circulating_supply.txt")
	if err := os.WriteFile(supplyFile, []byte(supplyData), 0644); err != nil {
		return fmt.Errorf("failed to write supply file: %v", err)
	}

	fmt.Printf("Saved balances for %d addresses\n", len(bt.balances))
	return nil
}

// LoadBalances loads token balances from a JSON file
func (bt *BinomToken) LoadBalances(dataDir string) error {
	balancesDir := filepath.Join(dataDir, "balances")
	balancesFile := filepath.Join(balancesDir, "token_balances.json")

	// Check if file exists
	if _, err := os.Stat(balancesFile); os.IsNotExist(err) {
		fmt.Printf("Balances file not found, using default initialization\n")
		return nil
	}

	// Read balances file
	data, err := os.ReadFile(balancesFile)
	if err != nil {
		return fmt.Errorf("failed to read balances file: %v", err)
	}

	// Unmarshal JSON
	bt.mu.Lock()
	if err := json.Unmarshal(data, &bt.balances); err != nil {
		bt.mu.Unlock()
		return fmt.Errorf("failed to unmarshal balances: %v", err)
	}

	// Read circulating supply
	supplyFile := filepath.Join(balancesDir, "circulating_supply.txt")
	if _, err := os.Stat(supplyFile); !os.IsNotExist(err) {
		supplyData, err := os.ReadFile(supplyFile)
		if err == nil {
			fmt.Sscanf(string(supplyData), "%f", &bt.circulatingSupply)
		}
	}

	bt.mu.Unlock()
	fmt.Printf("Loaded balances for %d addresses\n", len(bt.balances))
	return nil
}
