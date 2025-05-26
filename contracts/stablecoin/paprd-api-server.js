// PAPRD Token API Server
// Provides REST endpoints for UI integration

const express = require('express');
const cors = require('cors');
const fs = require('fs');
const path = require('path');

const app = express();
const PORT = process.env.PORT || 3001;

// Enable CORS for your UI
app.use(cors({
    origin: ['https://www.binomchainapp.fyi', 'http://localhost:3000'],
    credentials: true
}));

app.use(express.json());

// PAPRD Ledger file path
const LEDGER_PATH = path.join(__dirname, 'paprd-ledger.json');

// Helper function to read ledger
function readLedger() {
    try {
        const data = fs.readFileSync(LEDGER_PATH, 'utf8');
        return JSON.parse(data);
    } catch (error) {
        // Initialize default ledger
        const defaultLedger = {
            contract_id: "AdNec919d52b5eedb9662999ec97cae93677440cb0cea9e11aed5cee637ee7f0",
            token_name: "Paper Dollar Stablecoin",
            token_symbol: "PAPRD",
            token_decimals: 18,
            total_supply: "100000000000000000000000000",
            owner: "AdNe6c3ce54e4371d056c7c566675ba16909eb2e9534",
            balances: {
                "AdNe6c3ce54e4371d056c7c566675ba16909eb2e9534": "100000000000000000000000000"
            },
            minters: ["AdNe6c3ce54e4371d056c7c566675ba16909eb2e9534"],
            blacklisted: [],
            paused: false,
            transactions: []
        };
        writeLedger(defaultLedger);
        return defaultLedger;
    }
}

// Helper function to write ledger
function writeLedger(ledger) {
    fs.writeFileSync(LEDGER_PATH, JSON.stringify(ledger, null, 2));
}

// Helper function to convert from 18 decimals
function fromDecimals(amount) {
    return (BigInt(amount) / BigInt(10 ** 18)).toString();
}

// Helper function to convert to 18 decimals
function toDecimals(amount) {
    return (BigInt(amount) * BigInt(10 ** 18)).toString();
}

// 📊 GET /paprd/info - Token information
app.get('/paprd/info', (req, res) => {
    const ledger = readLedger();
    res.json({
        name: ledger.token_name,
        symbol: ledger.token_symbol,
        decimals: ledger.token_decimals,
        totalSupply: fromDecimals(ledger.total_supply),
        owner: ledger.owner,
        contract: ledger.contract_id,
        paused: ledger.paused
    });
});

// 💰 GET /paprd/balance/:address - Get PAPRD balance
app.get('/paprd/balance/:address', (req, res) => {
    const { address } = req.params;
    const ledger = readLedger();
    
    const balance = ledger.balances[address] || "0";
    const readableBalance = fromDecimals(balance);
    
    res.json({
        address,
        balance: readableBalance,
        balanceWei: balance,
        symbol: "PAPRD"
    });
});

// 🔄 POST /paprd/transfer - Transfer PAPRD tokens
app.post('/paprd/transfer', (req, res) => {
    const { from, to, amount, privateKey } = req.body;
    
    // Validate request
    if (!from || !to || !amount) {
        return res.status(400).json({ 
            error: "Missing required fields: from, to, amount" 
        });
    }

    // TODO: Verify signature with privateKey in production
    
    const ledger = readLedger();
    
    // Check if contract is paused
    if (ledger.paused) {
        return res.status(400).json({ error: "Contract is paused" });
    }
    
    // Check blacklist
    if (ledger.blacklisted.includes(from) || ledger.blacklisted.includes(to)) {
        return res.status(400).json({ error: "Address is blacklisted" });
    }
    
    const amountWei = toDecimals(amount);
    const fromBalance = BigInt(ledger.balances[from] || "0");
    
    // Check balance
    if (fromBalance < BigInt(amountWei)) {
        return res.status(400).json({ error: "Insufficient balance" });
    }
    
    // Perform transfer
    ledger.balances[from] = (fromBalance - BigInt(amountWei)).toString();
    ledger.balances[to] = (BigInt(ledger.balances[to] || "0") + BigInt(amountWei)).toString();
    
    // Record transaction
    const tx = {
        id: `tx_${Date.now()}`,
        type: "transfer",
        from,
        to,
        amount: amountWei,
        timestamp: new Date().toISOString(),
        block: Math.floor(Date.now() / 1000),
        status: "confirmed"
    };
    
    ledger.transactions.push(tx);
    writeLedger(ledger);
    
    res.json({
        success: true,
        transaction: tx,
        newBalance: fromDecimals(ledger.balances[from])
    });
});

