package main

func luaL_error(L *lua_State, fmt ...interface{}) int {
	//  va_list argp;
	//  va_start(argp, fmt);
	//  luaL_where(L, 1);
	//  lua_pushvfstring(L, fmt, argp);
	//  va_end(argp);
	//  lua_concat(L, 2);
	//  return lua_error(L);
	return 1
}

/*
** Ensures the stack has at least 'space' extra slots, raising an error
** if it cannot fulfill the request. (The error handling needs a few
** extra slots to format the error message. In case of an error without
** this extra space, Lua will generate the same 'stack overflow' error,
** but without 'msg'.)
 */
func luaL_checkstack(L *lua_State, space int, msg string) {
	if lua_checkstack(L, space) == 0 {
		if msg != "" {
			luaL_error(L, "stack overflow (%s)", msg)
		} else {
			luaL_error(L, "stack overflow")
		}
	}
}

type LoadS struct {
	s    string
	size size_t
}

func getS(L *lua_State, ud interface{}, size *size_t) string {
	var ls *LoadS = ud.(*LoadS)
	if ls.size == 0 {
		return ""
	}
	*size = ls.size
	ls.size = 0
	return ls.s
}
func luaL_loadbufferx(L *lua_State, buff string, size size_t, name string, mode string) int {
	var ls LoadS
	ls.s = buff
	ls.size = size
	return lua_load(L, getS, &ls, name, mode)
}

var l_alloc = func(ud, ptr interface{}, osize, nsize size_t) interface{} {
	if nsize == 0 {
		//free  ptr
		return nil
	} else {
		//return realloc(ptr, nsize)
		return make([]byte, nsize)
	}
}

func Panic(L *lua_State) int {
	lua_writestringerror("PANIC: unprotected error in call to Lua API (%s)\n", lua_tostring(L, -1))
	return 0
}

func luaL_newstate() *lua_State {
	L := lua_newstate(l_alloc, nil)
	if L != nil {
		lua_atpanic(L, Panic)
	}
	return L
}
func luaL_checkversion_(L *lua_State, ver lua_Number, sz size_t) {
	var v *lua_Number = lua_version(L)
	if sz != LUAL_NUMSIZES { /* check numeric types */
		luaL_error(L, "core and library have incompatible numeric types")
	}
	if v != lua_version(nil) {
		luaL_error(L, "multiple Lua VMs detected")
	} else if (*v) != ver {
		luaL_error(L, "version mismatch: app. needs %f, Lua core provides %f", LUAI_UACNUMBER(ver), LUAI_UACNUMBER(*v))
	}
}
