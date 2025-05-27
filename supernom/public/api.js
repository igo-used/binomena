// SuperNom API Configuration
const SUPERNOM_API = 'https://supernom-vpn.onrender.com';
const BINOMENA_API = 'https://binomena-node.onrender.com';

// SuperNom API Helper Functions
class SuperNomAPI {
    static async checkHealth() {
        try {
            const response = await fetch(`${SUPERNOM_API}/health`);
            return await response.json();
        } catch (error) {
            console.error('API Error:', error);
            return null;
        }
    }
    
    static async getSystemStatus() {
        try {
            const response = await fetch(`${SUPERNOM_API}/status`);
            return await response.json();
        } catch (error) {
            console.error('API Error:', error);
            return null;
        }
    }
    
    static async checkVPNAuth(walletAddress) {
        try {
            const response = await fetch(`${SUPERNOM_API}/auth/check?wallet=${walletAddress}`);
            return await response.json();
        } catch (error) {
            console.error('API Error:', error);
            return null;
        }
    }
    
    static async purchaseVPN(walletAddress, amount, duration, zone, signature) {
        try {
            const response = await fetch(`${SUPERNOM_API}/auth/purchase`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    walletAddress,
                    tokenType: 'BNM',
                    amount,
                    duration,
                    geographicZone: zone,
                    paymentSignature: signature
                })
            });
            return await response.json();
        } catch (error) {
            console.error('API Error:', error);
            return null;
        }
    }
    
    static async getBNMBalance(address) {
        try {
            const response = await fetch(`${BINOMENA_API}/balance/${address}`);
            return await response.json();
        } catch (error) {
            console.error('API Error:', error);
            return null;
        }
    }
}

// Initialize when page loads
document.addEventListener('DOMContentLoaded', async function() {
    console.log('🚀 SuperNom API Ready');
    
    // Check system status
    const status = await SuperNomAPI.getSystemStatus();
    if (status && status.success) {
        console.log('✅ SuperNom VPN Online:', status.data);
        
        // Update live stats on page if elements exist
        const statsElements = document.querySelectorAll('[data-stat]');
        statsElements.forEach(el => {
            const stat = el.dataset.stat;
            if (status.data.vpn && status.data.vpn[stat] !== undefined) {
                el.textContent = status.data.vpn[stat];
            }
        });
    }
}); 