package main

func api_incr_top(L *lua_State) {
	L.top++
	api_check(L, L.top <= L.ci.top, "stack overflow")
}

func adjustresults(L *lua_State, nres int) {
	if (nres) == LUA_MULTRET && L.ci.top < L.top {
		L.ci.top = L.top
	}
}

func api_checknelems(L *lua_State, n int) {
	api_check(L, (n) < (L.top-L.ci.Func), "not enough elements in the stack")
}
