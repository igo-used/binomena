use std::collections::HashMap;
use std::sync::Mutex;
use serde::{Serialize, Deserialize};
use borsh::{BorshSerialize, BorshDeserialize};

// Storage abstraction for the contract
// In a real blockchain implementation, this would interface with the blockchain's storage layer
static STORAGE: Mutex<HashMap<String, Vec<u8>>> = Mutex::new(HashMap::new());

pub struct Storage;

impl Storage {
    /// Get a value from storage with a default if not found
    pub fn get<T>(key: &str, default_value: T) -> T 
    where 
        T: BorshDeserialize + Clone
    {
        let storage = STORAGE.lock().unwrap();
        
        if let Some(data) = storage.get(key) {
            match T::try_from_slice(data) {
                Ok(value) => value,
                Err(_) => default_value,
            }
        } else {
            default_value
        }
    }

    /// Set a value in storage
    pub fn set<T>(key: &str, value: T) 
    where 
        T: BorshSerialize
    {
        let mut storage = STORAGE.lock().unwrap();
        
        if let Ok(serialized) = value.try_to_vec() {
            storage.insert(key.to_string(), serialized);
        }
    }

    /// Get a string from storage with default
    pub fn get_string(key: &str, default_value: &str) -> String {
        let storage = STORAGE.lock().unwrap();
        
        if let Some(data) = storage.get(key) {
            String::from_utf8(data.clone()).unwrap_or_else(|_| default_value.to_string())
        } else {
            default_value.to_string()
        }
    }

    /// Set a string in storage
    pub fn set_string(key: &str, value: &str) {
        let mut storage = STORAGE.lock().unwrap();
        storage.insert(key.to_string(), value.as_bytes().to_vec());
    }

    /// Get a u64 from storage with default
    pub fn get_u64(key: &str, default_value: u64) -> u64 {
        Self::get(key, default_value)
    }

    /// Set a u64 in storage
    pub fn set_u64(key: &str, value: u64) {
        Self::set(key, value)
    }

    /// Get a boolean from storage with default
    pub fn get_bool(key: &str, default_value: bool) -> bool {
        Self::get(key, default_value)
    }

    /// Set a boolean in storage
    pub fn set_bool(key: &str, value: bool) {
        Self::set(key, value)
    }

    /// Check if a key exists in storage
    pub fn contains_key(key: &str) -> bool {
        let storage = STORAGE.lock().unwrap();
        storage.contains_key(key)
    }

    /// Remove a key from storage
    pub fn remove(key: &str) {
        let mut storage = STORAGE.lock().unwrap();
        storage.remove(key);
    }

    /// Clear all storage (useful for testing)
    #[cfg(test)]
    pub fn clear() {
        let mut storage = STORAGE.lock().unwrap();
        storage.clear();
    }
}

// Storage keys constants
pub struct StorageKeys;

impl StorageKeys {
    pub const OWNER: &'static str = "owner";
    pub const TOTAL_SUPPLY: &'static str = "totalSupply";
    pub const BALANCES_PREFIX: &'static str = "balance_";
    pub const PAUSED: &'static str = "paused";
    pub const BLACKLIST_PREFIX: &'static str = "blacklist_";
    pub const MINTER_PREFIX: &'static str = "minter_";
    pub const COLLATERAL_PREFIX: &'static str = "collateral_";
    pub const COLLATERAL_TYPE_PREFIX: &'static str = "collateral_type_";
    pub const BINOM_TOKEN_ADDRESS: &'static str = "binom_token_address";
    pub const COLLATERAL_RATIO: &'static str = "collateral_ratio";
    pub const FIAT_RESERVE: &'static str = "fiat_reserve";
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_storage_operations() {
        Storage::clear();
        
        // Test u64 operations
        Storage::set_u64("test_u64", 42);
        assert_eq!(Storage::get_u64("test_u64", 0), 42);
        assert_eq!(Storage::get_u64("nonexistent", 100), 100);

        // Test string operations
        Storage::set_string("test_string", "hello");
        assert_eq!(Storage::get_string("test_string", "default"), "hello");
        assert_eq!(Storage::get_string("nonexistent", "default"), "default");

        // Test bool operations
        Storage::set_bool("test_bool", true);
        assert_eq!(Storage::get_bool("test_bool", false), true);
        assert_eq!(Storage::get_bool("nonexistent", false), false);

        // Test contains_key
        assert!(Storage::contains_key("test_u64"));
        assert!(!Storage::contains_key("nonexistent"));

        // Test remove
        Storage::remove("test_u64");
        assert!(!Storage::contains_key("test_u64"));
    }
} 