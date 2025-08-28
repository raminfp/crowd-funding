package main

import (
	"bufio"
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
)

const (
	ProgramID = "3r5NUnG85XtVExb1234ZYYyUazjchqjfYknnQATyCDzp"
	Network   = rpc.DevNet_RPC
)

// generateDiscriminator creates an 8-byte discriminator for Anchor instructions
func generateDiscriminator(namespace, name string) []byte {
	preimage := fmt.Sprintf("%s:%s", namespace, name)
	hash := sha256.Sum256([]byte(preimage))
	return hash[:8]
}

// Campaign represents the campaign account structure
type Campaign struct {
	Admin         solana.PublicKey `json:"admin"`
	Name          string           `json:"name"`
	Description   string           `json:"description"`
	AmountDonated uint64           `json:"amount_donated"`
	Bump          uint8            `json:"bump"`
}

// SolanaDApp represents our dApp instance
type SolanaDApp struct {
	client          *rpc.Client
	wsClient        *ws.Client
	wallet          *Wallet
	programID       solana.PublicKey
	campaignAddress *solana.PublicKey // Current campaign address
	campaignName    string            // Current campaign name
}

// Wallet represents a Solana wallet
type Wallet struct {
	PublicKey  solana.PublicKey
	PrivateKey ed25519.PrivateKey
}

// WalletData represents the wallet file format
type WalletData struct {
	PublicKey  string `json:"publicKey,omitempty"`
	PrivateKey string `json:"privateKey,omitempty"`
}

// NewWallet creates a new wallet from a private key file or generates one
func NewWallet(keyPath string) (*Wallet, error) {
	var privateKey ed25519.PrivateKey

	if keyPath != "" {
		// Load existing key
		keyData, err := os.ReadFile(keyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read key file: %w", err)
		}

		// Try to parse as wallet data with base58 keys first
		var walletData WalletData
		if err := json.Unmarshal(keyData, &walletData); err == nil && walletData.PrivateKey != "" {
			// Parse base58 private key
			privKeyBytes, err := solana.PrivateKeyFromBase58(walletData.PrivateKey)
			if err != nil {
				return nil, fmt.Errorf("failed to parse base58 private key: %w", err)
			}
			privateKey = ed25519.PrivateKey(privKeyBytes)
		} else {
			// Try to parse as byte array (legacy format)
			var keyArray []byte
			if err := json.Unmarshal(keyData, &keyArray); err != nil {
				return nil, fmt.Errorf("failed to parse key file: %w", err)
			}

			if len(keyArray) != 64 {
				return nil, fmt.Errorf("invalid key length: expected 64, got %d", len(keyArray))
			}

			privateKey = ed25519.PrivateKey(keyArray)
		}
	} else {
		// Generate new key
		_, privateKey, _ = ed25519.GenerateKey(nil)

		// Save key to file
		keyBytes, _ := json.Marshal([]byte(privateKey))
		if err := os.WriteFile("wallet.json", keyBytes, 0600); err != nil {
			log.Printf("Warning: failed to save wallet key: %v", err)
		} else {
			fmt.Println("New wallet saved to wallet.json")
		}
	}

	publicKey := privateKey.Public().(ed25519.PublicKey)

	return &Wallet{
		PublicKey:  solana.PublicKeyFromBytes(publicKey),
		PrivateKey: privateKey,
	}, nil
}

// NewSolanaDApp creates a new instance of the Solana dApp
func NewSolanaDApp(keyPath string) (*SolanaDApp, error) {
	client := rpc.New(Network)
	wsClient, err := ws.Connect(context.Background(), rpc.DevNet_WS)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to WebSocket: %w", err)
	}

	wallet, err := NewWallet(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create wallet: %w", err)
	}

	programID := solana.MustPublicKeyFromBase58(ProgramID)

	app := &SolanaDApp{
		client:    client,
		wsClient:  wsClient,
		wallet:    wallet,
		programID: programID,
	}

	// Try to load saved campaign address
	app.loadSavedCampaign()

	// Note: We can't check for existing campaigns here without campaign name
	// Users will need to provide campaign name to check for existing campaigns

	return app, nil
}

