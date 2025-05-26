# Paper Dollar Stablecoin - Rust Smart Contract

A decentralized stablecoin implementation written in Rust, compiled to WebAssembly for blockchain deployment. This contract provides a collateral-backed stablecoin with advanced features including multi-collateral support, minting controls, and administrative functions.

## Features

### Core Functionality
- **ERC20-like Token**: Standard transfer, mint, and burn operations
- **Collateral Management**: Support for multiple collateral types (Fiat and BINOM tokens)
- **Collateralization Ratio**: Configurable over-collateralization requirements
- **Minting Control**: Role-based minting with collateral verification

### Administrative Features
- **Owner Controls**: Contract ownership with transfer capability
- **Minter Management**: Add/remove authorized minters
- **Blacklist System**: Block specific addresses from operations
- **Pause Mechanism**: Emergency pause/unpause functionality
- **Collateral Ratio Adjustment**: Dynamic collateral requirements

### Security Features
- **Safe Arithmetic**: Overflow/underflow protection
- **Access Control**: Role-based permission system
- **Input Validation**: Comprehensive parameter checking
- **Event Logging**: Full audit trail of operations

## Architecture

### Modules Structure
```
src/
├── lib.rs              # Main library and WASM exports
├── stablecoin.rs       # Core stablecoin contract logic
├── storage.rs          # Persistent storage abstraction
├── context.rs          # Execution context management
├── events.rs           # Event definitions and emission
└── errors.rs           # Error types and safe arithmetic
```

### Key Components

#### Storage Layer
- Persistent key-value storage abstraction
- Type-safe serialization with Borsh
- Efficient storage key management

#### Context System
- Execution context tracking (caller, block info)
- Thread-safe global state management
- Test utilities for context manipulation

#### Event System
- Comprehensive event logging
- Structured event data with timestamps
- Cross-platform event emission (WASM/native)

#### Error Handling
- Custom error types with descriptive messages
- Safe arithmetic operations with overflow protection
- Result-based error propagation

## Contract Interface

### View Functions

```rust
// Balance and supply information
get_balance(address: &str) -> u64
get_total_supply() -> u64

// Collateral information
get_collateral_balance(address: &str, collateral_type: CollateralType) -> u64
get_collateral_ratio() -> u64
get_fiat_reserve() -> u64

// Access control checks
is_blacklisted(address: &str) -> bool
is_minter(address: &str) -> bool
is_paused() -> bool
get_owner() -> String
```

### State-Changing Functions

```rust
// Token operations
transfer(to: &str, amount: u64) -> Result<bool>
mint(to: &str, amount: u64) -> Result<bool>          // Minters only
burn(amount: u64) -> Result<bool>

// Collateral management
add_collateral(amount: u64, collateral_type: CollateralType) -> Result<bool>
remove_collateral(amount: u64) -> Result<bool>

// Administrative functions (Owner only)
add_minter(address: &str) -> Result<bool>
remove_minter(address: &str) -> Result<bool>
blacklist(address: &str) -> Result<bool>
unblacklist(address: &str) -> Result<bool>
pause() -> Result<bool>
unpause() -> Result<bool>
set_collateral_ratio(ratio: u64) -> Result<bool>
transfer_ownership(new_owner: &str) -> Result<bool>
```

## Collateral Types

The contract supports two types of collateral:

### Fiat Collateral (Type 0)
- Traditional fiat currency backing
- Managed through fiat reserves
- Suitable for stable, low-volatility backing

### BINOM Token Collateral (Type 1)
- Cryptocurrency token backing
- Integration with BINOM token ecosystem
- Dynamic value based on token price

## Collateralization Requirements

- **Default Ratio**: 150% (configurable by owner)
- **Minimum Ratio**: 100% (enforced by contract)
- **Over-collateralization**: Required for all minting operations
- **Collateral Checks**: Enforced on minting and collateral removal

## Build and Deployment

### Prerequisites

1. **Rust Toolchain**:
   ```bash
   curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
   rustup target add wasm32-unknown-unknown
   ```

