package smartcontract

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// ContractStorage handles persistence of smart contracts
type ContractStorage struct {
	storagePath string
	mu          sync.RWMutex
}

// NewContractStorage creates a new contract storage
func NewContractStorage(storagePath string) (*ContractStorage, error) {
	// Create storage directory if it doesn't exist
	if err := os.MkdirAll(storagePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %v", err)
	}
	
	return &ContractStorage{
		storagePath: storagePath,
	}, nil
}

// SaveContract saves a contract to storage
func (s *ContractStorage) SaveContract(contract *Contract) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Marshal contract to JSON
	data, err := json.Marshal(contract)
	if err != nil {
		return fmt.Errorf("failed to marshal contract: %v", err)
	}
	
	// Save to file
	filePath := filepath.Join(s.storagePath, contract.ID+".json")
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write contract file: %v", err)
	}
	
	return nil
}

// LoadContract loads a contract from storage
func (s *ContractStorage) LoadContract(contractID string) (*Contract, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	// Load from file
	filePath := filepath.Join(s.storagePath, contractID+".json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read contract file: %v", err)
	}
	
	// Unmarshal contract from JSON
	var contract Contract
	if err := json.Unmarshal(data, &contract); err != nil {
		return nil, fmt.Errorf("failed to unmarshal contract: %v", err)
	}
	
	return &contract, nil
}

// LoadAllContracts loads all contracts from storage
func (s *ContractStorage) LoadAllContracts() ([]*Contract, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	// Get all contract files
	pattern := filepath.Join(s.storagePath, "*.json")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to list contract files: %v", err)
	}
	
	// Load each contract
	contracts := make([]*Contract, 0, len(files))
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			continue // Skip files that can't be read
		}
		
		var contract Contract
		if err := json.Unmarshal(data, &contract); err != nil {
			continue // Skip files that can't be unmarshaled
		}
		
		contracts = append(contracts, &contract)
	}
	
	return contracts, nil
}

// DeleteContract deletes a contract from storage
func (s *ContractStorage) DeleteContract(contractID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Delete file
	filePath := filepath.Join(s.storagePath, contractID+".json")
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete contract file: %v", err)
	}
	
	return nil
}

// ContractExists checks if a contract exists in storage
func (s *ContractStorage) ContractExists(contractID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	// Check if file exists
	filePath := filepath.Join(s.storagePath, contractID+".json")
	_, err := os.Stat(filePath)
	return err == nil
}
