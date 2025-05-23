package core

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/igo-used/binomena/database"
	"gorm.io/gorm"
)

// BlockchainDB represents the database-backed blockchain
type BlockchainDB struct {
	transactions []Transaction
	mu           sync.RWMutex
}

// NewBlockchainWithDB creates a new database-backed blockchain with a genesis block
func NewBlockchainWithDB() *BlockchainDB {
	bc := &BlockchainDB{
		transactions: []Transaction{},
	}

	// Check if genesis block exists
	var count int64
	database.DB.Model(&database.Block{}).Count(&count)

	if count == 0 {
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

		// Save genesis block to database
		if err := bc.saveBlockToDB(genesisBlock); err != nil {
			log.Printf("Error creating genesis block: %v", err)
		} else {
			log.Println("Genesis block created successfully")
		}
	}

	return bc
}

// saveBlockToDB saves a block to the database
func (bc *BlockchainDB) saveBlockToDB(block Block) error {
	// Serialize transactions to JSON
	transactionsJSON, err := json.Marshal(block.Data)
	if err != nil {
		return fmt.Errorf("failed to serialize transactions: %v", err)
	}

	dbBlock := database.Block{
		Index:        block.Index,
		PreviousHash: block.PreviousHash,
		Timestamp:    block.Timestamp,
		Data:         string(transactionsJSON),
		Hash:         block.Hash,
		Validator:    block.Validator,
		Signature:    block.Signature,
	}

	if err := database.DB.Create(&dbBlock).Error; err != nil {
		return fmt.Errorf("failed to save block to database: %v", err)
	}

	// Save individual transactions
	for _, tx := range block.Data {
		dbTx := database.Transaction{
			TxID:      tx.ID,
			FromAddr:  tx.From,
			ToAddr:    tx.To,
			Amount:    tx.Amount,
			Timestamp: tx.Timestamp,
			Signature: tx.Signature,
			BlockID:   &dbBlock.ID,
		}

		if err := database.DB.Create(&dbTx).Error; err != nil {
			log.Printf("Warning: Failed to save transaction %s: %v", tx.ID, err)
		}
	}

	return nil
}

// loadBlockFromDB loads a block from the database
func (bc *BlockchainDB) loadBlockFromDB(dbBlock database.Block) (Block, error) {
	// Deserialize transactions from JSON
	var transactions []Transaction
	if err := json.Unmarshal([]byte(dbBlock.Data), &transactions); err != nil {
		return Block{}, fmt.Errorf("failed to deserialize transactions: %v", err)
	}

	return Block{
		Index:        dbBlock.Index,
		PreviousHash: dbBlock.PreviousHash,
		Timestamp:    dbBlock.Timestamp,
		Data:         transactions,
		Hash:         dbBlock.Hash,
		Validator:    dbBlock.Validator,
		Signature:    dbBlock.Signature,
	}, nil
}

// AddBlock adds a new block to the blockchain database
func (bc *BlockchainDB) AddBlock(block Block) error {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	// Get the last block from database
	var lastDBBlock database.Block
	result := database.DB.Order("index desc").First(&lastDBBlock)
	if result.Error != nil {
		return fmt.Errorf("failed to get last block: %v", result.Error)
	}

	// Validate block
	if block.Index != lastDBBlock.Index+1 {
		return fmt.Errorf("invalid block index")
	}

	if block.PreviousHash != lastDBBlock.Hash {
		return fmt.Errorf("invalid previous hash")
	}

	// Verify block hash
	calculatedHash := CalculateHash(block)
	if calculatedHash != block.Hash {
		return fmt.Errorf("invalid block hash")
	}

	// Verify transaction prefixes
	for _, tx := range block.Data {
		if len(tx.ID) < 4 || tx.ID[:4] != "AdNe" {
			return fmt.Errorf("transaction ID must start with 'AdNe'")
		}
	}

	// Save block to database
	if err := bc.saveBlockToDB(block); err != nil {
		return err
	}

	// Clear pending transactions
	bc.transactions = []Transaction{}

	return nil
}

