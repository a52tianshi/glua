package main

const (
	LUA_FLOORN2I = 0
)

//#if !defined(LUA_NOCVTN2S)
//#define cvt2str(o)	ttisnumber(o)
//#else
//#define cvt2str(o)	0	/* no conversion from numbers to strings */
//#endif
//LUA_NOCVTN2S
//这里是定义宏 表示 数字是否禁止默认转成字符串
//本代码不禁止

func cvt2str(o *TValue) bool {
	return ttisnumber(o)
}

//同上
func cvt2num(o *TValue) bool {
	return ttisstring(o)
}

func tointeger(o *TValue, i *lua_Integer) int {
	if ttisinteger(o) {
		*i = ivalue(o)
		return 1
	} else {
		return luaV_tointeger(o, i, LUA_FLOORN2I)
	}
}

func luaV_rawequalobj(t1 *TValue, t2 *TValue) bool {
	//test
	return false
}
