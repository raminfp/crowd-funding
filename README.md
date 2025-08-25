# Solana Crowdfunding dApp

[![Solana](https://img.shields.io/badge/Solana-Blockchain-blue)](https://solana.com/)
[![Anchor](https://img.shields.io/badge/Anchor-Framework-purple)](https://anchor-lang.com/)
[![React](https://img.shields.io/badge/React-Frontend-blue)](https://reactjs.org/)
[![Go](https://img.shields.io/badge/Go-Client-cyan)](https://golang.org/)
[![Rust](https://img.shields.io/badge/Rust-Client-orange)](https://rust-lang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A comprehensive decentralized crowdfunding application built on the Solana blockchain using the Anchor framework. This project demonstrates how to create, deploy, and interact with Solana programs through multiple client interfaces including React web app, Go CLI, and Rust CLI.

> **⚠️ Disclaimer**: This is a educational/demonstration project. Use in production environments at your own risk and ensure proper security audits.

## 📚 Table of Contents

- [Features](#-features)
- [Project Structure](#-project-structure)
- [Prerequisites](#-prerequisites)
- [Installation](#-installation)
- [Quick Start](#-quick-start)
- [Anchor Program Architecture](#-anchor-program-architecture)
- [Code Examples](#-code-examples)
- [Deployment](#-deployment)
- [Client Interfaces](#-client-interfaces)
- [API Documentation](#-api-documentation)
- [Troubleshooting](#-troubleshooting)
- [Contributing](#-contributing)

## 🚀 Features

- **Decentralized Crowdfunding**: Create and manage crowdfunding campaigns on Solana
- **Multi-Client Support**: Interact via React frontend, Go CLI, or JavaScript
- **Secure Transactions**: Built with Anchor framework for enhanced security
- **Real-time Updates**: Live campaign status and balance updates
- **User Account Management**: Complete user account system with transfers
- **Campaign Management**: Create, donate, withdraw, and monitor campaigns

## 📁 Project Structure

```
crowdfunding/
├── programs/
│   └── crowdfunding/         # Anchor Solana program (Rust)
│       ├── src/
│       │   ├── lib.rs       # Main program entry point
│       │   ├── instructions.rs # Business logic implementations
│       │   ├── state.rs     # Account structures and data models
│       │   └── errors.rs    # Custom error definitions
│       └── Cargo.toml       # Rust dependencies
├── frontend/                 # React web application
│   ├── src/
│   │   ├── App.js          # Main React component
│   │   └── idl.json        # Program interface definition
│   └── package.json        # Node.js dependencies
├── go_client/               # Go CLI client
│   ├── main.go             # CLI implementation
│   ├── wallet.example.json # Example wallet configuration
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

Before starting this tutorial, ensure you have:

- **Node.js** (v16 or higher)
- **Rust** (latest stable)
- **Solana CLI** (v1.14 or higher)
- **Anchor CLI** (v0.28 or higher)
- **Go** (v1.19 or higher) - for Go client
- A Solana wallet with some SOL for testing

## 🔐 Security Notice

**IMPORTANT**: This repository does not contain private keys or wallet files for security reasons. You must create your own wallet files for testing.

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
   solana-keygen pubkey ~/.config/solana/id.json --keypair
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

## 🏗 Anchor Program Architecture

### Basic Structure of an Anchor Program
An Anchor program consists of several key components:

- **`declare_id!` macro**: Declares the program's on-chain ID
- **`#[program]` module**: Contains the core business logic functions (instructions)
- **`#[derive(Accounts)]` structs**: Define the accounts required for each instruction
- **`#[account]` structs**: Define the custom data structures stored in program accounts
- **Error handling**: Custom error types for robust program execution

## 🚀 Deployment

### Successful Deployment Example

Here's what a successful deployment looks like:

```bash
❯ anchor deploy
Deploying cluster: https://api.devnet.solana.com
Upgrade authority: /home/raminfp/.config/solana/id.json
Deploying program "crowdfunding"...
Program path: /home/raminfp/Projects/crowdfunding/target/deploy/crowdfunding.so...
Program Id: 3r5NUnG85XtVExb1234ZYYyUazjchqjfYknnQATyCDzp

Signature: JByaWbxEr74GRQiMEZhTUBYgQAg6uLCTynURQWr6mTD6y7kLbCKZWA5vP3SBGAHqWVVd4DZnBLBa28WB2N62xGA

Deploy success
```

**Key Points:**
- Program ID uniquely identifies your program on the blockchain
- The signature confirms the successful transaction
- Deployment creates a `.so` file (Solana program binary)

## 📝 Code Examples

### 1. User Account Management Program

This example demonstrates a complete user account system with balance management and secure transfers:

```rust
use anchor_lang::prelude::*;

declare_id!("Fg6PaFpoGXkYsidMpWTK6W2BeZ7FEfcYkg476zPFsLnS");

#[program]
pub mod user_account_program {
    use super::*;

    pub fn initialize_account(ctx: Context<InitializeAccount>, username: String) -> Result<()> {
        let user_account = &mut ctx.accounts.user_account;
        user_account.owner = *ctx.accounts.user.key;
        user_account.balance = 0;
        user_account.is_active = true;
        user_account.creation_date = Clock::get()?.unix_timestamp;
        user_account.username = username;
        Ok(())
    }

    pub fn deposit(ctx: Context<ModifyBalance>, amount: u64) -> Result<()> {
        let user_account = &mut ctx.accounts.user_account;
        // Only the owner can modify the balance
        require!(user_account.owner == *ctx.accounts.user.key, CustomError::Unauthorized);
        user_account.balance = user_account.balance.checked_add(amount).ok_or(CustomError::Overflow)?;
        Ok(())
    }

    pub fn transfer(ctx: Context<Transfer>, amount: u64) -> Result<()> {
        let from = &mut ctx.accounts.from_account;
        let to = &mut ctx.accounts.to_account;
        // Only the source account owner can perform transfer
        require!(from.owner == *ctx.accounts.from_signer.key, CustomError::Unauthorized);
        // Check sufficient balance
        require!(from.balance >= amount, CustomError::InsufficientFunds);
        from.balance -= amount;
        to.balance = to.balance.checked_add(amount).ok_or(CustomError::Overflow)?;
        Ok(())
    }
}

#[derive(Accounts)]
pub struct InitializeAccount<'info> {
    #[account(init, payer = user, space = 8 + 32 + 8 + 1 + 8 + 4 + 32)] // Storage space estimation
    pub user_account: Account<'info, UserAccount>,
    #[account(mut)]
    pub user: Signer<'info>,
    pub system_program: Program<'info, System>,
}

#[derive(Accounts)]
pub struct ModifyBalance<'info> {
    #[account(mut)]
    pub user_account: Account<'info, UserAccount>,
    pub user: Signer<'info>,
}

#[derive(Accounts)]
pub struct Transfer<'info> {
    #[account(mut)]
    pub from_account: Account<'info, UserAccount>,
    #[account(mut)]
    pub to_account: Account<'info, UserAccount>,
    pub from_signer: Signer<'info>,
}

#[account]
pub struct UserAccount {
    pub owner: Pubkey,
    pub balance: u64,
    pub is_active: bool,
    pub creation_date: i64,
    pub username: String,
}

#[error_code]
pub enum CustomError {
    #[msg("Unauthorized Access")]
    Unauthorized,
    #[msg("Insufficient Funds")]
    InsufficientFunds,
    #[msg("Balance Overflow")]
    Overflow,
}
```

#### 📖 Method Explanations:

**`initialize_account`**: Creates a new user account with username and sets initial values like zero balance and creation date.

**`deposit`**: Allows account owner to add amount to their balance.

**`transfer`**: Enables source account owner to transfer amount to another account. First ensures the owner has permission and sufficient balance.

#### 🔐 Security Features:

- **Ownership verification**: All operations verify ownership through public key comparison with signer
- **Space allocation**: Account space is estimated and should be adjusted based on string sizes and data
- **Custom errors**: Defined to allow the program to return appropriate errors according to business logic
- **Overflow protection**: Uses `checked_add` and `checked_sub` to prevent integer overflow attacks


### 2. Core Crowdfunding Program Methods

Here are the essential methods for the crowdfunding functionality:

```rust
pub fn initialize(ctx: Context<Initialize>, ...) -> Result<()> { 
    // Initialize crowdfunding campaign with target goal and deadline
    ...
}

pub fn deposit(ctx: Context<Deposit>, amount: u64) -> Result<()> { 
    // Allow users to contribute to campaigns
    ...
}

pub fn withdraw(ctx: Context<Withdraw>, amount: u64) -> Result<()> {
    // Secure withdrawal with balance validation
    require!(ctx.accounts.user_account.balance >= amount, ErrorCode::InsufficientFunds);
    // Reentrancy protection and updates with checked_sub and checked_add
    ...
}

pub fn transfer(ctx: Context<Transfer>, amount: u64) -> Result<()> {
    // Owner and balance verification with secure transfer
}

pub fn pause(ctx: Context<Pause>) -> Result<()> {
    // Only administrators can pause the contract
}
```

#### 🔒 Security Considerations:
- **Reentrancy protection**: All state changes occur before external calls
- **Balance validation**: Always check sufficient funds before transfers
- **Access control**: Verify user permissions for each operation
- **Integer overflow protection**: Use `checked_add` and `checked_sub`

## 🖥 Client Interfaces

### Go CLI Client

The Go CLI provides a user-friendly command-line interface for interacting with the crowdfunding program:

```bash
❯ go run main.go my_wallet.json
🚀 Solana dApp CLI Starting...
📋 Loaded saved campaign: HeHiRzgqE18tasfssX1BAruFrzLP5rBzA2BqvQB4sVAe
✅ Connected to Solana devnet
💳 Wallet loaded: 7gkrxUoUVa1aQoYJwv8RYHWjjcb7Vc8KZ5Soy3CyV822
💰 Current balance: 1.9354 SOL

=== Solana dApp CLI ===
Wallet: 7gkrxUoUVa1aQoYJwv8RYHWjjcb7Vc8KZ5Soy3CyV822
Balance: 1.9354 SOL
Current Campaign: HeHiRzgqE18tasfssX1BAruFrzLP5rBzA2BqvQB4sVAe

Options:
1. Request Airdrop (2 SOL)
2. Create Campaign
3. Donate to Campaign ⭐
4. Withdraw from Campaign ⭐
5. Check Balance
6. Check Campaign Status
7. Exit

Choose an option (1-7): 2
Campaign name: ramin
Campaign description: ramin
✅ Found properly initialized campaign at HeHiRzgqE18tasfssX1BAruFrzLP5rBzA2BqvQB4sVAe
✅ Campaign already exists at: HeHiRzgqE18tasfssX1BAruFrzLP5rBzA2BqvQB4sVAe
📋 Using existing campaign for future operations!

=== Solana dApp CLI ===
Wallet: 7gkrxUoUVa1aQoYJwv8RYHWjjcb7Vc8KZ5Soy3CyV822
Balance: 1.9354 SOL
Current Campaign: HeHiRzgqE18tasfssX1BAruFrzLP5rBzA2BqvQB4sVAe

Options:
1. Request Airdrop (2 SOL)
2. Create Campaign
3. Donate to Campaign ⭐
4. Withdraw from Campaign ⭐
5. Check Balance
6. Check Campaign Status
7. Exit

Choose an option (1-7): 3   
Use current campaign (HeHiRzgqE18tasfssX1BAruFrzLP5rBzA2BqvQB4sVAe)? (y/n): y
Amount (lamports): 100000000
Donating 100000000 lamports to campaign HeHiRzgqE18tasfssX1BAruFrzLP5rBzA2BqvQB4sVAe
Transaction sent: 2E2wia7KFqb5BzNBvPuY3tCVvS82iGdND7KXdk9Ucb9K2V5QqLx9sqgnHiNNQFnGHan48tN2aodYvbbAd5Vy7AXF
✅ Successfully donated 100000000 lamports!

Choose an option (1-7): 6

🔍 Campaign Status for Wallet: 7gkrxUoUVa1aQoYJwv8RYHWjjcb7Vc8KZ5Soy3CyV822
📍 Expected Campaign Address: HeHiRzgqE18tasfssX1BAruFrzLP5rBzA2BqvQB4sVAe
🔗 Explorer Link: https://explorer.solana.com/address/HeHiRzgqE18tasfssX1BAruFrzLP5rBzA2BqvQB4sVAe?cluster=devnet
📊 Account Info:
   Owner: 
   Data Size: 9000 bytes
   Lamports: 164530880
✅ Account is properly owned by the crowdfunding program
✅ Account appears to have campaign data

Choose an option (1-7): 5
Current balance: 1.8354 SOL



https://explorer.solana.com/tx/2E2wia7KFqb5BzNBvPuY3tCVvS82iGdND7KXdk9Ucb9K2V5QqLx9sqgnHiNNQFnGHan48tN2aodYvbbAd5Vy7AXF?cluster=devnet

```

#### 🎯 CLI Features:
- **Interactive menu**: Easy-to-use command selection
- **Campaign management**: Create, monitor, and interact with campaigns
- **Wallet integration**: Secure wallet loading and balance checking
- **Real-time feedback**: Live status updates and transaction confirmations
- **Error handling**: Graceful error messages and recovery

## 📚 API Documentation

### Core Instructions

| Instruction | Description | Parameters |
|------------|-------------|------------|
| `initialize` | Create a new crowdfunding campaign | `name`, `description`, `target_amount`, `deadline` |
| `donate` | Contribute to a campaign | `campaign_id`, `amount` |
| `withdraw` | Withdraw funds from campaign | `campaign_id`, `amount` |
| `close_campaign` | Close an expired or successful campaign | `campaign_id` |
| `get_campaign_status` | Retrieve campaign information | `campaign_id` |

### Account Types

- **Campaign**: Stores campaign metadata, goal, current amount, and deadline
- **UserAccount**: Manages user balance and transaction history
- **DonationRecord**: Tracks individual contributions for transparency

## 🔧 Troubleshooting

### Common Issues

**1. Deployment Failures**
```bash
# Check Solana CLI configuration
solana config get

# Verify wallet has sufficient SOL
solana balance

# Check network connectivity
solana cluster-version
```

**2. Transaction Failures**
- Ensure wallet has enough SOL for transaction fees
- Verify program ID matches deployed program
- Check account ownership and permissions

**3. Client Connection Issues**
- Confirm RPC endpoint is accessible
- Verify wallet file exists and is properly formatted
- Check network configuration (devnet/mainnet)

### Debug Commands

```bash
# View program logs
solana logs <PROGRAM_ID>

# Check account information
solana account <ACCOUNT_ADDRESS>

# Verify program deployment
anchor verify <PROGRAM_ID>
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