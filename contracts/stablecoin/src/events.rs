use serde::{Serialize, Deserialize};
use borsh::{BorshSerialize, BorshDeserialize};

#[derive(Clone, Debug, PartialEq, Eq, BorshSerialize, BorshDeserialize)]
pub enum CollateralType {
    Fiat = 0,
    Binom = 1,
}

impl From<u8> for CollateralType {
    fn from(value: u8) -> Self {
        match value {
            0 => CollateralType::Fiat,
            1 => CollateralType::Binom,
            _ => CollateralType::Fiat, // Default to Fiat for invalid values
        }
    }
}

impl From<CollateralType> for u8 {
    fn from(collateral_type: CollateralType) -> Self {
        match collateral_type {
            CollateralType::Fiat => 0,
            CollateralType::Binom => 1,
        }
    }
}

impl ToString for CollateralType {
    fn to_string(&self) -> String {
        match self {
            CollateralType::Fiat => "0".to_string(),
            CollateralType::Binom => "1".to_string(),
        }
    }
}

// Event structures
#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct TransferEvent {
    pub from: String,
    pub to: String,
    pub amount: u64,
    pub timestamp: u64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct MintEvent {
    pub to: String,
    pub amount: u64,
    pub collateral_type: CollateralType,
    pub timestamp: u64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct BurnEvent {
    pub from: String,
    pub amount: u64,
    pub timestamp: u64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct CollateralAddedEvent {
    pub from: String,
    pub amount: u64,
    pub collateral_type: CollateralType,
    pub timestamp: u64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct CollateralRemovedEvent {
    pub to: String,
    pub amount: u64,
    pub collateral_type: CollateralType,
    pub timestamp: u64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct OwnershipTransferredEvent {
    pub previous_owner: String,
    pub new_owner: String,
    pub timestamp: u64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct PausedEvent {
    pub timestamp: u64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct UnpausedEvent {
    pub timestamp: u64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct BlacklistAddedEvent {
    pub address: String,
    pub timestamp: u64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct BlacklistRemovedEvent {
    pub address: String,
    pub timestamp: u64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct MinterAddedEvent {
    pub address: String,
    pub timestamp: u64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct MinterRemovedEvent {
    pub address: String,
    pub timestamp: u64,
}

// Event emission functions
pub fn emit_transfer(from: &str, to: &str, amount: u64) {
    let event = TransferEvent {
        from: from.to_string(),
        to: to.to_string(),
        amount,
        timestamp: crate::context::Context::block_timestamp(),
    };
    
    // In a real blockchain implementation, this would emit the event to the blockchain
    // For now, we'll log it
    log_event("Transfer", &event);
}

pub fn emit_mint(to: &str, amount: u64, collateral_type: CollateralType) {
    let event = MintEvent {
        to: to.to_string(),
        amount,
        collateral_type,
        timestamp: crate::context::Context::block_timestamp(),
    };
    
    log_event("Mint", &event);
}

pub fn emit_burn(from: &str, amount: u64) {
    let event = BurnEvent {
        from: from.to_string(),
        amount,
        timestamp: crate::context::Context::block_timestamp(),
    };
    
    log_event("Burn", &event);
}

pub fn emit_collateral_added(from: &str, amount: u64, collateral_type: CollateralType) {
    let event = CollateralAddedEvent {
        from: from.to_string(),
        amount,
        collateral_type,
        timestamp: crate::context::Context::block_timestamp(),
    };
    
    log_event("CollateralAdded", &event);
}

pub fn emit_collateral_removed(to: &str, amount: u64, collateral_type: CollateralType) {
    let event = CollateralRemovedEvent {
        to: to.to_string(),
        amount,
        collateral_type,
        timestamp: crate::context::Context::block_timestamp(),
    };
    
    log_event("CollateralRemoved", &event);
}

pub fn emit_ownership_transferred(previous_owner: &str, new_owner: &str) {
    let event = OwnershipTransferredEvent {
        previous_owner: previous_owner.to_string(),
        new_owner: new_owner.to_string(),
        timestamp: crate::context::Context::block_timestamp(),
    };
    
    log_event("OwnershipTransferred", &event);
}

pub fn emit_paused() {
    let event = PausedEvent {
        timestamp: crate::context::Context::block_timestamp(),
    };
    
    log_event("Paused", &event);
}

pub fn emit_unpaused() {
    let event = UnpausedEvent {
        timestamp: crate::context::Context::block_timestamp(),
    };
    
    log_event("Unpaused", &event);
}

pub fn emit_blacklist_added(address: &str) {
    let event = BlacklistAddedEvent {
        address: address.to_string(),
        timestamp: crate::context::Context::block_timestamp(),
    };
    
    log_event("BlacklistAdded", &event);
}

pub fn emit_blacklist_removed(address: &str) {
    let event = BlacklistRemovedEvent {
        address: address.to_string(),
        timestamp: crate::context::Context::block_timestamp(),
    };
    
    log_event("BlacklistRemoved", &event);
}

pub fn emit_minter_added(address: &str) {
    let event = MinterAddedEvent {
        address: address.to_string(),
        timestamp: crate::context::Context::block_timestamp(),
    };
    
    log_event("MinterAdded", &event);
}

pub fn emit_minter_removed(address: &str) {
    let event = MinterRemovedEvent {
        address: address.to_string(),
        timestamp: crate::context::Context::block_timestamp(),
    };
    
    log_event("MinterRemoved", &event);
}

// Helper function to log events (in a real implementation, this would emit to blockchain)
fn log_event<T: Serialize>(event_name: &str, event_data: &T) {
    if let Ok(json) = serde_json::to_string(event_data) {
        #[cfg(target_arch = "wasm32")]
        {
            crate::console_log!("Event {}: {}", event_name, json);
        }
        
        #[cfg(not(target_arch = "wasm32"))]
        {
            println!("Event {}: {}", event_name, json);
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_collateral_type_conversion() {
        assert_eq!(CollateralType::from(0), CollateralType::Fiat);
        assert_eq!(CollateralType::from(1), CollateralType::Binom);
        assert_eq!(CollateralType::from(99), CollateralType::Fiat); // Invalid defaults to Fiat

        assert_eq!(u8::from(CollateralType::Fiat), 0);
        assert_eq!(u8::from(CollateralType::Binom), 1);

        assert_eq!(CollateralType::Fiat.to_string(), "0");
        assert_eq!(CollateralType::Binom.to_string(), "1");
    }

    #[test]
    fn test_event_emission() {
        // These tests just verify that events can be created and don't panic
        emit_transfer("from_addr", "to_addr", 100);
        emit_mint("to_addr", 50, CollateralType::Fiat);
        emit_burn("from_addr", 25);
        emit_collateral_added("from_addr", 200, CollateralType::Binom);
        emit_collateral_removed("to_addr", 100, CollateralType::Fiat);
        emit_ownership_transferred("old_owner", "new_owner");
        emit_paused();
        emit_unpaused();
        emit_blacklist_added("bad_addr");
        emit_blacklist_removed("good_addr");
        emit_minter_added("minter_addr");
        emit_minter_removed("ex_minter_addr");
    }
} 