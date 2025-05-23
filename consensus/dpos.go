package consensus

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/igo-used/binomena/core"
	"github.com/igo-used/binomena/database"
	"gorm.io/gorm"
)

const (
	// DPoS Configuration
	MaxDelegates        = 21     // Maximum number of delegates
	MinDelegateStake    = 5000.0 // Minimum BNM required to become delegate
	BlockTime           = 3      // Seconds between blocks
	DelegateRewardRatio = 0.6    // 60% of fees go to delegates
	BurnRatio           = 0.3    // 30% of fees burned
	CommunityRatio      = 0.05   // 5% to community
	FounderRatio        = 0.05   // 5% to founder
)

// Delegate represents a DPoS delegate
type Delegate struct {
	ID             uint    `gorm:"primaryKey"`
	Address        string  `gorm:"size:66;uniqueIndex;not null"`
	Stake          float64 `gorm:"type:decimal(20,8);not null"`
	VotesReceived  float64 `gorm:"type:decimal(20,8);default:0"`
	IsActive       bool    `gorm:"default:true"`
	RegisteredAt   int64   `gorm:"not null"`
	LastBlockTime  int64   `gorm:"default:0"`
	BlocksProduced uint64  `gorm:"default:0"`
	TotalRewards   float64 `gorm:"type:decimal(20,8);default:0"`
	Commission     float64 `gorm:"type:decimal(5,4);default:0.1"` // 10% default commission
}

// Vote represents a vote for a delegate
type Vote struct {
	ID           uint    `gorm:"primaryKey"`
	VoterAddress string  `gorm:"size:66;not null;index"`
	DelegateID   uint    `gorm:"not null;index"`
	Amount       float64 `gorm:"type:decimal(20,8);not null"`
	Timestamp    int64   `gorm:"not null"`
}

// DPoSConsensus implements Delegated Proof of Stake
type DPoSConsensus struct {
	delegates        []Delegate
	currentProducer  int
	mu               sync.RWMutex
	lastBlockTime    int64
	founderAddress   string
	communityAddress string
}

// NewDPoSConsensus creates a new DPoS consensus mechanism
func NewDPoSConsensus(founderAddress, communityAddress string) *DPoSConsensus {
	dpos := &DPoSConsensus{
		delegates:        []Delegate{},
		currentProducer:  0,
		lastBlockTime:    time.Now().Unix(),
		founderAddress:   founderAddress,
		communityAddress: communityAddress,
	}

	// Migrate delegate tables
	if err := database.DB.AutoMigrate(&Delegate{}, &Vote{}); err != nil {
		log.Printf("Failed to migrate DPoS tables: %v", err)
	}

	// Load existing delegates
	dpos.loadDelegates()

	return dpos
}

// RegisterDelegate registers a new delegate
func (d *DPoSConsensus) RegisterDelegate(address string, stake float64) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	// Check minimum stake requirement
	if stake < MinDelegateStake {
		return fmt.Errorf("minimum stake required: %.2f BNM", MinDelegateStake)
	}

	// Check if already registered
	var existing Delegate
	result := database.DB.Where("address = ?", address).First(&existing)
	if result.Error != gorm.ErrRecordNotFound {
		return fmt.Errorf("delegate already registered")
	}

	// Check delegate limit
	var count int64
	database.DB.Model(&Delegate{}).Where("is_active = ?", true).Count(&count)
	if count >= MaxDelegates {
		return fmt.Errorf("maximum delegates reached (%d)", MaxDelegates)
	}

	// Create new delegate
	delegate := Delegate{
		Address:       address,
		Stake:         stake,
		VotesReceived: stake, // Self-vote
		IsActive:      true,
		RegisteredAt:  time.Now().Unix(),
		Commission:    0.1, // 10% default commission
	}

	if err := database.DB.Create(&delegate).Error; err != nil {
		return fmt.Errorf("failed to register delegate: %v", err)
	}

	// Add self-vote
	vote := Vote{
		VoterAddress: address,
		DelegateID:   delegate.ID,
		Amount:       stake,
		Timestamp:    time.Now().Unix(),
	}

	if err := database.DB.Create(&vote).Error; err != nil {
		log.Printf("Failed to create self-vote: %v", err)
	}

	// Reload delegates
	d.loadDelegates()

	log.Printf("Delegate registered: %s with stake %.2f BNM", address, stake)
	return nil
}

// VoteForDelegate allows voting for a delegate
func (d *DPoSConsensus) VoteForDelegate(voterAddress, delegateAddress string, amount float64) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	// Find delegate
	var delegate Delegate
	result := database.DB.Where("address = ? AND is_active = ?", delegateAddress, true).First(&delegate)
	if result.Error != nil {
		return fmt.Errorf("delegate not found or inactive")
	}

	// Check if voter already voted for this delegate
	var existingVote Vote
	result = database.DB.Where("voter_address = ? AND delegate_id = ?", voterAddress, delegate.ID).First(&existingVote)
	if result.Error == nil {
		// Update existing vote
		existingVote.Amount += amount
		existingVote.Timestamp = time.Now().Unix()
		if err := database.DB.Save(&existingVote).Error; err != nil {
			return fmt.Errorf("failed to update vote: %v", err)
		}
	} else {
		// Create new vote
		vote := Vote{
			VoterAddress: voterAddress,
			DelegateID:   delegate.ID,
			Amount:       amount,
			Timestamp:    time.Now().Unix(),
		}
		if err := database.DB.Create(&vote).Error; err != nil {
			return fmt.Errorf("failed to create vote: %v", err)
		}
	}

	// Update delegate's total votes
	delegate.VotesReceived += amount
	if err := database.DB.Save(&delegate).Error; err != nil {
		return fmt.Errorf("failed to update delegate votes: %v", err)
	}

	// Reload delegates
	d.loadDelegates()

	log.Printf("Vote recorded: %s voted %.2f BNM for delegate %s", voterAddress, amount, delegateAddress)
	return nil
}

