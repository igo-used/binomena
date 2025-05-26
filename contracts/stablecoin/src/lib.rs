use wasm_bindgen::prelude::*;

mod storage;
mod context;
mod stablecoin;
mod events;
mod errors;

pub use storage::Storage;
pub use context::Context;
pub use stablecoin::PaperDollar;
pub use events::*;
pub use errors::*;

// Enable logging in wasm
#[cfg(target_arch = "wasm32")]
#[wasm_bindgen]
extern "C" {
    #[wasm_bindgen(js_namespace = console)]
    fn log(s: &str);
}

#[macro_export]
macro_rules! console_log {
    ($($t:tt)*) => {
        #[cfg(target_arch = "wasm32")]
        {
            web_sys::console::log_1(&format_args!($($t)*).to_string().into());
        }
        #[cfg(not(target_arch = "wasm32"))]
        {
            println!("{}", format_args!($($t)*))
        }
    }
}

// Export the contract interface
#[wasm_bindgen]
pub struct ContractInstance {
    contract: PaperDollar,
}

#[wasm_bindgen]
impl ContractInstance {
    #[wasm_bindgen(constructor)]
    pub fn new() -> ContractInstance {
        ContractInstance {
            contract: PaperDollar::new(),
        }
    }

    /// Initialize contract with founder address as owner
    #[wasm_bindgen]
    pub fn new_with_founder() -> ContractInstance {
        // Set the founder as the caller
        Context::set_caller("AdNe6c3ce54e4371d056c7c566675ba16909eb2e9534".to_string());
        
        ContractInstance {
            contract: PaperDollar::new(),
        }
    }

    /// Get the current owner of the contract
    #[wasm_bindgen]
    pub fn get_owner(&self) -> String {
        self.contract.get_owner()
    }

    /// Check if the contract is paused
    #[wasm_bindgen]
    pub fn is_paused(&self) -> bool {
        self.contract.is_paused()
    }

    /// Get token name
    #[wasm_bindgen]
    pub fn get_name(&self) -> String {
        self.contract.get_name()
    }

    /// Get token symbol
    #[wasm_bindgen]
    pub fn get_symbol(&self) -> String {
        self.contract.get_symbol()
    }

    /// Get token decimals
    #[wasm_bindgen]
    pub fn get_decimals(&self) -> u8 {
        self.contract.get_decimals()
    }

    // View functions
    #[wasm_bindgen]
    pub fn get_balance(&self, address: &str) -> u64 {
        self.contract.get_balance(address)
    }

    #[wasm_bindgen]
    pub fn get_total_supply(&self) -> u64 {
        self.contract.get_total_supply()
    }

    #[wasm_bindgen]
    pub fn get_collateral_balance(&self, address: &str, collateral_type: u8) -> u64 {
        self.contract.get_collateral_balance(address, collateral_type.into())
    }

    #[wasm_bindgen]
    pub fn get_collateral_ratio(&self) -> u64 {
        self.contract.get_collateral_ratio()
    }

    #[wasm_bindgen]
    pub fn get_fiat_reserve(&self) -> u64 {
        self.contract.get_fiat_reserve()
    }

    #[wasm_bindgen]
    pub fn is_blacklisted(&self, address: &str) -> bool {
        self.contract.is_blacklisted(address)
    }

    #[wasm_bindgen]
    pub fn is_minter(&self, address: &str) -> bool {
        self.contract.is_minter(address)
    }

    // Mutating functions
    #[wasm_bindgen]
    pub fn transfer(&mut self, to: &str, amount: u64) -> Result<bool, JsValue> {
        self.contract.transfer(to, amount)
            .map_err(|e| JsValue::from_str(&e.to_string()))
    }

    #[wasm_bindgen]
    pub fn mint(&mut self, to: &str, amount: u64) -> Result<bool, JsValue> {
        self.contract.mint(to, amount)
            .map_err(|e| JsValue::from_str(&e.to_string()))
    }

    #[wasm_bindgen]
    pub fn burn(&mut self, amount: u64) -> Result<bool, JsValue> {
        self.contract.burn(amount)
            .map_err(|e| JsValue::from_str(&e.to_string()))
    }

    #[wasm_bindgen]
    pub fn add_collateral(&mut self, amount: u64, collateral_type: u8) -> Result<bool, JsValue> {
        self.contract.add_collateral(amount, collateral_type.into())
            .map_err(|e| JsValue::from_str(&e.to_string()))
    }

    #[wasm_bindgen]
    pub fn remove_collateral(&mut self, amount: u64) -> Result<bool, JsValue> {
        self.contract.remove_collateral(amount)
            .map_err(|e| JsValue::from_str(&e.to_string()))
    }

    // Owner functions
    #[wasm_bindgen]
    pub fn add_minter(&mut self, address: &str) -> Result<bool, JsValue> {
        self.contract.add_minter(address)
            .map_err(|e| JsValue::from_str(&e.to_string()))
    }

    #[wasm_bindgen]
    pub fn remove_minter(&mut self, address: &str) -> Result<bool, JsValue> {
        self.contract.remove_minter(address)
            .map_err(|e| JsValue::from_str(&e.to_string()))
    }

    #[wasm_bindgen]
    pub fn blacklist(&mut self, address: &str) -> Result<bool, JsValue> {
        self.contract.blacklist(address)
            .map_err(|e| JsValue::from_str(&e.to_string()))
    }

    #[wasm_bindgen]
    pub fn unblacklist(&mut self, address: &str) -> Result<bool, JsValue> {
        self.contract.unblacklist(address)
            .map_err(|e| JsValue::from_str(&e.to_string()))
    }

    #[wasm_bindgen]
    pub fn pause(&mut self) -> Result<bool, JsValue> {
        self.contract.pause()
            .map_err(|e| JsValue::from_str(&e.to_string()))
    }

    #[wasm_bindgen]
    pub fn unpause(&mut self) -> Result<bool, JsValue> {
        self.contract.unpause()
            .map_err(|e| JsValue::from_str(&e.to_string()))
    }

    #[wasm_bindgen]
    pub fn set_collateral_ratio(&mut self, ratio: u64) -> Result<bool, JsValue> {
        self.contract.set_collateral_ratio(ratio)
            .map_err(|e| JsValue::from_str(&e.to_string()))
    }

    #[wasm_bindgen]
    pub fn transfer_ownership(&mut self, new_owner: &str) -> Result<bool, JsValue> {
        self.contract.transfer_ownership(new_owner)
            .map_err(|e| JsValue::from_str(&e.to_string()))
    }
} 