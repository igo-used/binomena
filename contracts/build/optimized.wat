(module
 (type $0 (func (param i32 i32) (result i32)))
 (type $1 (func (param i32)))
 (type $2 (func (result i32)))
 (global $assembly/index/storedValue (mut i32) (i32.const 0))
 (memory $0 0)
 (export "add" (func $assembly/index/add))
 (export "multiply" (func $assembly/index/multiply))
 (export "store" (func $assembly/index/store))
 (export "retrieve" (func $assembly/index/retrieve))
 (export "memory" (memory $0))
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
)
