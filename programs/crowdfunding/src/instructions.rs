use anchor_lang::prelude::*;
use crate::{Campaign, CampaignError, Create, Withdraw, Donate};

pub fn create(ctx: Context<Create>, name: String, description: String) -> Result<()> {
    let campaign = &mut ctx.accounts.campaign;
    campaign.name = name;
    campaign.description = description;
    campaign.amount_donated = 0;
    campaign.admin = *ctx.accounts.user.key;
    campaign.bump = ctx.bumps.campaign;
    Ok(())
}

pub fn withdraw(ctx: Context<Withdraw>, name: String, amount: u64) -> Result<()> {
    let campaign = &mut ctx.accounts.campaign;
    let user = &mut ctx.accounts.user;
    
    if campaign.admin != *user.key {
        return Err(CampaignError::Unauthorized.into());
    }

    let rent_balance = Rent::get()?.minimum_balance(campaign.to_account_info().data_len());
    
    if **campaign.to_account_info().lamports.borrow() - rent_balance < amount {
        return Err(CampaignError::InsufficientFunds.into());
    }

    // Manual lamport transfer from PDA to user
    **campaign.to_account_info().try_borrow_mut_lamports()? -= amount;
    **user.to_account_info().try_borrow_mut_lamports()? += amount;

    Ok(())
}

pub fn donate(ctx: Context<Donate>, name: String, amount: u64) -> Result<()> {
    let ix = anchor_lang::solana_program::system_instruction::transfer(
        &ctx.accounts.user.key(),
        &ctx.accounts.campaign.key(),
        amount,
    );
    
    anchor_lang::solana_program::program::invoke(
        &ix,
        &[
            ctx.accounts.user.to_account_info(),
            ctx.accounts.campaign.to_account_info(),
            ctx.accounts.system_program.to_account_info()
        ]
    )?;
    
    (&mut ctx.accounts.campaign).amount_donated += amount;
    Ok(())
}
