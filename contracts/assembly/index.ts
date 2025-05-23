// AssemblyScript smart contract

import { PaperDollar } from "./stablecoin";

export { PaperDollar };

// Global contract instance
let contractInstance: PaperDollar = new PaperDollar();

// Contract factory function
export function createPaperDollar(): PaperDollar {
  return contractInstance;
}

// Helper function to get contract instance
function getContract(): PaperDollar {
  return contractInstance;
}

// Export all stablecoin functions

// Balance and supply functions
export function getBalance(address: string): u64 {
  return getContract().getBalance(address);
}

export function getTotalSupply(): u64 {
  return getContract().getTotalSupply();
}

// Transfer functions
export function transfer(to: string, amount: u64): boolean {
  return getContract().transfer(to, amount);
}

// Minting and burning
export function mint(to: string, amount: u64): boolean {
  return getContract().mint(to, amount);
}

export function burn(amount: u64): boolean {
  return getContract().burn(amount);
}

// Minter management
export function addMinter(address: string): boolean {
  return getContract().addMinter(address);
}

export function removeMinter(address: string): boolean {
  return getContract().removeMinter(address);
}

export function isMinter(address: string): boolean {
  return getContract().isMinter(address);
}

// Blacklist management
export function blacklist(address: string): boolean {
  return getContract().blacklist(address);
}

export function unblacklist(address: string): boolean {
  return getContract().unblacklist(address);
}

export function isBlacklisted(address: string): boolean {
  return getContract().isBlacklisted(address);
}

// Collateral management
export function addCollateral(amount: u64, collateralType: u32): boolean {
  return getContract().addCollateral(amount, collateralType);
}

export function removeCollateral(amount: u64): boolean {
  return getContract().removeCollateral(amount);
}

export function getCollateralBalance(address: string, collateralType: u32): u64 {
  return getContract().getCollateralBalance(address, collateralType);
}

export function getCollateralType(address: string): u32 {
  return getContract().getCollateralType(address);
}

// Configuration functions
export function setCollateralRatio(ratio: u64): boolean {
  return getContract().setCollateralRatio(ratio);
}

export function getCollateralRatio(): u64 {
  return getContract().getCollateralRatio();
}

export function setBinomTokenAddress(address: string): boolean {
  return getContract().setBinomTokenAddress(address);
}

export function getFiatReserve(): u64 {
  return getContract().getFiatReserve();
}

// Pause/unpause functions
export function pause(): boolean {
  return getContract().pause();
}

export function unpause(): boolean {
  return getContract().unpause();
}

// Ownership functions
export function transferOwnership(newOwner: string): boolean {
  return getContract().transferOwnership(newOwner);
}

// Test functions (keep for compatibility)
export function add(a: i32, b: i32): i32 {
    return a + b;
}

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