# 🚀 Binomena Blockchain Platform

**The Next-Generation Blockchain for Decentralized Finance and Smart Contract Innovation**

Binomena is a high-performance blockchain platform built with Go, featuring advanced Delegated Proof of Stake (DPoS) consensus, WebAssembly smart contracts, and comprehensive DeFi capabilities including our flagship **PaperDollar USD Stablecoin**.

---

## 🌟 Key Features

### 🏛️ **Advanced Consensus Mechanism**
- **Delegated Proof of Stake (DPoS)** with 21 validator delegates
- **Dynamic reputation scoring** for validator selection
- **Fast block finality** and high throughput
- **Byzantine fault tolerance** with 2/3+1 majority consensus

### 💻 **WebAssembly Smart Contracts**
- **AssemblyScript/Rust/C++** contract development
- **Gas metering** for fair resource allocation
- **Memory management** with proper isolation
- **Event emission system** for real-time monitoring

### 💰 **Native BNM Token Economics**
- **Fixed supply**: 1 billion BNM tokens
- **Deflationary mechanism** through transaction fees
- **Staking rewards** for validators and delegators
- **Governance participation** for network decisions

### 🏦 **DeFi Infrastructure**
- **PaperDollar (USD)** - Production-ready stablecoin
- **Collateral management** with 150% collateralization ratio
- **Dual collateral support** (FIAT + BNM tokens)
- **Emergency controls** and compliance features

---

## 🎯 Live Deployments

### PaperDollar USD Stablecoin ✅
- **Contract ID**: `AdNe235e7e48b029a244c36daf0e9ef918e765b849bab4a97443ba820408b9bd`
- **Type**: USD-pegged stablecoin
- **Collateralization**: 150% minimum ratio
- **Features**: Minting, burning, blacklist/whitelist, emergency pause
- **Status**: ✅ **LIVE AND TESTED**

---

## 🏗️ Architecture Overview

### Core Components

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   P2P Network   │    │   DPoS Engine   │    │  WASM Runtime   │
│   (libp2p)      │◄──►│  (21 Delegates) │◄──►│  (Wasmer VM)    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  HTTP API/RPC   │    │  State Manager  │    │  Smart Contract │
│  (Gin Router)   │◄──►│  (PostgreSQL)   │◄──►│   Registry      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### Technology Stack

| Component | Technology | Purpose |
|-----------|------------|---------|
| **Backend** | Go 1.18+ | Core blockchain implementation |
| **Consensus** | DPoS | Fast, efficient block production |
| **Smart Contracts** | AssemblyScript/WASM | High-performance contract execution |
| **Database** | PostgreSQL | Persistent state storage |
| **Networking** | libp2p | Peer-to-peer communication |
| **Crypto** | ECDSA P-256 | Digital signatures and security |
| **API** | Gin Framework | RESTful web interface |

---

## 🚀 Quick Start

### Prerequisites

```bash
# Required software
- Go 1.18 or higher
- Node.js 16+ (for smart contracts)
- PostgreSQL 12+ (for state persistence)
- Git
```

### Installation

```bash
# Clone the repository
git clone https://github.com/igo-used/binomena.git
cd binomena

# Build the blockchain
go build -o app

# Setup smart contract environment
cd contracts
npm install
npm run build
cd ..

# Start the blockchain
./app
```

### Verify Installation

```bash
# Check blockchain status
curl http://localhost:8080/status

# Verify PaperDollar deployment
curl -X POST http://localhost:8080/contracts/call \
  -H "Content-Type: application/json" \
  -d '{
    "contractId": "AdNe235e7e48b029a244c36daf0e9ef918e765b849bab4a97443ba820408b9bd",
    "function": "getTotalSupply",
    "parameters": []
  }'
```

---

## 💼 Wallet Operations

### Create New Wallet

```bash
curl -X POST http://localhost:8080/wallet
```

**Response:**
```json
{
  "address": "AdNe6c3ce54e4371d056c7c566675ba16909eb2e9534",
  "privateKey": "your-private-key-here",
  "publicKey": "your-public-key-here"
}
```

### Check Balance

```bash
curl http://localhost:8080/balance/AdNe6c3ce54e4371d056c7c566675ba16909eb2e9534
```

### Send Transaction

