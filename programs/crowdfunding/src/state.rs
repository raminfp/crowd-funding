use anchor_lang::prelude::*;

#[derive(Accounts)]
pub struct Create<'info> {
    #[account(
        init,
        payer = user,
        space = 9000,
        seeds = [b"CAMPAIGN_DEMO".as_ref(), user.key().as_ref()],
        bump
    )]
    pub campaign: Account<'info, Campaign>,
    #[account(mut)]
    pub user: Signer<'info>,
    pub system_program: Program<'info, System>,
}

#[derive(Accounts)]
pub struct Withdraw<'info> {
    #[account(mut)]
    pub campaign: Account<'info, Campaign>,
    #[account(mut)]
    pub user: Signer<'info>,
}

#[derive(Accounts)]
pub struct Donate<'info> {
    #[account(mut)]
    pub campaign: Account<'info, Campaign>,
    #[account(mut)]
    pub user: Signer<'info>,
    pub system_program: Program<'info, System>,
}

#[account]
pub struct Campaign {
    pub admin: Pubkey,        // 32 bytes
    pub name: String,         // dynamic
    pub description: String,  // dynamic
    pub amount_donated: u64,  // 8 bytes
}
