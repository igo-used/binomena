// Storage interface for the contract
export namespace storage {
    export function get<T>(key: string, defaultValue: T): T {
        return defaultValue; // Implementation provided by VM
    }

    export function set<T>(key: string, value: T): void {
        // Implementation provided by VM
    }
}

// Context information for the current execution
export namespace Context {
    let _caller: string = ""; // Internal storage for caller
    
    // Get the current caller address
    export function getCaller(): string {
        if (_caller === "") {
            // Try to read caller from memory location set by VM
            _caller = readCallerFromMemory();
        }
        return _caller;
    }
    
    // Set the caller address (called by VM)
    export function setCaller(value: string): void {
        _caller = value;
    }
    
    // For compatibility, export as property
    export let caller: string = "";
    
    // Update caller property when needed
    export function updateCaller(): void {
        caller = getCaller();
    }
}

// Function to set caller context (called by VM)
export function setCaller(callerPointer: i32): void {
    const callerAddress = readStringFromPointer(callerPointer);
    Context.setCaller(callerAddress);
    Context.caller = callerAddress; // Update the property too
    Context.updateCaller();
}

// Read a string from WASM memory given a pointer
function readStringFromPointer(pointer: i32): string {
    if (pointer == 0) return "";
    
    // For simplified string handling, convert hash back to address format
    // This is a simplified approach for the string parameter issue
    const hashValue = pointer;
    
    // For wallet addresses, we can reconstruct a valid format
    if (hashValue > 0) {
        // Generate a mock address based on the hash for testing
        // In a real implementation, you might have a lookup table
        return "AdNe" + hashValue.toString(16).padStart(56, "0");
    }
    
    return "";
}

// Read caller from fixed memory location (fallback)
function readCallerFromMemory(): string {
    // For the simplified approach, try to read from known memory locations
    // where the VM stores string data
    let result = "";
    
    // Try to read from memory locations 200-300 where VM stores strings
    for (let offset = 200; offset < 300; offset += 20) {
        let found = "";
        for (let i = 0; i < 16; i++) {
            const byteVal = readByteFromMemory(offset + i);
            if (byteVal > 0 && byteVal < 128) { // Valid ASCII
                found += String.fromCharCode(byteVal);
            } else {
                break;
            }
        }
        if (found.length >= 4 && found.startsWith("AdNe")) {
            result = found;
            break;
        }
    }
    
    return result;
}

// Read a byte from memory at specified offset
function readByteFromMemory(offset: i32): i32 {
    // This is a placeholder - in real implementation you'd use load<u8>
    // For now, return 0 to avoid memory access issues
    return 0;
}

// Read int32 from memory at specified offset
function readInt32FromMemory(offset: i32): i32 {
    // Read 4 bytes from memory and construct int32
    // This is a placeholder - in real implementation you'd use load<i32>
    return 0; // Simplified for now
}

// Contract decorator
export function contract(target: Function): void {
    // Implementation provided by VM
}

// View decorator for read-only functions
export function view(target: any, propertyKey: string): void {
    // Implementation provided by VM
}

// Mutate decorator for state-changing functions
export function mutate(target: any, propertyKey: string): void {
    // Implementation provided by VM
}

// Assert function for validation
export function assert(condition: boolean, message: string = "Assertion failed"): void {
    if (!condition) {
        throw new Error(message);
    }
} 