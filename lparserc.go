package main

import (
	"github.com/golang/glog"
)

/*
** nodes for block list (list of active blocks)
 */
type BlockCnt struct {
	previous   *BlockCnt /* chain */
	firstlabel int       /* index of first label in this block */
	firstgoto  int       /* index of first pending goto in this block */
	nactvar    byte      /* # active locals outside the block */
	upval      byte      /* true if some variable in the block is an upvalue */
	isloop     byte      /* true if 'block' is a loop */
}

func error_expected(ls *LexState, token int) {
	luaX_syntaxerror(ls, luaO_pushfstring(ls.L, "%s expected", luaX_token2str(ls, token)))
}
func errorlimit(fs *FuncState, limit int, what string) {
	var L *lua_State = fs.ls.L
	var msg string
	var line = fs.f.linedefined
	var where string
	if line == 0 {
		where = "main function"
	} else {
		where = luaO_pushfstring(L, "function at line %d", line)
	}
	msg = luaO_pushfstring(L, "too many %s (limit is %d) in %s", what, limit, where)
	luaX_syntaxerror(fs.ls, msg)
}
func checklimit(fs *FuncState, v int, l int, what string) {
	if v > l {
		errorlimit(fs, l, what)
	}
}
func check(ls *LexState, c int) {
	if ls.t.token != c {
		error_expected(ls, c)
	}
}
func init_exp(e *expdesc, k expkind, i int) {
	e.f = NO_JUMP
	e.t = NO_JUMP
	e.k = k
	e.u.ival = lua_Integer(i)
}

func newupvalue(fs *FuncState, name *TString, v *expdesc) int {
	var f *Proto = fs.f
	var oldsize int = f.sizeupvalues
	checklimit(fs, int(fs.nups+1), MAXUPVAL, "upvalues")
	//	  luaM_growvector(fs->ls->L, f->upvalues, fs->nups, f->sizeupvalues,
	//                  Upvaldesc, MAXUPVAL, "upvalues");
	for ; oldsize < f.sizeupvalues; oldsize++ {
		f.upvalues[oldsize].name = nil
	}
	//  f->upvalues[fs->nups].instack = (v->k == VLOCAL);
	//  f->upvalues[fs->nups].idx = cast_byte(v->u.info);
	//  f->upvalues[fs->nups].name = name;
	//  luaC_objbarrier(fs->ls->L, f, name);
	fs.nups++
	return int(fs.nups - 1)
}

func leavelevel(ls *LexState) { ls.L.nCcalls-- }

func open_func(ls *LexState, fs *FuncState, bl *BlockCnt) {
	var f *Proto
	fs.prev = ls.fs /* linked list of funcstates */
	fs.ls = ls
	ls.fs = fs
	fs.pc = 0
	f = fs.f
	f.source = ls.source
	f.maxstacksize = 2 /* registers 0/1 are always valid */
	//enterblock(fs,bl,0)
}
func close_func(ls *LexState) {
	var L *lua_State = ls.L
	var fs *FuncState = ls.fs
	var f *Proto = fs.f
	luaK_ret(fs, 0, 0) /* final return */
	//	  leaveblock(fs);
	//  luaM_reallocvector(L, f->code, f->sizecode, fs->pc, Instruction);
	f.sizecode = fs.pc
	//  luaM_reallocvector(L, f->lineinfo, f->sizelineinfo, fs->pc, int);
	//  f->sizelineinfo = fs->pc;
	//  luaM_reallocvector(L, f->k, f->sizek, fs->nk, TValue);
	//  f->sizek = fs->nk;
	//  luaM_reallocvector(L, f->p, f->sizep, fs->np, Proto *);
	//  f->sizep = fs->np;
	//  luaM_reallocvector(L, f->locvars, f->sizelocvars, fs->nlocvars, LocVar);
	//  f->sizelocvars = fs->nlocvars;
	//  luaM_reallocvector(L, f->upvalues, f->sizeupvalues, fs->nups, Upvaldesc);
	//  f->sizeupvalues = fs->nups;
	//  lua_assert(fs->bl == NULL);
	//  ls->fs = fs->prev;
	luaC_checkGC(L)
}

/*============================================================*/
/* GRAMMAR RULES */
/*============================================================*/

/*
** check whether current token is in the follow set of a block.
** 'until' closes syntactical blocks, but do not close scope,
** so it is handled in separate.
 */
func block_follow(ls *LexState, withuntil int) int {
	switch ls.t.token {
	case TK_ELSE, TK_ELSEIF, TK_END, TK_EOS:
		return 1
	case TK_UNTIL:
		return withuntil
	default:
		return 0
	}
}
func statlist(ls *LexState) {
	/* statlist -> { stat [';'] } */
	for block_follow(ls, 1) == 0 {
		if ls.t.token == TK_RETURN {
			statement(ls)
			return /* 'return' must be last statement */
		}
		statement(ls)
	}
}
func statement(ls *LexState) {
	var line int = ls.linenumber /* may be needed for error messages */
	switch ls.t.token {
	case ';':
		luaX_next(ls)
	default:
		print(line)
	}
	assert(ls.fs.f.maxstacksize >= ls.fs.freereg && ls.fs.freereg >= ls.fs.nactvar)
	ls.fs.freereg = ls.fs.nactvar
	leavelevel(ls)
}

/*
** compiles the main function, which is a regular vararg function with an
** upvalue named LUA_ENV
 */
func mainfunc(ls *LexState, fs *FuncState) {
	var bl BlockCnt
	var v expdesc

	open_func(ls, fs, &bl)
	fs.f.is_vararg = 1          /* main function is always declared vararg */
	init_exp(&v, VLOCAL, 0)     /* create and... */
	newupvalue(fs, ls.envn, &v) /* ...set environment upvalue */
	//assert(false)
	luaX_next(ls) /* read first token */
	glog.Infoln(ls)
	assert(false)
	statlist(ls) /* parse main body */
	assert(false)
	check(ls, TK_EOS)
	close_func(ls)
}

func luaY_parser(L *lua_State, z *ZIO, buff *Mbuffer, dyd *Dyndata, name string, firstchar int) *LClosure {
	var lexstate LexState
	var funcstate FuncState
	var cl *LClosure = luaF_newLclosure(L, 1) /* create main closure */
	setclLvalue(L, &L.stack[L.top], cl)       ///* anchor it (to avoid being collected) */
	luaD_inctop(L)
	lexstate.h = luaH_new(L)                  /* create table for scanner */
	sethvalue(L, &L.stack[L.top], lexstate.h) /* anchor it */
	luaD_inctop(L)
	cl.p = luaF_newproto(L)
	funcstate.f = cl.p
	funcstate.f.source = luaS_new(L, name) /* create and anchor TString */
	assert(iswhite(funcstate.f))           /* do not need barrier here */
	lexstate.buff = buff
	lexstate.dyd = dyd
	dyd.actvar.n = 0
	dyd.gt.n = 0
	dyd.label.n = 0
	luaX_setinput(L, &lexstate, z, funcstate.f.source, firstchar)
	mainfunc(&lexstate, &funcstate)
	assert(funcstate.prev == nil && funcstate.nups == 1 && lexstate.fs == nil)
	/* all scopes should be correctly finished */
	assert(dyd.actvar.n == 0 && dyd.gt.n == 0 && dyd.label.n == 0)
	L.top--
	return cl
}
