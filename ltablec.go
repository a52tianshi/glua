package main

import (
	"math"
)

var dummynode *Node = &dummynode_
var dummynode_ Node

//func init() {
//	//dummynode_.i_val.value_ = nil      /* value */
//	dummynode_.i_val.tt_ = LUA_TNIL    /* value */
//	//dummynode_.i_key.nk.value_ = nil   /* key */
//	dummynode_.i_key.nk.tt_ = LUA_TNIL /* key */
//	dummynode_.i_key.nk.next = 0       /* key */
//	dummynode_.i_key.tvk = nil         /* key */
//	dummynode_.i_key.tvk = LUA_TNIL    /* key */
//}

func setnodevector(L *lua_State, t *Table, size uint) {
	if size == 0 {
		t.node = dummynode
		t.lsizenode = 0
		t.lastfree = nil
	}
}

/*
** }=============================================================
 */

func luaH_new(L *lua_State) *Table {
	var o GCObject = luaC_newobj(L, LUA_TTABLE /*, sizeof(Table)*/)
	var t *Table = gco2t(o)
	t.metatable = nil
	t.flags = math.MaxUint8
	t.array = nil
	t.sizearray = 0
	setnodevector(L, t, 0)
	return t
}

/*
** search function for integers
 */
func luaH_getint(t *Table, key lua_Integer) *TValue {
	/* (1 <= key && key <= t->sizearray) */
	if uint(key)-1 < t.sizearray {
		return &t.array[key-1]
	} else {
		return luaO_nilobject
	}
}
