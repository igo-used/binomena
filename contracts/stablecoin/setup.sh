#!/bin/bash

# Setup script for Paper Dollar Stablecoin Rust Contract

set -e

echo "🦀 Setting up Rust development environment..."

# Check if Rust is already installed
if command -v cargo &> /dev/null; then
    echo "✅ Rust is already installed"
    rustc --version
    cargo --version
else
    echo "📦 Installing Rust toolchain..."
    curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y
    source ~/.cargo/env
    echo "✅ Rust installed successfully"
fi

# Add WASM target
echo "🎯 Adding WebAssembly target..."
rustup target add wasm32-unknown-unknown

# Install wasm-pack if not already installed
if command -v wasm-pack &> /dev/null; then
    echo "✅ wasm-pack is already installed"
else
    echo "📦 Installing wasm-pack..."
    curl https://rustwasm.github.io/wasm-pack/installer/init.sh -sSf | sh
    echo "✅ wasm-pack installed successfully"
fi

# Install useful development tools
echo "🔧 Installing development tools..."
rustup component add clippy rustfmt

echo "🧪 Running initial checks..."

# Check compilation
echo "Checking if code compiles..."
cargo check

# Run tests
echo "Running tests..."
cargo test

# Check formatting
echo "Checking code formatting..."
cargo fmt --check || (echo "⚠️  Code needs formatting. Run 'cargo fmt' to fix." && cargo fmt)

# Run clippy lints
echo "Running clippy lints..."
cargo clippy -- -W clippy::all

echo ""
echo "🎉 Setup completed successfully!"
echo ""
echo "Available commands:"
echo "  cargo check          - Check compilation"
echo "  cargo test           - Run tests"
echo "  cargo fmt            - Format code"
echo "  cargo clippy         - Run lints"
echo "  ./build.sh           - Build WASM packages"
echo ""
echo "📖 See README.md for detailed usage instructions" 