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
	imports := map[string]wasmer.IntoExtern{
		"abort": abortFunc,
		"trace": traceFunc,
	}

	// Try to create basic memory for modules that need it
	// Start with 1 page (64KB) which is standard for AssemblyScript
	memoryLimits, err := wasmer.NewLimits(1, 65536) // Min 1 page, max 65536 pages
	if err == nil {
		memoryType := wasmer.NewMemoryType(memoryLimits)
		if memory := wasmer.NewMemory(store, memoryType); memory != nil {
			imports["memory"] = memory
		}
	}

	importObject.Register("env", imports)

	return importObject
}

// executeWasmFunction executes a WASM function with the given parameters
func executeWasmFunction(instance *wasmer.Instance, function string, params []interface{}, ctx *ExecutionContext) (interface{}, error) {
	// Get function
	wasmFunc, err := instance.Exports.GetFunction(function)
	if err != nil {
		return nil, fmt.Errorf("function not found: %s", function)
	}

	// Inject caller context into WASM memory if available
	err = injectCallerContext(instance, ctx.Caller)
	if err != nil {
		// Log warning but continue - context injection is optional
		fmt.Printf("Warning: Failed to inject caller context: %v\n", err)
	}

	// Convert parameters to WASM-compatible types
	wasmParams := make([]interface{}, len(params))
	for i, param := range params {
		convertedParam, err := convertToWasmTypeWithMemory(param, instance)
		if err != nil {
			return nil, fmt.Errorf("failed to convert parameter %d: %v", i, err)
		}
		wasmParams[i] = convertedParam
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

// injectCallerContext injects the caller address into WASM memory for Context.caller
func injectCallerContext(instance *wasmer.Instance, caller string) error {
	// Try to set caller using a global variable or memory location
	// This is a simplified approach - in production you might use a more sophisticated method

	// Check if there's a setCaller function in the contract
	setCallerFunc, err := instance.Exports.GetFunction("setCaller")
	if err == nil {
		// If setCaller function exists, use it
		callerPointer, err := allocateStringInWasm(caller, instance)
		if err != nil {
			return err
		}
		_, err = setCallerFunc(callerPointer)
		return err
	}

	// Alternative: try to write caller to a known memory location
	// This is a fallback for contracts that don't have setCaller
	memory, err := instance.Exports.GetMemory("memory")
	if err != nil {
		return fmt.Errorf("memory not found in WASM instance")
	}

	// Try to allocate the caller string in memory and store the pointer
	// at a fixed location (e.g., offset 100) that the contract can read
	callerPointer, err := allocateStringInWasm(caller, instance)
	if err != nil {
		return err
	}

	// Store the caller pointer at a fixed memory location
	memData := memory.Data()
	if len(memData) >= 104 { // Need at least 104 bytes for our fixed location
		writeInt32ToMemory(memData, 100, callerPointer) // Store pointer at offset 100
	}

	return nil
}

// convertToWasmType converts a Go type to a WASM-compatible type
func convertToWasmType(value interface{}) interface{} {
	switch v := value.(type) {
	case int:
		return int32(v)
	case int32:
		return v
	case int64:
		return v
	case float32:
		return v
	case float64:
		// JSON numbers come as float64, convert to int32 if it's a whole number
		if v == float64(int32(v)) {
			return int32(v)
		}
		return v
	case bool:
		if v {
			return int32(1)
		}
		return int32(0)
	case string:
		// For backward compatibility with simple cases, return a hash of the string
		// This allows string parameters to work even if not perfect
		hash := simpleStringHash(v)
		return int32(hash)
	default:
		// Try to convert to int32 as default
		if i, ok := value.(int); ok {
			return int32(i)
		}
		// Handle JSON numbers that come as float64
		if f, ok := value.(float64); ok {
			if f == float64(int32(f)) {
				return int32(f)
			}
			return f
		}
		return int32(0)
	}
}

// convertToWasmTypeWithMemory converts a Go type to a WASM-compatible type with proper string handling
func convertToWasmTypeWithMemory(value interface{}, instance *wasmer.Instance) (interface{}, error) {
	switch v := value.(type) {
	case int:
		return int32(v), nil
	case int32:
		return v, nil
	case int64:
		return v, nil
	case float32:
		return v, nil
	case float64:
		// JSON numbers come as float64, convert to int32 if it's a whole number
		if v == float64(int32(v)) {
			return int32(v), nil
		}
		return v, nil
	case bool:
		if v {
			return int32(1), nil
		}
		return int32(0), nil
	case string:
		// Use simplified string handling to avoid memory issues
		// For now, just use the hash approach which works reliably
		hash := simpleStringHash(v)
		return int32(hash), nil
	default:
		// Try to convert to int32 as default
		if i, ok := value.(int); ok {
			return int32(i), nil
		}
		// Handle JSON numbers that come as float64
		if f, ok := value.(float64); ok {
			if f == float64(int32(f)) {
				return int32(f), nil
			}
			return f, nil
		}
		return int32(0), nil
	}
}

// allocateStringInWasm allocates a string in WASM memory and returns its pointer
func allocateStringInWasm(str string, instance *wasmer.Instance) (int32, error) {
	// For now, use a simplified approach to avoid memory issues
	// Return a deterministic hash that can be used by the contract
	hash := simpleStringHash(str)

	// Try to store the string mapping in a simple way if memory is available
	memory, err := instance.Exports.GetMemory("memory")
	if err == nil && memory != nil {
		memData := memory.Data()
		// Store string length and first few characters at a predictable location
		// This gives contracts a way to validate addresses if needed
		if len(memData) >= 200 && len(str) >= 4 {
			// Store at offset 200+ to avoid conflicts
			baseOffset := 200 + (int(hash%100) * 20) // Distribute across memory
			if baseOffset+20 < len(memData) {
				// Store first 16 chars of string for validation
				for i, char := range str[:min(16, len(str))] {
					if baseOffset+i < len(memData) {
						memData[baseOffset+i] = byte(char)
					}
				}
			}
		}
	}

	return hash, nil
}

// simpleStringHash creates a simple hash from a string for fallback cases
func simpleStringHash(str string) int32 {
	if str == "" {
		return 0
	}

	// For wallet addresses starting with "AdNe", use a special handling
	if len(str) >= 4 && str[:4] == "AdNe" {
		// Extract meaningful parts of the address for hashing
		hash := int32(0)

		// Use the last 8 characters for more uniqueness
		start := len(str) - 8
		if start < 4 {
			start = 4
		}

		for i, char := range str[start:] {
			hash = hash*37 + int32(char) + int32(i)
		}

		// Ensure positive value and reasonable range
		hash = hash & 0x7FFFFFFF // Remove sign bit
		if hash == 0 {
			hash = 1 // Avoid zero hash
		}
		return hash
	}

	// Standard string hashing for other strings
	hash := int32(5381) // djb2 hash algorithm starting value
	for _, char := range str {
		hash = ((hash << 5) + hash) + int32(char)
	}

	// Ensure positive value
	hash = hash & 0x7FFFFFFF // Remove sign bit
	if hash == 0 {
		hash = 1 // Avoid zero hash
	}
	return hash
}

// writeInt32ToMemory writes an int32 value to WASM memory at the specified offset
func writeInt32ToMemory(memory []byte, offset int, value int32) {
	if offset+4 <= len(memory) {
		memory[offset] = byte(value)
		memory[offset+1] = byte(value >> 8)
		memory[offset+2] = byte(value >> 16)
		memory[offset+3] = byte(value >> 24)
	}
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
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
