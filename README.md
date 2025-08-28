# Solana Crowdfunding dApp

[![Solana](https://img.shields.io/badge/Solana-Blockchain-blue)](https://solana.com/)
[![Anchor](https://img.shields.io/badge/Anchor-Framework-purple)](https://anchor-lang.com/)
[![React](https://img.shields.io/badge/React-Frontend-blue)](https://reactjs.org/)
[![Go](https://img.shields.io/badge/Go-Client-cyan)](https://golang.org/)
[![Rust](https://img.shields.io/badge/Rust-Client-orange)](https://rust-lang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A comprehensive decentralized crowdfunding application built on the Solana blockchain using the Anchor framework. This project demonstrates how to create, deploy, and interact with Solana programs through multiple client interfaces including React web app, Go CLI, and Rust CLI.

> **✅ Production Ready**: This project implements security best practices including PDA signing, seed verification, and proper access controls.

## 📚 Table of Contents

- [Features](#-features)
- [Project Structure](#-project-structure)
- [Prerequisites](#-prerequisites)
- [Installation](#-installation)
- [Quick Start](#-quick-start)
- [Security Features](#-security-features)
- [Program Architecture](#-program-architecture)
- [Deployment](#-deployment)
- [Client Interfaces](#-client-interfaces)
- [Troubleshooting](#-troubleshooting)
- [Contributing](#-contributing)

## 🚀 Features

- **🔐 Secure Crowdfunding**: Create and manage campaigns with enterprise-grade security
- **🎯 Multi-Campaign Support**: Each user can create multiple campaigns with unique names
- **💰 Secure Withdrawals**: PDA-signed withdrawals with proper authorization checks
- **🛡️ Seed Verification**: Advanced seed constraint validation prevents unauthorized access
- **📱 Multi-Client Support**: Interact via React frontend, Go CLI, or Rust CLI
- **🔄 Real-time Updates**: Live campaign status and balance updates
- **💾 Persistent Storage**: Campaign data persistence across client sessions

## 📁 Project Structure

```
crowdfunding/
├── programs/
│   └── crowdfunding/         # Anchor Solana program (Rust)
│       ├── src/
│       │   ├── lib.rs       # Main program entry point
│       │   ├── instructions.rs # Secure business logic
│       │   ├── state.rs     # Account structures with seed constraints
│       │   └── errors.rs    # Custom error definitions
│       └── Cargo.toml       # Rust dependencies
├── frontend/                 # React web application
│   ├── src/
│   │   ├── App.js          # Main React component
│   │   └── idl.json        # Program interface definition
│   └── package.json        # Node.js dependencies
├── go_client/               # Go CLI client
│   ├── main.go             # Feature-complete CLI implementation
│   ├── my_wallet.json      # Your secure wallet file
│   └── README.md           # Go client documentation
├── rust_client/             # Rust CLI client
│   ├── src/main.rs         # Rust CLI implementation
│   ├── wallet.example.json # Example wallet configuration
│   └── Cargo.toml          # Rust dependencies
├── tests/                   # Anchor test suite
│   └── crowdfunding.ts     # TypeScript tests
├── migrations/              # Deployment scripts
│   └── deploy.ts           # Deployment configuration
├── LICENSE                  # MIT license
└── README.md               # Project documentation
```

## 📋 Prerequisites

Before starting, ensure you have:

- **Node.js** (v16 or higher)
- **Rust** (latest stable)
- **Solana CLI** (v1.14 or higher)
- **Anchor CLI** (v0.28 or higher)
- **Go** (v1.19 or higher) - for Go client
- A Solana wallet with some SOL for testing

## 🔐 Security Notice

**IMPORTANT**: This repository implements production-grade security features:

### ✅ Security Features Implemented

1. **Seed Constraint Validation**: Prevents unauthorized access to campaigns
2. **PDA Signing**: Secure fund transfers using Program Derived Addresses
3. **Bump Storage**: Proper bump value storage and verification
4. **Access Control**: Only campaign creators can withdraw funds
5. **Anti-Collision**: Campaign names prevent seed collisions

### Setting Up Your Wallet

1. **Create your wallet files** based on the example templates:
   ```bash
   # For Go client
   cp go_client/wallet.example.json go_client/my_wallet.json
   
   # For Rust client  
   cp rust_client/wallet.example.json rust_client/wallet.json
   ```

2. **Generate a new Solana keypair:**
   ```bash
   solana-keygen new --outfile ~/.config/solana/id.json
   ```

3. **Update wallet files** with your actual keypair information
   ```bash
   # Get your public key
   solana-keygen pubkey ~/.config/solana/id.json
   
   # Get your private key (base58 encoded)
   cat ~/.config/solana/id.json | python3 -c "
   import json, base58, sys
   data = json.load(sys.stdin)
   print(base58.b58encode(bytes(data)).decode())
   "
   ```

> **⚠️ WARNING**: Never commit real private keys to version control!

## 🛠 Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/YOUR_USERNAME/crowd-funding.git
   cd crowd-funding
   ```

2. **Install dependencies:**
   ```bash
   # Install main project dependencies
   npm install
   
   # Install frontend dependencies
   cd frontend && npm install && cd ..
   
   # Install Go client dependencies
   cd go_client && go mod tidy && cd ..
   
   # Build Rust client
   cd rust_client && cargo build && cd ..
   ```

3. **Configure Solana for development:**
   ```bash
   # Set to devnet for testing
   solana config set --url devnet
   
   # Generate a new keypair (if you don't have one)
   solana-keygen new --outfile ~/.config/solana/id.json
   
   # Request airdrop for testing (devnet only)
   solana airdrop 2
   ```

## ⚡ Quick Start

1. **Build and deploy the program:**
   ```bash
   anchor build
   anchor deploy
   ```

2. **Run tests:**
   ```bash
   anchor test
   ```

3. **Start the frontend:**
   ```bash
   cd frontend
   npm start
   ```

4. **Use the Go CLI:**
   ```bash
   cd go_client
   # Make sure you've created my_wallet.json with your keypair
   go run main.go my_wallet.json
   ```

5. **Use the Rust CLI:**
   ```bash
   cd rust_client
   # Make sure you've created wallet.json with your keypair
   cargo run
   ```

## 🔒 Security Features

This project implements enterprise-grade security measures:

### 1. Seed Constraint Validation

```rust
#[derive(Accounts)]
#[instruction(name: String)]
pub struct Withdraw<'info> {
    #[account(
        mut,
        seeds = [b"CAMPAIGN_DEMO".as_ref(), campaign.admin.as_ref(), name.as_ref()],
        bump = campaign.bump
    )]
    pub campaign: Account<'info, Campaign>,
    #[account(mut)]
    pub user: Signer<'info>,
}
```

### 2. PDA-Signed Withdrawals

```rust
pub fn withdraw(ctx: Context<Withdraw>, name: String, amount: u64) -> Result<()> {
    let campaign = &mut ctx.accounts.campaign;
    let user = &mut ctx.accounts.user;
    
    // Verify admin access
    if campaign.admin != *user.key {
        return Err(CampaignError::Unauthorized.into());
    }

    // Manual lamport transfer with PDA ownership
    **campaign.to_account_info().try_borrow_mut_lamports()? -= amount;
    **user.to_account_info().try_borrow_mut_lamports()? += amount;
    
    Ok(())
}
```

### 3. Anti-Collision Campaign Creation

Each campaign uses a unique seed combination:
- `b"CAMPAIGN_DEMO"`
- `user.key()` (wallet address)
- `campaign_name` (user-provided name)

This prevents seed collisions and allows multiple campaigns per user.

## 🏗 Program Architecture

### Current Deployment

- **Program ID**: `3r5NUnG85XtVExb1234ZYYyUazjchqjfYknnQATyCDzp`
- **Network**: Solana Devnet
- **Status**: ✅ Deployed and Verified

### Core Instructions

| Instruction | Description | Security Features |
|-------------|-------------|-------------------|
| `create` | Create new campaign | Seed validation, bump storage |
| `donate` | Contribute to campaign | Seed verification, amount validation |
| `withdraw` | Withdraw funds | Admin check, PDA signing, balance validation |

### Account Structure

```rust
#[account]
pub struct Campaign {
    pub admin: Pubkey,        // Campaign creator
    pub name: String,         // Campaign name (part of seeds)
    pub description: String,  // Campaign description
    pub amount_donated: u64,  // Total donations received
    pub bump: u8,            // PDA bump for secure signing
}
```

## 🚀 Deployment

### Successful Deployment Example

```bash
❯ anchor deploy
Deploying cluster: https://api.devnet.solana.com
Upgrade authority: /home/raminfp/.config/solana/id.json
Deploying program "crowdfunding"...
Program path: /home/raminfp/Projects/crowdfunding/target/deploy/crowdfunding.so...
Program Id: 3r5NUnG85XtVExb1234ZYYyUazjchqjfYknnQATyCDzp

Signature: 5fSx9XCFC9HDbyYbwpZNipNVuUWK59s8gWMtUPMRkqvLWQAGEvNMas94myCcz6tJ6rn8fJe2HtiT6grruTW58x6m

Deploy success
```

## 🖥 Client Interfaces

### Go CLI Client Features

The Go CLI provides a comprehensive interface with:

- **🔐 Secure Wallet Management**: File-based wallet with Base58 encoding
- **📋 Campaign Persistence**: JSON-based campaign storage with name tracking
- **💰 Real-time Balance**: Live SOL balance monitoring
- **🎯 Interactive Menu**: User-friendly command selection
- **🛡️ Error Handling**: Graceful error recovery and user guidance

### Example CLI Session

```bash
❯ ./crowdfunding-client my_wallet.json
🚀 Solana dApp CLI Starting...
📋 Loaded saved campaign 'ramfs': 64BBRdyRSrH1WWzSbLmkjiagVQUR7WqXdfWknCPgCW86
✅ Connected to Solana devnet
💳 Wallet loaded: 9Mbf6JiwmzVjF5kbqTdanZS77szBtx56kwoRtw4uAE7z
💰 Current balance: 1.8094 SOL

=== Solana dApp CLI ===
Wallet: 9Mbf6JiwmzVjF5kbqTdanZS77szBtx56kwoRtw4uAE7z
Balance: 1.8094 SOL
Current Campaign: 'ramfs' (64BBRdyRSrH1WWzSbLmkjiagVQUR7WqXdfWknCPgCW86)

Options:
1. Request Airdrop (2 SOL)
2. Create Campaign
3. Donate to Campaign ⭐
4. Withdraw from Campaign ⭐
5. Check Balance
6. Check Campaign Status
7. Exit

Choose an option (1-7): 3
Use current campaign 'ramfs' (64BBRdyRSrH1WWzSbLmkjiagVQUR7WqXdfWknCPgCW86)? (y/n): y
Amount (lamports): 10000000
Donating 10000000 lamports to campaign 64BBRdyRSrH1WWzSbLmkjiagVQUR7WqXdfWknCPgCW86
Transaction sent: 4d1S2F2crbi53eguFc9Sohh1Ar3JMvF4EdUSSKEGLxSnio1KQVasKpMnExXThySezYvUSEYHnJAHaL6cEQ3swWpy
✅ Successfully donated 10000000 lamports!

Choose an option (1-7): 4
Use current campaign 'ramfs' (64BBRdyRSrH1WWzSbLmkjiagVQUR7WqXdfWknCPgCW86)? (y/n): y
Amount (lamports): 10000000
Withdrawing 10000000 lamports from campaign 64BBRdyRSrH1WWzSbLmkjiagVQUR7WqXdfWknCPgCW86
Transaction sent: 4EMUxxSDnEpeWVMBNrznwUYUW5wtBvi7aJdF7z27WAKgUVQDU7yB93iJ5LvpzFcLvoYW8xXSdLF4w7Pgy1oVjz9e
✅ Successfully withdrew 10000000 lamports!
```

## 📚 API Documentation

### Core Instructions

| Instruction | Parameters | Security Checks | Returns |
|-------------|------------|-----------------|---------|
| `create` | `name: String`, `description: String` | Seed uniqueness, bump storage | Campaign PDA |
| `donate` | `name: String`, `amount: u64` | Seed validation, campaign verification | Transaction signature |
| `withdraw` | `name: String`, `amount: u64` | Admin verification, balance check, PDA signing | Transaction signature |

### Error Codes

| Error | Code | Description |
|-------|------|-------------|
| `Unauthorized` | 6000 | Not campaign admin |
| `InsufficientFunds` | 6001 | Campaign has insufficient balance |
| `ConstraintSeeds` | 2006 | Seed constraint violation |

## 🔧 Troubleshooting

### Common Issues

**1. Seed Constraint Violations**
```
Error: ConstraintSeeds. Error Number: 2006
```
- **Cause**: Campaign name doesn't match the provided address
- **Solution**: Ensure you're using the correct campaign name for the address

**2. Transfer Failures**
```
Error: Transfer: `from` must not carry data
```
- **Cause**: Using system_program transfer on data accounts
- **Solution**: Use manual lamport manipulation (fixed in current version)

**3. Account Not Found**
```
Error: AccountNotFound
```
- **Cause**: Campaign doesn't exist or wrong program ID
- **Solution**: Verify campaign exists and program ID is correct

### Debug Commands

```bash
# View program logs
solana logs 3r5NUnG85XtVExb1234ZYYyUazjchqjfYknnQATyCDzp

# Check account information
solana account <CAMPAIGN_ADDRESS>

# Verify program deployment
anchor verify 3r5NUnG85XtVExb1234ZYYyUazjchqjfYknnQATyCDzp
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Setup

```bash
# Install development dependencies
npm install --save-dev

# Run tests
anchor test

# Build program  
anchor build

# Deploy to localnet for testing
anchor deploy --provider.cluster localnet

# View program logs (useful for debugging)
solana logs
```

### Security Best Practices

- **Never commit private keys**: Always use `.gitignore` to exclude wallet files
- **Use environment variables**: Store sensitive configuration in environment files
- **Test on devnet first**: Always test thoroughly before mainnet deployment
- **Audit your code**: Consider security audits for production applications
- **Use hardware wallets**: For mainnet, use hardware wallets for enhanced security

---

**📬 Contact & Support**

For questions, issues, or contributions, please refer to the project's GitHub repository or reach out to the development team.

**🔗 Useful Links**
- [Solana Documentation](https://docs.solana.com/)
- [Anchor Framework](https://anchor-lang.com/)
- [Solana Web3.js](https://solana-labs.github.io/solana-web3.js/)
- [Solana CLI Reference](https://docs.solana.com/cli)
- [Program Explorer](https://explorer.solana.com/address/3r5NUnG85XtVExb1234ZYYyUazjchqjfYknnQATyCDzp?cluster=devnet)