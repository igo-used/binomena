package core

import (
	"fmt"
	"sync"
	"time"
)

// Node represents a node in the Binomena network
type Node struct {
	blockchain *Blockchain
	consensus  Consensus
	token      Token
	peers      map[string]Peer
	isRunning  bool
	mu         sync.RWMutex
	stopChan   chan struct{}
}

// Consensus interface for consensus mechanisms
type Consensus interface {
	ValidateBlock(block Block) bool
	SelectValidator(validators []string, stakes map[string]float64) string
}

// Token interface for token operations
type Token interface {
	Transfer(from, to string, amount float64) error
	GetBalance(address string) float64
	GetCirculatingSupply() float64
	Burn(amount float64)
}

// Peer represents a peer node in the network
type Peer struct {
	ID      string
	Address string
}

// NewNode creates a new node
func NewNode(blockchain *Blockchain, consensus Consensus, token Token) *Node {
	return &Node{
		blockchain: blockchain,
		consensus:  consensus,
		token:      token,
		peers:      make(map[string]Peer),
		isRunning:  false,
		stopChan:   make(chan struct{}),
	}
}

// Start starts the node
func (n *Node) Start() {
	n.mu.Lock()
	if n.isRunning {
		n.mu.Unlock()
		return
	}
	n.isRunning = true
	n.mu.Unlock()

	// Start block creation loop
	go n.blockCreationLoop()
}

// Stop stops the node
func (n *Node) Stop() {
	n.mu.Lock()
	defer n.mu.Unlock()

	if !n.isRunning {
		return
	}

	close(n.stopChan)
	n.isRunning = false
}

// SubmitTransaction submits a transaction to the blockchain
func (n *Node) SubmitTransaction(tx Transaction) error {
	// Validate transaction
	if tx.From == "" || tx.To == "" || tx.Amount <= 0 {
		return fmt.Errorf("invalid transaction")
	}

	// Check if sender has enough balance
	senderBalance := n.token.GetBalance(tx.From)
	if senderBalance < tx.Amount {
		return fmt.Errorf("insufficient balance")
	}

	// Apply burn ratio (0.1%)
	burnAmount := tx.Amount * 0.001
	transferAmount := tx.Amount - burnAmount

	// Transfer tokens
	if err := n.token.Transfer(tx.From, tx.To, transferAmount); err != nil {
		return err
	}

	// Burn tokens
	n.token.Burn(burnAmount)

	// Add transaction to blockchain
	n.blockchain.AddTransaction(tx)

	return nil
}

// GetPeerCount returns the number of peers
func (n *Node) GetPeerCount() int {
	n.mu.RLock()
	defer n.mu.RUnlock()

	return len(n.peers)
}

// blockCreationLoop continuously creates new blocks
func (n *Node) blockCreationLoop() {
	ticker := time.NewTicker(10 * time.Second) // Create a block every 10 seconds
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			n.createNewBlock()
		case <-n.stopChan:
			return
		}
	}
}

// createNewBlock creates a new block and adds it to the blockchain
func (n *Node) createNewBlock() {
	// Get pending transactions
	transactions := n.blockchain.GetPendingTransactions()
	if len(transactions) == 0 {
		return // No transactions to process
	}

	// Get the last block
	lastBlock := n.blockchain.GetLastBlock()

	// Create new block
	newBlock := Block{
		Index:        lastBlock.Index + 1,
		PreviousHash: lastBlock.Hash,
		Timestamp:    time.Now().Unix(),
		Data:         transactions,
		Validator:    "self", // In a real implementation, this would be determined by the consensus mechanism
	}

	// Calculate hash
	newBlock.Hash = CalculateHash(newBlock)

	// Add block to blockchain
	if err := n.blockchain.AddBlock(newBlock); err != nil {
		fmt.Printf("Error adding block: %v\n", err)
	} else {
		fmt.Printf("Block #%d added to the blockchain\n", newBlock.Index)
	}
}
