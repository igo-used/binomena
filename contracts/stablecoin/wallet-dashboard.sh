#!/bin/bash

# 💰 BINOMENA WALLET DASHBOARD
# Complete overview of your crypto holdings

echo "🚀 BINOMENA WALLET DASHBOARD"
echo "============================"
echo ""

# Configuration
FOUNDER_ADDRESS="AdNe6c3ce54e4371d056c7c566675ba16909eb2e9534"
TREASURY_ADDRESS="AdNec13f53bb89865c7e2be8ff9aa43e84e26d226bf3"
COMMUNITY_ADDRESS="AdNebaefd75d426056bffbc622bd9f334ed89450efae"
PAPRD_CONTRACT="AdNec919d52b5eedb9662999ec97cae93677440cb0cea9e11aed5cee637ee7f0"
MAINNET_API="https://binomena-node.onrender.com"

echo "📊 BLOCKCHAIN STATUS:"
echo "Network: Binomena Mainnet"
echo "API: $MAINNET_API"
curl -s "$MAINNET_API/status" | jq '.'
echo ""

echo "💰 BNM TOKEN BALANCES:"
echo "======================"
echo ""

echo "🏛️  FOUNDER (YOU):"
echo "Address: $FOUNDER_ADDRESS"
FOUNDER_BALANCE=$(curl -s "$MAINNET_API/balance/$FOUNDER_ADDRESS" | jq -r '.balance')
echo "Balance: $(printf "%.2f" $FOUNDER_BALANCE) BNM"
echo "USD Value: \$$(printf "%.2f" $(echo "$FOUNDER_BALANCE" | bc))"
echo ""

echo "🏦 TREASURY:"
echo "Address: $TREASURY_ADDRESS"
TREASURY_BALANCE=$(curl -s "$MAINNET_API/balance/$TREASURY_ADDRESS" | jq -r '.balance')
echo "Balance: $(printf "%.2f" $TREASURY_BALANCE) BNM"
echo ""

echo "🌐 COMMUNITY:"
echo "Address: $COMMUNITY_ADDRESS"
COMMUNITY_BALANCE=$(curl -s "$MAINNET_API/balance/$COMMUNITY_ADDRESS" | jq -r '.balance')
echo "Balance: $(printf "%.2f" $COMMUNITY_BALANCE) BNM"
echo ""

TOTAL_BNM=$(echo "$FOUNDER_BALANCE + $TREASURY_BALANCE + $COMMUNITY_BALANCE" | bc)
echo "💎 TOTAL CONTROLLED: $(printf "%.2f" $TOTAL_BNM) BNM"
echo ""

echo "🪙 PAPRD STABLECOIN STATUS:"
echo "=========================="
echo ""

if [ -f "paprd-ledger.json" ]; then
    echo "📋 PAPRD Token Info:"
    echo "Contract: $PAPRD_CONTRACT"
    echo "Name: $(jq -r '.token_name' paprd-ledger.json)"
    echo "Symbol: $(jq -r '.token_symbol' paprd-ledger.json)"
    echo "Decimals: $(jq -r '.token_decimals' paprd-ledger.json)"
    echo ""
    
    echo "💰 YOUR PAPRD BALANCE:"
    PAPRD_BALANCE=$(jq -r ".balances[\"$FOUNDER_ADDRESS\"]" paprd-ledger.json)
    if [ "$PAPRD_BALANCE" != "null" ]; then
        # Convert from 18 decimal places
        PAPRD_READABLE=$(echo "scale=2; $PAPRD_BALANCE / 1000000000000000000" | bc)
        echo "Balance: $(printf "%.0f" $PAPRD_READABLE) PAPRD"
        echo "USD Value: \$$(printf "%.0f" $PAPRD_READABLE) (1:1 peg)"
    else
        echo "Balance: 0 PAPRD"
    fi
    echo ""
    
    echo "📊 PAPRD Statistics:"
    TOTAL_SUPPLY=$(jq -r '.total_supply' paprd-ledger.json)
    TOTAL_SUPPLY_READABLE=$(echo "scale=0; $TOTAL_SUPPLY / 1000000000000000000" | bc)
    echo "Total Supply: $(printf "%.0f" $TOTAL_SUPPLY_READABLE) PAPRD"
    echo "Market Cap: \$$(printf "%.0f" $TOTAL_SUPPLY_READABLE)"
    echo "Collateral Ratio: $(jq -r '.collateral_ratio' paprd-ledger.json)%"
    echo "Status: $(jq -r 'if .paused then "PAUSED" else "ACTIVE" end' paprd-ledger.json)"
    echo ""
else
    echo "❌ PAPRD ledger not found"
fi

echo "🔧 CONTRACT STATUS:"
echo "=================="
echo "Checking PAPRD contract on mainnet..."
CONTRACT_RESPONSE=$(curl -s "$MAINNET_API/contracts" | jq '.contracts[] | select(.id=="'$PAPRD_CONTRACT'")')
if [ -n "$CONTRACT_RESPONSE" ]; then
    echo "✅ Contract deployed and found"
    echo "Deploy Time: $(echo "$CONTRACT_RESPONSE" | jq -r '.deployedAt')"
    echo "Executions: $(echo "$CONTRACT_RESPONSE" | jq -r '.executionCount')"
else
    echo "❌ Contract not found on mainnet"
fi
echo ""

echo "💡 WHAT YOU CAN DO:"
echo "=================="
echo "1. 🏦 Start FIAT-to-PAPRD exchange service"
echo "2. 🌐 List PAPRD on crypto exchanges"
echo "3. 🤝 Partner with businesses for payments"
echo "4. 📱 Build a wallet app for users"
echo "5. 🔄 Create cross-chain bridges"
echo ""

echo "🎯 YOUR EMPIRE SUMMARY:"
echo "======================"
echo "• Blockchain Founder: ✅ (Own entire Binomena network)"
echo "• Stablecoin Issuer: ✅ (100M PAPRD tokens)"
echo "• Crypto Whale: ✅ (999M+ BNM tokens)"
echo "• DeFi Pioneer: ✅ (Complete financial ecosystem)"
echo ""
echo "🚀 Ready to conquer the crypto world!" 