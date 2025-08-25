# Solana Crowdfunding Go Client

A command-line interface for interacting with the Solana crowdfunding program, converted from the React frontend.

## Features

- 🔐 **Wallet Management**: Automatic wallet creation and persistent storage
- 💰 **Balance Checking**: Real-time SOL balance display
- 🚁 **Airdrop**: Request SOL from devnet faucet
- 📋 **Campaign Management**: Create and manage crowdfunding campaigns
- 💸 **Donations**: Donate SOL to campaigns
- 💵 **Withdrawals**: Withdraw funds (campaign admin only)
- 💾 **Persistence**: Saves campaign addresses for easy reuse

## Setup

1. **Create Wallet File**:
   ```bash
   cd go_client
   # Copy the example wallet template
   cp wallet.example.json my_wallet.json
   
   # Edit my_wallet.json with your actual keypair information
   # You can generate a keypair using: solana-keygen new
   ```

2. **Install Dependencies**:
   ```bash
   go mod tidy
   ```

3. **Run the Application**:
   ```bash
   go run main.go my_wallet.json
   ```

   > **Note**: The wallet file is required and must contain your Solana keypair.

## Usage

### First Time Setup

1. **Prepare Wallet**: Create `my_wallet.json` with your keypair (see Setup above)
2. **Start Application**: Run with `go run main.go my_wallet.json`
3. **Request Airdrop**: Use option 1 to get SOL for transaction fees (devnet only)
4. **Create Campaign**: Use option 2 to create your first campaign
5. **Interact**: Donate to or withdraw from campaigns

### Menu Options

1. **Request Airdrop (2 SOL)**: Get SOL from devnet faucet
2. **Create Campaign**: Create a new crowdfunding campaign
3. **Donate to Campaign**: Send SOL to a campaign
4. **Withdraw from Campaign**: Withdraw funds (admin only)
5. **Check Balance**: Display current SOL balance
6. **Exit**: Close the application

### Smart Features

- **Campaign Persistence**: Created campaigns are automatically saved and suggested for future operations
- **Wallet Persistence**: Your wallet is saved to `wallet.json` for reuse
- **Auto-Loading**: Previously used campaign addresses are loaded on startup
- **Error Handling**: User-friendly error messages for common issues

## Troubleshooting

### "Airdrop failed" or Rate Limit Errors
- **Cause**: Solana devnet faucet has rate limits
- **Solution**: Wait a few minutes and try again, or use an alternative faucet

### "Account Not Found" on Campaign Creation
- **Cause**: Insufficient SOL for transaction fees
- **Solution**: Request an airdrop first (option 1)

### "Unauthorized" on Withdrawal
- **Cause**: Only campaign creators can withdraw funds
- **Solution**: Ensure you're using the same wallet that created the campaign

## Files Created

- `my_wallet.json`: Your wallet's private key (keep secure!)
- `campaign.txt`: Last used campaign address
- `main`: Compiled binary (if you use `go build`)

## Program Details

- **Program ID**: `7XJkGrdSHn3chc7rsv1xdZEKtwP9w5rSx1sHohzM5skv`
- **Network**: Solana Devnet
- **RPC Endpoint**: Solana devnet RPC

## Security Notes

⚠️ **CRITICAL SECURITY WARNINGS**:
- **Never commit wallet files**: Keep `my_wallet.json` secure and never commit to git
- **Devnet only**: This setup is for development/testing on devnet only
- **Private key safety**: Your private key grants full access to your wallet
- **Backup important**: Back up your wallet file securely
- **Production use**: For mainnet, use hardware wallets and proper key management

## Conversion from React

This Go client provides the same functionality as the React frontend (`frontend/src/App.js`) but as a CLI tool:

| React Feature | Go Equivalent |
|---------------|---------------|
| Phantom Wallet | File-based wallet |
| React State | Persistent file storage |
| useEffect hooks | Auto-loading on startup |
| Error alerts | CLI error messages |
| UI buttons | Menu options |
| Real-time updates | Balance refresh |

---

**Happy crowdfunding!** 🚀