// GetActiveProducer returns the current block producer
func (d *DPoSConsensus) GetActiveProducer() string {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if len(d.delegates) == 0 {
		return d.founderAddress // Fallback to founder if no delegates
	}

	currentTime := time.Now().Unix()
	timeSinceLastBlock := currentTime - d.lastBlockTime

	// If enough time has passed, move to next producer
	if timeSinceLastBlock >= BlockTime {
		d.mu.RUnlock()
		d.mu.Lock()
		d.currentProducer = (d.currentProducer + 1) % len(d.delegates)
		d.lastBlockTime = currentTime
		d.mu.Unlock()
		d.mu.RLock()
	}

	return d.delegates[d.currentProducer].Address
}

// DistributeFees distributes transaction fees according to DPoS rules
func (d *DPoSConsensus) DistributeFees(totalFees float64, tokenSystem interface{}) error {
	if totalFees <= 0 {
		return nil
	}

	// Calculate fee distribution
	delegateReward := totalFees * DelegateRewardRatio // 60%
	burnAmount := totalFees * BurnRatio               // 30%
	communityReward := totalFees * CommunityRatio     // 5%
	founderReward := totalFees * FounderRatio         // 5%

	// Distribute to active delegates
	if len(d.delegates) > 0 {
		rewardPerDelegate := delegateReward / float64(len(d.delegates))
		for _, delegate := range d.delegates {
			if delegate.IsActive {
				// Transfer reward to delegate
				if err := d.transferReward(delegate.Address, rewardPerDelegate, tokenSystem); err != nil {
					log.Printf("Failed to reward delegate %s: %v", delegate.Address, err)
				} else {
					// Update delegate stats
					d.updateDelegateRewards(delegate.ID, rewardPerDelegate)
				}
			}
		}
	}

	// Burn tokens
	if burner, ok := tokenSystem.(interface{ Burn(float64) }); ok {
		burner.Burn(burnAmount)
	}

	// Reward community
	if err := d.transferReward(d.communityAddress, communityReward, tokenSystem); err != nil {
		log.Printf("Failed to reward community: %v", err)
	}

	// Reward founder
	if err := d.transferReward(d.founderAddress, founderReward, tokenSystem); err != nil {
		log.Printf("Failed to reward founder: %v", err)
	}

	log.Printf("Fees distributed: %.6f to delegates, %.6f burned, %.6f to community, %.6f to founder",
		delegateReward, burnAmount, communityReward, founderReward)

	return nil
}

// GetDelegates returns all active delegates sorted by votes
func (d *DPoSConsensus) GetDelegates() []Delegate {
	d.mu.RLock()
	defer d.mu.RUnlock()

	// Return copy of delegates
	result := make([]Delegate, len(d.delegates))
	copy(result, d.delegates)
	return result
}

// loadDelegates loads delegates from database and sorts by votes
func (d *DPoSConsensus) loadDelegates() {
	var delegates []Delegate
	database.DB.Where("is_active = ?", true).Order("votes_received DESC").Find(&delegates)

	// Limit to max delegates
	if len(delegates) > MaxDelegates {
		delegates = delegates[:MaxDelegates]
	}

	d.delegates = delegates
	log.Printf("Loaded %d active delegates", len(delegates))
}

// transferReward transfers reward tokens
func (d *DPoSConsensus) transferReward(address string, amount float64, tokenSystem interface{}) error {
	if transferer, ok := tokenSystem.(interface {
		Transfer(string, string, float64) error
	}); ok {
		return transferer.Transfer("treasury", address, amount)
	}
	return fmt.Errorf("token system does not support transfers")
}

// updateDelegateRewards updates delegate reward statistics
func (d *DPoSConsensus) updateDelegateRewards(delegateID uint, reward float64) {
	var delegate Delegate
	if err := database.DB.First(&delegate, delegateID).Error; err != nil {
		return
	}

	delegate.TotalRewards += reward
	delegate.BlocksProduced++
	database.DB.Save(&delegate)
}

// ValidateBlock validates a block (satisfies core.Consensus interface)
func (d *DPoSConsensus) ValidateBlock(block core.Block) bool {
	// Basic validation - in production, add more sophisticated checks
	// Check if block producer is an active delegate
	producer := block.Validator
	for _, delegate := range d.delegates {
		if delegate.Address == producer && delegate.IsActive {
			return true
		}
	}

	// Allow founder to produce blocks if no delegates
	if len(d.delegates) == 0 && producer == d.founderAddress {
		return true
	}

	return false
}

// SelectValidator selects next validator (satisfies core.Consensus interface)
func (d *DPoSConsensus) SelectValidator(validators []string, stakes map[string]float64) string {
	return d.GetActiveProducer()
}
