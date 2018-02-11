package main

import (
	"fmt"
	"os"
)

const LUAL_NUMSIZES = 8*16 + 8 //(sizeof(lua_Integer)*16 + sizeof(lua_Number))
func luaL_checkversion(L *lua_State) {
	luaL_checkversion_(L, LUA_VERSION_NUM, LUAL_NUMSIZES)
}
func luaL_loadbuffer(L *lua_State, s string, sz size_t, n string) int {
	return luaL_loadbufferx(L, s, sz, n, "")
}

func lua_writestring(s string) {
	fmt.Fprint(os.Stdout, s)
}

func lua_writeline() {
	fmt.Fprint(os.Stdout, "\n")
}

func lua_writestringerror(s string, p interface{}) {
	fmt.Fprintf(os.Stderr, s, p)
	os.Stderr.Sync()
}
