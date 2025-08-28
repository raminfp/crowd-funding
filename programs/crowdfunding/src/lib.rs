use anchor_lang::prelude::*;

pub mod instructions;
pub mod state;
pub mod errors;

use instructions::*;
use state::*;
use errors::*;

declare_id!("3r5NUnG85XtVExb1234ZYYyUazjchqjfYknnQATyCDzp");

#[program]
pub mod crowdfunding {
    use super::*;

    pub fn create(ctx: Context<Create>, name: String, description: String) -> Result<()> {
        instructions::create(ctx, name, description)
    }

    pub fn withdraw(ctx: Context<Withdraw>, name: String, amount: u64) -> Result<()> {
        instructions::withdraw(ctx, name, amount)
    }

    pub fn donate(ctx: Context<Donate>, name: String, amount: u64) -> Result<()> {
        instructions::donate(ctx, name, amount)
    }
}