// SavedCampaign represents saved campaign data
type SavedCampaign struct {
	Address string `json:"address"`
	Name    string `json:"name"`
}

// loadSavedCampaign tries to load a previously saved campaign address and name
func (app *SolanaDApp) loadSavedCampaign() {
	data, err := os.ReadFile("campaign.txt")
	if err != nil {
		return // No saved campaign, which is fine
	}

	campaignStr := strings.TrimSpace(string(data))
	if campaignStr == "" {
		return
	}

	// Try to parse as JSON first (new format)
	var savedCampaign SavedCampaign
	if err := json.Unmarshal([]byte(campaignStr), &savedCampaign); err == nil {
		// New format with name
		campaignPubkey, err := solana.PublicKeyFromBase58(savedCampaign.Address)
		if err != nil {
			log.Printf("Warning: invalid saved campaign address: %v", err)
			return
		}
		app.campaignAddress = &campaignPubkey
		app.campaignName = savedCampaign.Name
		fmt.Printf("📋 Loaded saved campaign '%s': %s\n", savedCampaign.Name, savedCampaign.Address)
	} else {
		// Old format - just address
		campaignPubkey, err := solana.PublicKeyFromBase58(campaignStr)
		if err != nil {
			log.Printf("Warning: invalid saved campaign address: %v", err)
			return
		}
		app.campaignAddress = &campaignPubkey
		app.campaignName = "" // Unknown name for old saves
		fmt.Printf("📋 Loaded saved campaign: %s (name unknown)\n", campaignStr)
	}
}

// saveCampaign saves the current campaign address and name to a file
func (app *SolanaDApp) saveCampaign() {
	if app.campaignAddress == nil {
		return
	}

	savedCampaign := SavedCampaign{
		Address: app.campaignAddress.String(),
		Name:    app.campaignName,
	}

	data, err := json.Marshal(savedCampaign)
	if err != nil {
		log.Printf("Warning: failed to marshal campaign data: %v", err)
		return
	}

	err = os.WriteFile("campaign.txt", data, 0644)
	if err != nil {
		log.Printf("Warning: failed to save campaign data: %v", err)
	}
}

// GetBalance returns the wallet's SOL balance
func (app *SolanaDApp) GetBalance() (float64, error) {
	balance, err := app.client.GetBalance(
		context.Background(),
		app.wallet.PublicKey,
		rpc.CommitmentFinalized,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to get balance: %w", err)
	}

	return float64(balance.Value) / float64(solana.LAMPORTS_PER_SOL), nil
}

// RequestAirdrop requests SOL from the devnet faucet
func (app *SolanaDApp) RequestAirdrop() error {
	fmt.Println("Requesting airdrop...")

	sig, err := app.client.RequestAirdrop(
		context.Background(),
		app.wallet.PublicKey,
		2*solana.LAMPORTS_PER_SOL, // 2 SOL
		rpc.CommitmentFinalized,
	)
	if err != nil {
		return fmt.Errorf("failed to request airdrop: %w", err)
	}

	fmt.Printf("Airdrop requested. Transaction signature: %s\n", sig)
	fmt.Println("Waiting for confirmation...")

	// Wait for confirmation
	ctx := context.Background()
	status, err := app.client.GetSignatureStatuses(ctx, true, sig)
	if err != nil {
		return fmt.Errorf("failed to confirm airdrop: %w", err)
	}

	if len(status.Value) > 0 && status.Value[0] != nil && status.Value[0].Err != nil {
		return fmt.Errorf("airdrop transaction failed: %v", status.Value[0].Err)
	}

	fmt.Println("✅ Airdrop confirmed!")
	return nil
}

// CreateCampaignPDA generates the Program Derived Address for a campaign
func (app *SolanaDApp) CreateCampaignPDA(campaignName string) (solana.PublicKey, uint8, error) {
	seeds := [][]byte{
		[]byte("CAMPAIGN_DEMO"),
		app.wallet.PublicKey.Bytes(),
		[]byte(campaignName),
	}

	return solana.FindProgramAddress(seeds, app.programID)
}

