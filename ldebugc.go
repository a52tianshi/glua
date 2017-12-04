package main

func ci_func(ci *CallInfo, L *lua_State) *LClosure {
	return clLvalue(&L.stack[ci.Func])
}
func currentpc(ci *CallInfo, L *lua_State) int {
	assert(isLua(ci))
	return pcRel(ci.u.l.savedpc, ci_func(ci, L).p)
}
func currentline(ci *CallInfo, L *lua_State) int {
	return getfuncline(ci_func(ci, L).p, currentpc(ci, L))
}

func luaG_addinfo(L *lua_State, msg string, src *TString, line int) string {
	var buff [LUA_IDSIZE]byte
	if src != nil {
		assert(false)
		//luaO_chunkid(buff, getstr(src), LUA_IDSIZE);
	} else { /* no source available; use "?" instead */
		buff[0] = '?'
		buff[1] = 0
	}
	return luaO_pushfstring(L, "%s:%d: %s", buff, line, msg)
}

func luaG_errormsg(L *lua_State) {
	if L.errfunc != 0 {
		var errfunc StkId = restorestack(L, L.errfunc)
		setobjs2s(L, &L.stack[L.top], &L.stack[L.top-1]) /* move argument */
		setobjs2s(L, &L.stack[L.top-1], errfunc)         /* push function */
		L.top++
		luaD_callnoyield(L, L.top-2, 1)
	}
	luaD_throw(L, LUA_ERRRUN)
}
func luaG_runerror(L *lua_State, fmt ...interface{}) {
	var ci *CallInfo = L.ci
	var msg string
	msg = luaO_pushvfstring(L, fmt[0].(string), fmt[1:])
	if isLua(ci) { /* if Lua function, add source:line information */
		luaG_addinfo(L, msg, ci_func(ci, L).p.source, currentline(ci, L))
	}
	luaG_errormsg(L)
}
