use crate::storage::{Storage, StorageKeys};
use crate::context::Context;
use crate::events::*;
use crate::errors::*;
use crate::{require, safe_add, safe_sub, safe_mul};

pub struct PaperDollar;

impl PaperDollar {
    /// Initialize the contract
    pub fn new() -> Self {
        // Set the caller as the owner if no owner exists
        if !Storage::exists(StorageKeys::OWNER) {
            let caller = Context::caller();
            Storage::set_string(StorageKeys::OWNER, &caller);
            Storage::set_u64(StorageKeys::TOTAL_SUPPLY, 0);
            Storage::set_bool(StorageKeys::PAUSED, false);
            Storage::set_u64(StorageKeys::COLLATERAL_RATIO, 150); // Default 150%
            Storage::set_u64(StorageKeys::FIAT_RESERVE, 0);
        }
        
        PaperDollar
    }

    // Modifier functions for access control
    fn check_owner(&self) -> ContractResult<()> {
        let owner = Storage::get_string(StorageKeys::OWNER, "");
        require!(owner == Context::caller(), ContractError::OnlyOwner);
        Ok(())
    }

    fn check_not_paused(&self) -> ContractResult<()> {
        require!(!Storage::get_bool(StorageKeys::PAUSED, false), ContractError::ContractPaused);
        Ok(())
    }

    fn check_not_blacklisted(&self, address: &str) -> ContractResult<()> {
        require!(!self.is_blacklisted(address), ContractError::AddressBlacklisted);
        Ok(())
    }

    fn check_minter(&self) -> ContractResult<()> {
        require!(self.is_minter(&Context::caller()), ContractError::OnlyMinters);
        Ok(())
    }

    // View functions
    pub fn get_balance(&self, address: &str) -> u64 {
        let key = format!("{}{}", StorageKeys::BALANCES_PREFIX, address);
        Storage::get_u64(&key, 0)
    }

    pub fn get_total_supply(&self) -> u64 {
        Storage::get_u64(StorageKeys::TOTAL_SUPPLY, 0)
    }

    pub fn get_collateral_balance(&self, address: &str, collateral_type: CollateralType) -> u64 {
        let key = format!("{}{}_{}",
            StorageKeys::COLLATERAL_PREFIX,
            collateral_type.to_string(),
            address
        );
        Storage::get_u64(&key, 0)
    }

    pub fn get_collateral_type(&self, address: &str) -> CollateralType {
        let key = format!("{}{}", StorageKeys::COLLATERAL_TYPE_PREFIX, address);
        let value = Storage::get_u64(&key, 0);
        CollateralType::from(value as u8)
    }

    pub fn get_collateral_ratio(&self) -> u64 {
        Storage::get_u64(StorageKeys::COLLATERAL_RATIO, 150)
    }

    pub fn get_fiat_reserve(&self) -> u64 {
        Storage::get_u64(StorageKeys::FIAT_RESERVE, 0)
    }

    pub fn is_blacklisted(&self, address: &str) -> bool {
        let key = format!("{}{}", StorageKeys::BLACKLIST_PREFIX, address);
        Storage::get_bool(&key, false)
    }

    pub fn is_minter(&self, address: &str) -> bool {
        let key = format!("{}{}", StorageKeys::MINTER_PREFIX, address);
        Storage::get_bool(&key, false)
    }

    pub fn get_owner(&self) -> String {
        Storage::get_string(StorageKeys::OWNER, "")
    }

    pub fn is_paused(&self) -> bool {
        Storage::get_bool(StorageKeys::PAUSED, false)
    }

    // Mutating functions

    /// Set BINOM token address (only owner)
    pub fn set_binom_token_address(&mut self, address: &str) -> ContractResult<bool> {
        self.check_owner()?;
        Storage::set_string(StorageKeys::BINOM_TOKEN_ADDRESS, address);
        Ok(true)
    }

    /// Set collateral ratio (only owner)
    pub fn set_collateral_ratio(&mut self, ratio: u64) -> ContractResult<bool> {
        self.check_owner()?;
        require!(ratio >= 100, ContractError::InvalidCollateralRatio);
        Storage::set_u64(StorageKeys::COLLATERAL_RATIO, ratio);
        Ok(true)
    }

