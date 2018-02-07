package main

var log_2 [256]byte = [256]byte{ /* log_2[i] = ceil(log2(i - 1)) */
	0, 1, 2, 2, 3, 3, 3, 3, 4, 4, 4, 4, 4, 4, 4, 4, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5,
	6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6,
	7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
	7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
	8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8,
	8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8,
	8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8,
	8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8}

func luaO_ceillog2(x uint) int {

	var l int
	x--
	for x >= 256 {
		l += 8
		x >>= 8
	}
	return l + int(log_2[x])
}
func pushstr(L *lua_State, str string, l size_t) {
	setsvalue2s(L, L.top, luaS_newlstr(L, []byte(str), l))
	luaD_inctop(L)
}

/*
** this function handles only '%d', '%c', '%f', '%p', and '%s'
   conventional formats, plus Lua-specific '%I' and '%U'
*/
func luaO_pushvfstring(L *lua_State, fmt string, va_list []interface{}) string {
	var n int = 0
	for {
		//改写
		var idx int = strchr(fmt, '%')
		if idx == -1 {
			break
		}
		pushstr(L, fmt, size_t(idx))
		switch []byte(fmt)[idx+1] {
		case 's':
			if len(va_list) >= 1 {
				pushstr(L, va_list[0].(string), size_t(len(va_list[0].(string))))
			} else {
				pushstr(L, "(null)", size_t(len("(null)")))
			}
		}
		n += 2
		fmt = fmt[idx+2:]
	}
	luaD_checkstack(L, 1)
	pushstr(L, fmt, size_t(len(fmt)))
	if n > 0 {
		luaV_concat(L, n+1)
	}
	return svalue(&L.stack[L.top-1])
}
func luaO_pushfstring(L *lua_State, fmt string, argp ...interface{}) string {
	return luaO_pushvfstring(L, fmt, argp)
}
