# 🚀 SuperNom - Decentralized VPN on Binomena Blockchain

**SuperNom** is a revolutionary decentralized VPN system that uses the Binomena blockchain as the access control layer, replacing traditional SIM card trust models with blockchain-based authentication. This is the future of internet access - where your wallet becomes your identity and tokens become your access pass.

## 🌟 Features

### **Core Functionality**
- 🔐 **Blockchain Authentication**: Pay with BNM tokens to access VPN services
- 🏆 **Reputation System**: Build trust through successful sessions and token staking
- 🗳️ **Democratic Governance**: Delegate voting for blacklisting, whitelisting, and system changes
- ⚡ **Smart Contracts**: Automated session management and payment processing
- 🌍 **Geographic Coverage**: Multiple VPN nodes across different zones
- 📊 **Tiered Access**: Basic, Standard, and Premium access levels

### **Security Features**
- 🛡️ **Anti-Abuse System**: Reputation-based access control and rate limiting
- 🚨 **Emergency Controls**: Delegate-triggered emergency shutdown capabilities
- 🔒 **Stake Requirements**: New users must stake tokens for access
- 📝 **Audit Trail**: All transactions and access logged on blockchain
- 🚫 **Blacklist Management**: Community-governed ban system with appeals

### **Technical Architecture**
- 🔗 **Binomena Integration**: Seamless integration with existing DPoS blockchain
- 🌐 **WireGuard VPN**: Modern, secure VPN protocol
- 🎯 **RESTful API**: Easy integration for developers and applications
- 📱 **Web Client**: Beautiful web interface for managing VPN access
- 🔄 **Real-time Sync**: Continuous synchronization with blockchain state

## 🏗️ Architecture Overview

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Web Client    │───▶│  SuperNom API   │───▶│ Smart Contracts │
│   (HTML/JS)     │    │   (Gateway)     │    │  (VPN + Gov)    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │                        │
                                ▼                        ▼
                       ┌─────────────────┐    ┌─────────────────┐
                       │ VPN Infrastructure│    │ Binomena Chain │
                       │   (WireGuard)   │    │   (DPoS + WASM) │
                       └─────────────────┘    └─────────────────┘
