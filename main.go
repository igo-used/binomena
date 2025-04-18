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
	"github.com/igo-used/binomena/p2p"
	"github.com/igo-used/binomena/smartcontract"
	"github.com/igo-used/binomena/token"
	"github.com/igo-used/binomena/wallet"
)

func main() {
	// Parse command line flags
	apiPort := flag.Int("api-port", 8080, "API server port")
	p2pPort := flag.Int("p2p-port", 9000, "P2P server port")
	bootstrapNode := flag.String("bootstrap", "", "Bootstrap node address (optional)")
	nodeID := flag.String("id", "", "Node identifier (optional)")
	flag.Parse()

	// Set node identifier
	nodeName := *nodeID
	if nodeName == "" {
		nodeName = fmt.Sprintf("node-%d", *p2pPort)
	}

	// Initialize the blockchain
	blockchain := core.NewBlockchain()

	// Get data directory from environment variable
	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = "./data" // Fallback to local directory if env var not set
	}

	log.Printf("Using data directory: %s", dataDir)

	// Load blockchain from disk if available
	if err := blockchain.LoadChain(dataDir); err != nil {
		log.Printf("Warning: Failed to load blockchain: %v", err)
	} else {
		log.Printf("Successfully loaded blockchain from disk")
	}

	// Set up periodic saving of blockchain
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			if err := blockchain.SaveChain(dataDir); err != nil {
				log.Printf("Warning: Failed to save blockchain: %v", err)
			} else {
				log.Printf("Blockchain saved successfully")
			}
		}
	}()

	// Initialize the NodeSwift consensus mechanism
	nodeSwift := consensus.NewNodeSwift()

	// Initialize the Binom token
	binomToken := token.NewBinomToken()

	// Load token balances from disk if available
	if err := binomToken.LoadBalances(dataDir); err != nil {
		log.Printf("Warning: Failed to load token balances: %v", err)
	} else {
		log.Printf("Successfully loaded token balances from disk")
	}

	// Set up periodic saving of token balances
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			if err := binomToken.SaveBalances(dataDir); err != nil {
				log.Printf("Warning: Failed to save token balances: %v", err)
			} else {
				log.Printf("Token balances saved successfully")
			}
		}
	}()

	// Initialize the smart contract system
	contractStorage, err := smartcontract.NewContractStorage("./contracts")
	if err != nil {
		log.Fatalf("Failed to initialize contract storage: %v", err)
	}

	contractState, err := smartcontract.NewContractState("./contracts")
	if err != nil {
		log.Fatalf("Failed to initialize contract state: %v", err)
	}

	wasmVM, err := smartcontract.NewWasmVM(binomToken, blockchain)
	if err != nil {
		log.Fatalf("Failed to initialize WASM VM: %v", err)
	}

	// Load existing contracts
	contracts, err := contractStorage.LoadAllContracts()
	if err != nil {
		log.Printf("Warning: Failed to load contracts: %v", err)
	} else {
		for _, contract := range contracts {
			// Add contract to VM
			wasmVM.AddContract(contract)
		}
		log.Printf("Loaded %d contracts", len(contracts))
	}

	// Create contract API
	contractAPI := smartcontract.NewContractAPI(wasmVM, contractStorage, contractState, binomToken)

	// Initialize the audit service
	auditService := audit.NewAuditService(blockchain)

	// Create a new node
	node := core.NewNode(blockchain, nodeSwift, binomToken)

	// Start the P2P network
	p2pAddress := fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", *p2pPort)
	p2pNode, err := p2p.NewP2PNode(blockchain, p2pAddress)
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
		auditService.LogEvent(
			audit.InfoLevel,
			"WalletCreated",
			fmt.Sprintf("New wallet created with address %s", newWallet.Address),
			nil,
		)
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
		auditService.LogEvent(
			audit.InfoLevel,
			"WalletImported",
			fmt.Sprintf("Wallet imported with address %s", importedWallet.Address),
			nil,
		)
	})

	// NEW ENDPOINT: Get wallet balance
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

		// Verify admin key (use environment variable with fallback)
		adminKey := os.Getenv("BINOMENA_ADMIN_KEY")
		if adminKey == "" {
			adminKey = "binomena-founder-key-2025" // Fallback key if env var is not set
		}

		if request.AdminKey != adminKey {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
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

		// Save token balances immediately after distribution
		if err := binomToken.SaveBalances(dataDir); err != nil {
			log.Printf("Warning: Failed to save token balances after distribution: %v", err)
		} else {
			log.Printf("Token balances saved successfully after distribution")
		}

		// Log the distribution
		auditService.LogEvent(
			audit.InfoLevel,
			"InitialTokenDistribution",
			fmt.Sprintf("Distributed tokens: %f to founder, %f to treasury, %f to community",
				founderAmount, treasuryAmount, communityAmount),
			map[string]interface{}{
				"founder":   request.FounderAddress,
				"treasury":  request.TreasuryAddress,
				"community": request.CommunityAddress,
			},
		)

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

	// Faucet endpoint to request initial tokens
	router.POST("/faucet", func(c *gin.Context) {
		var request struct {
			Address  string  `json:"address"`
			Amount   float64 `json:"amount"`
			AdminKey string  `json:"adminKey"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate address
		if request.Address[:4] != "AdNe" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address format"})
			return
		}

		// Get admin key from environment variable with fallback
		adminKey := os.Getenv("BINOMENA_ADMIN_KEY")
		if adminKey == "" {
			adminKey = "binomena-founder-key-2025" // Fallback key if env var is not set
		}

		// Check if admin key is valid
		if request.AdminKey != adminKey {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Unauthorized. Tokens are not available for free distribution.",
			})
			return
		}

		// Validate amount
		if request.Amount <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Amount must be greater than 0"})
			return
		}

		// Transfer tokens from treasury to the address
		err := binomToken.Transfer("treasury", request.Address, request.Amount)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Save token balances immediately after faucet distribution
		if err := binomToken.SaveBalances(dataDir); err != nil {
			log.Printf("Warning: Failed to save token balances after faucet: %v", err)
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": fmt.Sprintf("Transferred %f BNM to %s", request.Amount, request.Address),
			"balance": binomToken.GetBalance(request.Address),
		})

		// Log faucet request
		auditService.LogEvent(
			audit.InfoLevel,
			"FaucetRequest",
			fmt.Sprintf("Transferred %f BNM to %s", request.Amount, request.Address),
			nil,
		)
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

		// Check balance
		balance := binomToken.GetBalance(request.From)
		if balance < request.Amount {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":    "insufficient balance",
				"balance":  balance,
				"required": request.Amount,
			})
			return
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

		// Save token balances and blockchain after transaction
		if err := binomToken.SaveBalances(dataDir); err != nil {
			log.Printf("Warning: Failed to save token balances after transaction: %v", err)
		}

		if err := blockchain.SaveChain(dataDir); err != nil {
			log.Printf("Warning: Failed to save blockchain after transaction: %v", err)
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "transaction submitted",
			"txId":   tx.ID,
			"node":   nodeName,
		})

		// Log transaction
		auditService.LogEvent(
			audit.InfoLevel,
			"TransactionSubmitted",
			fmt.Sprintf("Transaction %s submitted from %s to %s for %f BNM", tx.ID, tx.From, tx.To, tx.Amount),
			tx,
		)
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

			// Save the updated blockchain and token balances
			if err := blockchain.SaveChain(dataDir); err != nil {
				log.Printf("Warning: Failed to save blockchain after sync: %v", err)
			}

			if err := binomToken.SaveBalances(dataDir); err != nil {
				log.Printf("Warning: Failed to save token balances after sync: %v", err)
			}

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

			// Save the updated blockchain and token balances
			if err := blockchain.SaveChain(dataDir); err != nil {
				log.Printf("Warning: Failed to save blockchain after sync: %v", err)
			}

			if err := binomToken.SaveBalances(dataDir); err != nil {
				log.Printf("Warning: Failed to save token balances after sync: %v", err)
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
			"events": auditService.GetEvents(),
		})
	})

	router.GET("/audit/security", func(c *gin.Context) {
		// Perform a full blockchain audit
		auditService.AuditBlockchain()

		// Get critical events
		criticalEvents := auditService.GetEventsByLevel(audit.CriticalLevel)

		c.JSON(http.StatusOK, gin.H{
			"status": "completed",
			"issues": len(criticalEvents),
			"events": criticalEvents,
		})
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

	// Save token balances before shutdown
	if err := binomToken.SaveBalances(dataDir); err != nil {
		log.Printf("Warning: Failed to save token balances on shutdown: %v", err)
	} else {
		log.Printf("Token balances saved successfully on shutdown")
	}

	// Save blockchain before shutdown
	if err := blockchain.SaveChain(dataDir); err != nil {
		log.Printf("Warning: Failed to save blockchain on shutdown: %v", err)
	} else {
		log.Printf("Blockchain saved successfully on shutdown")
	}

	node.Stop()
	p2pNode.Stop()
	time.Sleep(time.Second)
	fmt.Println("Node stopped")
}
