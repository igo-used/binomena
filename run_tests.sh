#!/bin/bash

set -e  # Exit on first error
echo "🔧 Building AssemblyScript smart contracts..."

# Navigate to contracts folder and build the WASM contract
cd contracts
npm install
npm run asbuild

echo "✅ AssemblyScript contract built successfully."

# Navigate back to root
cd ..

echo "🧪 Running Go tests..."

# Run wallet module tests
echo "➡️  Wallet module tests:"
go test ./wallet -v

# Run core module tests
echo "➡️  Core module tests:"
go test ./core -v

# Run smartcontract module tests
echo "➡️  Smartcontract module tests:"
go test ./smartcontract -v

# Run integration tests
echo "➡️  Integration tests:"
go test ./tests -v

# Run the test application
echo "🚀 Running the blockchain test application..."
go run cmd/test_app/main.go contracts/build/optimized.wasm

echo "✅ All tests completed!"
