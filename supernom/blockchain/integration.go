package blockchain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"supernom/contracts"
)

// BinomenaIntegration handles communication with the Binomena blockchain
type BinomenaIntegration struct {
	NodeURL       string
	VPNContract   *contracts.VPNAccessContract
	GovContract   *contracts.GovernanceContract
	LastSyncBlock int64
	SyncInterval  time.Duration
}

// BlockchainTransaction represents a transaction on Binomena blockchain
type BlockchainTransaction struct {
	ID        string  `json:"id"`
	From      string  `json:"from"`
	To        string  `json:"to"`
	Amount    float64 `json:"amount"`
	TokenType string  `json:"tokenType"`
	Data      string  `json:"data"`
	Signature string  `json:"signature"`
	Timestamp int64   `json:"timestamp"`
}

// VPNPaymentData represents VPN payment transaction data
type VPNPaymentData struct {
	Action         string `json:"action"`   // "purchase_vpn", "stake_tokens"
	Duration       int64  `json:"duration"` // seconds
	GeographicZone string `json:"geographicZone"`
	SessionID      string `json:"sessionId,omitempty"`
}

// NewBinomenaIntegration creates a new blockchain integration instance
func NewBinomenaIntegration(nodeURL string) *BinomenaIntegration {
	return &BinomenaIntegration{
		NodeURL:      nodeURL,
		SyncInterval: 30 * time.Second, // Sync every 30 seconds
	}
}

// InitializeContracts deploys SuperNom contracts to the blockchain
func (b *BinomenaIntegration) InitializeContracts(ownerAddress string) error {
	log.Println("🚀 Initializing SuperNom contracts on Binomena blockchain...")

	// Create VPN access contract
	b.VPNContract = contracts.NewVPNAccessContract(ownerAddress)
	log.Printf("✅ VPN Access Contract initialized: %s", b.VPNContract.ContractID)

	// Get current delegates from blockchain
	delegates, err := b.getCurrentDelegates()
	if err != nil {
		log.Printf("⚠️  Warning: Could not fetch delegates, using default: %v", err)
		delegates = []string{ownerAddress} // fallback to owner
	}

	// Create governance contract
	b.GovContract = contracts.NewGovernanceContract(ownerAddress, delegates)
	log.Printf("✅ Governance Contract initialized: %s", b.GovContract.ContractID)

	// Start syncing with blockchain
	go b.startSyncLoop()

	return nil
}

// ProcessVPNPayment handles VPN access payments from blockchain
func (b *BinomenaIntegration) ProcessVPNPayment(tx *BlockchainTransaction) error {
	log.Printf("💰 Processing VPN payment: %s -> %s (%.2f %s)",
		tx.From, tx.To, tx.Amount, tx.TokenType)

	// Parse VPN payment data
	var paymentData VPNPaymentData
	if err := json.Unmarshal([]byte(tx.Data), &paymentData); err != nil {
		return fmt.Errorf("invalid VPN payment data: %v", err)
	}

	switch paymentData.Action {
	case "purchase_vpn":
		return b.processPurchaseVPN(tx, &paymentData)
	case "stake_tokens":
		return b.processStakeTokens(tx, &paymentData)
	default:
		return fmt.Errorf("unknown VPN action: %s", paymentData.Action)
	}
}

// processPurchaseVPN handles VPN access purchase
func (b *BinomenaIntegration) processPurchaseVPN(tx *BlockchainTransaction, data *VPNPaymentData) error {
	session, err := b.VPNContract.PurchaseAccess(
		tx.From,
		tx.TokenType,
		tx.Amount,
		data.Duration,
		data.GeographicZone,
	)
	if err != nil {
		return fmt.Errorf("VPN purchase failed: %v", err)
	}

	log.Printf("✅ VPN session created: %s (expires: %s)",
		session.SessionID,
		time.Unix(session.ExpiresAt, 0).Format(time.RFC3339))

	// Store session ID in blockchain transaction data for future reference
	data.SessionID = session.SessionID

	return nil
}

// processStakeTokens handles token staking for reputation
func (b *BinomenaIntegration) processStakeTokens(tx *BlockchainTransaction, data *VPNPaymentData) error {
	// Get or create reputation
	reputation := b.VPNContract.Reputation[tx.From]
	if reputation == nil {
		reputation = &contracts.UserReputation{
			Address:         tx.From,
			TrustScore:      100,
			AccessTier:      1,
			ReputationLevel: "new",
			LastActivity:    time.Now().Unix(),
		}
		b.VPNContract.Reputation[tx.From] = reputation
	}

	// Update stake
	reputation.StakeLocked += tx.Amount
	reputation.LastActivity = time.Now().Unix()

	log.Printf("🔒 Tokens staked: %s staked %.2f %s (total: %.2f)",
		tx.From, tx.Amount, tx.TokenType, reputation.StakeLocked)

	return nil
}

