(module
 (type $0 (func (param i32 i32) (result i32)))
 (type $1 (func (param i32) (result i32)))
 (type $2 (func (param i32 i32)))
 (type $3 (func (param i32)))
 (type $4 (func (result i32)))
 (type $5 (func (param i32 i64) (result i32)))
 (type $6 (func (param i32) (result i64)))
 (type $7 (func (param i32 i64 i32)))
 (type $8 (func))
 (type $9 (func (result i64)))
 (type $10 (func (param i64) (result i32)))
 (type $11 (func (param i32 i64)))
 (type $12 (func (param i32 i32) (result i64)))
 (type $13 (func (param i32 i32 i64) (result i32)))
 (type $14 (func (param i64 i32) (result i32)))
 (type $15 (func (param i32 i32 i32 i32)))
 (type $16 (func (param i32 i64) (result i64)))
 (type $17 (func (param i32 i32 i64)))
 (type $18 (func (param i32 i32 i32)))
 (type $19 (func (param i32 i64 i32 i32)))
 (type $20 (func (param i32 i32 i32) (result i64)))
 (type $21 (func (param i32 i32 i32 i32 i32) (result i32)))
 (type $22 (func (param i32 i64 i32) (result i32)))
 (import "env" "memory" (memory $0 1))
 (import "env" "abort" (func $~lib/builtins/abort (param i32 i32 i32 i32)))
 (global $assembly/contract/Context.caller (mut i32) (i32.const 32))
 (global $assembly/stablecoin/OWNER i32 (i32.const 64))
 (global $assembly/stablecoin/TOTAL_SUPPLY i32 (i32.const 96))
 (global $assembly/stablecoin/BALANCES_PREFIX i32 (i32.const 144))
 (global $assembly/stablecoin/PAUSED i32 (i32.const 192))
 (global $assembly/stablecoin/BLACKLIST_PREFIX i32 (i32.const 224))
 (global $assembly/stablecoin/MINTER_PREFIX i32 (i32.const 272))
 (global $assembly/stablecoin/COLLATERAL_PREFIX i32 (i32.const 320))
 (global $assembly/stablecoin/COLLATERAL_TYPE_PREFIX i32 (i32.const 368))
 (global $assembly/stablecoin/BINOM_TOKEN_ADDRESS i32 (i32.const 432))
 (global $assembly/stablecoin/COLLATERAL_RATIO i32 (i32.const 496))
 (global $assembly/stablecoin/FIAT_RESERVE i32 (i32.const 560))
 (global $assembly/stablecoin/CollateralType.FIAT i32 (i32.const 0))
 (global $assembly/stablecoin/CollateralType.BINOM i32 (i32.const 1))
 (global $~lib/shared/runtime/Runtime.Stub i32 (i32.const 0))
 (global $~lib/shared/runtime/Runtime.Minimal i32 (i32.const 1))
 (global $~lib/shared/runtime/Runtime.Incremental i32 (i32.const 2))
 (global $~lib/rt/stub/startOffset (mut i32) (i32.const 0))
 (global $~lib/rt/stub/offset (mut i32) (i32.const 0))
 (global $assembly/index/contractInstance (mut i32) (i32.const 0))
 (global $assembly/index/storedValue (mut i32) (i32.const 0))
 (global $~lib/native/ASC_SHRINK_LEVEL i32 (i32.const 0))
 (global $~lib/rt/__rtti_base i32 (i32.const 3216))
 (global $~lib/memory/__heap_base i32 (i32.const 3240))
 (data $0 (i32.const 12) "\1c\00\00\00\00\00\00\00\00\00\00\00\02\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00")
 (data $1 (i32.const 44) "\1c\00\00\00\00\00\00\00\00\00\00\00\02\00\00\00\n\00\00\00o\00w\00n\00e\00r\00\00\00")
 (data $2 (i32.const 76) ",\00\00\00\00\00\00\00\00\00\00\00\02\00\00\00\16\00\00\00t\00o\00t\00a\00l\00S\00u\00p\00p\00l\00y\00\00\00\00\00\00\00")
 (data $3 (i32.const 124) ",\00\00\00\00\00\00\00\00\00\00\00\02\00\00\00\10\00\00\00b\00a\00l\00a\00n\00c\00e\00_\00\00\00\00\00\00\00\00\00\00\00\00\00")
 (data $4 (i32.const 172) "\1c\00\00\00\00\00\00\00\00\00\00\00\02\00\00\00\0c\00\00\00p\00a\00u\00s\00e\00d\00")
 (data $5 (i32.const 204) ",\00\00\00\00\00\00\00\00\00\00\00\02\00\00\00\14\00\00\00b\00l\00a\00c\00k\00l\00i\00s\00t\00_\00\00\00\00\00\00\00\00\00")
 (data $6 (i32.const 252) ",\00\00\00\00\00\00\00\00\00\00\00\02\00\00\00\0e\00\00\00m\00i\00n\00t\00e\00r\00_\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00")
 (data $7 (i32.const 300) ",\00\00\00\00\00\00\00\00\00\00\00\02\00\00\00\16\00\00\00c\00o\00l\00l\00a\00t\00e\00r\00a\00l\00_\00\00\00\00\00\00\00")
 (data $8 (i32.const 348) "<\00\00\00\00\00\00\00\00\00\00\00\02\00\00\00 \00\00\00c\00o\00l\00l\00a\00t\00e\00r\00a\00l\00_\00t\00y\00p\00e\00_\00\00\00\00\00\00\00\00\00\00\00\00\00")
 (data $9 (i32.const 412) "<\00\00\00\00\00\00\00\00\00\00\00\02\00\00\00&\00\00\00b\00i\00n\00o\00m\00_\00t\00o\00k\00e\00n\00_\00a\00d\00d\00r\00e\00s\00s\00\00\00\00\00\00\00")
 (data $10 (i32.const 476) "<\00\00\00\00\00\00\00\00\00\00\00\02\00\00\00 \00\00\00c\00o\00l\00l\00a\00t\00e\00r\00a\00l\00_\00r\00a\00t\00i\00o\00\00\00\00\00\00\00\00\00\00\00\00\00")
 (data $11 (i32.const 540) ",\00\00\00\00\00\00\00\00\00\00\00\02\00\00\00\18\00\00\00f\00i\00a\00t\00_\00r\00e\00s\00e\00r\00v\00e\00\00\00\00\00")
 (data $12 (i32.const 588) "<\00\00\00\00\00\00\00\00\00\00\00\02\00\00\00(\00\00\00A\00l\00l\00o\00c\00a\00t\00i\00o\00n\00 \00t\00o\00o\00 \00l\00a\00r\00g\00e\00\00\00\00\00")
 (data $13 (i32.const 652) "<\00\00\00\00\00\00\00\00\00\00\00\02\00\00\00\1e\00\00\00~\00l\00i\00b\00/\00r\00t\00/\00s\00t\00u\00b\00.\00t\00s\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00")
 (data $14 (i32.const 716) "<\00\00\00\00\00\00\00\00\00\00\00\02\00\00\00$\00\00\00C\00o\00n\00t\00r\00a\00c\00t\00 \00i\00s\00 \00p\00a\00u\00s\00e\00d\00\00\00\00\00\00\00\00\00")
 (data $15 (i32.const 780) "<\00\00\00\00\00\00\00\00\00\00\00\02\00\00\00(\00\00\00a\00s\00s\00e\00m\00b\00l\00y\00/\00c\00o\00n\00t\00r\00a\00c\00t\00.\00t\00s\00\00\00\00\00")
 (data $16 (i32.const 844) "<\00\00\00\00\00\00\00\00\00\00\00\02\00\00\00,\00\00\00A\00d\00d\00r\00e\00s\00s\00 \00i\00s\00 \00b\00l\00a\00c\00k\00l\00i\00s\00t\00e\00d\00")
 (data $17 (i32.const 908) "<\00\00\00\00\00\00\00\00\00\00\00\02\00\00\00(\00\00\00I\00n\00s\00u\00f\00f\00i\00c\00i\00e\00n\00t\00 \00b\00a\00l\00a\00n\00c\00e\00\00\00\00\00")
 (data $18 (i32.const 972) "<\00\00\00\00\00\00\00\00\00\00\00\02\00\00\00*\00\00\00O\00n\00l\00y\00 \00m\00i\00n\00t\00e\00r\00s\00 \00c\00a\00n\00 \00m\00i\00n\00t\00\00\00")
 (data $19 (i32.const 1036) "|\00\00\00\00\00\00\00\00\00\00\00\02\00\00\00d\00\00\00t\00o\00S\00t\00r\00i\00n\00g\00(\00)\00 \00r\00a\00d\00i\00x\00 \00a\00r\00g\00u\00m\00e\00n\00t\00 \00m\00u\00s\00t\00 \00b\00e\00 \00b\00e\00t\00w\00e\00e\00n\00 \002\00 \00a\00n\00d\00 \003\006\00\00\00\00\00\00\00\00\00")
 (data $20 (i32.const 1164) "<\00\00\00\00\00\00\00\00\00\00\00\02\00\00\00&\00\00\00~\00l\00i\00b\00/\00u\00t\00i\00l\00/\00n\00u\00m\00b\00e\00r\00.\00t\00s\00\00\00\00\00\00\00")
 (data $21 (i32.const 1228) "\1c\00\00\00\00\00\00\00\00\00\00\00\02\00\00\00\02\00\00\000\00\00\00\00\00\00\00\00\00\00\00")
 (data $22 (i32.const 1260) "0\000\000\001\000\002\000\003\000\004\000\005\000\006\000\007\000\008\000\009\001\000\001\001\001\002\001\003\001\004\001\005\001\006\001\007\001\008\001\009\002\000\002\001\002\002\002\003\002\004\002\005\002\006\002\007\002\008\002\009\003\000\003\001\003\002\003\003\003\004\003\005\003\006\003\007\003\008\003\009\004\000\004\001\004\002\004\003\004\004\004\005\004\006\004\007\004\008\004\009\005\000\005\001\005\002\005\003\005\004\005\005\005\006\005\007\005\008\005\009\006\000\006\001\006\002\006\003\006\004\006\005\006\006\006\007\006\008\006\009\007\000\007\001\007\002\007\003\007\004\007\005\007\006\007\007\007\008\007\009\008\000\008\001\008\002\008\003\008\004\008\005\008\006\008\007\008\008\008\009\009\000\009\001\009\002\009\003\009\004\009\005\009\006\009\007\009\008\009\009\00")
 (data $23 (i32.const 1660) "\1c\04\00\00\00\00\00\00\00\00\00\00\02\00\00\00\00\04\00\000\000\000\001\000\002\000\003\000\004\000\005\000\006\000\007\000\008\000\009\000\00a\000\00b\000\00c\000\00d\000\00e\000\00f\001\000\001\001\001\002\001\003\001\004\001\005\001\006\001\007\001\008\001\009\001\00a\001\00b\001\00c\001\00d\001\00e\001\00f\002\000\002\001\002\002\002\003\002\004\002\005\002\006\002\007\002\008\002\009\002\00a\002\00b\002\00c\002\00d\002\00e\002\00f\003\000\003\001\003\002\003\003\003\004\003\005\003\006\003\007\003\008\003\009\003\00a\003\00b\003\00c\003\00d\003\00e\003\00f\004\000\004\001\004\002\004\003\004\004\004\005\004\006\004\007\004\008\004\009\004\00a\004\00b\004\00c\004\00d\004\00e\004\00f\005\000\005\001\005\002\005\003\005\004\005\005\005\006\005\007\005\008\005\009\005\00a\005\00b\005\00c\005\00d\005\00e\005\00f\006\000\006\001\006\002\006\003\006\004\006\005\006\006\006\007\006\008\006\009\006\00a\006\00b\006\00c\006\00d\006\00e\006\00f\007\000\007\001\007\002\007\003\007\004\007\005\007\006\007\007\007\008\007\009\007\00a\007\00b\007\00c\007\00d\007\00e\007\00f\008\000\008\001\008\002\008\003\008\004\008\005\008\006\008\007\008\008\008\009\008\00a\008\00b\008\00c\008\00d\008\00e\008\00f\009\000\009\001\009\002\009\003\009\004\009\005\009\006\009\007\009\008\009\009\009\00a\009\00b\009\00c\009\00d\009\00e\009\00f\00a\000\00a\001\00a\002\00a\003\00a\004\00a\005\00a\006\00a\007\00a\008\00a\009\00a\00a\00a\00b\00a\00c\00a\00d\00a\00e\00a\00f\00b\000\00b\001\00b\002\00b\003\00b\004\00b\005\00b\006\00b\007\00b\008\00b\009\00b\00a\00b\00b\00b\00c\00b\00d\00b\00e\00b\00f\00c\000\00c\001\00c\002\00c\003\00c\004\00c\005\00c\006\00c\007\00c\008\00c\009\00c\00a\00c\00b\00c\00c\00c\00d\00c\00e\00c\00f\00d\000\00d\001\00d\002\00d\003\00d\004\00d\005\00d\006\00d\007\00d\008\00d\009\00d\00a\00d\00b\00d\00c\00d\00d\00d\00e\00d\00f\00e\000\00e\001\00e\002\00e\003\00e\004\00e\005\00e\006\00e\007\00e\008\00e\009\00e\00a\00e\00b\00e\00c\00e\00d\00e\00e\00e\00f\00f\000\00f\001\00f\002\00f\003\00f\004\00f\005\00f\006\00f\007\00f\008\00f\009\00f\00a\00f\00b\00f\00c\00f\00d\00f\00e\00f\00f\00\00\00\00\00\00\00\00\00\00\00\00\00")
 (data $24 (i32.const 2716) "\\\00\00\00\00\00\00\00\00\00\00\00\02\00\00\00H\00\00\000\001\002\003\004\005\006\007\008\009\00a\00b\00c\00d\00e\00f\00g\00h\00i\00j\00k\00l\00m\00n\00o\00p\00q\00r\00s\00t\00u\00v\00w\00x\00y\00z\00\00\00\00\00")
 (data $25 (i32.const 2812) "\1c\00\00\00\00\00\00\00\00\00\00\00\02\00\00\00\02\00\00\00_\00\00\00\00\00\00\00\00\00\00\00")
 (data $26 (i32.const 2844) "L\00\00\00\00\00\00\00\00\00\00\00\02\00\00\00.\00\00\00I\00n\00s\00u\00f\00f\00i\00c\00i\00e\00n\00t\00 \00c\00o\00l\00l\00a\00t\00e\00r\00a\00l\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00")
 (data $27 (i32.const 2924) "\\\00\00\00\00\00\00\00\00\00\00\00\02\00\00\00B\00\00\00O\00n\00l\00y\00 \00o\00w\00n\00e\00r\00 \00c\00a\00n\00 \00c\00a\00l\00l\00 \00t\00h\00i\00s\00 \00f\00u\00n\00c\00t\00i\00o\00n\00\00\00\00\00\00\00\00\00\00\00")
 (data $28 (i32.const 3020) "\\\00\00\00\00\00\00\00\00\00\00\00\02\00\00\00B\00\00\00C\00o\00l\00l\00a\00t\00e\00r\00a\00l\00 \00r\00a\00t\00i\00o\00 \00w\00o\00u\00l\00d\00 \00b\00e\00 \00t\00o\00o\00 \00l\00o\00w\00\00\00\00\00\00\00\00\00\00\00")
 (data $29 (i32.const 3116) "\\\00\00\00\00\00\00\00\00\00\00\00\02\00\00\00L\00\00\00C\00o\00l\00l\00a\00t\00e\00r\00a\00l\00 \00r\00a\00t\00i\00o\00 \00m\00u\00s\00t\00 \00b\00e\00 \00a\00t\00 \00l\00e\00a\00s\00t\00 \001\000\000\00%\00")
 (data $30 (i32.const 3216) "\05\00\00\00 \00\00\00 \00\00\00 \00\00\00\00\00\00\00 \00\00\00")
 (table $0 1 1 funcref)
 (elem $0 (i32.const 1))
 (export "createPaperDollar" (func $assembly/index/createPaperDollar))
 (export "getBalance" (func $assembly/index/getBalance))
 (export "getTotalSupply" (func $assembly/index/getTotalSupply))
 (export "transfer" (func $assembly/index/transfer))
 (export "mint" (func $assembly/index/mint))
 (export "burn" (func $assembly/index/burn))
 (export "addMinter" (func $assembly/index/addMinter))
 (export "removeMinter" (func $assembly/index/removeMinter))
 (export "isMinter" (func $assembly/index/isMinter))
 (export "blacklist" (func $assembly/index/blacklist))
 (export "unblacklist" (func $assembly/index/unblacklist))
 (export "isBlacklisted" (func $assembly/index/isBlacklisted))
 (export "addCollateral" (func $assembly/index/addCollateral))
 (export "removeCollateral" (func $assembly/index/removeCollateral))
 (export "getCollateralBalance" (func $assembly/index/getCollateralBalance))
 (export "getCollateralType" (func $assembly/index/getCollateralType))
 (export "setCollateralRatio" (func $assembly/index/setCollateralRatio))
 (export "getCollateralRatio" (func $assembly/index/getCollateralRatio))
 (export "setBinomTokenAddress" (func $assembly/index/setBinomTokenAddress))
 (export "getFiatReserve" (func $assembly/index/getFiatReserve))
 (export "pause" (func $assembly/index/pause))
 (export "unpause" (func $assembly/index/unpause))
 (export "transferOwnership" (func $assembly/index/transferOwnership))
 (export "add" (func $assembly/index/add))
 (export "multiply" (func $assembly/index/multiply))
 (export "store" (func $assembly/index/store))
 (export "retrieve" (func $assembly/index/retrieve))
 (export "emitTransfer" (func $assembly/stablecoin/emitTransfer))
 (export "emitMint" (func $assembly/stablecoin/emitMint))
 (export "emitBurn" (func $assembly/stablecoin/emitBurn))
 (export "emitCollateralAdded" (func $assembly/stablecoin/emitCollateralAdded))
 (export "emitCollateralRemoved" (func $assembly/stablecoin/emitCollateralRemoved))
 (export "__new" (func $~lib/rt/stub/__new))
 (export "__pin" (func $~lib/rt/stub/__pin))
 (export "__unpin" (func $~lib/rt/stub/__unpin))
 (export "__collect" (func $~lib/rt/stub/__collect))
 (export "__rtti_base" (global $~lib/rt/__rtti_base))
 (export "memory" (memory $0))
 (start $~start)
 (func $assembly/contract/storage.get<~lib/string/String> (param $key i32) (param $defaultValue i32) (result i32)
  local.get $defaultValue
  return
 )
 (func $~lib/rt/common/OBJECT#get:rtSize (param $this i32) (result i32)
  local.get $this
  i32.load offset=16
 )
 (func $~lib/string/String#get:length (param $this i32) (result i32)
  local.get $this
  i32.const 20
  i32.sub
  call $~lib/rt/common/OBJECT#get:rtSize
  i32.const 1
  i32.shr_u
  return
 )
 (func $~lib/string/String.__not (param $str i32) (result i32)
  local.get $str
  i32.const 0
  i32.eq
  if (result i32)
   i32.const 1
  else
   local.get $str
   call $~lib/string/String#get:length
   i32.eqz
  end
  return
 )
 (func $assembly/contract/storage.set<~lib/string/String> (param $key i32) (param $value i32)
 )
 (func $assembly/contract/storage.set<u64> (param $key i32) (param $value i64)
 )
 (func $assembly/contract/storage.set<bool> (param $key i32) (param $value i32)
 )
 (func $~lib/rt/stub/maybeGrowMemory (param $newOffset i32)
  (local $pagesBefore i32)
  (local $maxOffset i32)
  (local $pagesNeeded i32)
  (local $4 i32)
  (local $5 i32)
  (local $pagesWanted i32)
  memory.size
  local.set $pagesBefore
  local.get $pagesBefore
  i32.const 16
  i32.shl
  i32.const 15
  i32.add
  i32.const 15
  i32.const -1
  i32.xor
  i32.and
  local.set $maxOffset
  local.get $newOffset
  local.get $maxOffset
  i32.gt_u
  if
   local.get $newOffset
   local.get $maxOffset
   i32.sub
   i32.const 65535
   i32.add
   i32.const 65535
   i32.const -1
   i32.xor
   i32.and
   i32.const 16
   i32.shr_u
   local.set $pagesNeeded
   local.get $pagesBefore
   local.tee $4
   local.get $pagesNeeded
   local.tee $5
   local.get $4
   local.get $5
   i32.gt_s
   select
   local.set $pagesWanted
   local.get $pagesWanted
   memory.grow
   i32.const 0
   i32.lt_s
   if
    local.get $pagesNeeded
    memory.grow
    i32.const 0
    i32.lt_s
    if
     unreachable
    end
   end
  end
  local.get $newOffset
  global.set $~lib/rt/stub/offset
 )
 (func $~lib/rt/common/BLOCK#set:mmInfo (param $this i32) (param $mmInfo i32)
  local.get $this
  local.get $mmInfo
  i32.store
 )
 (func $~lib/rt/stub/__alloc (param $size i32) (result i32)
  (local $block i32)
  (local $ptr i32)
  (local $size|3 i32)
  (local $payloadSize i32)
  local.get $size
  i32.const 1073741820
  i32.gt_u
  if
   i32.const 608
   i32.const 672
   i32.const 33
   i32.const 29
   call $~lib/builtins/abort
   unreachable
  end
  global.get $~lib/rt/stub/offset
  local.set $block
  global.get $~lib/rt/stub/offset
  i32.const 4
  i32.add
  local.set $ptr
  block $~lib/rt/stub/computeSize|inlined.0 (result i32)
   local.get $size
   local.set $size|3
   local.get $size|3
   i32.const 4
   i32.add
   i32.const 15
   i32.add
   i32.const 15
   i32.const -1
   i32.xor
   i32.and
   i32.const 4
   i32.sub
   br $~lib/rt/stub/computeSize|inlined.0
  end
  local.set $payloadSize
  local.get $ptr
  local.get $payloadSize
  i32.add
  call $~lib/rt/stub/maybeGrowMemory
  local.get $block
  local.get $payloadSize
  call $~lib/rt/common/BLOCK#set:mmInfo
  local.get $ptr
  return
 )
 (func $~lib/rt/common/OBJECT#set:gcInfo (param $this i32) (param $gcInfo i32)
  local.get $this
  local.get $gcInfo
  i32.store offset=4
 )
 (func $~lib/rt/common/OBJECT#set:gcInfo2 (param $this i32) (param $gcInfo2 i32)
  local.get $this
  local.get $gcInfo2
  i32.store offset=8
 )
 (func $~lib/rt/common/OBJECT#set:rtId (param $this i32) (param $rtId i32)
  local.get $this
  local.get $rtId
  i32.store offset=12
 )
 (func $~lib/rt/common/OBJECT#set:rtSize (param $this i32) (param $rtSize i32)
  local.get $this
  local.get $rtSize
  i32.store offset=16
 )
 (func $~lib/rt/stub/__new (param $size i32) (param $id i32) (result i32)
  (local $ptr i32)
  (local $object i32)
  local.get $size
  i32.const 1073741804
  i32.gt_u
  if
   i32.const 608
   i32.const 672
   i32.const 86
   i32.const 30
   call $~lib/builtins/abort
   unreachable
  end
  i32.const 16
  local.get $size
  i32.add
  call $~lib/rt/stub/__alloc
  local.set $ptr
  local.get $ptr
  i32.const 4
  i32.sub
  local.set $object
  local.get $object
  i32.const 0
  call $~lib/rt/common/OBJECT#set:gcInfo
  local.get $object
  i32.const 0
  call $~lib/rt/common/OBJECT#set:gcInfo2
  local.get $object
  local.get $id
  call $~lib/rt/common/OBJECT#set:rtId
  local.get $object
  local.get $size
  call $~lib/rt/common/OBJECT#set:rtSize
  local.get $ptr
  i32.const 16
  i32.add
  return
 )
 (func $assembly/stablecoin/PaperDollar#constructor (param $this i32) (result i32)
  local.get $this
  i32.eqz
  if
   i32.const 0
   i32.const 4
   call $~lib/rt/stub/__new
   local.set $this
  end
  global.get $assembly/stablecoin/OWNER
  i32.const 32
  call $assembly/contract/storage.get<~lib/string/String>
  call $~lib/string/String.__not
  if
   global.get $assembly/stablecoin/OWNER
   global.get $assembly/contract/Context.caller
   call $assembly/contract/storage.set<~lib/string/String>
   global.get $assembly/stablecoin/TOTAL_SUPPLY
   i64.const 0
   call $assembly/contract/storage.set<u64>
   global.get $assembly/stablecoin/PAUSED
   i32.const 0
   call $assembly/contract/storage.set<bool>
   global.get $assembly/stablecoin/COLLATERAL_RATIO
   i64.const 150
   call $assembly/contract/storage.set<u64>
   global.get $assembly/stablecoin/FIAT_RESERVE
   i64.const 0
   call $assembly/contract/storage.set<u64>
  end
  local.get $this
 )
 (func $start:assembly/index
  global.get $~lib/memory/__heap_base
  i32.const 4
  i32.add
  i32.const 15
  i32.add
  i32.const 15
  i32.const -1
  i32.xor
  i32.and
  i32.const 4
  i32.sub
  global.set $~lib/rt/stub/startOffset
  global.get $~lib/rt/stub/startOffset
  global.set $~lib/rt/stub/offset
  i32.const 0
  call $assembly/stablecoin/PaperDollar#constructor
  global.set $assembly/index/contractInstance
 )
 (func $assembly/index/createPaperDollar (result i32)
  global.get $assembly/index/contractInstance
  return
 )
 (func $assembly/index/getContract (result i32)
  global.get $assembly/index/contractInstance
  return
 )
 (func $~lib/string/String#concat (param $this i32) (param $other i32) (result i32)
  (local $thisSize i32)
  (local $otherSize i32)
  (local $outSize i32)
  (local $out i32)
  local.get $this
  call $~lib/string/String#get:length
  i32.const 1
  i32.shl
  local.set $thisSize
  local.get $other
  call $~lib/string/String#get:length
  i32.const 1
  i32.shl
  local.set $otherSize
  local.get $thisSize
  local.get $otherSize
  i32.add
  local.set $outSize
  local.get $outSize
  i32.const 0
  i32.eq
  if
   i32.const 32
   return
  end
  local.get $outSize
  i32.const 2
  call $~lib/rt/stub/__new
  local.set $out
  local.get $out
  local.get $this
  local.get $thisSize
  memory.copy
  local.get $out
  local.get $thisSize
  i32.add
  local.get $other
  local.get $otherSize
  memory.copy
  local.get $out
  return
 )
 (func $~lib/string/String.__concat (param $left i32) (param $right i32) (result i32)
  local.get $left
  local.get $right
  call $~lib/string/String#concat
  return
 )
 (func $assembly/contract/storage.get<u64> (param $key i32) (param $defaultValue i64) (result i64)
  local.get $defaultValue
  return
 )
 (func $assembly/stablecoin/PaperDollar#getBalance (param $this i32) (param $address i32) (result i64)
  global.get $assembly/stablecoin/BALANCES_PREFIX
  local.get $address
  call $~lib/string/String.__concat
  i64.const 0
  call $assembly/contract/storage.get<u64>
  return
 )
 (func $assembly/index/getBalance (param $address i32) (result i64)
  call $assembly/index/getContract
  local.get $address
  call $assembly/stablecoin/PaperDollar#getBalance
  return
 )
 (func $assembly/stablecoin/PaperDollar#getTotalSupply (param $this i32) (result i64)
  global.get $assembly/stablecoin/TOTAL_SUPPLY
  i64.const 0
  call $assembly/contract/storage.get<u64>
  return
 )
 (func $assembly/index/getTotalSupply (result i64)
  call $assembly/index/getContract
  call $assembly/stablecoin/PaperDollar#getTotalSupply
  return
 )
 (func $assembly/contract/storage.get<bool> (param $key i32) (param $defaultValue i32) (result i32)
  local.get $defaultValue
  return
 )
 (func $assembly/contract/assert (param $condition i32) (param $message i32)
  local.get $condition
  i32.eqz
  if
   local.get $message
   i32.const 800
   i32.const 35
   i32.const 9
   call $~lib/builtins/abort
   unreachable
  end
 )
 (func $assembly/stablecoin/PaperDollar#checkNotPaused (param $this i32)
  global.get $assembly/stablecoin/PAUSED
  i32.const 0
  call $assembly/contract/storage.get<bool>
  i32.eqz
  i32.const 736
  call $assembly/contract/assert
 )
 (func $assembly/stablecoin/PaperDollar#isBlacklisted (param $this i32) (param $address i32) (result i32)
  global.get $assembly/stablecoin/BLACKLIST_PREFIX
  local.get $address
  call $~lib/string/String.__concat
  i32.const 0
  call $assembly/contract/storage.get<bool>
  i32.const 0
  i32.ne
  return
 )
 (func $assembly/stablecoin/PaperDollar#checkNotBlacklisted (param $this i32) (param $address i32)
  local.get $this
  local.get $address
  call $assembly/stablecoin/PaperDollar#isBlacklisted
  i32.eqz
  i32.const 864
  call $assembly/contract/assert
 )
 (func $assembly/stablecoin/emitTransfer (param $from i32) (param $to i32) (param $amount i64)
 )
 (func $assembly/stablecoin/PaperDollar#transfer (param $this i32) (param $to i32) (param $amount i64) (result i32)
  (local $from i32)
  (local $fromBalance i64)
  global.get $assembly/contract/Context.caller
  local.set $from
  local.get $this
  call $assembly/stablecoin/PaperDollar#checkNotPaused
  local.get $this
  local.get $from
  call $assembly/stablecoin/PaperDollar#checkNotBlacklisted
  local.get $this
  local.get $to
  call $assembly/stablecoin/PaperDollar#checkNotBlacklisted
  local.get $this
  local.get $from
  call $assembly/stablecoin/PaperDollar#getBalance
  local.set $fromBalance
  local.get $fromBalance
  local.get $amount
  i64.ge_u
  i32.const 928
  call $assembly/contract/assert
  global.get $assembly/stablecoin/BALANCES_PREFIX
  local.get $from
  call $~lib/string/String.__concat
  local.get $fromBalance
  local.get $amount
  i64.sub
  call $assembly/contract/storage.set<u64>
  global.get $assembly/stablecoin/BALANCES_PREFIX
  local.get $to
  call $~lib/string/String.__concat
  local.get $this
  local.get $to
  call $assembly/stablecoin/PaperDollar#getBalance
  local.get $amount
  i64.add
  call $assembly/contract/storage.set<u64>
  local.get $from
  local.get $to
  local.get $amount
  call $assembly/stablecoin/emitTransfer
  i32.const 1
  return
 )
 (func $assembly/index/transfer (param $to i32) (param $amount i64) (result i32)
  call $assembly/index/getContract
  local.get $to
  local.get $amount
  call $assembly/stablecoin/PaperDollar#transfer
  return
 )
 (func $assembly/stablecoin/PaperDollar#isMinter (param $this i32) (param $address i32) (result i32)
  global.get $assembly/stablecoin/MINTER_PREFIX
  local.get $address
  call $~lib/string/String.__concat
  i32.const 0
  call $assembly/contract/storage.get<bool>
  i32.const 0
  i32.ne
  return
 )
 (func $assembly/contract/storage.get<i32> (param $key i32) (param $defaultValue i32) (result i32)
  local.get $defaultValue
  return
 )
 (func $assembly/stablecoin/PaperDollar#getCollateralType (param $this i32) (param $address i32) (result i32)
  global.get $assembly/stablecoin/COLLATERAL_TYPE_PREFIX
  local.get $address
  call $~lib/string/String.__concat
  global.get $assembly/stablecoin/CollateralType.FIAT
  call $assembly/contract/storage.get<i32>
  return
 )
 (func $~lib/util/number/decimalCount32 (param $value i32) (result i32)
  local.get $value
  i32.const 100000
  i32.lt_u
  if
   local.get $value
   i32.const 100
   i32.lt_u
   if
    i32.const 1
    local.get $value
    i32.const 10
    i32.ge_u
    i32.add
    return
   else
    i32.const 3
    local.get $value
    i32.const 10000
    i32.ge_u
    i32.add
    local.get $value
    i32.const 1000
    i32.ge_u
    i32.add
    return
   end
   unreachable
  else
   local.get $value
   i32.const 10000000
   i32.lt_u
   if
    i32.const 6
    local.get $value
    i32.const 1000000
    i32.ge_u
    i32.add
    return
   else
    i32.const 8
    local.get $value
    i32.const 1000000000
    i32.ge_u
    i32.add
    local.get $value
    i32.const 100000000
    i32.ge_u
    i32.add
    return
   end
   unreachable
  end
  unreachable
 )
 (func $~lib/util/number/utoa32_dec_lut (param $buffer i32) (param $num i32) (param $offset i32)
  (local $t i32)
  (local $r i32)
  (local $d1 i32)
  (local $d2 i32)
  (local $digits1 i64)
  (local $digits2 i64)
  (local $t|9 i32)
  (local $d1|10 i32)
  (local $digits i32)
  (local $digits|12 i32)
  (local $digit i32)
  loop $while-continue|0
   local.get $num
   i32.const 10000
   i32.ge_u
   if
    local.get $num
    i32.const 10000
    i32.div_u
    local.set $t
    local.get $num
    i32.const 10000
    i32.rem_u
    local.set $r
    local.get $t
    local.set $num
    local.get $r
    i32.const 100
    i32.div_u
    local.set $d1
    local.get $r
    i32.const 100
    i32.rem_u
    local.set $d2
    i32.const 1260
    local.get $d1
    i32.const 2
    i32.shl
    i32.add
    i64.load32_u
    local.set $digits1
    i32.const 1260
    local.get $d2
    i32.const 2
    i32.shl
    i32.add
    i64.load32_u
    local.set $digits2
    local.get $offset
    i32.const 4
    i32.sub
    local.set $offset
    local.get $buffer
    local.get $offset
    i32.const 1
    i32.shl
    i32.add
    local.get $digits1
    local.get $digits2
    i64.const 32
    i64.shl
    i64.or
    i64.store
    br $while-continue|0
   end
  end
  local.get $num
  i32.const 100
  i32.ge_u
  if
   local.get $num
   i32.const 100
   i32.div_u
   local.set $t|9
   local.get $num
   i32.const 100
   i32.rem_u
   local.set $d1|10
   local.get $t|9
   local.set $num
   local.get $offset
   i32.const 2
   i32.sub
   local.set $offset
   i32.const 1260
   local.get $d1|10
   i32.const 2
   i32.shl
   i32.add
   i32.load
   local.set $digits
   local.get $buffer
   local.get $offset
   i32.const 1
   i32.shl
   i32.add
   local.get $digits
   i32.store
  end
  local.get $num
  i32.const 10
  i32.ge_u
  if
   local.get $offset
   i32.const 2
   i32.sub
   local.set $offset
   i32.const 1260
   local.get $num
   i32.const 2
   i32.shl
   i32.add
   i32.load
   local.set $digits|12
   local.get $buffer
   local.get $offset
   i32.const 1
   i32.shl
   i32.add
   local.get $digits|12
   i32.store
  else
   local.get $offset
   i32.const 1
   i32.sub
   local.set $offset
   i32.const 48
   local.get $num
   i32.add
   local.set $digit
   local.get $buffer
   local.get $offset
   i32.const 1
   i32.shl
   i32.add
   local.get $digit
   i32.store16
  end
 )
 (func $~lib/util/number/utoa_hex_lut (param $buffer i32) (param $num i64) (param $offset i32)
  loop $while-continue|0
   local.get $offset
   i32.const 2
   i32.ge_u
   if
    local.get $offset
    i32.const 2
    i32.sub
    local.set $offset
    local.get $buffer
    local.get $offset
    i32.const 1
    i32.shl
    i32.add
    i32.const 1680
    local.get $num
    i32.wrap_i64
    i32.const 255
    i32.and
    i32.const 2
    i32.shl
    i32.add
    i32.load
    i32.store
    local.get $num
    i64.const 8
    i64.shr_u
    local.set $num
    br $while-continue|0
   end
  end
  local.get $offset
  i32.const 1
  i32.and
  if
   local.get $buffer
   i32.const 1680
   local.get $num
   i32.wrap_i64
   i32.const 6
   i32.shl
   i32.add
   i32.load16_u
   i32.store16
  end
 )
 (func $~lib/util/number/ulog_base (param $num i64) (param $base i32) (result i32)
  (local $value i32)
  (local $b64 i64)
  (local $b i64)
  (local $e i32)
  block $~lib/util/number/isPowerOf2<i32>|inlined.0 (result i32)
   local.get $base
   local.set $value
   local.get $value
   i32.popcnt
   i32.const 1
   i32.eq
   br $~lib/util/number/isPowerOf2<i32>|inlined.0
  end
  if
   i32.const 63
   local.get $num
   i64.clz
   i32.wrap_i64
   i32.sub
   i32.const 31
   local.get $base
   i32.clz
   i32.sub
   i32.div_u
   i32.const 1
   i32.add
   return
  end
  local.get $base
  i64.extend_i32_s
  local.set $b64
  local.get $b64
  local.set $b
  i32.const 1
  local.set $e
  loop $while-continue|0
   local.get $num
   local.get $b
   i64.ge_u
   if
    local.get $num
    local.get $b
    i64.div_u
    local.set $num
    local.get $b
    local.get $b
    i64.mul
    local.set $b
    local.get $e
    i32.const 1
    i32.shl
    local.set $e
    br $while-continue|0
   end
  end
  loop $while-continue|1
   local.get $num
   i64.const 1
   i64.ge_u
   if
    local.get $num
    local.get $b64
    i64.div_u
    local.set $num
    local.get $e
    i32.const 1
    i32.add
    local.set $e
    br $while-continue|1
   end
  end
  local.get $e
  i32.const 1
  i32.sub
  return
 )
 (func $~lib/util/number/utoa64_any_core (param $buffer i32) (param $num i64) (param $offset i32) (param $radix i32)
  (local $base i64)
  (local $shift i64)
  (local $mask i64)
  (local $q i64)
  local.get $radix
  i64.extend_i32_s
  local.set $base
  local.get $radix
  local.get $radix
  i32.const 1
  i32.sub
  i32.and
  i32.const 0
  i32.eq
  if
   local.get $radix
   i32.ctz
   i32.const 7
   i32.and
   i64.extend_i32_s
   local.set $shift
   local.get $base
   i64.const 1
   i64.sub
   local.set $mask
   loop $do-loop|0
    local.get $offset
    i32.const 1
    i32.sub
    local.set $offset
    local.get $buffer
    local.get $offset
    i32.const 1
    i32.shl
    i32.add
    i32.const 2736
    local.get $num
    local.get $mask
    i64.and
    i32.wrap_i64
    i32.const 1
    i32.shl
    i32.add
    i32.load16_u
    i32.store16
    local.get $num
    local.get $shift
    i64.shr_u
    local.set $num
    local.get $num
    i64.const 0
    i64.ne
    br_if $do-loop|0
   end
  else
   loop $do-loop|1
    local.get $offset
    i32.const 1
    i32.sub
    local.set $offset
    local.get $num
    local.get $base
    i64.div_u
    local.set $q
    local.get $buffer
    local.get $offset
    i32.const 1
    i32.shl
    i32.add
    i32.const 2736
    local.get $num
    local.get $q
    local.get $base
    i64.mul
    i64.sub
    i32.wrap_i64
    i32.const 1
    i32.shl
    i32.add
    i32.load16_u
    i32.store16
    local.get $q
    local.set $num
    local.get $num
    i64.const 0
    i64.ne
    br_if $do-loop|1
   end
  end
 )
 (func $~lib/util/number/itoa32 (param $value i32) (param $radix i32) (result i32)
  (local $sign i32)
  (local $out i32)
  (local $decimals i32)
  (local $buffer i32)
  (local $num i32)
  (local $offset i32)
  (local $decimals|8 i32)
  (local $buffer|9 i32)
  (local $num|10 i32)
  (local $offset|11 i32)
  (local $val32 i32)
  (local $decimals|13 i32)
  local.get $radix
  i32.const 2
  i32.lt_s
  if (result i32)
   i32.const 1
  else
   local.get $radix
   i32.const 36
   i32.gt_s
  end
  if
   i32.const 1056
   i32.const 1184
   i32.const 373
   i32.const 5
   call $~lib/builtins/abort
   unreachable
  end
  local.get $value
  i32.eqz
  if
   i32.const 1248
   return
  end
  local.get $value
  i32.const 31
  i32.shr_u
  i32.const 1
  i32.shl
  local.set $sign
  local.get $sign
  if
   i32.const 0
   local.get $value
   i32.sub
   local.set $value
  end
  local.get $radix
  i32.const 10
  i32.eq
  if
   local.get $value
   call $~lib/util/number/decimalCount32
   local.set $decimals
   local.get $decimals
   i32.const 1
   i32.shl
   local.get $sign
   i32.add
   i32.const 2
   call $~lib/rt/stub/__new
   local.set $out
   local.get $out
   local.get $sign
   i32.add
   local.set $buffer
   local.get $value
   local.set $num
   local.get $decimals
   local.set $offset
   i32.const 0
   i32.const 1
   i32.ge_s
   drop
   local.get $buffer
   local.get $num
   local.get $offset
   call $~lib/util/number/utoa32_dec_lut
  else
   local.get $radix
   i32.const 16
   i32.eq
   if
    i32.const 31
    local.get $value
    i32.clz
    i32.sub
    i32.const 2
    i32.shr_s
    i32.const 1
    i32.add
    local.set $decimals|8
    local.get $decimals|8
    i32.const 1
    i32.shl
    local.get $sign
    i32.add
    i32.const 2
    call $~lib/rt/stub/__new
    local.set $out
    local.get $out
    local.get $sign
    i32.add
    local.set $buffer|9
    local.get $value
    local.set $num|10
    local.get $decimals|8
    local.set $offset|11
    i32.const 0
    i32.const 1
    i32.ge_s
    drop
    local.get $buffer|9
    local.get $num|10
    i64.extend_i32_u
    local.get $offset|11
    call $~lib/util/number/utoa_hex_lut
   else
    local.get $value
    local.set $val32
    local.get $val32
    i64.extend_i32_u
    local.get $radix
    call $~lib/util/number/ulog_base
    local.set $decimals|13
    local.get $decimals|13
    i32.const 1
    i32.shl
    local.get $sign
    i32.add
    i32.const 2
    call $~lib/rt/stub/__new
    local.set $out
    local.get $out
    local.get $sign
    i32.add
    local.get $val32
    i64.extend_i32_u
    local.get $decimals|13
    local.get $radix
    call $~lib/util/number/utoa64_any_core
   end
  end
  local.get $sign
  if
   local.get $out
   i32.const 45
   i32.store16
  end
  local.get $out
  return
 )
 (func $~lib/number/I32#toString (param $this i32) (param $radix i32) (result i32)
  local.get $this
  local.get $radix
  call $~lib/util/number/itoa32
  return
 )
 (func $assembly/stablecoin/PaperDollar#getCollateralBalance (param $this i32) (param $address i32) (param $collateralType i32) (result i64)
  global.get $assembly/stablecoin/COLLATERAL_PREFIX
  local.get $collateralType
  i32.const 10
  call $~lib/number/I32#toString
  call $~lib/string/String.__concat
  i32.const 2832
  call $~lib/string/String.__concat
  local.get $address
  call $~lib/string/String.__concat
  i64.const 0
  call $assembly/contract/storage.get<u64>
  return
 )
 (func $assembly/stablecoin/PaperDollar#getCollateralRatio (param $this i32) (result i64)
  global.get $assembly/stablecoin/COLLATERAL_RATIO
  i64.const 150
  call $assembly/contract/storage.get<u64>
  return
 )
 (func $assembly/stablecoin/emitMint (param $to i32) (param $amount i64) (param $collateralType i32)
 )
 (func $assembly/stablecoin/PaperDollar#mint (param $this i32) (param $to i32) (param $amount i64) (result i32)
  (local $collateralType i32)
  (local $collateral i64)
  (local $currentBalance i64)
  (local $newBalance i64)
  (local $newTotalSupply i64)
  local.get $this
  global.get $assembly/contract/Context.caller
  call $assembly/stablecoin/PaperDollar#isMinter
  i32.const 992
  call $assembly/contract/assert
  local.get $this
  call $assembly/stablecoin/PaperDollar#checkNotPaused
  local.get $this
  local.get $to
  call $assembly/stablecoin/PaperDollar#checkNotBlacklisted
  local.get $this
  local.get $to
  call $assembly/stablecoin/PaperDollar#getCollateralType
  local.set $collateralType
  local.get $this
  local.get $to
  local.get $collateralType
  call $assembly/stablecoin/PaperDollar#getCollateralBalance
  local.set $collateral
  local.get $this
  local.get $to
  call $assembly/stablecoin/PaperDollar#getBalance
  local.set $currentBalance
  local.get $currentBalance
  local.get $amount
  i64.add
  local.set $newBalance
  local.get $collateral
  i64.const 100
  i64.mul
  local.get $newBalance
  local.get $this
  call $assembly/stablecoin/PaperDollar#getCollateralRatio
  i64.mul
  i64.ge_u
  i32.const 2864
  call $assembly/contract/assert
  local.get $this
  call $assembly/stablecoin/PaperDollar#getTotalSupply
  local.get $amount
  i64.add
  local.set $newTotalSupply
  global.get $assembly/stablecoin/TOTAL_SUPPLY
  local.get $newTotalSupply
  call $assembly/contract/storage.set<u64>
  global.get $assembly/stablecoin/BALANCES_PREFIX
  local.get $to
  call $~lib/string/String.__concat
  local.get $newBalance
  call $assembly/contract/storage.set<u64>
  local.get $to
  local.get $amount
  local.get $collateralType
  call $assembly/stablecoin/emitMint
  i32.const 1
  return
 )
 (func $assembly/index/mint (param $to i32) (param $amount i64) (result i32)
  call $assembly/index/getContract
  local.get $to
  local.get $amount
  call $assembly/stablecoin/PaperDollar#mint
  return
 )
 (func $assembly/stablecoin/emitBurn (param $from i32) (param $amount i64)
 )
 (func $assembly/stablecoin/PaperDollar#burn (param $this i32) (param $amount i64) (result i32)
  (local $from i32)
  (local $fromBalance i64)
  global.get $assembly/contract/Context.caller
  local.set $from
  local.get $this
  call $assembly/stablecoin/PaperDollar#checkNotPaused
  local.get $this
  local.get $from
  call $assembly/stablecoin/PaperDollar#checkNotBlacklisted
  local.get $this
  local.get $from
  call $assembly/stablecoin/PaperDollar#getBalance
  local.set $fromBalance
  local.get $fromBalance
  local.get $amount
  i64.ge_u
  i32.const 928
  call $assembly/contract/assert
  global.get $assembly/stablecoin/BALANCES_PREFIX
  local.get $from
  call $~lib/string/String.__concat
  local.get $fromBalance
  local.get $amount
  i64.sub
  call $assembly/contract/storage.set<u64>
  global.get $assembly/stablecoin/TOTAL_SUPPLY
  local.get $this
  call $assembly/stablecoin/PaperDollar#getTotalSupply
  local.get $amount
  i64.sub
  call $assembly/contract/storage.set<u64>
  local.get $from
  local.get $amount
  call $assembly/stablecoin/emitBurn
  i32.const 1
  return
 )
 (func $assembly/index/burn (param $amount i64) (result i32)
  call $assembly/index/getContract
  local.get $amount
  call $assembly/stablecoin/PaperDollar#burn
  return
 )
 (func $~lib/util/string/compareImpl (param $str1 i32) (param $index1 i32) (param $str2 i32) (param $index2 i32) (param $len i32) (result i32)
  (local $ptr1 i32)
  (local $ptr2 i32)
  (local $7 i32)
  (local $a i32)
  (local $b i32)
  local.get $str1
  local.get $index1
  i32.const 1
  i32.shl
  i32.add
  local.set $ptr1
  local.get $str2
  local.get $index2
  i32.const 1
  i32.shl
  i32.add
  local.set $ptr2
  i32.const 0
  i32.const 2
  i32.lt_s
  drop
  local.get $len
  i32.const 4
  i32.ge_u
  if (result i32)
   local.get $ptr1
   i32.const 7
   i32.and
   local.get $ptr2
   i32.const 7
   i32.and
   i32.or
   i32.eqz
  else
   i32.const 0
  end
  if
   block $do-break|0
    loop $do-loop|0
     local.get $ptr1
     i64.load
     local.get $ptr2
     i64.load
     i64.ne
     if
      br $do-break|0
     end
     local.get $ptr1
     i32.const 8
     i32.add
     local.set $ptr1
     local.get $ptr2
     i32.const 8
     i32.add
     local.set $ptr2
     local.get $len
     i32.const 4
     i32.sub
     local.set $len
     local.get $len
     i32.const 4
     i32.ge_u
     br_if $do-loop|0
    end
   end
  end
  loop $while-continue|1
   local.get $len
   local.tee $7
   i32.const 1
   i32.sub
   local.set $len
   local.get $7
   if
    local.get $ptr1
    i32.load16_u
    local.set $a
    local.get $ptr2
    i32.load16_u
    local.set $b
    local.get $a
    local.get $b
    i32.ne
    if
     local.get $a
     local.get $b
     i32.sub
     return
    end
    local.get $ptr1
    i32.const 2
    i32.add
    local.set $ptr1
    local.get $ptr2
    i32.const 2
    i32.add
    local.set $ptr2
    br $while-continue|1
   end
  end
  i32.const 0
  return
 )
 (func $~lib/string/String.__eq (param $left i32) (param $right i32) (result i32)
  (local $leftLength i32)
  local.get $left
  local.get $right
  i32.eq
  if
   i32.const 1
   return
  end
  local.get $left
  i32.const 0
  i32.eq
  if (result i32)
   i32.const 1
  else
   local.get $right
   i32.const 0
   i32.eq
  end
  if
   i32.const 0
   return
  end
  local.get $left
  call $~lib/string/String#get:length
  local.set $leftLength
  local.get $leftLength
  local.get $right
  call $~lib/string/String#get:length
  i32.ne
  if
   i32.const 0
   return
  end
  local.get $left
  i32.const 0
  local.get $right
  i32.const 0
  local.get $leftLength
  call $~lib/util/string/compareImpl
  i32.eqz
  return
 )
 (func $assembly/stablecoin/PaperDollar#checkOwner (param $this i32)
  global.get $assembly/stablecoin/OWNER
  i32.const 32
  call $assembly/contract/storage.get<~lib/string/String>
  global.get $assembly/contract/Context.caller
  call $~lib/string/String.__eq
  i32.const 2944
  call $assembly/contract/assert
 )
 (func $assembly/stablecoin/PaperDollar#addMinter (param $this i32) (param $address i32) (result i32)
  local.get $this
  call $assembly/stablecoin/PaperDollar#checkOwner
  global.get $assembly/stablecoin/MINTER_PREFIX
  local.get $address
  call $~lib/string/String.__concat
  i32.const 1
  call $assembly/contract/storage.set<bool>
  i32.const 1
  return
 )
 (func $assembly/index/addMinter (param $address i32) (result i32)
  call $assembly/index/getContract
  local.get $address
  call $assembly/stablecoin/PaperDollar#addMinter
  return
 )
 (func $assembly/stablecoin/PaperDollar#removeMinter (param $this i32) (param $address i32) (result i32)
  local.get $this
  call $assembly/stablecoin/PaperDollar#checkOwner
  global.get $assembly/stablecoin/MINTER_PREFIX
  local.get $address
  call $~lib/string/String.__concat
  i32.const 0
  call $assembly/contract/storage.set<bool>
  i32.const 1
  return
 )
 (func $assembly/index/removeMinter (param $address i32) (result i32)
  call $assembly/index/getContract
  local.get $address
  call $assembly/stablecoin/PaperDollar#removeMinter
  return
 )
 (func $assembly/index/isMinter (param $address i32) (result i32)
  call $assembly/index/getContract
  local.get $address
  call $assembly/stablecoin/PaperDollar#isMinter
  return
 )
 (func $assembly/stablecoin/PaperDollar#blacklist (param $this i32) (param $address i32) (result i32)
  local.get $this
  call $assembly/stablecoin/PaperDollar#checkOwner
  global.get $assembly/stablecoin/BLACKLIST_PREFIX
  local.get $address
  call $~lib/string/String.__concat
  i32.const 1
  call $assembly/contract/storage.set<bool>
  i32.const 1
  return
 )
 (func $assembly/index/blacklist (param $address i32) (result i32)
  call $assembly/index/getContract
  local.get $address
  call $assembly/stablecoin/PaperDollar#blacklist
  return
 )
 (func $assembly/stablecoin/PaperDollar#unblacklist (param $this i32) (param $address i32) (result i32)
  local.get $this
  call $assembly/stablecoin/PaperDollar#checkOwner
  global.get $assembly/stablecoin/BLACKLIST_PREFIX
  local.get $address
  call $~lib/string/String.__concat
  i32.const 0
  call $assembly/contract/storage.set<bool>
  i32.const 1
  return
 )
 (func $assembly/index/unblacklist (param $address i32) (result i32)
  call $assembly/index/getContract
  local.get $address
  call $assembly/stablecoin/PaperDollar#unblacklist
  return
 )
 (func $assembly/index/isBlacklisted (param $address i32) (result i32)
  call $assembly/index/getContract
  local.get $address
  call $assembly/stablecoin/PaperDollar#isBlacklisted
  return
 )
 (func $assembly/stablecoin/PaperDollar#getFiatReserve (param $this i32) (result i64)
  global.get $assembly/stablecoin/FIAT_RESERVE
  i64.const 0
  call $assembly/contract/storage.get<u64>
  return
 )
 (func $assembly/contract/storage.set<i32> (param $key i32) (param $value i32)
 )
 (func $assembly/stablecoin/emitCollateralAdded (param $from i32) (param $amount i64) (param $collateralType i32)
 )
 (func $assembly/stablecoin/PaperDollar#addCollateral (param $this i32) (param $amount i64) (param $collateralType i32) (result i32)
  (local $from i32)
  (local $currentCollateral i64)
  global.get $assembly/contract/Context.caller
  local.set $from
  local.get $this
  call $assembly/stablecoin/PaperDollar#checkNotPaused
  local.get $this
  local.get $from
  call $assembly/stablecoin/PaperDollar#checkNotBlacklisted
  local.get $collateralType
  global.get $assembly/stablecoin/CollateralType.FIAT
  i32.eq
  if
   global.get $assembly/stablecoin/FIAT_RESERVE
   local.get $this
   call $assembly/stablecoin/PaperDollar#getFiatReserve
   local.get $amount
   i64.add
   call $assembly/contract/storage.set<u64>
  end
  local.get $this
  local.get $from
  local.get $collateralType
  call $assembly/stablecoin/PaperDollar#getCollateralBalance
  local.set $currentCollateral
  global.get $assembly/stablecoin/COLLATERAL_PREFIX
  local.get $collateralType
  i32.const 10
  call $~lib/number/I32#toString
  call $~lib/string/String.__concat
  i32.const 2832
  call $~lib/string/String.__concat
  local.get $from
  call $~lib/string/String.__concat
  local.get $currentCollateral
  local.get $amount
  i64.add
  call $assembly/contract/storage.set<u64>
  global.get $assembly/stablecoin/COLLATERAL_TYPE_PREFIX
  local.get $from
  call $~lib/string/String.__concat
  local.get $collateralType
  call $assembly/contract/storage.set<i32>
  local.get $from
  local.get $amount
  local.get $collateralType
  call $assembly/stablecoin/emitCollateralAdded
  i32.const 1
  return
 )
 (func $assembly/index/addCollateral (param $amount i64) (param $collateralType i32) (result i32)
  call $assembly/index/getContract
  local.get $amount
  local.get $collateralType
  call $assembly/stablecoin/PaperDollar#addCollateral
  return
 )
 (func $assembly/stablecoin/emitCollateralRemoved (param $to i32) (param $amount i64) (param $collateralType i32)
 )
 (func $assembly/stablecoin/PaperDollar#removeCollateral (param $this i32) (param $amount i64) (result i32)
  (local $from i32)
  (local $collateralType i32)
  (local $currentCollateral i64)
  (local $userBalance i64)
  (local $remainingCollateral i64)
  global.get $assembly/contract/Context.caller
  local.set $from
  local.get $this
  call $assembly/stablecoin/PaperDollar#checkNotPaused
  local.get $this
  local.get $from
  call $assembly/stablecoin/PaperDollar#checkNotBlacklisted
  local.get $this
  local.get $from
  call $assembly/stablecoin/PaperDollar#getCollateralType
  local.set $collateralType
  local.get $this
  local.get $from
  local.get $collateralType
  call $assembly/stablecoin/PaperDollar#getCollateralBalance
  local.set $currentCollateral
  local.get $currentCollateral
  local.get $amount
  i64.ge_u
  i32.const 2864
  call $assembly/contract/assert
  local.get $this
  local.get $from
  call $assembly/stablecoin/PaperDollar#getBalance
  local.set $userBalance
  local.get $currentCollateral
  local.get $amount
  i64.sub
  local.set $remainingCollateral
  local.get $remainingCollateral
  i64.const 100
  i64.mul
  local.get $userBalance
  local.get $this
  call $assembly/stablecoin/PaperDollar#getCollateralRatio
  i64.mul
  i64.ge_u
  i32.const 3040
  call $assembly/contract/assert
  local.get $collateralType
  global.get $assembly/stablecoin/CollateralType.FIAT
  i32.eq
  if
   global.get $assembly/stablecoin/FIAT_RESERVE
   local.get $this
   call $assembly/stablecoin/PaperDollar#getFiatReserve
   local.get $amount
   i64.sub
   call $assembly/contract/storage.set<u64>
  end
  global.get $assembly/stablecoin/COLLATERAL_PREFIX
  local.get $collateralType
  i32.const 10
  call $~lib/number/I32#toString
  call $~lib/string/String.__concat
  i32.const 2832
  call $~lib/string/String.__concat
  local.get $from
  call $~lib/string/String.__concat
  local.get $remainingCollateral
  call $assembly/contract/storage.set<u64>
  local.get $from
  local.get $amount
  local.get $collateralType
  call $assembly/stablecoin/emitCollateralRemoved
  i32.const 1
  return
 )
 (func $assembly/index/removeCollateral (param $amount i64) (result i32)
  call $assembly/index/getContract
  local.get $amount
  call $assembly/stablecoin/PaperDollar#removeCollateral
  return
 )
 (func $assembly/index/getCollateralBalance (param $address i32) (param $collateralType i32) (result i64)
  call $assembly/index/getContract
  local.get $address
  local.get $collateralType
  call $assembly/stablecoin/PaperDollar#getCollateralBalance
  return
 )
 (func $assembly/index/getCollateralType (param $address i32) (result i32)
  call $assembly/index/getContract
  local.get $address
  call $assembly/stablecoin/PaperDollar#getCollateralType
  return
 )
 (func $assembly/stablecoin/PaperDollar#setCollateralRatio (param $this i32) (param $ratio i64) (result i32)
  local.get $this
  call $assembly/stablecoin/PaperDollar#checkOwner
  local.get $ratio
  i64.const 100
  i64.ge_u
  i32.const 3136
  call $assembly/contract/assert
  global.get $assembly/stablecoin/COLLATERAL_RATIO
  local.get $ratio
  call $assembly/contract/storage.set<u64>
  i32.const 1
  return
 )
 (func $assembly/index/setCollateralRatio (param $ratio i64) (result i32)
  call $assembly/index/getContract
  local.get $ratio
  call $assembly/stablecoin/PaperDollar#setCollateralRatio
  return
 )
 (func $assembly/index/getCollateralRatio (result i64)
  call $assembly/index/getContract
  call $assembly/stablecoin/PaperDollar#getCollateralRatio
  return
 )
 (func $assembly/stablecoin/PaperDollar#setBinomTokenAddress (param $this i32) (param $address i32) (result i32)
  local.get $this
  call $assembly/stablecoin/PaperDollar#checkOwner
  global.get $assembly/stablecoin/BINOM_TOKEN_ADDRESS
  local.get $address
  call $assembly/contract/storage.set<~lib/string/String>
  i32.const 1
  return
 )
 (func $assembly/index/setBinomTokenAddress (param $address i32) (result i32)
  call $assembly/index/getContract
  local.get $address
  call $assembly/stablecoin/PaperDollar#setBinomTokenAddress
  return
 )
 (func $assembly/index/getFiatReserve (result i64)
  call $assembly/index/getContract
  call $assembly/stablecoin/PaperDollar#getFiatReserve
  return
 )
 (func $assembly/stablecoin/PaperDollar#pause (param $this i32) (result i32)
  local.get $this
  call $assembly/stablecoin/PaperDollar#checkOwner
  global.get $assembly/stablecoin/PAUSED
  i32.const 1
  call $assembly/contract/storage.set<bool>
  i32.const 1
  return
 )
 (func $assembly/index/pause (result i32)
  call $assembly/index/getContract
  call $assembly/stablecoin/PaperDollar#pause
  return
 )
 (func $assembly/stablecoin/PaperDollar#unpause (param $this i32) (result i32)
  local.get $this
  call $assembly/stablecoin/PaperDollar#checkOwner
  global.get $assembly/stablecoin/PAUSED
  i32.const 0
  call $assembly/contract/storage.set<bool>
  i32.const 1
  return
 )
 (func $assembly/index/unpause (result i32)
  call $assembly/index/getContract
  call $assembly/stablecoin/PaperDollar#unpause
  return
 )
 (func $assembly/stablecoin/PaperDollar#transferOwnership (param $this i32) (param $newOwner i32) (result i32)
  local.get $this
  call $assembly/stablecoin/PaperDollar#checkOwner
  global.get $assembly/stablecoin/OWNER
  local.get $newOwner
  call $assembly/contract/storage.set<~lib/string/String>
  i32.const 1
  return
 )
 (func $assembly/index/transferOwnership (param $newOwner i32) (result i32)
  call $assembly/index/getContract
  local.get $newOwner
  call $assembly/stablecoin/PaperDollar#transferOwnership
  return
 )
 (func $assembly/index/add (param $a i32) (param $b i32) (result i32)
  local.get $a
  local.get $b
  i32.add
  return
 )
 (func $assembly/index/multiply (param $a i32) (param $b i32) (result i32)
  local.get $a
  local.get $b
  i32.mul
  return
 )
 (func $assembly/index/store (param $value i32)
  local.get $value
  global.set $assembly/index/storedValue
 )
 (func $assembly/index/retrieve (result i32)
  global.get $assembly/index/storedValue
  return
 )
 (func $~lib/rt/stub/__pin (param $ptr i32) (result i32)
  local.get $ptr
  return
 )
 (func $~lib/rt/stub/__unpin (param $ptr i32)
 )
 (func $~lib/rt/stub/__collect
 )
 (func $~start
  call $start:assembly/index
 )
)
