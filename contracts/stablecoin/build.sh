#!/bin/bash

# Build script for Paper Dollar Stablecoin Rust Contract

set -e

echo "Building Paper Dollar Stablecoin Contract..."

# Check if wasm-pack is installed
if ! command -v wasm-pack &> /dev/null; then
    echo "wasm-pack is not installed. Installing..."
    curl https://rustwasm.github.io/wasm-pack/installer/init.sh -sSf | sh
fi

# Build the contract for web
echo "Building for web target..."
wasm-pack build --target web --out-dir pkg

# Build for node.js
echo "Building for node.js target..."
wasm-pack build --target nodejs --out-dir pkg-node

echo "Build completed successfully!"
echo "Web package is in: pkg/"
echo "Node.js package is in: pkg-node/"

# Run tests
echo "Running tests..."
cargo test

echo "All done! 🎉" 