// CheckExistingCampaign checks if a properly initialized campaign already exists for this wallet and campaign name
func (app *SolanaDApp) CheckExistingCampaign(campaignName string) (*solana.PublicKey, error) {
	campaignPDA, _, err := app.CreateCampaignPDA(campaignName)
	if err != nil {
		return nil, fmt.Errorf("failed to create campaign PDA: %w", err)
	}

	// Check if the account exists and is properly initialized
	accountInfo, err := app.client.GetAccountInfo(context.Background(), campaignPDA)
	if err != nil {
		return nil, nil // Account doesn't exist
	}

	if accountInfo.Value == nil {
		return nil, nil // Account doesn't exist
	}

	// Check if the account is owned by our program (not just allocated by system program)
	if !accountInfo.Value.Owner.Equals(app.programID) {
		fmt.Printf("⚠️  Found uninitialized account at %s (owned by %s, not %s)\n",
			campaignPDA.String(), accountInfo.Value.Owner.String(), app.programID.String())
		return nil, nil // Account exists but not initialized by our program
	}

	// Check if the account has data (properly initialized campaign)
	if len(accountInfo.Value.Data.GetBinary()) < 32 { // Minimum size for a campaign account
		fmt.Printf("⚠️  Found account with insufficient data at %s\n", campaignPDA.String())
		return nil, nil // Account exists but not properly initialized
	}

	fmt.Printf("✅ Found properly initialized campaign at %s\n", campaignPDA.String())
	return &campaignPDA, nil
}

// CheckCampaignStatus provides detailed status information about the campaign account
func (app *SolanaDApp) CheckCampaignStatus(campaignName string) error {
	campaignPDA, _, err := app.CreateCampaignPDA(campaignName)
	if err != nil {
		return fmt.Errorf("failed to create campaign PDA: %w", err)
	}

	fmt.Printf("\n🔍 Campaign Status for Wallet: %s\n", app.wallet.PublicKey.String())
	fmt.Printf("📍 Expected Campaign Address: %s\n", campaignPDA.String())
	fmt.Printf("🔗 Explorer Link: https://explorer.solana.com/address/%s?cluster=devnet\n", campaignPDA.String())

	// Get account info
	accountInfo, err := app.client.GetAccountInfo(context.Background(), campaignPDA)
	if err != nil {
		fmt.Printf("❌ Account does not exist or error fetching: %v\n", err)
		fmt.Println("✅ You can create a new campaign!")
		return nil
	}

	if accountInfo.Value == nil {
		fmt.Println("❌ Account does not exist")
		fmt.Println("✅ You can create a new campaign!")
		return nil
	}

	fmt.Printf("📊 Account Info:\n")
	fmt.Printf("   Owner: %s\n", accountInfo.Value.Owner.String())
	fmt.Printf("   Data Size: %d bytes\n", len(accountInfo.Value.Data.GetBinary()))
	fmt.Printf("   Lamports: %d\n", accountInfo.Value.Lamports)

	if accountInfo.Value.Owner.Equals(solana.SystemProgramID) {
		fmt.Println("⚠️  Account is allocated but NOT initialized by the crowdfunding program")
		fmt.Println("💡 This means a previous campaign creation failed partway through")
		fmt.Println("🔧 The account exists but has no campaign data")
		fmt.Println("❗ You'll need to use a different wallet or wait for the account to be reclaimed")
	} else if accountInfo.Value.Owner.Equals(app.programID) {
		fmt.Println("✅ Account is properly owned by the crowdfunding program")
		if len(accountInfo.Value.Data.GetBinary()) >= 32 {
			fmt.Println("✅ Account appears to have campaign data")
			app.campaignAddress = &campaignPDA
			app.campaignName = campaignName
			app.saveCampaign()
		} else {
			fmt.Println("⚠️  Account is owned by program but has insufficient data")
		}
	} else {
		fmt.Printf("❓ Account is owned by unknown program: %s\n", accountInfo.Value.Owner.String())
	}

	return nil
}

