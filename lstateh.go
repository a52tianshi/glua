package main

type l_signalT int //atomic

const (
	EXTRA_STACK      = 5
	BASIC_STACK_SIZE = 2 * LUA_MINSTACK

	/* kinds of Garbage Collection */
	KGC_NORMAL    = 0
	KGC_EMERGENCY = 1 /* gc was forced by an allocation failure */
)

type stringtable struct {
	hash []*TString
	nuse int
	size int
}

type CallInfo struct {
	Func     int //stack[Func] //StkId
	top      int //stack[top] //StkId
	previous *CallInfo
	next     *CallInfo
	u        struct {
		l struct {
			base    StkId
			savedpc *Instruction
		}
		c struct {
			k           lua_KFunction
			old_errfunc ptrdiff_t
			ctx         lua_KContext
		}
	}
	extra      ptrdiff_t
	nresults   int16
	callstatus uint16
}

/*
** Bits in CallInfo status
 */
const (
	CIST_OAH       = 1 << 0
	CIST_LUA       = 1 << 1
	CIST_HOOKED    = 1 << 2
	CIST_FRESH     = 1 << 3
	CIST_YPCALL    = 1 << 4
	CIST_TAIL      = 1 << 5
	CIST_HOOKYIELD = 1 << 6
	CIST_LEQ       = 1 << 7
	CIST_FIN       = 1 << 8
)

func isLua(ci *CallInfo) bool {
	return (ci.callstatus & CIST_LUA) != 0
}

type global_State struct {
	//	frealloc     lua_Alloc
	//	ud           interface{}
	totalbytes l_mem
	GCdebt     l_mem
	//	GCmemtrav    lu_mem
	//	GCestimate   lu_mem
	strt         stringtable
	l_registry   TValue
	seed         uint
	currentwhite lu_byte
	gcstate      lu_byte
	//	gckind       lu_byte
	gcrunning lu_byte
	allgc     GCObject
	//	sweepgc      **GCObject
	//	finobj       *GCObject
	//	gray         *GCObject
	//	grayagain    *GCObject
	//	weak         *GCObject
	//	ephemeron    *GCObject
	//	allweak      *GCObject
	//	tobefnz      *GCObject
	fixedgc GCObject
	//	twups        *lua_State
	//	gcfinnum     uint
	//	gcpause      int
	//	gcstepmul    int
	Panic      lua_CFunction
	mainthread *lua_State
	version    *lua_Number
	memerrmsg  *TString
	//	tmname       [TM_N]*TString
	//	mt           [LUA_NUMTAGS]*Table
	strcache [STRCACHE_N][STRCACHE_M]*TString
}

type lua_State struct {
	CommonHeader
	nci           uint16
	status        lu_byte
	top           int //StkId  //stack[top]   ///* first free slot in the stack */
	l_G           *global_State
	ci            *CallInfo
	oldpc         *Instruction
	stack_last    int //StkId  //stack[stack_last]
	stack         []TValue
	openupval     *UpVal
	gclist        *GCObject
	twups         *lua_State
	errorJmp      *lua_longjmp
	base_ci       CallInfo
	hook          lua_Hook //volatile
	errfunc       ptrdiff_t
	stacksize     int
	basehookcount int
	hookcount     int
	nny           uint16
	nCcalls       uint16
	hookmask      l_signalT
	allowhoot     lu_byte
}

/* macros to convert a GCObject into a specific value */
func gco2ts(o GCObject) *TString {
	assert(novariant(int(o.Tt())) == LUA_TSTRING)
	return o.(*TString)
}
func gco2lcl(o GCObject) *LClosure {
	assert(o.Tt() == LUA_TLCL)
	return o.(*LClosure)
}
func gco2t(o GCObject) *Table {
	assert(o.Tt() == LUA_TTABLE)
	return o.(*Table)
}
func gco2p(o GCObject) *Proto {
	assert(o.Tt() == LUA_TPROTO)
	return o.(*Proto)
}

/* macro to convert a Lua object into a GCObject */
func obj2gco(v GCObject) GCObject {
	assert(novariant(int(v.Tt())) < LUA_TDEADKEY)
	return v
}

/* actual number of total bytes allocated */
func gettotalbytes(g *global_State) lu_mem {
	return lu_mem(g.totalbytes + g.GCdebt)
}
