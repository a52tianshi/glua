package main

import (
	"github.com/golang/glog"
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
		glog.Infoln(tag)
		glog.Infoln("cqtest fail")
		return nil
	}
}
