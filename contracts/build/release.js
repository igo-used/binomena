async function instantiate(module, imports = {}) {
  const adaptedImports = {
    env: Object.assign(Object.create(globalThis), imports.env || {}, {
      abort(message, fileName, lineNumber, columnNumber) {
        // ~lib/builtins/abort(~lib/string/String | null?, ~lib/string/String | null?, u32?, u32?) => void
        message = __liftString(message >>> 0);
        fileName = __liftString(fileName >>> 0);
        lineNumber = lineNumber >>> 0;
        columnNumber = columnNumber >>> 0;
        (() => {
          // @external.js
          throw Error(`${message} in ${fileName}:${lineNumber}:${columnNumber}`);
        })();
      },
    }),
  };
  const { exports } = await WebAssembly.instantiate(module, adaptedImports);
  const memory = exports.memory || imports.env.memory;
  const adaptedExports = Object.setPrototypeOf({
    createPaperDollar() {
      // assembly/index/createPaperDollar() => assembly/stablecoin/PaperDollar
      return __liftInternref(exports.createPaperDollar() >>> 0);
    },
    getBalance(address) {
      // assembly/index/getBalance(~lib/string/String) => u64
      address = __lowerString(address) || __notnull();
      return BigInt.asUintN(64, exports.getBalance(address));
    },
    getTotalSupply() {
      // assembly/index/getTotalSupply() => u64
      return BigInt.asUintN(64, exports.getTotalSupply());
    },
    transfer(to, amount) {
      // assembly/index/transfer(~lib/string/String, u64) => bool
      to = __lowerString(to) || __notnull();
      amount = amount || 0n;
      return exports.transfer(to, amount) != 0;
    },
    mint(to, amount) {
      // assembly/index/mint(~lib/string/String, u64) => bool
      to = __lowerString(to) || __notnull();
      amount = amount || 0n;
      return exports.mint(to, amount) != 0;
    },
    burn(amount) {
      // assembly/index/burn(u64) => bool
      amount = amount || 0n;
      return exports.burn(amount) != 0;
    },
    addMinter(address) {
      // assembly/index/addMinter(~lib/string/String) => bool
      address = __lowerString(address) || __notnull();
      return exports.addMinter(address) != 0;
    },
    removeMinter(address) {
      // assembly/index/removeMinter(~lib/string/String) => bool
      address = __lowerString(address) || __notnull();
      return exports.removeMinter(address) != 0;
    },
    isMinter(address) {
      // assembly/index/isMinter(~lib/string/String) => bool
      address = __lowerString(address) || __notnull();
      return exports.isMinter(address) != 0;
    },
    blacklist(address) {
      // assembly/index/blacklist(~lib/string/String) => bool
      address = __lowerString(address) || __notnull();
      return exports.blacklist(address) != 0;
    },
    unblacklist(address) {
      // assembly/index/unblacklist(~lib/string/String) => bool
      address = __lowerString(address) || __notnull();
      return exports.unblacklist(address) != 0;
    },
    isBlacklisted(address) {
      // assembly/index/isBlacklisted(~lib/string/String) => bool
      address = __lowerString(address) || __notnull();
      return exports.isBlacklisted(address) != 0;
    },
    addCollateral(amount, collateralType) {
      // assembly/index/addCollateral(u64, u32) => bool
      amount = amount || 0n;
      return exports.addCollateral(amount, collateralType) != 0;
    },
    removeCollateral(amount) {
      // assembly/index/removeCollateral(u64) => bool
      amount = amount || 0n;
      return exports.removeCollateral(amount) != 0;
    },
    getCollateralBalance(address, collateralType) {
      // assembly/index/getCollateralBalance(~lib/string/String, u32) => u64
      address = __lowerString(address) || __notnull();
      return BigInt.asUintN(64, exports.getCollateralBalance(address, collateralType));
    },
    getCollateralType(address) {
      // assembly/index/getCollateralType(~lib/string/String) => u32
      address = __lowerString(address) || __notnull();
      return exports.getCollateralType(address) >>> 0;
    },
    setCollateralRatio(ratio) {
      // assembly/index/setCollateralRatio(u64) => bool
      ratio = ratio || 0n;
      return exports.setCollateralRatio(ratio) != 0;
    },
    getCollateralRatio() {
      // assembly/index/getCollateralRatio() => u64
      return BigInt.asUintN(64, exports.getCollateralRatio());
    },
    setBinomTokenAddress(address) {
      // assembly/index/setBinomTokenAddress(~lib/string/String) => bool
      address = __lowerString(address) || __notnull();
      return exports.setBinomTokenAddress(address) != 0;
    },
    getFiatReserve() {
      // assembly/index/getFiatReserve() => u64
      return BigInt.asUintN(64, exports.getFiatReserve());
    },
    pause() {
      // assembly/index/pause() => bool
      return exports.pause() != 0;
    },
    unpause() {
      // assembly/index/unpause() => bool
      return exports.unpause() != 0;
    },
    transferOwnership(newOwner) {
      // assembly/index/transferOwnership(~lib/string/String) => bool
      newOwner = __lowerString(newOwner) || __notnull();
      return exports.transferOwnership(newOwner) != 0;
    },
    emitTransfer(from, to, amount) {
      // assembly/stablecoin/emitTransfer(~lib/string/String, ~lib/string/String, u64) => void
      from = __retain(__lowerString(from) || __notnull());
      to = __lowerString(to) || __notnull();
      amount = amount || 0n;
      try {
        exports.emitTransfer(from, to, amount);
      } finally {
        __release(from);
      }
    },
    emitMint(to, amount, collateralType) {
      // assembly/stablecoin/emitMint(~lib/string/String, u64, i32) => void
      to = __lowerString(to) || __notnull();
      amount = amount || 0n;
      exports.emitMint(to, amount, collateralType);
    },
    emitBurn(from, amount) {
      // assembly/stablecoin/emitBurn(~lib/string/String, u64) => void
      from = __lowerString(from) || __notnull();
      amount = amount || 0n;
      exports.emitBurn(from, amount);
    },
    emitCollateralAdded(from, amount, collateralType) {
      // assembly/stablecoin/emitCollateralAdded(~lib/string/String, u64, i32) => void
      from = __lowerString(from) || __notnull();
      amount = amount || 0n;
      exports.emitCollateralAdded(from, amount, collateralType);
    },
    emitCollateralRemoved(to, amount, collateralType) {
      // assembly/stablecoin/emitCollateralRemoved(~lib/string/String, u64, i32) => void
      to = __lowerString(to) || __notnull();
      amount = amount || 0n;
      exports.emitCollateralRemoved(to, amount, collateralType);
    },
  }, exports);
  function __liftString(pointer) {
    if (!pointer) return null;
    const
      end = pointer + new Uint32Array(memory.buffer)[pointer - 4 >>> 2] >>> 1,
      memoryU16 = new Uint16Array(memory.buffer);
    let
      start = pointer >>> 1,
      string = "";
    while (end - start > 1024) string += String.fromCharCode(...memoryU16.subarray(start, start += 1024));
    return string + String.fromCharCode(...memoryU16.subarray(start, end));
  }
  function __lowerString(value) {
    if (value == null) return 0;
    const
      length = value.length,
      pointer = exports.__new(length << 1, 2) >>> 0,
      memoryU16 = new Uint16Array(memory.buffer);
    for (let i = 0; i < length; ++i) memoryU16[(pointer >>> 1) + i] = value.charCodeAt(i);
    return pointer;
  }
  class Internref extends Number {}
  const registry = new FinalizationRegistry(__release);
  function __liftInternref(pointer) {
    if (!pointer) return null;
    const sentinel = new Internref(__retain(pointer));
    registry.register(sentinel, pointer);
    return sentinel;
  }
  const refcounts = new Map();
  function __retain(pointer) {
    if (pointer) {
      const refcount = refcounts.get(pointer);
      if (refcount) refcounts.set(pointer, refcount + 1);
      else refcounts.set(exports.__pin(pointer), 1);
    }
    return pointer;
  }
  function __release(pointer) {
    if (pointer) {
      const refcount = refcounts.get(pointer);
      if (refcount === 1) exports.__unpin(pointer), refcounts.delete(pointer);
      else if (refcount) refcounts.set(pointer, refcount - 1);
      else throw Error(`invalid refcount '${refcount}' for reference '${pointer}'`);
    }
  }
  function __notnull() {
    throw TypeError("value must not be null");
  }
  return adaptedExports;
}
export const {
  memory,
  createPaperDollar,
  getBalance,
  getTotalSupply,
  transfer,
  mint,
  burn,
  addMinter,
  removeMinter,
  isMinter,
  blacklist,
  unblacklist,
  isBlacklisted,
  addCollateral,
  removeCollateral,
  getCollateralBalance,
  getCollateralType,
  setCollateralRatio,
  getCollateralRatio,
  setBinomTokenAddress,
  getFiatReserve,
  pause,
  unpause,
  transferOwnership,
  add,
  multiply,
  store,
  retrieve,
  emitTransfer,
  emitMint,
  emitBurn,
  emitCollateralAdded,
  emitCollateralRemoved,
} = await (async url => instantiate(
  await (async () => {
    const isNodeOrBun = typeof process != "undefined" && process.versions != null && (process.versions.node != null || process.versions.bun != null);
    if (isNodeOrBun) { return globalThis.WebAssembly.compile(await (await import("node:fs/promises")).readFile(url)); }
    else { return await globalThis.WebAssembly.compileStreaming(globalThis.fetch(url)); }
  })(), {
  }
))(new URL("release.wasm", import.meta.url));
