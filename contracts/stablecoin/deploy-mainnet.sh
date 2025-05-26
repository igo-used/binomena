#!/bin/bash

# PAPRD Stablecoin MAINNET Deployment Script for Binomena Blockchain
# Deploys to production blockchain on Render: https://binomena-node.onrender.com

set -e

echo "🚀 PAPRD Stablecoin MAINNET Deployment"
echo "======================================"
echo "🌐 Target: https://binomena-node.onrender.com"
echo ""

# Configuration
FOUNDER_ADDRESS="AdNe6c3ce54e4371d056c7c566675ba16909eb2e9534"
FOUNDER_PRIVATE_KEY="508f988eaa468e8092ea01cbb484d474b705b45047f3bcae748749722f637fdf"
CONTRACT_NAME="Paper Dollar Stablecoin (PAPRD)"
DEPLOYMENT_FEE="1.0"  # 1 BNM for deployment
MAINNET_API="https://binomena-node.onrender.com"

echo "📋 MAINNET Deployment Configuration:"
echo "Founder Address: $FOUNDER_ADDRESS"
echo "Contract Name: $CONTRACT_NAME"
echo "Deployment Fee: $DEPLOYMENT_FEE BNM"
echo "Mainnet API: $MAINNET_API"
echo ""

# Pre-deployment checks
echo "🔍 Pre-deployment verification..."

