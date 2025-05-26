#!/bin/bash

# PAPRD Stablecoin Administration Script
# For Binomena Blockchain Founder
# Contract: AdNec919d52b5eedb9662999ec97cae93677440cb0cea9e11aed5cee637ee7f0

set -e

# Configuration
FOUNDER_ADDRESS="AdNe6c3ce54e4371d056c7c566675ba16909eb2e9534"
FOUNDER_PRIVATE_KEY="508f988eaa468e8092ea01cbb484d474b705b45047f3bcae748749722f637fdf"
CONTRACT_ID="AdNec919d52b5eedb9662999ec97cae93677440cb0cea9e11aed5cee637ee7f0"
MAINNET_API="https://binomena-node.onrender.com"

echo "🏛️  PAPRD Stablecoin Administration Panel"
echo "========================================"
echo "Contract ID: $CONTRACT_ID"
echo "Founder: $FOUNDER_ADDRESS"
echo "Network: Binomena Mainnet"
echo ""

# Function to check contract status
check_status() {
    echo "📊 Contract Status:"
    curl -s "$MAINNET_API/contracts" | grep -o '"count":[0-9]*' | cut -d':' -f2
    echo " contracts deployed on mainnet"
    echo ""
}

# Function to check founder balance
check_balance() {
    echo "💰 Founder Balance:"
    BALANCE=$(curl -s "$MAINNET_API/balance/$FOUNDER_ADDRESS" | grep -o '"balance":[0-9.]*' | cut -d':' -f2)
    echo "BNM Balance: $BALANCE"
    echo ""
}

# Function to mint PAPRD tokens
mint_paprd() {
    local to_address=$1
    local amount=$2
    local collateral_type=${3:-"BNM"}
    
    echo "🏦 Minting $amount PAPRD to $to_address"
    echo "Collateral Type: $collateral_type"
    
    # Add collateral first (simplified for demo)
    echo "Adding collateral..."
    curl -X POST "$MAINNET_API/contracts/$CONTRACT_ID/execute" \
        -H "Content-Type: application/json" \
        -d "{
            \"caller\": \"$FOUNDER_ADDRESS\",
            \"function\": \"add_collateral\",
            \"params\": [\"$collateral_type\", \"$amount\"],
            \"fee\": 0.01,
            \"privateKey\": \"$FOUNDER_PRIVATE_KEY\"
        }"
    
    echo ""
    echo "Minting tokens..."
    curl -X POST "$MAINNET_API/contracts/$CONTRACT_ID/execute" \
        -H "Content-Type: application/json" \
        -d "{
            \"caller\": \"$FOUNDER_ADDRESS\",
            \"function\": \"mint\",
            \"params\": [\"$to_address\", \"$amount\", \"$collateral_type\"],
            \"fee\": 0.01,
            \"privateKey\": \"$FOUNDER_PRIVATE_KEY\"
        }"
    echo ""
}

# Function to add a minter
add_minter() {
    local minter_address=$1
    
    echo "👤 Adding minter: $minter_address"
    curl -X POST "$MAINNET_API/contracts/$CONTRACT_ID/execute" \
        -H "Content-Type: application/json" \
        -d "{
            \"caller\": \"$FOUNDER_ADDRESS\",
            \"function\": \"add_minter\",
            \"params\": [\"$minter_address\"],
            \"fee\": 0.01,
            \"privateKey\": \"$FOUNDER_PRIVATE_KEY\"
        }"
    echo ""
}

# Function to blacklist an address
blacklist_address() {
    local address=$1
    
    echo "🚫 Blacklisting address: $address"
    curl -X POST "$MAINNET_API/contracts/$CONTRACT_ID/execute" \
        -H "Content-Type: application/json" \
        -d "{
            \"caller\": \"$FOUNDER_ADDRESS\",
            \"function\": \"blacklist\",
            \"params\": [\"$address\"],
            \"fee\": 0.01,
            \"privateKey\": \"$FOUNDER_PRIVATE_KEY\"
        }"
    echo ""
}

