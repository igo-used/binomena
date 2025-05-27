# 🚀 SuperNom - Revolutionary Blockchain VPN

**World's First Blockchain-Based VPN Service - Mobile Network Operators (MNO) Killer**

SuperNom replaces traditional ISPs and mobile carriers with blockchain-powered internet access. Pay with BNM tokens instead of monthly subscriptions, access global internet without restrictions, and own your digital privacy.

**Founded by Juxhino Kapllanaj (uJ1N0)** - Creator of Binomena Blockchain and pioneer of decentralized telecommunications.

---

## 🌟 Revolutionary Features

### 📡 **Replace Your SIM Card with Your Wallet**
- **No SIM cards needed** - Your blockchain wallet IS your identity
- **Pay-per-use access** - No monthly subscriptions or contracts
- **Global roaming** - Works everywhere without carrier restrictions
- **Complete privacy** - No personal data collection or tracking
- **Instant activation** - Buy BNM tokens, get internet immediately

### 🔐 **Privacy-First Technology**
- **Zero-logs architecture** - No browsing history stored
- **End-to-end encryption** - WireGuard protocol for maximum security
- **Blockchain verification** - All payments transparent and immutable
- **Censorship resistant** - Access the open internet globally
- **Democratic governance** - Community-controlled network

### 💰 **Token Economics**
- **Basic**: 10 BNM per hour (5GB bandwidth)
- **Standard**: 50 BNM per 6 hours (15GB bandwidth)  
- **Premium**: 150 BNM per 24 hours (50GB bandwidth)
- **Staking rewards** - Lock BNM tokens for better rates
- **Token appreciation** - Value increases with network growth

---

## 🏗️ Architecture

SuperNom consists of three main components that work together to create the world's first blockchain ISP:

### 🔗 **Smart Contracts** (`contracts/`)
- **VPNAccessContract** - Manages VPN purchases and sessions
- **GovernanceContract** - Democratic voting and proposals
- **ReputationContract** - User trust scores and staking

### 🌐 **API Gateway** (`api/`)
- RESTful API for all VPN operations
- CORS-enabled for web interface integration
- Real-time session management
- Payment processing and verification

### 💻 **Web Interface** (`client/` & `public/`)
- Professional marketing website
- Interactive VPN demo interface
- Live blockchain integration
- User-friendly token purchasing

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Web Client    │    │   API Gateway   │    │ Smart Contracts │
│  (React/HTML)   │◄──►│   (Go/Gin)      │◄──►│  (Rust/WASM)    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  User Payments  │    │ VPN Servers     │    │ Binomena Chain │
│   (BNM Tokens)  │◄──►│  (WireGuard)    │◄──►│  (Blockchain)   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

---

## 🚀 Quick Start

### Prerequisites
- Go 1.18+ (for API server)
- Rust 1.70+ (for smart contracts)
- BNM tokens (for VPN access)

### Installation

```bash
# Clone SuperNom
git clone https://github.com/igo-used/binomena.git
cd binomena/supernom

# Install dependencies
go mod download

# Build SuperNom
go build -o supernom-vpn

# Start the service
./supernom-vpn
```

### Quick Test

```bash
# Check system status
curl http://localhost:8080/health

# Check VPN access for wallet
curl "http://localhost:8080/auth/check?wallet=YOUR_WALLET_ADDRESS"

# Purchase VPN access
curl -X POST http://localhost:8080/auth/purchase \
  -H "Content-Type: application/json" \
  -d '{
    "walletAddress": "YOUR_WALLET",
    "tokenType": "BNM",
    "amount": 10,
    "duration": 3600,
    "geographicZone": "global",
    "paymentSignature": "demo_signature"
  }'
```

---

## 📋 API Endpoints

### 🔐 Authentication
- `GET /auth/check` - Verify VPN access for wallet
- `POST /auth/purchase` - Buy VPN access with BNM tokens
- `GET /auth/config` - Download WireGuard configuration
- `POST /auth/revoke` - Terminate VPN session

### 👤 User Management  
- `GET /user/reputation` - Get user trust score
- `GET /user/sessions` - View session history
- `POST /user/stake` - Stake BNM tokens for better rates

### 🗳️ Governance
- `GET /governance/proposals` - View community proposals
- `POST /governance/vote` - Vote on proposals (delegates only)
- `POST /governance/emergency` - Trigger emergency actions

