package main

func luaF_newLclosure(L *lua_State, n int) *LClosure {
	var o GCObject = luaC_newobj(L, LUA_TLCL)
	var c *LClosure = gco2lcl(o)
	c.p = nil
	c.nupvalues = lu_byte(n)
	for n > 0 {
		n--
		c.upvals[n] = nil
	}
	return c
}

func luaF_initupvals(L *lua_State, cl *LClosure) {
	var i lu_byte
	for i = 0; i < cl.nupvalues; i++ {
		var uv *UpVal = new(UpVal)
		uv.refcount = 1
		uv.v = &uv.u.value
		setnilvalue(uv.v)
		cl.upvals[i] = uv
	}
}

func luaF_close(L *lua_State, level StkId) {

}

func luaF_newproto(L *lua_State) *Proto {
	var o GCObject = luaC_newobj(L, LUA_TPROTO)
	var f *Proto = gco2p(o)
	f.k = nil
	return f
}
