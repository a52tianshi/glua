package main

func luaM_toobig(L *lua_State) {
	luaG_runerror(L, "memory allocation error: block too big")
}
