// PAPRD Token Integration for Binomena UI
// Add this to your existing JavaScript files

class PaprdAPI {
    constructor() {
        this.apiUrl = 'https://your-paprd-api.herokuapp.com'; // Replace with your deployed API
        this.bnmApiUrl = 'https://binomena-node.onrender.com';
    }

    // 📊 Get PAPRD token info
    async getTokenInfo() {
        try {
            const response = await fetch(`${this.apiUrl}/paprd/info`);
            return await response.json();
        } catch (error) {
            console.error('Error fetching PAPRD info:', error);
            return null;
        }
    }

    // 💰 Get PAPRD balance for address
    async getPaprdBalance(address) {
        try {
            const response = await fetch(`${this.apiUrl}/paprd/balance/${address}`);
            return await response.json();
        } catch (error) {
            console.error('Error fetching PAPRD balance:', error);
            return { balance: '0' };
        }
    }

    // 💰 Get BNM balance for address
    async getBnmBalance(address) {
        try {
            const response = await fetch(`${this.bnmApiUrl}/balance/${address}`);
            return await response.json();
        } catch (error) {
            console.error('Error fetching BNM balance:', error);
            return { balance: 0 };
        }
    }

    // 🏦 Get combined wallet data
    async getWalletData(address) {
        try {
            const [paprdBalance, bnmBalance] = await Promise.all([
                this.getPaprdBalance(address),
                this.getBnmBalance(address)
            ]);

            return {
                address,
                balances: {
                    BNM: {
                        balance: bnmBalance.balance || 0,
                        symbol: 'BNM',
                        name: 'Binomena Token'
                    },
                    PAPRD: {
                        balance: paprdBalance.balance || 0,
                        symbol: 'PAPRD',
                        name: 'Paper Dollar Stablecoin'
                    }
                },
                totalValueUSD: (parseFloat(bnmBalance.balance || 0) + parseFloat(paprdBalance.balance || 0))
            };
        } catch (error) {
            console.error('Error fetching wallet data:', error);
            return null;
        }
    }

    // 🔄 Transfer PAPRD tokens
    async transferPaprd(from, to, amount, privateKey) {
        try {
            const response = await fetch(`${this.apiUrl}/paprd/transfer`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    from,
                    to,
                    amount,
                    privateKey
                })
            });

            const result = await response.json();
            if (!response.ok) {
                throw new Error(result.error || 'Transfer failed');
            }
            return result;
        } catch (error) {
            console.error('Error transferring PAPRD:', error);
            throw error;
        }
    }

    // 🔄 Transfer BNM tokens (existing functionality)
    async transferBnm(from, to, amount, privateKey) {
        try {
            const response = await fetch(`${this.bnmApiUrl}/transfer`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    from,
                    to,
                    amount,
                    privateKey
                })
            });

            return await response.json();
        } catch (error) {
            console.error('Error transferring BNM:', error);
            throw error;
        }
    }

    // 📋 Get transaction history
    async getTransactionHistory(address) {
        try {
            const response = await fetch(`${this.apiUrl}/paprd/transactions/${address}`);
            return await response.json();
        } catch (error) {
            console.error('Error fetching transactions:', error);
            return { transactions: [] };
        }
    }
}

// Initialize PAPRD API
const paprdAPI = new PaprdAPI();

// 🎨 UI Components for PAPRD Integration

// Enhanced Wallet Balance Display
function updateWalletDisplay(address) {
    paprdAPI.getWalletData(address).then(walletData => {
        if (!walletData) return;

        const walletContainer = document.getElementById('wallet-balance') || createWalletContainer();
        
        walletContainer.innerHTML = `
            <div class="wallet-header">
                <h3>💰 Your Wallet</h3>
                <p class="wallet-address">${address.substring(0, 10)}...${address.substring(address.length - 8)}</p>
            </div>
            
            <div class="token-balances">
                <div class="token-balance bnm">
                    <div class="token-icon">🪙</div>
                    <div class="token-info">
                        <h4>Binomena Token</h4>
                        <p class="balance">${parseFloat(walletData.balances.BNM.balance).toLocaleString()} BNM</p>
                        <p class="value">$${parseFloat(walletData.balances.BNM.balance).toLocaleString()}</p>
                    </div>
                </div>
                
                <div class="token-balance paprd">
                    <div class="token-icon">💵</div>
                    <div class="token-info">
                        <h4>Paper Dollar Stablecoin</h4>
                        <p class="balance">${parseFloat(walletData.balances.PAPRD.balance).toLocaleString()} PAPRD</p>
                        <p class="value">$${parseFloat(walletData.balances.PAPRD.balance).toLocaleString()}</p>
                    </div>
                </div>
            </div>
            
            <div class="wallet-total">
                <h4>Total Value: $${walletData.totalValueUSD.toLocaleString()}</h4>
            </div>
            
            <div class="wallet-actions">
                <button onclick="showSendModal('BNM')" class="btn-send">Send BNM</button>
                <button onclick="showSendModal('PAPRD')" class="btn-send">Send PAPRD</button>
            </div>
        `;
    });
}

