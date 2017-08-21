package main

func luaC_newobj(L *lua_State, tt int) GCObject {
	var g *global_State = L.l_G
	var o GCObject = luaM_newobject(L, novariant(tt))
	o.SetMarked(luaC_white(g))
	o.SetTt(lu_byte(tt))
	o.SetNext(g.allgc)
	g.allgc = o
	return o
}
