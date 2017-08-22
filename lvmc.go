package main

import (
	"fmt"
)

func luaV_tointeger(obj *TValue, p *lua_Integer, mode int) int {
	//again:
	return 0
}
func tostring(L *lua_State, o *TValue) bool {
	return ttisstring(o)
}

//重写
func copy2buff(L *lua_State, top, n int, buff []byte) {
	var tl size_t = 0
	for ; n > 0; n-- {
		var l size_t = vslen(&L.stack[top-n])
		copy(buff[tl:], []byte(svalue(&L.stack[top-n])))
		tl += l
	}
}
func luaV_concat(L *lua_State, total int) {
	assert(total >= 2)
	fmt.Println("cqq ", L.stack[:L.top])
	for total > 1 {
		//var top StkId = &L.stack[L.top]
		var n int = 2 /* number of elements handled in this pass (at least 2) */
		if false {

		} else {
			/* at least two non-empty string values; get as many as possible */
			var tl size_t = vslen(&L.stack[L.top-1])
			var ts *TString
			/* collect total length and number of strings */
			for n = 1; n < total && tostring(L, &L.stack[L.top-n-1]); n++ {
				var l size_t = vslen(&L.stack[L.top-n-1])
				if l >= MAX_SIZE-tl {

				}
				tl += l
			}
			if tl <= LUAI_MAXSHORTLEN {
				var buff []byte = make([]byte, LUAI_MAXSHORTLEN)
				copy2buff(L, L.top, n, buff)
				ts = luaS_newlstr(L, buff, tl)
			} else {

			}
			setsvalue2s(L, L.top-n, ts)
		}
		total -= n - 1
		L.top -= n - 1
	}
}
func luaV_execute(L *lua_State) {

}
