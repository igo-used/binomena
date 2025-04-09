package smartcontract

import (
	"fmt"
	"time"

	"github.com/igo-used/binomena/wallet"
)

// ContractTransaction represents a transaction for a smart contract
type ContractTransaction struct {
	ID            string      `json:"id"`
	Type          string      `json:"type"`
	ContractID    string      `json:"contractId"`
	Caller        string      `json:"caller"`
	Function      string      `json:"function,omitempty"`
	Params        interface{} `json:"params,omitempty"`
	Fee           float64     `json:"fee"`
	GasUsed       float64     `json:"gasUsed"`
	ExecutionTime int64       `json:"executionTime"`
	Timestamp     int64       `json:"timestamp"`
	Status        string      `json:"status"`
	Error         string      `json:"error,omitempty"`
}

// ContractTransactionType defines the type of contract transaction
type ContractTransactionType string

const (
	// DeployTransaction represents a contract deployment transaction
	DeployTransaction ContractTransactionType = "deploy"

	// ExecuteTransaction represents a contract execution transaction
	ExecuteTransaction ContractTransactionType = "execute"
)

// CreateDeployTransaction creates a transaction for contract deployment
func CreateDeployTransaction(contractID, caller string, fee, gasUsed float64, status string, err error) *ContractTransaction {
	tx := &ContractTransaction{
		ID:         generateTransactionID(),
		Type:       string(DeployTransaction),
		ContractID: contractID,
		Caller:     caller,
		Fee:        fee,
		GasUsed:    gasUsed,
		Timestamp:  time.Now().Unix(),
		Status:     status,
	}

	if err != nil {
		tx.Error = err.Error()
	}

	return tx
}

// CreateExecuteTransaction creates a transaction for contract execution
func CreateExecuteTransaction(contractID, caller, function string, params interface{}, fee, gasUsed float64, executionTime int64, status string, err error) *ContractTransaction {
	tx := &ContractTransaction{
		ID:            generateTransactionID(),
		Type:          string(ExecuteTransaction),
		ContractID:    contractID,
		Caller:        caller,
		Function:      function,
		Params:        params,
		Fee:           fee,
		GasUsed:       gasUsed,
		ExecutionTime: executionTime,
		Timestamp:     time.Now().Unix(),
		Status:        status,
	}

	if err != nil {
		tx.Error = err.Error()
	}

	return tx
}

// VerifyContractTransaction verifies a contract transaction
func VerifyContractTransaction(tx *ContractTransaction, senderWallet *wallet.Wallet) (bool, error) {
	// Verify transaction signature
	// In a real implementation, you would verify the signature

	return true, nil
}

// generateTransactionID generates a transaction ID for a contract transaction
func generateTransactionID() string {
	// Create a unique ID with "AdNe" prefix
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("AdNe%d", timestamp)
}
