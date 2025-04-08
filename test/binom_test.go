package test

import (
	"testing"

	"github.com/igo-used/binomena/token"
)

func TestBinomToken(t *testing.T) {
	// Create a new Binom token
	binomToken := token.NewBinomToken()
	
	// Check max supply
	maxSupply := 1000000000.0 // 1 billion
	circulatingSupply := binomToken.GetCirculatingSupply()
	if circulatingSupply != maxSupply {
		t.Errorf("Expected circulating supply to be %f, got %f", maxSupply, circulatingSupply)
	}
	
	// Test token transfer
	// First, allocate some tokens to alice
	err := binomToken.Transfer("treasury", "alice", 1000.0)
	if err != nil {
		t.Errorf("Failed to transfer tokens: %v", err)
	}
	
	// Check alice's balance
	aliceBalance := binomToken.GetBalance("alice")
	if aliceBalance != 1000.0 {
		t.Errorf("Expected alice's balance to be 1000.0, got %f", aliceBalance)
	}
	
	// Test token transfer from alice to bob
	err = binomToken.Transfer("alice", "bob", 500.0)
	if err != nil {
		t.Errorf("Failed to transfer tokens: %v", err)
	}
	
	// Check balances
	aliceBalance = binomToken.GetBalance("alice")
	bobBalance := binomToken.GetBalance("bob")
	if aliceBalance != 500.0 {
		t.Errorf("Expected alice's balance to be 500.0, got %f", aliceBalance)
	}
	if bobBalance != 500.0 {
		t.Errorf("Expected bob's balance to be 500.0, got %f", bobBalance)
	}
	
	// Test token burn
	initialSupply := binomToken.GetCirculatingSupply()
	burnAmount := 100.0
	binomToken.Burn(burnAmount)
	
	// Check circulating supply after burn
	newSupply := binomToken.GetCirculatingSupply()
	expectedSupply := initialSupply - burnAmount
	if newSupply != expectedSupply {
		t.Errorf("Expected circulating supply to be %f, got %f", expectedSupply, newSupply)
	}
}
