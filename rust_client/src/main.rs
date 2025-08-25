use anyhow::{anyhow, Result};
use base58::FromBase58;
use serde::{Deserialize, Serialize};
use sha2::{Digest, Sha256};
use solana_client::rpc_client::RpcClient;
use solana_program::pubkey::Pubkey;
use solana_sdk::{
    commitment_config::CommitmentConfig,
    instruction::{AccountMeta, Instruction},
    signature::{Keypair, Signature, Signer},
    system_program,
    transaction::Transaction,
};
use std::fs;
use std::io::{self, Write};
use std::str::FromStr;

const PROGRAM_ID: &str = "7GMjTXTH1KS1Q46ngEnUYakAJi4xb2KJ3JsbJW2UNpHC";
const NETWORK: &str = "https://api.devnet.solana.com";
const LAMPORTS_PER_SOL: u64 = 1_000_000_000;

/// Campaign account structure matching the Solana program
#[derive(Debug, Serialize, Deserialize)]
pub struct Campaign {
    pub admin: Pubkey,
    pub name: String,
    pub description: String,
    pub amount_donated: u64,
}

/// Wallet data format for JSON serialization (Go client format)
#[derive(Serialize, Deserialize)]
struct WalletData {
    #[serde(rename = "publicKey")]
    public_key: String,
    #[serde(rename = "privateKey")]  
    private_key: String,
}

/// Main Solana dApp client
pub struct SolanaDApp {
    client: RpcClient,
    wallet: Keypair,
    program_id: Pubkey,
    campaign_address: Option<Pubkey>,
}

impl SolanaDApp {
    /// Create a new SolanaDApp instance
    pub fn new(key_path: Option<&str>) -> Result<Self> {
        let client = RpcClient::new_with_commitment(NETWORK, CommitmentConfig::finalized());
        let wallet = Self::load_or_create_wallet(key_path)?;
        let program_id = Pubkey::from_str(PROGRAM_ID)?;
        let mut app = Self {
            client,
            wallet,
            program_id,
            campaign_address: None,
        };

        // Try to load saved campaign
        app.load_saved_campaign();
        
        // Check for existing campaign if none saved
        if app.campaign_address.is_none() {
            // Note: check_existing_campaign is async, so we skip auto-detection in new()
            // Users can manually call check_existing_campaign() later if needed
            println!("💡 No saved campaign found. Use 'check' command to look for existing campaigns.");
        }

        Ok(app)
    }

    /// Load existing wallet or create new one
    fn load_or_create_wallet(key_path: Option<&str>) -> Result<Keypair> {
        if let Some(path) = key_path {
            // Try to load existing wallet
            let data = fs::read_to_string(path)?;
            
            // Try wallet data format first (Go client base58 keys format)
            if let Ok(wallet_data) = serde_json::from_str::<WalletData>(&data) {
                let private_key_vec = wallet_data.private_key.from_base58()
                    .map_err(|e| anyhow!("Failed to decode base58: {:?}", e))?;
                let private_key_bytes: [u8; 64] = private_key_vec.try_into()
                    .map_err(|_| anyhow!("Invalid private key length, expected 64 bytes"))?;
                return Ok(Keypair::from_bytes(&private_key_bytes)?);
            }

            // Try byte array format (legacy)
            if let Ok(key_array) = serde_json::from_str::<Vec<u8>>(&data) {
                return Ok(Keypair::from_bytes(&key_array)?);
            }

            return Err(anyhow!("Failed to parse key file"));
        }

        // Generate new wallet
        let keypair = Keypair::new();
        let key_bytes = keypair.to_bytes().to_vec();
        if let Ok(key_json) = serde_json::to_string(&key_bytes) {
            if fs::write("wallet.json", &key_json).is_ok() {
                println!("New wallet saved to wallet.json");
            }
        }
        
        Ok(keypair)
    }

