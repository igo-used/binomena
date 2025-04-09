// Token contract example in AssemblyScript
// This can be compiled to WASM using the AssemblyScript compiler

// Import necessary types
type u64 = number
type u8 = number

// State keys
const TOTAL_SUPPLY_KEY = "totalSupply"
const BALANCES_PREFIX = "balance:"
const ALLOWANCES_PREFIX = "allowance:"

// Token parameters
const TOKEN_NAME = "BinomenaToken"
const TOKEN_SYMBOL = "BMT"
const TOKEN_DECIMALS = 18

// Read total supply from state
export function totalSupply(): u64 {
  const totalSupply = Storage.get(TOTAL_SUPPLY_KEY)
  return totalSupply ? U64.parseInt(totalSupply) : 0
}

// Read balance of an account
export function balanceOf(account: string): u64 {
  const balance = Storage.get(BALANCES_PREFIX + account)
  return balance ? U64.parseInt(balance) : 0
}

// Transfer tokens to another account
export function transfer(to: string, amount: u64): boolean {
  const caller = Context.caller()

  // Check if sender has enough balance
  const senderBalance = balanceOf(caller)
  if (senderBalance < amount) {
    return false
  }

  // Update balances
  Storage.set(BALANCES_PREFIX + caller, (senderBalance - amount).toString())

  const recipientBalance = balanceOf(to)
  Storage.set(BALANCES_PREFIX + to, (recipientBalance + amount).toString())

  // Emit transfer event
  Context.emit("Transfer", [caller, to, amount.toString()])

  return true
}

// Approve another account to spend tokens
export function approve(spender: string, amount: u64): boolean {
  const owner = Context.caller()

  // Set allowance
  Storage.set(ALLOWANCES_PREFIX + owner + ":" + spender, amount.toString())

  // Emit approval event
  Context.emit("Approval", [owner, spender, amount.toString()])

  return true
}

// Get allowance
export function allowance(owner: string, spender: string): u64 {
  const allowance = Storage.get(ALLOWANCES_PREFIX + owner + ":" + spender)
  return allowance ? U64.parseInt(allowance) : 0
}

// Transfer tokens from one account to another (requires approval)
export function transferFrom(from: string, to: string, amount: u64): boolean {
  const caller = Context.caller()

  // Check allowance
  const currentAllowance = allowance(from, caller)
  if (currentAllowance < amount) {
    return false
  }

  // Check if sender has enough balance
  const senderBalance = balanceOf(from)
  if (senderBalance < amount) {
    return false
  }

  // Update allowance
  Storage.set(ALLOWANCES_PREFIX + from + ":" + caller, (currentAllowance - amount).toString())

  // Update balances
  Storage.set(BALANCES_PREFIX + from, (senderBalance - amount).toString())

  const recipientBalance = balanceOf(to)
  Storage.set(BALANCES_PREFIX + to, (recipientBalance + amount).toString())

  // Emit transfer event
  Context.emit("Transfer", [from, to, amount.toString()])

  return true
}

// Initialize the contract (called once when deployed)
export function initialize(initialSupply: u64): boolean {
  // Check if already initialized
  if (totalSupply() > 0) {
    return false
  }

  const caller = Context.caller()

  // Set total supply
  Storage.set(TOTAL_SUPPLY_KEY, initialSupply.toString())

  // Assign all tokens to contract creator
  Storage.set(BALANCES_PREFIX + caller, initialSupply.toString())

  // Emit transfer event (from zero address)
  Context.emit("Transfer", ["0x0000000000000000000000000000000000000000", caller, initialSupply.toString()])

  return true
}

// Get token name
export function name(): string {
  return TOKEN_NAME
}

// Get token symbol
export function symbol(): string {
  return TOKEN_SYMBOL
}

// Get token decimals
export function decimals(): u8 {
  return TOKEN_DECIMALS
}

// Storage namespace for state operations
namespace Storage {
  // Get a value from storage
  export function get(key: string): string | null {
    return null // Implemented by the runtime
  }

  // Set a value in storage
  export function set(key: string, value: string): void {
    // Implemented by the runtime
  }
}

// Context namespace for execution context
namespace Context {
  // Get the caller of the current function
  export function caller(): string {
    return "" // Implemented by the runtime
  }

  // Emit an event
  export function emit(name: string, args: string[]): void {
    // Implemented by the runtime
  }
}

// U64 namespace for parsing strings to u64
namespace U64 {
  export function parseInt(str: string): u64 {
    return parseInt(str)
  }
}