    /// Add collateral
    pub fn add_collateral(&mut self, amount: u64, collateral_type: CollateralType) -> ContractResult<bool> {
        let from = Context::caller();
        require!(amount > 0, ContractError::ZeroAmount);
        
        self.check_not_paused()?;
        self.check_not_blacklisted(&from)?;

        if collateral_type == CollateralType::Fiat {
            let new_reserve = safe_add(self.get_fiat_reserve(), amount)?;
            Storage::set_u64(StorageKeys::FIAT_RESERVE, new_reserve);
        }

        let current_collateral = self.get_collateral_balance(&from, collateral_type);
        let new_collateral = safe_add(current_collateral, amount)?;
        
        let collateral_key = format!("{}{}_{}",
            StorageKeys::COLLATERAL_PREFIX,
            collateral_type.to_string(),
            from
        );
        Storage::set_u64(&collateral_key, new_collateral);

        let type_key = format!("{}{}", StorageKeys::COLLATERAL_TYPE_PREFIX, from);
        Storage::set_u64(&type_key, collateral_type as u64);

        emit_collateral_added(&from, amount, collateral_type);
        Ok(true)
    }

    /// Remove collateral
    pub fn remove_collateral(&mut self, amount: u64) -> ContractResult<bool> {
        let from = Context::caller();
        require!(amount > 0, ContractError::ZeroAmount);
        
        self.check_not_paused()?;
        self.check_not_blacklisted(&from)?;

        let collateral_type = self.get_collateral_type(&from);
        let current_collateral = self.get_collateral_balance(&from, collateral_type);
        require!(current_collateral >= amount, ContractError::InsufficientCollateral);

        // Check if removing collateral would break the collateral ratio
        let user_balance = self.get_balance(&from);
        let remaining_collateral = safe_sub(current_collateral, amount)?;
        let required_collateral = safe_mul(user_balance, self.get_collateral_ratio())?;
        let collateral_value = safe_mul(remaining_collateral, 100)?;
        
        require!(collateral_value >= required_collateral, ContractError::CollateralRatioTooLow);

        if collateral_type == CollateralType::Fiat {
            let new_reserve = safe_sub(self.get_fiat_reserve(), amount)?;
            Storage::set_u64(StorageKeys::FIAT_RESERVE, new_reserve);
        }

        let collateral_key = format!("{}{}_{}",
            StorageKeys::COLLATERAL_PREFIX,
            collateral_type.to_string(),
            from
        );
        Storage::set_u64(&collateral_key, remaining_collateral);

        emit_collateral_removed(&from, amount, collateral_type);
        Ok(true)
    }

    /// Transfer tokens
    pub fn transfer(&mut self, to: &str, amount: u64) -> ContractResult<bool> {
        let from = Context::caller();
        require!(amount > 0, ContractError::ZeroAmount);
        require!(from != to, ContractError::TransferToSelf);
        
        self.check_not_paused()?;
        self.check_not_blacklisted(&from)?;
        self.check_not_blacklisted(to)?;

        let from_balance = self.get_balance(&from);
        require!(from_balance >= amount, ContractError::InsufficientBalance);

        let new_from_balance = safe_sub(from_balance, amount)?;
        let to_balance = self.get_balance(to);
        let new_to_balance = safe_add(to_balance, amount)?;

        let from_key = format!("{}{}", StorageKeys::BALANCES_PREFIX, from);
        let to_key = format!("{}{}", StorageKeys::BALANCES_PREFIX, to);
        
        Storage::set_u64(&from_key, new_from_balance);
        Storage::set_u64(&to_key, new_to_balance);

        emit_transfer(&from, to, amount);
        Ok(true)
    }

