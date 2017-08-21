package main

const (
	LUA_FLOORN2I = 0
)

func tointeger(o *TValue, i *lua_Integer) int {
	if ttisinteger(o) {
		*i = ivalue(o)
		return 1
	} else {
		return luaV_tointeger(o, i, LUA_FLOORN2I)
	}
}
