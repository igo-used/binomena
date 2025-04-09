package smartcontract

import (
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/igo-used/binomena/token"
	"github.com/igo-used/binomena/wallet"
)

// ContractAPI handles API endpoints for smart contracts
type ContractAPI struct {
	vm      *WasmVM
	storage *ContractStorage
	state   *ContractState
	token   *token.BinomToken
}

// NewContractAPI creates a new contract API
func NewContractAPI(vm *WasmVM, storage *ContractStorage, state *ContractState, token *token.BinomToken) *ContractAPI {
	return &ContractAPI{
		vm:      vm,
		storage: storage,
		state:   state,
		token:   token,
	}
}

// RegisterRoutes registers API routes for smart contracts
func (api *ContractAPI) RegisterRoutes(router *gin.Engine) {
	contracts := router.Group("/contracts")
	{
		// Deploy a new contract
		contracts.POST("/deploy", api.DeployContract)

		// Execute a contract
		contracts.POST("/:id/execute", api.ExecuteContract)

		// Get contract details
		contracts.GET("/:id", api.GetContract)

		// List all contracts
		contracts.GET("", api.ListContracts)

		// Get contract state
		contracts.GET("/:id/state/:key", api.GetContractState)

		// Set contract state (admin only)
		contracts.POST("/:id/state", api.SetContractState)
	}
}

// DeployContract handles contract deployment
func (api *ContractAPI) DeployContract(c *gin.Context) {
	var request struct {
		Owner      string  `json:"owner" binding:"required"`
		Name       string  `json:"name" binding:"required"`
		Code       string  `json:"code" binding:"required"` // Base64 encoded WASM
		Fee        float64 `json:"fee" binding:"required"`
		PrivateKey string  `json:"privateKey" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Decode WASM code
	code, err := base64.StdEncoding.DecodeString(request.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid WASM code: " + err.Error()})
		return
	}

	// Verify owner's wallet
	ownerWallet, err := wallet.ImportPrivateKey(request.PrivateKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid private key"})
		return
	}

	// Check if wallet address matches
	if ownerWallet.Address != request.Owner {
		c.JSON(http.StatusBadRequest, gin.H{"error": "private key does not match owner address"})
		return
	}

	// Check if owner has enough balance
	balance := api.token.GetBalance(request.Owner)
	if balance < request.Fee {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":    "insufficient balance",
			"balance":  balance,
			"required": request.Fee,
		})
		return
	}

	// Deploy contract
	contractID, err := api.vm.DeployContract(request.Owner, request.Name, code, request.Fee)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get contract
	contract, err := api.vm.GetContract(contractID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Save contract to storage
	if err := api.storage.SaveContract(contract); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create transaction
	tx := CreateDeployTransaction(contractID, request.Owner, request.Fee, 0, "success", nil)

	c.JSON(http.StatusOK, gin.H{
		"contractId":  contractID,
		"owner":       request.Owner,
		"name":        request.Name,
		"deployedAt":  contract.DeployedAt,
		"transaction": tx,
	})
}

// ExecuteContract handles contract execution
func (api *ContractAPI) ExecuteContract(c *gin.Context) {
	contractID := c.Param("id")

	var request struct {
		Caller     string        `json:"caller" binding:"required"`
		Function   string        `json:"function" binding:"required"`
		Params     []interface{} `json:"params"`
		Fee        float64       `json:"fee" binding:"required"`
		PrivateKey string        `json:"privateKey" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify caller's wallet
	callerWallet, err := wallet.ImportPrivateKey(request.PrivateKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid private key"})
		return
	}

	// Check if wallet address matches
	if callerWallet.Address != request.Caller {
		c.JSON(http.StatusBadRequest, gin.H{"error": "private key does not match caller address"})
		return
	}

	// Check if caller has enough balance
	balance := api.token.GetBalance(request.Caller)
	if balance < request.Fee {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":    "insufficient balance",
			"balance":  balance,
			"required": request.Fee,
		})
		return
	}

	// Execute contract
	result, err := api.vm.ExecuteContract(contractID, request.Function, request.Params, request.Caller, request.Fee)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create transaction
	tx := CreateExecuteTransaction(
		contractID,
		request.Caller,
		request.Function,
		request.Params,
		request.Fee,
		result.GasUsed,
		result.ExecutionTime.Milliseconds(),
		"success",
		nil,
	)

	c.JSON(http.StatusOK, gin.H{
		"result":      result,
		"transaction": tx,
	})
}

// GetContract handles getting contract details
func (api *ContractAPI) GetContract(c *gin.Context) {
	contractID := c.Param("id")

	// Get contract
	contract, err := api.vm.GetContract(contractID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, contract)
}

// ListContracts handles listing all contracts
func (api *ContractAPI) ListContracts(c *gin.Context) {
	// List contracts
	contracts := api.vm.ListContracts()

	c.JSON(http.StatusOK, gin.H{
		"contracts": contracts,
		"count":     len(contracts),
	})
}

// GetContractState handles getting contract state
func (api *ContractAPI) GetContractState(c *gin.Context) {
	contractID := c.Param("id")
	key := c.Param("key")

	// Get state
	value, err := api.state.GetState(contractID, key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if value == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "state key not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"contractId": contractID,
		"key":        key,
		"value":      value,
	})
}

// SetContractState handles setting contract state (admin only)
func (api *ContractAPI) SetContractState(c *gin.Context) {
	contractID := c.Param("id")

	var request struct {
		Key        string      `json:"key" binding:"required"`
		Value      interface{} `json:"value"`
		Caller     string      `json:"caller" binding:"required"`
		PrivateKey string      `json:"privateKey" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify caller's wallet
	callerWallet, err := wallet.ImportPrivateKey(request.PrivateKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid private key"})
		return
	}

	// Check if wallet address matches
	if callerWallet.Address != request.Caller {
		c.JSON(http.StatusBadRequest, gin.H{"error": "private key does not match caller address"})
		return
	}

	// Get contract
	contract, err := api.vm.GetContract(contractID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Check if caller is contract owner
	if contract.Owner != request.Caller {
		c.JSON(http.StatusForbidden, gin.H{"error": "only contract owner can set state"})
		return
	}

	// Set state
	if err := api.state.SetState(contractID, request.Key, request.Value); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"contractId": contractID,
		"key":        request.Key,
		"value":      request.Value,
	})
}
