package main

import (
	"github.com/golang/glog"
)

/* value at a non-valid index */
var NONVALIDVALUE = luaO_nilobject

/* test for pseudo index */
func ispseudo(i int) bool { return i <= LUA_REGISTRYINDEX }
func index2addr(L *lua_State, idx int) *TValue {
	var ci *CallInfo = L.ci
	if idx > 0 {
		var o *TValue = &L.stack[ci.Func+idx]
		api_check(L, idx <= ci.top-(ci.Func+1), "unacceptable index")
		if ci.Func+idx >= L.top { // (o >= L->top)
			return NONVALIDVALUE
		}
		return o

	} else if !ispseudo(idx) { /* negative index */
		api_check(L, idx != 0 && -idx <= L.top-(ci.Func+1), "invalid index")
		return &L.stack[L.top+idx]
	}
	return nil
}

func lua_checkstack(L *lua_State, n int) int {
	var res int
	//  CallInfo *ci = L->ci;
	lua_lock(L)
	//  api_check(L, n >= 0, "negative 'n'");
	//  if (L->stack_last - L->top > n)  /* stack large enough? */
	//    res = 1;  /* yes; check is OK */
	//  else {  /* no; need to grow stack */
	//    int inuse = cast_int(L->top - L->stack) + EXTRA_STACK;
	//    if (inuse > LUAI_MAXSTACK - n)  /* can grow without overflow? */
	//      res = 0;  /* no */
	//    else  /* try to grow stack */
	//      res = (luaD_rawrunprotected(L, &growstack, &n) == LUA_OK);
	//  }
	//  if (res && ci->top < L->top + n)
	//    ci->top = L->top + n;  /* adjust frame top */
	lua_unlock(L)
	return res
}

func lua_atpanic(L *lua_State, panicf lua_CFunction) lua_CFunction {
	var old lua_CFunction
	lua_lock(L)
	old = L.l_G.Panic
	L.l_G.Panic = panicf
	lua_unlock(L)
	return old
}
func lua_version(L *lua_State) *lua_Number {
	var version lua_Number = LUA_VERSION_NUM
	if L == nil {
		return &version
	} else {
		return L.l_G.version
	}
}
func lua_gettop(L *lua_State) int {
	return L.top - (L.ci.Func + 1)
}

func lua_settop(L *lua_State, idx int) {
	//	var func_ StkId = L.ci.func_

	//	var temp int
	//	for i, v := range L.stack {
	//		v = func_
	//		temp = i
	//	}

	lua_lock(L)
	//	if idx >= 0 {
	//		api_check(L, idx <= L.stack_last-(temp+1), "new top too large")
	//		for L.top < temp+1+idx {
	//			(*TValue)(L.stack[L.top]).setnilvalue()
	//			L.top++
	//		}
	//		L.top = temp + 1 + idx
	//	} else {
	//		api_check(L, -(idx+1) <= (L.top-(temp+1)), "invalid new top")
	//		L.top += idx + 1
	//	}
	lua_unlock(L)
}

/*
** Reverse the stack segment from 'from' to 'to'
** (auxiliary to 'lua_rotate')
 */
func reverse(L *lua_State, from, to StkId) {

}

/*
** Let x = AB, where A is a prefix of length 'n'. Then,
** rotate x n == BA. But BA == (A^r . B^r)^r.
 */
func lua_rotate(L *lua_State, idx int, n int) {
	var p, t, m StkId
	lua_lock(L)
	t = &L.stack[L.top-1]  /* end of stack segment being rotated */
	p = index2addr(L, idx) /* start of segment */
	//api_checkstackindex(L, idx, p)
	//api_check(L, ITE_int(n >= 0, n, -n) <= (t-p+1), "invalid 'n'")

	// m = (n >= 0 ? t - n : p - n - 1);  /* end of prefix */
	if n >= 0 {
		m = &L.stack[L.top-1-n]
	} else {
		m = GetStackByOpPtr(L, -n-1, p)
	}
	reverse(L, p, m)                        /* reverse the prefix with length 'n' */
	reverse(L, GetStackByOpPtr(L, 1, m), t) /* reverse the suffix */
	reverse(L, p, t)                        /* reverse the entire segment */
	lua_unlock(L)
}
func lua_tointegerx(L *lua_State, idx int, pisnum *int) lua_Integer {
	var res lua_Integer
	var o *TValue = index2addr(L, idx)
	var isnum int = tointeger(o, &res)
	if isnum == 0 {
		res = 0
	}
	if pisnum != nil {
		*pisnum = isnum
	}
	return res
}
func lua_toboolean(L *lua_State, idx int) int {
	var o *TValue = index2addr(L, idx)
	if !l_isfalse(o) {
		return 1
	}
	return 0
}

