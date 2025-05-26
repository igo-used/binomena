# 🚀 PAPRD Token - Blockchain Deployment Guide

## 💰 **Token Information**
- **Name**: Paper Dollar
- **Symbol**: **PAPRD** 
- **Decimals**: 18
- **Type**: Collateral-backed stablecoin
- **Founder/Admin**: `AdNe6c3ce54e4371d056c7c566675ba16909eb2e9534`

---

## 🌐 **Blockchain Deployment Options**

### **Option 1: NEAR Protocol (Recommended)**

#### **Why NEAR?**
- ✅ Native WebAssembly support
- ✅ Low transaction fees
- ✅ Fast finality (1-2 seconds)
- ✅ Developer-friendly
- ✅ Built-in wallet integration

#### **Deploy to NEAR Testnet:**
```bash
# 1. Install NEAR CLI
npm install -g near-cli

# 2. Create NEAR account (if you don't have one)
near login

# 3. Create contract account
near create-account paprd-stablecoin.your-account.testnet --masterAccount your-account.testnet

# 4. Deploy the contract
near deploy --accountId paprd-stablecoin.your-account.testnet --wasmFile pkg/paper_dollar_stablecoin_bg.wasm

# 5. Initialize with founder
near call paprd-stablecoin.your-account.testnet new_with_founder --accountId AdNe6c3ce54e4371d056c7c566675ba16909eb2e9534
```

#### **Deploy to NEAR Mainnet:**
```bash
# Same process but use mainnet accounts
near create-account paprd-stablecoin.near --masterAccount your-account.near
near deploy --accountId paprd-stablecoin.near --wasmFile pkg/paper_dollar_stablecoin_bg.wasm
```

### **Option 2: Ethereum/Polygon (EVM Compatible)**

#### **Using Hardhat:**
```javascript
// hardhat.config.js
require("@nomiclabs/hardhat-waffle");

module.exports = {
  solidity: "0.8.19",
  networks: {
    polygon: {
      url: "https://polygon-rpc.com",
      accounts: ["YOUR_PRIVATE_KEY_HERE"] // Use your founder's private key
    }
  }
};
```

```javascript
// deploy.js
const { ethers } = require("hardhat");

async function main() {
  // Deploy using a WASM wrapper contract
  const PaprdToken = await ethers.getContractFactory("WasmWrapper");
  const paprd = await PaprdToken.deploy(
    "Paper Dollar",
    "PAPRD",
    18,
    "AdNe6c3ce54e4371d056c7c566675ba16909eb2e9534"
  );
  
  console.log("PAPRD deployed to:", paprd.address);
}
```

### **Option 3: Polkadot/Substrate**

```bash
# For Substrate-based chains with WASM support
substrate-contracts-node --dev --tmp

# Deploy using cargo-contract
cargo contract deploy --suri //Alice --constructor new_with_founder
```

---

## 💼 **Getting PAPRD Tokens in Your Wallet**

### **🔑 Can You Use Your Existing Wallet?**

**YES!** You can use your existing wallet depending on the blockchain:

#### **For NEAR Protocol:**
- **NEAR Wallet** (wallet.near.org)
- **MyNearWallet** 
- **Meteor Wallet**
- **Sender Wallet**

#### **For Ethereum/Polygon:**
- **MetaMask** ✅
- **Trust Wallet** ✅  
- **Coinbase Wallet** ✅
- **WalletConnect** compatible wallets ✅

#### **For Polkadot:**
- **Polkadot.js Extension** ✅
- **Talisman Wallet** ✅
- **SubWallet** ✅

### **🎯 How to Get PAPRD Tokens**

#### **Method 1: Direct Minting (If you're authorized)**
```javascript
// 1. Founder adds you as a minter
await contract.add_minter("your_wallet_address");

// 2. You deposit collateral first
await contract.add_collateral(1500, 0); // $1,500 fiat collateral

// 3. Mint tokens for yourself
await contract.mint("your_wallet_address", 1000); // Get 1000 PAPRD
```

#### **Method 2: Get Tokens from Founder**
```javascript
// Founder can mint tokens and send them to you
// 1. Founder deposits collateral
await contract.add_collateral(15000, 0); // $15,000 collateral

// 2. Founder mints tokens
await contract.mint("founder_address", 10000); // 10,000 PAPRD

// 3. Founder transfers to you
await contract.transfer("your_wallet_address", 1000); // You get 1000 PAPRD
```

