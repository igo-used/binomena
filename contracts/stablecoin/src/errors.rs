use thiserror::Error;

#[derive(Error, Debug, Clone)]
pub enum ContractError {
    #[error("Only owner can call this function")]
    OnlyOwner,
    
    #[error("Contract is paused")]
    ContractPaused,
    
    #[error("Address is blacklisted")]
    AddressBlacklisted,
    
    #[error("Only minters can mint")]
    OnlyMinters,
    
    #[error("Insufficient balance")]
    InsufficientBalance,
    
    #[error("Insufficient collateral")]
    InsufficientCollateral,
    
    #[error("Collateral ratio must be at least 100%")]
    InvalidCollateralRatio,
    
    #[error("Collateral ratio would be too low")]
    CollateralRatioTooLow,
    
    #[error("Invalid amount: {0}")]
    InvalidAmount(String),
    
    #[error("Invalid address: {0}")]
    InvalidAddress(String),
    
    #[error("Storage error: {0}")]
    StorageError(String),
    
    #[error("Arithmetic overflow")]
    ArithmeticOverflow,
    
    #[error("Arithmetic underflow")]
    ArithmeticUnderflow,
    
    #[error("Division by zero")]
    DivisionByZero,
    
    #[error("Invalid collateral type")]
    InvalidCollateralType,
    
    #[error("Transfer to self not allowed")]
    TransferToSelf,
    
    #[error("Zero amount not allowed")]
    ZeroAmount,
}

pub type ContractResult<T> = Result<T, ContractError>;

// Helper functions for safe arithmetic operations
pub fn safe_add(a: u64, b: u64) -> ContractResult<u64> {
    a.checked_add(b).ok_or(ContractError::ArithmeticOverflow)
}

pub fn safe_sub(a: u64, b: u64) -> ContractResult<u64> {
    a.checked_sub(b).ok_or(ContractError::ArithmeticUnderflow)
}

pub fn safe_mul(a: u64, b: u64) -> ContractResult<u64> {
    a.checked_mul(b).ok_or(ContractError::ArithmeticOverflow)
}

pub fn safe_div(a: u64, b: u64) -> ContractResult<u64> {
    if b == 0 {
        return Err(ContractError::DivisionByZero);
    }
    Ok(a / b)
}

// Macro for assertions with custom error messages
#[macro_export]
macro_rules! require {
    ($condition:expr, $error:expr) => {
        if !$condition {
            return Err($error);
        }
    };
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_safe_arithmetic() {
        assert_eq!(safe_add(5, 3).unwrap(), 8);
        assert!(safe_add(u64::MAX, 1).is_err());

        assert_eq!(safe_sub(5, 3).unwrap(), 2);
        assert!(safe_sub(3, 5).is_err());

        assert_eq!(safe_mul(5, 3).unwrap(), 15);
        assert!(safe_mul(u64::MAX, 2).is_err());

        assert_eq!(safe_div(6, 3).unwrap(), 2);
        assert!(safe_div(6, 0).is_err());
    }

    #[test]
    fn test_error_display() {
        let error = ContractError::OnlyOwner;
        assert_eq!(error.to_string(), "Only owner can call this function");

        let error = ContractError::InvalidAmount("negative value".to_string());
        assert_eq!(error.to_string(), "Invalid amount: negative value");
    }

    #[test]
    fn test_require_macro() {
        fn test_function(value: u64) -> ContractResult<u64> {
            require!(value > 0, ContractError::ZeroAmount);
            Ok(value * 2)
        }

        assert_eq!(test_function(5).unwrap(), 10);
        assert!(matches!(test_function(0).unwrap_err(), ContractError::ZeroAmount));
    }
} 