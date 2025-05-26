#!/bin/bash

# PAPRD Stablecoin Mock Deployment Script
# This script simulates the deployment of PAPRD stablecoin for testing

set -e

echo "🚀 PAPRD Stablecoin Mock Deployment"
echo "=================================="
echo ""

# Configuration
FOUNDER_ADDRESS="AdNe6c3ce54e4371d056c7c566675ba16909eb2e9534"
FOUNDER_PRIVATE_KEY="508f988eaa468e8092ea01cbb484d474b705b45047f3bcae748749722f637fdf"
CONTRACT_NAME="Paper Dollar Stablecoin (PAPRD)"
DEPLOYMENT_FEE="1.0"

echo "📋 Deployment Configuration:"
echo "Founder Address: $FOUNDER_ADDRESS"
echo "Contract Name: $CONTRACT_NAME"
echo "Deployment Fee: $DEPLOYMENT_FEE BNM"
echo ""

# Build the contract
echo "🔨 Building PAPRD contract..."
if [ ! -f "pkg/paper_dollar_stablecoin_bg.wasm" ]; then
    echo "Contract not built. Building now..."
    ./deploy.sh
fi

# Convert WASM to base64
echo "📦 Converting WASM to base64..."
WASM_BASE64=$(base64 -i pkg/paper_dollar_stablecoin_bg.wasm)
WASM_SIZE=$(wc -c < pkg/paper_dollar_stablecoin_bg.wasm)

echo "✅ Contract built successfully!"
echo "WASM Size: $WASM_SIZE bytes"
echo "Base64 Size: ${#WASM_BASE64} characters"
echo ""

# Generate mock contract ID
CONTRACT_ID="AdNe$(echo -n "$CONTRACT_NAME$FOUNDER_ADDRESS$(date)" | shasum -a 256 | cut -c1-60)"

echo "🎯 Mock Deployment Results:"
echo "Contract ID: $CONTRACT_ID"
echo "Owner: $FOUNDER_ADDRESS"
echo "Name: $CONTRACT_NAME"
echo "Symbol: PAPRD"
echo "Decimals: 18"
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
  "wasmSize": $WASM_SIZE,
  "base64Size": ${#WASM_BASE64},
  "deploymentStatus": "mock_success",
  "blockchainAPI": "http://localhost:8080"
}
EOF

echo "📄 Contract info saved to: paprd_contract_info.json"
echo ""

# Save WASM for deployment
echo "$WASM_BASE64" > paprd_contract.base64
echo "💾 WASM base64 saved to: paprd_contract.base64"
echo ""

echo "🧪 Mock Contract Function Tests:"
echo "✅ get_symbol() -> PAPRD"
echo "✅ get_name() -> Paper Dollar"
echo "✅ get_decimals() -> 18"
echo "✅ get_total_supply() -> 0"
echo "✅ get_administrator() -> $FOUNDER_ADDRESS"
echo ""

echo "🎉 PAPRD Stablecoin Mock Deployment Completed!"
echo ""
echo "📖 Next Steps for Real Deployment:"
echo "1. Fix binomena blockchain database issues"
echo "2. Start blockchain with: ./app --api-port 8080 --p2p-port 9000 --id paprd-node --use-db=false"
echo "3. Run: ./deploy-binomena.sh"
echo ""
echo "🔗 Mock Contract Details:"
echo "Contract ID: $CONTRACT_ID"
echo "Ready for deployment to binomena blockchain!"
echo ""

# Create deployment payload for when blockchain is ready
cat > deployment_payload.json <<EOF
{
  "owner": "$FOUNDER_ADDRESS",
  "name": "$CONTRACT_NAME",
  "code": "$WASM_BASE64",
  "fee": $DEPLOYMENT_FEE,
  "privateKey": "$FOUNDER_PRIVATE_KEY"
}
EOF

echo "📋 Deployment payload saved to: deployment_payload.json"
echo "Use this payload when the blockchain is running!" 