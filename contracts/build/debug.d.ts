/** Exported memory */
export declare const memory: WebAssembly.Memory;
/**
 * assembly/index/createPaperDollar
 * @returns `assembly/stablecoin/PaperDollar`
 */
export declare function createPaperDollar(): __Internref4;
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
