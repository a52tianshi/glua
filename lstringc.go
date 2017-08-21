package main

func luaS_new(L *lua_State, str string) *TString {
	var ret *TString = luaC_newobj(L, LUA_TSTRING).(*TString)
	ret.data = str
	return ret
}
func luaS_newlstr(L *lua_State, str []byte, l size_t) *TString {
	var ret *TString = luaC_newobj(L, LUA_TSTRING).(*TString)
	ret.data = string(str[:l])
	return ret
}

//const (
//	LUAI_HASHLIMIT = 5
//)

//func luaS_hash(str []byte, l size_t, seed uint) uint {
//	var h uint = seed ^ uint(l)
//	var step = l>>LUAI_HASHLIMIT + 1
//	for ; l >= step; l -= step {
//		h ^= (h << 5) + (h >> 2) + uint(str[l-1])
//	}
//	return h
//}

///*
//** checks whether short string exists and reuses it or creates a new one
// */
//func internshrstr(L *lua_State, str []byte, l size_t) *TString {
//	var ts *TString
//	var g *global_State = L.l_G
//	var h uint = luaS_hash(str, l, g.seed)
//	var list **TString = &g.strt.hash[lmod(h, g.strt.size)]
//	assert(str != nil)
//	for ts = *list; ts != nil; ts = ts.u.hnext {
//		if l == size_t(ts.shrlen) && (memcmp(str, getstr(ts), l*1)) {
//			/* found! */
//			if isdead(g, ts) {
//				changewhite(ts)
//			}
//			return ts
//		}
//	}
//	assert(false)
//	return ts
//}

///*
//** new string (with explicit length)
// */
//func luaS_newlstr(L *lua_State, str []byte, l size_t) *TString {
//	if l <= LUAI_MAXSHORTLEN {
//		return internshrstr(L, str, l)
//	} else {
//		//var ts *TString

//	}
//	assert(false)
//	return nil
//}

////缓存
///*
//** Create or reuse a zero-terminated string, first checking in the
//** cache (using the string address as a key). The cache can contain
//** only zero-terminated strings, so it is safe to use 'strcmp' to
//** check hits.
// */
//func luaS_new(L *lua_State, str []byte) *TString {
//	var i uint = point2uint(str) % STRCACHE_N
//	var j int
//	var p *([STRCACHE_M]*TString) = &L.l_G.strcache[i]
//	for j = 0; j < STRCACHE_M; j++ {
//		if str == getstr(p[j]) {
//			return p[j]
//		}
//	}
//	//不命中
//	for j = STRCACHE_M - 1; j > 0; j-- {
//		p[j] = p[j-1]
//	}
//	p[0] = luaS_newlstr(L, str, strlen(str))
//	return p[0]
//}
