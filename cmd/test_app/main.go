package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/igo-used/binomena/core"
	"github.com/igo-used/binomena/smartcontract"
	"github.com/igo-used/binomena/token"
	"github.com/igo-used/binomena/wallet"
)

func main() {
	fmt.Println("Binomena Blockchain Test Application")
	fmt.Println("===================================")

	// Step 1: Create wallets
	fmt.Println("\n1. Creating wallets...")
	aliceWallet, err := wallet.NewWallet()
	if err != nil {
		log.Fatalf("Failed to create Alice's wallet: %v", err)
	}

	bobWallet, err := wallet.NewWallet()
	if err != nil {
		log.Fatalf("Failed to create Bob's wallet: %v", err)
	}

	fmt.Printf("Alice's wallet address: %s\n", aliceWallet.Address)
	fmt.Printf("Bob's wallet address: %s\n", bobWallet.Address)

	// Step 2: Create a transaction
	fmt.Println("\n2. Creating a transaction...")
	tx, err := core.NewTransaction(
		aliceWallet.Address,
		bobWallet.Address,
		100.0,
		aliceWallet,
	)

	if err != nil {
		log.Fatalf("Failed to create transaction: %v", err)
	}

	fmt.Printf("Transaction created: %s\n", tx.ID)
	fmt.Printf("  From: %s\n", tx.From)
	fmt.Printf("  To: %s\n", tx.To)
	fmt.Printf("  Amount: %.2f\n", tx.Amount)
	fmt.Printf("  Fee: %.4f\n", tx.CalculateFee())

	// Verify transaction
	isValid := core.VerifyTransaction(tx, aliceWallet.PublicKey)
	fmt.Printf("Transaction signature valid: %v\n", isValid)

	// Step 3: Initialize blockchain components
	fmt.Println("\n3. Initializing blockchain components...")
	binomToken := &token.BinomToken{}
	blockchain := &core.Blockchain{}

	// Step 4: Create WebAssembly VM
	fmt.Println("\n4. Creating WebAssembly VM...")
	vm, err := smartcontract.NewWasmVM(binomToken, blockchain)
	if err != nil {
		log.Fatalf("Failed to create WasmVM: %v", err)
	}

	// Step 5: Deploy a smart contract (if WASM file is available)
	fmt.Println("\n5. Deploying a smart contract...")

	// Check if a WASM file is provided as argument
	var wasmCode []byte
	if len(os.Args) > 1 {
		wasmFilePath := os.Args[1]
		var err error
		wasmCode, err = ioutil.ReadFile(wasmFilePath)
		if err != nil {
			log.Printf("Warning: Failed to read WASM file: %v", err)
			fmt.Println("   Skipping contract deployment (no valid WASM file)")
		}
	} else {
		fmt.Println("   Skipping contract deployment (no WASM file provided)")
		fmt.Println("   Usage: test_app <path-to-wasm-file>")
	}

	if len(wasmCode) > 0 {
		contractID, err := vm.DeployContract(
			aliceWallet.Address,
			"TestContract",
			wasmCode,
			0.5, // fee
		)

		if err != nil {
			log.Printf("Warning: Failed to deploy contract: %v", err)
		} else {
			fmt.Printf("Contract deployed: %s\n", contractID)

			// Step 6: Execute the smart contract
			fmt.Println("\n6. Executing the smart contract...")
			result, err := vm.ExecuteContract(
				contractID,
				"add",                // function name
				[]interface{}{40, 2}, // parameters
				aliceWallet.Address,
				0.01, // fee
			)

			if err != nil {
				log.Printf("Warning: Failed to execute contract: %v", err)
			} else {
				fmt.Printf("Execution result: %+v\n", result)
			}
		}
	}

	fmt.Println("\nTest application completed!")
}