2. **wasm-pack** (for WebAssembly builds):
   ```bash
   curl https://rustwasm.github.io/wasm-pack/installer/init.sh -sSf | sh
   ```

### Building

#### Quick Build
```bash
./build.sh
```

#### Manual Build Options

**For Web Deployment**:
```bash
wasm-pack build --target web --out-dir pkg
```

**For Node.js**:
```bash
wasm-pack build --target nodejs --out-dir pkg-node
```

**For Testing**:
```bash
cargo test
```

**For Development**:
```bash
cargo check
cargo clippy
cargo fmt
```

### Build Outputs

- `pkg/`: Web-compatible WASM package
- `pkg-node/`: Node.js-compatible WASM package
- `target/`: Rust compilation artifacts

## Usage Examples

### JavaScript Integration (Web)

```javascript
import init, { ContractInstance } from './pkg/paper_dollar_stablecoin.js';

async function main() {
    await init();
    
    // Create contract instance
    const contract = new ContractInstance();
    
    // View operations
    const totalSupply = contract.get_total_supply();
    const balance = contract.get_balance("user_address");
    
    // State-changing operations
    try {
        await contract.add_collateral(1000, 0); // Add 1000 fiat collateral
        await contract.mint("recipient", 500);   // Mint 500 tokens
        await contract.transfer("recipient", 100); // Transfer 100 tokens
    } catch (error) {
        console.error("Transaction failed:", error);
    }
}
```

### Node.js Integration

```javascript
const { ContractInstance } = require('./pkg-node/paper_dollar_stablecoin.js');

const contract = new ContractInstance();

// Add collateral and mint tokens
contract.add_collateral(1500, 0); // 1500 fiat collateral
contract.mint("user1", 1000);     // Mint 1000 tokens

console.log("Total supply:", contract.get_total_supply());
console.log("User balance:", contract.get_balance("user1"));
```

## Testing

The contract includes comprehensive unit tests covering:

- **Core Functionality**: Transfer, mint, burn operations
- **Collateral Management**: Adding/removing collateral with ratio checks
- **Access Control**: Owner, minter, and blacklist permissions
- **Pause Mechanism**: Contract pause/unpause functionality
- **Edge Cases**: Overflow protection, invalid inputs, boundary conditions

### Running Tests

```bash
cargo test                    # Run all tests
cargo test --lib             # Run unit tests only
cargo test test_mint         # Run specific test
cargo test -- --nocapture   # Show test output
```

### Test Coverage

```bash
cargo install cargo-tarpaulin
cargo tarpaulin --out Html
```

## Security Considerations

### Implemented Protections

1. **Integer Overflow Protection**: Safe arithmetic operations throughout
2. **Access Control**: Strict role-based permissions
3. **Input Validation**: Comprehensive parameter checking
4. **Reentrancy Protection**: Single-threaded execution model
5. **Collateral Verification**: Enforced over-collateralization

### Audit Recommendations

1. **Formal Verification**: Consider mathematical proof of contract properties
2. **External Audit**: Professional security audit before mainnet deployment
3. **Gradual Rollout**: Start with limited functionality and expand
4. **Monitoring**: Implement real-time contract monitoring
5. **Emergency Procedures**: Test pause/unpause mechanisms thoroughly

## Roadmap

### Phase 1: Core Implementation ✅
- Basic token functionality
- Collateral management
- Access controls

### Phase 2: Enhanced Features 🚧
- Oracle integration for price feeds
- Dynamic collateral ratios
- Multi-signature support

### Phase 3: Ecosystem Integration 📋
- DEX integration
- Governance token support
- Cross-chain compatibility

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

### Code Style

- Follow Rust standard formatting (`cargo fmt`)
- Use clippy for linting (`cargo clippy`)
- Maintain test coverage above 90%
- Document public APIs

## License

MIT License - see LICENSE file for details.

## Support

For questions, issues, or contributions:
- Create an issue on GitHub
- Join our Discord community
- Check our documentation wiki

---

**Note**: This is a financial smart contract. Always conduct thorough testing and security audits before deploying to production networks. 