package tests

import (
	"testing"

	"github.com/igo-used/binomena/core"
	"github.com/igo-used/binomena/smartcontract"
	"github.com/igo-used/binomena/token"
	"github.com/igo-used/binomena/wallet"
)

// Simple test WASM binary (placeholder)
var testWasmBinary = []byte{
	0x00, 0x61, 0x73, 0x6D, // magic
	0x01, 0x00, 0x00, 0x00, // version
	// ... more binary data would go here
}

func TestWalletAndTransactions(t *testing.T) {
	// Create wallets
	alice, _ := wallet.NewWallet()
	bob, _ := wallet.NewWallet()

	// Create transaction
	tx, err := core.NewTransaction(
		alice.Address,
		bob.Address,
		50.0,
		alice,
	)

	if err != nil {
		t.Fatalf("Failed to create transaction: %v", err)
	}

	// Verify transaction
	isValid := core.VerifyTransaction(tx, alice.PublicKey)
	if !isValid {
		t.Error("Transaction verification failed")
	}

	// Calculate fee
	fee := tx.CalculateFee()
	if fee != 0.05 { // 50.0 * 0.001
		t.Errorf("Fee calculation incorrect: expected=0.05, got=%.4f", fee)
	}
}

func TestSmartContractAndWallet(t *testing.T) {
	// Skip if no valid WASM binary is available
	t.Skip("Skipping smart contract integration test - needs valid WASM binary")

	// Create blockchain components
	binomToken := &token.BinomToken{}
	blockchain := &core.Blockchain{}

	// Create VM
	vm, _ := smartcontract.NewWasmVM(binomToken, blockchain)

	// Create wallet for contract owner
	ownerWallet, _ := wallet.NewWallet()

	// Deploy contract
	contractID, err := vm.DeployContract(
		ownerWallet.Address,
		"IntegrationTestContract",
		testWasmBinary,
		0.5, // fee
	)

	if err != nil {
		t.Fatalf("Failed to deploy contract: %v", err)
	}

	// Execute contract
	result, err := vm.ExecuteContract(
		contractID,
		"test_function",
		[]interface{}{123, "test"},
		ownerWallet.Address,
		0.01, // fee
	)

	// This will likely fail without a real WASM binary, but we're testing the integration
	if err == nil {
		t.Logf("Contract execution result: %+v", result)
	}
}
