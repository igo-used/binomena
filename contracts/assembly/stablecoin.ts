import { storage, Context, contract, view, mutate, assert } from "./contract";

// Contract storage keys
const OWNER: string = "owner";
const TOTAL_SUPPLY: string = "totalSupply";
const BALANCES_PREFIX: string = "balance_";
const PAUSED: string = "paused";
const BLACKLIST_PREFIX: string = "blacklist_";
const MINTER_PREFIX: string = "minter_";
const COLLATERAL_PREFIX: string = "collateral_";
const COLLATERAL_TYPE_PREFIX: string = "collateral_type_";
const BINOM_TOKEN_ADDRESS: string = "binom_token_address";
const COLLATERAL_RATIO: string = "collateral_ratio";
const FIAT_RESERVE: string = "fiat_reserve";

// Collateral types
enum CollateralType {
    FIAT,
    BINOM
}

// Events
export function emitTransfer(from: string, to: string, amount: u64): void {
    // Event emission will be handled by the VM
}

export function emitMint(to: string, amount: u64, collateralType: CollateralType): void {
    // Event emission will be handled by the VM
}

export function emitBurn(from: string, amount: u64): void {
    // Event emission will be handled by the VM
}

export function emitCollateralAdded(from: string, amount: u64, collateralType: CollateralType): void {
    // Event emission will be handled by the VM
}

export function emitCollateralRemoved(to: string, amount: u64, collateralType: CollateralType): void {
    // Event emission will be handled by the VM
}

@contract
export class PaperDollar {
    // Initialize the contract
    constructor() {
        if (!storage.get<string>(OWNER, "")) {
            storage.set<string>(OWNER, Context.caller);
            storage.set<u64>(TOTAL_SUPPLY, 0);
            storage.set<boolean>(PAUSED, false);
            storage.set<u64>(COLLATERAL_RATIO, 150); // 150% collateralization ratio
            storage.set<u64>(FIAT_RESERVE, 0);
        }
    }

    // Only owner modifier
    private checkOwner(): void {
        assert(storage.get<string>(OWNER, "") == Context.caller, "Only owner can call this function");
    }

    // Not paused modifier
    private checkNotPaused(): void {
        assert(!storage.get<boolean>(PAUSED, false), "Contract is paused");
    }

    // Not blacklisted modifier
    private checkNotBlacklisted(address: string): void {
        assert(!this.isBlacklisted(address), "Address is blacklisted");
    }

    // Get balance of an address
    @view
    getBalance(address: string): u64 {
        return storage.get<u64>(BALANCES_PREFIX + address, 0);
    }

    // Get total supply
    @view
    getTotalSupply(): u64 {
        return storage.get<u64>(TOTAL_SUPPLY, 0);
    }

    // Get collateral balance
    @view
    getCollateralBalance(address: string, collateralType: CollateralType): u64 {
        return storage.get<u64>(COLLATERAL_PREFIX + collateralType.toString() + "_" + address, 0);
    }

    // Get collateral type
    @view
    getCollateralType(address: string): CollateralType {
        return storage.get<CollateralType>(COLLATERAL_TYPE_PREFIX + address, CollateralType.FIAT);
    }

    // Get collateral ratio
    @view
    getCollateralRatio(): u64 {
        return storage.get<u64>(COLLATERAL_RATIO, 150);
    }

    // Get fiat reserve
    @view
    getFiatReserve(): u64 {
        return storage.get<u64>(FIAT_RESERVE, 0);
    }

    // Check if address is blacklisted
    @view
    isBlacklisted(address: string): boolean {
        return storage.get<boolean>(BLACKLIST_PREFIX + address, false);
    }

    // Check if address is a minter
    @view
    isMinter(address: string): boolean {
        return storage.get<boolean>(MINTER_PREFIX + address, false);
    }

    // Set BINOM token address (only owner)
    @mutate
    setBinomTokenAddress(address: string): boolean {
        this.checkOwner();
        storage.set<string>(BINOM_TOKEN_ADDRESS, address);
        return true;
    }

    // Set collateral ratio (only owner)
    @mutate
    setCollateralRatio(ratio: u64): boolean {
        this.checkOwner();
        assert(ratio >= 100, "Collateral ratio must be at least 100%");
        storage.set<u64>(COLLATERAL_RATIO, ratio);
        return true;
    }

