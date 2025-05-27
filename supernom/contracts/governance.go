package contracts

import (
	"fmt"
	"time"
)

// GovernanceContract manages delegate voting and system governance for SuperNom
type GovernanceContract struct {
	// Voting proposals
	Proposals map[string]*Proposal `json:"proposals"`

	// Delegate management
	AuthorizedDelegates map[string]*Delegate `json:"authorizedDelegates"`

	// Emergency controls
	EmergencySettings *EmergencySettings `json:"emergencySettings"`

	// Contract metadata
	ContractID string `json:"contractId"`
	Owner      string `json:"owner"`
	CreatedAt  int64  `json:"createdAt"`
}

// Proposal represents a governance proposal for voting
type Proposal struct {
	ProposalID    string                 `json:"proposalId"`
	ProposalType  string                 `json:"proposalType"` // "blacklist", "whitelist", "settings", "emergency"
	Title         string                 `json:"title"`
	Description   string                 `json:"description"`
	TargetAddress string                 `json:"targetAddress,omitempty"`
	ProposedBy    string                 `json:"proposedBy"`
	CreatedAt     int64                  `json:"createdAt"`
	ExpiresAt     int64                  `json:"expiresAt"`
	VotesFor      map[string]bool        `json:"votesFor"` // delegate address -> vote
	VotesAgainst  map[string]bool        `json:"votesAgainst"`
	Status        string                 `json:"status"`        // "active", "passed", "rejected", "executed"
	RequiredVotes int                    `json:"requiredVotes"` // minimum votes needed
	ExecutedAt    int64                  `json:"executedAt"`
	ProposalData  map[string]interface{} `json:"proposalData"` // additional data for the proposal
}

// Delegate represents an authorized delegate in the governance system
type Delegate struct {
	Address       string  `json:"address"`
	Name          string  `json:"name"`
	StakeAmount   float64 `json:"stakeAmount"`
	VotingPower   int     `json:"votingPower"`
	LastActivity  int64   `json:"lastActivity"`
	Active        bool    `json:"active"`
	VotingHistory int     `json:"votingHistory"` // number of votes cast
}

// EmergencySettings contains emergency control parameters
type EmergencySettings struct {
	EmergencyActive      bool     `json:"emergencyActive"`
	EmergencyActivatedBy string   `json:"emergencyActivatedBy"`
	EmergencyActivatedAt int64    `json:"emergencyActivatedAt"`
	EmergencyReason      string   `json:"emergencyReason"`
	MinEmergencyVotes    int      `json:"minEmergencyVotes"` // minimum votes to trigger emergency
	EmergencyDuration    int64    `json:"emergencyDuration"` // maximum emergency duration
	AutoResolveEmergency bool     `json:"autoResolveEmergency"`
	EmergencyContacts    []string `json:"emergencyContacts"`
}

// NewGovernanceContract creates a new governance contract
func NewGovernanceContract(owner string, initialDelegates []string) *GovernanceContract {
	contract := &GovernanceContract{
		Proposals:           make(map[string]*Proposal),
		AuthorizedDelegates: make(map[string]*Delegate),
		ContractID:          generateContractID(owner),
		Owner:               owner,
		CreatedAt:           time.Now().Unix(),
		EmergencySettings: &EmergencySettings{
			EmergencyActive:      false,
			MinEmergencyVotes:    3,         // 3 out of 21 delegates can trigger emergency
			EmergencyDuration:    86400 * 7, // 7 days maximum
			AutoResolveEmergency: true,
		},
	}

	// Initialize delegates
	for i, address := range initialDelegates {
		contract.AuthorizedDelegates[address] = &Delegate{
			Address:      address,
			Name:         fmt.Sprintf("Delegate-%d", i+1),
			VotingPower:  1,
			Active:       true,
			LastActivity: time.Now().Unix(),
		}
	}

	return contract
}

