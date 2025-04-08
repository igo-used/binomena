# Binomena Blockchain

Binomena is a private blockchain implemented in Go with a custom "NodeSwift" consensus mechanism and a native token called "Binom" (BNM).

## Features

- **NodeSwift Consensus**: A custom Proof of Stake (PoS) consensus mechanism designed for security and fast transaction validation
- **Binom Token (BNM)**: Native token with a maximum supply of 1 billion
- **Deflationary Mechanism**: 0.1% burn ratio for each transaction
- **Fast Transaction Processing**: Optimized for quick validation and confirmation
- **P2P Wallet Generation**: Create and manage wallets with public/private key cryptography
- **"AdNe" Transaction Prefix**: All transactions and wallet addresses start with "AdNe" for easy identification
- **Security Audit System**: Built-in security monitoring and validation

## Architecture

The Binomena blockchain consists of the following components:

1. **Core**: Implements the basic blockchain structure, blocks, and transactions
2. **Consensus**: Implements the NodeSwift consensus mechanism
3. **Token**: Implements the Binom token with its economic model
4. **Wallet**: Manages cryptographic key pairs and addresses
5. **P2P**: Handles peer-to-peer communication and transaction propagation
6. **Audit**: Provides security monitoring and validation
7. **Node**: Manages the blockchain node, peer connections, and transaction processing

## Getting Started

### Prerequisites

- Go 1.16 or higher
- Git

### Installation

1. Clone the repository:
   \`\`\`
   git clone [https://github.com/yourusername/binomena.git](https://github.com/igo-used/binomena)
   cd binomena
   \`\`\`

2. Build the project:
   \`\`\`
   go build -o binomena
   \`\`\`

3. Run the node:
   \`\`\`
   ./binomena
   \`\`\`

## API Endpoints

The Binomena node exposes the following API endpoints:

- `GET /status`: Returns the current status of the node
- `POST /wallet`: Creates a new wallet
- `POST /wallet/import`: Imports a wallet from a private key
- `POST /transaction`: Submits a new transaction to the blockchain
- `GET /audit`: Returns all audit events
- `GET /audit/security`: Performs a security audit of the blockchain

## NodeSwift Consensus

NodeSwift is a custom Proof of Stake consensus mechanism that prioritizes security and transaction speed. Key features include:

- Minimum stake requirement for validators
- Reputation-based validator selection
- Short validation windows for faster transaction confirmation
- Weighted random selection based on stake and reputation

## Binom Token (BNM)

The Binom token (BNM) is the native token of the Binomena blockchain with the following characteristics:

- Maximum supply: 1 billion BNM
- Deflationary: 0.1% of each transaction is burned
- Used for transaction fees and staking

## Wallet System

The Binomena wallet system provides:

- Secure key generation using ECDSA
- "AdNe" prefixed addresses for easy identification
- Transaction signing and verification
- Import/export functionality

## Security Audit

The Binomena blockchain includes a comprehensive security audit system:

- Real-time monitoring of blockchain activity
- Detection of invalid blocks and transactions
- Periodic full blockchain validation
- Security event logging and reporting

## License

This project is licensed under the MIT License - see the LICENSE file for details.
