// Copyright [2025] [uJ1NO (Juxhino Kapllanaj)] [binomena.com] [adaneural.com].
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/igo-used/binomena/audit"
	"github.com/igo-used/binomena/consensus"
	"github.com/igo-used/binomena/core"
	"github.com/igo-used/binomena/database"
	"github.com/igo-used/binomena/p2p"
	"github.com/igo-used/binomena/smartcontract"
	"github.com/igo-used/binomena/token"
	"github.com/igo-used/binomena/wallet"
)

// Helper functions for audit service calls
func logAuditEvent(auditService any, level audit.SecurityLevel, eventType, message string, data interface{}) {
	if fileAudit, ok := auditService.(*audit.AuditService); ok {
		fileAudit.LogEvent(level, eventType, message, data)
	} else if dbAudit, ok := auditService.(*audit.AuditServiceDB); ok {
		dbAudit.LogEvent(level, eventType, message, data)
	}
}

func getAuditEvents(auditService any) []audit.AuditEvent {
	if fileAudit, ok := auditService.(*audit.AuditService); ok {
		return fileAudit.GetEvents()
	} else if dbAudit, ok := auditService.(*audit.AuditServiceDB); ok {
		return dbAudit.GetEvents()
	}
	return []audit.AuditEvent{}
}

func getAuditEventsByLevel(auditService any, level audit.SecurityLevel) []audit.AuditEvent {
	if fileAudit, ok := auditService.(*audit.AuditService); ok {
		return fileAudit.GetEventsByLevel(level)
	} else if dbAudit, ok := auditService.(*audit.AuditServiceDB); ok {
		return dbAudit.GetEventsByLevel(level)
	}
	return []audit.AuditEvent{}
}

func auditBlockchain(auditService any) {
	if fileAudit, ok := auditService.(*audit.AuditService); ok {
		fileAudit.AuditBlockchain()
	} else if dbAudit, ok := auditService.(*audit.AuditServiceDB); ok {
		dbAudit.AuditBlockchain()
	}
}

