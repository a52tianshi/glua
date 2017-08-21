package main

//版本
const (
	LUA_VERSION_MAJOR   = "5"
	LUA_VERSION_MINOR   = "3"
	LUA_VERSION_NUM     = 503
	LUA_VERSION_RELEASE = "4"

	LUA_VERSION   = "Lua " + LUA_VERSION_MAJOR + "." + LUA_VERSION_MINOR
	LUA_RELEASE   = LUA_VERSION + "." + LUA_VERSION_RELEASE
	LUA_COPYRIGHT = LUA_RELEASE + "  Copyright (C) 1994-2017 Lua.org, PUC-Rio"
	LUA_AUTHORS   = "R. Ierusalimschy, L. H. de Figueiredo, W. Celes"
)

/* option for multiple returns in 'lua_pcall' and 'lua_call' */
const LUA_MULTRET = (-1)

/*
** Pseudo-indices
** (-LUAI_MAXSTACK is the minimum valid index; we keep some free empty
** space after that to help overflow detection)
 */
const (
	LUA_REGISTRYINDEX = -LUAI_MAXSTACK - 1000
)

//线程状态
const (
	LUA_OK        = 0
	LUA_YIELD     = 1
	LUA_ERRRUN    = 2
	LUA_ERRSYNTAX = 3
	LUA_ERRMEM    = 4
	LUA_ERRGCMM   = 5
	LUA_ERRERR    = 6
)

//基础类型
const (
	LUA_TNONE          = -1
	LUA_TNIL           = 0
	LUA_TBOOLEAN       = 1
	LUA_TLIGHTUSERDATA = 2
	LUA_TNUMBER        = 3
	LUA_TSTRING        = 4
	LUA_TTABLE         = 5
	LUA_TFUNCTION      = 6
	LUA_TUSERDATA      = 7
	LUA_TTHREAD        = 8
	LUA_NUMTAGS        = 9

	LUA_MINSTACK = 20

	LUA_RIDX_MAINTHREAD = 1
	LUA_RIDX_GLOBALS    = 2
	LUA_RIDX_LAST       = LUA_RIDX_GLOBALS
)

type lua_CFunction func(L *lua_State) int
type lua_KFunction func(L *lua_State, status int, ctx ptrdiff_t) int
type lua_Hook func(L *lua_State, ar *lua_Debug)
type lua_Debug struct {
	event           int
	name            string
	namewhat        string
	what            string
	source          string
	currentline     int
	linedefined     int
	lastlinedefined int
	nups            byte
	nparams         byte
	isvararg        bool
	istaiicall      bool
	short_src       [LUA_IDSIZE]byte

	//private part
	i_ci *CallInfo
}

/*
** 'load' and 'call' functions (load and run Lua code)
 */

func lua_pcall(L *lua_State, n, r, f int) int {
	return lua_pcallk(L, n, r, f, 0, nil)
}

type lua_Alloc func(ud, ptr interface{}, osize, nsize size_t) interface{}

/*
** {==============================================================
** some useful macros
** ===============================================================
 */
func lua_tointeger(L *lua_State, i int) lua_Integer {
	return lua_tointegerx(L, i, nil)
}
func lua_pop(L *lua_State, n int) {
	lua_settop(L, -(n)-1)
}
func lua_pushcfunction(L *lua_State, f lua_CFunction) {
	lua_pushcclosure(L, f, 0)
}
func lua_tostring(L *lua_State, i int) string {
	return lua_tolstring(L, i, nil)
}
func lua_insert(L *lua_State, idx int) {
	lua_rotate(L, idx, 1)
}
func lua_remove(L *lua_State, idx int) {
	lua_rotate(L, idx, -1)
	lua_pop(L, 1)
}
