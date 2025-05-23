package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Block model for PostgreSQL
type Block struct {
	ID           uint   `gorm:"primaryKey"`
	Index        uint64 `gorm:"uniqueIndex;not null"`
	PreviousHash string `gorm:"size:64;not null"`
	Timestamp    int64  `gorm:"not null"`
	Data         string `gorm:"type:jsonb"` // Store transactions as JSON
	Hash         string `gorm:"size:64;uniqueIndex;not null"`
	Validator    string `gorm:"size:66;not null"`
	Signature    string `gorm:"size:144;not null"`
}

// Wallet model for PostgreSQL
type Wallet struct {
	ID      uint    `gorm:"primaryKey"`
	Address string  `gorm:"size:66;uniqueIndex;not null"`
	Balance float64 `gorm:"type:decimal(20,8);default:0"`
}

// Transaction model for PostgreSQL
type Transaction struct {
	ID        uint    `gorm:"primaryKey"`
	TxID      string  `gorm:"size:66;uniqueIndex;not null"`
	FromAddr  string  `gorm:"size:66;not null;index"`
	ToAddr    string  `gorm:"size:66;not null;index"`
	Amount    float64 `gorm:"type:decimal(20,8);not null"`
	Timestamp int64   `gorm:"not null"`
	Signature string  `gorm:"size:144;not null"`
	BlockID   *uint   `gorm:"index"` // Reference to block
}

// Contract model for PostgreSQL
type Contract struct {
	ID             uint    `gorm:"primaryKey"`
	ContractID     string  `gorm:"size:66;uniqueIndex;not null"`
	Owner          string  `gorm:"size:66;not null;index"`
	Name           string  `gorm:"size:100;not null"`
	Code           []byte  `gorm:"type:bytea;not null"`
	DeployedAt     int64   `gorm:"not null"`
	LastExecuted   int64   `gorm:"default:0"`
	ExecutionCount uint64  `gorm:"default:0"`
	TotalGasUsed   float64 `gorm:"type:decimal(20,8);default:0"`
}

// AuditEvent model for PostgreSQL
type AuditEvent struct {
	ID        uint   `gorm:"primaryKey"`
	EventID   string `gorm:"size:66;uniqueIndex;not null"`
	Timestamp int64  `gorm:"not null;index"`
	Level     int    `gorm:"not null;index"`
	Type      string `gorm:"size:50;not null;index"`
	Message   string `gorm:"type:text;not null"`
	Data      string `gorm:"type:jsonb"` // Store additional data as JSON
}

// TokenBalance model for tracking token balances
type TokenBalance struct {
	ID      uint    `gorm:"primaryKey"`
	Address string  `gorm:"size:66;uniqueIndex;not null"`
	Balance float64 `gorm:"type:decimal(20,8);not null;default:0"`
}

// SystemState model for storing system-wide state
type SystemState struct {
	ID          uint   `gorm:"primaryKey"`
	Key         string `gorm:"size:50;uniqueIndex;not null"`
	Value       string `gorm:"type:text;not null"`
	LastUpdated int64  `gorm:"not null"`
}

// ConnectDatabase connects to PostgreSQL database
func ConnectDatabase() error {
	// Get database URL from environment
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return fmt.Errorf("DATABASE_URL environment variable is required")
	}

	// Configure GORM logger
	gormLogger := logger.Default
	if os.Getenv("DEBUG") == "true" {
		gormLogger = logger.Default.LogMode(logger.Info)
	}

	// Connect to database
	var err error
	DB, err = gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	log.Println("Successfully connected to PostgreSQL database")
	return nil
}

// MigrateDatabase creates tables and runs migrations
func MigrateDatabase() error {
	err := DB.AutoMigrate(
		&Block{},
		&Wallet{},
		&Transaction{},
		&Contract{},
		&AuditEvent{},
		&TokenBalance{},
		&SystemState{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %v", err)
	}

	log.Println("Database migration completed successfully")
	return nil
}

// InitializeSystemState initializes system state with default values
func InitializeSystemState() error {
	// Initialize circulating supply
	var supply SystemState
	result := DB.Where("key = ?", "circulating_supply").First(&supply)
	if result.Error == gorm.ErrRecordNotFound {
		supply = SystemState{
			Key:         "circulating_supply",
			Value:       "1000000000.0", // 1 billion BNM
			LastUpdated: 0,
		}
		if err := DB.Create(&supply).Error; err != nil {
			return fmt.Errorf("failed to initialize circulating supply: %v", err)
		}
	}

	// Initialize treasury balance
	var treasury TokenBalance
	result = DB.Where("address = ?", "treasury").First(&treasury)
	if result.Error == gorm.ErrRecordNotFound {
		treasury = TokenBalance{
			Address: "treasury",
			Balance: 1000000000.0, // All tokens start in treasury
		}
		if err := DB.Create(&treasury).Error; err != nil {
			return fmt.Errorf("failed to initialize treasury: %v", err)
		}
	}

	log.Println("System state initialized successfully")
	return nil
}

// CloseDatabase closes the database connection
func CloseDatabase() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
