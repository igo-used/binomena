package smartcontract

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// ContractState manages the state of smart contracts
type ContractState struct {
	storagePath string
	states      map[string]map[string]interface{}
	mu          sync.RWMutex
}

// NewContractState creates a new contract state manager
func NewContractState(storagePath string) (*ContractState, error) {
	// Create storage directory if it doesn't exist
	statePath := filepath.Join(storagePath, "state")
	if err := os.MkdirAll(statePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create state directory: %v", err)
	}
	
	return &ContractState{
		storagePath: statePath,
		states:      make(map[string]map[string]interface{}),
	}, nil
}

// GetState gets a state value for a contract
func (s *ContractState) GetState(contractID string, key string) (interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	// Load state if not loaded
	if _, exists := s.states[contractID]; !exists {
		if err := s.loadState(contractID); err != nil {
			return nil, err
		}
	}
	
	// Get value
	value, exists := s.states[contractID][key]
	if !exists {
		return nil, nil // Key not found, return nil
	}
	
	return value, nil
}

// SetState sets a state value for a contract
func (s *ContractState) SetState(contractID string, key string, value interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Load state if not loaded
	if _, exists := s.states[contractID]; !exists {
		if err := s.loadState(contractID); err != nil {
			// If state doesn't exist, create it
			s.states[contractID] = make(map[string]interface{})
		}
	}
	
	// Set value
	s.states[contractID][key] = value
	
	// Save state
	return s.saveState(contractID)
}

// DeleteState deletes a state value for a contract
func (s *ContractState) DeleteState(contractID string, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Load state if not loaded
	if _, exists := s.states[contractID]; !exists {
		if err := s.loadState(contractID); err != nil {
			return nil // State doesn't exist, nothing to delete
		}
	}
	
	// Delete value
	delete(s.states[contractID], key)
	
	// Save state
	return s.saveState(contractID)
}

// ClearState clears all state for a contract
func (s *ContractState) ClearState(contractID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Delete state
	delete(s.states, contractID)
	
	// Delete state file
	filePath := filepath.Join(s.storagePath, contractID+".json")
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete state file: %v", err)
	}
	
	return nil
}

// loadState loads the state for a contract
func (s *ContractState) loadState(contractID string) error {
	// Load from file
	filePath := filepath.Join(s.storagePath, contractID+".json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// State doesn't exist, create empty state
			s.states[contractID] = make(map[string]interface{})
			return nil
		}
		return fmt.Errorf("failed to read state file: %v", err)
	}
	
	// Unmarshal state from JSON
	var state map[string]interface{}
	if err := json.Unmarshal(data, &state); err != nil {
		return fmt.Errorf("failed to unmarshal state: %v", err)
	}
	
	// Store state
	s.states[contractID] = state
	
	return nil
}

// saveState saves the state for a contract
func (s *ContractState) saveState(contractID string) error {
	// Get state
	state, exists := s.states[contractID]
	if !exists {
		return fmt.Errorf("state not found for contract: %s", contractID)
	}
	
	// Marshal state to JSON
	data, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("failed to marshal state: %v", err)
	}
	
	// Save to file
	filePath := filepath.Join(s.storagePath, contractID+".json")
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write state file: %v", err)
	}
	
	return nil
}
