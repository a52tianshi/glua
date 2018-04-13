package main

import (
	"github.com/golang/glog"
)

func white2gray(x GCObject) {
	x.SetMarked(resetbits(x.Marked(), WHITEBITS))
}
func black2gray(x GCObject) {
	x.SetMarked(resetbits(x.Marked(), BLACKBIT))
}
func markobject(g *global_State, t GCObject) {

}

/*
** barrier that moves collector backward, that is, mark the black object
** pointing to a white object as gray again.
 */
func luaC_barrier_(L *lua_State, o GCObject, v GCObject) {
	var g *global_State = L.l_G
	assert(isblack(o) && iswhite(v) && !isdead(g, v) && !isdead(g, o))
	if keepinvariant(g) {
		//reallymarkobject(g, v)
	} else {
		//assert(issweepphase(g))
		//makewhite(g, o)
	}
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

func luaC_fix(L *lua_State, o GCObject) {
	var g *global_State = L.l_G
	assert(g.allgc == o) /* object must be 1st in 'allgc' list! */
	white2gray(o)        /* they will be gray forever */
	g.allgc = o.Next()   /* remove object from 'allgc' list */
	o.SetNext(g.fixedgc) /* link it to 'fixedgc' list */
	g.fixedgc = o
}

func luaC_newobj(L *lua_State, tt int) GCObject {
	glog.Infoln("token type", tt)
	var g *global_State = L.l_G
	var o GCObject = luaM_newobject(L, novariant(tt))
	o.SetMarked(luaC_white(g))
	o.SetTt(byte(tt))
	o.SetNext(g.allgc)
	g.allgc = o
	return o
}

func checkSizes(L *lua_State, g *global_State) {
	if g.gckind != KGC_EMERGENCY {
		var olddebt l_mem = g.GCdebt
		if g.strt.nuse < g.strt.size/4 { /* string table too big? */
			luaS_resize(L, g.strt.size/2) /* shrink it a little */
		}
		g.GCestimate += lu_mem(g.GCdebt - olddebt) /* update estimate */
	}
}
