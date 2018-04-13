package main

import (
	"fmt"
)

const (
	LUAI_BITSINT = 32

	LUA_INT_INT      = 1
	LUA_INT_LONG     = 2
	LUA_INT_LONGLONG = 3

	LUA_FLOAT_FLOAT      = 1
	LUA_FLOAT_DOUBLE     = 2
	LUA_FLOAT_LONGDOUBLE = 3

	//64
	LUA_INT_TYPE   = LUA_INT_LONGLONG
	LUA_FLOAT_TYPE = LUA_FLOAT_DOUBLE

	LUA_PATH_SEP  = ";"
	LUA_PATH_MARK = "?"
	LUA_EXEC_DIR  = "!"

	LUA_IDSIZE = 60
)

func lua_number2str(s []byte, sz uint, n lua_Number) size_t {
	return l_sprintf((s), sz, LUA_INTEGER_FMT, LUAI_UACNUMBER(n))
}

type LUAI_UACNUMBER float64
type lua_Integer int64
type lua_Number float64

const LUA_INTEGER_FMT = "%d" //golang 基本常识

type LUAI_UACINT int64

func lua_integer2str(s []byte, sz uint, n lua_Integer) size_t {
	return l_sprintf((s), sz, LUA_INTEGER_FMT, LUAI_UACINT(n))
}
func l_sprintf(s []byte, sz uint, f string, i interface{}) size_t {
	var tempstr = fmt.Sprintf(f, i)
	copy(s, []byte(tempstr))
	return size_t(len(tempstr))
}

type lua_KContext ptrdiff_t

func lua_getlocaledecpoint() byte {
	return '.'
}

const (
	LUAI_MAXSTACK  = 1000000
	LUA_EXTRASPACE = 8 //(sizeof(void *))
)