// CreateProposal creates a new governance proposal
func (g *GovernanceContract) CreateProposal(proposalType, title, description, targetAddress, proposedBy string, additionalData map[string]interface{}) (*Proposal, error) {
	// Validate proposer is authorized delegate
	delegate, exists := g.AuthorizedDelegates[proposedBy]
	if !exists || !delegate.Active {
		return nil, fmt.Errorf("unauthorized proposer or inactive delegate")
	}

	// Generate proposal ID
	proposalID := generateSessionID(proposedBy, time.Now().Unix())

	// Determine required votes based on proposal type
	requiredVotes := g.calculateRequiredVotes(proposalType)

	proposal := &Proposal{
		ProposalID:    proposalID,
		ProposalType:  proposalType,
		Title:         title,
		Description:   description,
		TargetAddress: targetAddress,
		ProposedBy:    proposedBy,
		CreatedAt:     time.Now().Unix(),
		ExpiresAt:     time.Now().Unix() + 86400*3, // 3 days voting period
		VotesFor:      make(map[string]bool),
		VotesAgainst:  make(map[string]bool),
		Status:        "active",
		RequiredVotes: requiredVotes,
		ProposalData:  additionalData,
	}

	g.Proposals[proposalID] = proposal
	return proposal, nil
}

// CastVote allows delegates to vote on proposals
func (g *GovernanceContract) CastVote(proposalID, delegateAddress string, vote bool) error {
	// Validate delegate
	delegate, exists := g.AuthorizedDelegates[delegateAddress]
	if !exists || !delegate.Active {
		return fmt.Errorf("unauthorized delegate or inactive")
	}

	// Get proposal
	proposal, exists := g.Proposals[proposalID]
	if !exists {
		return fmt.Errorf("proposal not found")
	}

	// Check if proposal is still active
	if proposal.Status != "active" || proposal.ExpiresAt < time.Now().Unix() {
		return fmt.Errorf("proposal voting period has ended")
	}

	// Check if delegate already voted
	if _, votedFor := proposal.VotesFor[delegateAddress]; votedFor {
		return fmt.Errorf("delegate already voted for this proposal")
	}
	if _, votedAgainst := proposal.VotesAgainst[delegateAddress]; votedAgainst {
		return fmt.Errorf("delegate already voted against this proposal")
	}

	// Record vote
	if vote {
		proposal.VotesFor[delegateAddress] = true
	} else {
		proposal.VotesAgainst[delegateAddress] = true
	}

	// Update delegate activity
	delegate.LastActivity = time.Now().Unix()
	delegate.VotingHistory++

	// Check if proposal reached required votes
	g.checkProposalStatus(proposal)

	return nil
}

// ExecuteProposal executes a passed proposal
func (g *GovernanceContract) ExecuteProposal(proposalID string, vpnContract *VPNAccessContract) error {
	proposal, exists := g.Proposals[proposalID]
	if !exists {
		return fmt.Errorf("proposal not found")
	}

	if proposal.Status != "passed" {
		return fmt.Errorf("proposal has not passed voting")
	}

	// Execute based on proposal type
	switch proposal.ProposalType {
	case "blacklist":
		return g.executeBlacklistProposal(proposal, vpnContract)
	case "whitelist":
		return g.executeWhitelistProposal(proposal, vpnContract)
	case "settings":
		return g.executeSettingsProposal(proposal, vpnContract)
	case "emergency":
		return g.executeEmergencyProposal(proposal)
	default:
		return fmt.Errorf("unknown proposal type")
	}
}

// TriggerEmergency allows delegates to trigger emergency shutdown
func (g *GovernanceContract) TriggerEmergency(reason string, triggeredBy string) error {
	// Validate delegate
	delegate, exists := g.AuthorizedDelegates[triggeredBy]
	if !exists || !delegate.Active {
		return fmt.Errorf("unauthorized delegate")
	}

	// Check if emergency already active
	if g.EmergencySettings.EmergencyActive {
		return fmt.Errorf("emergency already active")
	}

	// For now, allow any delegate to trigger emergency (can be enhanced with voting)
	g.EmergencySettings.EmergencyActive = true
	g.EmergencySettings.EmergencyActivatedBy = triggeredBy
	g.EmergencySettings.EmergencyActivatedAt = time.Now().Unix()
	g.EmergencySettings.EmergencyReason = reason

	return nil
}

