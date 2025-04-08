package core

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

// Block represents a block in the blockchain
type Block struct {
	Index        uint64        `json:"index"`
	PreviousHash string        `json:"previousHash"`
	Timestamp    int64         `json:"timestamp"`
	Data         []Transaction `json:"data"`
	Hash         string        `json:"hash"`
	Validator    string        `json:"validator"`
	Signature    string        `json:"signature"`
}

// Blockchain represents the blockchain
type Blockchain struct {
	chain        []Block
	transactions []Transaction
	mu           sync.RWMutex
}

// NewBlockchain creates a new blockchain with a genesis block
func NewBlockchain() *Blockchain {
	bc := &Blockchain{
		chain:        []Block{},
		transactions: []Transaction{},
	}

	// Create genesis block
	genesisBlock := Block{
		Index:        0,
		PreviousHash: "0",
		Timestamp:    time.Now().Unix(),
		Data:         []Transaction{},
		Hash:         "",
		Validator:    "genesis",
		Signature:    "genesis",
	}

	genesisBlock.Hash = CalculateHash(genesisBlock)
	bc.chain = append(bc.chain, genesisBlock)

	return bc
}

// NewBlockchainWithGenesis creates a new blockchain with a specific genesis block
func NewBlockchainWithGenesis(genesisBlock Block) *Blockchain {
	bc := &Blockchain{
		chain:        []Block{},
		transactions: []Transaction{},
	}

	// Add the genesis block
	bc.chain = append(bc.chain, genesisBlock)

	return bc
}

// ReplaceChain safely replaces the blockchain's chain with a new one
func (bc *Blockchain) ReplaceChain(newChain []Block) {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	bc.chain = make([]Block, len(newChain))
	copy(bc.chain, newChain)
	bc.transactions = []Transaction{} // Clear pending transactions
}

// AddBlock adds a new block to the blockchain
func (bc *Blockchain) AddBlock(block Block) error {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	lastBlock := bc.chain[len(bc.chain)-1]

	// Validate block
	if block.Index != lastBlock.Index+1 {
		return fmt.Errorf("invalid block index")
	}

	if block.PreviousHash != lastBlock.Hash {
		return fmt.Errorf("invalid previous hash")
	}

	// Verify block hash
	calculatedHash := CalculateHash(block)
	if calculatedHash != block.Hash {
		return fmt.Errorf("invalid block hash")
	}

	// Verify transaction prefixes
	for _, tx := range block.Data {
		if tx.ID[:4] != "AdNe" {
			return fmt.Errorf("transaction ID must start with 'AdNe'")
		}
	}

	// Add block to chain
	bc.chain = append(bc.chain, block)

	// Remove transactions that are now in the block
	bc.transactions = []Transaction{}

	return nil
}

// GetLastBlock returns the last block in the blockchain
func (bc *Blockchain) GetLastBlock() Block {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	return bc.chain[len(bc.chain)-1]
}

// AddTransaction adds a new transaction to the pending transactions
func (bc *Blockchain) AddTransaction(tx Transaction) error {
	// Validate transaction prefix
	if tx.ID[:4] != "AdNe" {
		return fmt.Errorf("transaction ID must start with 'AdNe'")
	}

	bc.mu.Lock()
	defer bc.mu.Unlock()

	bc.transactions = append(bc.transactions, tx)
	return nil
}

// GetPendingTransactions returns all pending transactions
func (bc *Blockchain) GetPendingTransactions() []Transaction {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	return bc.transactions
}

// GetBlockCount returns the number of blocks in the blockchain
func (bc *Blockchain) GetBlockCount() int {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	return len(bc.chain)
}

// GetChain returns a copy of the blockchain
func (bc *Blockchain) GetChain() []Block {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	// Create a copy of the chain
	chainCopy := make([]Block, len(bc.chain))
	copy(chainCopy, bc.chain)

	return chainCopy
}

// GetBlockByIndex returns a block by its index
func (bc *Blockchain) GetBlockByIndex(index uint64) (Block, error) {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	if int(index) >= len(bc.chain) {
		return Block{}, fmt.Errorf("block index out of range")
	}

	return bc.chain[index], nil
}

// CalculateHash calculates the hash of a block
func CalculateHash(block Block) string {
	record := fmt.Sprintf("%d%s%d%v%s", block.Index, block.PreviousHash, block.Timestamp, block.Data, block.Validator)
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}
