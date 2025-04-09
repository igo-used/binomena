package smartcontract

import (
	"testing"

	"github.com/igo-used/binomena/core"
	"github.com/igo-used/binomena/token"
)

// Simple WebAssembly binary that exports an "add" function
var simpleWasmBinary = []byte{
	0x00, 0x61, 0x73, 0x6D, // magic
	0x01, 0x00, 0x00, 0x00, // version
	// ... more binary data would go here
	// This is a placeholder - you would need a real compiled WASM binary
}

func TestWasmVMCreation(t *testing.T) {
	// Create mock dependencies
	mockToken := &token.BinomToken{}
	mockBlockchain := &core.Blockchain{}

	// Create VM
	vm, err := NewWasmVM(mockToken, mockBlockchain)
	if err != nil {
		t.Fatalf("Failed to create WasmVM: %v", err)
	}

	// Check if VM was initialized correctly
	if vm.contracts == nil {
		t.Error("Contracts map not initialized")
	}

	if vm.instances == nil {
		t.Error("Instances map not initialized")
	}

	if vm.store == nil {
		t.Error("WASM store not initialized")
	}

	if vm.engine == nil {
		t.Error("WASM engine not initialized")
	}
}

func TestContractDeployment(t *testing.T) {
	// Skip this test if no valid WASM binary is available
	t.Skip("Skipping contract deployment test - needs valid WASM binary")

	// Create mock dependencies
	mockToken := &token.BinomToken{}
	mockBlockchain := &core.Blockchain{}

	// Create VM
	vm, _ := NewWasmVM(mockToken, mockBlockchain)

	// Deploy contract
	contractID, err := vm.DeployContract(
		"AdNeTestOwner123",
		"TestContract",
		simpleWasmBinary,
		0.5, // fee
	)

	if err != nil {
		t.Fatalf("Failed to deploy contract: %v", err)
	}

	// Check contract ID format
	if len(contractID) < 5 || contractID[:4] != "AdNe" {
		t.Errorf("Invalid contract ID format: %s", contractID)
	}

	// Check if contract was stored
	contract, err := vm.GetContract(contractID)
	if err != nil {
		t.Fatalf("Failed to get deployed contract: %v", err)
	}

	if contract.Name != "TestContract" {
		t.Errorf("Contract name mismatch: expected=%s, got=%s",
			"TestContract", contract.Name)
	}

	if contract.Owner != "AdNeTestOwner123" {
		t.Errorf("Contract owner mismatch: expected=%s, got=%s",
			"AdNeTestOwner123", contract.Owner)
	}
}