// Create wallet container if it doesn't exist
function createWalletContainer() {
    const container = document.createElement('div');
    container.id = 'wallet-balance';
    container.className = 'wallet-container';
    
    // Add to existing UI
    const mainContent = document.querySelector('.main-content') || document.body;
    mainContent.appendChild(container);
    
    return container;
}

// Send Token Modal
function showSendModal(tokenType) {
    const modal = document.createElement('div');
    modal.className = 'send-modal';
    modal.innerHTML = `
        <div class="modal-content">
            <div class="modal-header">
                <h3>Send ${tokenType} Tokens</h3>
                <button onclick="closeSendModal()" class="close-btn">×</button>
            </div>
            
            <form id="send-form">
                <div class="form-group">
                    <label>To Address:</label>
                    <input type="text" id="send-to" placeholder="Enter recipient address" required>
                </div>
                
                <div class="form-group">
                    <label>Amount:</label>
                    <input type="number" id="send-amount" placeholder="0.00" min="0" step="0.01" required>
                </div>
                
                <div class="form-group">
                    <label>Private Key:</label>
                    <input type="password" id="send-private-key" placeholder="Your private key" required>
                </div>
                
                <div class="form-actions">
                    <button type="button" onclick="closeSendModal()">Cancel</button>
                    <button type="submit" onclick="sendTokens('${tokenType}')">Send ${tokenType}</button>
                </div>
            </form>
        </div>
    `;
    
    document.body.appendChild(modal);
}

// Close send modal
function closeSendModal() {
    const modal = document.querySelector('.send-modal');
    if (modal) {
        modal.remove();
    }
}

// Send tokens function
async function sendTokens(tokenType) {
    event.preventDefault();
    
    const to = document.getElementById('send-to').value;
    const amount = document.getElementById('send-amount').value;
    const privateKey = document.getElementById('send-private-key').value;
    
    if (!to || !amount || !privateKey) {
        alert('Please fill in all fields');
        return;
    }
    
    // Get sender address (you'll need to implement this based on your wallet system)
    const from = getCurrentWalletAddress(); // Implement this function
    
    try {
        let result;
        if (tokenType === 'PAPRD') {
            result = await paprdAPI.transferPaprd(from, to, amount, privateKey);
        } else {
            result = await paprdAPI.transferBnm(from, to, amount, privateKey);
        }
        
        if (result.success) {
            alert(`Successfully sent ${amount} ${tokenType} to ${to}`);
            closeSendModal();
            // Refresh wallet display
            updateWalletDisplay(from);
        } else {
            alert(`Transfer failed: ${result.error}`);
        }
    } catch (error) {
        alert(`Transfer failed: ${error.message}`);
    }
}

// Token selection for send/receive
function createTokenSelector() {
    return `
        <div class="token-selector">
            <button class="token-option" onclick="selectToken('BNM')">
                <span class="token-icon">🪙</span>
                <span class="token-name">BNM</span>
            </button>
            <button class="token-option" onclick="selectToken('PAPRD')">
                <span class="token-icon">💵</span>
                <span class="token-name">PAPRD</span>
            </button>
        </div>
    `;
}

// Transaction history display
async function showTransactionHistory(address) {
    const paprdTx = await paprdAPI.getTransactionHistory(address);
    
    const historyContainer = document.createElement('div');
    historyContainer.className = 'transaction-history';
    historyContainer.innerHTML = `
        <h3>📋 Recent Transactions</h3>
        <div class="transactions-list">
            ${paprdTx.transactions.map(tx => `
                <div class="transaction-item">
                    <div class="tx-type">${tx.type === 'transfer' ? '🔄' : '🪙'}</div>
                    <div class="tx-details">
                        <p class="tx-amount">${tx.amount} PAPRD</p>
                        <p class="tx-address">${tx.type === 'transfer' ? 
                            (tx.from === address ? `To: ${tx.to.substring(0, 10)}...` : `From: ${tx.from.substring(0, 10)}...`) 
                            : 'Minted'}</p>
                        <p class="tx-time">${new Date(tx.timestamp).toLocaleString()}</p>
                    </div>
                    <div class="tx-status">${tx.status === 'confirmed' ? '✅' : '⏳'}</div>
                </div>
            `).join('')}
        </div>
    `;
    
    return historyContainer;
}

