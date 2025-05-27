package contracts

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// VPNAccessContract manages decentralized VPN access through blockchain payments
type VPNAccessContract struct {
	// Contract state
	Sessions       map[string]*VPNSession     `json:"sessions"`
	Reputation     map[string]*UserReputation `json:"reputation"`
	GlobalSettings *GlobalVPNSettings         `json:"globalSettings"`
	Blacklist      map[string]*BlacklistEntry `json:"blacklist"`

	// Contract metadata
	ContractID string `json:"contractId"`
	Owner      string `json:"owner"`
	Version    string `json:"version"`
	CreatedAt  int64  `json:"createdAt"`
}

// VPNSession represents an active VPN session
type VPNSession struct {
	WalletAddress  string  `json:"walletAddress"`
	SessionID      string  `json:"sessionId"`
	ExpiresAt      int64   `json:"expiresAt"`
	CreatedAt      int64   `json:"createdAt"`
	PaidAmount     float64 `json:"paidAmount"`
	TokenType      string  `json:"tokenType"`
	AccessLevel    int     `json:"accessLevel"` // 1=basic, 2=standard, 3=premium
	VPNNodeID      string  `json:"vpnNodeId"`
	SessionHash    string  `json:"sessionHash"`
	Active         bool    `json:"active"`
	BandwidthLimit int64   `json:"bandwidthLimit"` // bytes
	BandwidthUsed  int64   `json:"bandwidthUsed"`  // bytes
	GeographicZone string  `json:"geographicZone"`
}

// UserReputation tracks user behavior and trust score
type UserReputation struct {
	Address         string  `json:"address"`
	TrustScore      int     `json:"trustScore"` // 0-1000
	TotalSessions   int     `json:"totalSessions"`
	CleanSessions   int     `json:"cleanSessions"`
	Violations      int     `json:"violations"`
	StakeLocked     float64 `json:"stakeLocked"`
	LastActivity    int64   `json:"lastActivity"`
	AccessTier      int     `json:"accessTier"`      // 1=new, 2=verified, 3=trusted
	ReputationLevel string  `json:"reputationLevel"` // "new", "good", "excellent", "flagged"
}

// GlobalVPNSettings contains system-wide configuration
type GlobalVPNSettings struct {
	MaxSessionDuration  int64    `json:"maxSessionDuration"`  // seconds
	BasicAccessPrice    float64  `json:"basicAccessPrice"`    // BNM
	StandardAccessPrice float64  `json:"standardAccessPrice"` // BNM
	PremiumAccessPrice  float64  `json:"premiumAccessPrice"`  // BNM
	MinStakeRequired    float64  `json:"minStakeRequired"`    // BNM
	ReputationDecayRate float64  `json:"reputationDecayRate"` // per day
	MaxSessionsPerDay   int      `json:"maxSessionsPerDay"`
	ViolationPenalty    int      `json:"violationPenalty"`  // trust score reduction
	TrustedThreshold    int      `json:"trustedThreshold"`  // trust score for tier 3
	VerifiedThreshold   int      `json:"verifiedThreshold"` // trust score for tier 2
	EmergencyShutdown   bool     `json:"emergencyShutdown"`
	EnableGeofencing    bool     `json:"enableGeofencing"`
	RestrictedCountries []string `json:"restrictedCountries"`
}

// BlacklistEntry represents a banned or flagged address
type BlacklistEntry struct {
	Address     string `json:"address"`
	Reason      string `json:"reason"`
	BannedAt    int64  `json:"bannedAt"`
	BannedBy    string `json:"bannedBy"` // delegate address
	Permanent   bool   `json:"permanent"`
	ExpiresAt   int64  `json:"expiresAt"`
	AppealCount int    `json:"appealCount"`
}

