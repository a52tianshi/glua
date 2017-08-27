package main

func luaS_newliteral(L *lua_State, s string) *TString {
	return luaS_newlstr(L, []byte(s), size_t(len(s)))
}

/*
** test whether a string is a reserved word
 */
//测试这个字符串是不是关键字
func isreserved(s *TString) bool {
	return (s.tt == LUA_TSHRSTR && s.extra > 0)
}