    /// Load saved campaign address from file
    fn load_saved_campaign(&mut self) {
        if let Ok(data) = fs::read_to_string("campaign.txt") {
            let campaign_str = data.trim();
            if !campaign_str.is_empty() {
                if let Ok(pubkey) = Pubkey::from_str(campaign_str) {
                    self.campaign_address = Some(pubkey);
                    println!("📋 Loaded saved campaign: {}", campaign_str);
                }
            }
        }
    }

    /// Save current campaign address to file
    fn save_campaign(&self) {
        if let Some(campaign) = self.campaign_address {
            if let Err(e) = fs::write("campaign.txt", campaign.to_string()) {
                eprintln!("Warning: failed to save campaign address: {}", e);
            }
        }
    }

    /// Generate discriminator for Anchor instructions
    fn generate_discriminator(namespace: &str, name: &str) -> [u8; 8] {
        let preimage = format!("{}:{}", namespace, name);
        let mut hasher = Sha256::new();
        hasher.update(preimage.as_bytes());
        let hash = hasher.finalize();
        let mut discriminator = [0u8; 8];
        discriminator.copy_from_slice(&hash[0..8]);
        discriminator
    }

    /// Get wallet SOL balance
    pub async fn get_balance(&self) -> Result<f64> {
        let balance = self.client.get_balance(&self.wallet.pubkey())?;
        Ok(balance as f64 / LAMPORTS_PER_SOL as f64)
    }

    /// Request SOL airdrop from devnet faucet
    pub async fn request_airdrop(&self) -> Result<()> {
        println!("Requesting airdrop...");
        
        let signature = self.client.request_airdrop(
            &self.wallet.pubkey(),
            2 * LAMPORTS_PER_SOL,
        )?;

        println!("Airdrop requested. Transaction signature: {}", signature);
        println!("Waiting for confirmation...");

        // Wait for confirmation
        self.client.confirm_transaction(&signature)?;
        println!("✅ Airdrop confirmed!");
        
        Ok(())
    }

    /// Generate Program Derived Address for campaign
    fn create_campaign_pda(&self) -> Result<(Pubkey, u8)> {
        let wallet_pubkey = self.wallet.pubkey();
        let seeds = &[
            b"CAMPAIGN_DEMO",
            wallet_pubkey.as_ref(),
        ];
        
        let (pda, bump) = Pubkey::find_program_address(seeds, &self.program_id);
        Ok((pda, bump))
    }

    /// Check if campaign already exists for this wallet
    pub async fn check_existing_campaign(&self) -> Result<Option<Pubkey>> {
        let (campaign_pda, _) = self.create_campaign_pda()?;
        
        match self.client.get_account(&campaign_pda) {
            Ok(account) => {
                // Check if owned by our program and has sufficient data
                if account.owner == self.program_id && account.data.len() >= 32 {
                    println!("✅ Found properly initialized campaign at {}", campaign_pda);
                    return Ok(Some(campaign_pda));
                } else if account.owner == system_program::id() {
                    println!("⚠️  Found uninitialized account at {}", campaign_pda);
                }
            }
            Err(_) => {
                // Account doesn't exist, which is fine
            }
        }
        
        Ok(None)
    }

