use std::collections::HashMap;
use std::sync::{Mutex, LazyLock};
use serde::{Serialize, Deserialize};
use borsh::{BorshSerialize, BorshDeserialize};

// Use LazyLock for static initialization
static STORAGE: LazyLock<Mutex<HashMap<String, Vec<u8>>>> = LazyLock::new(|| {
    Mutex::new(HashMap::new())
});

pub struct Storage;

impl Storage {
    /// Get a value from storage
    pub fn get<T>(key: &str) -> Option<T> 
    where 
        T: BorshDeserialize,
    {
        let storage = STORAGE.lock().unwrap();
        if let Some(bytes) = storage.get(key) {
            T::try_from_slice(bytes).ok()
        } else {
            None
        }
    }

    /// Set a value in storage
    pub fn set<T>(key: &str, value: T) 
    where 
        T: BorshSerialize,
    {
        let mut storage = STORAGE.lock().unwrap();
        if let Ok(serialized) = borsh::to_vec(&value) {
            storage.insert(key.to_string(), serialized);
        }
    }

    /// Get a string from storage with default
    pub fn get_string(key: &str, default_value: &str) -> String {
        Self::get(key).unwrap_or_else(|| default_value.to_string())
    }

    /// Set a string in storage
    pub fn set_string(key: &str, value: &str) {
        Self::set(key, value.to_string());
    }

    /// Get a u64 from storage with default
    pub fn get_u64(key: &str, default_value: u64) -> u64 {
        Self::get(key).unwrap_or(default_value)
    }

    /// Set a u64 in storage
    pub fn set_u64(key: &str, value: u64) {
        Self::set(key, value);
    }

    /// Get a boolean from storage with default
    pub fn get_bool(key: &str, default_value: bool) -> bool {
        Self::get(key).unwrap_or(default_value)
    }

    /// Set a boolean in storage
    pub fn set_bool(key: &str, value: bool) {
        Self::set(key, value);
    }

    /// Check if a key exists in storage
    pub fn exists(key: &str) -> bool {
        let storage = STORAGE.lock().unwrap();
        storage.contains_key(key)
    }

    /// Remove a key from storage
    pub fn remove(key: &str) {
        let mut storage = STORAGE.lock().unwrap();
        storage.remove(key);
    }

    /// Clear all storage (testing only)
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
        assert!(Storage::exists("test_u64"));
        assert!(!Storage::exists("nonexistent"));

        // Test remove
        Storage::remove("test_u64");
        assert!(!Storage::exists("test_u64"));
    }
} 