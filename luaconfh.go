package main

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

type lua_Integer int64
type lua_Number float64

const (
	LUAI_MAXSTACK  = 1000000
	LUA_EXTRASPACE = 8 //(sizeof(void *))
)

type lua_KContext ptrdiff_t
