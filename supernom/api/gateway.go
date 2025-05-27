package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"supernom/contracts"

	"github.com/gorilla/mux"
)

// SuperNomGateway handles API requests for the decentralized VPN system
type SuperNomGateway struct {
	VPNContract        *contracts.VPNAccessContract
	GovernanceContract *contracts.GovernanceContract
	Port               string
	ServerMux          *mux.Router
}

// APIResponse represents a standardized API response
type APIResponse struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data,omitempty"`
	Error     string      `json:"error,omitempty"`
	Timestamp int64       `json:"timestamp"`
	RequestID string      `json:"requestId"`
}

// VPNConfigResponse contains VPN configuration for client
type VPNConfigResponse struct {
	SessionID      string `json:"sessionId"`
	VPNNodeID      string `json:"vpnNodeId"`
	Config         string `json:"config"` // WireGuard config
	ExpiresAt      int64  `json:"expiresAt"`
	BandwidthLimit int64  `json:"bandwidthLimit"`
	AccessLevel    int    `json:"accessLevel"`
}

// PurchaseRequest represents a VPN access purchase request
type PurchaseRequest struct {
	WalletAddress    string  `json:"walletAddress"`
	TokenType        string  `json:"tokenType"`
	Amount           float64 `json:"amount"`
	Duration         int64   `json:"duration"` // seconds
	GeographicZone   string  `json:"geographicZone"`
	PaymentSignature string  `json:"paymentSignature"` // blockchain transaction signature
}

// NewSuperNomGateway creates a new API gateway instance
func NewSuperNomGateway(port string) *SuperNomGateway {
	// Initialize contracts (in production, these would be loaded from blockchain)
	vpnContract := contracts.NewVPNAccessContract("AdNe1234567890123456789012345678901234567890123456789012345678901234")
	governanceContract := contracts.NewGovernanceContract(
		"AdNe1234567890123456789012345678901234567890123456789012345678901234",
		[]string{"AdNe1234567890123456789012345678901234567890123456789012345678901234"}, // initial delegates
	)

	gateway := &SuperNomGateway{
		VPNContract:        vpnContract,
		GovernanceContract: governanceContract,
		Port:               port,
		ServerMux:          mux.NewRouter(),
	}

	gateway.setupRoutes()
	return gateway
}

// setupRoutes configures all API endpoints
func (g *SuperNomGateway) setupRoutes() {
	// VPN Access endpoints
	g.ServerMux.HandleFunc("/auth/check", g.checkAuthHandler).Methods("GET")
	g.ServerMux.HandleFunc("/auth/purchase", g.purchaseAccessHandler).Methods("POST")
	g.ServerMux.HandleFunc("/auth/config", g.getVPNConfigHandler).Methods("GET")
	g.ServerMux.HandleFunc("/auth/revoke", g.revokeAccessHandler).Methods("POST")
	g.ServerMux.HandleFunc("/auth/status", g.getSessionStatusHandler).Methods("GET")

	// User management endpoints
	g.ServerMux.HandleFunc("/user/reputation", g.getUserReputationHandler).Methods("GET")
	g.ServerMux.HandleFunc("/user/sessions", g.getUserSessionsHandler).Methods("GET")
	g.ServerMux.HandleFunc("/user/stake", g.stakeTokensHandler).Methods("POST")

	// Governance endpoints
	g.ServerMux.HandleFunc("/governance/proposals", g.getProposalsHandler).Methods("GET")
	g.ServerMux.HandleFunc("/governance/proposal", g.createProposalHandler).Methods("POST")
	g.ServerMux.HandleFunc("/governance/vote", g.castVoteHandler).Methods("POST")
	g.ServerMux.HandleFunc("/governance/emergency", g.triggerEmergencyHandler).Methods("POST")

	// System status endpoints
	g.ServerMux.HandleFunc("/status", g.systemStatusHandler).Methods("GET")
	g.ServerMux.HandleFunc("/stats", g.systemStatsHandler).Methods("GET")
	g.ServerMux.HandleFunc("/health", g.healthCheckHandler).Methods("GET")

	// Add CORS middleware
	g.ServerMux.Use(g.corsMiddleware)
	g.ServerMux.Use(g.rateLimitMiddleware)
	g.ServerMux.Use(g.loggingMiddleware)
}