    /// Mint new tokens (only minters)
    pub fn mint(&mut self, to: &str, amount: u64) -> ContractResult<bool> {
        require!(amount > 0, ContractError::ZeroAmount);
        
        self.check_minter()?;
        self.check_not_paused()?;
        self.check_not_blacklisted(to)?;

        let collateral_type = self.get_collateral_type(to);
        let collateral = self.get_collateral_balance(to, collateral_type);
        let current_balance = self.get_balance(to);
        let new_balance = safe_add(current_balance, amount)?;

        // Check collateral ratio
        let required_collateral = safe_mul(new_balance, self.get_collateral_ratio())?;
        let collateral_value = safe_mul(collateral, 100)?;
        require!(collateral_value >= required_collateral, ContractError::InsufficientCollateral);

        let new_total_supply = safe_add(self.get_total_supply(), amount)?;
        Storage::set_u64(StorageKeys::TOTAL_SUPPLY, new_total_supply);

        let balance_key = format!("{}{}", StorageKeys::BALANCES_PREFIX, to);
        Storage::set_u64(&balance_key, new_balance);

        emit_mint(to, amount, collateral_type);
        Ok(true)
    }

    /// Burn tokens
    pub fn burn(&mut self, amount: u64) -> ContractResult<bool> {
        let from = Context::caller();
        require!(amount > 0, ContractError::ZeroAmount);
        
        self.check_not_paused()?;
        self.check_not_blacklisted(&from)?;

        let from_balance = self.get_balance(&from);
        require!(from_balance >= amount, ContractError::InsufficientBalance);

        let new_balance = safe_sub(from_balance, amount)?;
        let new_total_supply = safe_sub(self.get_total_supply(), amount)?;

        let balance_key = format!("{}{}", StorageKeys::BALANCES_PREFIX, from);
        Storage::set_u64(&balance_key, new_balance);
        Storage::set_u64(StorageKeys::TOTAL_SUPPLY, new_total_supply);

        emit_burn(&from, amount);
        Ok(true)
    }

    /// Add a minter (only owner)
    pub fn add_minter(&mut self, address: &str) -> ContractResult<bool> {
        self.check_owner()?;
        let key = format!("{}{}", StorageKeys::MINTER_PREFIX, address);
        Storage::set_bool(&key, true);
        emit_minter_added(address);
        Ok(true)
    }

    /// Remove a minter (only owner)
    pub fn remove_minter(&mut self, address: &str) -> ContractResult<bool> {
        self.check_owner()?;
        let key = format!("{}{}", StorageKeys::MINTER_PREFIX, address);
        Storage::set_bool(&key, false);
        emit_minter_removed(address);
        Ok(true)
    }

    /// Blacklist an address (only owner)
    pub fn blacklist(&mut self, address: &str) -> ContractResult<bool> {
        self.check_owner()?;
        let key = format!("{}{}", StorageKeys::BLACKLIST_PREFIX, address);
        Storage::set_bool(&key, true);
        emit_blacklist_added(address);
        Ok(true)
    }

    /// Remove address from blacklist (only owner)
    pub fn unblacklist(&mut self, address: &str) -> ContractResult<bool> {
        self.check_owner()?;
        let key = format!("{}{}", StorageKeys::BLACKLIST_PREFIX, address);
        Storage::set_bool(&key, false);
        emit_blacklist_removed(address);
        Ok(true)
    }

    /// Pause the contract (only owner)
    pub fn pause(&mut self) -> ContractResult<bool> {
        self.check_owner()?;
        Storage::set_bool(StorageKeys::PAUSED, true);
        emit_paused();
        Ok(true)
    }

    /// Unpause the contract (only owner)
    pub fn unpause(&mut self) -> ContractResult<bool> {
        self.check_owner()?;
        Storage::set_bool(StorageKeys::PAUSED, false);
        emit_unpaused();
        Ok(true)
    }