    /// Check detailed campaign status
    pub async fn check_campaign_status(&self) -> Result<()> {
        let (campaign_pda, _) = self.create_campaign_pda()?;
        
        println!("\n🔍 Campaign Status for Wallet: {}", self.wallet.pubkey());
        println!("📍 Expected Campaign Address: {}", campaign_pda);
        println!("🔗 Explorer Link: https://explorer.solana.com/address/{}?cluster=devnet", campaign_pda);

        match self.client.get_account(&campaign_pda) {
            Ok(account) => {
                println!("📊 Account Info:");
                println!("   Owner: {}", account.owner);
                println!("   Data Size: {} bytes", account.data.len());
                println!("   Lamports: {}", account.lamports);

                if account.owner == system_program::id() {
                    println!("⚠️  Account is allocated but NOT initialized by the crowdfunding program");
                    println!("💡 This means a previous campaign creation failed partway through");
                    println!("🔧 The account exists but has no campaign data");
                    println!("❗ You'll need to use a different wallet or wait for the account to be reclaimed");
                } else if account.owner == self.program_id {
                    println!("✅ Account is properly owned by the crowdfunding program");
                    if account.data.len() >= 32 {
                        println!("✅ Account appears to have campaign data");
                    } else {
                        println!("⚠️  Account is owned by program but has insufficient data");
                    }
                } else {
                    println!("❓ Account is owned by unknown program: {}", account.owner);
                }
            }
            Err(_) => {
                println!("❌ Account does not exist");
                println!("✅ You can create a new campaign!");
            }
        }

        Ok(())
    }

    /// Create a new campaign
    pub async fn create_campaign(&mut self, name: &str, description: &str) -> Result<()> {
        // Check for existing campaign first
        if let Ok(Some(existing)) = self.check_existing_campaign().await {
            println!("✅ Campaign already exists at: {}", existing);
            self.campaign_address = Some(existing);
            self.save_campaign();
            println!("📋 Using existing campaign for future operations!");
            return Ok(());
        }

        println!("Creating campaign: {}", name);

        let (campaign_pda, _) = self.create_campaign_pda()?;

        // Build instruction data
        let mut instruction_data = Self::generate_discriminator("global", "create").to_vec();
        
        // Serialize name (u32 length + bytes)
        let name_bytes = name.as_bytes();
        instruction_data.extend_from_slice(&(name_bytes.len() as u32).to_le_bytes());
        instruction_data.extend_from_slice(name_bytes);
        
        // Serialize description (u32 length + bytes)
        let desc_bytes = description.as_bytes();
        instruction_data.extend_from_slice(&(desc_bytes.len() as u32).to_le_bytes());
        instruction_data.extend_from_slice(desc_bytes);

        let instruction = Instruction::new_with_bytes(
            self.program_id,
            &instruction_data,
            vec![
                AccountMeta::new(campaign_pda, false),
                AccountMeta::new(self.wallet.pubkey(), true),
                AccountMeta::new_readonly(system_program::id(), false),
            ],
        );

        let signature = self.send_transaction(&[instruction]).await?;
        
        println!("Campaign created! Transaction: {}", signature);
        println!("Campaign address: {}", campaign_pda);
        
        self.campaign_address = Some(campaign_pda);
        self.save_campaign();
        println!("✅ Campaign address saved for quick access!");

        Ok(())
    }

    /// Donate to a campaign
    pub async fn donate_to_campaign(&self, campaign_address: &str, amount: u64) -> Result<()> {
        println!("Donating {} lamports to campaign {}", amount, campaign_address);

        let campaign_pubkey = Pubkey::from_str(campaign_address)?;
        
        let mut instruction_data = Self::generate_discriminator("global", "donate").to_vec();
        instruction_data.extend_from_slice(&amount.to_le_bytes());

        let instruction = Instruction::new_with_bytes(
            self.program_id,
            &instruction_data,
            vec![
                AccountMeta::new(campaign_pubkey, false),
                AccountMeta::new(self.wallet.pubkey(), true),
                AccountMeta::new_readonly(system_program::id(), false),
            ],
        );

        let signature = self.send_transaction(&[instruction]).await?;
        println!("Transaction sent: {}", signature);
        
        Ok(())
    }

    /// Withdraw from a campaign (admin only)
    pub async fn withdraw_from_campaign(&self, campaign_address: &str, amount: u64) -> Result<()> {
        println!("Withdrawing {} lamports from campaign {}", amount, campaign_address);

        let campaign_pubkey = Pubkey::from_str(campaign_address)?;
        
        let mut instruction_data = Self::generate_discriminator("global", "withdraw").to_vec();
        instruction_data.extend_from_slice(&amount.to_le_bytes());

        let instruction = Instruction::new_with_bytes(
            self.program_id,
            &instruction_data,
            vec![
                AccountMeta::new(campaign_pubkey, false),
                AccountMeta::new(self.wallet.pubkey(), true),
            ],
        );

        let signature = self.send_transaction(&[instruction]).await?;
        println!("Transaction sent: {}", signature);
        
        Ok(())
    }

