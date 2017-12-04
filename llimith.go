package main

import (
	"math"
	"unsafe"
)

type size_t uint64
type lu_mem uint64 //size_t
type l_mem int64   //ptrdiff_t
type ptrdiff_t int64

type lu_byte byte

const LUAI_MAXCCALLS = 200

type Instruction uint64

const (
	LUAI_MAXSHORTLEN = 40
	MAX_SIZET        = math.MaxUint64
	MAX_SIZE         = math.MaxUint64
)

func point2uint(p interface{}) uint {
	return uint(uintptr(unsafe.Pointer(&p))) // & unsafe.Pointer(math.MaxUint32)
}

func luai_apicheck(l *lua_State, e interface{}) {
	assert(e)
}

func api_check(l *lua_State, e bool, msg string) {
	luai_apicheck(l, (e) && msg != "")
}

const (
	STRCACHE_N = 53
	STRCACHE_M = 2

	LUA_MINBUFFER = 32
)

/*
** macros that are executed whenever program enters the Lua core
** ('lua_lock') and leaves the core ('lua_unlock')
 */
func lua_lock(L *lua_State) {
}
func lua_unlock(L *lua_State) {
}
func luai_numeq(a, b lua_Number) bool {
	return a == b
}
func luai_numisnan(a lua_Number) bool {
	return !(luai_numeq(a, a))
}