func main() {
	// Parse command line flags
	apiPort := flag.Int("api-port", 8080, "API server port")
	p2pPort := flag.Int("p2p-port", 9000, "P2P server port")
	bootstrapNode := flag.String("bootstrap", "", "Bootstrap node address (optional)")
	nodeID := flag.String("id", "", "Node identifier (optional)")
	useDB := flag.Bool("use-db", true, "Use database backend (default: true)")
	flag.Parse()

	// Check for PORT environment variable (required for Render deployment)
	if portEnv := os.Getenv("PORT"); portEnv != "" {
		if port, err := strconv.Atoi(portEnv); err == nil {
			*apiPort = port
			log.Printf("Using PORT environment variable: %d", port)
		}
	}

	// Set node identifier
	nodeName := *nodeID
	if nodeName == "" {
		nodeName = fmt.Sprintf("node-%d", *p2pPort)
	}

	// Initialize database connection if using DB backend
	var useDatabase bool
	if *useDB {
		if err := database.ConnectDatabase(); err != nil {
			log.Printf("Failed to connect to database: %v", err)
			log.Println("Falling back to file-based storage")
			useDatabase = false
		} else {
			log.Println("Connected to database successfully")

			// Run database migrations
			if err := database.MigrateDatabase(); err != nil {
				log.Printf("Failed to migrate database: %v", err)
				log.Println("Falling back to file-based storage")
				useDatabase = false
			} else {
				log.Println("Database migration completed")

				// Initialize system state
				if err := database.InitializeSystemState(); err != nil {
					log.Printf("Failed to initialize system state: %v", err)
				}

				useDatabase = true
			}
		}
	}

	// Initialize blockchain based on backend choice
	var blockchain core.BlockchainInterface
	if useDatabase {
		blockchain = core.NewBlockchainWithDB()
		log.Println("Using database-backed blockchain")
	} else {
		blockchain = core.NewBlockchain()
		log.Println("Using file-backed blockchain")
	}

	// Initialize token based on backend choice
	var binomToken core.TokenInterface
	if useDatabase {
		binomToken = token.NewBinomTokenWithDB()
		log.Println("Using database-backed token system")
	} else {
		binomToken = token.NewBinomToken()
		log.Println("Using file-backed token system")
	}

	// Initialize the DPoS consensus mechanism with founder and community addresses
	founderAddress := "AdNe6c3ce54e4371d056c7c566675ba16909eb2e9534"
	communityAddress := "AdNebaefd75d426056bffbc622bd9f334ed89450efae"

	dposConsensus := consensus.NewDPoSConsensus(founderAddress, communityAddress)

	// Register founder as the first delegate with their 400M BNM stake
	if err := dposConsensus.RegisterDelegate(founderAddress, 400000000.0); err != nil {
		log.Printf("Warning: Failed to register founder as delegate: %v", err)
	} else {
		log.Println("Founder registered as first delegate with 400M BNM stake")
	}

	// Initialize smart contract system based on backend choice
	var wasmVM *smartcontract.WasmVM
	var contractStorage interface{}
	var contractState interface{}

	if useDatabase {
		dbContractStorage, err := smartcontract.NewContractStorageWithDB()
		if err != nil {
			log.Fatalf("Failed to initialize database contract storage: %v", err)
		}
		contractStorage = dbContractStorage

		dbContractState, err := smartcontract.NewContractStateWithDB()
		if err != nil {
			log.Fatalf("Failed to initialize database contract state: %v", err)
		}
		contractState = dbContractState

		// For database backend, we need to work with concrete types
		if _, ok := binomToken.(*token.BinomTokenDB); ok {
			if _, ok := blockchain.(*core.BlockchainDB); ok {
				// Create a simple wrapper VM that can work with database types
				// For now, we'll use file-based VM as the database VM doesn't exist
				log.Println("Database VM not implemented yet, falling back to file VM for smart contracts")

				// Create temporary file-based implementations for VM
				tempToken := token.NewBinomToken()
				tempBlockchain := core.NewBlockchain()

				var err error
				wasmVM, err = smartcontract.NewWasmVM(tempToken, tempBlockchain)
				if err != nil {
					log.Fatalf("Failed to initialize WASM VM: %v", err)
				}
			} else {
				log.Fatalf("Expected database blockchain implementation")
			}
		} else {
			log.Fatalf("Expected database token implementation")
		}

		log.Println("Using database-backed smart contract storage with file-based VM")
	} else {
		fileContractStorage, err := smartcontract.NewContractStorage("./contracts")
		if err != nil {
			log.Fatalf("Failed to initialize file contract storage: %v", err)
		}
		contractStorage = fileContractStorage

		fileContractState, err := smartcontract.NewContractState("./contracts")
		if err != nil {
			log.Fatalf("Failed to initialize file contract state: %v", err)
		}
		contractState = fileContractState

		// For file backend, use the original types
		if fileToken, ok := binomToken.(*token.BinomToken); ok {
			if fileBlockchain, ok := blockchain.(*core.Blockchain); ok {
				wasmVM, err = smartcontract.NewWasmVM(fileToken, fileBlockchain)
				if err != nil {
					log.Fatalf("Failed to initialize WASM VM: %v", err)
				}
			} else {
				log.Fatalf("Expected file blockchain implementation")
			}
		} else {
			log.Fatalf("Expected file token implementation")
		}

		log.Println("Using file-backed smart contract system")
	}

	// Load existing contracts directly from storage
	var contracts []*smartcontract.Contract
	var err error
	if fileStorage, ok := contractStorage.(*smartcontract.ContractStorage); ok {
		contracts, err = fileStorage.LoadAllContracts()
	} else if dbStorage, ok := contractStorage.(*smartcontract.ContractStorageDB); ok {
		contracts, err = dbStorage.LoadAllContracts()
	}

	if err != nil {
		log.Printf("Warning: Failed to load contracts: %v", err)
	} else {
		for _, contract := range contracts {
			wasmVM.AddContract(contract)
		}
		log.Printf("Loaded %d contracts", len(contracts))
	}

	// Create contract API with concrete types
	var contractAPI *smartcontract.ContractAPI
	if fileStorage, ok := contractStorage.(*smartcontract.ContractStorage); ok {
		if fileState, ok := contractState.(*smartcontract.ContractState); ok {
			if fileToken, ok := binomToken.(*token.BinomToken); ok {
				contractAPI = smartcontract.NewContractAPI(wasmVM, fileStorage, fileState, fileToken)
			}
		}
	}

	// For database storage, create a working API with database backend
	if contractAPI == nil {
		log.Println("Creating contract API for database backend")
		// Create temporary file storage for API until database VM is implemented
		tempStorage, err := smartcontract.NewContractStorage("./temp_contracts")
		if err != nil {
			log.Printf("Warning: Failed to create temp contract storage: %v", err)
		}
		tempState, err := smartcontract.NewContractState("./temp_contracts")
		if err != nil {
			log.Printf("Warning: Failed to create temp contract state: %v", err)
		}

		// Use the actual database token system instead of temporary
		// Initialize with database token for balance access
		if tempStorage != nil && tempState != nil {
			// Create a wrapper for database token that implements the file token interface
			if dbToken, ok := binomToken.(*token.BinomTokenDB); ok {
				// Create temporary file token and sync balances for contract API
				tempToken := token.NewBinomToken()

				// Sync key balances from database to temp token for contract operations
				founderBalance := dbToken.GetBalance("AdNe6c3ce54e4371d056c7c566675ba16909eb2e9534")
				treasuryBalance := dbToken.GetBalance("AdNec13f53bb89865c7e2be8ff9aa43e84e26d226bf3")
				communityBalance := dbToken.GetBalance("AdNebaefd75d426056bffbc622bd9f334ed89450efae")

				// Transfer from treasury to sync balances in temp token
				if founderBalance > 0 {
					tempToken.Transfer("treasury", "AdNe6c3ce54e4371d056c7c566675ba16909eb2e9534", founderBalance)
				}
				if treasuryBalance > 0 {
					tempToken.Transfer("treasury", "AdNec13f53bb89865c7e2be8ff9aa43e84e26d226bf3", treasuryBalance)
				}
				if communityBalance > 0 {
					tempToken.Transfer("treasury", "AdNebaefd75d426056bffbc622bd9f334ed89450efae", communityBalance)
				}

				contractAPI = smartcontract.NewContractAPI(wasmVM, tempStorage, tempState, tempToken)
				log.Printf("Contract API initialized with synced balances: Founder=%.2f, Treasury=%.2f", founderBalance, treasuryBalance)
			} else {
				// Fallback to temporary token
				tempToken := token.NewBinomToken()
				contractAPI = smartcontract.NewContractAPI(wasmVM, tempStorage, tempState, tempToken)
				log.Println("Contract API initialized with temporary file-based storage")
			}
		} else {
			log.Println("Warning: Failed to initialize contract API")
		}
	}

	// Initialize audit service
	var auditService any

	if useDatabase {
		auditService = audit.NewAuditServiceWithDB(blockchain)
		log.Println("Using database-backed audit service")
	} else {
		if fileBlockchain, ok := blockchain.(*core.Blockchain); ok {
			auditService = audit.NewAuditService(fileBlockchain)
		} else {
			log.Fatalf("Expected file blockchain implementation for audit service")
		}
		log.Println("Using file-backed audit service")
	}

	// Create node
	node := core.NewNode(blockchain, dposConsensus, binomToken, "genesis")

	// Start the P2P network
	p2pAddress := fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", *p2pPort)
	var p2pNode *p2p.P2PNode
	if fileBlockchain, ok := blockchain.(*core.Blockchain); ok {
		p2pNode, err = p2p.NewP2PNode(fileBlockchain, p2pAddress)
	} else {
		// For database blockchain, create a temporary file blockchain for P2P
		tempBlockchain := core.NewBlockchain()
		p2pNode, err = p2p.NewP2PNode(tempBlockchain, p2pAddress)
		log.Println("Warning: Using temporary blockchain for P2P due to database backend")
	}
	if err != nil {
		log.Fatalf("Failed to start P2P node: %v", err)
	}

	// Connect to bootstrap node if provided
	if *bootstrapNode != "" {
		if err := p2pNode.ConnectToPeer(*bootstrapNode); err != nil {
			log.Printf("Warning: Failed to connect to bootstrap node: %v", err)
		} else {
			log.Printf("Connected to bootstrap node: %s", *bootstrapNode)
		}
	}

	// Start the node
	go node.Start()

	// Setup API server
	router := gin.Default()
	// Add CORS middleware
	router.Use(corsMiddleware())

	// Add security headers middleware
	router.Use(func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Next()
	})

	// Health check endpoint for Render
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().Unix(),
			"node":      nodeName,
		})
	})

	// Health check endpoint with trailing space (Render compatibility)
	router.GET("/health ", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().Unix(),
			"node":      nodeName,
		})
	})

	// API endpoints
	router.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"nodeId":      nodeName,
			"status":      "running",
			"blocks":      blockchain.GetBlockCount(),
			"peers":       p2pNode.GetPeerCount(),
			"wallets":     p2pNode.GetWalletCount(),
			"tokenSupply": binomToken.GetCirculatingSupply(),
		})
	})

	// Create wallet endpoint
	router.POST("/wallet", func(c *gin.Context) {
		newWallet, err := wallet.NewWallet()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Announce wallet to the network
		err = p2pNode.AnnounceWallet(newWallet.Address)
		if err != nil {
			log.Printf("Failed to announce wallet: %v", err)
		}

		c.JSON(http.StatusOK, gin.H{
			"address":    newWallet.Address,
			"privateKey": newWallet.ExportPrivateKey(),
		})

		// Log wallet creation
		logAuditEvent(auditService, audit.InfoLevel, "WalletCreated", fmt.Sprintf("New wallet created with address %s", newWallet.Address), nil)
	})

	// Import wallet endpoint
	router.POST("/wallet/import", func(c *gin.Context) {
		var request struct {
			PrivateKey string `json:"privateKey"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		importedWallet, err := wallet.ImportPrivateKey(request.PrivateKey)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Announce wallet to the network
		err = p2pNode.AnnounceWallet(importedWallet.Address)
		if err != nil {
			log.Printf("Failed to announce wallet: %v", err)
		}

		c.JSON(http.StatusOK, gin.H{
			"address": importedWallet.Address,
		})

		// Log wallet import
		logAuditEvent(auditService, audit.InfoLevel, "WalletImported", fmt.Sprintf("Wallet imported with address %s", importedWallet.Address), nil)
	})

	// Get wallet balance
	router.GET("/balance/:address", func(c *gin.Context) {
		address := c.Param("address")

		// Validate address format
		if len(address) < 4 || address[:4] != "AdNe" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address format"})
			return
		}

		balance := binomToken.GetBalance(address)

		c.JSON(http.StatusOK, gin.H{
			"address": address,
			"balance": balance,
		})
	})

	// Simple faucet endpoint (no admin key required for basic testing)
	router.POST("/faucet", func(c *gin.Context) {
		var request struct {
			Address string  `json:"address"`
			Amount  float64 `json:"amount"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate address
		if len(request.Address) < 4 || request.Address[:4] != "AdNe" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address format"})
			return
		}

		// Validate amount (limit to reasonable amounts for testing)
		if request.Amount <= 0 || request.Amount > 10000 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Amount must be between 0 and 10000"})
			return
		}

		// Transfer tokens from treasury to the address
		err := binomToken.Transfer("treasury", request.Address, request.Amount)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": fmt.Sprintf("Transferred %f BNM to %s", request.Amount, request.Address),
			"balance": binomToken.GetBalance(request.Address),
		})

		// Log faucet request
		logAuditEvent(auditService, audit.InfoLevel, "FaucetRequest", fmt.Sprintf("Transferred %f BNM to %s", request.Amount, request.Address), nil)
	})

	// NEW ENDPOINT: Distribute initial tokens to three wallets
	router.POST("/admin/distribute-initial-tokens", func(c *gin.Context) {
		var request struct {
			AdminKey         string  `json:"adminKey"`
			FounderAddress   string  `json:"founderAddress"`
			TreasuryAddress  string  `json:"treasuryAddress"`
			CommunityAddress string  `json:"communityAddress"`
			FounderPercent   float64 `json:"founderPercent"`
			TreasuryPercent  float64 `json:"treasuryPercent"`
			CommunityPercent float64 `json:"communityPercent"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Verify percentages add up to 100
		totalPercent := request.FounderPercent + request.TreasuryPercent + request.CommunityPercent
		if totalPercent != 100.0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Percentages must add up to 100"})
			return
		}

		// Get total supply from treasury
		totalSupply := binomToken.GetBalance("treasury")

		// Calculate token amounts
		founderAmount := totalSupply * (request.FounderPercent / 100.0)
		treasuryAmount := totalSupply * (request.TreasuryPercent / 100.0)
		communityAmount := totalSupply * (request.CommunityPercent / 100.0)

		// Transfer tokens to founder
		if err := binomToken.Transfer("treasury", request.FounderAddress, founderAmount); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to transfer to founder: %v", err)})
			return
		}

		// Transfer tokens to new treasury
		if err := binomToken.Transfer("treasury", request.TreasuryAddress, treasuryAmount); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to transfer to treasury: %v", err)})
			return
		}

		// Transfer tokens to community
		if err := binomToken.Transfer("treasury", request.CommunityAddress, communityAmount); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to transfer to community: %v", err)})
			return
		}

		// Log the distribution
		logAuditEvent(auditService, audit.InfoLevel, "InitialTokenDistribution", fmt.Sprintf("Distributed tokens: %f to founder, %f to treasury, %f to community",
			founderAmount, treasuryAmount, communityAmount), map[string]interface{}{
			"founder":   request.FounderAddress,
			"treasury":  request.TreasuryAddress,
			"community": request.CommunityAddress,
		})

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"distribution": map[string]interface{}{
				"founder": map[string]interface{}{
					"address": request.FounderAddress,
					"amount":  founderAmount,
					"percent": request.FounderPercent,
				},
				"treasury": map[string]interface{}{
					"address": request.TreasuryAddress,
					"amount":  treasuryAmount,
					"percent": request.TreasuryPercent,
				},
				"community": map[string]interface{}{
					"address": request.CommunityAddress,
					"amount":  communityAmount,
					"percent": request.CommunityPercent,
				},
			},
		})
	})

	// Transaction endpoint
	router.POST("/transaction", func(c *gin.Context) {
		var request struct {
			From       string  `json:"from"`
			To         string  `json:"to"`
			Amount     float64 `json:"amount"`
			PrivateKey string  `json:"privateKey"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Import wallet from private key
		senderWallet, err := wallet.ImportPrivateKey(request.PrivateKey)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid private key"})
			return
		}

		// Verify wallet address matches
		if senderWallet.Address != request.From {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Private key does not match sender address"})
			return
		}

		// Calculate fee (0.1% of transaction amount)
		feeRate := 0.001
		transactionFee := request.Amount * feeRate
		totalRequired := request.Amount + transactionFee

		// Check balance (sender pays both amount and fee)
		balance := binomToken.GetBalance(request.From)
		if balance < totalRequired {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":    "insufficient balance",
				"balance":  balance,
				"required": totalRequired,
				"amount":   request.Amount,
				"fee":      transactionFee,
			})
			return
		}

		// Transfer the exact amount to receiver
		if err := binomToken.Transfer(request.From, request.To, request.Amount); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Collect fee from sender separately
		if err := binomToken.Transfer(request.From, "treasury", transactionFee); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to collect fee: " + err.Error()})
			return
		}

		// Distribute fees according to DPoS rules
		if err := dposConsensus.DistributeFees(transactionFee, binomToken); err != nil {
			log.Printf("Failed to distribute fees: %v", err)
		}

		// Create transaction
		tx, err := core.NewTransaction(request.From, request.To, request.Amount, senderWallet)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Submit transaction
		if err := node.SubmitTransaction(*tx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Broadcast transaction to the network
		if err := p2pNode.BroadcastTransaction(*tx); err != nil {
			log.Printf("Failed to broadcast transaction: %v", err)
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "transaction submitted",
			"txId":   tx.ID,
			"amount": request.Amount,
			"fee":    transactionFee,
			"feeDistribution": gin.H{
				"delegates": transactionFee * 0.6,
				"burned":    transactionFee * 0.3,
				"community": transactionFee * 0.05,
				"founder":   transactionFee * 0.05,
			},
			"node": nodeName,
		})

		// Log transaction
		logAuditEvent(auditService, audit.InfoLevel, "TransactionSubmitted",
			fmt.Sprintf("Transaction %s: %s sent %.6f BNM to %s (fee: %.6f BNM)", tx.ID, tx.From, tx.Amount, tx.To, transactionFee), tx)
	})

	// Get peers endpoint
	router.GET("/peers", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"peers": p2pNode.GetPeers(),
			"count": p2pNode.GetPeerCount(),
		})
	})

	// Connect to peer endpoint
	router.POST("/peers", func(c *gin.Context) {
		var request struct {
			Address string `json:"address"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := p2pNode.ConnectToPeer(request.Address); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "connected",
			"peer":   request.Address,
		})
	})

	// Blockchain synchronization endpoints

	// Get block by index
	router.GET("/blocks/:index", func(c *gin.Context) {
		indexStr := c.Param("index")
		index, err := strconv.ParseUint(indexStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid block index"})
			return
		}

		block, err := blockchain.GetBlockByIndex(index)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, block)
	})

	// Get all blocks
	router.GET("/blocks", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"blocks": blockchain.GetChain(),
			"count":  blockchain.GetBlockCount(),
		})
	})

	// Sync blockchain with a peer
	router.POST("/sync", func(c *gin.Context) {
		var request struct {
			PeerAddress string `json:"peerAddress"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Get the peer's blockchain
		resp, err := http.Get(fmt.Sprintf("http://%s/blocks", request.PeerAddress))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Failed to connect to peer: %v", err)})
			return
		}
		defer resp.Body.Close()

		// Parse the response
		var peerBlockchain struct {
			Blocks []core.Block `json:"blocks"`
			Count  int          `json:"count"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&peerBlockchain); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Failed to parse peer blockchain: %v", err)})
			return
		}

		// Check if peer has more blocks
		localBlockCount := blockchain.GetBlockCount()
		if peerBlockchain.Count <= localBlockCount {
			c.JSON(http.StatusOK, gin.H{
				"status":      "no sync needed",
				"localBlocks": localBlockCount,
				"peerBlocks":  peerBlockchain.Count,
			})
			return
		}

		// Check if genesis blocks are different
		localGenesis, _ := blockchain.GetBlockByIndex(0)
		peerGenesis := peerBlockchain.Blocks[0]

		if localGenesis.Hash != peerGenesis.Hash {
			// Genesis blocks are different, we need to replace the entire chain
			log.Printf("Different genesis blocks detected. Replacing local chain with peer chain.")

			// Create a new blockchain with the peer's genesis block
			newBlockchain := core.NewBlockchainWithGenesis(peerGenesis)

			// Add all blocks from the peer
			for i := 1; i < peerBlockchain.Count; i++ {
				block := peerBlockchain.Blocks[i]
				if err := newBlockchain.AddBlock(block); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error":       fmt.Sprintf("Failed to add block %d: %v", i, err),
						"syncedUntil": i - 1,
					})
					return
				}

				// Process transactions in the block
				for _, tx := range block.Data {
					// Apply transaction effects (transfer tokens, burn fee)
					fee := tx.Amount * 0.001 // 0.1% fee
					transferAmount := tx.Amount - fee

					// Transfer tokens
					if err := binomToken.Transfer(tx.From, tx.To, transferAmount); err != nil {
						log.Printf("Warning: Failed to apply transaction effect: %v", err)
					}

					// Burn fee
					binomToken.Burn(fee)
				}
			}

			// Replace the blockchain safely
			blockchain.ReplaceChain(newBlockchain.GetChain())

			c.JSON(http.StatusOK, gin.H{
				"status":        "full chain replacement completed",
				"blocksAdded":   peerBlockchain.Count - 1,
				"newBlockCount": blockchain.GetBlockCount(),
			})
			return
		} else {
			// Genesis blocks are the same, just add missing blocks
			for i := localBlockCount; i < peerBlockchain.Count; i++ {
				block := peerBlockchain.Blocks[i]
				if err := blockchain.AddBlock(block); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error":       fmt.Sprintf("Failed to add block %d: %v", i, err),
						"syncedUntil": i - 1,
					})
					return
				}

				// Process transactions in the block
				for _, tx := range block.Data {
					// Apply transaction effects (transfer tokens, burn fee)
					fee := tx.Amount * 0.001 // 0.1% fee
					transferAmount := tx.Amount - fee

					// Transfer tokens
					if err := binomToken.Transfer(tx.From, tx.To, transferAmount); err != nil {
						log.Printf("Warning: Failed to apply transaction effect: %v", err)
					}

					// Burn fee
					binomToken.Burn(fee)
				}
			}

			c.JSON(http.StatusOK, gin.H{
				"status":        "sync completed",
				"blocksAdded":   peerBlockchain.Count - localBlockCount,
				"newBlockCount": blockchain.GetBlockCount(),
			})
			return
		}
	})

	// Audit endpoints
	router.GET("/audit", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"events": getAuditEvents(auditService),
		})
	})

	router.GET("/audit/security", func(c *gin.Context) {
		// Perform a full blockchain audit
		auditBlockchain(auditService)

		// Get critical events
		criticalEvents := getAuditEventsByLevel(auditService, audit.CriticalLevel)

		c.JSON(http.StatusOK, gin.H{
			"status": "completed",
			"issues": len(criticalEvents),
			"events": criticalEvents,
		})
	})

	// PAPRD Stablecoin endpoints
	// Read PAPRD ledger helper function
	readPAPRDLedger := func() (map[string]interface{}, error) {
		data, err := os.ReadFile("contracts/stablecoin/paprd-ledger.json")
		if err != nil {
			return nil, err
		}
		var ledger map[string]interface{}
		err = json.Unmarshal(data, &ledger)
		return ledger, err
	}

	// Write PAPRD ledger helper function
	writePAPRDLedger := func(ledger map[string]interface{}) error {
		data, err := json.MarshalIndent(ledger, "", "  ")
		if err != nil {
			return err
		}
		return os.WriteFile("contracts/stablecoin/paprd-ledger.json", data, 0644)
	}

	// Convert from 18 decimals
	fromDecimals := func(amount string) string {
		if amount == "" || amount == "0" {
			return "0"
		}
		// Simple conversion for display - divide by 10^18
		// For production, use big.Int for precision
		val, _ := strconv.ParseFloat(amount, 64)
		return fmt.Sprintf("%.0f", val/1e18)
	}

	// Convert to 18 decimals
	toDecimals := func(amount string) string {
		val, _ := strconv.ParseFloat(amount, 64)
		return fmt.Sprintf("%.0f", val*1e18)
	}

	// 📊 GET /paprd/info - Token information
	router.GET("/paprd/info", func(c *gin.Context) {
		ledger, err := readPAPRDLedger()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read PAPRD ledger"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"name":        ledger["token_name"],
			"symbol":      ledger["token_symbol"],
			"decimals":    ledger["token_decimals"],
			"totalSupply": fromDecimals(ledger["total_supply"].(string)),
			"owner":       ledger["owner"],
			"contract":    ledger["contract_id"],
			"paused":      ledger["paused"],
			"status":      "live",
		})
	})

	// 💰 GET /paprd/balance/:address - Get PAPRD balance
	router.GET("/paprd/balance/:address", func(c *gin.Context) {
		address := c.Param("address")

		ledger, err := readPAPRDLedger()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read PAPRD ledger"})
			return
		}

		balances := ledger["balances"].(map[string]interface{})
		balance := "0"
		if bal, exists := balances[address]; exists {
			balance = bal.(string)
		}

		readableBalance := fromDecimals(balance)

		c.JSON(http.StatusOK, gin.H{
			"address":    address,
			"balance":    readableBalance,
			"balanceWei": balance,
			"symbol":     "PAPRD",
		})
	})

	// 🔄 POST /paprd/transfer - Transfer PAPRD tokens
	router.POST("/paprd/transfer", func(c *gin.Context) {
		var request struct {
			From       string `json:"from" binding:"required"`
			To         string `json:"to" binding:"required"`
			Amount     string `json:"amount" binding:"required"`
			PrivateKey string `json:"privateKey"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ledger, err := readPAPRDLedger()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read PAPRD ledger"})
			return
		}

		// Check if contract is paused
		if ledger["paused"].(bool) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Contract is paused"})
			return
		}

		// Check blacklist
		blacklisted := ledger["blacklisted"].([]interface{})
		for _, addr := range blacklisted {
			if addr.(string) == request.From || addr.(string) == request.To {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Address is blacklisted"})
				return
			}
		}

		amountWei := toDecimals(request.Amount)
		balances := ledger["balances"].(map[string]interface{})

		fromBalanceStr := "0"
		if bal, exists := balances[request.From]; exists {
			fromBalanceStr = bal.(string)
		}

		fromBalance, _ := strconv.ParseFloat(fromBalanceStr, 64)
		amountFloat, _ := strconv.ParseFloat(amountWei, 64)

		// Check balance
		if fromBalance < amountFloat {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance"})
			return
		}

		// Perform transfer
		toBalanceStr := "0"
		if bal, exists := balances[request.To]; exists {
			toBalanceStr = bal.(string)
		}
		toBalance, _ := strconv.ParseFloat(toBalanceStr, 64)

		balances[request.From] = fmt.Sprintf("%.0f", fromBalance-amountFloat)
		balances[request.To] = fmt.Sprintf("%.0f", toBalance+amountFloat)

		// Record transaction
		transactions := ledger["transactions"].([]interface{})
		tx := map[string]interface{}{
			"id":        fmt.Sprintf("tx_%d", time.Now().UnixNano()),
			"type":      "transfer",
			"from":      request.From,
			"to":        request.To,
			"amount":    amountWei,
			"timestamp": time.Now().Format(time.RFC3339),
			"block":     time.Now().Unix(),
			"status":    "confirmed",
		}

		transactions = append(transactions, tx)
		ledger["transactions"] = transactions

		if err := writePAPRDLedger(ledger); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save transaction"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success":     true,
			"transaction": tx,
			"newBalance":  fromDecimals(balances[request.From].(string)),
		})

		// Log the transfer
		logAuditEvent(auditService, audit.InfoLevel, "PAPRDTransfer",
			fmt.Sprintf("PAPRD transfer: %s PAPRD from %s to %s", request.Amount, request.From, request.To), tx)
	})

	// 🪙 POST /paprd/mint - Mint PAPRD tokens (owner only)
	router.POST("/paprd/mint", func(c *gin.Context) {
		var request struct {
			To         string `json:"to" binding:"required"`
			Amount     string `json:"amount" binding:"required"`
			Caller     string `json:"caller" binding:"required"`
			PrivateKey string `json:"privateKey"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ledger, err := readPAPRDLedger()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read PAPRD ledger"})
			return
		}

		// Check if caller is owner
		if request.Caller != ledger["owner"].(string) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Only owner can mint tokens"})
			return
		}

		amountWei := toDecimals(request.Amount)
		balances := ledger["balances"].(map[string]interface{})

		toBalanceStr := "0"
		if bal, exists := balances[request.To]; exists {
			toBalanceStr = bal.(string)
		}
		toBalance, _ := strconv.ParseFloat(toBalanceStr, 64)
		amountFloat, _ := strconv.ParseFloat(amountWei, 64)

		balances[request.To] = fmt.Sprintf("%.0f", toBalance+amountFloat)

		// Update total supply
		totalSupply, _ := strconv.ParseFloat(ledger["total_supply"].(string), 64)
		ledger["total_supply"] = fmt.Sprintf("%.0f", totalSupply+amountFloat)

		// Record transaction
		transactions := ledger["transactions"].([]interface{})
		tx := map[string]interface{}{
			"id":        fmt.Sprintf("tx_%d", time.Now().UnixNano()),
			"type":      "mint",
			"from":      "0x0000000000000000000000000000000000000000",
			"to":        request.To,
			"amount":    amountWei,
			"timestamp": time.Now().Format(time.RFC3339),
			"block":     time.Now().Unix(),
			"status":    "confirmed",
		}

		transactions = append(transactions, tx)
		ledger["transactions"] = transactions

		if err := writePAPRDLedger(ledger); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save mint transaction"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success":     true,
			"transaction": tx,
			"newBalance":  fromDecimals(balances[request.To].(string)),
			"totalSupply": fromDecimals(ledger["total_supply"].(string)),
		})

		// Log the mint
		logAuditEvent(auditService, audit.InfoLevel, "PAPRDMint",
			fmt.Sprintf("PAPRD mint: %s PAPRD to %s by %s", request.Amount, request.To, request.Caller), tx)
	})

	// 📋 GET /paprd/transactions/:address - Get transaction history
	router.GET("/paprd/transactions/:address", func(c *gin.Context) {
		address := c.Param("address")

		ledger, err := readPAPRDLedger()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read PAPRD ledger"})
			return
		}

		transactions := ledger["transactions"].([]interface{})
		var userTxs []interface{}

		for _, txInterface := range transactions {
			tx := txInterface.(map[string]interface{})
			if tx["from"].(string) == address || tx["to"].(string) == address {
				userTxs = append(userTxs, tx)
			}
		}

		// Sort by timestamp (most recent first) and limit to 50
		if len(userTxs) > 50 {
			userTxs = userTxs[len(userTxs)-50:]
		}

		c.JSON(http.StatusOK, gin.H{
			"address":      address,
			"transactions": userTxs,
			"count":        len(userTxs),
		})
	})

	// 👛 GET /paprd/wallet/:address - Get wallet info
	router.GET("/paprd/wallet/:address", func(c *gin.Context) {
		address := c.Param("address")

		ledger, err := readPAPRDLedger()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read PAPRD ledger"})
			return
		}

		balances := ledger["balances"].(map[string]interface{})
		balance := "0"
		if bal, exists := balances[address]; exists {
			balance = bal.(string)
		}

		// Get BNM balance from main token system
		bnmBalance := binomToken.GetBalance(address)

		// Check if address is owner/minter
		isOwner := address == ledger["owner"].(string)
		isMinter := false
		minters := ledger["minters"].([]interface{})
		for _, minter := range minters {
			if minter.(string) == address {
				isMinter = true
				break
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"address":      address,
			"paprdBalance": fromDecimals(balance),
			"bnmBalance":   bnmBalance,
			"isOwner":      isOwner,
			"isMinter":     isMinter,
			"permissions": map[string]bool{
				"transfer": true,
				"mint":     isOwner || isMinter,
				"burn":     isOwner,
			},
		})
	})

	// DPoS endpoints
	router.POST("/delegates/register", func(c *gin.Context) {
		var request struct {
			Address    string  `json:"address" binding:"required"`
			Stake      float64 `json:"stake" binding:"required"`
			PrivateKey string  `json:"privateKey" binding:"required"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Verify wallet ownership
		senderWallet, err := wallet.ImportPrivateKey(request.PrivateKey)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid private key"})
			return
		}

		if senderWallet.Address != request.Address {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Private key does not match address"})
			return
		}

		// Check if address has enough balance for stake
		balance := binomToken.GetBalance(request.Address)
		if balance < request.Stake {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":    "Insufficient balance for stake",
				"balance":  balance,
				"required": request.Stake,
			})
			return
		}

		// Register delegate
		if err := dposConsensus.RegisterDelegate(request.Address, request.Stake); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": fmt.Sprintf("Delegate registered with %.2f BNM stake", request.Stake),
			"address": request.Address,
		})

		logAuditEvent(auditService, audit.InfoLevel, "DelegateRegistered",
			fmt.Sprintf("New delegate registered: %s with stake %.2f BNM", request.Address, request.Stake), nil)
	})

	router.POST("/delegates/vote", func(c *gin.Context) {
		var request struct {
			VoterAddress    string  `json:"voterAddress" binding:"required"`
			DelegateAddress string  `json:"delegateAddress" binding:"required"`
			Amount          float64 `json:"amount" binding:"required"`
			PrivateKey      string  `json:"privateKey" binding:"required"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Verify wallet ownership
		voterWallet, err := wallet.ImportPrivateKey(request.PrivateKey)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid private key"})
			return
		}

		if voterWallet.Address != request.VoterAddress {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Private key does not match voter address"})
			return
		}

		// Check if voter has enough balance
		balance := binomToken.GetBalance(request.VoterAddress)
		if balance < request.Amount {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":    "Insufficient balance for vote",
				"balance":  balance,
				"required": request.Amount,
			})
			return
		}

		// Vote for delegate
		if err := dposConsensus.VoteForDelegate(request.VoterAddress, request.DelegateAddress, request.Amount); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": fmt.Sprintf("Voted %.2f BNM for delegate %s", request.Amount, request.DelegateAddress),
		})

		logAuditEvent(auditService, audit.InfoLevel, "DelegateVote",
			fmt.Sprintf("Vote cast: %s voted %.2f BNM for delegate %s", request.VoterAddress, request.Amount, request.DelegateAddress), nil)
	})

	router.GET("/delegates", func(c *gin.Context) {
		delegates := dposConsensus.GetDelegates()

		c.JSON(http.StatusOK, gin.H{
			"delegates":    delegates,
			"count":        len(delegates),
			"maxDelegates": 21,
		})
	})

	router.GET("/delegates/:address", func(c *gin.Context) {
		address := c.Param("address")
		delegates := dposConsensus.GetDelegates()

		for _, delegate := range delegates {
			if delegate.Address == address {
				c.JSON(http.StatusOK, delegate)
				return
			}
		}

		c.JSON(http.StatusNotFound, gin.H{"error": "Delegate not found"})
	})

	// Register contract API routes
	contractAPI.RegisterRoutes(router)

	// Start the API server
	apiAddress := fmt.Sprintf(":%d", *apiPort)
	go func() {
		if err := router.Run(apiAddress); err != nil {
			log.Fatalf("Failed to start API server: %v", err)
		}
	}()

	fmt.Printf("Binomena blockchain node '%s' started\n", nodeName)
	fmt.Printf("API server running on http://localhost:%d\n", *apiPort)
	fmt.Printf("P2P node running on %s\n", p2pAddress)

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down Binomena node...")

	node.Stop()
	p2pNode.Stop()
	time.Sleep(time.Second)
	fmt.Println("Node stopped")
}