# Check mainnet blockchain status
echo "Checking mainnet blockchain status..."
BLOCKCHAIN_STATUS=$(curl -s "$MAINNET_API/status" || echo '{"status": "error"}')
if echo "$BLOCKCHAIN_STATUS" | grep -q '"status":"running"'; then
    echo "✅ Mainnet blockchain is running"
    echo "Status: $(echo $BLOCKCHAIN_STATUS | grep -o '"status":"[^"]*"' | cut -d'"' -f4)"
    echo "Blocks: $(echo $BLOCKCHAIN_STATUS | grep -o '"blocks":[0-9]*' | cut -d':' -f2)"
    echo "Node ID: $(echo $BLOCKCHAIN_STATUS | grep -o '"nodeId":"[^"]*"' | cut -d'"' -f4)"
else
    echo "❌ Mainnet blockchain is not responding!"
    echo "Response: $BLOCKCHAIN_STATUS"
    exit 1
fi

# Check for existing contracts
echo ""
echo "Checking existing contracts on mainnet..."
EXISTING_CONTRACTS=$(curl -s "$MAINNET_API/contracts" || echo '{"count": -1}')
CONTRACT_COUNT=$(echo $EXISTING_CONTRACTS | grep -o '"count":[0-9]*' | cut -d':' -f2 || echo "-1")

if [ "$CONTRACT_COUNT" = "0" ]; then
    echo "✅ Mainnet is clean - no existing contracts"
else
    echo "⚠️ Found $CONTRACT_COUNT existing contracts on mainnet"
    echo "Existing contracts: $EXISTING_CONTRACTS"
    read -p "Continue with deployment? (y/n): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Deployment cancelled by user"
        exit 1
    fi
fi

# Check founder balance on mainnet
echo ""
echo "💰 Checking founder balance on mainnet..."
BALANCE_RESPONSE=$(curl -s "$MAINNET_API/balance/$FOUNDER_ADDRESS" || echo '{"balance": 0}')
FOUNDER_BALANCE=$(echo $BALANCE_RESPONSE | grep -o '"balance":[0-9.]*' | cut -d':' -f2 || echo "0")

echo "Founder balance on mainnet: $FOUNDER_BALANCE BNM"

if (( $(echo "$FOUNDER_BALANCE < $DEPLOYMENT_FEE" | bc -l) )); then
    echo "❌ Insufficient balance for deployment!"
    echo "Required: $DEPLOYMENT_FEE BNM"
    echo "Available: $FOUNDER_BALANCE BNM"
    exit 1
fi

echo "✅ Sufficient balance for deployment"

# Verify contract build
echo ""
echo "🔨 Verifying PAPRD contract build..."
if [ ! -f "pkg/paper_dollar_stablecoin_bg.wasm" ]; then
    echo "Contract not built. Building now..."
    ./deploy.sh
fi

# Get WASM file info
WASM_SIZE=$(wc -c < pkg/paper_dollar_stablecoin_bg.wasm)
echo "✅ Contract WASM found"
echo "WASM Size: $WASM_SIZE bytes"

# Convert WASM to base64
echo ""
echo "📦 Converting WASM to base64 for mainnet deployment..."
WASM_BASE64=$(base64 -i pkg/paper_dollar_stablecoin_bg.wasm)
BASE64_SIZE=${#WASM_BASE64}
echo "Base64 Size: $BASE64_SIZE characters"

# Final confirmation
echo ""
echo "🎯 MAINNET DEPLOYMENT SUMMARY:"
echo "=============================="
echo "Target Blockchain: BINOMENA MAINNET"
echo "API Endpoint: $MAINNET_API"
echo "Contract: $CONTRACT_NAME"
echo "Owner: $FOUNDER_ADDRESS"
echo "Deployment Fee: $DEPLOYMENT_FEE BNM"
echo "Founder Balance: $FOUNDER_BALANCE BNM"
echo "Contract Size: $WASM_SIZE bytes"
echo ""

read -p "🚨 Deploy PAPRD to MAINNET? This is PRODUCTION! (y/n): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Mainnet deployment cancelled"
    exit 1
fi

# Deploy contract to mainnet
echo ""
echo "🚀 Deploying PAPRD contract to MAINNET..."
echo "📡 Sending deployment transaction to production blockchain..."

DEPLOY_PAYLOAD=$(cat <<EOF
{
  "owner": "$FOUNDER_ADDRESS",
  "name": "$CONTRACT_NAME",
  "code": "$WASM_BASE64",
  "fee": $DEPLOYMENT_FEE,
  "privateKey": "$FOUNDER_PRIVATE_KEY"
}
EOF
)

DEPLOY_RESPONSE=$(curl -s -X POST "$MAINNET_API/contracts/deploy" \
  -H "Content-Type: application/json" \
  -d "$DEPLOY_PAYLOAD")

# Check deployment result
if echo "$DEPLOY_RESPONSE" | grep -q "contractId"; then
    CONTRACT_ID=$(echo "$DEPLOY_RESPONSE" | grep -o '"contractId":"[^"]*"' | cut -d'"' -f4)
    echo ""
    echo "🎉 PAPRD CONTRACT DEPLOYED TO MAINNET! 🎉"
    echo "=========================================="
    echo ""
    echo "🎯 MAINNET Contract Details:"
    echo "Contract ID: $CONTRACT_ID"
    echo "Owner: $FOUNDER_ADDRESS"
    echo "Name: $CONTRACT_NAME"
    echo "Symbol: PAPRD"
    echo "Decimals: 18"
    echo "Blockchain: BINOMENA MAINNET"
    echo "API: $MAINNET_API"
    echo "Deployment Fee: $DEPLOYMENT_FEE BNM"
    echo ""
    
    # Save mainnet contract info
    cat > paprd_mainnet_contract.json <<EOF
{
  "contractId": "$CONTRACT_ID",
  "owner": "$FOUNDER_ADDRESS",
  "name": "$CONTRACT_NAME",
  "symbol": "PAPRD",
  "decimals": 18,
  "blockchain": "BINOMENA_MAINNET",
  "apiEndpoint": "$MAINNET_API",
  "deployedAt": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
  "deploymentFee": "$DEPLOYMENT_FEE BNM",
  "wasmSize": $WASM_SIZE,
  "network": "PRODUCTION"
}
EOF
    
    echo "📄 Mainnet contract info saved to: paprd_mainnet_contract.json"
    echo ""
    
    # Test basic contract functions on mainnet
    echo "🧪 Testing contract on mainnet..."
    
    # Test get_symbol
    echo "Testing get_symbol() on mainnet..."
    SYMBOL_PAYLOAD=$(cat <<EOF
{
  "caller": "$FOUNDER_ADDRESS",
  "function": "get_symbol",
  "params": [],
  "fee": 0.001,
  "privateKey": "$FOUNDER_PRIVATE_KEY"
}
EOF
)
    
    SYMBOL_RESPONSE=$(curl -s -X POST "$MAINNET_API/contracts/$CONTRACT_ID/execute" \
      -H "Content-Type: application/json" \
      -d "$SYMBOL_PAYLOAD")
    
    if echo "$SYMBOL_RESPONSE" | grep -q "PAPRD"; then
        echo "✅ Symbol test passed on MAINNET: PAPRD"
    else
        echo "⚠️ Symbol test result: $SYMBOL_RESPONSE"
    fi
    
    echo ""
    echo "🎊 PAPRD STABLECOIN SUCCESSFULLY DEPLOYED TO BINOMENA MAINNET! 🎊"
    echo ""
    echo "🌐 MAINNET ACCESS INFORMATION:"
    echo "Contract ID: $CONTRACT_ID"
    echo "Mainnet API: $MAINNET_API/contracts/$CONTRACT_ID"
    echo "Status Check: $MAINNET_API/contracts"
    echo ""
    echo "📱 WALLET INTEGRATION:"
    echo "Add PAPRD to wallets using:"
    echo "Contract Address: $CONTRACT_ID"
    echo "Symbol: PAPRD"
    echo "Decimals: 18"
    echo "Network: Binomena Mainnet"
    echo ""
    echo "🔗 PRODUCTION LINKS:"
    echo "Blockchain Explorer: $MAINNET_API"
    echo "Contract Info: $MAINNET_API/contracts/$CONTRACT_ID"
    echo "Balance Check: $MAINNET_API/balance/[ADDRESS]"
    
else
    echo ""
    echo "❌ MAINNET DEPLOYMENT FAILED!"
    echo "Response: $DEPLOY_RESPONSE"
    echo ""
    echo "🔍 Troubleshooting:"
    echo "1. Check mainnet connectivity: curl $MAINNET_API/status"
    echo "2. Verify founder balance: curl $MAINNET_API/balance/$FOUNDER_ADDRESS"
    echo "3. Check deployment logs above for errors"
    exit 1
fi 