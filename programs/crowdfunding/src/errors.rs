use anchor_lang::prelude::*;


#[error_code]
pub enum CampaignError {
    #[msg("You are not the admin of this campaign.")]
    Unauthorized,
    #[msg("Insufficient funds to perform this action.")]
    InsufficientFunds,
}
