package main

func bitmask(b lu_byte) lu_byte       { return 1 << b }
func bit2mask(b1, b2 lu_byte) lu_byte { return bitmask(b1) | bitmask(b2) }

const (
	WHITE0BIT     = 0
	WHITE1BIT     = 1
	BLACKBIT      = 2
	FINALIZEDEBIT = 3
)

var WHITEBITS = bit2mask(WHITE0BIT, WHITE1BIT)

func otherwhite(g *global_State) lu_byte {
	return g.currentwhite ^ WHITEBITS
}
func isdeadm(ow lu_byte, m lu_byte) bool {
	return (m^WHITEBITS)&ow == 0
}
func isdead(g *global_State, v GCObject) bool {
	return isdeadm(otherwhite(g), v.Marked())
}
func changewhite(x GCObject) {
	x.SetMarked(x.Marked() ^ WHITEBITS)
}

func luaC_white(g *global_State) lu_byte {
	return g.currentwhite & WHITEBITS
}

/*
** Does one step of collection when debt becomes positive. 'pre'/'pos'
** allows some adjustments to be done only when needed. macro
** 'condchangemem' is used only for heavy tests (forcing a full
** GC cycle on every opportunity)
 */
func luaC_condGC(L *lua_State, pre, pos int) {
	if L.l_G.GCdebt > 0 {
		//
	}
	//condchangemem
}

/* more often than not, 'pre'/'pos' are empty */
func luaC_checkGC(L *lua_State) {
	luaC_condGC(L, 0, 0)
}
