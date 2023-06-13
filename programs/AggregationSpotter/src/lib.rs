use anchor_lang::{
    prelude::*,
    solana_program,
    solana_program::{pubkey, pubkey::Pubkey},
};

declare_id!("3XQdG1Zpk151xuGHSd6DUkNuh9m9i3M8ptxJEHNfZdJ2");
pub const OWNER: Pubkey = pubkey!("CuK4CzZFFQaK2KaUYYyNodQ6ZG6PTv1jjYKvqHUx7P5Y");
pub const KEEPER_SEED: &str = "keeper";
pub const CONTRACT_SEED: &str = "keeper";
pub const DESCRIMINATOR_LEN: usize = 8;

#[program]
pub mod aggregation_spotter {
    use super::*;

    pub fn initialize(ctx: Context<Initialize>, init_addr: [Pubkey; 2]) -> Result<()> {
        ctx.accounts.spotter.is_initialized = true;
        ctx.accounts.spotter.admin = init_addr[0];
        ctx.accounts.spotter.executor = init_addr[1];
        ctx.accounts.spotter.rate_decimals = 8;
        ctx.accounts.spotter.number_of_keepers = 0;

        Ok(())
    }

    pub fn add_keeper(ctx: Context<RemoveKeeper>, keeper: Pubkey) -> Result<()> {
        ctx.accounts.keeper_acc.key = keeper;
        ctx.accounts.keeper_acc.is_allowed = true;
        ctx.accounts.spotter.number_of_keepers += 1;

        Ok(())
    }

    pub fn remove_keeper(ctx: Context<AddKeeper>, keeper: Pubkey) -> Result<()> {
        if ctx.accounts.keeper_acc.key != keeper {
            return Err(SpotterError::WrongKeeperAccount.into());
        }

        ctx.accounts.keeper_acc.is_allowed = false;
        ctx.accounts.spotter.number_of_keepers -= 1;

        Ok(())
    }
}

#[error_code]
pub enum SpotterError {
    Unauthorized,
    WrongKeeperAccount
}

#[account]
pub struct Keeper {
    pub is_allowed: bool, // 1
    pub key: Pubkey,      // 32
}

impl Keeper {
    pub const MAXIMUM_SIZE: usize = DESCRIMINATOR_LEN + 1 + 32;
}

#[account]
pub struct AggregationSpotter {
    pub is_initialized: bool,   // 1
    pub admin: Pubkey,          // 32
    pub executor: Pubkey,       // 32
    pub number_of_keepers: u64, // 8
    pub rate_decimals: u8,      // 1
}

impl AggregationSpotter {
    pub const MAXIMUM_SIZE: usize = DESCRIMINATOR_LEN + 1 + 32 + 32 + 8 + 1;
}

#[derive(Accounts)]
pub struct Initialize<'info> {
    #[account(init, payer = admin, space = AggregationSpotter::MAXIMUM_SIZE)]
    pub spotter: Account<'info, AggregationSpotter>,
    #[account(mut, constraint = admin.key() == OWNER)]
    pub admin: Signer<'info>,
    pub system_program: Program<'info, System>,
}

#[derive(Accounts)]
pub struct AddKeeper<'info> {
    #[account(mut, has_one = admin)]
    pub spotter: Account<'info, AggregationSpotter>,
    #[account(mut, constraint = admin.key() == OWNER)]
    pub admin: Signer<'info>,
    #[account(init, payer = admin, space = Keeper::MAXIMUM_SIZE, seeds=[KEEPER_SEED.as_ref(), keeper_acc.key.as_ref()], bump)]
    pub keeper_acc: Account<'info, Keeper>,
    pub system_program: Program<'info, System>,
}

#[derive(Accounts)]
pub struct RemoveKeeper<'info> {
    #[account(mut, has_one = admin)]
    pub spotter: Account<'info, AggregationSpotter>,
    #[account(constraint = admin.key() == OWNER)]
    pub admin: Signer<'info>,
    #[account(mut, seeds=[KEEPER_SEED.as_ref(), keeper_acc.key.as_ref()], bump)]
    pub keeper_acc: Account<'info, Keeper>,
    pub system_program: Program<'info, System>,
}
