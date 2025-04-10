### Binomena Blockchain

Binomena is a next-generation blockchain platform featuring NodeSwift consensus (a reputation-based Proof of Stake system), WebAssembly smart contracts, and a native BNM token with deflationary economics.

## Features

- **NodeSwift Consensus**: Efficient Proof of Stake with validator reputation scoring
- **WebAssembly Smart Contracts**: Write contracts in AssemblyScript, Rust, or any WASM-compatible language
- **Integrated Security**: Built-in audit framework for continuous blockchain validation
- **BNM Token**: Native cryptocurrency with fixed supply and deflationary mechanism
- **P2P Network**: Robust peer-to-peer communication using libp2p
- **Comprehensive API**: RESTful interface for all blockchain operations


## Technology Stack

- **Backend**: Go
- **Smart Contracts**: AssemblyScript/WebAssembly
- **Networking**: libp2p
- **WebAssembly Runtime**: Wasmer
- **API Framework**: Gin
- **Cryptography**: ECDSA with P-256 curve


## Installation

### Prerequisites

- Go 1.18 or higher
- Node.js 16 or higher (for smart contract development)


### Building from Source

```shellscript
# Clone the repository
git clone https://github.com/your-org/binomena.git
cd binomena

# Build the blockchain
go build -o binomena

# Install smart contract dependencies
cd contracts
npm install
```

## Running a Node

```shellscript
# Start a single node
./binomena --api-port 8080 --p2p-port 9000 --id "my-node"

# Start a node with a bootstrap peer
./binomena --api-port 8081 --p2p-port 9001 --bootstrap "/ip4/127.0.0.1/tcp/9000/p2p/QmBootstrapPeerID" --id "peer-node"
```

## Smart Contract Development

```shellscript
# Navigate to contracts directory
cd contracts

# Create a new contract
mkdir -p assembly/mycontract
touch assembly/mycontract/index.ts

# Build the contract
npm run asbuild
```

Example contract (AssemblyScript):

```typescript
// assembly/mycontract/index.ts
export function add(a: i32, b: i32): i32 {
  return a + b;
}

export function getGreeting(name: string): string {
  return "Hello, " + name + "!";
}
```

## API Examples

### Create a Wallet

```shellscript
curl -X POST http://localhost:8080/wallet
```

### Check Balance

```shellscript
curl -X GET http://localhost:8080/balance/AdNe123456789...
```

### Submit a Transaction

```shellscript
curl -X POST http://localhost:8080/transaction \
  -H "Content-Type: application/json" \
  -d '{
    "from": "AdNe123456789...",
    "to": "AdNe987654321...",
    "amount": 100.0,
    "privateKey": "your-private-key"
  }'
```

### Deploy a Smart Contract

```shellscript
curl -X POST http://localhost:8080/contracts/deploy \
  -H "Content-Type: application/json" \
  -d '{
    "owner": "AdNe123456789...",
    "name": "MyContract",
    "code": "base64-encoded-wasm-binary",
    "fee": 1.0,
    "privateKey": "your-private-key"
  }'
```

## Project Structure


## Testing

```shellscript
# Run all tests
./run_tests.sh

# Run specific module tests
go test ./wallet -v
go test ./core -v
go test ./smartcontract -v
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request


## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## White Paper

For a detailed technical overview of Binomena, please refer to our [white paper](https://example.com/binomena-whitepaper.pdf).

## Contact

- Website: [https://binomena.io]([https://binomena.io](https://binomena.com/))
- Email: [info@binomena.io](team@binomena.com)
- Founder Email: [@binomena](juxhino.kap@yahoo.com)
- Telegram: [Binomena Community](https://t.me/binomchain)