// GetLastBlock returns the last block in the blockchain
func (bc *BlockchainDB) GetLastBlock() Block {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	var dbBlock database.Block
	result := database.DB.Order("index desc").First(&dbBlock)
	if result.Error != nil {
		log.Printf("Error getting last block: %v", result.Error)
		// Return genesis block as fallback
		return Block{
			Index:        0,
			PreviousHash: "0",
			Timestamp:    time.Now().Unix(),
			Data:         []Transaction{},
			Hash:         "genesis",
			Validator:    "genesis",
			Signature:    "genesis",
		}
	}

	block, err := bc.loadBlockFromDB(dbBlock)
	if err != nil {
		log.Printf("Error loading block from DB: %v", err)
		// Return genesis block as fallback
		return Block{
			Index:        0,
			PreviousHash: "0",
			Timestamp:    time.Now().Unix(),
			Data:         []Transaction{},
			Hash:         "genesis",
			Validator:    "genesis",
			Signature:    "genesis",
		}
	}

	return block
}

// AddTransaction adds a new transaction to the pending transactions
func (bc *BlockchainDB) AddTransaction(tx Transaction) error {
	// Validate transaction prefix
	if len(tx.ID) < 4 || tx.ID[:4] != "AdNe" {
		return fmt.Errorf("transaction ID must start with 'AdNe'")
	}

	bc.mu.Lock()
	defer bc.mu.Unlock()

	bc.transactions = append(bc.transactions, tx)
	return nil
}

// GetPendingTransactions returns all pending transactions
func (bc *BlockchainDB) GetPendingTransactions() []Transaction {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	return bc.transactions
}

// GetBlockCount returns the number of blocks in the blockchain
func (bc *BlockchainDB) GetBlockCount() int {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	var count int64
	database.DB.Model(&database.Block{}).Count(&count)
	return int(count)
}

// GetChain returns a copy of the blockchain
func (bc *BlockchainDB) GetChain() []Block {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	var dbBlocks []database.Block
	database.DB.Order("index asc").Find(&dbBlocks)

	blocks := make([]Block, len(dbBlocks))
	for i, dbBlock := range dbBlocks {
		block, err := bc.loadBlockFromDB(dbBlock)
		if err != nil {
			log.Printf("Error loading block %d: %v", dbBlock.Index, err)
			continue
		}
		blocks[i] = block
	}

	return blocks
}

// GetBlockByIndex returns a block by its index
func (bc *BlockchainDB) GetBlockByIndex(index uint64) (Block, error) {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	var dbBlock database.Block
	result := database.DB.Where("index = ?", index).First(&dbBlock)
	if result.Error == gorm.ErrRecordNotFound {
		return Block{}, fmt.Errorf("block index out of range")
	}
	if result.Error != nil {
		return Block{}, fmt.Errorf("failed to get block: %v", result.Error)
	}

	return bc.loadBlockFromDB(dbBlock)
}

// ReplaceChain safely replaces the blockchain's chain with a new one
func (bc *BlockchainDB) ReplaceChain(newChain []Block) {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	// Start database transaction
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Delete all existing blocks
	if err := tx.Exec("DELETE FROM blocks").Error; err != nil {
		tx.Rollback()
		log.Printf("Error deleting blocks: %v", err)
		return
	}

	// Delete all existing transactions
	if err := tx.Exec("DELETE FROM transactions").Error; err != nil {
		tx.Rollback()
		log.Printf("Error deleting transactions: %v", err)
		return
	}

	// Insert new blocks
	for _, block := range newChain {
		// Serialize transactions to JSON
		transactionsJSON, err := json.Marshal(block.Data)
		if err != nil {
			tx.Rollback()
			log.Printf("Error serializing transactions: %v", err)
			return
		}

		dbBlock := database.Block{
			Index:        block.Index,
			PreviousHash: block.PreviousHash,
			Timestamp:    block.Timestamp,
			Data:         string(transactionsJSON),
			Hash:         block.Hash,
			Validator:    block.Validator,
			Signature:    block.Signature,
		}

		if err := tx.Create(&dbBlock).Error; err != nil {
			tx.Rollback()
			log.Printf("Error creating block: %v", err)
			return
		}

		// Save individual transactions
		for _, txData := range block.Data {
			dbTx := database.Transaction{
				TxID:      txData.ID,
				FromAddr:  txData.From,
				ToAddr:    txData.To,
				Amount:    txData.Amount,
				Timestamp: txData.Timestamp,
				Signature: txData.Signature,
				BlockID:   &dbBlock.ID,
			}

			if err := tx.Create(&dbTx).Error; err != nil {
				log.Printf("Warning: Failed to save transaction %s: %v", txData.ID, err)
			}
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		log.Printf("Error committing chain replacement: %v", err)
		return
	}

	// Clear pending transactions
	bc.transactions = []Transaction{}
}
