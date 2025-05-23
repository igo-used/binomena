package smartcontract

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/igo-used/binomena/core"
	"github.com/igo-used/binomena/token"
	wasmer "github.com/wasmerio/wasmer-go/wasmer"
)

// Gas pricing constants
const (
	// Base fee for contract execution (in BNM)
	BaseExecutionFee = 0.001

	// Gas per instruction (in BNM)
	GasPerInstruction = 0.0000001

	// Gas limit per contract execution
	DefaultGasLimit = 10000000

	// Deployment fee per byte (in BNM)
	DeploymentFeePerByte = 0.0000005

	// Minimum deployment fee (in BNM)
	MinimumDeploymentFee = 0.1
)

// SecurityLevel defines the security level for contract execution
type SecurityLevel int

const (
	// LowSecurity allows all WASM instructions
	LowSecurity SecurityLevel = iota

	// MediumSecurity restricts some potentially dangerous operations
	MediumSecurity

	// HighSecurity applies strict limitations on contract execution
	HighSecurity
)

// WasmVM represents the WebAssembly virtual machine for smart contract execution
type WasmVM struct {
	contracts     map[string]*Contract
	instances     map[string]*wasmer.Instance
	store         *wasmer.Store
	engine        *wasmer.Engine
	securityLevel SecurityLevel
	mu            sync.RWMutex
	binomToken    *token.BinomToken
	blockchain    *core.Blockchain
}

// Contract represents a smart contract
type Contract struct {
	ID             string    `json:"id"`
	Owner          string    `json:"owner"`
	Code           []byte    `json:"code"`
	Name           string    `json:"name"`
	DeployedAt     time.Time `json:"deployedAt"`
	LastExecuted   time.Time `json:"lastExecuted"`
	ExecutionCount uint64    `json:"executionCount"`
	TotalGasUsed   float64   `json:"totalGasUsed"`
	AverageGasUsed float64   `json:"averageGasUsed"`
}

// ExecutionResult represents the result of a contract execution
type ExecutionResult struct {
	Success       bool          `json:"success"`
	GasUsed       float64       `json:"gasUsed"`
	ExecutionFee  float64       `json:"executionFee"`
	ReturnValue   interface{}   `json:"returnValue"`
	Error         string        `json:"error,omitempty"`
	ExecutionTime time.Duration `json:"executionTime"`
}

// NewWasmVM creates a new WebAssembly virtual machine
func NewWasmVM(binomToken *token.BinomToken, blockchain *core.Blockchain) (*WasmVM, error) {
	// Create a new WASM engine
	engine := wasmer.NewEngine()
	store := wasmer.NewStore(engine)

	return &WasmVM{
		contracts:     make(map[string]*Contract),
		instances:     make(map[string]*wasmer.Instance),
		store:         store,
		engine:        engine,
		securityLevel: HighSecurity, // Default to high security
		binomToken:    binomToken,
		blockchain:    blockchain,
	}, nil
}

// AddContract adds an existing contract to the VM
func (vm *WasmVM) AddContract(contract *Contract) error {
	vm.mu.Lock()
	defer vm.mu.Unlock()

	// Store contract
	vm.contracts[contract.ID] = contract

	return nil
}

// DeployContract deploys a new smart contract
func (vm *WasmVM) DeployContract(owner string, name string, code []byte, fee float64) (string, error) {
	vm.mu.Lock()
	defer vm.mu.Unlock()

	// Calculate required fee
	requiredFee := calculateDeploymentFee(code)

	// Check if fee is sufficient
	if fee < requiredFee {
		return "", fmt.Errorf("insufficient fee: required %.6f BNM, got %.6f BNM", requiredFee, fee)
	}

	// Validate WASM code
	if err := validateWasmCode(code, vm.securityLevel); err != nil {
		return "", fmt.Errorf("invalid WASM code: %v", err)
	}

	// Generate contract ID with "AdNe" prefix
	hash := sha256.Sum256(append(code, []byte(owner+name+time.Now().String())...))
	contractID := "AdNe" + hex.EncodeToString(hash[:])[:60]

	// Create contract
	contract := &Contract{
		ID:             contractID,
		Owner:          owner,
		Code:           code,
		Name:           name,
		DeployedAt:     time.Now(),
		ExecutionCount: 0,
		TotalGasUsed:   0,
		AverageGasUsed: 0,
	}

	// Store contract
	vm.contracts[contractID] = contract

	// Burn the fee
	vm.binomToken.Burn(fee)

	// Compile and validate the contract
	_, err := vm.compileContract(contract)
	if err != nil {
		delete(vm.contracts, contractID)
		return "", fmt.Errorf("failed to compile contract: %v", err)
	}

	return contractID, nil
}

