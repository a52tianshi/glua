package main

import (
	"github.com/golang/glog"
)

//luaM_realloc_ 封禁不用
func luaM_reallocv(L *lua_State, b interface{}, on, n int, e interface{}) {

}

const MINSIZEARRAY = 4

//内存管理全部重写
func luaM_reallocvchar(L *lua_State, b *[]byte, on, n size_t) {
	temp := make([]byte, n)
	copy(temp, *b)
	*b = temp
}
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
func luaM_growvector(L *lua_State, v interface{}, nelems byte, size int, type_ interface{}, limit int, e string) {
	if int(nelems)+1 > size {
		var newsize int
		if size >= limit/2 {
			if size >= limit {
				luaG_runerror(L, "too many %s (limit is %d)", e, limit)
			}
			newsize = limit
		} else {
			newsize = size * 2
			if size < MINSIZEARRAY {
				newsize = MINSIZEARRAY
			}
		}
		//上述就是调整size的范围 不过大不过小
		glog.Info(newsize)
		switch v.(type) {
		case *[]Upvaldesc:
			switch type_.(type) {
			case *Upvaldesc:
				temp := make([]Upvaldesc, newsize)
				copy(temp, *v.(*[]Upvaldesc))
				*(v.(*[]Upvaldesc)) = temp
			default:
				assert(false)
			}
		default:
			assert(false)
		}
	}
}
func luaM_reallocvector(L *lua_State, v interface{}, oldn, n uint, type_ interface{}) {
	switch type_.(type) {
	case Instruction:
		if n+1 > (MAX_SIZET / 8) {
			luaM_toobig(L)
		}
		var temp []Instruction = make([]Instruction, n)
		copy(temp, *(v.(*[]Instruction)))
		*(v.(*[]Instruction)) = temp
	case **TString:
		if n+1 > (MAX_SIZET / 8) {
			luaM_toobig(L)
		}
		var temp []*TString = make([]*TString, n)
		copy(temp, *(v.(*[]*TString)))
		*(v.(*[]*TString)) = temp
	case *TValue:
		if n+1 > (MAX_SIZET / 8) {
			luaM_toobig(L)
		}
		var temp []TValue = make([]TValue, n)
		copy(temp, *(v.(*[]TValue)))
		*(v.(*[]TValue)) = temp
		//cqtest
	}
}
