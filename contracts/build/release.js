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
