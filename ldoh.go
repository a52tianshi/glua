package main

func savestack(L *lua_State, p int) ptrdiff_t {
	return ptrdiff_t(p) * 8 //sizeof(StkId)
}

type Pfunc func(L *lua_State, ud interface{})
