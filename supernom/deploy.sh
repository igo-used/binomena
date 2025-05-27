#!/bin/bash

echo "🚀 Deploying SuperNom - Revolutionary Blockchain VPN to Render..."
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Step 1: Commit latest changes
echo "📦 Committing latest SuperNom changes..."
git add .
git commit -m "🚀 SuperNom ready for Render deployment - Blockchain VPN going global"

# Step 2: Push to GitHub (triggers auto-deploy)
echo "🌍 Pushing to GitHub (triggers Render auto-deploy)..."
git push origin main

echo "✅ SuperNom deployment initiated!"
echo ""
echo "🎯 NEXT STEPS:"
echo "1. Go to https://dashboard.render.com"
echo "2. Create new Web Service"
echo "3. Connect your GitHub repo: igo-used/binomena"
echo "4. Set build command: cd supernom && go mod download && go build -o supernom-app"
echo "5. Set start command: cd supernom && ./supernom-app"
echo "6. Set environment variables:"
echo "   - PORT=10000"
echo "   - SUPERNOM_PORT=10000" 
echo "   - BINOMENA_NODE_URL=https://binomena-node.onrender.com"
echo ""
echo "🚀 After deployment, SuperNom will be live at:"
echo "   https://supernom-vpn.onrender.com"
echo ""
echo "💰 THIS WILL MASSIVELY INCREASE YOUR BNM TOKEN VALUE!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" 