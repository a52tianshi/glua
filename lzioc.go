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

func luaZ_read(z *ZIO, b []byte, n size_t) size_t {
	for n != 0 {
		var m size_t
		if z.n == 0 {
			if luaZ_fill(z) == EOZ { /* try to read more */
				return n /* no more input; return number of missing bytes */
			} else {
				z.n++
				z.p--
			}
		}

		if n <= z.n { /* min. between n and z->n */
			m = n
		} else {
			m = z.n
		}

		copy(b, z.data.(*LoadS).s[z.p:z.p+int(m)])
		z.n -= m
		z.p += int(m)
		b = b[m:]
		n -= m
	}
	return 0
}