    /// Transfer ownership (only owner)
    pub fn transfer_ownership(&mut self, new_owner: &str) -> ContractResult<bool> {
        self.check_owner()?;
        let old_owner = self.get_owner();
        Storage::set_string(StorageKeys::OWNER, new_owner);
        emit_ownership_transferred(&old_owner, new_owner);
        Ok(true)
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::with_caller;

    fn setup_contract() -> PaperDollar {
        Storage::clear();
        Context::reset();
        with_caller!("owner", {
            PaperDollar::new()
        })
    }

    #[test]
    fn test_initialization() {
        let contract = setup_contract();
        assert_eq!(contract.get_owner(), "owner");
        assert_eq!(contract.get_total_supply(), 0);
        assert_eq!(contract.get_collateral_ratio(), 150);
        assert!(!contract.is_paused());
    }

    #[test]
    fn test_add_collateral() {
        let mut contract = setup_contract();
        
        with_caller!("user1", {
            let result = contract.add_collateral(1000, CollateralType::Fiat);
            assert!(result.is_ok());
            assert_eq!(contract.get_collateral_balance("user1", CollateralType::Fiat), 1000);
            assert_eq!(contract.get_fiat_reserve(), 1000);
        });
    }

    #[test]
    fn test_mint_with_collateral() {
        let mut contract = setup_contract();
        
        // Add minter
        with_caller!("owner", {
            contract.add_minter("minter1").unwrap();
        });

        // User adds collateral
        with_caller!("user1", {
            contract.add_collateral(1500, CollateralType::Fiat).unwrap();
        });

        // Minter mints tokens
        with_caller!("minter1", {
            let result = contract.mint("user1", 1000);
            assert!(result.is_ok());
            assert_eq!(contract.get_balance("user1"), 1000);
            assert_eq!(contract.get_total_supply(), 1000);
        });
    }

    #[test]
    fn test_transfer() {
        let mut contract = setup_contract();
        
        // Setup: mint tokens to user1
        with_caller!("owner", {
            contract.add_minter("minter1").unwrap();
        });

        with_caller!("user1", {
            contract.add_collateral(1500, CollateralType::Fiat).unwrap();
        });

        with_caller!("minter1", {
            contract.mint("user1", 1000).unwrap();
        });

        // Transfer tokens
        with_caller!("user1", {
            let result = contract.transfer("user2", 500);
            assert!(result.is_ok());
            assert_eq!(contract.get_balance("user1"), 500);
            assert_eq!(contract.get_balance("user2"), 500);
        });
    }

    #[test]
    fn test_burn() {
        let mut contract = setup_contract();
        
        // Setup: mint tokens to user1
        with_caller!("owner", {
            contract.add_minter("minter1").unwrap();
        });

        with_caller!("user1", {
            contract.add_collateral(1500, CollateralType::Fiat).unwrap();
        });

        with_caller!("minter1", {
            contract.mint("user1", 1000).unwrap();
        });

        // Burn tokens
        with_caller!("user1", {
            let result = contract.burn(300);
            assert!(result.is_ok());
            assert_eq!(contract.get_balance("user1"), 700);
            assert_eq!(contract.get_total_supply(), 700);
        });
    }

    #[test]
    fn test_access_control() {
        let mut contract = setup_contract();
        
        // Non-owner cannot add minter
        with_caller!("user1", {
            let result = contract.add_minter("user2");
            assert!(matches!(result.unwrap_err(), ContractError::OnlyOwner));
        });

        // Non-minter cannot mint
        with_caller!("user1", {
            let result = contract.mint("user2", 100);
            assert!(matches!(result.unwrap_err(), ContractError::OnlyMinters));
        });
    }

    #[test]
    fn test_pause_functionality() {
        let mut contract = setup_contract();
        
        // Owner pauses contract
        with_caller!("owner", {
            contract.pause().unwrap();
            assert!(contract.is_paused());
        });

        // Operations should fail when paused
        with_caller!("user1", {
            let result = contract.add_collateral(1000, CollateralType::Fiat);
            assert!(matches!(result.unwrap_err(), ContractError::ContractPaused));
        });

        // Owner unpauses contract
        with_caller!("owner", {
            contract.unpause().unwrap();
            assert!(!contract.is_paused());
        });

        // Operations should work again
        with_caller!("user1", {
            let result = contract.add_collateral(1000, CollateralType::Fiat);
            assert!(result.is_ok());
        });
    }
} 