// checkAuthHandler verifies if a wallet has active VPN access
func (g *SuperNomGateway) checkAuthHandler(w http.ResponseWriter, r *http.Request) {
	wallet := r.URL.Query().Get("wallet")
	if wallet == "" {
		g.sendErrorResponse(w, "wallet parameter required", http.StatusBadRequest)
		return
	}

	session, authorized := g.VPNContract.CheckAuthorization(wallet)
	if !authorized {
		g.sendResponse(w, map[string]interface{}{
			"authorized": false,
			"message":    "No active session found",
		})
		return
	}

	g.sendResponse(w, map[string]interface{}{
		"authorized":     true,
		"sessionId":      session.SessionID,
		"expiresAt":      session.ExpiresAt,
		"accessLevel":    session.AccessLevel,
		"bandwidthLimit": session.BandwidthLimit,
		"bandwidthUsed":  session.BandwidthUsed,
		"vpnNodeId":      session.VPNNodeID,
		"geographicZone": session.GeographicZone,
	})
}

// purchaseAccessHandler handles VPN access purchase requests
func (g *SuperNomGateway) purchaseAccessHandler(w http.ResponseWriter, r *http.Request) {
	var req PurchaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		g.sendErrorResponse(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.WalletAddress == "" || req.Amount <= 0 || req.Duration <= 0 {
		g.sendErrorResponse(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// In production, verify payment signature against blockchain transaction
	if !g.verifyPaymentSignature(req) {
		g.sendErrorResponse(w, "Invalid payment signature", http.StatusUnauthorized)
		return
	}

	// Purchase access through smart contract
	session, err := g.VPNContract.PurchaseAccess(
		req.WalletAddress,
		req.TokenType,
		req.Amount,
		req.Duration,
		req.GeographicZone,
	)
	if err != nil {
		g.sendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	g.sendResponse(w, map[string]interface{}{
		"sessionId":      session.SessionID,
		"expiresAt":      session.ExpiresAt,
		"accessLevel":    session.AccessLevel,
		"bandwidthLimit": session.BandwidthLimit,
		"vpnNodeId":      session.VPNNodeID,
		"geographicZone": session.GeographicZone,
		"message":        "VPN access purchased successfully",
	})
}

// getVPNConfigHandler provides WireGuard configuration for authorized users
func (g *SuperNomGateway) getVPNConfigHandler(w http.ResponseWriter, r *http.Request) {
	wallet := r.URL.Query().Get("wallet")
	if wallet == "" {
		g.sendErrorResponse(w, "wallet parameter required", http.StatusBadRequest)
		return
	}

	session, authorized := g.VPNContract.CheckAuthorization(wallet)
	if !authorized {
		g.sendErrorResponse(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	// Generate WireGuard configuration
	config := g.generateWireGuardConfig(session)

	response := VPNConfigResponse{
		SessionID:      session.SessionID,
		VPNNodeID:      session.VPNNodeID,
		Config:         config,
		ExpiresAt:      session.ExpiresAt,
		BandwidthLimit: session.BandwidthLimit,
		AccessLevel:    session.AccessLevel,
	}

	g.sendResponse(w, response)
}

// revokeAccessHandler terminates a VPN session
func (g *SuperNomGateway) revokeAccessHandler(w http.ResponseWriter, r *http.Request) {
	var reqData map[string]string
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		g.sendErrorResponse(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	sessionID := reqData["sessionId"]
	reason := reqData["reason"]
	revokedBy := reqData["revokedBy"]

	if sessionID == "" {
		g.sendErrorResponse(w, "sessionId required", http.StatusBadRequest)
		return
	}

	err := g.VPNContract.RevokeAccess(sessionID, reason, revokedBy)
	if err != nil {
		g.sendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	g.sendResponse(w, map[string]interface{}{
		"message": "Session revoked successfully",
	})
}

// getSessionStatusHandler returns detailed session information
func (g *SuperNomGateway) getSessionStatusHandler(w http.ResponseWriter, r *http.Request) {
	wallet := r.URL.Query().Get("wallet")
	if wallet == "" {
		g.sendErrorResponse(w, "wallet parameter required", http.StatusBadRequest)
		return
	}

	session, err := g.VPNContract.GetActiveSession(wallet)
	if err != nil {
		g.sendErrorResponse(w, err.Error(), http.StatusNotFound)
		return
	}

	g.sendResponse(w, session)
}

// getUserReputationHandler returns user reputation information
func (g *SuperNomGateway) getUserReputationHandler(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	if address == "" {
		g.sendErrorResponse(w, "address parameter required", http.StatusBadRequest)
		return
	}

	reputation, err := g.VPNContract.GetUserReputation(address)
	if err != nil {
		g.sendErrorResponse(w, err.Error(), http.StatusNotFound)
		return
	}

	g.sendResponse(w, reputation)
}

// getUserSessionsHandler returns user session history
func (g *SuperNomGateway) getUserSessionsHandler(w http.ResponseWriter, r *http.Request) {
	wallet := r.URL.Query().Get("wallet")
	if wallet == "" {
		g.sendErrorResponse(w, "wallet parameter required", http.StatusBadRequest)
		return
	}

	// Filter sessions for the specific wallet
	var userSessions []*contracts.VPNSession
	for _, session := range g.VPNContract.Sessions {
		if session.WalletAddress == wallet {
			userSessions = append(userSessions, session)
		}
	}

	g.sendResponse(w, map[string]interface{}{
		"sessions": userSessions,
		"count":    len(userSessions),
	})
}

// systemStatusHandler returns overall system status
func (g *SuperNomGateway) systemStatusHandler(w http.ResponseWriter, r *http.Request) {
	vpnStats := g.VPNContract.GetContractStats()
	govStats := g.GovernanceContract.GetGovernanceStats()

	status := map[string]interface{}{
		"service":    "SuperNom VPN Gateway",
		"version":    "1.0.0",
		"status":     "operational",
		"timestamp":  time.Now().Unix(),
		"vpn":        vpnStats,
		"governance": govStats,
		"emergency":  g.GovernanceContract.EmergencySettings.EmergencyActive,
	}

	g.sendResponse(w, status)
}

// systemStatsHandler returns detailed system statistics
func (g *SuperNomGateway) systemStatsHandler(w http.ResponseWriter, r *http.Request) {
	vpnStats := g.VPNContract.GetContractStats()
	govStats := g.GovernanceContract.GetGovernanceStats()

	// Calculate additional metrics
	var totalBandwidthUsed int64
	activeSessionsCount := 0
	for _, session := range g.VPNContract.Sessions {
		totalBandwidthUsed += session.BandwidthUsed
		if session.Active && session.ExpiresAt > time.Now().Unix() {
			activeSessionsCount++
		}
	}

	stats := map[string]interface{}{
		"vpnContract":        vpnStats,
		"governanceContract": govStats,
		"totalBandwidthUsed": totalBandwidthUsed,
		"activeSessionsNow":  activeSessionsCount,
		"uptime":             time.Now().Unix() - g.VPNContract.CreatedAt,
	}

	g.sendResponse(w, stats)
}

// healthCheckHandler provides health check endpoint
func (g *SuperNomGateway) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	health := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
		"service":   "SuperNom Gateway",
		"emergency": g.GovernanceContract.EmergencySettings.EmergencyActive,
	}

	g.sendResponse(w, health)
}

// Helper functions

// generateWireGuardConfig creates a WireGuard configuration for a session
func (g *SuperNomGateway) generateWireGuardConfig(session *contracts.VPNSession) string {
	// In production, this would generate actual WireGuard keys and configuration
	// This is a placeholder configuration
	config := fmt.Sprintf(`[Interface]
PrivateKey = CLIENT_PRIVATE_KEY_PLACEHOLDER
Address = 10.0.%d.2/24
DNS = 1.1.1.1, 8.8.8.8

[Peer]
PublicKey = SERVER_PUBLIC_KEY_PLACEHOLDER
Endpoint = %s:51820
AllowedIPs = 0.0.0.0/0
PersistentKeepalive = 25

# Session Info
# SessionID: %s
# AccessLevel: %d
# ExpiresAt: %s
# BandwidthLimit: %d bytes
`,
		session.AccessLevel,                                  // subnet based on access level
		g.getVPNNodeEndpoint(session.VPNNodeID),              // VPN server endpoint
		session.SessionID,                                    // session identifier
		session.AccessLevel,                                  // access level
		time.Unix(session.ExpiresAt, 0).Format(time.RFC3339), // expiration
		session.BandwidthLimit,                               // bandwidth limit
	)

	return config
}

// getVPNNodeEndpoint returns the endpoint for a VPN node
func (g *SuperNomGateway) getVPNNodeEndpoint(nodeID string) string {
	// In production, this would query actual VPN node registry
	// For now, return placeholder based on node ID
	if strings.Contains(nodeID, "level1") {
		return "vpn1.supernom.network"
	} else if strings.Contains(nodeID, "level2") {
		return "vpn2.supernom.network"
	} else if strings.Contains(nodeID, "level3") {
		return "vpn3.supernom.network"
	}
	return "vpn.supernom.network"
}

// verifyPaymentSignature verifies blockchain payment signature
func (g *SuperNomGateway) verifyPaymentSignature(req PurchaseRequest) bool {
	// In production, this would verify the payment against blockchain
	// For now, return true if signature is provided
	return req.PaymentSignature != ""
}

// Middleware functions

// corsMiddleware handles CORS headers
func (g *SuperNomGateway) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// rateLimitMiddleware implements basic rate limiting
func (g *SuperNomGateway) rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// In production, implement proper rate limiting
		// For now, just pass through
		next.ServeHTTP(w, r)
	})
}

// loggingMiddleware logs requests
func (g *SuperNomGateway) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}

// Response helpers

// sendResponse sends a successful API response
func (g *SuperNomGateway) sendResponse(w http.ResponseWriter, data interface{}) {
	response := APIResponse{
		Success:   true,
		Data:      data,
		Timestamp: time.Now().Unix(),
		RequestID: generateSessionID("req", time.Now().Unix()),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// sendErrorResponse sends an error API response
func (g *SuperNomGateway) sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	response := APIResponse{
		Success:   false,
		Error:     message,
		Timestamp: time.Now().Unix(),
		RequestID: generateSessionID("err", time.Now().Unix()),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// Start starts the API gateway server
func (g *SuperNomGateway) Start() error {
	log.Printf("Starting SuperNom Gateway on port %s", g.Port)
	log.Printf("Available endpoints:")
	log.Printf("  GET  /auth/check - Check VPN authorization")
	log.Printf("  POST /auth/purchase - Purchase VPN access")
	log.Printf("  GET  /auth/config - Get VPN configuration")
	log.Printf("  GET  /status - System status")
	log.Printf("  GET  /health - Health check")

	return http.ListenAndServe(":"+g.Port, g.ServerMux)
}

// Governance-related handlers (additional endpoints)

// getProposalsHandler returns governance proposals
func (g *SuperNomGateway) getProposalsHandler(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	if status == "" {
		status = "active"
	}

	proposals := g.GovernanceContract.GetProposalsByStatus(status)
	g.sendResponse(w, map[string]interface{}{
		"proposals": proposals,
		"count":     len(proposals),
		"status":    status,
	})
}

// createProposalHandler creates a new governance proposal
func (g *SuperNomGateway) createProposalHandler(w http.ResponseWriter, r *http.Request) {
	var reqData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		g.sendErrorResponse(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	proposalType := reqData["type"].(string)
	title := reqData["title"].(string)
	description := reqData["description"].(string)
	targetAddress := ""
	if addr, exists := reqData["targetAddress"]; exists {
		targetAddress = addr.(string)
	}
	proposedBy := reqData["proposedBy"].(string)

	var additionalData map[string]interface{}
	if data, exists := reqData["data"]; exists {
		additionalData = data.(map[string]interface{})
	}

	proposal, err := g.GovernanceContract.CreateProposal(
		proposalType, title, description, targetAddress, proposedBy, additionalData,
	)
	if err != nil {
		g.sendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	g.sendResponse(w, proposal)
}

// castVoteHandler allows delegates to vote on proposals
func (g *SuperNomGateway) castVoteHandler(w http.ResponseWriter, r *http.Request) {
	var reqData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		g.sendErrorResponse(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	proposalID := reqData["proposalId"].(string)
	delegateAddress := reqData["delegateAddress"].(string)
	vote := reqData["vote"].(bool)

	err := g.GovernanceContract.CastVote(proposalID, delegateAddress, vote)
	if err != nil {
		g.sendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	g.sendResponse(w, map[string]interface{}{
		"message": "Vote cast successfully",
	})
}

// triggerEmergencyHandler triggers emergency shutdown
func (g *SuperNomGateway) triggerEmergencyHandler(w http.ResponseWriter, r *http.Request) {
	var reqData map[string]string
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		g.sendErrorResponse(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	reason := reqData["reason"]
	triggeredBy := reqData["triggeredBy"]

	err := g.GovernanceContract.TriggerEmergency(reason, triggeredBy)
	if err != nil {
		g.sendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	g.sendResponse(w, map[string]interface{}{
		"message": "Emergency triggered successfully",
	})
}

// stakeTokensHandler handles token staking for reputation
func (g *SuperNomGateway) stakeTokensHandler(w http.ResponseWriter, r *http.Request) {
	var reqData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		g.sendErrorResponse(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	address := reqData["address"].(string)
	amount := reqData["amount"].(float64)

	// Get or create reputation
	reputation := g.VPNContract.Reputation[address]
	if reputation == nil {
		reputation = &contracts.UserReputation{
			Address:         address,
			TrustScore:      100,
			AccessTier:      1,
			ReputationLevel: "new",
			LastActivity:    time.Now().Unix(),
		}
		g.VPNContract.Reputation[address] = reputation
	}

	// Update stake
	reputation.StakeLocked += amount
	reputation.LastActivity = time.Now().Unix()

	g.sendResponse(w, map[string]interface{}{
		"message":    "Tokens staked successfully",
		"totalStake": reputation.StakeLocked,
		"trustScore": reputation.TrustScore,
		"accessTier": reputation.AccessTier,
	})
}

// generateSessionID creates a unique session ID
func generateSessionID(prefix string, timestamp int64) string {
	hash := sha256.Sum256([]byte(prefix + strconv.FormatInt(timestamp, 10)))
	return hex.EncodeToString(hash[:])[:32]
}
