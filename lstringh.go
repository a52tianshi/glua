package main

func luaS_newliteral(L *lua_State, s string) *TString {
	return luaS_newlstr(L, []byte(s), size_t(len(s)))
}
