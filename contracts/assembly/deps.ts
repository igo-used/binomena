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
    export let caller: string = ""; // Will be set by VM at runtime
}

// Contract decorator
export function contract(_constructor: Function): void {
    // Implementation provided by VM
}

// View decorator for read-only functions
export function view(_target: any, _propertyKey: string): void {
    // Implementation provided by VM
}

// Mutate decorator for state-changing functions
export function mutate(_target: any, _propertyKey: string): void {
    // Implementation provided by VM
}

// Assert function for validation
export function assert(condition: boolean, message: string = "Assertion failed"): void {
    if (!condition) {
        throw new Error(message);
    }
} 