// NewVPNAccessContract creates a new VPN access contract instance
func NewVPNAccessContract(owner string) *VPNAccessContract {
	contractID := generateContractID(owner)

	return &VPNAccessContract{
		Sessions:   make(map[string]*VPNSession),
		Reputation: make(map[string]*UserReputation),
		Blacklist:  make(map[string]*BlacklistEntry),
		ContractID: contractID,
		Owner:      owner,
		Version:    "1.0.0",
		CreatedAt:  time.Now().Unix(),
		GlobalSettings: &GlobalVPNSettings{
			MaxSessionDuration:  86400, // 24 hours
			BasicAccessPrice:    10.0,  // 10 BNM for 1 hour
			StandardAccessPrice: 50.0,  // 50 BNM for 6 hours
			PremiumAccessPrice:  150.0, // 150 BNM for 24 hours
			MinStakeRequired:    100.0, // 100 BNM minimum stake
			ReputationDecayRate: 0.1,   // -0.1 trust score per day inactive
			MaxSessionsPerDay:   5,
			ViolationPenalty:    50,  // -50 trust score per violation
			TrustedThreshold:    800, // 800+ trust score for tier 3
			VerifiedThreshold:   400, // 400+ trust score for tier 2
			EmergencyShutdown:   false,
			EnableGeofencing:    true,
			RestrictedCountries: []string{"high-risk-zone-1", "high-risk-zone-2"},
		},
	}
}

// PurchaseAccess handles VPN access purchase with payment validation
func (c *VPNAccessContract) PurchaseAccess(walletAddress, tokenType string, amount float64, duration int64, geographicZone string) (*VPNSession, error) {
	// Emergency shutdown check
	if c.GlobalSettings.EmergencyShutdown {
		return nil, fmt.Errorf("VPN service temporarily unavailable")
	}

	// Validate wallet address format
	if !isValidAddress(walletAddress) {
		return nil, fmt.Errorf("invalid wallet address format")
	}

	// Check blacklist
	if entry, exists := c.Blacklist[walletAddress]; exists {
		if entry.Permanent || entry.ExpiresAt > time.Now().Unix() {
			return nil, fmt.Errorf("address is blacklisted: %s", entry.Reason)
		}
	}

	// Check geographic restrictions
	if c.GlobalSettings.EnableGeofencing {
		for _, restricted := range c.GlobalSettings.RestrictedCountries {
			if strings.EqualFold(geographicZone, restricted) {
				return nil, fmt.Errorf("geographic zone restricted")
			}
		}
	}

	// Initialize or get user reputation
	reputation := c.getOrCreateReputation(walletAddress)

	// Check daily session limits
	if c.countTodaySessions(walletAddress) >= c.GlobalSettings.MaxSessionsPerDay {
		return nil, fmt.Errorf("daily session limit exceeded")
	}

	// Validate payment amount
	accessLevel := c.determineAccessLevel(amount, duration)
	if accessLevel == 0 {
		return nil, fmt.Errorf("insufficient payment for requested duration")
	}

	// Check stake requirement for new users
	if reputation.AccessTier == 1 && reputation.StakeLocked < c.GlobalSettings.MinStakeRequired {
		return nil, fmt.Errorf("minimum stake required for new users: %.2f %s", c.GlobalSettings.MinStakeRequired, tokenType)
	}

	// Create session
	sessionID := generateSessionID(walletAddress, time.Now().Unix())
	sessionHash := generateSessionHash(walletAddress, sessionID)
	expiresAt := time.Now().Unix() + duration

	// Apply reputation-based limits
	if reputation.TrustScore < 200 {
		// Restrict new/low-trust users
		if duration > 3600 { // max 1 hour for low trust
			duration = 3600
			expiresAt = time.Now().Unix() + duration
		}
	}

	session := &VPNSession{
		WalletAddress:  walletAddress,
		SessionID:      sessionID,
		ExpiresAt:      expiresAt,
		CreatedAt:      time.Now().Unix(),
		PaidAmount:     amount,
		TokenType:      tokenType,
		AccessLevel:    accessLevel,
		VPNNodeID:      c.selectVPNNode(geographicZone, accessLevel),
		SessionHash:    sessionHash,
		Active:         true,
		BandwidthLimit: c.calculateBandwidthLimit(accessLevel, duration),
		BandwidthUsed:  0,
		GeographicZone: geographicZone,
	}

	// Store session
	c.Sessions[sessionID] = session

	// Update reputation
	reputation.TotalSessions++
	reputation.LastActivity = time.Now().Unix()

	return session, nil
}