#### **Method 3: Purchase from Others** (Future)
```javascript
// Once trading is established, you can buy from others
await contract.transfer("seller_address", purchase_amount);
```

---

## 📱 **Step-by-Step: Getting PAPRD in MetaMask**

### **For Ethereum/Polygon Deployment:**

#### **1. Add PAPRD Token to MetaMask:**
```javascript
// Contract details to add in MetaMask
Token Contract Address: "0x..." // (After deployment)
Token Symbol: PAPRD
Token Decimals: 18
```

#### **2. Add Token Manually:**
1. Open MetaMask
2. Click "Import tokens"
3. Select "Custom Token"
4. Enter contract address, PAPRD, 18
5. Click "Add Custom Token"

#### **3. Request Tokens:**
Contact the founder to mint tokens for you, or deposit collateral and request minting.

---

## 🔧 **Complete Deployment Example (NEAR)**

### **Step 1: Deploy Contract**
```bash
# Deploy to NEAR testnet
near deploy --accountId paprd.testnet --wasmFile pkg/paper_dollar_stablecoin_bg.wasm
```

### **Step 2: Initialize Contract**
```bash
# Initialize with founder
near call paprd.testnet new_with_founder --accountId paprd.testnet
```

### **Step 3: Verify Deployment**
```bash
# Check token info
near view paprd.testnet get_name
near view paprd.testnet get_symbol  # Should return "PAPRD"
near view paprd.testnet get_decimals
near view paprd.testnet get_owner   # Should be founder address
```

### **Step 4: Setup Minters**
```bash
# Founder adds authorized minters
near call paprd.testnet add_minter '{"address": "authorized-minter.testnet"}' --accountId AdNe6c3ce54e4371d056c7c566675ba16909eb2e9534
```

### **Step 5: First Token Creation**
```bash
# User deposits collateral
near call paprd.testnet add_collateral '{"amount": 1500, "collateral_type": 0}' --accountId user.testnet

# Minter creates tokens
near call paprd.testnet mint '{"to": "user.testnet", "amount": 1000}' --accountId authorized-minter.testnet
```

---

## 🏦 **Token Economics**

### **Initial Setup:**
- **Total Supply**: 0 PAPRD (grows with minting)
- **Collateral Ratio**: 150% (adjustable by founder)
- **Supported Collateral**: 
  - Fiat currency backing
  - BINOM tokens
- **Maximum Supply**: Unlimited (controlled by collateral)

### **Token Distribution:**
```
Founder: Controls all initial minting and management
Authorized Minters: Can create tokens for users with collateral
Users: Can hold, transfer, and burn tokens
```

---

## 📋 **Post-Deployment Checklist**

### **For Founder:**
- [ ] Deploy contract to chosen blockchain
- [ ] Verify founder address is owner
- [ ] Add initial authorized minters
- [ ] Set up collateral acceptance process
- [ ] Configure initial token distribution

### **For Users:**
- [ ] Choose compatible wallet for blockchain
- [ ] Add PAPRD token to wallet (contract address needed)
- [ ] Contact founder or authorized minter for initial tokens
- [ ] Understand collateral requirements if becoming a minter

---

## 🔐 **Security Considerations**

### **For Founder:**
- **Never share your private key**: `508f988eaa468e8092ea01cbb484d474b705b45047f3bcae748749722f637fdf`
- **Use hardware wallet** for mainnet deployments
- **Regularly backup** contract state and keys
- **Monitor contract** for unusual activity

### **For Users:**
- **Verify contract address** before adding to wallet
- **Start with small amounts** for testing
- **Understand collateral requirements** before minting
- **Keep wallet software updated**

---

## 📞 **Getting Help**

### **Contract Information:**
- **Name**: Paper Dollar
- **Symbol**: PAPRD  
- **Decimals**: 18
- **Admin**: AdNe6c3ce54e4371d056c7c566675ba16909eb2e9534

### **Support:**
- Technical docs: `README.md`
- Deployment guide: `DEPLOYMENT.md`
- Migration guide: `MIGRATION.md`

---

## 🎯 **Quick Start Commands**

```bash
# Test the contract locally
./deploy.sh

# Initialize with founder
./initialize-contract.sh

# Deploy to NEAR testnet
near deploy --accountId your-contract.testnet --wasmFile pkg/paper_dollar_stablecoin_bg.wasm

# Check contract info
near view your-contract.testnet get_symbol  # Returns: "PAPRD"
```

Your **PAPRD stablecoin** is now ready for blockchain deployment! 🚀 