// Add CSS styles for PAPRD UI components
function addPaprdStyles() {
    const style = document.createElement('style');
    style.textContent = `
        .wallet-container {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            border-radius: 15px;
            padding: 25px;
            margin: 20px 0;
            color: white;
            box-shadow: 0 10px 30px rgba(0,0,0,0.3);
        }
        
        .wallet-header h3 {
            margin: 0 0 10px 0;
            font-size: 1.5em;
        }
        
        .wallet-address {
            opacity: 0.8;
            font-family: monospace;
        }
        
        .token-balances {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 20px;
            margin: 20px 0;
        }
        
        .token-balance {
            background: rgba(255,255,255,0.1);
            border-radius: 10px;
            padding: 20px;
            display: flex;
            align-items: center;
            gap: 15px;
        }
        
        .token-icon {
            font-size: 2em;
        }
        
        .token-info h4 {
            margin: 0;
            font-size: 0.9em;
            opacity: 0.8;
        }
        
        .balance {
            font-size: 1.2em;
            font-weight: bold;
            margin: 5px 0;
        }
        
        .value {
            opacity: 0.7;
            font-size: 0.9em;
        }
        
        .wallet-total {
            text-align: center;
            border-top: 1px solid rgba(255,255,255,0.2);
            padding-top: 20px;
        }
        
        .wallet-actions {
            display: flex;
            gap: 10px;
            margin-top: 20px;
        }
        
        .btn-send {
            flex: 1;
            background: rgba(255,255,255,0.2);
            border: none;
            padding: 12px;
            border-radius: 8px;
            color: white;
            cursor: pointer;
            transition: all 0.3s;
        }
        
        .btn-send:hover {
            background: rgba(255,255,255,0.3);
            transform: translateY(-2px);
        }
        
        .send-modal {
            position: fixed;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;
            background: rgba(0,0,0,0.5);
            display: flex;
            justify-content: center;
            align-items: center;
            z-index: 1000;
        }
        
        .modal-content {
            background: white;
            border-radius: 15px;
            padding: 30px;
            width: 90%;
            max-width: 500px;
            box-shadow: 0 20px 50px rgba(0,0,0,0.3);
        }
        
        .modal-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 25px;
        }
        
        .close-btn {
            background: none;
            border: none;
            font-size: 1.5em;
            cursor: pointer;
        }
        
        .form-group {
            margin-bottom: 20px;
        }
        
        .form-group label {
            display: block;
            margin-bottom: 5px;
            font-weight: bold;
        }
        
        .form-group input {
            width: 100%;
            padding: 12px;
            border: 2px solid #ddd;
            border-radius: 8px;
            font-size: 1em;
        }
        
        .form-actions {
            display: flex;
            gap: 10px;
            justify-content: flex-end;
        }
        
        .form-actions button {
            padding: 12px 24px;
            border: none;
            border-radius: 8px;
            cursor: pointer;
            font-size: 1em;
        }
        
        .form-actions button[type="submit"] {
            background: #667eea;
            color: white;
        }
        
        .transaction-history {
            background: white;
            border-radius: 15px;
            padding: 25px;
            margin: 20px 0;
            box-shadow: 0 5px 15px rgba(0,0,0,0.1);
        }
        
        .transaction-item {
            display: flex;
            align-items: center;
            padding: 15px;
            border-bottom: 1px solid #eee;
            gap: 15px;
        }
        
        .tx-type {
            font-size: 1.5em;
        }
        
        .tx-details {
            flex: 1;
        }
        
        .tx-amount {
            font-weight: bold;
            margin: 0;
        }
        
        .tx-address, .tx-time {
            margin: 5px 0;
            opacity: 0.7;
            font-size: 0.9em;
        }
    `;
    
    document.head.appendChild(style);
}

// Initialize PAPRD integration when page loads
document.addEventListener('DOMContentLoaded', () => {
    addPaprdStyles();
    
    // Check if user has a wallet connected
    const connectedWallet = getCurrentWalletAddress();
    if (connectedWallet) {
        updateWalletDisplay(connectedWallet);
    }
});

// Helper function to get current wallet address
function getCurrentWalletAddress() {
    // Implement this based on your existing wallet system
    // This might be stored in localStorage, sessionStorage, or a global variable
    return localStorage.getItem('connectedWallet') || 
           window.currentWallet || 
           'AdNe6c3ce54e4371d056c7c566675ba16909eb2e9534'; // Your founder address as default
}

// Export for use in other parts of your application
window.PaprdAPI = paprdAPI;
window.updateWalletDisplay = updateWalletDisplay;
window.showSendModal = showSendModal; 