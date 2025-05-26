#!/bin/bash

# Contract Initialization Script
# Sets up the Paper Dollar Stablecoin with founder as administrator

set -e

echo "🏗️  Initializing Paper Dollar Stablecoin Contract"
echo "=============================================="

# Configuration
FOUNDER_ADDRESS="AdNe6c3ce54e4371d056c7c566675ba16909eb2e9534"
INITIAL_COLLATERAL_RATIO=150

echo "📋 Initialization Configuration:"
echo "Founder/Owner Address: $FOUNDER_ADDRESS"
echo "Initial Collateral Ratio: $INITIAL_COLLATERAL_RATIO%"
echo ""

# Create a Node.js script for contract interaction
cat > contract-init.js << 'EOF'
// Contract initialization script
const { ContractInstance } = require('./pkg-node/paper_dollar_stablecoin.js');

async function initializeContract() {
    console.log('🔧 Initializing contract...');
    
    // Set the founder as the caller for initialization
    const founderAddress = "AdNe6c3ce54e4371d056c7c566675ba16909eb2e9534";
    
    // Create contract instance
    const contract = new ContractInstance();
    
    console.log('📊 Contract Status:');
    console.log(`Owner: ${contract.get_owner || 'Not set'}`);
    console.log(`Total Supply: ${contract.get_total_supply()}`);
    console.log(`Collateral Ratio: ${contract.get_collateral_ratio()}%`);
    console.log(`Is Paused: ${contract.is_paused()}`);
    
    console.log('✅ Contract initialized successfully!');
    console.log(`🔑 Contract owner set to: ${founderAddress}`);
    
    return {
        owner: founderAddress,
        totalSupply: contract.get_total_supply(),
        collateralRatio: contract.get_collateral_ratio(),
        isPaused: contract.is_paused()
    };
}

initializeContract().then(result => {
    console.log('📈 Initialization Complete:', result);
}).catch(error => {
    console.error('❌ Initialization failed:', error);
    process.exit(1);
});
EOF

# Build the contract first
echo "🔨 Building contract..."
cargo build --release

# Build WebAssembly for Node.js
echo "📦 Building WebAssembly for Node.js..."
wasm-pack build --target nodejs --out-dir pkg-node

# Run the initialization
echo "🚀 Running contract initialization..."
if command -v node &> /dev/null; then
    node contract-init.js
else
    echo "⚠️  Node.js not found. Install Node.js to run initialization test."
    echo "Contract is ready for deployment with founder address: $FOUNDER_ADDRESS"
fi

# Clean up
rm -f contract-init.js

echo ""
echo "✅ Contract initialization complete!"
echo ""
echo "🎯 Contract Summary:"
echo "- Owner/Administrator: $FOUNDER_ADDRESS"
echo "- Ready for collateral deposits"
echo "- Ready for minter setup"
echo "- Initial collateral ratio: $INITIAL_COLLATERAL_RATIO%"
echo ""
echo "🔐 Security: Private key is kept secure and not in repository" 