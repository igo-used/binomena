# Paper Dollar Stablecoin - Deployment Guide

## 🚀 Production Deployment

This guide covers deploying the Paper Dollar Stablecoin to a blockchain network with the founder as the administrator.

### 📋 Pre-Deployment Configuration

**Founder/Administrator Details:**
- **Address**: `AdNe6c3ce54e4371d056c7c566675ba16909eb2e9534`
- **Role**: Contract Owner & Administrator
- **Permissions**: Full contract control, minter management, pause/unpause, blacklist management

### 🔧 Deployment Steps

#### 1. Environment Setup

```bash
# Clone and navigate to the contract
cd contracts/stablecoin

# Make scripts executable
chmod +x deploy.sh initialize-contract.sh build.sh setup.sh
```

#### 2. Build and Test

```bash
# Run the deployment script (builds, tests, and creates WASM)
./deploy.sh
```

This script will:
- ✅ Build the Rust contract
- ✅ Run all tests (15 tests should pass)
- ✅ Generate WebAssembly binaries
- ✅ Create deployment artifacts in `pkg/` directory

#### 3. Contract Initialization

```bash
# Initialize the contract with founder as owner
./initialize-contract.sh
```

This will:
- ✅ Set founder address as contract owner
- ✅ Initialize with 150% collateral ratio
- ✅ Verify contract functionality

### 📦 Deployment Artifacts

After successful deployment, you'll have:

```
contracts/stablecoin/
├── pkg/                    # Web deployment package
├── pkg-node/              # Node.js deployment package
├── target/release/        # Optimized Rust binary
└── src/                   # Source code
```

### 🔐 Security Configuration

#### Critical Security Notes:

1. **Private Key Security**: The founder's private key is configured for testing but should be managed securely in production
2. **Environment Variables**: Never commit private keys to version control
3. **Access Control**: Only the founder address can perform administrative functions

#### Founder Permissions:

- ✅ Add/remove minters
- ✅ Pause/unpause contract
- ✅ Blacklist/unblacklist addresses
- ✅ Set collateral ratios
- ✅ Transfer ownership
- ✅ Configure BINOM token address

### 🌐 Blockchain Deployment

#### For NEAR Protocol:
```bash
# Deploy to NEAR testnet
near deploy --accountId your-contract.testnet --wasmFile pkg/paper_dollar_stablecoin_bg.wasm

# Initialize contract
near call your-contract.testnet new_with_founder --accountId AdNe6c3ce54e4371d056c7c566675ba16909eb2e9534
```

#### For Ethereum/EVM:
```bash
# Use the generated WASM with your preferred deployment tool
# The contract is compatible with WASM-based EVM solutions
```

#### For Other Blockchains:
The contract is designed to be blockchain-agnostic. Adapt the deployment process according to your target blockchain's requirements.

### 📊 Post-Deployment Verification

#### 1. Verify Owner
```javascript
// Check that founder is set as owner
const owner = await contract.get_owner();
console.log("Owner:", owner); // Should be: AdNe6c3ce54e4371d056c7c566675ba16909eb2e9534
```

#### 2. Verify Initial State
```javascript
// Check initial configuration
const totalSupply = await contract.get_total_supply();
const collateralRatio = await contract.get_collateral_ratio();
const isPaused = await contract.is_paused();

console.log("Total Supply:", totalSupply);      // Should be: 0
console.log("Collateral Ratio:", collateralRatio); // Should be: 150
console.log("Is Paused:", isPaused);            // Should be: false
```

#### 3. Test Administrative Functions
```javascript
// Add a minter (only founder can do this)
await contract.add_minter("minter_address");

// Verify minter was added
const isMinter = await contract.is_minter("minter_address");
console.log("Is Minter:", isMinter); // Should be: true
```

### 🔄 Operational Procedures

#### Adding Minters
```javascript
// Only founder can add minters
await contract.add_minter("new_minter_address");
```

#### Managing Collateral
```javascript
// Users can add collateral
await contract.add_collateral(1000, 0); // 1000 units of Fiat collateral

// Check collateral balance
const balance = await contract.get_collateral_balance("user_address", 0);
```

#### Minting Tokens
```javascript
// Only authorized minters can mint
await contract.mint("user_address", 500); // Mint 500 tokens
```

#### Emergency Controls
```javascript
// Founder can pause contract in emergency
await contract.pause();

// Founder can unpause when safe
await contract.unpause();
```

### 📈 Monitoring and Maintenance

#### Key Metrics to Monitor:
- Total supply vs. total collateral
- Collateral ratio compliance
- Active minters
- Blacklisted addresses
- Contract pause status

#### Regular Maintenance:
- Monitor collateral ratios
- Review minter permissions
- Update blacklist as needed
- Backup contract state

### 🆘 Emergency Procedures

#### If Compromise Detected:
1. **Immediate**: Call `pause()` to halt all operations
2. **Assess**: Determine scope of compromise
3. **Communicate**: Notify stakeholders
4. **Recover**: Use `unpause()` when safe

#### If Founder Key Compromised:
1. **Immediate**: Transfer ownership to secure address
2. **Revoke**: Remove compromised minters
3. **Audit**: Review all recent transactions

### 📞 Support and Documentation

- **Technical Documentation**: See `README.md`
- **Migration Guide**: See `MIGRATION.md`
- **Quick Start**: See `QUICKSTART.md`
- **Source Code**: `src/` directory

### ✅ Deployment Checklist

- [ ] Environment setup complete
- [ ] All tests passing (15/15)
- [ ] WebAssembly build successful
- [ ] Founder address configured as owner
- [ ] Initial collateral ratio set (150%)
- [ ] Contract deployed to blockchain
- [ ] Post-deployment verification complete
- [ ] Administrative functions tested
- [ ] Monitoring systems in place
- [ ] Emergency procedures documented

---

**🔐 Security Reminder**: The founder address `AdNe6c3ce54e4371d056c7c566675ba16909eb2e9534` has full administrative control. Ensure the corresponding private key is stored securely and never exposed in public repositories or logs. 