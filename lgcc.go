package main

func markobject(g *global_State, t GCObject) {

}

/*
** barrier for assignments to closed upvalues. Because upvalues are
** shared among closures, it is impossible to know the color of all
** closures pointing to it. So, we assume that the object being assigned
** must be marked.
 */
func luaC_upvalbarrier_(L *lua_State, uv *UpVal) {
	var g *global_State = L.l_G
	var o GCObject = gcvalue(uv.v)
	assert(!upisopen(uv))
	if keepinvariant(g) {
		markobject(g, o)
	}
}

func luaC_newobj(L *lua_State, tt int) GCObject {
	var g *global_State = L.l_G
	var o GCObject = luaM_newobject(L, novariant(tt))
	o.SetMarked(luaC_white(g))
	o.SetTt(lu_byte(tt))
	o.SetNext(g.allgc)
	g.allgc = o
	return o
}
