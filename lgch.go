package main

const (
	//Possible states of the Garbage Collector
	GCSpropagate  = 0
	GCSatomic     = 1
	GCSswpallgc   = 2
	GCSswpfinobj  = 3
	GCSswptobefnz = 4
	GCSswpend     = 5
	GCScallfin    = 6
	GCSpause      = 7
)

func keepinvariant(g *global_State) bool {
	return g.gcstate <= GCSatomic
}

func bitmask(b lu_byte) lu_byte       { return 1 << b }
func bit2mask(b1, b2 lu_byte) lu_byte { return bitmask(b1) | bitmask(b2) }
func testbits(x, m lu_byte) lu_byte   { return x & m }

const (
	WHITE0BIT     = 0
	WHITE1BIT     = 1
	BLACKBIT      = 2
	FINALIZEDEBIT = 3
)

var WHITEBITS = bit2mask(WHITE0BIT, WHITE1BIT)

func iswhite(x GCObject) bool {
	return testbits(x.Marked(), WHITEBITS) != 0
}
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

func luaC_upvalbarrier(L *lua_State, uv *UpVal) {
	if iscollectable(uv.v) && !upisopen(uv) {
		luaC_upvalbarrier_(L, uv)
	}
}
