import './App.css';
import { useEffect, useState, useCallback } from "react";
import idl from "./idl.json";
import { Connection, PublicKey, clusterApiUrl } from "@solana/web3.js";
import { Program, AnchorProvider, web3, utils, BN } from "@coral-xyz/anchor";
import { Buffer } from 'buffer';

// Ensure Buffer is available globally
if (typeof window !== 'undefined') {
  window.Buffer = Buffer;
}

const programID = new PublicKey("7GMjTXTH1KS1Q46ngEnUYakAJi4xb2KJ3JsbJW2UNpHC");
const network = clusterApiUrl("devnet");
const opts = { preflightCommitment: "finalized" };
const { SystemProgram } = web3;

function App() {
  const [walletAddress, setWalletAddress] = useState(null);
  const [campaignAddress, setCampaignAddress] = useState(null);
  const [donationAmount, setDonationAmount] = useState("");
  const [withdrawAmount, setWithdrawAmount] = useState("");
  const [walletBalance, setWalletBalance] = useState(null);

  // Provider
  const getProvider = () => {
    const connection = new Connection(network, opts.preflightCommitment);
    const provider = new AnchorProvider(connection, window.solana, opts.preflightCommitment);
    return provider;
  };

  // Connect Phantom Wallet
  const connectWallet = async (onlyIfTrusted = false) => {
    try {
      if (window.solana && window.solana.isPhantom) {
        const resp = await window.solana.connect({ onlyIfTrusted });
        console.log("Wallet connected:", resp.publicKey.toString());
        setWalletAddress(resp.publicKey.toString());
      } else {
        alert("Phantom wallet not found! Please install it.");
      }
    } catch (err) {
      console.error("Wallet connection failed:", err);
      if (!onlyIfTrusted) {
        alert("Failed to connect wallet. Please try again.");
      }
    }
  };

  // Get wallet balance
  const getWalletBalance = useCallback(async () => {
    try {
      if (!window.solana || !window.solana.isConnected) {
        return;
      }

      const provider = getProvider();
      const balance = await provider.connection.getBalance(window.solana.publicKey);
      setWalletBalance(balance / web3.LAMPORTS_PER_SOL); // Convert to SOL
    } catch (error) {
      console.error("Error getting balance:", error);
    }
  }, []);

  // Request SOL from faucet (for devnet)
  const requestAirdrop = async () => {
    try {
      if (!window.solana || !window.solana.isConnected) {
        alert("Please connect your wallet first!");
        return;
      }

      const provider = getProvider();
      const signature = await provider.connection.requestAirdrop(
        window.solana.publicKey,
        2 * web3.LAMPORTS_PER_SOL // Request 2 SOL
      );
      
      await provider.connection.confirmTransaction(signature, opts.preflightCommitment);
      alert("Airdrop successful! You should receive 2 SOL shortly.");
      
      // Update balance after airdrop
      setTimeout(getWalletBalance, 2000);
    } catch (error) {
      console.error("Error requesting airdrop:", error);
      alert("Failed to request airdrop. Check console for details.");
    }
  };

  // Create Campaign
  const createCampaign = async () => {
    try {
      if (!window.solana || !window.solana.isConnected) {
        alert("Please connect your wallet first!");
        return;
      }

      const provider = getProvider();
      const program = new Program(idl, programID, provider);
      
      // Create the campaign PDA
      const [campaign] = await PublicKey.findProgramAddress(
          [
            utils.bytes.utf8.encode("CAMPAIGN_DEMO"),
            window.solana.publicKey.toBuffer(),
          ],
          programID
      );

      // Use the Anchor program interface to create the campaign
      await program.rpc.create("My Campaign", "This is a test campaign", {
        accounts: {
          campaign: campaign,
          user: window.solana.publicKey,
          systemProgram: SystemProgram.programId,
        },
      });
      
      console.log("Campaign created at:", campaign.toString());
      setCampaignAddress(campaign.toString());
    } catch (error) {
      console.error("Error creating campaign:", error);
      
      // Check if it's an insufficient SOL error
      if (error.message && error.message.includes("insufficient")) {
        alert("Insufficient SOL in your wallet. Please use the 'Get SOL (Airdrop)' button to get some SOL for transaction fees.");
      } else if (error.message && error.message.includes("Unexpected error")) {
        alert("Transaction failed. This might be due to insufficient SOL or network issues. Try getting SOL first with the airdrop button.");
      } else {
        alert("Failed to create campaign. Check console for details.");
      }
    }
  };

  // Donate to Campaign
  const donateToCampaign = async () => {
    if (!campaignAddress) return alert("No campaign selected!");
    try {
      const provider = getProvider();
      const program = new Program(idl, provider, programID);
      const amount = new BN(parseInt(donationAmount));

      await program.rpc.donate(amount, {
        accounts: {
          campaign: new PublicKey(campaignAddress),
          user: provider.wallet.publicKey,
          systemProgram: SystemProgram.programId,
        },
      });

      console.log(`Donated ${donationAmount} lamports to campaign ${campaignAddress}`);
      setDonationAmount("");
    } catch (error) {
      console.error("Error donating:", error);
      alert("Failed to donate. Check console for details.");
    }
  };

  // Withdraw from Campaign
  const withdrawFromCampaign = async () => {
    if (!campaignAddress) return alert("No campaign selected!");
    try {
      const provider = getProvider();
      const program = new Program(idl, provider, programID);
      const amount = new BN(parseInt(withdrawAmount));

      await program.rpc.withdraw(amount, {
        accounts: {
          campaign: new PublicKey(campaignAddress),
          user: provider.wallet.publicKey,
        },
      });

      console.log(`Withdrew ${withdrawAmount} lamports from campaign ${campaignAddress}`);
      setWithdrawAmount("");
    } catch (error) {
      console.error("Error withdrawing:", error);
      alert("Failed to withdraw. Check console for details.");
    }
  };

  useEffect(() => {
    const onLoad = async () => {
      await connectWallet(true);
    };
    onLoad();
  }, []);

  // Get balance when wallet connects
  useEffect(() => {
    if (walletAddress) {
      getWalletBalance();
    }
  }, [walletAddress, getWalletBalance]);

  return (
      <div className="App">
        <h1>Solana dApp</h1>
        {walletAddress ? (
            <div>
              <p>✅ Connected Wallet: {walletAddress}</p>
              {walletBalance !== null && (
                <p>💰 Balance: {walletBalance.toFixed(4)} SOL</p>
              )}
              <button onClick={requestAirdrop} style={{ marginRight: 10, backgroundColor: '#4CAF50', color: 'white' }}>
                Get SOL (Airdrop)
              </button>
              <button onClick={createCampaign}>Create Campaign</button>

              {campaignAddress && <p>Campaign: {campaignAddress}</p>}

              <div style={{ marginTop: 20 }}>
                <h3>Donate</h3>
                <input
                    type="number"
                    placeholder="Amount (lamports)"
                    value={donationAmount}
                    onChange={(e) => setDonationAmount(e.target.value)}
                />
                <button onClick={donateToCampaign}>Donate</button>
              </div>

              <div style={{ marginTop: 20 }}>
                <h3>Withdraw</h3>
                <input
                    type="number"
                    placeholder="Amount (lamports)"
                    value={withdrawAmount}
                    onChange={(e) => setWithdrawAmount(e.target.value)}
                />
                <button onClick={withdrawFromCampaign}>Withdraw</button>
              </div>
            </div>
        ) : (
            <button onClick={() => connectWallet(false)}>Connect Phantom Wallet</button>
        )}
      </div>
  );
}

export default App;