// ExecuteContract executes a smart contract
func (vm *WasmVM) ExecuteContract(contractID string, function string, params []interface{}, caller string, fee float64) (*ExecutionResult, error) {
	vm.mu.Lock()
	defer vm.mu.Unlock()

	// Get contract
	contract, exists := vm.contracts[contractID]
	if !exists {
		return nil, fmt.Errorf("contract not found: %s", contractID)
	}

	// Check if fee is sufficient for base execution
	if fee < BaseExecutionFee {
		return nil, fmt.Errorf("insufficient fee: required at least %.6f BNM", BaseExecutionFee)
	}

	// Get or create instance
	instance, err := vm.getContractInstance(contract)
	if err != nil {
		return nil, fmt.Errorf("failed to get contract instance: %v", err)
	}

	// Prepare execution context
	ctx := &ExecutionContext{
		ContractID: contractID,
		Caller:     caller,
		GasLimit:   DefaultGasLimit,
		GasUsed:    0,
		VM:         vm,
		Blockchain: vm.blockchain,
		BinomToken: vm.binomToken,
		StartTime:  time.Now(),
	}

	// Execute contract
	startTime := time.Now()
	result, err := executeWasmFunction(instance, function, params, ctx)
	executionTime := time.Since(startTime)

	// Calculate gas used
	gasUsed := calculateGasUsed(ctx.GasUsed)
	executionFee := BaseExecutionFee + gasUsed

	// Update contract statistics
	contract.LastExecuted = time.Now()
	contract.ExecutionCount++
	contract.TotalGasUsed += gasUsed
	contract.AverageGasUsed = contract.TotalGasUsed / float64(contract.ExecutionCount)

	// Burn the fee
	vm.binomToken.Burn(fee)

	// Prepare result
	executionResult := &ExecutionResult{
		Success:       err == nil,
		GasUsed:       gasUsed,
		ExecutionFee:  executionFee,
		ReturnValue:   result,
		ExecutionTime: executionTime,
	}

	if err != nil {
		executionResult.Error = err.Error()
	}

	return executionResult, nil
}

// GetContract returns a contract by ID
func (vm *WasmVM) GetContract(contractID string) (*Contract, error) {
	vm.mu.RLock()
	defer vm.mu.RUnlock()

	contract, exists := vm.contracts[contractID]
	if !exists {
		return nil, fmt.Errorf("contract not found: %s", contractID)
	}

	return contract, nil
}

// ListContracts returns all contracts
func (vm *WasmVM) ListContracts() []*Contract {
	vm.mu.RLock()
	defer vm.mu.RUnlock()

	contracts := make([]*Contract, 0, len(vm.contracts))
	for _, contract := range vm.contracts {
		contracts = append(contracts, contract)
	}

	return contracts
}

// SetSecurityLevel sets the security level for contract execution
func (vm *WasmVM) SetSecurityLevel(level SecurityLevel) {
	vm.mu.Lock()
	defer vm.mu.Unlock()

	vm.securityLevel = level
}

// compileContract compiles a contract and returns a module
func (vm *WasmVM) compileContract(contract *Contract) (*wasmer.Module, error) {
	// Compile WASM code
	module, err := wasmer.NewModule(vm.store, contract.Code)
	if err != nil {
		return nil, fmt.Errorf("failed to compile WASM code: %v", err)
	}

	return module, nil
}

// getContractInstance gets or creates a contract instance
func (vm *WasmVM) getContractInstance(contract *Contract) (*wasmer.Instance, error) {
	// Check if instance already exists
	instance, exists := vm.instances[contract.ID]
	if exists {
		return instance, nil
	}

	// Compile contract
	module, err := vm.compileContract(contract)
	if err != nil {
		return nil, err
	}

	// Create import object with security restrictions
	importObject := createSecureImportObject(vm.store, vm.securityLevel)

	// Instantiate module
	// Change from := to = since instance and err are already declared
	var newInstance *wasmer.Instance
	newInstance, err = wasmer.NewInstance(module, importObject)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate WASM module: %v", err)
	}

	// Store instance
	vm.instances[contract.ID] = newInstance

	return newInstance, nil
}

