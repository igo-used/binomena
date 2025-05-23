(module
 (type $0 (func (param i32 i32) (result i32)))
 (type $1 (func (result i32)))
 (type $2 (func (param i32)))
 (type $3 (func))
 (type $4 (func (param i32 i64)))
 (type $5 (func (param i32 i32 i32 i32)))
 (type $6 (func (param i32 i32 i64)))
 (type $7 (func (param i32 i64 i32)))
 (type $8 (func (param i32) (result i32)))
 (import "env" "memory" (memory $0 1))
 (import "env" "abort" (func $~lib/builtins/abort (param i32 i32 i32 i32)))
 (global $assembly/index/storedValue (mut i32) (i32.const 0))
 (global $~lib/rt/stub/offset (mut i32) (i32.const 0))
 (global $~lib/rt/__rtti_base i32 (i32.const 1744))
 (data $0 (i32.const 1036) "\1c\00\00\00\00\00\00\00\00\00\00\00\02\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00")
 (data $1 (i32.const 1068) "\1c\00\00\00\00\00\00\00\00\00\00\00\02\00\00\00\n\00\00\00o\00w\00n\00e\00r\00\00\00")
 (data $2 (i32.const 1100) ",\00\00\00\00\00\00\00\00\00\00\00\02\00\00\00\16\00\00\00t\00o\00t\00a\00l\00S\00u\00p\00p\00l\00y\00\00\00\00\00\00\00")
 (data $3 (i32.const 1148) ",\00\00\00\00\00\00\00\00\00\00\00\02\00\00\00\10\00\00\00b\00a\00l\00a\00n\00c\00e\00_\00\00\00\00\00\00\00\00\00\00\00\00\00")
 (data $4 (i32.const 1196) "\1c\00\00\00\00\00\00\00\00\00\00\00\02\00\00\00\0c\00\00\00p\00a\00u\00s\00e\00d\00")
 (data $5 (i32.const 1228) ",\00\00\00\00\00\00\00\00\00\00\00\02\00\00\00\14\00\00\00b\00l\00a\00c\00k\00l\00i\00s\00t\00_\00\00\00\00\00\00\00\00\00")
 (data $6 (i32.const 1276) ",\00\00\00\00\00\00\00\00\00\00\00\02\00\00\00\0e\00\00\00m\00i\00n\00t\00e\00r\00_\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00")
 (data $7 (i32.const 1324) ",\00\00\00\00\00\00\00\00\00\00\00\02\00\00\00\16\00\00\00c\00o\00l\00l\00a\00t\00e\00r\00a\00l\00_\00\00\00\00\00\00\00")
 (data $8 (i32.const 1372) "<\00\00\00\00\00\00\00\00\00\00\00\02\00\00\00 \00\00\00c\00o\00l\00l\00a\00t\00e\00r\00a\00l\00_\00t\00y\00p\00e\00_\00\00\00\00\00\00\00\00\00\00\00\00\00")
 (data $9 (i32.const 1436) "<\00\00\00\00\00\00\00\00\00\00\00\02\00\00\00&\00\00\00b\00i\00n\00o\00m\00_\00t\00o\00k\00e\00n\00_\00a\00d\00d\00r\00e\00s\00s\00\00\00\00\00\00\00")
 (data $10 (i32.const 1500) "<\00\00\00\00\00\00\00\00\00\00\00\02\00\00\00 \00\00\00c\00o\00l\00l\00a\00t\00e\00r\00a\00l\00_\00r\00a\00t\00i\00o\00\00\00\00\00\00\00\00\00\00\00\00\00")
 (data $11 (i32.const 1564) ",\00\00\00\00\00\00\00\00\00\00\00\02\00\00\00\18\00\00\00f\00i\00a\00t\00_\00r\00e\00s\00e\00r\00v\00e\00\00\00\00\00")
 (data $12 (i32.const 1612) "<\00\00\00\00\00\00\00\00\00\00\00\02\00\00\00(\00\00\00A\00l\00l\00o\00c\00a\00t\00i\00o\00n\00 \00t\00o\00o\00 \00l\00a\00r\00g\00e\00\00\00\00\00")
 (data $13 (i32.const 1676) "<\00\00\00\00\00\00\00\00\00\00\00\02\00\00\00\1e\00\00\00~\00l\00i\00b\00/\00r\00t\00/\00s\00t\00u\00b\00.\00t\00s\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00")
 (data $14 (i32.const 1744) "\05\00\00\00 \00\00\00 \00\00\00 \00\00\00\00\00\00\00 \00\00\00")
 (export "createPaperDollar" (func $assembly/index/createPaperDollar))
 (export "add" (func $assembly/index/add))
 (export "multiply" (func $assembly/index/multiply))
 (export "store" (func $assembly/index/store))
 (export "retrieve" (func $assembly/index/retrieve))
 (export "emitTransfer" (func $assembly/stablecoin/emitTransfer))
 (export "emitMint" (func $assembly/stablecoin/emitMint))
 (export "emitBurn" (func $assembly/contract/storage.set<u64>))
 (export "emitCollateralAdded" (func $assembly/stablecoin/emitMint))
 (export "emitCollateralRemoved" (func $assembly/stablecoin/emitMint))
 (export "__new" (func $~lib/rt/stub/__new))
 (export "__pin" (func $~lib/rt/stub/__pin))
 (export "__unpin" (func $~lib/rt/stub/__unpin))
 (export "__collect" (func $~lib/rt/stub/__collect))
 (export "__rtti_base" (global $~lib/rt/__rtti_base))
 (export "memory" (memory $0))
 (start $~start)
 (func $assembly/contract/storage.set<u64> (param $0 i32) (param $1 i64)
 )
 (func $~lib/rt/stub/__new (param $0 i32) (param $1 i32) (result i32)
  (local $2 i32)
  (local $3 i32)
  (local $4 i32)
  (local $5 i32)
  (local $6 i32)
  (local $7 i32)
  local.get $0
  i32.const 1073741804
  i32.gt_u
  if
   i32.const 1632
   i32.const 1696
   i32.const 86
   i32.const 30
   call $~lib/builtins/abort
   unreachable
  end
  local.get $0
  i32.const 16
  i32.add
  local.tee $4
  i32.const 1073741820
  i32.gt_u
  if
   i32.const 1632
   i32.const 1696
   i32.const 33
   i32.const 29
   call $~lib/builtins/abort
   unreachable
  end
  global.get $~lib/rt/stub/offset
  local.set $3
  global.get $~lib/rt/stub/offset
  i32.const 4
  i32.add
  local.tee $2
  local.get $4
  i32.const 19
  i32.add
  i32.const -16
  i32.and
  i32.const 4
  i32.sub
  local.tee $4
  i32.add
  local.tee $5
  memory.size
  local.tee $6
  i32.const 16
  i32.shl
  i32.const 15
  i32.add
  i32.const -16
  i32.and
  local.tee $7
  i32.gt_u
  if
   local.get $6
   local.get $5
   local.get $7
   i32.sub
   i32.const 65535
   i32.add
   i32.const -65536
   i32.and
   i32.const 16
   i32.shr_u
   local.tee $7
   local.get $6
   local.get $7
   i32.gt_s
   select
   memory.grow
   i32.const 0
   i32.lt_s
   if
    local.get $7
    memory.grow
    i32.const 0
    i32.lt_s
    if
     unreachable
    end
   end
  end
  local.get $5
  global.set $~lib/rt/stub/offset
  local.get $3
  local.get $4
  i32.store
  local.get $2
  i32.const 4
  i32.sub
  local.tee $3
  i32.const 0
  i32.store offset=4
  local.get $3
  i32.const 0
  i32.store offset=8
  local.get $3
  local.get $1
  i32.store offset=12
  local.get $3
  local.get $0
  i32.store offset=16
  local.get $2
  i32.const 16
  i32.add
 )
 (func $assembly/index/createPaperDollar (result i32)
  (local $0 i32)
  i32.const 0
  i32.const 4
  call $~lib/rt/stub/__new
  local.set $0
  i32.const 1052
  i32.load
  drop
  local.get $0
 )
 (func $assembly/index/add (param $0 i32) (param $1 i32) (result i32)
  local.get $0
  local.get $1
  i32.add
 )
 (func $assembly/index/multiply (param $0 i32) (param $1 i32) (result i32)
  local.get $0
  local.get $1
  i32.mul
 )
 (func $assembly/index/store (param $0 i32)
  local.get $0
  global.set $assembly/index/storedValue
 )
 (func $assembly/index/retrieve (result i32)
  global.get $assembly/index/storedValue
 )
 (func $assembly/stablecoin/emitTransfer (param $0 i32) (param $1 i32) (param $2 i64)
 )
 (func $assembly/stablecoin/emitMint (param $0 i32) (param $1 i64) (param $2 i32)
 )
 (func $~lib/rt/stub/__pin (param $0 i32) (result i32)
  local.get $0
 )
 (func $~lib/rt/stub/__unpin (param $0 i32)
 )
 (func $~lib/rt/stub/__collect
 )
 (func $~start
  i32.const 1772
  global.set $~lib/rt/stub/offset
 )
)
