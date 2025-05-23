(module
 (type $0 (func (param i32 i32)))
 (type $1 (func (param i32) (result i32)))
 (type $2 (func (param i32 i32) (result i32)))
 (type $3 (func (param i32)))
 (type $4 (func (param i32 i64 i32)))
 (type $5 (func (param i32 i64)))
 (type $6 (func (result i32)))
 (type $7 (func))
 (type $8 (func (param i32 i32 i32 i32)))
 (type $9 (func (param i32 i32 i64)))
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
 (global $assembly/index/storedValue (mut i32) (i32.const 0))
 (global $~lib/shared/runtime/Runtime.Stub i32 (i32.const 0))
 (global $~lib/shared/runtime/Runtime.Minimal i32 (i32.const 1))
 (global $~lib/shared/runtime/Runtime.Incremental i32 (i32.const 2))
 (global $~lib/rt/stub/startOffset (mut i32) (i32.const 0))
 (global $~lib/rt/stub/offset (mut i32) (i32.const 0))
 (global $~lib/rt/__rtti_base i32 (i32.const 720))
 (global $~lib/memory/__heap_base i32 (i32.const 744))
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
 (data $14 (i32.const 720) "\05\00\00\00 \00\00\00 \00\00\00 \00\00\00\00\00\00\00 \00\00\00")
 (table $0 1 1 funcref)
 (elem $0 (i32.const 1))
 (export "createPaperDollar" (func $assembly/index/createPaperDollar))
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
 (func $assembly/index/createPaperDollar (result i32)
  i32.const 0
  call $assembly/stablecoin/PaperDollar#constructor
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
 (func $assembly/stablecoin/emitTransfer (param $from i32) (param $to i32) (param $amount i64)
 )
 (func $assembly/stablecoin/emitMint (param $to i32) (param $amount i64) (param $collateralType i32)
 )
 (func $assembly/stablecoin/emitBurn (param $from i32) (param $amount i64)
 )
 (func $assembly/stablecoin/emitCollateralAdded (param $from i32) (param $amount i64) (param $collateralType i32)
 )
 (func $assembly/stablecoin/emitCollateralRemoved (param $to i32) (param $amount i64) (param $collateralType i32)
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
 )
)
