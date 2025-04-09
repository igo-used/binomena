### Binomena Blockchain

## Overview

Binomena is a high-performance blockchain implemented in Go, featuring a custom Proof of Stake consensus mechanism called "NodeSwift" and a native token called "Binom" (BNM). This blockchain provides a complete ecosystem for secure transactions, wallet management, and token transfers with a deflationary economic model.

## Table of Contents

- [Features](#features)
- [Architecture](#architecture)
- [Installation](#installation)
- [Usage](#usage)
- [API Reference](#api-reference)
- [Development](#development)
- [License](#license)


## Features

### Core Blockchain

- **Custom Blockchain Implementation**: Built from scratch in Go
- **Block Structure**: Includes index, timestamp, transactions, validator signature, and hash
- **Transaction Validation**: Ensures all transactions are properly signed and valid
- **Genesis Block**: Automatic creation of genesis block on initialization


### NodeSwift Consensus

- **Proof of Stake (PoS)**: Energy-efficient consensus mechanism
- **Minimum Stake Requirement**: 1000 BNM minimum stake to participate in validation
- **Reputation System**: Validators gain or lose reputation based on performance
- **Fast Validation**: 5-second validation window for quick transaction confirmation
- **Weighted Selection**: Validators selected based on stake amount and reputation score


### Binom Token (BNM)

- **Maximum Supply**: Fixed at 1 billion BNM
- **Deflationary Mechanism**: 0.1% of each transaction amount is burned
- **Treasury Management**: Initial supply allocated to treasury for distribution
- **Balance Tracking**: Real-time balance tracking for all addresses


### Wallet System

- **Secure Key Generation**: ECDSA cryptographic key pairs
- **"AdNe" Prefixed Addresses**: Unique address format for easy identification
- **Transaction Signing**: Cryptographic signing of all transactions
- **Key Import/Export**: Support for wallet backup and recovery
- **Faucet Mechanism**: Request initial tokens for testing


### P2P Network

- **Peer Discovery**: Automatic peer discovery using mDNS
- **Transaction Broadcasting**: Propagation of transactions across the network
- **Block Synchronization**: Automatic blockchain synchronization between nodes
- **Wallet Announcements**: Network-wide wallet address announcements
- **Bootstrap Node Support**: Easy network initialization with bootstrap nodes


### Security Features

- **Audit System**: Comprehensive security monitoring
- **Event Logging**: Different severity levels for security events
- **Blockchain Validation**: Periodic validation of the entire blockchain
- **Hash Verification**: Ensures block integrity through hash verification
- **Signature Verification**: Cryptographic verification of all transactions


## Architecture

Binomena is built with a modular architecture consisting of the following components:

### Core Package

- `blockchain.go`: Implements the blockchain data structure and operations
- `node.go`: Manages the blockchain node and transaction processing
- `transaction.go`: Defines transaction structure and validation


### Consensus Package

- `nodeswift.go`: Implements the NodeSwift consensus mechanism


### Token Package

- `binom.go`: Implements the Binom token and its economic model


### Wallet Package

- `wallet.go`: Handles wallet creation, key management, and signing


### P2P Package

- `node.go`: Manages peer-to-peer communication and networking


### Audit Package

- `audit.go`: Provides security monitoring and validation


## Installation

### Prerequisites

- Go 1.16 or higher
- Git


### Building from Source

1. Clone the repository:


```shellscript
git clone https://github.com/igo-used/binomena.git
cd binomena
```

2. Build the project:


```shellscript
go build -o binomena
```

### Docker Deployment

1. Build the Docker image:


```shellscript
docker build -t binomena .
```

2. Run a container:


```shellscript
docker run -p 8080:8080 -p 9000:9000 binomena
```

## Usage

### Starting a Node

Run a standalone node:

```shellscript
./binomena --api-port 8080 --p2p-port 9000 --id "node-1"
```

Run a node with a bootstrap connection:

```shellscript
./binomena --api-port 8081 --p2p-port 9001 --bootstrap "/ip4/192.168.1.100/tcp/9000/p2p/QmNodeID" --id "node-2"
```

### Running Multiple Nodes

Use the provided script to run multiple nodes:

```shellscript
./run_nodes.sh
```

### Creating a Wallet

Create a new wallet via API:

```shellscript
curl -X POST http://localhost:8080/wallet
```

Import an existing wallet:

```shellscript
curl -X POST http://localhost:8080/wallet/import -d '{"privateKey":"your-private-key"}'
```

### Requesting Tokens from Faucet

```shellscript
curl -X POST http://localhost:8080/faucet -d '{"address":"AdNe123456789...","amount":1000}'
```

### Sending Transactions

```shellscript
curl -X POST http://localhost:8080/transaction -d '{
  "from": "AdNe123456789...",
  "to": "AdNe987654321...",
  "amount": 100,
  "privateKey": "your-private-key"
}'
```

### Checking Node Status

```shellscript
curl http://localhost:8080/status
```

### Viewing the Blockchain

```shellscript
curl http://localhost:8080/blocks
```

### Connecting to Peers

```shellscript
curl -X POST http://localhost:8080/peers -d '{"address":"/ip4/192.168.1.101/tcp/9001/p2p/QmNodeID"}'
```

## API Reference

### Node Management

- `GET /status`: Returns the current status of the node
- `GET /peers`: Lists all connected peers
- `POST /peers`: Connects to a new peer


### Wallet Operations

- `POST /wallet`: Creates a new wallet
- `POST /wallet/import`: Imports a wallet from a private key
- `POST /faucet`: Requests tokens from the faucet


### Transaction Operations

- `POST /transaction`: Submits a new transaction


### Blockchain Operations

- `GET /blocks`: Returns all blocks in the blockchain
- `GET /blocks/:index`: Returns a specific block by index
- `POST /sync`: Synchronizes the blockchain with a peer


### Security Operations

- `GET /audit`: Returns all audit events
- `GET /audit/security`: Performs a security audit of the blockchain


## Development

### Project Structure

```plaintext
binomena/
├── audit/
│   └── audit.go
├── consensus/
│   └── nodeswift.go
├── core/
│   ├── blockchain.go
│   ├── node.go
│   └── transaction.go
├── p2p/
│   └── node.go
├── test/
│   ├── blockchain_test.go
│   ├── binom_test.go
│   ├── nodeswift_test.go
│   ├── transaction_test.go
│   └── wallet_test.go
├── token/
│   └── binom.go
├── wallet/
│   └── wallet.go
├── main.go
├── run_nodes.sh
├── Dockerfile
├── go.mod
└── go.sum
```

### Running Tests

```shellscript
go test ./test/...
```

### Key Components

1. **Blockchain Core**:

1. Manages blocks and transactions
2. Validates blockchain integrity
3. Handles transaction processing



2. **NodeSwift Consensus**:

1. Implements validator selection
2. Manages reputation scores
3. Ensures fast transaction validation



3. **Binom Token**:

1. Manages token supply and balances
2. Implements deflationary mechanism
3. Handles token transfers



4. **Wallet System**:

1. Generates cryptographic keys
2. Creates and verifies signatures
3. Manages wallet addresses



5. **P2P Network**:

1. Handles peer discovery and connection
2. Broadcasts transactions and blocks
3. Synchronizes blockchain state



6. **Audit System**:

1. Monitors blockchain security
2. Logs security events
3. Performs periodic validations





## License

This project is licensed under the MIT License - see the LICENSE file for details.