### 📊 System Status
- `GET /health` - Health check
- `GET /status` - System operational status  
- `GET /stats` - Detailed network statistics

---

## 💡 Use Cases

### 🌍 **Global Digital Nomads**
- Work from anywhere with instant internet access
- No need for local SIM cards or contracts
- Pay only for what you use

### 🔒 **Privacy Advocates**
- Complete anonymity with crypto payments
- No personal data required
- Censorship-resistant browsing

### 🏢 **Remote Teams**
- Secure connections for distributed workers
- Enterprise-grade encryption
- Scalable token purchasing

### 🚀 **Crypto Enthusiasts**
- Native blockchain integration
- Token appreciation potential
- Community governance participation

---

## 🛡️ Security & Privacy

### 🔒 **Technical Security**
- WireGuard VPN protocol (state-of-the-art)
- ECDSA P-256 cryptographic signatures
- Zero-logs architecture (nothing stored)
- Perfect forward secrecy
- End-to-end encryption

### 📊 **Blockchain Security**
- Smart contract-based access control
- Immutable payment records
- Democratic governance
- Emergency pause mechanisms
- Multi-signature controls

### 🌐 **Operational Security**
- RAM-only VPN servers
- No user data retention
- Jurisdiction-aware compliance
- Regular security audits
- Open source transparency

---

## 🌟 Future Roadmap

### 📱 **Phase 1: VPN Service (Active)**
- ✅ Blockchain VPN infrastructure
- ✅ Smart contract integration
- ✅ Web interface and API
- 🔄 Beta testing and optimization

### 🌐 **Phase 2: ISP Replacement (2025)**
- 📋 Direct internet access without traditional ISPs
- 📋 Mesh networking with token incentives
- 📋 Mobile app for seamless usage
- 📋 Enterprise partnerships

### 🛰️ **Phase 3: Satellite Integration (2026)**
- 📋 Starlink and satellite internet support
- 📋 6G infrastructure tokenization
- 📋 Global mesh network deployment
- 📋 Complete MNO disruption

### 🚀 **Phase 4: Telecommunications Revolution (2027+)**
- 📋 Replace all traditional carriers
- 📋 Global decentralized internet
- 📋 Token-based economy standard
- 📋 Complete digital sovereignty

---

## 🤝 Contributing

SuperNom is open source and welcomes contributions from the community!

### 🔧 **Development Areas**
- Smart contract optimization
- VPN server infrastructure  
- Mobile app development
- Security auditing
- Documentation and tutorials

### 📋 **Contribution Guidelines**
1. **Code Quality** - Follow Go and Rust best practices
2. **Security First** - All changes must maintain security standards
3. **Privacy Focus** - Uphold zero-logs and anonymity principles
4. **Testing** - Add comprehensive tests for new features
5. **Documentation** - Update docs for all changes

---

## 📜 License

This project is licensed under the **Apache License 2.0** - see the [LICENSE](LICENSE) file for details.

### Copyright Notice
```
Copyright 2025 Juxhino Kapllanaj (uJ1N0)
Licensed under the Apache License, Version 2.0
```

---

## 📞 Contact & Community

- **Founder**: Juxhino Kapllanaj (uJ1N0)
- **Binomena Blockchain**: https://github.com/igo-used/binomena
- **Email**: juxhino.kap@yahoo.com
- **Community**: https://x.com/BinomChain

---

## 🎯 Why SuperNom Will Win

### 💰 **Economic Advantage**
- **Traditional ISP**: $50-100/month subscriptions
- **SuperNom**: $1-10 pay-per-use with token appreciation

### 🌍 **Global Accessibility**
- **Traditional Roaming**: $5-20 per MB in foreign countries
- **SuperNom**: Same low rates everywhere globally

### 🔐 **Privacy Superiority**
- **Traditional ISPs**: Track, log, and sell your data
- **SuperNom**: Zero logs, complete anonymity, blockchain verified

### 🚀 **Technology Leadership**
- **First mover** in blockchain ISP space
- **Proven technology** with live mainnet
- **Community governed** vs corporate controlled
- **Token economics** creating sustainable growth

---

*Built with ❤️ by Juxhino Kapllanaj (uJ1N0) - Disrupting telecommunications one block at a time* 