# Solana Crowdfunding Rust Client

A command-line interface written in Rust for interacting with the Solana crowdfunding program.

## Features

- 🦀 **Rust Performance**: Fast and memory-efficient CLI
- 🔐 **Wallet Management**: Secure keypair handling
- 💰 **Balance Checking**: Real-time SOL balance display
- 📋 **Campaign Management**: Create and manage crowdfunding campaigns
- 💸 **Donations**: Donate SOL to campaigns
- 💵 **Withdrawals**: Withdraw funds (campaign admin only)

## Setup

1. **Create Wallet File**:
   ```bash
   cd rust_client
   # Copy the example wallet template
   cp wallet.example.json wallet.json
   
   # Edit wallet.json with your actual keypair information
   # You can generate a keypair using: solana-keygen new
   ```

2. **Build the Application**:
   ```bash
   cargo build --release
   ```

3. **Run the Application**:
   ```bash
   cargo run
   # or run the built binary
   ./target/release/rust_client
   ```

   > **Note**: The wallet.json file is required and must contain your Solana keypair.

## Usage

### First Time Setup

1. **Prepare Wallet**: Create `wallet.json` with your keypair (see Setup above)
2. **Build & Run**: Use `cargo run` to start the application
3. **Request Airdrop**: Get SOL for transaction fees (devnet only)
4. **Create Campaign**: Create your first campaign
5. **Interact**: Donate to or withdraw from campaigns

## Program Details

- **Program ID**: `7XJkGrdSHn3chc7rsv1xdZEKtwP9w5rSx1sHohzM5skv`
- **Network**: Solana Devnet
- **RPC Endpoint**: Solana devnet RPC

## Security Notes

⚠️ **CRITICAL SECURITY WARNINGS**:
- **Never commit wallet files**: Keep `wallet.json` secure and never commit to git
- **Devnet only**: This setup is for development/testing on devnet only
- **Private key safety**: Your private key grants full access to your wallet
- **Backup important**: Back up your wallet file securely
- **Production use**: For mainnet, use hardware wallets and proper key management

## Development

```bash
# Format code
cargo fmt

# Run tests
cargo test

# Check for issues
cargo clippy

# Build optimized release
cargo build --release
```

---

**Happy crowdfunding with Rust!** 🦀🚀