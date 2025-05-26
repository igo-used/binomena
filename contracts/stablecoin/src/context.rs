use std::sync::Mutex;

// Global context state
static CALLER: Mutex<Option<String>> = Mutex::new(None);

pub struct Context;

impl Context {
    /// Get the current caller address
    pub fn caller() -> String {
        let caller = CALLER.lock().unwrap();
        caller.as_ref().unwrap_or(&"default_caller".to_string()).clone()
    }

    /// Set the caller address (used by the VM)
    pub fn set_caller(address: String) {
        let mut caller = CALLER.lock().unwrap();
        *caller = Some(address);
    }

    /// Get block timestamp (placeholder for future blockchain integration)
    pub fn block_timestamp() -> u64 {
        use std::time::{SystemTime, UNIX_EPOCH};
        SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .unwrap()
            .as_secs()
    }

    /// Get block height (placeholder for future blockchain integration)
    pub fn block_height() -> u64 {
        // In a real implementation, this would come from the blockchain
        1
    }

    /// Get current contract address (placeholder)
    pub fn current_account_id() -> String {
        "paper_dollar_contract".to_string()
    }

    /// Check if the caller is a specific address
    pub fn is_caller(address: &str) -> bool {
        Self::caller() == address
    }

    /// Reset context (useful for testing)
    #[cfg(test)]
    pub fn reset() {
        let mut caller = CALLER.lock().unwrap();
        *caller = None;
    }
}

// Macro for setting up test context
#[macro_export]
macro_rules! with_caller {
    ($caller:expr, $block:expr) => {{
        crate::context::Context::set_caller($caller.to_string());
        let result = $block;
        #[cfg(test)]
        crate::context::Context::reset();
        result
    }};
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_context_operations() {
        Context::reset();
        
        // Test default caller
        let default_caller = Context::caller();
        assert_eq!(default_caller, "default_caller");

        // Test setting caller
        Context::set_caller("test_address".to_string());
        assert_eq!(Context::caller(), "test_address");
        
        // Test is_caller
        assert!(Context::is_caller("test_address"));
        assert!(!Context::is_caller("other_address"));

        // Test block operations
        let timestamp = Context::block_timestamp();
        assert!(timestamp > 0);
        
        let height = Context::block_height();
        assert_eq!(height, 1);

        let account_id = Context::current_account_id();
        assert_eq!(account_id, "paper_dollar_contract");
    }

    #[test]
    fn test_with_caller_macro() {
        let result = with_caller!("test_caller", {
            assert_eq!(Context::caller(), "test_caller");
            42
        });
        
        assert_eq!(result, 42);
        // After macro, context should be reset in test mode
        assert_eq!(Context::caller(), "default_caller");
    }
} 