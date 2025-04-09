// AssemblyScript smart contract

// Export a function that adds two numbers
export function add(a: i32, b: i32): i32 {
    return a + b;
  }
  
  // Export a function that multiplies two numbers
  export function multiply(a: i32, b: i32): i32 {
    return a * b;
  }
  
  // Export a function that stores and retrieves a value
  let storedValue: i32 = 0;
  
  export function store(value: i32): void {
    storedValue = value;
  }
  
  export function retrieve(): i32 {
    return storedValue;
  }