func lua_tolstring(L *lua_State, idx int, len_ *size_t) string {
	var o StkId = index2addr(L, idx)
	if !ttisstring(o) {
		if !cvt2str(o) {
			if len_ != nil {
				*len_ = 0
			}
			return ""
		}
		lua_lock(L) /* 'luaO_tostring' may create a new string */
		luaO_tostring(L, o)
		luaC_checkGC(L)
		o = index2addr(L, idx) /* previous call may reallocate the stack */
		lua_unlock(L)
	}
	if len_ != nil {
		*len_ = vslen(o)
	}
	return svalue(o)
}

func lua_touserdata(L *lua_State, idx int) interface{} {
	var o StkId = index2addr(L, idx)
	switch ttnov(o) {
	case LUA_TUSERDATA:
		return nil
		//		return getudatamem(uvalue(o))
	case LUA_TLIGHTUSERDATA:
		return pvalue(o)
	default:
		return nil
	}
}

/*
** push functions (C -> stack)
 */

func lua_pushnil(L *lua_State) {
	lua_lock(L)
	//setnilvalue(L->top)
	//api_incr_top(L)
	lua_unlock(L)
}

func lua_pushnumber(L *lua_State, n lua_Number) {
	lua_lock(L)
	// setfltvalue(L->top, n)
	//api_incr_top(L)
	lua_unlock(L)
}

func lua_pushinteger(L *lua_State, n lua_Integer) {
	lua_lock(L)
	setivalue(&L.stack[L.top], n)
	api_incr_top(L)
	lua_unlock(L)
}

/*
** Pushes on the stack a string with given length. Avoid using 's' when
** 'len' == 0 (as 's' can be NULL in that case), due to later use of
** 'memcmp' and 'memcpy'.
 */
func lua_pushlstring(L *lua_State, s []byte, Len size_t) string {
	var ts *TString
	lua_lock(L)
	//ts = (len == 0) ? luaS_new(L, "") : luaS_newlstr(L, s, len);
	if Len == 0 {
		ts = luaS_new(L, "")
	} else {
		ts = luaS_newlstr(L, s, Len)
	}
	setsvalue2s(L, L.top, ts)
	api_incr_top(L)
	luaC_checkGC(L)
	lua_unlock(L)

	glog.Info(L.top, L.stack, L.stack[L.top-1].TValuefields.value_.gc.(*TString))
	return getstr(ts)
}

//func lua_pushvfstring(L *lua_State, fmt string, argp ...string) (ret string) {

//}
func lua_pushfstring(L *lua_State, fmt string, argp ...interface{}) (ret string) {
	lua_lock(L)
	ret = luaO_pushvfstring(L, fmt, argp)
	luaC_checkGC(L)
	lua_unlock(L)
	return
}

func lua_pushcclosure(L *lua_State, fn lua_CFunction, n int) {
	lua_lock(L)
	if n == 0 {
		setfvalue(&L.stack[L.top], fn)
	} else {
		//    CClosure *cl
		//    api_checknelems(L, n);
		//    api_check(L, n <= MAXUPVAL, "upvalue index too large");
		//    cl = luaF_newCclosure(L, n);
		//    cl->f = fn;
		//    L->top -= n;
		//    while (n--) {
		//      setobj2n(L, &cl->upvalue[n], L->top + n);
		//      /* does not need barrier because closure is white */
		//    }
		//    setclCvalue(L, L->top, cl);
	}
	api_incr_top(L)
	luaC_checkGC(L)
	lua_unlock(L)
}
func lua_pushlightuserdata(L *lua_State, p interface{}) {
	lua_lock(L)
	setpvalue(&L.stack[L.top], p)
	api_incr_top(L)
	lua_unlock(L)
}

/*
** get functions (Lua -> stack)
 */
