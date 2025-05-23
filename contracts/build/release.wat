(module
 (type $0 (func (param i32) (result i32)))
 (type $1 (func (param i32 i32) (result i32)))
 (type $2 (func (result i32)))
 (type $3 (func (param i64) (result i32)))
 (type $4 (func (result i64)))
 (type $5 (func (param i32 i64) (result i32)))
 (type $6 (func (param i32)))
 (type $7 (func))
 (type $8 (func (param i32 i64)))
 (type $9 (func (param i32 i32 i32 i32)))
 (type $10 (func (param i32) (result i64)))
 (type $11 (func (param i32 i32 i64)))
 (type $12 (func (param i32 i64 i32)))
 (type $13 (func (param i64 i32) (result i32)))
 (type $14 (func (param i32 i32) (result i64)))
 (import "env" "memory" (memory $0 1))
 (import "env" "abort" (func $~lib/builtins/abort (param i32 i32 i32 i32)))
 (global $~lib/rt/stub/offset (mut i32) (i32.const 0))
 (global $assembly/index/contractInstance (mut i32) (i32.const 0))
 (global $assembly/index/storedValue (mut i32) (i32.const 0))
 (global $~lib/rt/__rtti_base i32 (i32.const 4240))
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
 (data $14 (i32.const 1740) "<\00\00\00\00\00\00\00\00\00\00\00\02\00\00\00$\00\00\00C\00o\00n\00t\00r\00a\00c\00t\00 \00i\00s\00 \00p\00a\00u\00s\00e\00d\00\00\00\00\00\00\00\00\00")
 (data $15 (i32.const 1804) "<\00\00\00\00\00\00\00\00\00\00\00\02\00\00\00(\00\00\00a\00s\00s\00e\00m\00b\00l\00y\00/\00c\00o\00n\00t\00r\00a\00c\00t\00.\00t\00s\00\00\00\00\00")
 (data $16 (i32.const 1868) "<\00\00\00\00\00\00\00\00\00\00\00\02\00\00\00,\00\00\00A\00d\00d\00r\00e\00s\00s\00 \00i\00s\00 \00b\00l\00a\00c\00k\00l\00i\00s\00t\00e\00d\00")
 (data $17 (i32.const 1932) "<\00\00\00\00\00\00\00\00\00\00\00\02\00\00\00(\00\00\00I\00n\00s\00u\00f\00f\00i\00c\00i\00e\00n\00t\00 \00b\00a\00l\00a\00n\00c\00e\00\00\00\00\00")
 (data $18 (i32.const 1996) "<\00\00\00\00\00\00\00\00\00\00\00\02\00\00\00*\00\00\00O\00n\00l\00y\00 \00m\00i\00n\00t\00e\00r\00s\00 \00c\00a\00n\00 \00m\00i\00n\00t\00\00\00")
 (data $19 (i32.const 2060) "|\00\00\00\00\00\00\00\00\00\00\00\02\00\00\00d\00\00\00t\00o\00S\00t\00r\00i\00n\00g\00(\00)\00 \00r\00a\00d\00i\00x\00 \00a\00r\00g\00u\00m\00e\00n\00t\00 \00m\00u\00s\00t\00 \00b\00e\00 \00b\00e\00t\00w\00e\00e\00n\00 \002\00 \00a\00n\00d\00 \003\006\00\00\00\00\00\00\00\00\00")
 (data $20 (i32.const 2188) "<\00\00\00\00\00\00\00\00\00\00\00\02\00\00\00&\00\00\00~\00l\00i\00b\00/\00u\00t\00i\00l\00/\00n\00u\00m\00b\00e\00r\00.\00t\00s\00\00\00\00\00\00\00")
 (data $21 (i32.const 2252) "\1c\00\00\00\00\00\00\00\00\00\00\00\02\00\00\00\02\00\00\000\00\00\00\00\00\00\00\00\00\00\00")
 (data $22 (i32.const 2284) "0\000\000\001\000\002\000\003\000\004\000\005\000\006\000\007\000\008\000\009\001\000\001\001\001\002\001\003\001\004\001\005\001\006\001\007\001\008\001\009\002\000\002\001\002\002\002\003\002\004\002\005\002\006\002\007\002\008\002\009\003\000\003\001\003\002\003\003\003\004\003\005\003\006\003\007\003\008\003\009\004\000\004\001\004\002\004\003\004\004\004\005\004\006\004\007\004\008\004\009\005\000\005\001\005\002\005\003\005\004\005\005\005\006\005\007\005\008\005\009\006\000\006\001\006\002\006\003\006\004\006\005\006\006\006\007\006\008\006\009\007\000\007\001\007\002\007\003\007\004\007\005\007\006\007\007\007\008\007\009\008\000\008\001\008\002\008\003\008\004\008\005\008\006\008\007\008\008\008\009\009\000\009\001\009\002\009\003\009\004\009\005\009\006\009\007\009\008\009\009\00")
 (data $23 (i32.const 2684) "\1c\04\00\00\00\00\00\00\00\00\00\00\02\00\00\00\00\04\00\000\000\000\001\000\002\000\003\000\004\000\005\000\006\000\007\000\008\000\009\000\00a\000\00b\000\00c\000\00d\000\00e\000\00f\001\000\001\001\001\002\001\003\001\004\001\005\001\006\001\007\001\008\001\009\001\00a\001\00b\001\00c\001\00d\001\00e\001\00f\002\000\002\001\002\002\002\003\002\004\002\005\002\006\002\007\002\008\002\009\002\00a\002\00b\002\00c\002\00d\002\00e\002\00f\003\000\003\001\003\002\003\003\003\004\003\005\003\006\003\007\003\008\003\009\003\00a\003\00b\003\00c\003\00d\003\00e\003\00f\004\000\004\001\004\002\004\003\004\004\004\005\004\006\004\007\004\008\004\009\004\00a\004\00b\004\00c\004\00d\004\00e\004\00f\005\000\005\001\005\002\005\003\005\004\005\005\005\006\005\007\005\008\005\009\005\00a\005\00b\005\00c\005\00d\005\00e\005\00f\006\000\006\001\006\002\006\003\006\004\006\005\006\006\006\007\006\008\006\009\006\00a\006\00b\006\00c\006\00d\006\00e\006\00f\007\000\007\001\007\002\007\003\007\004\007\005\007\006\007\007\007\008\007\009\007\00a\007\00b\007\00c\007\00d\007\00e\007\00f\008\000\008\001\008\002\008\003\008\004\008\005\008\006\008\007\008\008\008\009\008\00a\008\00b\008\00c\008\00d\008\00e\008\00f\009\000\009\001\009\002\009\003\009\004\009\005\009\006\009\007\009\008\009\009\009\00a\009\00b\009\00c\009\00d\009\00e\009\00f\00a\000\00a\001\00a\002\00a\003\00a\004\00a\005\00a\006\00a\007\00a\008\00a\009\00a\00a\00a\00b\00a\00c\00a\00d\00a\00e\00a\00f\00b\000\00b\001\00b\002\00b\003\00b\004\00b\005\00b\006\00b\007\00b\008\00b\009\00b\00a\00b\00b\00b\00c\00b\00d\00b\00e\00b\00f\00c\000\00c\001\00c\002\00c\003\00c\004\00c\005\00c\006\00c\007\00c\008\00c\009\00c\00a\00c\00b\00c\00c\00c\00d\00c\00e\00c\00f\00d\000\00d\001\00d\002\00d\003\00d\004\00d\005\00d\006\00d\007\00d\008\00d\009\00d\00a\00d\00b\00d\00c\00d\00d\00d\00e\00d\00f\00e\000\00e\001\00e\002\00e\003\00e\004\00e\005\00e\006\00e\007\00e\008\00e\009\00e\00a\00e\00b\00e\00c\00e\00d\00e\00e\00e\00f\00f\000\00f\001\00f\002\00f\003\00f\004\00f\005\00f\006\00f\007\00f\008\00f\009\00f\00a\00f\00b\00f\00c\00f\00d\00f\00e\00f\00f\00\00\00\00\00\00\00\00\00\00\00\00\00")
 (data $24 (i32.const 3740) "\\\00\00\00\00\00\00\00\00\00\00\00\02\00\00\00H\00\00\000\001\002\003\004\005\006\007\008\009\00a\00b\00c\00d\00e\00f\00g\00h\00i\00j\00k\00l\00m\00n\00o\00p\00q\00r\00s\00t\00u\00v\00w\00x\00y\00z\00\00\00\00\00")
 (data $25 (i32.const 3836) "\1c\00\00\00\00\00\00\00\00\00\00\00\02\00\00\00\02\00\00\00_\00\00\00\00\00\00\00\00\00\00\00")
 (data $26 (i32.const 3868) "L\00\00\00\00\00\00\00\00\00\00\00\02\00\00\00.\00\00\00I\00n\00s\00u\00f\00f\00i\00c\00i\00e\00n\00t\00 \00c\00o\00l\00l\00a\00t\00e\00r\00a\00l\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00")
 (data $27 (i32.const 3948) "\\\00\00\00\00\00\00\00\00\00\00\00\02\00\00\00B\00\00\00O\00n\00l\00y\00 \00o\00w\00n\00e\00r\00 \00c\00a\00n\00 \00c\00a\00l\00l\00 \00t\00h\00i\00s\00 \00f\00u\00n\00c\00t\00i\00o\00n\00\00\00\00\00\00\00\00\00\00\00")
 (data $28 (i32.const 4044) "\\\00\00\00\00\00\00\00\00\00\00\00\02\00\00\00B\00\00\00C\00o\00l\00l\00a\00t\00e\00r\00a\00l\00 \00r\00a\00t\00i\00o\00 \00w\00o\00u\00l\00d\00 \00b\00e\00 \00t\00o\00o\00 \00l\00o\00w\00\00\00\00\00\00\00\00\00\00\00")
 (data $29 (i32.const 4140) "\\\00\00\00\00\00\00\00\00\00\00\00\02\00\00\00L\00\00\00C\00o\00l\00l\00a\00t\00e\00r\00a\00l\00 \00r\00a\00t\00i\00o\00 \00m\00u\00s\00t\00 \00b\00e\00 \00a\00t\00 \00l\00e\00a\00s\00t\00 \001\000\000\00%\00")
 (data $30 (i32.const 4240) "\05\00\00\00 \00\00\00 \00\00\00 \00\00\00\00\00\00\00 \00\00\00")
 (export "createPaperDollar" (func $assembly/index/createPaperDollar))
 (export "getBalance" (func $assembly/index/getBalance))
 (export "getTotalSupply" (func $assembly/index/getTotalSupply))
 (export "transfer" (func $assembly/index/transfer))
 (export "mint" (func $assembly/index/mint))
 (export "burn" (func $assembly/index/burn))
 (export "addMinter" (func $assembly/index/addMinter))
 (export "removeMinter" (func $assembly/index/addMinter))
 (export "isMinter" (func $assembly/index/isMinter))
 (export "blacklist" (func $assembly/index/blacklist))
 (export "unblacklist" (func $assembly/index/blacklist))
 (export "isBlacklisted" (func $assembly/index/isBlacklisted))
 (export "addCollateral" (func $assembly/index/addCollateral))
 (export "removeCollateral" (func $assembly/index/removeCollateral))
 (export "getCollateralBalance" (func $assembly/index/getCollateralBalance))
 (export "getCollateralType" (func $assembly/index/getCollateralType))
 (export "setCollateralRatio" (func $assembly/index/setCollateralRatio))
 (export "getCollateralRatio" (func $assembly/index/getCollateralRatio))
 (export "setBinomTokenAddress" (func $assembly/index/setBinomTokenAddress))
 (export "getFiatReserve" (func $assembly/index/getTotalSupply))
 (export "pause" (func $assembly/index/pause))
 (export "unpause" (func $assembly/index/pause))
 (export "transferOwnership" (func $assembly/index/setBinomTokenAddress))
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
  global.get $assembly/index/contractInstance
 )
 (func $~lib/string/String.__concat (param $0 i32) (param $1 i32) (result i32)
  (local $2 i32)
  (local $3 i32)
  (local $4 i32)
  (local $5 i32)
  i32.const 1056
  local.set $2
  local.get $0
  i32.const 20
  i32.sub
  i32.load offset=16
  i32.const -2
  i32.and
  local.tee $3
  local.get $1
  i32.const 20
  i32.sub
  i32.load offset=16
  i32.const -2
  i32.and
  local.tee $4
  i32.add
  local.tee $5
  if
   local.get $5
   i32.const 2
   call $~lib/rt/stub/__new
   local.tee $2
   local.get $0
   local.get $3
   memory.copy
   local.get $2
   local.get $3
   i32.add
   local.get $1
   local.get $4
   memory.copy
  end
  local.get $2
 )
 (func $assembly/index/getBalance (param $0 i32) (result i64)
  i32.const 1168
  local.get $0
  call $~lib/string/String.__concat
  drop
  i64.const 0
 )
 (func $assembly/index/getTotalSupply (result i64)
  i64.const 0
 )
 (func $assembly/stablecoin/emitTransfer (param $0 i32) (param $1 i32) (param $2 i64)
 )
 (func $assembly/index/transfer (param $0 i32) (param $1 i64) (result i32)
  i32.const 1248
  i32.const 1056
  call $~lib/string/String.__concat
  drop
  i32.const 1248
  local.get $0
  call $~lib/string/String.__concat
  drop
  i32.const 1168
  i32.const 1056
  call $~lib/string/String.__concat
  drop
  local.get $1
  i64.const 0
  i64.ne
  if
   i32.const 1952
   i32.const 1824
   i32.const 123
   i32.const 9
   call $~lib/builtins/abort
   unreachable
  end
  i32.const 1168
  i32.const 1056
  call $~lib/string/String.__concat
  drop
  i32.const 1168
  local.get $0
  call $~lib/string/String.__concat
  drop
  i32.const 1168
  local.get $0
  call $~lib/string/String.__concat
  drop
  i32.const 1
 )
 (func $~lib/util/number/itoa32 (param $0 i32) (result i32)
  (local $1 i32)
  (local $2 i32)
  (local $3 i32)
  (local $4 i32)
  (local $5 i32)
  local.get $0
  i32.eqz
  if
   i32.const 2272
   return
  end
  i32.const 0
  local.get $0
  i32.sub
  local.get $0
  local.get $0
  i32.const 31
  i32.shr_u
  i32.const 1
  i32.shl
  local.tee $4
  select
  local.tee $1
  i32.const 100000
  i32.lt_u
  if (result i32)
   local.get $1
   i32.const 100
   i32.lt_u
   if (result i32)
    local.get $1
    i32.const 10
    i32.ge_u
    i32.const 1
    i32.add
   else
    local.get $1
    i32.const 10000
    i32.ge_u
    i32.const 3
    i32.add
    local.get $1
    i32.const 1000
    i32.ge_u
    i32.add
   end
  else
   local.get $1
   i32.const 10000000
   i32.lt_u
   if (result i32)
    local.get $1
    i32.const 1000000
    i32.ge_u
    i32.const 6
    i32.add
   else
    local.get $1
    i32.const 1000000000
    i32.ge_u
    i32.const 8
    i32.add
    local.get $1
    i32.const 100000000
    i32.ge_u
    i32.add
   end
  end
  local.tee $0
  i32.const 1
  i32.shl
  local.get $4
  i32.add
  i32.const 2
  call $~lib/rt/stub/__new
  local.tee $3
  local.get $4
  i32.add
  local.set $5
  loop $while-continue|0
   local.get $1
   i32.const 10000
   i32.ge_u
   if
    local.get $1
    i32.const 10000
    i32.rem_u
    local.set $2
    local.get $1
    i32.const 10000
    i32.div_u
    local.set $1
    local.get $5
    local.get $0
    i32.const 4
    i32.sub
    local.tee $0
    i32.const 1
    i32.shl
    i32.add
    local.get $2
    i32.const 100
    i32.div_u
    i32.const 2
    i32.shl
    i32.const 2284
    i32.add
    i64.load32_u
    local.get $2
    i32.const 100
    i32.rem_u
    i32.const 2
    i32.shl
    i32.const 2284
    i32.add
    i64.load32_u
    i64.const 32
    i64.shl
    i64.or
    i64.store
    br $while-continue|0
   end
  end
  local.get $1
  i32.const 100
  i32.ge_u
  if
   local.get $5
   local.get $0
   i32.const 2
   i32.sub
   local.tee $0
   i32.const 1
   i32.shl
   i32.add
   local.get $1
   i32.const 100
   i32.rem_u
   i32.const 2
   i32.shl
   i32.const 2284
   i32.add
   i32.load
   i32.store
   local.get $1
   i32.const 100
   i32.div_u
   local.set $1
  end
  local.get $1
  i32.const 10
  i32.ge_u
  if
   local.get $5
   local.get $0
   i32.const 2
   i32.sub
   i32.const 1
   i32.shl
   i32.add
   local.get $1
   i32.const 2
   i32.shl
   i32.const 2284
   i32.add
   i32.load
   i32.store
  else
   local.get $5
   local.get $0
   i32.const 1
   i32.sub
   i32.const 1
   i32.shl
   i32.add
   local.get $1
   i32.const 48
   i32.add
   i32.store16
  end
  local.get $4
  if
   local.get $3
   i32.const 45
   i32.store16
  end
  local.get $3
 )
 (func $assembly/stablecoin/emitMint (param $0 i32) (param $1 i64) (param $2 i32)
 )
 (func $assembly/index/mint (param $0 i32) (param $1 i64) (result i32)
  i32.const 1296
  i32.const 1056
  call $~lib/string/String.__concat
  drop
  i32.const 2016
  i32.const 1824
  i32.const 123
  i32.const 9
  call $~lib/builtins/abort
  unreachable
 )
 (func $assembly/index/burn (param $0 i64) (result i32)
  i32.const 1248
  i32.const 1056
  call $~lib/string/String.__concat
  drop
  i32.const 1168
  i32.const 1056
  call $~lib/string/String.__concat
  drop
  local.get $0
  i64.const 0
  i64.ne
  if
   i32.const 1952
   i32.const 1824
   i32.const 123
   i32.const 9
   call $~lib/builtins/abort
   unreachable
  end
  i32.const 1168
  i32.const 1056
  call $~lib/string/String.__concat
  drop
  i32.const 1
 )
 (func $assembly/index/addMinter (param $0 i32) (result i32)
  i32.const 1296
  local.get $0
  call $~lib/string/String.__concat
  drop
  i32.const 1
 )
 (func $assembly/index/isMinter (param $0 i32) (result i32)
  i32.const 1296
  local.get $0
  call $~lib/string/String.__concat
  drop
  i32.const 0
 )
 (func $assembly/index/blacklist (param $0 i32) (result i32)
  i32.const 1248
  local.get $0
  call $~lib/string/String.__concat
  drop
  i32.const 1
 )
 (func $assembly/index/isBlacklisted (param $0 i32) (result i32)
  i32.const 1248
  local.get $0
  call $~lib/string/String.__concat
  drop
  i32.const 0
 )
 (func $assembly/index/addCollateral (param $0 i64) (param $1 i32) (result i32)
  i32.const 1248
  i32.const 1056
  call $~lib/string/String.__concat
  drop
  i32.const 1344
  local.get $1
  call $~lib/util/number/itoa32
  call $~lib/string/String.__concat
  i32.const 3856
  call $~lib/string/String.__concat
  i32.const 1056
  call $~lib/string/String.__concat
  drop
  i32.const 1344
  local.get $1
  call $~lib/util/number/itoa32
  call $~lib/string/String.__concat
  i32.const 3856
  call $~lib/string/String.__concat
  i32.const 1056
  call $~lib/string/String.__concat
  drop
  i32.const 1392
  i32.const 1056
  call $~lib/string/String.__concat
  drop
  i32.const 1
 )
 (func $assembly/index/removeCollateral (param $0 i64) (result i32)
  i32.const 1248
  i32.const 1056
  call $~lib/string/String.__concat
  drop
  i32.const 1392
  i32.const 1056
  call $~lib/string/String.__concat
  drop
  i32.const 1344
  i32.const 0
  call $~lib/util/number/itoa32
  call $~lib/string/String.__concat
  i32.const 3856
  call $~lib/string/String.__concat
  i32.const 1056
  call $~lib/string/String.__concat
  drop
  local.get $0
  i64.eqz
  if
   i32.const 1168
   i32.const 1056
   call $~lib/string/String.__concat
   drop
   i32.const 1344
   i32.const 0
   call $~lib/util/number/itoa32
   call $~lib/string/String.__concat
   i32.const 3856
   call $~lib/string/String.__concat
   i32.const 1056
   call $~lib/string/String.__concat
   drop
   i32.const 1
   return
  end
  i32.const 3888
  i32.const 1824
  i32.const 123
  i32.const 9
  call $~lib/builtins/abort
  unreachable
 )
 (func $assembly/index/getCollateralBalance (param $0 i32) (param $1 i32) (result i64)
  i32.const 1344
  local.get $1
  call $~lib/util/number/itoa32
  call $~lib/string/String.__concat
  i32.const 3856
  call $~lib/string/String.__concat
  local.get $0
  call $~lib/string/String.__concat
  drop
  i64.const 0
 )
 (func $assembly/index/getCollateralType (param $0 i32) (result i32)
  i32.const 1392
  local.get $0
  call $~lib/string/String.__concat
  drop
  i32.const 0
 )
 (func $assembly/index/setCollateralRatio (param $0 i64) (result i32)
  local.get $0
  i64.const 100
  i64.lt_u
  if
   i32.const 4160
   i32.const 1824
   i32.const 123
   i32.const 9
   call $~lib/builtins/abort
   unreachable
  end
  i32.const 1
 )
 (func $assembly/index/getCollateralRatio (result i64)
  i64.const 150
 )
 (func $assembly/index/setBinomTokenAddress (param $0 i32) (result i32)
  i32.const 1
 )
 (func $assembly/index/pause (result i32)
  i32.const 1
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
 (func $~lib/rt/stub/__pin (param $0 i32) (result i32)
  local.get $0
 )
 (func $~lib/rt/stub/__unpin (param $0 i32)
 )
 (func $~lib/rt/stub/__collect
 )
 (func $~start
  (local $0 i32)
  i32.const 4268
  global.set $~lib/rt/stub/offset
  i32.const 0
  i32.const 4
  call $~lib/rt/stub/__new
  local.set $0
  i32.const 1052
  i32.load
  drop
  local.get $0
  global.set $assembly/index/contractInstance
 )
)