    /// Helper method to send transactions
    async fn send_transaction(&self, instructions: &[Instruction]) -> Result<Signature> {
        let recent_blockhash = self.client.get_latest_blockhash()?;
        
        let transaction = Transaction::new_signed_with_payer(
            instructions,
            Some(&self.wallet.pubkey()),
            &[&self.wallet],
            recent_blockhash,
        );

        let signature = self.client.send_and_confirm_transaction(&transaction)?;
        Ok(signature)
    }

    /// Show interactive menu
    fn show_menu(&self) {
        println!("\n=== Solana dApp CLI ===");
        println!("Wallet: {}", self.wallet.pubkey());

        // Show current campaign
        if let Some(campaign) = self.campaign_address {
            println!("Current Campaign: {}", campaign);
        } else {
            println!("Current Campaign: None");
        }

        println!("\nOptions:");
        println!("1. Request Airdrop (2 SOL)");
        println!("2. Create Campaign");
        if self.campaign_address.is_some() {
            println!("3. Donate to Campaign ⭐");
            println!("4. Withdraw from Campaign ⭐");
        } else {
            println!("3. Donate to Campaign");
            println!("4. Withdraw from Campaign");
        }
        println!("5. Check Balance");
        println!("6. Check Campaign Status");
        println!("7. Exit");
        print!("\nChoose an option (1-7): ");
        io::stdout().flush().unwrap();
    }