// CreateCampaign creates a new fundraising campaign
func (app *SolanaDApp) CreateCampaign(name, description string) error {
	// First, check if a campaign already exists
	existingCampaign, err := app.CheckExistingCampaign(name)
	if err != nil {
		return fmt.Errorf("failed to check existing campaign: %w", err)
	}

	if existingCampaign != nil {
		fmt.Printf("✅ Campaign already exists at: %s\n", existingCampaign.String())
		app.campaignAddress = existingCampaign
		app.campaignName = name
		app.saveCampaign()
		fmt.Println("📋 Using existing campaign for future operations!")
		return nil
	}

	fmt.Printf("Creating campaign: %s\n", name)

	campaignPDA, _, err := app.CreateCampaignPDA(name)
	if err != nil {
		return fmt.Errorf("failed to create campaign PDA: %w", err)
	}

	// Build the instruction data for Anchor program
	// Generate the correct discriminator for the "create" instruction
	instructionData := generateDiscriminator("global", "create")

	// Serialize name length and name (u32 + string)
	nameLen := uint32(len(name))
	nameLenBytes := make([]byte, 4)
	for i := 0; i < 4; i++ {
		nameLenBytes[i] = byte(nameLen >> (i * 8))
	}
	instructionData = append(instructionData, nameLenBytes...)
	instructionData = append(instructionData, []byte(name)...)

	// Serialize description length and description (u32 + string)
	descLen := uint32(len(description))
	descLenBytes := make([]byte, 4)
	for i := 0; i < 4; i++ {
		descLenBytes[i] = byte(descLen >> (i * 8))
	}
	instructionData = append(instructionData, descLenBytes...)
	instructionData = append(instructionData, []byte(description)...)

	instruction := &solana.GenericInstruction{
		ProgID: app.programID,
		AccountValues: solana.AccountMetaSlice{
			{
				PublicKey:  campaignPDA,
				IsWritable: true,
				IsSigner:   false,
			},
			{
				PublicKey:  app.wallet.PublicKey,
				IsWritable: true,
				IsSigner:   true,
			},
			{
				PublicKey:  solana.SystemProgramID,
				IsWritable: false,
				IsSigner:   false,
			},
		},
		DataBytes: instructionData,
	}

	// Get latest blockhash
	recent, err := app.client.GetLatestBlockhash(context.Background(), rpc.CommitmentFinalized)
	if err != nil {
		return fmt.Errorf("failed to get latest blockhash: %w", err)
	}

	// Create transaction
	tx, err := solana.NewTransaction(
		[]solana.Instruction{instruction},
		recent.Value.Blockhash,
		solana.TransactionPayer(app.wallet.PublicKey),
	)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	// Sign transaction
	privKey := solana.PrivateKey(app.wallet.PrivateKey)
	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if key.Equals(app.wallet.PublicKey) {
			return &privKey
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Send transaction
	sig, err := app.client.SendTransaction(context.Background(), tx)
	if err != nil {
		return fmt.Errorf("failed to send transaction: %w", err)
	}

	fmt.Printf("Campaign created! Transaction: %s\n", sig)
	fmt.Printf("Campaign address: %s\n", campaignPDA.String())

	// Store the campaign address and name for future use
	app.campaignAddress = &campaignPDA
	app.campaignName = name
	app.saveCampaign()
	fmt.Printf("✅ Campaign address and name saved for quick access!\n")

	return nil
}

// DonateToCampaign donates SOL to a campaign
func (app *SolanaDApp) DonateToCampaign(campaignName, campaignAddress string, amount uint64) error {
	fmt.Printf("Donating %d lamports to campaign %s\n", amount, campaignAddress)

	campaignPubkey := solana.MustPublicKeyFromBase58(campaignAddress)

	// Build donate instruction with proper discriminator
	instructionData := generateDiscriminator("global", "donate")
	// Add name length and name (u32 + string)
	nameLen := uint32(len(campaignName))
	nameLenBytes := make([]byte, 4)
	for i := 0; i < 4; i++ {
		nameLenBytes[i] = byte(nameLen >> (i * 8))
	}
	instructionData = append(instructionData, nameLenBytes...)
	instructionData = append(instructionData, []byte(campaignName)...)
	// Add amount as 8 bytes (little endian)
	amountBytes := make([]byte, 8)
	for i := 0; i < 8; i++ {
		amountBytes[i] = byte(amount >> (i * 8))
	}
	instructionData = append(instructionData, amountBytes...)

	instruction := &solana.GenericInstruction{
		ProgID: app.programID,
		AccountValues: solana.AccountMetaSlice{
			{
				PublicKey:  campaignPubkey,
				IsWritable: true,
				IsSigner:   false,
			},
			{
				PublicKey:  app.wallet.PublicKey,
				IsWritable: true,
				IsSigner:   true,
			},
			{
				PublicKey:  solana.SystemProgramID,
				IsWritable: false,
				IsSigner:   false,
			},
		},
		DataBytes: instructionData,
	}

	// Get recent blockhash and send transaction
	return app.sendTransaction([]solana.Instruction{instruction})
}

// WithdrawFromCampaign withdraws SOL from a campaign (only campaign admin can do this)
func (app *SolanaDApp) WithdrawFromCampaign(campaignName, campaignAddress string, amount uint64) error {
	fmt.Printf("Withdrawing %d lamports from campaign %s\n", amount, campaignAddress)

	campaignPubkey := solana.MustPublicKeyFromBase58(campaignAddress)

	// Build withdraw instruction with proper discriminator
	instructionData := generateDiscriminator("global", "withdraw")
	// Add name length and name (u32 + string)
	nameLen := uint32(len(campaignName))
	nameLenBytes := make([]byte, 4)
	for i := 0; i < 4; i++ {
		nameLenBytes[i] = byte(nameLen >> (i * 8))
	}
	instructionData = append(instructionData, nameLenBytes...)
	instructionData = append(instructionData, []byte(campaignName)...)
	// Add amount as 8 bytes (little endian)
	amountBytes := make([]byte, 8)
	for i := 0; i < 8; i++ {
		amountBytes[i] = byte(amount >> (i * 8))
	}
	instructionData = append(instructionData, amountBytes...)

	instruction := &solana.GenericInstruction{
		ProgID: app.programID,
		AccountValues: solana.AccountMetaSlice{
			{
				PublicKey:  campaignPubkey,
				IsWritable: true,
				IsSigner:   false,
			},
			{
				PublicKey:  app.wallet.PublicKey,
				IsWritable: true,
				IsSigner:   true,
			},
		},
		DataBytes: instructionData,
	}

	return app.sendTransaction([]solana.Instruction{instruction})
}

// sendTransaction is a helper method to send transactions
func (app *SolanaDApp) sendTransaction(instructions []solana.Instruction) error {
	recent, err := app.client.GetLatestBlockhash(context.Background(), rpc.CommitmentFinalized)
	if err != nil {
		return fmt.Errorf("failed to get latest blockhash: %w", err)
	}

	tx, err := solana.NewTransaction(
		instructions,
		recent.Value.Blockhash,
		solana.TransactionPayer(app.wallet.PublicKey),
	)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	privKey := solana.PrivateKey(app.wallet.PrivateKey)
	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if key.Equals(app.wallet.PublicKey) {
			return &privKey
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %w", err)
	}

	sig, err := app.client.SendTransaction(context.Background(), tx)
	if err != nil {
		return fmt.Errorf("failed to send transaction: %w", err)
	}

	fmt.Printf("Transaction sent: %s\n", sig)
	return nil
}

// ShowMenu displays the interactive menu
func (app *SolanaDApp) ShowMenu() {
	fmt.Println("\n=== Solana dApp CLI ===")
	fmt.Printf("Wallet: %s\n", app.wallet.PublicKey.String())

	balance, err := app.GetBalance()
	if err != nil {
		fmt.Printf("Balance: Error getting balance (%v)\n", err)
	} else {
		fmt.Printf("Balance: %.4f SOL\n", balance)
	}

	// Show current campaign if available
	if app.campaignAddress != nil {
		if app.campaignName != "" {
			fmt.Printf("Current Campaign: '%s' (%s)\n", app.campaignName, app.campaignAddress.String())
		} else {
			fmt.Printf("Current Campaign: %s (name unknown)\n", app.campaignAddress.String())
		}
	} else {
		fmt.Println("Current Campaign: None")
	}

	fmt.Println("\nOptions:")
	fmt.Println("1. Request Airdrop (2 SOL)")
	fmt.Println("2. Create Campaign")
	if app.campaignAddress != nil {
		fmt.Println("3. Donate to Campaign ⭐")
		fmt.Println("4. Withdraw from Campaign ⭐")
	} else {
		fmt.Println("3. Donate to Campaign")
		fmt.Println("4. Withdraw from Campaign")
	}
	fmt.Println("5. Check Balance")
	fmt.Println("6. Check Campaign Status")
	fmt.Println("7. Exit")
	fmt.Print("\nChoose an option (1-7): ")
}

// Run starts the interactive CLI
func (app *SolanaDApp) Run() {
	reader := bufio.NewReader(os.Stdin)

	for {
		app.ShowMenu()

		input, _ := reader.ReadString('\n')
		choice := strings.TrimSpace(input)

		switch choice {
		case "1":
			if err := app.RequestAirdrop(); err != nil {
				if strings.Contains(err.Error(), "airdrop") {
					fmt.Println("❌ Airdrop failed. You may have reached the rate limit. Try again later.")
				} else {
					fmt.Printf("❌ Error requesting airdrop: %v\n", err)
				}
			}
		case "2":
			fmt.Print("Campaign name: ")
			name, _ := reader.ReadString('\n')
			name = strings.TrimSpace(name)

			fmt.Print("Campaign description: ")
			description, _ := reader.ReadString('\n')
			description = strings.TrimSpace(description)

			if err := app.CreateCampaign(name, description); err != nil {
				if strings.Contains(err.Error(), "insufficient") {
					fmt.Println("❌ Insufficient SOL in your wallet. Please use option 1 to get SOL via airdrop.")
				} else {
					fmt.Printf("❌ Error creating campaign: %v\n", err)
				}
			}
		case "3":
			var address string
			var campaignName string

			if app.campaignAddress != nil && app.campaignName != "" {
				fmt.Printf("Use current campaign '%s' (%s)? (y/n): ", app.campaignName, app.campaignAddress.String())
				response, _ := reader.ReadString('\n')
				if strings.TrimSpace(strings.ToLower(response)) == "y" {
					address = app.campaignAddress.String()
					campaignName = app.campaignName
				} else {
					fmt.Print("Campaign address: ")
					address, _ = reader.ReadString('\n')
					address = strings.TrimSpace(address)

					fmt.Print("Campaign name: ")
					campaignName, _ = reader.ReadString('\n')
					campaignName = strings.TrimSpace(campaignName)
				}
			} else {
				fmt.Print("Campaign address: ")
				address, _ = reader.ReadString('\n')
				address = strings.TrimSpace(address)

				fmt.Print("Campaign name: ")
				campaignName, _ = reader.ReadString('\n')
				campaignName = strings.TrimSpace(campaignName)
			}

			if campaignName == "" {
				fmt.Println("❌ Campaign name cannot be empty.")
				continue
			}

			fmt.Print("Amount (lamports): ")
			amountStr, _ := reader.ReadString('\n')
			amountStr = strings.TrimSpace(amountStr)

			amount, err := strconv.ParseUint(amountStr, 10, 64)
			if err != nil {
				fmt.Println("❌ Invalid amount. Please enter a valid number.")
				continue
			}
			if amount == 0 {
				fmt.Println("❌ Amount must be greater than 0.")
				continue
			}

			if err := app.DonateToCampaign(campaignName, address, amount); err != nil {
				if strings.Contains(err.Error(), "insufficient") {
					fmt.Println("❌ Insufficient SOL for donation. Please check your balance or request an airdrop.")
				} else {
					fmt.Printf("❌ Error donating: %v\n", err)
				}
			} else {
				fmt.Printf("✅ Successfully donated %d lamports!\n", amount)
			}
		case "4":
			var address string
			var campaignName string

			if app.campaignAddress != nil && app.campaignName != "" {
				fmt.Printf("Use current campaign '%s' (%s)? (y/n): ", app.campaignName, app.campaignAddress.String())
				response, _ := reader.ReadString('\n')
				if strings.TrimSpace(strings.ToLower(response)) == "y" {
					address = app.campaignAddress.String()
					campaignName = app.campaignName
				} else {
					fmt.Print("Campaign address: ")
					address, _ = reader.ReadString('\n')
					address = strings.TrimSpace(address)

					fmt.Print("Campaign name: ")
					campaignName, _ = reader.ReadString('\n')
					campaignName = strings.TrimSpace(campaignName)
				}
			} else {
				fmt.Print("Campaign address: ")
				address, _ = reader.ReadString('\n')
				address = strings.TrimSpace(address)

				fmt.Print("Campaign name: ")
				campaignName, _ = reader.ReadString('\n')
				campaignName = strings.TrimSpace(campaignName)
			}

			if campaignName == "" {
				fmt.Println("❌ Campaign name cannot be empty.")
				continue
			}

			fmt.Print("Amount (lamports): ")
			amountStr, _ := reader.ReadString('\n')
			amountStr = strings.TrimSpace(amountStr)

			amount, err := strconv.ParseUint(amountStr, 10, 64)
			if err != nil {
				fmt.Println("❌ Invalid amount. Please enter a valid number.")
				continue
			}
			if amount == 0 {
				fmt.Println("❌ Amount must be greater than 0.")
				continue
			}

			if err := app.WithdrawFromCampaign(campaignName, address, amount); err != nil {
				if strings.Contains(err.Error(), "Unauthorized") || strings.Contains(err.Error(), "6000") {
					fmt.Println("❌ Unauthorized: You are not the admin of this campaign.")
				} else if strings.Contains(err.Error(), "InsufficientFunds") || strings.Contains(err.Error(), "6001") {
					fmt.Println("❌ Insufficient funds in the campaign to withdraw this amount.")
				} else {
					fmt.Printf("❌ Error withdrawing: %v\n", err)
				}
			} else {
				fmt.Printf("✅ Successfully withdrew %d lamports!\n", amount)
			}
		case "5":
			balance, err := app.GetBalance()
			if err != nil {
				fmt.Printf("Error getting balance: %v\n", err)
			} else {
				fmt.Printf("Current balance: %.4f SOL\n", balance)
			}
		case "6":
			fmt.Print("Campaign name: ")
			campaignName, _ := reader.ReadString('\n')
			campaignName = strings.TrimSpace(campaignName)
			if campaignName == "" {
				fmt.Println("❌ Campaign name cannot be empty.")
				continue
			}
			if err := app.CheckCampaignStatus(campaignName); err != nil {
				fmt.Printf("❌ Error checking campaign status: %v\n", err)
			}
		case "7":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("❌ Invalid choice. Please enter a number between 1-7.")
		}

		fmt.Print("\nPress Enter to continue...")
		reader.ReadString('\n')
	}
}

func main() {
	var keyPath string
	if len(os.Args) > 1 {
		keyPath = os.Args[1]
	}

	fmt.Println("🚀 Solana dApp CLI Starting...")

	app, err := NewSolanaDApp(keyPath)
	if err != nil {
		log.Fatalf("Failed to initialize dApp: %v", err)
	}
	defer app.wsClient.Close()

	fmt.Printf("✅ Connected to Solana devnet\n")
	fmt.Printf("💳 Wallet loaded: %s\n", app.wallet.PublicKey.String())

	// Show initial balance
	if balance, err := app.GetBalance(); err == nil {
		fmt.Printf("💰 Current balance: %.4f SOL\n", balance)
		if balance < 0.01 {
			fmt.Println("⚠️  Low balance! You may want to request an airdrop.")
		}
	}

	app.Run()
}