```bash
curl -X POST http://localhost:8080/transaction \
  -H "Content-Type: application/json" \
  -d '{
    "from": "AdNe6c3ce54e4371d056c7c566675ba16909eb2e9534",
    "to": "AdNe987654321...",
    "amount": 100.0,
    "privateKey": "your-private-key"
  }'
```

---

## 🤖 Smart Contract Development

### PaperDollar Stablecoin Functions

```typescript
// Get total supply of stablecoins
export function getTotalSupply(): u64

// Check balance of an address
export function getBalance(address: string): u64

// Transfer tokens
export function transfer(to: string, amount: u64): boolean

// Mint new tokens (minter only)
export function mint(to: string, amount: u64): boolean

// Burn tokens
export function burn(amount: u64): boolean

// Add collateral (FIAT or BNM)
export function addCollateral(amount: u64, collateralType: u32): boolean

// Check collateralization ratio
export function getCollateralRatio(): u64

// Emergency controls
export function pause(): boolean
export function unpause(): boolean
```

### Deploy Your Own Contract

```bash
# Create contract directory
mkdir contracts/assembly/mytoken
cat > contracts/assembly/mytoken.ts << 'EOF'
export class MyToken {
  private balances: Map<string, u64> = new Map();
  private totalSupply: u64 = 1000000;

  constructor() {
    this.balances.set("owner", this.totalSupply);
  }

  getBalance(address: string): u64 {
    return this.balances.has(address) ? this.balances.get(address) : 0;
  }

  transfer(to: string, amount: u64): boolean {
    // Implementation here
    return true;
  }
}

// Export functions
let token = new MyToken();
export function getBalance(address: string): u64 {
  return token.getBalance(address);
}
EOF

# Build the contract
cd contracts && npm run build

# Deploy via API
curl -X POST http://localhost:8080/contracts/deploy \
  -H "Content-Type: application/json" \
  -d '{
    "owner": "AdNe6c3ce54e4371d056c7c566675ba16909eb2e9534",
    "name": "MyToken",
    "wasmPath": "contracts/build/release.wasm",
    "fee": 5.0,
    "privateKey": "your-private-key"
  }'
```

---

## 🏦 DeFi Operations

### PaperDollar Stablecoin Usage

#### 1. Add Collateral (Required before minting)

```bash
# Add FIAT collateral
curl -X POST http://localhost:8080/contracts/call \
  -H "Content-Type: application/json" \
  -d '{
    "contractId": "AdNe235e7e48b029a244c36daf0e9ef918e765b849bab4a97443ba820408b9bd",
    "function": "addCollateral",
    "parameters": [100000, 0],
    "caller": "AdNe6c3ce54e4371d056c7c566675ba16909eb2e9534",
    "privateKey": "your-private-key"
  }'
```

#### 2. Mint PaperDollar Tokens

```bash
# Mint 50,000 PaperDollar (requires 75,000+ collateral at 150% ratio)
curl -X POST http://localhost:8080/contracts/call \
  -H "Content-Type: application/json" \
  -d '{
    "contractId": "AdNe235e7e48b029a244c36daf0e9ef918e765b849bab4a97443ba820408b9bd",
    "function": "mint",
    "parameters": ["AdNe6c3ce54e4371d056c7c566675ba16909eb2e9534", 50000],
    "caller": "AdNe6c3ce54e4371d056c7c566675ba16909eb2e9534",
    "privateKey": "your-private-key"
  }'
```

#### 3. Transfer Stablecoins

```bash
curl -X POST http://localhost:8080/contracts/call \
  -H "Content-Type: application/json" \
  -d '{
    "contractId": "AdNe235e7e48b029a244c36daf0e9ef918e765b849bab4a97443ba820408b9bd",
    "function": "transfer",
    "parameters": ["AdNe987654321...", 1000],
    "caller": "AdNe6c3ce54e4371d056c7c566675ba16909eb2e9534",
    "privateKey": "your-private-key"
  }'
```

---

## 🎯 Network Statistics

### Token Distribution (1B BNM Total)
- **Founder**: 400M BNM (40%)
- **Treasury**: 400M BNM (40%)
- **Community**: 200M BNM (20%)

### Network Performance
- **Block Time**: ~3-5 seconds
- **TPS**: 1000+ transactions per second
- **Finality**: 2-3 blocks (~6-15 seconds)
- **Validators**: 21 active delegates

