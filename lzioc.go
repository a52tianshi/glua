package main

func luaZ_fill(z *ZIO) int {
	var size size_t
	var L *lua_State = z.L
	var buff string
	lua_unlock(L)
	buff = z.reader(L, z.data, &size)
	lua_lock(L)
	if buff == "" || size == 0 {
		return EOZ
	}
	z.n = size - 1
	z.p = 0
	z.p++
	return int(buff[0])
}
func luaZ_init(L *lua_State, z *ZIO, reader lua_Reader, data interface{}) {
	z.L = L
	z.reader = reader
	z.data = data
	z.n = 0
	z.p = 0
}
