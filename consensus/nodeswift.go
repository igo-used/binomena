package consensus

import (
	"crypto/rand"
	"math/big"
	"time"

	"github.com/igo-used/binomena/core"
)

// NodeSwift implements a custom Proof of Stake consensus mechanism
// that prioritizes security and transaction speed
type NodeSwift struct {
	// Minimum stake required to participate in validation
	minimumStake float64

	// Time window for block validation (in seconds)
	validationWindow int64

	// Reputation scores for validators
	validatorScores map[string]float64
}

// NewNodeSwift creates a new NodeSwift consensus mechanism
func NewNodeSwift() *NodeSwift {
	return &NodeSwift{
		minimumStake:     1000.0, // 1000 BNM minimum stake
		validationWindow: 5,      // 5 second validation window
		validatorScores:  make(map[string]float64),
	}
}

// ValidateBlock validates a block using the NodeSwift consensus rules
func (ns *NodeSwift) ValidateBlock(block core.Block) bool {
	// In a real implementation, this would validate:
	// 1. Block structure and hash
	// 2. Transaction signatures
	// 3. That the block was created by the selected validator
	// 4. That the block was created within the validation window

	// For this example, we'll just return true
	return true
}

// SelectValidator selects a validator for the next block
// based on stake amount and reputation score
func (ns *NodeSwift) SelectValidator(validators []string, stakes map[string]float64) string {
	if len(validators) == 0 {
		return ""
	}

	// Filter validators with minimum stake
	eligibleValidators := []string{}
	totalWeight := 0.0

	for _, validator := range validators {
		stake := stakes[validator]
		if stake >= ns.minimumStake {
			// Calculate validator weight based on stake and reputation
			reputation := ns.validatorScores[validator]
			if reputation == 0 {
				reputation = 1.0 // Default reputation
			}

			eligibleValidators = append(eligibleValidators, validator)
			totalWeight += stake * reputation
		}
	}

	if len(eligibleValidators) == 0 {
		return ""
	}

	// Select validator based on weighted probability
	// This is a simplified implementation
	randomValue, _ := rand.Int(rand.Reader, big.NewInt(1000))
	randomFloat := float64(randomValue.Int64()) / 1000.0 * totalWeight

	cumulativeWeight := 0.0
	for _, validator := range eligibleValidators {
		stake := stakes[validator]
		reputation := ns.validatorScores[validator]
		if reputation == 0 {
			reputation = 1.0
		}

		weight := stake * reputation
		cumulativeWeight += weight

		if randomFloat <= cumulativeWeight {
			return validator
		}
	}

	// Fallback to the first validator if something goes wrong
	return eligibleValidators[0]
}

// UpdateValidatorScore updates the reputation score of a validator
func (ns *NodeSwift) UpdateValidatorScore(validator string, successful bool) {
	currentScore, exists := ns.validatorScores[validator]
	if !exists {
		currentScore = 1.0
	}

	if successful {
		// Increase score for successful validation
		ns.validatorScores[validator] = currentScore * 1.01
	} else {
		// Decrease score for failed validation
		ns.validatorScores[validator] = currentScore * 0.5
	}
}

// GetValidationDeadline returns the deadline for block validation
func (ns *NodeSwift) GetValidationDeadline() time.Time {
	return time.Now().Add(time.Duration(ns.validationWindow) * time.Second)
}