// VerifyPaymentOnChain verifies a payment transaction exists on blockchain
func (b *BinomenaIntegration) VerifyPaymentOnChain(signature string, amount float64, from string) (bool, error) {
	// Query blockchain for transaction with this signature
	url := fmt.Sprintf("%s/transaction/%s", b.NodeURL, signature)

	resp, err := http.Get(url)
	if err != nil {
		return false, fmt.Errorf("failed to query blockchain: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("transaction not found on blockchain")
	}

	var tx BlockchainTransaction
	if err := json.NewDecoder(resp.Body).Decode(&tx); err != nil {
		return false, fmt.Errorf("failed to decode transaction: %v", err)
	}

	// Verify transaction details
	if tx.From != from || tx.Amount != amount {
		return false, fmt.Errorf("transaction details do not match")
	}

	return true, nil
}

// SyncWithBlockchain synchronizes SuperNom state with blockchain
func (b *BinomenaIntegration) SyncWithBlockchain() error {
	log.Println("🔄 Syncing SuperNom with Binomena blockchain...")

	// Get latest transactions related to SuperNom
	transactions, err := b.getSupernomTransactions()
	if err != nil {
		return fmt.Errorf("failed to fetch transactions: %v", err)
	}

	processedCount := 0
	for _, tx := range transactions {
		// Skip already processed transactions
		if tx.Timestamp <= b.LastSyncBlock {
			continue
		}

		// Process VPN-related transactions
		if b.isVPNTransaction(&tx) {
			if err := b.ProcessVPNPayment(&tx); err != nil {
				log.Printf("❌ Failed to process VPN transaction %s: %v", tx.ID, err)
				continue
			}
			processedCount++
		}

		// Update last processed block
		if tx.Timestamp > b.LastSyncBlock {
			b.LastSyncBlock = tx.Timestamp
		}
	}

	if processedCount > 0 {
		log.Printf("✅ Processed %d SuperNom transactions", processedCount)
	}

	return nil
}

// startSyncLoop starts the continuous sync with blockchain
func (b *BinomenaIntegration) startSyncLoop() {
	ticker := time.NewTicker(b.SyncInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := b.SyncWithBlockchain(); err != nil {
				log.Printf("⚠️  Sync error: %v", err)
			}
		}
	}
}

// Helper functions

// getCurrentDelegates fetches current delegates from Binomena blockchain
func (b *BinomenaIntegration) getCurrentDelegates() ([]string, error) {
	url := fmt.Sprintf("%s/delegates", b.NodeURL)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var delegateResponse struct {
		Delegates []struct {
			Address string `json:"address"`
		} `json:"delegates"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&delegateResponse); err != nil {
		return nil, err
	}

	var addresses []string
	for _, delegate := range delegateResponse.Delegates {
		addresses = append(addresses, delegate.Address)
	}

	return addresses, nil
}

// getSupernomTransactions fetches SuperNom-related transactions
func (b *BinomenaIntegration) getSupernomTransactions() ([]BlockchainTransaction, error) {
	// Query blockchain for transactions with SuperNom data
	url := fmt.Sprintf("%s/transactions/search?query=supernom&limit=100", b.NodeURL)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var txResponse struct {
		Transactions []BlockchainTransaction `json:"transactions"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&txResponse); err != nil {
		return nil, err
	}

	return txResponse.Transactions, nil
}

// isVPNTransaction checks if a transaction is VPN-related
func (b *BinomenaIntegration) isVPNTransaction(tx *BlockchainTransaction) bool {
	// Check if transaction data contains VPN actions
	var paymentData VPNPaymentData
	if err := json.Unmarshal([]byte(tx.Data), &paymentData); err != nil {
		return false
	}

	return paymentData.Action == "purchase_vpn" || paymentData.Action == "stake_tokens"
}

// SubmitToBlockchain submits SuperNom data to blockchain
func (b *BinomenaIntegration) SubmitToBlockchain(data interface{}) error {
	// Convert data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}

	// Submit to blockchain (this would be a real transaction in production)
	url := fmt.Sprintf("%s/transactions", b.NodeURL)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to submit to blockchain: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("blockchain submission failed with status: %d", resp.StatusCode)
	}

	return nil
}

// GetContractStats returns statistics about SuperNom contracts
func (b *BinomenaIntegration) GetContractStats() map[string]interface{} {
	vpnStats := b.VPNContract.GetContractStats()
	govStats := b.GovContract.GetGovernanceStats()

	return map[string]interface{}{
		"vpnContract":        vpnStats,
		"governanceContract": govStats,
		"lastSyncBlock":      b.LastSyncBlock,
		"nodeURL":            b.NodeURL,
		"syncInterval":       b.SyncInterval.String(),
	}
}

// Emergency functions

// TriggerEmergencyShutdown triggers emergency shutdown via governance
func (b *BinomenaIntegration) TriggerEmergencyShutdown(reason, triggeredBy string) error {
	log.Printf("🚨 EMERGENCY SHUTDOWN TRIGGERED: %s (by: %s)", reason, triggeredBy)

	// Trigger emergency in governance contract
	if err := b.GovContract.TriggerEmergency(reason, triggeredBy); err != nil {
		return err
	}

	// Update VPN contract settings
	b.VPNContract.GlobalSettings.EmergencyShutdown = true

	// Submit emergency event to blockchain
	emergencyData := map[string]interface{}{
		"action":      "emergency_shutdown",
		"reason":      reason,
		"triggeredBy": triggeredBy,
		"timestamp":   time.Now().Unix(),
	}

	return b.SubmitToBlockchain(emergencyData)
}

// ResolveEmergency resolves emergency shutdown
func (b *BinomenaIntegration) ResolveEmergency(resolvedBy string) error {
	log.Printf("✅ Emergency resolved by: %s", resolvedBy)

	// Resolve emergency in governance contract
	if err := b.GovContract.ResolveEmergency(resolvedBy); err != nil {
		return err
	}

	// Update VPN contract settings
	b.VPNContract.GlobalSettings.EmergencyShutdown = false

	// Submit resolution to blockchain
	resolutionData := map[string]interface{}{
		"action":     "emergency_resolved",
		"resolvedBy": resolvedBy,
		"timestamp":  time.Now().Unix(),
	}

	return b.SubmitToBlockchain(resolutionData)
}
