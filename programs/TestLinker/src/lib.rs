use anchor_lang::prelude::*;

declare_id!("FBXArrxAnAZpmGmVYD8jQVauohCrXHaJBgTXzUPKRyHa");

#[program]
pub mod test_linker {
    use super::*;

    pub fn test_link(_ctx: Context<TestLink>) -> Result<()> {
        msg!("test_link");
        Ok(())
    }
}

#[derive(Accounts)]
pub struct TestLink<'info> {
    /// CHECK: We need this account only to create the PDA
    #[account(mut)]
    pub executer: AccountInfo<'info>,
}
