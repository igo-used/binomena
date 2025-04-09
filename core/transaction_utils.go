package core

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// GenerateTransactionID generates a unique transaction ID with "AdNe" prefix
func GenerateTransactionID() string {
	// Create a unique string based on current time and a random component
	data := fmt.Sprintf("%d-%d", time.Now().UnixNano(), time.Now().Unix())
	
	// Hash the data
	hash := sha256.Sum256([]byte(data))
	
	// Return the ID with "AdNe" prefix
	return "AdNe" + hex.EncodeToString(hash[:])[:60]
}