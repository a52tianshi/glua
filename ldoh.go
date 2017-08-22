package main

func luaD_checkstackaux(L *lua_State, n int, pre, pos interface{}) {
	if L.stack_last-L.top <= n {
		// pre; luaD_growstack(L, n); pos;
	} else {
		//condmovestack(L,pre,pos)
	}
}
func luaD_checkstack(L *lua_State, n int) {
	luaD_checkstackaux(L, n, nil, nil)
}
func savestack(L *lua_State, p int) ptrdiff_t {
	return ptrdiff_t(p) * 8 //sizeof(StkId)
}

type Pfunc func(L *lua_State, ud interface{})
