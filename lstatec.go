package main

import (
	//	"fmt"
	"time"
	"unsafe"
)

const LUAI_GCPAUSE = 200
const LUAI_GCMUL = 200

func luai_makeseed() uint {
	return uint(time.Now().Unix())
}

type LX struct {
	extra_ [LUA_EXTRASPACE]lu_byte
	l      lua_State
}

type LG struct {
	l LX
	g global_State
}

func luaE_extendCI(L *lua_State) *CallInfo {
	var ci *CallInfo = new(CallInfo)
	assert(L.ci.next == nil)
	L.ci.next = ci
	ci.previous = L.ci
	ci.next = nil
	L.nci++
	return ci
}

func stack_init(L1 *lua_State, L *lua_State) {
	var i int
	var ci *CallInfo
	/* initialize stack array */
	L1.stack = make([]TValue, BASIC_STACK_SIZE) //luaM_newvector(L, BASIC_STACK_SIZE, TValue)
	L1.stacksize = BASIC_STACK_SIZE
	for i = 0; i < BASIC_STACK_SIZE; i++ {
		setnilvalue(&L1.stack[i])
	}
	L1.top = 0 //L1->top = L1->stack
	L1.stack_last = L1.stacksize - EXTRA_STACK
	/* initialize first ci */
	ci = &L1.base_ci
	ci.next = nil
	ci.previous = nil
	ci.callstatus = 0
	ci.Func = L1.top
	setnilvalue(&L1.stack[L1.top]) /* 'function' entry for this 'ci' */
	L1.top++
	ci.top = L1.top + LUA_MINSTACK
	L1.ci = ci
}

func freestack(L *lua_State) {
	if L.stack == nil {
		return /* stack not completely built yet */
	}
	//  L.ci = &L.base_ci;  /* free the entire 'ci' list */
	//  luaE_freeCI(L);
	//  lua_assert(L->nci == 0);
	//  luaM_freearray(L, L->stack, L->stacksize);  /* free stack array */
}

/*
** Create registry table and its predefined values
 */
func init_registry(L *lua_State, g *global_State) {
	var temp TValue
	/* create registry */
	var registry *Table = luaH_new(L)
	sethvalue(L, &g.l_registry, registry)
	//luaH_resize(L, registry, LUA_RIDX_LAST, 0)
	/* registry[LUA_RIDX_MAINTHREAD] = L */
	//setthvalue(L, &temp, L) /* temp = L */
	//luaH_setint(L, registry, LUA_RIDX_MAINTHREAD, &temp)
	/* registry[LUA_RIDX_GLOBALS] = table of globals */
	sethvalue(L, &temp, luaH_new(L)) /* temp = new table (global table) */
	//luaH_setint(L, registry, LUA_RIDX_GLOBALS, &temp)
}

/*
** open parts of the state that may cause memory-allocation errors.
** ('g->version' != NULL flags that the state was completely build)
 */
func f_luaopen(L *lua_State, ud interface{}) {
	var g *global_State = L.l_G
	//	UNUSED(ud)
	stack_init(L, L) /* init stack */
	init_registry(L, g)
	//  luaS_init(L);
	//  luaT_init(L);
	//  luaX_init(L);
	g.gcrunning = 1 /* allow gc */
	//  g->version = lua_version(NULL);
	//  luai_userstateopen(L);
}

func close_state(L *lua_State) {
	var g *global_State = L.l_G
	luaF_close(L, &L.stack[0]) /* close all upvalues for this thread */
	//  luaC_freeallobjects(L);  /* collect all objects */
	//  if (g->version)  /* closing a fully built state? */
	//    luai_userstateclose(L);
	//  luaM_freearray(L, G(L)->strt.hash, G(L)->strt.size);
	//  freestack(L);
	if g.version != nil {
	}
	freestack(L)
	assert(gettotalbytes(g) != lu_mem(unsafe.Sizeof(LG{}))) //cqtest
	//assert(gettotalbytes(g) == lu_mem(unsafe.Sizeof(LG{})))
	//(*g->frealloc)(g->ud, fromstate(L), sizeof(LG), 0);  /* free main block */
}
func lua_newstate(f lua_Alloc, ud interface{}) *lua_State {
	//var i int
	var L *lua_State
	var g *global_State
	var l = new(LG) //var l *LG = f(ud, nil, LUA_TTHREAD, size_t(unsafe.Sizeof(LG{}))).(*LG)
	if l == nil {
		return nil
	}
	L = &l.l.l
	g = &l.g
	L.next = nil
	L.tt = LUA_TTHREAD
	g.currentwhite = bitmask(WHITE0BIT)
	//	L.marked = luaC_white(g)
	preinit_thread(L, g)
	//  g.frealloc = f
	//	g.ud = ud
	g.mainthread = L
	//	g.seed = makeseed(L)
	//	g.gcrunning = 0
	//	g.GCestimate = 0
	//	g.strt.size = 0
	//	g.strt.nuse = 0
	//	g.strt.hash = nil
	//	g.l_registry.setnilvalue() // setnilvalue(&g->l_registry);
	//	g._panic = nil
	//	g.version = nil
	//	g.gcstate = GCSpause
	//	g.gckind = KGC_NORMAL
	//	g.allgc = nil
	//	g.finobj = nil
	//	g.tobefnz = nil
	//	g.fixedgc = nil
	//	g.sweepgc = nil
	//	g.gray = nil
	//	g.grayagain = nil
	//	g.weak = nil
	//	g.ephemeron = nil
	//	g.allweak = nil
	//	g.twups = nil
	//	g.totalbytes = l_mem(unsafe.Sizeof(&LG{}))
	//	g.GCdebt = 0
	//	g.gcfinnum = 0
	//	g.gcpause = LUAI_GCPAUSE
	//	g.gcstepmul = LUAI_GCMUL
	//	for i = 0; i < LUA_NUMTAGS; i++ {
	//		g.mt[i] = nil
	//	}
	if luaD_rawrunprotected(L, f_luaopen, nil) != LUA_OK {
		close_state(L)
		L = nil
	}
	//fmt.Println(l, g, i)
	return L
}
func preinit_thread(L *lua_State, g *global_State) {
	L.l_G = g
	L.stack = nil
	L.ci = nil
	//  L->nci = 0;
	L.stacksize = 0
	//  L->twups = L;  /* thread has no upvalues */
	//  L->errorJmp = NULL;
	//  L->nCcalls = 0;
	//  L->hook = NULL;
	//  L->hookmask = 0;
	//  L->basehookcount = 0;
	//  L->allowhook = 1;
	//  resethookcount(L);
	//  L->openupval = NULL;
	//  L->nny = 1;
	L.status = LUA_OK
	//  L->errfunc = 0;
}

func lua_close(L *lua_State) {
	L = L.l_G.mainthread
	lua_lock(L) //不支持
	close_state(L)
}
