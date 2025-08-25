// use anchor_lang::prelude::*;
//
// declare_id!("7XJkGrdSHn3chc7rsv1xDzEKtwP9w5rSx1shohzM5skv");
//
// #[program]
// pub mod crowdfunding {
//     use anchor_lang::error::Error::ProgramError;
//     use super::*;
//
//     pub fn create(
//         ctx: Context<Create>,
//         name: String,
//         description: String,
//     ) -> Result<()> {
//         let campaign = &mut ctx.accounts.campaign;
//         campaign.name = name;
//         campaign.description = description;
//         campaign.amount_donated = 0;
//         campaign.admin = *ctx.accounts.user.key;
//         Ok(())
//     }
//     pub fn withdraw(ctx: Context<Withdraw>, amount: u64) -> Result<()> {
//         let campaign = &mut ctx.accounts.campaign;
//         let user = &mut ctx.accounts.user;
//         if campaign.admin != *user.key {
//             return Err(CampaignError::Unauthorized.into());
//         }
//         let rent_balance = Rent::get()?.minimum_balance(campaign.to_account_info().data_len());
//
//         if **campaign.to_account_info().lamports.borrow() - rent_balance <  amount {
//             return Err(CampaignError::InsufficientFunds.into());
//         }
//         **campaign.to_account_info().try_borrow_mut_lamports()? -= amount;
//         **user.to_account_info().try_borrow_mut_lamports()? += amount;
//         Ok(())
//     }
//
//     pub fn donate(ctx: Context<Donate>, amount: u64) -> Result<()> {
//         let ix = anchor_lang::solana_program::system_instruction::transfer(
//             &ctx.accounts.user.key(),
//             &ctx.accounts.campaign.key(),
//             amount,
//         );
//         let _ = anchor_lang::solana_program::program::invoke(
//             &ix,
//             &[
//                 ctx.accounts.user.to_account_info(),
//                 ctx.accounts.campaign.to_account_info()
//             ]
//         );
//         (&mut ctx.accounts.campaign).amount_donated += amount;
//         Ok(())
//     }
// }
//
// #[derive(Accounts)]
// pub struct Create<'info> {
//     #[account(
//         init,
//         payer = user,
//         space = 9000,
//         seeds = [b"CAMPAIGN_DEMO".as_ref(), user.key().as_ref()],
//         bump
//     )]
//     pub campaign: Account<'info, Campaign>,
//     #[account(mut)]
//     pub user: Signer<'info>,
//     pub system_program: Program<'info, System>,
// }
//
// #[derive(Accounts)]
// pub struct Withdraw<'info> {
//     #[account(mut)]
//     pub campaign: Account<'info, Campaign>,
//     #[account(mut)]
//     pub user : Signer<'info>,
//
// }
//
// #[derive(Accounts)]
// pub struct Donate<'info> {
//     #[account(mut)]
//     pub campaign: Account<'info, Campaign>,
//     #[account(mut)]
//     pub user : Signer<'info>,
//     pub system_program: Program<'info, System>,
// }
//
// #[account]
// pub struct Campaign {
//     pub admin: Pubkey,        // 32 bytes
//     pub name: String,         // dynamic
//     pub description: String,  // dynamic
//     pub amount_donated: u64,  // 8 bytes
// }
//
// #[error_code]
// pub enum CampaignError {
//     #[msg("You are not the admin of this campaign.")]
//     Unauthorized,
//     #[msg("Insufficient funds to perform this action.")]
//     InsufficientFunds,
// }

use anchor_lang::prelude::*;

pub mod instructions;
pub mod state;
pub mod errors;

use instructions::*;
use state::*;
use errors::*;

declare_id!("7XJkGrdSHn3chc7rsv1xDzEKtwP9w5rSx1shohzM5skv");

#[program]
pub mod crowdfunding {
    use super::*;

    pub fn create(ctx: Context<Create>, name: String, description: String) -> Result<()> {
        instructions::create(ctx, name, description)
    }

    pub fn withdraw(ctx: Context<Withdraw>, amount: u64) -> Result<()> {
        instructions::withdraw(ctx, amount)
    }

    pub fn donate(ctx: Context<Donate>, amount: u64) -> Result<()> {
        instructions::donate(ctx, amount)
    }
}