// CheckAuthorization verifies if a wallet has active VPN access
func (c *VPNAccessContract) CheckAuthorization(walletAddress string) (*VPNSession, bool) {
	// Check blacklist first
	if entry, exists := c.Blacklist[walletAddress]; exists {
		if entry.Permanent || entry.ExpiresAt > time.Now().Unix() {
			return nil, false
		}
	}

	// Find active session
	for _, session := range c.Sessions {
		if session.WalletAddress == walletAddress &&
			session.Active &&
			session.ExpiresAt > time.Now().Unix() {
			return session, true
		}
	}

	return nil, false
}

// RevokeAccess immediately terminates a session
func (c *VPNAccessContract) RevokeAccess(sessionID, reason string, revokedBy string) error {
	session, exists := c.Sessions[sessionID]
	if !exists {
		return fmt.Errorf("session not found")
	}

	session.Active = false

	// Update reputation if revoked due to violation
	if reason != "user_requested" {
		reputation := c.getOrCreateReputation(session.WalletAddress)
		reputation.Violations++
		reputation.TrustScore -= c.GlobalSettings.ViolationPenalty
		if reputation.TrustScore < 0 {
			reputation.TrustScore = 0
		}
	}

	return nil
}

// UpdateBandwidthUsage tracks bandwidth consumption
func (c *VPNAccessContract) UpdateBandwidthUsage(sessionID string, bytesUsed int64) error {
	session, exists := c.Sessions[sessionID]
	if !exists {
		return fmt.Errorf("session not found")
	}

	session.BandwidthUsed += bytesUsed

	// Check bandwidth limit
	if session.BandwidthUsed > session.BandwidthLimit {
		session.Active = false
		return fmt.Errorf("bandwidth limit exceeded")
	}

	return nil
}

// CompleteSession marks a session as completed successfully
func (c *VPNAccessContract) CompleteSession(sessionID string) error {
	session, exists := c.Sessions[sessionID]
	if !exists {
		return fmt.Errorf("session not found")
	}

	session.Active = false

	// Update reputation positively
	reputation := c.getOrCreateReputation(session.WalletAddress)
	reputation.CleanSessions++

	// Increase trust score for successful sessions
	if session.ExpiresAt <= time.Now().Unix() { // completed full duration
		reputation.TrustScore += 5
		if reputation.TrustScore > 1000 {
			reputation.TrustScore = 1000
		}
	}

	// Update access tier based on trust score
	c.updateAccessTier(reputation)

	return nil
}

// Helper functions

func (c *VPNAccessContract) getOrCreateReputation(address string) *UserReputation {
	if rep, exists := c.Reputation[address]; exists {
		return rep
	}

	reputation := &UserReputation{
		Address:         address,
		TrustScore:      100, // Starting score
		AccessTier:      1,   // New user
		ReputationLevel: "new",
		LastActivity:    time.Now().Unix(),
	}

	c.Reputation[address] = reputation
	return reputation
}

func (c *VPNAccessContract) countTodaySessions(address string) int {
	count := 0
	today := time.Now().Unix() - 86400 // last 24 hours

	for _, session := range c.Sessions {
		if session.WalletAddress == address && session.CreatedAt > today {
			count++
		}
	}

	return count
}