    // Add collateral
    @mutate
    addCollateral(amount: u64, collateralType: CollateralType): boolean {
        const from = Context.caller;
        this.checkNotPaused();
        this.checkNotBlacklisted(from);

        if (collateralType == CollateralType.FIAT) {
            storage.set<u64>(FIAT_RESERVE, this.getFiatReserve() + amount);
        }

        const currentCollateral = this.getCollateralBalance(from, collateralType);
        storage.set<u64>(COLLATERAL_PREFIX + collateralType.toString() + "_" + from, currentCollateral + amount);
        storage.set<CollateralType>(COLLATERAL_TYPE_PREFIX + from, collateralType);

        emitCollateralAdded(from, amount, collateralType);
        return true;
    }

    // Remove collateral
    @mutate
    removeCollateral(amount: u64): boolean {
        const from = Context.caller;
        this.checkNotPaused();
        this.checkNotBlacklisted(from);

        const collateralType = this.getCollateralType(from);
        const currentCollateral = this.getCollateralBalance(from, collateralType);
        assert(currentCollateral >= amount, "Insufficient collateral");

        // Check if removing collateral would break the collateral ratio
        const userBalance = this.getBalance(from);
        const remainingCollateral = currentCollateral - amount;
        assert(remainingCollateral * 100 >= userBalance * this.getCollateralRatio(), "Collateral ratio would be too low");

        if (collateralType == CollateralType.FIAT) {
            storage.set<u64>(FIAT_RESERVE, this.getFiatReserve() - amount);
        }

        storage.set<u64>(COLLATERAL_PREFIX + collateralType.toString() + "_" + from, remainingCollateral);

        emitCollateralRemoved(from, amount, collateralType);
        return true;
    }

    // Transfer tokens
    @mutate
    transfer(to: string, amount: u64): boolean {
        const from = Context.caller;
        this.checkNotPaused();
        this.checkNotBlacklisted(from);
        this.checkNotBlacklisted(to);

        const fromBalance = this.getBalance(from);
        assert(fromBalance >= amount, "Insufficient balance");

        storage.set<u64>(BALANCES_PREFIX + from, fromBalance - amount);
        storage.set<u64>(BALANCES_PREFIX + to, this.getBalance(to) + amount);

        emitTransfer(from, to, amount);
        return true;
    }

    // Mint new tokens (only minters)
    @mutate
    mint(to: string, amount: u64): boolean {
        assert(this.isMinter(Context.caller), "Only minters can mint");
        this.checkNotPaused();
        this.checkNotBlacklisted(to);

        const collateralType = this.getCollateralType(to);
        const collateral = this.getCollateralBalance(to, collateralType);
        const currentBalance = this.getBalance(to);
        const newBalance = currentBalance + amount;

        // Check collateral ratio
        assert(collateral * 100 >= newBalance * this.getCollateralRatio(), "Insufficient collateral");

        const newTotalSupply = this.getTotalSupply() + amount;
        storage.set<u64>(TOTAL_SUPPLY, newTotalSupply);
        storage.set<u64>(BALANCES_PREFIX + to, newBalance);

        emitMint(to, amount, collateralType);
        return true;
    }

    // Burn tokens
    @mutate
    burn(amount: u64): boolean {
        const from = Context.caller;
        this.checkNotPaused();
        this.checkNotBlacklisted(from);

        const fromBalance = this.getBalance(from);
        assert(fromBalance >= amount, "Insufficient balance");

        storage.set<u64>(BALANCES_PREFIX + from, fromBalance - amount);
        storage.set<u64>(TOTAL_SUPPLY, this.getTotalSupply() - amount);

        emitBurn(from, amount);
        return true;
    }

    // Add a minter (only owner)
    @mutate
    addMinter(address: string): boolean {
        this.checkOwner();
        storage.set<boolean>(MINTER_PREFIX + address, true);
        return true;
    }

    // Remove a minter (only owner)
    @mutate
    removeMinter(address: string): boolean {
        this.checkOwner();
        storage.set<boolean>(MINTER_PREFIX + address, false);
        return true;
    }

    // Blacklist an address (only owner)
    @mutate
    blacklist(address: string): boolean {
        this.checkOwner();
        storage.set<boolean>(BLACKLIST_PREFIX + address, true);
        return true;
    }

    // Remove address from blacklist (only owner)
    @mutate
    unblacklist(address: string): boolean {
        this.checkOwner();
        storage.set<boolean>(BLACKLIST_PREFIX + address, false);
        return true;
    }

    // Pause the contract (only owner)
    @mutate
    pause(): boolean {
        this.checkOwner();
        storage.set<boolean>(PAUSED, true);
        return true;
    }

    // Unpause the contract (only owner)
    @mutate
    unpause(): boolean {
        this.checkOwner();
        storage.set<boolean>(PAUSED, false);
        return true;
    }

    // Transfer ownership (only owner)
    @mutate
    transferOwnership(newOwner: string): boolean {
        this.checkOwner();
        storage.set<string>(OWNER, newOwner);
        return true;
    }
} 