// calculateDeploymentFee calculates the fee required to deploy a contract
func calculateDeploymentFee(code []byte) float64 {
	fee := float64(len(code)) * DeploymentFeePerByte
	if fee < MinimumDeploymentFee {
		fee = MinimumDeploymentFee
	}
	return fee
}

// calculateGasUsed calculates the gas used in BNM
func calculateGasUsed(instructions uint64) float64 {
	return float64(instructions) * GasPerInstruction
}

// validateWasmCode validates WASM code based on security level
func validateWasmCode(code []byte, level SecurityLevel) error {
	// Basic validation: check if it's a valid WASM binary
	if len(code) < 8 {
		return errors.New("invalid WASM binary: too short")
	}

	// Check WASM magic number
	if code[0] != 0x00 || code[1] != 0x61 || code[2] != 0x73 || code[3] != 0x6D {
		return errors.New("invalid WASM binary: wrong magic number")
	}

	// Check WASM version
	if code[4] != 0x01 || code[5] != 0x00 || code[6] != 0x00 || code[7] != 0x00 {
		return errors.New("invalid WASM binary: unsupported version")
	}

	// Additional security checks based on security level
	if level == HighSecurity {
		// Implement more sophisticated validation for high security
		// This would typically involve analyzing the WASM sections and instructions
		// For a production system, consider using a dedicated WASM validator
	}

	return nil
}

// createSecureImportObject creates a secure import object for WASM execution
func createSecureImportObject(store *wasmer.Store, level SecurityLevel) *wasmer.ImportObject {
	importObject := wasmer.NewImportObject()

	// Add abort function required by AssemblyScript
	abortFunc := wasmer.NewFunction(
		store,
		wasmer.NewFunctionType(
			wasmer.NewValueTypes(wasmer.I32, wasmer.I32, wasmer.I32, wasmer.I32),
			wasmer.NewValueTypes(),
		),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			// Handle abort - log the error and continue
			message := args[0].I32()
			fileName := args[1].I32()
			lineNumber := args[2].I32()
			columnNumber := args[3].I32()

			fmt.Printf("Contract abort: message=%d, file=%d, line=%d, col=%d\n",
				message, fileName, lineNumber, columnNumber)

			return []wasmer.Value{}, nil
		},
	)

	// Add trace function for debugging (optional)
	traceFunc := wasmer.NewFunction(
		store,
		wasmer.NewFunctionType(
			wasmer.NewValueTypes(wasmer.I32, wasmer.I32, wasmer.I32, wasmer.I32),
			wasmer.NewValueTypes(),
		),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			// Handle trace function
			return []wasmer.Value{}, nil
		},
	)

	// Register functions in the env namespace
	importObject.Register("env", map[string]wasmer.IntoExtern{
		"abort": abortFunc,
		"trace": traceFunc,
	})

	return importObject
}

// executeWasmFunction executes a WASM function with the given parameters
func executeWasmFunction(instance *wasmer.Instance, function string, params []interface{}, ctx *ExecutionContext) (interface{}, error) {
	// Get function
	wasmFunc, err := instance.Exports.GetFunction(function)
	if err != nil {
		return nil, fmt.Errorf("function not found: %s", function)
	}

	// Convert parameters to WASM-compatible types
	wasmParams := make([]interface{}, len(params))
	for i, param := range params {
		wasmParams[i] = convertToWasmType(param)
	}

	// Execute function with gas metering
	// In a real implementation, you would use a custom gas metering mechanism
	result, err := wasmFunc(wasmParams...)
	if err != nil {
		return nil, fmt.Errorf("execution failed: %v", err)
	}

	// Simulate gas usage based on execution time
	// In a real implementation, you would count actual WASM instructions
	ctx.GasUsed = uint64(time.Since(ctx.StartTime).Microseconds())

	return result, nil
}

// convertToWasmType converts a Go type to a WASM-compatible type
func convertToWasmType(value interface{}) interface{} {
	// Implement type conversion logic
	// This is a simplified version; a real implementation would handle more types
	return value
}

// ExecutionContext represents the context for contract execution
type ExecutionContext struct {
	ContractID string
	Caller     string
	GasLimit   uint64
	GasUsed    uint64
	VM         *WasmVM
	Blockchain *core.Blockchain
	BinomToken *token.BinomToken
	StartTime  time.Time
}