// 🪙 POST /paprd/mint - Mint PAPRD tokens (owner only)
app.post('/paprd/mint', (req, res) => {
    const { to, amount, caller, privateKey } = req.body;
    
    const ledger = readLedger();
    
    // Check if caller is owner
    if (caller !== ledger.owner) {
        return res.status(403).json({ error: "Only owner can mint tokens" });
    }
    
    const amountWei = toDecimals(amount);
    ledger.balances[to] = (BigInt(ledger.balances[to] || "0") + BigInt(amountWei)).toString();
    ledger.total_supply = (BigInt(ledger.total_supply) + BigInt(amountWei)).toString();
    
    // Record transaction
    const tx = {
        id: `tx_${Date.now()}`,
        type: "mint",
        from: "0x0000000000000000000000000000000000000000",
        to,
        amount: amountWei,
        timestamp: new Date().toISOString(),
        block: Math.floor(Date.now() / 1000),
        status: "confirmed"
    };
    
    ledger.transactions.push(tx);
    writeLedger(ledger);
    
    res.json({
        success: true,
        transaction: tx,
        newBalance: fromDecimals(ledger.balances[to]),
        totalSupply: fromDecimals(ledger.total_supply)
    });
});

// 📋 GET /paprd/transactions/:address - Get transaction history
app.get('/paprd/transactions/:address', (req, res) => {
    const { address } = req.params;
    const ledger = readLedger();
    
    const transactions = ledger.transactions
        .filter(tx => tx.from === address || tx.to === address)
        .sort((a, b) => new Date(b.timestamp) - new Date(a.timestamp))
        .slice(0, 50); // Last 50 transactions
    
    res.json({
        address,
        transactions: transactions.map(tx => ({
            ...tx,
            amount: fromDecimals(tx.amount)
        }))
    });
});

// 🏦 GET /paprd/wallet/:address - Combined wallet info (BNM + PAPRD)
app.get('/paprd/wallet/:address', async (req, res) => {
    const { address } = req.params;
    
    try {
        // Get PAPRD balance
        const paprdBalance = await fetch(`http://localhost:${PORT}/paprd/balance/${address}`).then(r => r.json());
        
        // Get BNM balance from mainnet
        const bnmResponse = await fetch(`https://binomena-node.onrender.com/balance/${address}`);
        const bnmBalance = await bnmResponse.json();
        
        res.json({
            address,
            balances: {
                BNM: {
                    balance: bnmBalance.balance || 0,
                    symbol: "BNM",
                    name: "Binomena Token"
                },
                PAPRD: {
                    balance: paprdBalance.balance || 0,
                    symbol: "PAPRD",
                    name: "Paper Dollar Stablecoin"
                }
            },
            totalValueUSD: (parseFloat(bnmBalance.balance || 0) + parseFloat(paprdBalance.balance || 0))
        });
    } catch (error) {
        res.status(500).json({ error: "Failed to fetch wallet data" });
    }
});

// 📊 GET /paprd/stats - Overall statistics
app.get('/paprd/stats', (req, res) => {
    const ledger = readLedger();
    
    const totalHolders = Object.keys(ledger.balances).filter(addr => 
        BigInt(ledger.balances[addr]) > 0
    ).length;
    
    res.json({
        totalSupply: fromDecimals(ledger.total_supply),
        totalHolders,
        totalTransactions: ledger.transactions.length,
        paused: ledger.paused,
        contract: ledger.contract_id
    });
});

// Health check
app.get('/health', (req, res) => {
    res.json({ status: 'healthy', service: 'PAPRD API', timestamp: new Date().toISOString() });
});

app.listen(PORT, () => {
    console.log(`🚀 PAPRD API Server running on port ${PORT}`);
    console.log(`📊 Health check: http://localhost:${PORT}/health`);
    console.log(`💰 Example: http://localhost:${PORT}/paprd/balance/AdNe6c3ce54e4371d056c7c566675ba16909eb2e9534`);
});

module.exports = app; 