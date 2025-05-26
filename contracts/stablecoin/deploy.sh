#!/bin/bash

# Paper Dollar Stablecoin Deployment Script
# This script deploys the stablecoin with the founder as administrator

set -e

echo "🚀 Paper Dollar Stablecoin Deployment"
echo "====================================="

# Security warning
echo "⚠️  SECURITY NOTICE:"
echo "This script will set up the contract with the founder address as administrator."
echo "Make sure you're running this in a secure environment."
echo ""

# Founder address (will be set as contract owner)
FOUNDER_ADDRESS="AdNe6c3ce54e4371d056c7c566675ba16909eb2e9534"

echo "📋 Deployment Configuration:"
echo "Founder Address: $FOUNDER_ADDRESS"
echo "Contract: Paper Dollar Stablecoin"
echo "Initial Collateral Ratio: 150%"
echo ""

# Check if Rust is installed
if ! command -v cargo &> /dev/null; then
    echo "❌ Rust not found. Installing..."
    ./setup.sh
fi

# Build the contract
echo "🔨 Building contract..."
cargo build --release

# Run tests to ensure everything works
echo "🧪 Running tests..."
cargo test -- --test-threads=1

# Build WebAssembly
echo "📦 Building WebAssembly..."
if ! command -v wasm-pack &> /dev/null; then
    echo "Installing wasm-pack..."
    curl https://rustwasm.github.io/wasm-pack/installer/init.sh -sSf | sh
fi

wasm-pack build --target web --out-dir pkg

echo ""
echo "✅ Contract built successfully!"
echo ""
echo "📝 Next steps for blockchain deployment:"
echo "1. The contract is ready for deployment"
echo "2. Founder address will be set as owner: $FOUNDER_ADDRESS"
echo "3. Deploy using your blockchain's deployment tools"
echo "4. Initialize with ./initialize-contract.sh after deployment"
echo ""
echo "🔐 Security reminder: Keep private keys secure and never commit them!" 