func (c *VPNAccessContract) determineAccessLevel(amount float64, duration int64) int {
	settings := c.GlobalSettings

	if duration <= 3600 && amount >= settings.BasicAccessPrice {
		return 1 // Basic: 1 hour
	}
	if duration <= 21600 && amount >= settings.StandardAccessPrice {
		return 2 // Standard: 6 hours
	}
	if duration <= 86400 && amount >= settings.PremiumAccessPrice {
		return 3 // Premium: 24 hours
	}

	return 0 // Insufficient payment
}

func (c *VPNAccessContract) selectVPNNode(zone string, accessLevel int) string {
	// For now, return a deterministic node ID based on zone and level
	// In production, this would query available VPN nodes
	return fmt.Sprintf("vpn-node-%s-level%d", zone, accessLevel)
}

func (c *VPNAccessContract) calculateBandwidthLimit(accessLevel int, duration int64) int64 {
	// GB per hour based on access level
	var gbPerHour int64
	switch accessLevel {
	case 1:
		gbPerHour = 5 // 5 GB/hour basic
	case 2:
		gbPerHour = 15 // 15 GB/hour standard
	case 3:
		gbPerHour = 50 // 50 GB/hour premium
	default:
		gbPerHour = 1
	}

	hours := duration / 3600
	if hours == 0 {
		hours = 1
	}

	return gbPerHour * hours * 1024 * 1024 * 1024 // Convert to bytes
}

func (c *VPNAccessContract) updateAccessTier(reputation *UserReputation) {
	if reputation.TrustScore >= c.GlobalSettings.TrustedThreshold {
		reputation.AccessTier = 3
		reputation.ReputationLevel = "excellent"
	} else if reputation.TrustScore >= c.GlobalSettings.VerifiedThreshold {
		reputation.AccessTier = 2
		reputation.ReputationLevel = "good"
	} else if reputation.TrustScore < 100 {
		reputation.AccessTier = 1
		reputation.ReputationLevel = "flagged"
	} else {
		reputation.AccessTier = 1
		reputation.ReputationLevel = "new"
	}
}

// Utility functions

func generateContractID(owner string) string {
	hash := sha256.Sum256([]byte(owner + strconv.FormatInt(time.Now().Unix(), 10)))
	return "vpn-" + hex.EncodeToString(hash[:])[:16]
}

func generateSessionID(address string, timestamp int64) string {
	hash := sha256.Sum256([]byte(address + strconv.FormatInt(timestamp, 10)))
	return hex.EncodeToString(hash[:])[:32]
}

func generateSessionHash(address, sessionID string) string {
	hash := sha256.Sum256([]byte(address + sessionID))
	return hex.EncodeToString(hash[:])
}

func isValidAddress(address string) bool {
	return len(address) == 66 && strings.HasPrefix(address, "AdNe")
}

// GetActiveSession returns session info for API queries
func (c *VPNAccessContract) GetActiveSession(walletAddress string) (*VPNSession, error) {
	session, authorized := c.CheckAuthorization(walletAddress)
	if !authorized {
		return nil, fmt.Errorf("no active session found")
	}
	return session, nil
}

// GetUserReputation returns reputation info for an address
func (c *VPNAccessContract) GetUserReputation(address string) (*UserReputation, error) {
	reputation, exists := c.Reputation[address]
	if !exists {
		return nil, fmt.Errorf("no reputation data found")
	}
	return reputation, nil
}

// GetContractStats returns overall contract statistics
func (c *VPNAccessContract) GetContractStats() map[string]interface{} {
	activeSessions := 0
	totalRevenue := 0.0

	for _, session := range c.Sessions {
		if session.Active && session.ExpiresAt > time.Now().Unix() {
			activeSessions++
		}
		totalRevenue += session.PaidAmount
	}

	return map[string]interface{}{
		"activeSessions":   activeSessions,
		"totalSessions":    len(c.Sessions),
		"totalUsers":       len(c.Reputation),
		"totalRevenue":     totalRevenue,
		"blacklistedUsers": len(c.Blacklist),
		"contractVersion":  c.Version,
	}
}
