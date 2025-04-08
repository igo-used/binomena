package p2p

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/igo-used/binomena/core"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	discovery "github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"github.com/multiformats/go-multiaddr"
)

const (
	// Protocol IDs
	transactionProtocolID = "/binomena/tx/1.0.0"
	blockProtocolID       = "/binomena/block/1.0.0"
	walletDiscoveryID     = "/binomena/wallet/1.0.0"

	// Discovery service tag
	discoveryServiceTag = "binomena"
)

// P2PNode represents a P2P node in the Binomena network
type P2PNode struct {
	host         host.Host
	blockchain   *core.Blockchain
	knownPeers   map[peer.ID]peer.AddrInfo
	knownWallets map[string]string // address -> peer ID
	mu           sync.RWMutex
}

// discoveryNotifee gets notified when we find a new peer via mDNS discovery
type discoveryNotifee struct {
	node *P2PNode
}

// HandlePeerFound connects to peers discovered via mDNS
func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	log.Printf("Discovered new peer %s\n", pi.ID.String())
	err := n.node.host.Connect(context.Background(), pi)
	if err != nil {
		log.Printf("Error connecting to peer %s: %v\n", pi.ID.String(), err)
		return
	}

	// Add to known peers
	n.node.mu.Lock()
	n.node.knownPeers[pi.ID] = pi
	n.node.mu.Unlock()
}

// NewP2PNode creates a new P2P node
func NewP2PNode(blockchain *core.Blockchain, listenAddr string) (*P2PNode, error) {
	// Parse the multiaddress
	addr, err := multiaddr.NewMultiaddr(listenAddr)
	if err != nil {
		return nil, err
	}

	// Create a new libp2p host
	host, err := libp2p.New(
		libp2p.ListenAddrs(addr),
	)
	if err != nil {
		return nil, err
	}

	// Create the P2P node
	node := &P2PNode{
		host:         host,
		blockchain:   blockchain,
		knownPeers:   make(map[peer.ID]peer.AddrInfo),
		knownWallets: make(map[string]string),
	}

	// Set up protocol handlers
	host.SetStreamHandler(protocol.ID(transactionProtocolID), node.handleTransactionStream)
	host.SetStreamHandler(protocol.ID(blockProtocolID), node.handleBlockStream)
	host.SetStreamHandler(protocol.ID(walletDiscoveryID), node.handleWalletDiscoveryStream)

	// Setup local mDNS discovery
	notifee := &discoveryNotifee{node: node}
	disc := discovery.NewMdnsService(host, discoveryServiceTag, notifee)
	if disc == nil {
		return nil, fmt.Errorf("failed to initialize mDNS discovery service")
	}

	log.Printf("P2P node started with ID: %s", host.ID().String())
	log.Printf("Listening on: %s", host.Addrs()[0].String())

	return node, nil
}

// handleTransactionStream handles incoming transaction streams
func (n *P2PNode) handleTransactionStream(stream network.Stream) {
	defer stream.Close()

	// Decode the transaction
	var tx core.Transaction
	decoder := json.NewDecoder(stream)
	if err := decoder.Decode(&tx); err != nil {
		log.Printf("Error decoding transaction: %v", err)
		return
	}

	// Validate the transaction prefix
	if tx.ID[:4] != "AdNe" {
		log.Printf("Invalid transaction prefix: %s", tx.ID)
		return
	}

	// Add the transaction to the blockchain
	n.blockchain.AddTransaction(tx)

	log.Printf("Received transaction %s from peer %s", tx.ID, stream.Conn().RemotePeer().String())
}

// handleBlockStream handles incoming block streams
func (n *P2PNode) handleBlockStream(stream network.Stream) {
	defer stream.Close()

	// Decode the block
	var block core.Block
	decoder := json.NewDecoder(stream)
	if err := decoder.Decode(&block); err != nil {
		log.Printf("Error decoding block: %v", err)
		return
	}

	// Add the block to the blockchain
	if err := n.blockchain.AddBlock(block); err != nil {
		log.Printf("Error adding block: %v", err)
		return
	}

	log.Printf("Received block %d from peer %s", block.Index, stream.Conn().RemotePeer().String())
}

// handleWalletDiscoveryStream handles wallet discovery streams
func (n *P2PNode) handleWalletDiscoveryStream(stream network.Stream) {
	defer stream.Close()

	// Decode the wallet address
	var walletInfo struct {
		Address string `json:"address"`
	}
	decoder := json.NewDecoder(stream)
	if err := decoder.Decode(&walletInfo); err != nil {
		log.Printf("Error decoding wallet info: %v", err)
		return
	}

	// Store the wallet address and peer ID
	n.mu.Lock()
	n.knownWallets[walletInfo.Address] = stream.Conn().RemotePeer().String()
	n.mu.Unlock()

	log.Printf("Discovered wallet %s from peer %s", walletInfo.Address, stream.Conn().RemotePeer().String())
}