# Function to pause the contract
pause_contract() {
    echo "⏸️  Pausing PAPRD contract..."
    curl -X POST "$MAINNET_API/contracts/$CONTRACT_ID/execute" \
        -H "Content-Type: application/json" \
        -d "{
            \"caller\": \"$FOUNDER_ADDRESS\",
            \"function\": \"pause\",
            \"params\": [],
            \"fee\": 0.01,
            \"privateKey\": \"$FOUNDER_PRIVATE_KEY\"
        }"
    echo ""
}

# Function to unpause the contract
unpause_contract() {
    echo "▶️  Unpausing PAPRD contract..."
    curl -X POST "$MAINNET_API/contracts/$CONTRACT_ID/execute" \
        -H "Content-Type: application/json" \
        -d "{
            \"caller\": \"$FOUNDER_ADDRESS\",
            \"function\": \"unpause\",
            \"params\": [],
            \"fee\": 0.01,
            \"privateKey\": \"$FOUNDER_PRIVATE_KEY\"
        }"
    echo ""
}

# Function to set collateral ratio
set_collateral_ratio() {
    local ratio=$1
    
    echo "📊 Setting collateral ratio to $ratio%"
    curl -X POST "$MAINNET_API/contracts/$CONTRACT_ID/execute" \
        -H "Content-Type: application/json" \
        -d "{
            \"caller\": \"$FOUNDER_ADDRESS\",
            \"function\": \"set_collateral_ratio\",
            \"params\": [\"$ratio\"],
            \"fee\": 0.01,
            \"privateKey\": \"$FOUNDER_PRIVATE_KEY\"
        }"
    echo ""
}

# Function to transfer tokens
transfer_paprd() {
    local to_address=$1
    local amount=$2
    
    echo "💸 Transferring $amount PAPRD to $to_address"
    curl -X POST "$MAINNET_API/contracts/$CONTRACT_ID/execute" \
        -H "Content-Type: application/json" \
        -d "{
            \"caller\": \"$FOUNDER_ADDRESS\",
            \"function\": \"transfer\",
            \"params\": [\"$to_address\", \"$amount\"],
            \"fee\": 0.01,
            \"privateKey\": \"$FOUNDER_PRIVATE_KEY\"
        }"
    echo ""
}

# Main menu
case "$1" in
    "status")
        check_status
        ;;
    "balance")
        check_balance
        ;;
    "mint")
        mint_paprd "$2" "$3" "$4"
        ;;
    "add-minter")
        add_minter "$2"
        ;;
    "blacklist")
        blacklist_address "$2"
        ;;
    "pause")
        pause_contract
        ;;
    "unpause")
        unpause_contract
        ;;
    "set-ratio")
        set_collateral_ratio "$2"
        ;;
    "transfer")
        transfer_paprd "$2" "$3"
        ;;
    *)
        echo "PAPRD Stablecoin Administration Commands:"
        echo "========================================="
        echo ""
        echo "📊 Information:"
        echo "  ./paprd-admin.sh status          - Check contract status"
        echo "  ./paprd-admin.sh balance         - Check founder balance"
        echo ""
        echo "🏦 Minting Operations:"
        echo "  ./paprd-admin.sh mint [address] [amount] [type]  - Mint PAPRD tokens"
        echo "  ./paprd-admin.sh transfer [address] [amount]     - Transfer PAPRD"
        echo ""
        echo "👥 Administrative:"
        echo "  ./paprd-admin.sh add-minter [address]      - Add authorized minter"
        echo "  ./paprd-admin.sh blacklist [address]       - Blacklist address"
        echo "  ./paprd-admin.sh pause                     - Emergency pause"
        echo "  ./paprd-admin.sh unpause                   - Resume operations"
        echo "  ./paprd-admin.sh set-ratio [percentage]    - Set collateral ratio"
        echo ""
        echo "Examples:"
        echo "  ./paprd-admin.sh mint AdNe123... 1000 BNM"
        echo "  ./paprd-admin.sh add-minter AdNe456..."
        echo "  ./paprd-admin.sh set-ratio 150"
        echo ""
        ;;
esac 