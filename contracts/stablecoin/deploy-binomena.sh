#!/bin/bash

# PAPRD Stablecoin Deployment Script for Binomena Blockchain
# This script deploys the Paper Dollar (PAPRD) stablecoin to the binomena blockchain

set -e

echo "🚀 PAPRD Stablecoin Deployment to Binomena Blockchain"
echo "==================================================="
echo ""

# Configuration
FOUNDER_ADDRESS="AdNe6c3ce54e4371d056c7c566675ba16909eb2e9534"
FOUNDER_PRIVATE_KEY="508f988eaa468e8092ea01cbb484d474b705b45047f3bcae748749722f637fdf"
CONTRACT_NAME="Paper Dollar Stablecoin (PAPRD)"
DEPLOYMENT_FEE="1.0"  # 1 BNM for deployment
BLOCKCHAIN_API="http://localhost:8080"

echo "📋 Deployment Configuration:"
echo "Founder Address: $FOUNDER_ADDRESS"
echo "Contract Name: $CONTRACT_NAME"
echo "Deployment Fee: $DEPLOYMENT_FEE BNM"
echo "Blockchain API: $BLOCKCHAIN_API"
echo ""

# Check if binomena blockchain is running
echo "🔍 Checking binomena blockchain status..."
if ! curl -s "$BLOCKCHAIN_API/health" >/dev/null 2>&1; then
    echo "❌ Binomena blockchain is not running!"
    echo "Please start the blockchain first:"
    echo "  cd /Users/macbook/binomena"
    echo "  ./app --api-port 8080 --p2p-port 9000 --id node1"
    echo ""
    read -p "Would you like me to start the blockchain now? (y/n): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo "🚀 Starting binomena blockchain..."
        cd /Users/macbook/binomena
        ./app --api-port 8080 --p2p-port 9000 --id node1 &
        BLOCKCHAIN_PID=$!
        echo "Blockchain started with PID: $BLOCKCHAIN_PID"
        echo "Waiting 10 seconds for blockchain to initialize..."
        sleep 10
        cd contracts/stablecoin
    else
        echo "Deployment cancelled. Please start the blockchain and try again."
        exit 1
    fi
fi

# Build the contract
echo "🔨 Building PAPRD contract..."
if [ ! -f "pkg/paper_dollar_stablecoin_bg.wasm" ]; then
    echo "Contract not built. Building now..."
    ./deploy.sh
fi

# Convert WASM to base64
echo "📦 Converting WASM to base64..."
WASM_BASE64=$(base64 -i pkg/paper_dollar_stablecoin_bg.wasm)

# Check founder balance
echo "💰 Checking founder balance..."
BALANCE_RESPONSE=$(curl -s "$BLOCKCHAIN_API/token/balance/$FOUNDER_ADDRESS" || echo '{"balance": 0}')
FOUNDER_BALANCE=$(echo $BALANCE_RESPONSE | grep -o '"balance":[0-9.]*' | cut -d':' -f2 || echo "0")

echo "Founder balance: $FOUNDER_BALANCE BNM"

if (( $(echo "$FOUNDER_BALANCE < $DEPLOYMENT_FEE" | bc -l) )); then
    echo "❌ Insufficient balance for deployment!"
    echo "Required: $DEPLOYMENT_FEE BNM"
    echo "Available: $FOUNDER_BALANCE BNM"
    echo ""
    echo "💡 You can get BNM tokens by:"
    echo "1. Mining blocks on the blockchain"
    echo "2. Requesting tokens from the faucet"
    echo "3. Transferring from another wallet"
    exit 1
fi

# Deploy contract
echo "🚀 Deploying PAPRD contract to binomena blockchain..."

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

echo "📡 Sending deployment transaction..."
DEPLOY_RESPONSE=$(curl -s -X POST "$BLOCKCHAIN_API/contracts/deploy" \
  -H "Content-Type: application/json" \
  -d "$DEPLOY_PAYLOAD")

# Check deployment result
if echo "$DEPLOY_RESPONSE" | grep -q "contractId"; then
    CONTRACT_ID=$(echo "$DEPLOY_RESPONSE" | grep -o '"contractId":"[^"]*"' | cut -d'"' -f4)
    echo ""
    echo "✅ PAPRD Contract deployed successfully!"
    echo ""
    echo "🎯 Contract Details:"
    echo "Contract ID: $CONTRACT_ID"
    echo "Owner: $FOUNDER_ADDRESS"
    echo "Name: $CONTRACT_NAME"
    echo "Deployment Fee: $DEPLOYMENT_FEE BNM"
    echo ""
    
    # Save contract info
    cat > paprd_contract_info.json <<EOF
{
  "contractId": "$CONTRACT_ID",
  "owner": "$FOUNDER_ADDRESS",
  "name": "$CONTRACT_NAME",
  "symbol": "PAPRD",
  "decimals": 18,
  "deployedAt": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
  "blockchainAPI": "$BLOCKCHAIN_API"
}
EOF
    
    echo "📄 Contract info saved to: paprd_contract_info.json"
    echo ""
    
    # Test contract functions
    echo "🧪 Testing contract functions..."
    
    # Test get_symbol
    echo "Testing get_symbol()..."
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
    
    SYMBOL_RESPONSE=$(curl -s -X POST "$BLOCKCHAIN_API/contracts/$CONTRACT_ID/execute" \
      -H "Content-Type: application/json" \
      -d "$SYMBOL_PAYLOAD")
    
    if echo "$SYMBOL_RESPONSE" | grep -q "PAPRD"; then
        echo "✅ Symbol test passed: PAPRD"
    else
        echo "⚠️ Symbol test failed"
    fi
    
    # Test get_name
    echo "Testing get_name()..."
    NAME_PAYLOAD=$(cat <<EOF
{
  "caller": "$FOUNDER_ADDRESS",
  "function": "get_name",
  "params": [],
  "fee": 0.001,
  "privateKey": "$FOUNDER_PRIVATE_KEY"
}
EOF
)
    
    NAME_RESPONSE=$(curl -s -X POST "$BLOCKCHAIN_API/contracts/$CONTRACT_ID/execute" \
      -H "Content-Type: application/json" \
      -d "$NAME_PAYLOAD")
    
    if echo "$NAME_RESPONSE" | grep -q "Paper Dollar"; then
        echo "✅ Name test passed: Paper Dollar"
    else
        echo "⚠️ Name test failed"
    fi
    
    echo ""
    echo "🎉 PAPRD Stablecoin successfully deployed to Binomena Blockchain!"
    echo ""
    echo "📖 Next Steps:"
    echo "1. Add authorized minters using: curl -X POST $BLOCKCHAIN_API/contracts/$CONTRACT_ID/execute"
    echo "2. Set up collateral management"
    echo "3. Start minting PAPRD tokens"
    echo ""
    echo "🔗 Contract API URL: $BLOCKCHAIN_API/contracts/$CONTRACT_ID"
    echo "📱 You can now add PAPRD to wallets using Contract ID: $CONTRACT_ID"
    
else
    echo ""
    echo "❌ Contract deployment failed!"
    echo "Response: $DEPLOY_RESPONSE"
    exit 1
fi 