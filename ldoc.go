package main

import (
	"github.com/golang/glog"
)

func LUAI_THROW(L *lua_State, c *lua_longjmp) {
	panic(c)
}

type luai_jmpbuf int

type lua_longjmp struct {
	previous *lua_longjmp
	b        luai_jmpbuf
	status   int /* error code */
}

func luaD_throw(L *lua_State, errcode int) {
	if L.errorJmp != nil {
		L.errorJmp.status = errcode
		LUAI_THROW(L, L.errorJmp)
	} else {
		var g *global_State = L.l_G
		L.status = byte(errcode)
		if g.mainthread.errorJmp != nil {
			setobjs2s(L, &g.mainthread.stack[g.mainthread.top], &L.stack[L.top-1])
			g.mainthread.top++                /* copy error obj. */
			luaD_throw(g.mainthread, errcode) /* re-throw in main thread */
		} else {
			if g.Panic != nil {

			}
			abort()
		}
	}
}

//妖风写法
func luaD_rawrunprotected(L *lua_State, f Pfunc, ud interface{}) (ret int) {
	var oldnCcalls uint16 = L.nCcalls
	var lj lua_longjmp
	lj.status = LUA_OK
	lj.previous = L.errorJmp /* chain new error handler */
	L.errorJmp = &lj

	defer func() {
		//err := recover()
		//先让它崩溃
		var err error
		if err != nil {
			glog.Infoln("err", err)
			if lj.status == 0 {
				lj.status = -1
			}
		}
		L.errorJmp = lj.previous
		L.nCcalls = oldnCcalls
		ret = lj.status
	}()
	f(L, ud)
	return ret
}

func luaD_inctop(L *lua_State) {
	luaD_checkstack(L, 1)
	L.top++
}

func next_ci(L *lua_State) *CallInfo {
	if L.ci.next != nil {
		L.ci = L.ci.next
	} else {
		L.ci = luaE_extendCI(L)
	}
	return L.ci
}

/*
** Prepares a function call: checks the stack, creates a new CallInfo
** entry, fills in the relevant information, calls hook if needed.
** If function is a C function, does the call, too. (Otherwise, leave
** the execution ('luaV_execute') to the caller, to allow stackless
** calls.) Returns true iff function has been executed (C function).
 */
func luaD_precall(L *lua_State, Func int, nResults int) bool {
	var f lua_CFunction
	var ci *CallInfo
	switch ttype(&L.stack[Func]) {
	case LUA_TCCL:
	case LUA_TLCF: /* light C function */
		f = fvalue(&L.stack[Func])
		var n int /* number of returns */

		ci = next_ci(L)
		ci.nresults = int16(nResults)
		ci.Func = Func
		ci.top = L.top + LUA_MINSTACK
		assert(ci.top <= L.stack_last)
		ci.callstatus = 0
		lua_unlock(L)
		n = f(L)
		lua_lock(L)
		api_checknelems(L, n)
		//luaD_poscall(L, ci, L->top - n, n)
	}
	return true
}

/*
** Call a function (C or Lua). The function to be called is at *func.
** The arguments are on the stack, right after the function.
** When returns, all the results are on the stack, starting at the original
** function position.
 */
func luaD_call(L *lua_State, Func int, nResults int) {
	L.nCcalls++
	if L.nCcalls >= LUAI_MAXCCALLS {
		//stackerror(L)
	}
	if !luaD_precall(L, Func, nResults) {
		luaV_execute(L)
	}
	L.nCcalls--
}

/*
** Similar to 'luaD_call', but does not allow yields during the call
 */
func luaD_callnoyield(L *lua_State, Func int, nResults int) {
	L.nny++
	luaD_call(L, Func, nResults)
	L.nny--
}

func luaD_pcall(L *lua_State, Func Pfunc, u interface{}, old_top ptrdiff_t, ef ptrdiff_t) int {
	var status int
	//	var old_ci *CallInfo = L.ci
	//	var old_allowhooks byte = L.allowhoot
	//	var old_nny uint16 = L.nny
	var old_errfunc ptrdiff_t = L.errfunc
	L.errfunc = ef
	status = luaD_rawrunprotected(L, Func, u)
	if status != LUA_OK {
		glog.Infoln("cqerr")
	}
	L.errfunc = old_errfunc
	return status
}

/*
** Execute a protected parser.
 */
//词法解析
type SParser struct {
	z    *ZIO
	buff Mbuffer /* dynamic structure used by the scanner */
	dyd  Dyndata /* dynamic structures used by the parser */
	mode string
	name string
}

func checkmode(L *lua_State, mode, x string) {
	if mode != "" && strchr(mode, []byte(x)[0]) == -1 {
		luaO_pushfstring(L, "attempt to load a %s chunk (mode is '%s')", x, mode)
		panic(LUA_ERRSYNTAX)
	}
}
func f_parser(L *lua_State, ud interface{}) {
	var cl *LClosure
	var p *SParser = ud.(*SParser)
	var c int = zgetc(p.z) /* read first character */

	glog.Infoln("cq", c, cl, p.buff)
	if byte(c) == LUA_SIGNATURE[0] {
		assert(false)
		checkmode(L, p.mode, "binary") //暂时不会进入二进制模式(话说我也没见过lua的二进制文本 不常见)
		//cl = luaU_undump(L, p.z, p.name)
	} else {
		checkmode(L, p.mode, "text")
		cl = luaY_parser(L, p.z, &p.buff, &p.dyd, p.name, c)
	}
	glog.Infoln("cq", cl)
	assert(cl.nupvalues == byte(cl.p.sizeupvalues))
	luaF_initupvals(L, cl)
}
func luaD_protectedparser(L *lua_State, z *ZIO, name string, mode string) int {
	var p SParser
	var status int
	L.nny++ /* cannot yield during parsing */ //不能释放协程
	p.z = z
	p.name = name
	p.mode = mode
	p.dyd.actvar.arr = nil
	p.dyd.actvar.size = 0
	p.dyd.gt.arr = nil
	p.dyd.gt.size = 0
	p.dyd.label.arr = nil
	p.dyd.label.size = 0
	luaZ_initbuffer(L, &p.buff)
	status = luaD_pcall(L, f_parser, &p, savestack(L, L.top), L.errfunc)
	luaZ_freebuffer(L, &p.buff)
	//	luaM_freearray(L, p.dyd.actvar.arr, p.dyd.actvar.size)
	p.dyd.actvar.arr = nil
	//	luaM_freearray(L, p.dyd.gt.arr, p.dyd.gt.size)
	p.dyd.gt.arr = nil
	//	luaM_freearray(L, p.dyd.label.arr, p.dyd.label.size)
	p.dyd.label.arr = nil
	L.nny--
	return status
}
