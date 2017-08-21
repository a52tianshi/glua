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
	default:
		fmt.Println("cqtest fail")
		return nil
	}
}
