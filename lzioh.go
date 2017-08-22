package main

const EOZ = -1

func zgetc(z *ZIO) int { /* read first character */
	if z.n > 0 {
		z.n--
		z.p++
		var temp []byte = z.data.([]byte)
		return int(temp[z.p-1])
	} else {
		z.n--
		return luaZ_fill(z)
	}
}

type Mbuffer struct {
	buffer   []byte
	n        size_t
	buffsize size_t
}

func luaZ_initbuffer(L *lua_State, buff *Mbuffer) {
	buff.buffer = nil
	buff.buffsize = 0
}

func luaZ_resizebuffer(L *lua_State, buff *Mbuffer, size size_t) {
	buff.buffer = make([]byte, size)
	buff.buffsize = size
}

func luaZ_freebuffer(L *lua_State, buff *Mbuffer) {
	luaZ_resizebuffer(L, buff, 0)
}

type ZIO struct {
	n      size_t /* bytes still unread */
	p      int    /* current position in buffer */
	reader lua_Reader
	data   interface{} /* additional data */
	L      *lua_State  /* Lua state (for reader) */
}
