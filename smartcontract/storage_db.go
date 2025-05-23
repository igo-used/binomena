package smartcontract

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/igo-used/binomena/database"
	"gorm.io/gorm"
)

// ContractStorageDB represents database-backed contract storage
type ContractStorageDB struct {
}

// NewContractStorageWithDB creates a new database-backed contract storage
func NewContractStorageWithDB() (*ContractStorageDB, error) {
	return &ContractStorageDB{}, nil
}

// SaveContract saves a contract to the database
func (cs *ContractStorageDB) SaveContract(contract *Contract) error {
	dbContract := database.Contract{
		ContractID:     contract.ID,
		Owner:          contract.Owner,
		Name:           contract.Name,
		Code:           contract.Code,
		DeployedAt:     contract.DeployedAt.Unix(),
		LastExecuted:   contract.LastExecuted.Unix(),
		ExecutionCount: contract.ExecutionCount,
		TotalGasUsed:   contract.TotalGasUsed,
	}

	// Check if contract already exists
	var existing database.Contract
	result := database.DB.Where("contract_id = ?", contract.ID).First(&existing)
	if result.Error == gorm.ErrRecordNotFound {
		// Create new contract
		if err := database.DB.Create(&dbContract).Error; err != nil {
			return fmt.Errorf("failed to save contract: %v", err)
		}
	} else if result.Error != nil {
		return fmt.Errorf("failed to check existing contract: %v", result.Error)
	} else {
		// Update existing contract
		if err := database.DB.Model(&existing).Updates(dbContract).Error; err != nil {
			return fmt.Errorf("failed to update contract: %v", err)
		}
	}

	return nil
}

// LoadContract loads a contract from the database
func (cs *ContractStorageDB) LoadContract(contractID string) (*Contract, error) {
	var dbContract database.Contract
	result := database.DB.Where("contract_id = ?", contractID).First(&dbContract)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("contract not found")
	}
	if result.Error != nil {
		return nil, fmt.Errorf("failed to load contract: %v", result.Error)
	}

	contract := &Contract{
		ID:             dbContract.ContractID,
		Owner:          dbContract.Owner,
		Name:           dbContract.Name,
		Code:           dbContract.Code,
		DeployedAt:     time.Unix(dbContract.DeployedAt, 0),
		LastExecuted:   time.Unix(dbContract.LastExecuted, 0),
		ExecutionCount: dbContract.ExecutionCount,
		TotalGasUsed:   dbContract.TotalGasUsed,
		AverageGasUsed: 0,
	}

	if contract.ExecutionCount > 0 {
		contract.AverageGasUsed = contract.TotalGasUsed / float64(contract.ExecutionCount)
	}

	return contract, nil
}

// LoadAllContracts loads all contracts from the database
func (cs *ContractStorageDB) LoadAllContracts() ([]*Contract, error) {
	var dbContracts []database.Contract
	if err := database.DB.Find(&dbContracts).Error; err != nil {
		return nil, fmt.Errorf("failed to load contracts: %v", err)
	}

	contracts := make([]*Contract, len(dbContracts))
	for i, dbContract := range dbContracts {
		contract := &Contract{
			ID:             dbContract.ContractID,
			Owner:          dbContract.Owner,
			Name:           dbContract.Name,
			Code:           dbContract.Code,
			DeployedAt:     time.Unix(dbContract.DeployedAt, 0),
			LastExecuted:   time.Unix(dbContract.LastExecuted, 0),
			ExecutionCount: dbContract.ExecutionCount,
			TotalGasUsed:   dbContract.TotalGasUsed,
			AverageGasUsed: 0,
		}

		if contract.ExecutionCount > 0 {
			contract.AverageGasUsed = contract.TotalGasUsed / float64(contract.ExecutionCount)
		}

		contracts[i] = contract
	}

	return contracts, nil
}

// ContractStateDB represents database-backed contract state
type ContractStateDB struct {
}

// NewContractStateWithDB creates a new database-backed contract state
func NewContractStateWithDB() (*ContractStateDB, error) {
	return &ContractStateDB{}, nil
}

// GetState gets a state value for a contract
func (cs *ContractStateDB) GetState(contractID, key string) (interface{}, error) {
	var state database.SystemState
	stateKey := fmt.Sprintf("contract_%s_%s", contractID, key)

	result := database.DB.Where("key = ?", stateKey).First(&state)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("state not found")
	}
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get state: %v", result.Error)
	}

	// Try to parse as JSON first
	var value interface{}
	if err := json.Unmarshal([]byte(state.Value), &value); err != nil {
		// If JSON parsing fails, return as string
		return state.Value, nil
	}

	return value, nil
}

// SetState sets a state value for a contract
func (cs *ContractStateDB) SetState(contractID, key string, value interface{}) error {
	stateKey := fmt.Sprintf("contract_%s_%s", contractID, key)

	// Serialize value to JSON
	valueBytes, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to serialize state value: %v", err)
	}

	// Check if state already exists
	var existing database.SystemState
	result := database.DB.Where("key = ?", stateKey).First(&existing)
	if result.Error == gorm.ErrRecordNotFound {
		// Create new state
		newState := database.SystemState{
			Key:         stateKey,
			Value:       string(valueBytes),
			LastUpdated: time.Now().Unix(),
		}
		if err := database.DB.Create(&newState).Error; err != nil {
			return fmt.Errorf("failed to create state: %v", err)
		}
	} else if result.Error != nil {
		return fmt.Errorf("failed to check existing state: %v", result.Error)
	} else {
		// Update existing state
		existing.Value = string(valueBytes)
		existing.LastUpdated = time.Now().Unix()
		if err := database.DB.Save(&existing).Error; err != nil {
			return fmt.Errorf("failed to update state: %v", err)
		}
	}

	return nil
}

// DeleteState deletes a state value for a contract
func (cs *ContractStateDB) DeleteState(contractID, key string) error {
	stateKey := fmt.Sprintf("contract_%s_%s", contractID, key)

	if err := database.DB.Where("key = ?", stateKey).Delete(&database.SystemState{}).Error; err != nil {
		return fmt.Errorf("failed to delete state: %v", err)
	}

	return nil
}
