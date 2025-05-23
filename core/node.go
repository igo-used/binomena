package core

import (
	"fmt"
	"sync"
	"time"
)

// BlockchainInterface defines the interface for blockchain implementations
type BlockchainInterface interface {
	GetLastBlock() Block
	AddBlock(block Block) error
	GetBlockCount() int
	GetChain() []Block
	GetBlockByIndex(index uint64) (Block, error)
	AddTransaction(tx Transaction) error
	GetPendingTransactions() []Transaction
	ReplaceChain(newChain []Block)
}

// TokenInterface defines the interface for token implementations
type TokenInterface interface {
	Transfer(from, to string, amount float64) error
	GetBalance(address string) float64
	GetCirculatingSupply() float64
	Burn(amount float64)
}

// Node represents a node in the Binomena network
type Node struct {
	blockchain       BlockchainInterface
	consensus        Consensus
	token            TokenInterface
	peers            map[string]Peer
	isRunning        bool
	mu               sync.RWMutex
	stopChan         chan struct{}
	validatorAddress string
}

// Consensus interface for consensus mechanisms
type Consensus interface {
	ValidateBlock(block Block) bool
	SelectValidator(validators []string, stakes map[string]float64) string
}

// Token interface for token operations (deprecated, use TokenInterface)
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
func NewNode(blockchain BlockchainInterface, consensus Consensus, token TokenInterface, validatorAddress string) *Node {
	return &Node{
		blockchain:       blockchain,
		consensus:        consensus,
		token:            token,
		peers:            make(map[string]Peer),
		isRunning:        false,
		stopChan:         make(chan struct{}),
		validatorAddress: validatorAddress,
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
		Validator:    n.validatorAddress, // Use the founder's address as validator
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
