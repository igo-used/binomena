package token

import (
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/igo-used/binomena/database"
	"gorm.io/gorm"
)

// BinomTokenDB represents the database-backed Binom (BNM) token
type BinomTokenDB struct {
	maxSupply float64
	mu        sync.RWMutex
}

// NewBinomTokenWithDB creates a new database-backed Binom token
func NewBinomTokenWithDB() *BinomTokenDB {
	return &BinomTokenDB{
		maxSupply: 1000000000.0, // 1 billion
	}
}

// Transfer transfers tokens from one address to another using database
func (bt *BinomTokenDB) Transfer(from, to string, amount float64) error {
	bt.mu.Lock()
	defer bt.mu.Unlock()

	// Start a database transaction
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get sender balance
	var fromBalance database.TokenBalance
	result := tx.Where("address = ?", from).First(&fromBalance)
	if result.Error == gorm.ErrRecordNotFound {
		tx.Rollback()
		return fmt.Errorf("sender address not found")
	}
	if result.Error != nil {
		tx.Rollback()
		return fmt.Errorf("failed to get sender balance: %v", result.Error)
	}

	// Check if sender has enough balance
	if fromBalance.Balance < amount {
		tx.Rollback()
		return fmt.Errorf("insufficient balance")
	}

	// Get or create receiver balance
	var toBalance database.TokenBalance
	result = tx.Where("address = ?", to).First(&toBalance)
	if result.Error == gorm.ErrRecordNotFound {
		toBalance = database.TokenBalance{
			Address: to,
			Balance: 0,
		}
		if err := tx.Create(&toBalance).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create receiver balance: %v", err)
		}
	} else if result.Error != nil {
		tx.Rollback()
		return fmt.Errorf("failed to get receiver balance: %v", result.Error)
	}

	// Update balances
	fromBalance.Balance -= amount
	toBalance.Balance += amount

	// Save updated balances
	if err := tx.Save(&fromBalance).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update sender balance: %v", err)
	}

	if err := tx.Save(&toBalance).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update receiver balance: %v", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

// GetBalance returns the balance of an address from database
func (bt *BinomTokenDB) GetBalance(address string) float64 {
	bt.mu.RLock()
	defer bt.mu.RUnlock()

	var balance database.TokenBalance
	result := database.DB.Where("address = ?", address).First(&balance)
	if result.Error == gorm.ErrRecordNotFound {
		return 0.0
	}
	if result.Error != nil {
		log.Printf("Error getting balance for %s: %v", address, result.Error)
		return 0.0
	}

	return balance.Balance
}

// GetCirculatingSupply returns the circulating supply from database
func (bt *BinomTokenDB) GetCirculatingSupply() float64 {
	bt.mu.RLock()
	defer bt.mu.RUnlock()

	var supply database.SystemState
	result := database.DB.Where("key = ?", "circulating_supply").First(&supply)
	if result.Error != nil {
		log.Printf("Error getting circulating supply: %v", result.Error)
		return bt.maxSupply
	}

	supplyFloat, err := strconv.ParseFloat(supply.Value, 64)
	if err != nil {
		log.Printf("Error parsing circulating supply: %v", err)
		return bt.maxSupply
	}

	return supplyFloat
}

// Burn burns tokens, reducing the circulating supply in database
func (bt *BinomTokenDB) Burn(amount float64) {
	bt.mu.Lock()
	defer bt.mu.Unlock()

	// Get current circulating supply
	var supply database.SystemState
	result := database.DB.Where("key = ?", "circulating_supply").First(&supply)
	if result.Error != nil {
		log.Printf("Error getting circulating supply for burn: %v", result.Error)
		return
	}

	supplyFloat, err := strconv.ParseFloat(supply.Value, 64)
	if err != nil {
		log.Printf("Error parsing circulating supply for burn: %v", err)
		return
	}

	// Burn tokens
	newSupply := supplyFloat - amount
	supply.Value = fmt.Sprintf("%.8f", newSupply)

	// Save updated supply
	if err := database.DB.Save(&supply).Error; err != nil {
		log.Printf("Error saving burned supply: %v", err)
		return
	}

	log.Printf("Burned %.2f BNM tokens. New circulating supply: %.2f", amount, newSupply)
}

// Mint mints new tokens, increasing the circulating supply
func (bt *BinomTokenDB) Mint(to string, amount float64) error {
	bt.mu.Lock()
	defer bt.mu.Unlock()

	// Get current circulating supply
	var supply database.SystemState
	result := database.DB.Where("key = ?", "circulating_supply").First(&supply)
	if result.Error != nil {
		return fmt.Errorf("failed to get circulating supply: %v", result.Error)
	}

	supplyFloat, err := strconv.ParseFloat(supply.Value, 64)
	if err != nil {
		return fmt.Errorf("failed to parse circulating supply: %v", err)
	}

	// Check if minting would exceed max supply
	if supplyFloat+amount > bt.maxSupply {
		return fmt.Errorf("minting would exceed max supply")
	}

	// Start database transaction
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get or create receiver balance
	var toBalance database.TokenBalance
	result = tx.Where("address = ?", to).First(&toBalance)
	if result.Error == gorm.ErrRecordNotFound {
		toBalance = database.TokenBalance{
			Address: to,
			Balance: amount,
		}
		if err := tx.Create(&toBalance).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create balance for minting: %v", err)
		}
	} else if result.Error != nil {
		tx.Rollback()
		return fmt.Errorf("failed to get balance for minting: %v", result.Error)
	} else {
		toBalance.Balance += amount
		if err := tx.Save(&toBalance).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update balance for minting: %v", err)
		}
	}

	// Update circulating supply
	newSupply := supplyFloat + amount
	supply.Value = fmt.Sprintf("%.8f", newSupply)
	if err := tx.Save(&supply).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update circulating supply: %v", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit minting transaction: %v", err)
	}

	return nil
}