### Smart Contract Metrics
- **Gas Limit**: Configurable per contract
- **Languages**: AssemblyScript, Rust, C++
- **Runtime**: WebAssembly (WASM)
- **State**: Persistent PostgreSQL storage

---

## 🛠️ Advanced Features

### Consensus & Validation

```bash
# Check validator status
curl http://localhost:8080/validators

# View delegate information
curl http://localhost:8080/delegates

# Get network consensus state
curl http://localhost:8080/consensus/status
```

### Audit & Security

```bash
# Run security audit
curl http://localhost:8080/audit/run

# Check audit results
curl http://localhost:8080/audit/report

# Validate blockchain integrity
curl http://localhost:8080/validate/chain
```

### Network Administration

```bash
# Get node information
curl http://localhost:8080/node/info

# Check peer connections
curl http://localhost:8080/peers

# Monitor resource usage
curl http://localhost:8080/metrics
```

---

## 🧪 Testing & Development

### Run Test Suite

```bash
# Full test suite
./run_tests.sh

# Specific module tests
go test ./wallet -v
go test ./core -v
go test ./smartcontract -v
go test ./consensus -v

# Smart contract tests
cd contracts && npm test
```

### Development Tools

```bash
# Format code
go fmt ./...

# Lint code
golangci-lint run

# Build for production
go build -ldflags="-s -w" -o binomena-prod

# Cross-compile
GOOS=linux GOARCH=amd64 go build -o binomena-linux
```

---

## 📊 Monitoring & Analytics

### Blockchain Explorer

- **Blocks**: View all blocks with transactions
- **Transactions**: Search and analyze transaction history
- **Smart Contracts**: Browse deployed contracts
- **Validators**: Monitor validator performance

### API Endpoints for Analytics

```bash
# Get blockchain statistics
curl http://localhost:8080/stats

# Transaction volume metrics
curl http://localhost:8080/metrics/transactions

# Smart contract usage
curl http://localhost:8080/metrics/contracts

# Network health
curl http://localhost:8080/health
```

---

## 🔐 Security & Compliance

### Security Features

- ✅ **ECDSA P-256** cryptographic signatures
- ✅ **Byzantine fault tolerance** (BFT)
- ✅ **Smart contract sandboxing** via WASM
- ✅ **Gas metering** for DoS protection
- ✅ **Multi-signature** wallet support
- ✅ **Emergency pause** mechanisms

### Compliance Tools

- **Blacklist/Whitelist**: Address-based controls
- **KYC Integration**: Ready for compliance modules
- **Audit Trails**: Complete transaction history
- **Regulatory Reporting**: Exportable transaction data

---



## 🤝 Contributing

We welcome contributions from the community! Here's how to get started:

### Development Setup

```bash
# Fork the repository
git fork https://github.com/igo-used/binomena.git

# Clone your fork
git clone https://github.com/your-username/binomena.git
cd binomena

# Create feature branch
git checkout -b feature/amazing-feature

# Make changes and test
go test ./...
npm test --prefix contracts

# Commit and push
git commit -m "Add amazing feature"
git push origin feature/amazing-feature

# Open Pull Request
```

### Contribution Guidelines

1. **Code Quality**: Follow Go and AssemblyScript best practices
2. **Testing**: Add tests for new features
3. **Documentation**: Update README and code comments
4. **Security**: Consider security implications
5. **Performance**: Optimize for blockchain efficiency

---

## 📜 License

This project is licensed under the **Apache License 2.0** - see the [LICENSE](LICENSE) file for details.

---

## 📞 Contact

Website: https://binomena.com
Email: info@binomena.com
Founder Email: contact@Founder
X: Binomena Community

---

## 🎉 Success Stories

### PaperDollar Stablecoin Achievement

🏆 **Successfully deployed and tested USD-pegged stablecoin with:**
- ✅ 150% collateralization ratio
- ✅ Dual collateral support (FIAT + BNM)
- ✅ Professional-grade governance controls
- ✅ Emergency pause/unpause mechanisms
- ✅ Event emission for complete audit trails
- ✅ Production-ready on mainnet

**Contract ID**: `AdNe235e7e48b029a244c36daf0e9ef918e765b849bab4a97443ba820408b9bd`

---

*Built with ❤️ by the Binomena team for the decentralized future*
