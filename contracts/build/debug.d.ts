/** Exported memory */
export declare const memory: WebAssembly.Memory;
/**
 * assembly/index/createPaperDollar
 * @returns `assembly/stablecoin/PaperDollar`
 */
export declare function createPaperDollar(): __Internref4;
/**
 * assembly/index/getBalance
 * @param address `~lib/string/String`
 * @returns `u64`
 */
export declare function getBalance(address: string): bigint;
/**
 * assembly/index/getTotalSupply
 * @returns `u64`
 */
export declare function getTotalSupply(): bigint;
/**
 * assembly/index/transfer
 * @param to `~lib/string/String`
 * @param amount `u64`
 * @returns `bool`
 */
export declare function transfer(to: string, amount: bigint): boolean;
/**
 * assembly/index/mint
 * @param to `~lib/string/String`
 * @param amount `u64`
 * @returns `bool`
 */
export declare function mint(to: string, amount: bigint): boolean;
/**
 * assembly/index/burn
 * @param amount `u64`
 * @returns `bool`
 */
export declare function burn(amount: bigint): boolean;
/**
 * assembly/index/addMinter
 * @param address `~lib/string/String`
 * @returns `bool`
 */
export declare function addMinter(address: string): boolean;
/**
 * assembly/index/removeMinter
 * @param address `~lib/string/String`
 * @returns `bool`
 */
export declare function removeMinter(address: string): boolean;
/**
 * assembly/index/isMinter
 * @param address `~lib/string/String`
 * @returns `bool`
 */
export declare function isMinter(address: string): boolean;
/**
 * assembly/index/blacklist
 * @param address `~lib/string/String`
 * @returns `bool`
 */
export declare function blacklist(address: string): boolean;
/**
 * assembly/index/unblacklist
 * @param address `~lib/string/String`
 * @returns `bool`
 */
export declare function unblacklist(address: string): boolean;
/**
 * assembly/index/isBlacklisted
 * @param address `~lib/string/String`
 * @returns `bool`
 */
export declare function isBlacklisted(address: string): boolean;
/**
 * assembly/index/addCollateral
 * @param amount `u64`
 * @param collateralType `u32`
 * @returns `bool`
 */
export declare function addCollateral(amount: bigint, collateralType: number): boolean;
/**
 * assembly/index/removeCollateral
 * @param amount `u64`
 * @returns `bool`
 */
export declare function removeCollateral(amount: bigint): boolean;
/**
 * assembly/index/getCollateralBalance
 * @param address `~lib/string/String`
 * @param collateralType `u32`
 * @returns `u64`
 */
export declare function getCollateralBalance(address: string, collateralType: number): bigint;
/**
 * assembly/index/getCollateralType
 * @param address `~lib/string/String`
 * @returns `u32`
 */
export declare function getCollateralType(address: string): number;
/**
 * assembly/index/setCollateralRatio
 * @param ratio `u64`
 * @returns `bool`
 */
export declare function setCollateralRatio(ratio: bigint): boolean;
/**
 * assembly/index/getCollateralRatio
 * @returns `u64`
 */
export declare function getCollateralRatio(): bigint;
/**
 * assembly/index/setBinomTokenAddress
 * @param address `~lib/string/String`
 * @returns `bool`
 */
export declare function setBinomTokenAddress(address: string): boolean;
/**
 * assembly/index/getFiatReserve
 * @returns `u64`
 */
export declare function getFiatReserve(): bigint;
/**
 * assembly/index/pause
 * @returns `bool`
 */
export declare function pause(): boolean;
/**
 * assembly/index/unpause
 * @returns `bool`
 */
export declare function unpause(): boolean;
/**
 * assembly/index/transferOwnership
 * @param newOwner `~lib/string/String`
 * @returns `bool`
 */
export declare function transferOwnership(newOwner: string): boolean;
/**
 * assembly/index/add
 * @param a `i32`
 * @param b `i32`
 * @returns `i32`
 */
export declare function add(a: number, b: number): number;
/**
 * assembly/index/multiply
 * @param a `i32`
 * @param b `i32`
 * @returns `i32`
 */
export declare function multiply(a: number, b: number): number;
/**
 * assembly/index/store
 * @param value `i32`
 */
export declare function store(value: number): void;
/**
 * assembly/index/retrieve
 * @returns `i32`
 */
export declare function retrieve(): number;
/**
 * assembly/stablecoin/emitTransfer
 * @param from `~lib/string/String`
 * @param to `~lib/string/String`
 * @param amount `u64`
 */
export declare function emitTransfer(from: string, to: string, amount: bigint): void;
/**
 * assembly/stablecoin/emitMint
 * @param to `~lib/string/String`
 * @param amount `u64`
 * @param collateralType `i32`
 */
export declare function emitMint(to: string, amount: bigint, collateralType: number): void;
/**
 * assembly/stablecoin/emitBurn
 * @param from `~lib/string/String`
 * @param amount `u64`
 */
export declare function emitBurn(from: string, amount: bigint): void;
/**
 * assembly/stablecoin/emitCollateralAdded
 * @param from `~lib/string/String`
 * @param amount `u64`
 * @param collateralType `i32`
 */
export declare function emitCollateralAdded(from: string, amount: bigint, collateralType: number): void;
/**
 * assembly/stablecoin/emitCollateralRemoved
 * @param to `~lib/string/String`
 * @param amount `u64`
 * @param collateralType `i32`
 */
export declare function emitCollateralRemoved(to: string, amount: bigint, collateralType: number): void;
/** assembly/stablecoin/PaperDollar */
declare class __Internref4 extends Number {
  private __nominal4: symbol;
  private __nominal0: symbol;
}
