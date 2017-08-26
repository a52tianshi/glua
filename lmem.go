package main

import (
	"fmt"
)

//内存管理全部重写

func luaM_newobject(L *lua_State, tag int) GCObject {
	switch tag {
	case LUA_TTABLE:
		return new(Table)
	case LUA_TSTRING:
		return new(TString)
	case LUA_TFUNCTION:
		return new(LClosure)
	case LUA_TPROTO:
		return new(Proto)
	default:
		fmt.Println(tag)
		fmt.Println("cqtest fail")
		return nil
	}
}