// BroadcastTransaction broadcasts a transaction to all known peers
func (n *P2PNode) BroadcastTransaction(tx core.Transaction) error {
	// Ensure transaction has the correct prefix
	if tx.ID[:4] != "AdNe" {
		return fmt.Errorf("transaction ID must start with 'AdNe'")
	}

	// Marshal the transaction to JSON
	txJSON, err := json.Marshal(tx)
	if err != nil {
		return err
	}

	// Broadcast to all known peers
	n.mu.RLock()
	peers := make([]peer.ID, 0, len(n.knownPeers))
	for id := range n.knownPeers {
		peers = append(peers, id)
	}
	n.mu.RUnlock()

	for _, peerID := range peers {
		// Open a stream to the peer
		stream, err := n.host.NewStream(context.Background(), peerID, protocol.ID(transactionProtocolID))
		if err != nil {
			log.Printf("Error opening stream to peer %s: %v", peerID.String(), err)
			continue
		}

		// Write the transaction to the stream
		_, err = stream.Write(txJSON)
		if err != nil {
			log.Printf("Error writing to stream: %v", err)
			stream.Close()
			continue
		}

		stream.Close()
	}

	return nil
}

// BroadcastBlock broadcasts a block to all known peers
func (n *P2PNode) BroadcastBlock(block core.Block) error {
	// Marshal the block to JSON
	blockJSON, err := json.Marshal(block)
	if err != nil {
		return err
	}

	// Broadcast to all known peers
	n.mu.RLock()
	peers := make([]peer.ID, 0, len(n.knownPeers))
	for id := range n.knownPeers {
		peers = append(peers, id)
	}
	n.mu.RUnlock()

	for _, peerID := range peers {
		// Open a stream to the peer
		stream, err := n.host.NewStream(context.Background(), peerID, protocol.ID(blockProtocolID))
		if err != nil {
			log.Printf("Error opening stream to peer %s: %v", peerID.String(), err)
			continue
		}

		// Write the block to the stream
		_, err = stream.Write(blockJSON)
		if err != nil {
			log.Printf("Error writing to stream: %v", err)
			stream.Close()
			continue
		}

		stream.Close()
	}

	return nil
}

// AnnounceWallet announces a wallet address to the network
func (n *P2PNode) AnnounceWallet(address string) error {
	// Ensure wallet address has the correct prefix
	if address[:4] != "AdNe" {
		return fmt.Errorf("wallet address must start with 'AdNe'")
	}

	// Create wallet info
	walletInfo := struct {
		Address string `json:"address"`
	}{
		Address: address,
	}

	// Marshal the wallet info to JSON
	walletJSON, err := json.Marshal(walletInfo)
	if err != nil {
		return err
	}

	// Broadcast to all known peers
	n.mu.RLock()
	peers := make([]peer.ID, 0, len(n.knownPeers))
	for id := range n.knownPeers {
		peers = append(peers, id)
	}
	n.mu.RUnlock()

	for _, peerID := range peers {
		// Open a stream to the peer
		stream, err := n.host.NewStream(context.Background(), peerID, protocol.ID(walletDiscoveryID))
		if err != nil {
			log.Printf("Error opening stream to peer %s: %v", peerID.String(), err)
			continue
		}

		// Write the wallet info to the stream
		_, err = stream.Write(walletJSON)
		if err != nil {
			log.Printf("Error writing to stream: %v", err)
			stream.Close()
			continue
		}

		stream.Close()
	}

	return nil
}

// FindWalletPeer finds the peer ID for a wallet address
func (n *P2PNode) FindWalletPeer(address string) (string, bool) {
	n.mu.RLock()
	defer n.mu.RUnlock()

	peerID, found := n.knownWallets[address]
	return peerID, found
}

// GetPeerCount returns the number of known peers
func (n *P2PNode) GetPeerCount() int {
	n.mu.RLock()
	defer n.mu.RUnlock()

	return len(n.knownPeers)
}

// GetWalletCount returns the number of known wallets
func (n *P2PNode) GetWalletCount() int {
	n.mu.RLock()
	defer n.mu.RUnlock()

	return len(n.knownWallets)
}

// Stop stops the P2P node
func (n *P2PNode) Stop() error {
	return n.host.Close()
}

// ConnectToPeer connects to a peer using its multiaddress
func (n *P2PNode) ConnectToPeer(peerAddr string) error {
	// Parse the multiaddress
	addr, err := multiaddr.NewMultiaddr(peerAddr)
	if err != nil {
		return fmt.Errorf("invalid peer address: %v", err)
	}

	// Extract the peer ID from the multiaddress
	info, err := peer.AddrInfoFromP2pAddr(addr)
	if err != nil {
		return fmt.Errorf("invalid peer info: %v", err)
	}

	// Connect to the peer
	if err := n.host.Connect(context.Background(), *info); err != nil {
		return fmt.Errorf("failed to connect to peer: %v", err)
	}

	// Add to known peers
	n.mu.Lock()
	n.knownPeers[info.ID] = *info
	n.mu.Unlock()

	log.Printf("Connected to peer: %s", info.ID.String())
	return nil
}

// GetPeers returns a list of connected peers
func (n *P2PNode) GetPeers() []string {
	n.mu.RLock()
	defer n.mu.RUnlock()

	peers := make([]string, 0, len(n.knownPeers))
	for _, info := range n.knownPeers {
		for _, addr := range info.Addrs {
			peers = append(peers, addr.String()+"/p2p/"+info.ID.String())
		}
	}

	return peers
}

// GetAddress returns the node's address
func (n *P2PNode) GetAddress() string {
	addrs := n.host.Addrs()
	if len(addrs) == 0 {
		return ""
	}

	return addrs[0].String() + "/p2p/" + n.host.ID().String()
}