    /// Run interactive CLI
    pub async fn run(&mut self) -> Result<()> {
        // Show initial balance
        if let Ok(balance) = self.get_balance().await {
            println!("💰 Current balance: {:.4} SOL", balance);
            if balance < 0.01 {
                println!("⚠️  Low balance! You may want to request an airdrop.");
            }
        }

        loop {
            self.show_menu();

            let mut input = String::new();
            io::stdin().read_line(&mut input)?;
            let choice = input.trim();

            match choice {
                "1" => {
                    if let Err(e) = self.request_airdrop().await {
                        if e.to_string().contains("airdrop") {
                            println!("❌ Airdrop failed. You may have reached the rate limit. Try again later.");
                        } else {
                            println!("❌ Error requesting airdrop: {}", e);
                        }
                    }
                }
                "2" => {
                    print!("Campaign name: ");
                    io::stdout().flush()?;
                    let mut name = String::new();
                    io::stdin().read_line(&mut name)?;
                    let name = name.trim();

                    print!("Campaign description: ");
                    io::stdout().flush()?;
                    let mut description = String::new();
                    io::stdin().read_line(&mut description)?;
                    let description = description.trim();

                    if let Err(e) = self.create_campaign(name, description).await {
                        if e.to_string().contains("insufficient") {
                            println!("❌ Insufficient SOL in your wallet. Please use option 1 to get SOL via airdrop.");
                        } else {
                            println!("❌ Error creating campaign: {}", e);
                        }
                    }
                }
                "3" => {
                    let address = if let Some(campaign) = self.campaign_address {
                        print!("Use current campaign ({})? (y/n): ", campaign);
                        io::stdout().flush()?;
                        let mut response = String::new();
                        io::stdin().read_line(&mut response)?;
                        
                        if response.trim().to_lowercase() == "y" {
                            campaign.to_string()
                        } else {
                            print!("Campaign address: ");
                            io::stdout().flush()?;
                            let mut addr = String::new();
                            io::stdin().read_line(&mut addr)?;
                            addr.trim().to_string()
                        }
                    } else {
                        print!("Campaign address: ");
                        io::stdout().flush()?;
                        let mut addr = String::new();
                        io::stdin().read_line(&mut addr)?;
                        addr.trim().to_string()
                    };

                    print!("Amount (lamports): ");
                    io::stdout().flush()?;
                    let mut amount_str = String::new();
                    io::stdin().read_line(&mut amount_str)?;
                    
                    match amount_str.trim().parse::<u64>() {
                        Ok(amount) if amount > 0 => {
                            if let Err(e) = self.donate_to_campaign(&address, amount).await {
                                if e.to_string().contains("insufficient") {
                                    println!("❌ Insufficient SOL for donation. Please check your balance or request an airdrop.");
                                } else {
                                    println!("❌ Error donating: {}", e);
                                }
                            } else {
                                println!("✅ Successfully donated {} lamports!", amount);
                            }
                        }
                        _ => println!("❌ Invalid amount. Please enter a valid number greater than 0."),
                    }
                }
                "4" => {
                    let address = if let Some(campaign) = self.campaign_address {
                        print!("Use current campaign ({})? (y/n): ", campaign);
                        io::stdout().flush()?;
                        let mut response = String::new();
                        io::stdin().read_line(&mut response)?;
                        
                        if response.trim().to_lowercase() == "y" {
                            campaign.to_string()
                        } else {
                            print!("Campaign address: ");
                            io::stdout().flush()?;
                            let mut addr = String::new();
                            io::stdin().read_line(&mut addr)?;
                            addr.trim().to_string()
                        }
                    } else {
                        print!("Campaign address: ");
                        io::stdout().flush()?;
                        let mut addr = String::new();
                        io::stdin().read_line(&mut addr)?;
                        addr.trim().to_string()
                    };

                    print!("Amount (lamports): ");
                    io::stdout().flush()?;
                    let mut amount_str = String::new();
                    io::stdin().read_line(&mut amount_str)?;
                    
                    match amount_str.trim().parse::<u64>() {
                        Ok(amount) if amount > 0 => {
                            if let Err(e) = self.withdraw_from_campaign(&address, amount).await {
                                if e.to_string().contains("Unauthorized") || e.to_string().contains("6000") {
                                    println!("❌ Unauthorized: You are not the admin of this campaign.");
                                } else if e.to_string().contains("InsufficientFunds") || e.to_string().contains("6001") {
                                    println!("❌ Insufficient funds in the campaign to withdraw this amount.");
                                } else {
                                    println!("❌ Error withdrawing: {}", e);
                                }
                            } else {
                                println!("✅ Successfully withdrew {} lamports!", amount);
                            }
                        }
                        _ => println!("❌ Invalid amount. Please enter a valid number greater than 0."),
                    }
                }
                "5" => {
                    match self.get_balance().await {
                        Ok(balance) => println!("Current balance: {:.4} SOL", balance),
                        Err(e) => println!("Error getting balance: {}", e),
                    }
                }
                "6" => {
                    if let Err(e) = self.check_campaign_status().await {
                        println!("❌ Error checking campaign status: {}", e);
                    }
                }
                "7" => {
                    println!("Goodbye!");
                    return Ok(());
                }
                _ => println!("❌ Invalid choice. Please enter a number between 1-7."),
            }

            print!("\nPress Enter to continue...");
            io::stdout().flush()?;
            let mut _temp = String::new();
            io::stdin().read_line(&mut _temp)?;
        }
    }
}

#[tokio::main]
async fn main() -> Result<()> {
    // Use Go client wallet by default, or allow override via command line
    let key_path = std::env::args().nth(1)
        .unwrap_or_else(|| "../go_client/my_wallet.json".to_string());
    
    println!("🚀 Solana dApp CLI Starting...");
    
    let mut app = SolanaDApp::new(Some(&key_path))?;
    
    println!("✅ Connected to Solana devnet");
    println!("💳 Wallet loaded: {}", app.wallet.pubkey());
    
    app.run().await
}
