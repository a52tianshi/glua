package main

/*
** this function handles only '%d', '%c', '%f', '%p', and '%s'
   conventional formats, plus Lua-specific '%I' and '%U'
*/
func luaO_pushvfstring(L *lua_State, fmt string, va_list []string) string {
	return ""
}
