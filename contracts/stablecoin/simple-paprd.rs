// Simple PAPRD Stablecoin Contract - No wasm-bindgen dependencies
// Compatible with Binomena blockchain

use std::collections::HashMap;

#[derive(Debug, Clone)]
pub struct PaprdContract {
    owner: String,
    name: String,
    symbol: String,
    decimals: u8,
    total_supply: u64,
    balances: HashMap<String, u64>,
    allowances: HashMap<String, HashMap<String, u64>>,
    minters: Vec<String>,
    paused: bool,
    blacklisted: Vec<String>,
}

impl PaprdContract {
    pub fn new(owner: String) -> Self {
        let mut contract = PaprdContract {
            owner: owner.clone(),
            name: "Paper Dollar Stablecoin".to_string(),
            symbol: "PAPRD".to_string(),
            decimals: 18,
            total_supply: 0,
            balances: HashMap::new(),
            allowances: HashMap::new(),
            minters: vec![owner.clone()],
            paused: false,
            blacklisted: Vec::new(),
        };
        
        // Mint initial supply to owner
        contract.mint(&owner, 100_000_000 * 10u64.pow(18)).unwrap();
        contract
    }

    // View functions
    pub fn get_name(&self) -> String {
        self.name.clone()
    }

    pub fn get_symbol(&self) -> String {
        self.symbol.clone()
    }

    pub fn get_decimals(&self) -> u8 {
        self.decimals
    }

    pub fn get_total_supply(&self) -> u64 {
        self.total_supply
    }

    pub fn get_balance(&self, address: &str) -> u64 {
        *self.balances.get(address).unwrap_or(&0)
    }

    pub fn get_owner(&self) -> String {
        self.owner.clone()
    }

    pub fn is_paused(&self) -> bool {
        self.paused
    }

    pub fn is_blacklisted(&self, address: &str) -> bool {
        self.blacklisted.contains(&address.to_string())
    }

    pub fn is_minter(&self, address: &str) -> bool {
        self.minters.contains(&address.to_string())
    }

    // State-changing functions
    pub fn transfer(&mut self, from: &str, to: &str, amount: u64) -> Result<(), String> {
        if self.paused {
            return Err("Contract is paused".to_string());
        }
        
        if self.is_blacklisted(from) || self.is_blacklisted(to) {
            return Err("Address is blacklisted".to_string());
        }

        if from == to {
            return Err("Transfer to self not allowed".to_string());
        }

        if amount == 0 {
            return Err("Zero amount not allowed".to_string());
        }

        let from_balance = self.get_balance(from);
        if from_balance < amount {
            return Err("Insufficient balance".to_string());
        }

        self.balances.insert(from.to_string(), from_balance - amount);
        let to_balance = self.get_balance(to);
        self.balances.insert(to.to_string(), to_balance + amount);

        Ok(())
    }

    pub fn mint(&mut self, to: &str, amount: u64) -> Result<(), String> {
        if self.paused {
            return Err("Contract is paused".to_string());
        }

        if self.is_blacklisted(to) {
            return Err("Address is blacklisted".to_string());
        }

        if amount == 0 {
            return Err("Zero amount not allowed".to_string());
        }

        let to_balance = self.get_balance(to);
        self.balances.insert(to.to_string(), to_balance + amount);
        self.total_supply += amount;

        Ok(())
    }

    pub fn burn(&mut self, from: &str, amount: u64) -> Result<(), String> {
        if self.paused {
            return Err("Contract is paused".to_string());
        }

        if amount == 0 {
            return Err("Zero amount not allowed".to_string());
        }

        let from_balance = self.get_balance(from);
        if from_balance < amount {
            return Err("Insufficient balance".to_string());
        }

        self.balances.insert(from.to_string(), from_balance - amount);
        self.total_supply -= amount;

        Ok(())
    }

    // Admin functions
    pub fn add_minter(&mut self, caller: &str, minter: &str) -> Result<(), String> {
        if caller != self.owner {
            return Err("Only owner can call this function".to_string());
        }

        if !self.minters.contains(&minter.to_string()) {
            self.minters.push(minter.to_string());
        }

        Ok(())
    }

    pub fn blacklist(&mut self, caller: &str, address: &str) -> Result<(), String> {
        if caller != self.owner {
            return Err("Only owner can call this function".to_string());
        }

        if !self.blacklisted.contains(&address.to_string()) {
            self.blacklisted.push(address.to_string());
        }

        Ok(())
    }

    pub fn pause(&mut self, caller: &str) -> Result<(), String> {
        if caller != self.owner {
            return Err("Only owner can call this function".to_string());
        }

        self.paused = true;
        Ok(())
    }

    pub fn unpause(&mut self, caller: &str) -> Result<(), String> {
        if caller != self.owner {
            return Err("Only owner can call this function".to_string());
        }

        self.paused = false;
        Ok(())
    }
}

// Export functions for blockchain integration
#[no_mangle]
pub extern "C" fn paprd_new(owner_ptr: *const u8, owner_len: usize) -> *mut PaprdContract {
    let owner = unsafe {
        std::str::from_utf8(std::slice::from_raw_parts(owner_ptr, owner_len))
            .unwrap_or("unknown")
            .to_string()
    };
    
    Box::into_raw(Box::new(PaprdContract::new(owner)))
}

#[no_mangle]
pub extern "C" fn paprd_get_balance(
    contract: *mut PaprdContract,
    address_ptr: *const u8,
    address_len: usize,
) -> u64 {
    let contract = unsafe { &*contract };
    let address = unsafe {
        std::str::from_utf8(std::slice::from_raw_parts(address_ptr, address_len))
            .unwrap_or("")
    };
    
    contract.get_balance(address)
}

#[no_mangle]
pub extern "C" fn paprd_transfer(
    contract: *mut PaprdContract,
    from_ptr: *const u8,
    from_len: usize,
    to_ptr: *const u8,
    to_len: usize,
    amount: u64,
) -> i32 {
    let contract = unsafe { &mut *contract };
    let from = unsafe {
        std::str::from_utf8(std::slice::from_raw_parts(from_ptr, from_len))
            .unwrap_or("")
    };
    let to = unsafe {
        std::str::from_utf8(std::slice::from_raw_parts(to_ptr, to_len))
            .unwrap_or("")
    };
    
    match contract.transfer(from, to, amount) {
        Ok(()) => 0,
        Err(_) => 1,
    }
} 