func auxgetstr(L *lua_State, t *TValue, k string) int {
	//  const TValue *slot;
	//  TString *str = luaS_new(L, k);
	//  if (luaV_fastget(L, t, str, slot, luaH_getstr)) {
	//    setobj2s(L, L->top, slot);
	//    api_incr_top(L);
	//  }
	//  else {
	//    setsvalue2s(L, L->top, str);
	//    api_incr_top(L);
	//    luaV_finishget(L, t, L->top - 1, L->top - 1, slot);
	//  }
	//  lua_unlock(L);
	//  return ttnov(L->top - 1);
	return 1
}

func lua_getglobal(L *lua_State, name string) int {
	var reg *Table = hvalue(&L.l_G.l_registry)
	lua_lock(L)
	return auxgetstr(L, luaH_getint(reg, LUA_RIDX_GLOBALS), name)
}

func lua_createtable(L *lua_State, narray int, nrec int) {

}

/*
** 'load' and 'call' functions (run Lua code)
 */

func checkresults(L *lua_State, na, nr int) {
	api_check(L, (nr) == LUA_MULTRET || (L.ci.top-L.top >= (nr)-(na)), "results from function overflow current stack size")
}

type CallS struct {
	Func     int //StkId //L.stack[Func]
	nresults int
}

func f_call(L *lua_State, ud interface{}) {
	var c *CallS = ud.(*CallS)
	luaD_callnoyield(L, c.Func, c.nresults)
}

func lua_pcallk(L *lua_State, nargs, nresults, errfunc int, ctx ptrdiff_t, k lua_KFunction) int {
	var c CallS
	var status int
	var Func ptrdiff_t
	lua_lock(L)
	api_check(L, k == nil || !isLua(L.ci), "cannot use continuations inside hooks")
	api_checknelems(L, nargs+1)
	api_check(L, L.status == LUA_OK, "cannot do calls on non-normal thread")
	checkresults(L, nargs, nresults)
	if errfunc == 0 {
		Func = 0
	} else {
		//    StkId o = index2addr(L, errfunc);
		//    api_checkstackindex(L, errfunc, o);
		//    func = savestack(L, o);
	}
	c.Func = L.top - (nargs + 1) /* function to be called */
	if k == nil || L.nny > 0 {   /* no continuation or no yieldable? */
		c.nresults = nresults /* do a 'conventional' protected call */
		status = luaD_pcall(L, f_call, &c, savestack(L, c.Func), Func)
	}
	//  else {  /* prepare continuation (call is already protected by 'resume') */
	//    CallInfo *ci = L->ci;
	//    ci->u.c.k = k;  /* save continuation */
	//    ci->u.c.ctx = ctx;  /* save context */
	//    /* save information for error recovery */
	//    ci->extra = savestack(L, c.func);
	//    ci->u.c.old_errfunc = L->errfunc;
	//    L->errfunc = func;
	//    setoah(ci->callstatus, L->allowhook);  /* save value of 'allowhook' */
	//    ci->callstatus |= CIST_YPCALL;  /* function can do error recovery */
	//    luaD_call(L, c.func, nresults);  /* do the call */
	//    ci->callstatus &= ~CIST_YPCALL;
	//    L->errfunc = ci->u.c.old_errfunc;
	//    status = LUA_OK;  /* if it is here, there were no errors */
	//  }
	adjustresults(L, nresults)
	lua_unlock(L)
	return status
}

func lua_load(L *lua_State, reader lua_Reader, data interface{}, chunkname string, mode string) int {
	var z ZIO
	var status int
	lua_lock(L)
	if chunkname == "" {
		chunkname = "?"
	}
	luaZ_init(L, &z, reader, data)
	status = luaD_protectedparser(L, &z, chunkname, mode)
	if status == LUA_OK { /* no errors? */
		var f *LClosure = clLvalue(&L.stack[L.top-1])
		if f.nupvalues >= 1 { /* does it have an upvalue? */
			/* get global table from registry */
			var reg *Table = hvalue(&L.l_G.l_registry)
			var gt *TValue = luaH_getint(reg, LUA_RIDX_GLOBALS)
			/* set global table as 1st upvalue of 'f' (may be LUA_ENV) */
			setobj(L, f.upvals[0].v, gt)
			luaC_upvalbarrier(L, f.upvals[0])
		}
	}
	lua_unlock(L)
	return status
}