// ResolveEmergency resolves an active emergency
func (g *GovernanceContract) ResolveEmergency(resolvedBy string) error {
	// Validate delegate
	delegate, exists := g.AuthorizedDelegates[resolvedBy]
	if !exists || !delegate.Active {
		return fmt.Errorf("unauthorized delegate")
	}

	if !g.EmergencySettings.EmergencyActive {
		return fmt.Errorf("no active emergency")
	}

	g.EmergencySettings.EmergencyActive = false
	return nil
}

// Helper functions

func (g *GovernanceContract) calculateRequiredVotes(proposalType string) int {
	totalDelegates := len(g.AuthorizedDelegates)

	switch proposalType {
	case "blacklist":
		return (totalDelegates * 60) / 100 // 60% majority for blacklisting
	case "whitelist":
		return (totalDelegates * 51) / 100 // 51% majority for whitelisting
	case "settings":
		return (totalDelegates * 67) / 100 // 67% majority for settings changes
	case "emergency":
		return (totalDelegates * 75) / 100 // 75% majority for emergency actions
	default:
		return (totalDelegates * 51) / 100 // 51% majority for other proposals
	}
}

func (g *GovernanceContract) checkProposalStatus(proposal *Proposal) {
	votesFor := len(proposal.VotesFor)
	votesAgainst := len(proposal.VotesAgainst)
	totalVotes := votesFor + votesAgainst

	// Check if proposal has passed
	if votesFor >= proposal.RequiredVotes {
		proposal.Status = "passed"
		return
	}

	// Check if proposal cannot pass (too many against votes)
	remainingVotes := len(g.AuthorizedDelegates) - totalVotes
	if votesFor+remainingVotes < proposal.RequiredVotes {
		proposal.Status = "rejected"
		return
	}

	// Check if voting period expired
	if proposal.ExpiresAt < time.Now().Unix() {
		if votesFor >= proposal.RequiredVotes {
			proposal.Status = "passed"
		} else {
			proposal.Status = "rejected"
		}
	}
}

func (g *GovernanceContract) executeBlacklistProposal(proposal *Proposal, vpnContract *VPNAccessContract) error {
	if proposal.TargetAddress == "" {
		return fmt.Errorf("no target address specified for blacklist proposal")
	}

	// Add to blacklist in VPN contract
	blacklistEntry := &BlacklistEntry{
		Address:   proposal.TargetAddress,
		Reason:    proposal.Description,
		BannedAt:  time.Now().Unix(),
		BannedBy:  proposal.ProposedBy,
		Permanent: false,
		ExpiresAt: 0, // Indefinite unless specified
	}

	// Check if duration is specified in proposal data
	if duration, exists := proposal.ProposalData["duration"]; exists {
		if durationInt, ok := duration.(int64); ok {
			blacklistEntry.ExpiresAt = time.Now().Unix() + durationInt
		}
	}

	// Check if permanent ban
	if permanent, exists := proposal.ProposalData["permanent"]; exists {
		if permanentBool, ok := permanent.(bool); ok {
			blacklistEntry.Permanent = permanentBool
		}
	}

	vpnContract.Blacklist[proposal.TargetAddress] = blacklistEntry
	proposal.Status = "executed"
	proposal.ExecutedAt = time.Now().Unix()

	return nil
}

func (g *GovernanceContract) executeWhitelistProposal(proposal *Proposal, vpnContract *VPNAccessContract) error {
	if proposal.TargetAddress == "" {
		return fmt.Errorf("no target address specified for whitelist proposal")
	}

	// Remove from blacklist
	delete(vpnContract.Blacklist, proposal.TargetAddress)

	// Optionally boost reputation
	if reputation, exists := vpnContract.Reputation[proposal.TargetAddress]; exists {
		reputation.TrustScore += 100 // Boost trust score
		if reputation.TrustScore > 1000 {
			reputation.TrustScore = 1000
		}
	}

	proposal.Status = "executed"
	proposal.ExecutedAt = time.Now().Unix()

	return nil
}