```

## 🚦 Quick Start

### **Prerequisites**
- Go 1.19+
- Access to Binomena blockchain node
- WireGuard installed (for VPN usage)

### **Installation**

1. **Clone and Setup**
```bash
cd binomena/supernom
go mod download
```

2. **Configure Environment**
```bash
export SUPERNOM_PORT=8080
export BINOMENA_NODE_URL=https://binomena-node.onrender.com
```

3. **Start SuperNom**
```bash
go run main.go
```

4. **Access Web Interface**
```
http://localhost:8080/health
http://localhost:8080/status
```

## 📋 API Endpoints

### **Authentication & Access**
```
GET  /auth/check?wallet=AdNe...        - Check VPN authorization
POST /auth/purchase                    - Purchase VPN access
GET  /auth/config?wallet=AdNe...       - Get WireGuard config
POST /auth/revoke                      - Revoke VPN session
GET  /auth/status?wallet=AdNe...       - Get session status
```

### **User Management**
```
GET  /user/reputation?address=AdNe...  - Get user reputation
GET  /user/sessions?wallet=AdNe...     - Get session history
POST /user/stake                       - Stake tokens for reputation
```

### **Governance**
```
GET  /governance/proposals             - List governance proposals
POST /governance/proposal              - Create new proposal
POST /governance/vote                  - Vote on proposal
POST /governance/emergency             - Trigger emergency shutdown
```

### **System Status**
```
GET  /status                           - System status and stats
GET  /stats                            - Detailed system metrics
GET  /health                           - Health check
```

## 💰 Token Economics

### **Access Pricing (BNM)**
- **Basic (1 hour)**: 10 BNM - 5GB bandwidth
- **Standard (6 hours)**: 50 BNM - 15GB bandwidth  
- **Premium (24 hours)**: 150 BNM - 50GB bandwidth

### **Reputation System**
- **New Users**: Start with 100 trust score, require 100 BNM stake
- **Verified Users**: 400+ trust score, standard access
- **Trusted Users**: 800+ trust score, premium benefits
- **Violations**: -50 trust score penalty per violation

### **Governance Voting Thresholds**
- **Blacklist**: 60% delegate majority
- **Whitelist**: 51% delegate majority
- **Settings**: 67% delegate majority
- **Emergency**: 75% delegate majority

## 🔧 Smart Contract Architecture

### **VPN Access Contract**
```go
type VPNAccessContract struct {
    Sessions       map[string]*VPNSession
    Reputation     map[string]*UserReputation
    GlobalSettings *GlobalVPNSettings
    Blacklist      map[string]*BlacklistEntry
}
```

**Key Functions:**
- `PurchaseAccess()` - Handle VPN payments and create sessions
- `CheckAuthorization()` - Verify active VPN access
- `UpdateBandwidthUsage()` - Track and limit bandwidth
- `CompleteSession()` - Update reputation on successful completion

### **Governance Contract**
```go
type GovernanceContract struct {
    Proposals           map[string]*Proposal
    AuthorizedDelegates map[string]*Delegate
    EmergencySettings   *EmergencySettings
}
```

**Key Functions:**
- `CreateProposal()` - Submit governance proposals
- `CastVote()` - Delegate voting on proposals
- `ExecuteProposal()` - Apply approved changes
- `TriggerEmergency()` - Emergency system controls

## 🌐 Usage Examples

### **Purchase VPN Access**
```javascript
const response = await fetch('/auth/purchase', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
        walletAddress: 'AdNe1234567890...',
        tokenType: 'BNM',
        amount: 10,
        duration: 3600, // 1 hour
        geographicZone: 'us-east',
        paymentSignature: 'blockchain_tx_signature'
    })
});
```

### **Get VPN Configuration**
```javascript
const config = await fetch('/auth/config?wallet=AdNe1234567890...');
// Returns WireGuard configuration for direct import
```

### **Check Authorization**
```javascript
const auth = await fetch('/auth/check?wallet=AdNe1234567890...');
// Returns session details and authorization status
```

## 🔒 Security Considerations

### **Anti-Abuse Measures**
1. **Stake Requirements**: New users must stake 100 BNM minimum
2. **Session Limits**: Maximum 5 sessions per day per address
3. **Reputation Decay**: Inactive users lose trust score over time
4. **Geographic Restrictions**: Configurable country/region blocking
5. **Emergency Shutdown**: Delegate-controlled emergency stops

### **Privacy Protection**
- **No Deep Packet Inspection**: Traffic content never analyzed
- **Minimal Data Storage**: Only session metadata stored
- **Temporary Logs**: Connection logs deleted after 24-48 hours
- **Address Anonymity**: No KYC required, only wallet address

### **Compliance Framework**
- **Terms of Service**: Blockchain-signed user agreements
- **Prohibited Activities**: Clear usage guidelines
- **Law Enforcement**: Cooperation when legally required
- **Appeal Process**: Transparent ban appeals through governance

## 🎯 Future Roadmap

### **Phase 1: MVP (Current)**
- ✅ Basic VPN access via blockchain payments
- ✅ Smart contract governance
- ✅ Web client interface
- ✅ Reputation system

### **Phase 2: Enhanced Infrastructure**
- 🔲 Multiple VPN node deployment
- 🔲 Real WireGuard key generation
- 🔲 Mobile app development
- 🔲 QR code configuration sharing

### **Phase 3: 6G Integration**
- 🔲 Satellite internet compatibility
- 🔲 6G network integration
- 🔲 IoT device support
- 🔲 Global mesh networking

### **Phase 4: Ecosystem Expansion**
- 🔲 Delegate VPN node operators
- 🔲 Revenue sharing for node operators
- 🔲 Cross-chain compatibility
- 🔲 Enterprise features

## 🤝 Contributing

SuperNom is part of the Binomena ecosystem. To contribute:

1. **Report Issues**: Use GitHub issues for bugs and feature requests
2. **Code Contributions**: Submit pull requests with tests
3. **Governance**: Participate in delegate voting for system changes
4. **Documentation**: Help improve user guides and API docs

## 📜 License

SuperNom is open source software licensed under MIT License.

## 🆘 Support

- **Documentation**: See `/docs` directory
- **API Reference**: Available at `/api/docs` when running
- **Community**: Join Binomena Discord/Telegram
- **Issues**: Report bugs on GitHub

---

**SuperNom**: *Revolutionizing internet access through blockchain technology* 🚀

*Built with ❤️ on the Binomena blockchain* 