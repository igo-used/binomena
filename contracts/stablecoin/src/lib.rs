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
#[wasm_bindgen]
extern "C" {
    #[wasm_bindgen(js_namespace = console)]
    fn log(s: &str);
}

#[macro_export]
macro_rules! console_log {
    ($($t:tt)*) => (log(&format_args!($($t)*).to_string()))
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