func (g *GovernanceContract) executeSettingsProposal(proposal *Proposal, vpnContract *VPNAccessContract) error {
	// Apply settings changes from proposal data
	settings := vpnContract.GlobalSettings

	for key, value := range proposal.ProposalData {
		switch key {
		case "basicAccessPrice":
			if price, ok := value.(float64); ok {
				settings.BasicAccessPrice = price
			}
		case "standardAccessPrice":
			if price, ok := value.(float64); ok {
				settings.StandardAccessPrice = price
			}
		case "premiumAccessPrice":
			if price, ok := value.(float64); ok {
				settings.PremiumAccessPrice = price
			}
		case "maxSessionsPerDay":
			if sessions, ok := value.(int); ok {
				settings.MaxSessionsPerDay = sessions
			}
		case "minStakeRequired":
			if stake, ok := value.(float64); ok {
				settings.MinStakeRequired = stake
			}
		case "violationPenalty":
			if penalty, ok := value.(int); ok {
				settings.ViolationPenalty = penalty
			}
		}
	}

	proposal.Status = "executed"
	proposal.ExecutedAt = time.Now().Unix()

	return nil
}

func (g *GovernanceContract) executeEmergencyProposal(proposal *Proposal) error {
	// Handle emergency proposals (shutdown, restrictions, etc.)
	if action, exists := proposal.ProposalData["action"]; exists {
		switch action {
		case "shutdown":
			g.EmergencySettings.EmergencyActive = true
			g.EmergencySettings.EmergencyReason = proposal.Description
			g.EmergencySettings.EmergencyActivatedBy = proposal.ProposedBy
			g.EmergencySettings.EmergencyActivatedAt = time.Now().Unix()
		case "resolve":
			g.EmergencySettings.EmergencyActive = false
		}
	}

	proposal.Status = "executed"
	proposal.ExecutedAt = time.Now().Unix()

	return nil
}

// GetProposalsByStatus returns proposals filtered by status
func (g *GovernanceContract) GetProposalsByStatus(status string) []*Proposal {
	var proposals []*Proposal
	for _, proposal := range g.Proposals {
		if proposal.Status == status {
			proposals = append(proposals, proposal)
		}
	}
	return proposals
}

// GetDelegateVotingStats returns voting statistics for a delegate
func (g *GovernanceContract) GetDelegateVotingStats(address string) map[string]interface{} {
	delegate, exists := g.AuthorizedDelegates[address]
	if !exists {
		return nil
	}

	return map[string]interface{}{
		"address":       delegate.Address,
		"name":          delegate.Name,
		"votingHistory": delegate.VotingHistory,
		"lastActivity":  delegate.LastActivity,
		"active":        delegate.Active,
		"votingPower":   delegate.VotingPower,
	}
}

// GetGovernanceStats returns overall governance statistics
func (g *GovernanceContract) GetGovernanceStats() map[string]interface{} {
	activeProposals := 0
	passedProposals := 0
	rejectedProposals := 0

	for _, proposal := range g.Proposals {
		switch proposal.Status {
		case "active":
			activeProposals++
		case "passed", "executed":
			passedProposals++
		case "rejected":
			rejectedProposals++
		}
	}

	activeDelegates := 0
	for _, delegate := range g.AuthorizedDelegates {
		if delegate.Active {
			activeDelegates++
		}
	}

	return map[string]interface{}{
		"totalProposals":    len(g.Proposals),
		"activeProposals":   activeProposals,
		"passedProposals":   passedProposals,
		"rejectedProposals": rejectedProposals,
		"totalDelegates":    len(g.AuthorizedDelegates),
		"activeDelegates":   activeDelegates,
		"emergencyActive":   g.EmergencySettings.EmergencyActive,
	}
}
