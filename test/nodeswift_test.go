package test

import (
	"testing"

	"github.com/binomena/consensus"
)

func TestNodeSwift(t *testing.T) {
	// Create a new NodeSwift consensus mechanism
	nodeSwift := consensus.NewNodeSwift()
	
	// Test validator selection
	validators := []string{"validator1", "validator2", "validator3"}
	stakes := map[string]float64{
		"validator1": 5000.0,
		"validator2": 10000.0,
		"validator3": 500.0, // Below minimum stake
	}
	
	// Select validator
	validator := nodeSwift.SelectValidator(validators, stakes)
	
	// Check that a validator was selected
	if validator == "" {
		t.Errorf("Expected a validator to be selected, got empty string")
	}
	
	// Check that validator3 was not selected (below minimum stake)
	if validator == "validator3" {
		t.Errorf("Expected validator3 to not be selected (below minimum stake)")
	}
	
	// Test validator score updates
	nodeSwift.UpdateValidatorScore("validator1", true)
	nodeSwift.UpdateValidatorScore("validator2", false)
	
	// Select validator again after score updates
	validator = nodeSwift.SelectValidator(validators, stakes)
	
	// Check that a validator was selected
	if validator == "" {
		t.Errorf("Expected a validator to be selected, got empty string")
	}
}
