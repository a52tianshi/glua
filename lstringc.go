package main

import (
	"math"
)

const MEMERRMSG = "not enough memory"
const LUAI_HASHLIMIT = 5

func luaS_hash(str []byte, l size_t, seed uint) uint {
	var h uint = seed ^ uint(l)
	var step size_t = l>>LUAI_HASHLIMIT + 1
	for ; l >= step; l -= step {
		h ^= ((h << 5) + (h >> 2) + uint(str[l-1]))
	}
	return h
}

/*
** resizes the string table
 */
func luaS_resize(L *lua_State, newsize int) {
	var tb *stringtable = &L.l_G.strt
	if newsize > tb.size {
		tb.hash = append(tb.hash, make([]*TString, newsize-tb.size)...)
	}
	for i := 0; i < tb.size; i++ { /* rehash */
		var p *TString = tb.hash[i]
		tb.hash[i] = nil
		for p != nil { /* for each node in the list */
			var hnext *TString = p.u.hnext           /* save next */
			var h uint = uint(lmod(p.hash, newsize)) /* new position */
			p.u.hnext = tb.hash[h]                   /* chain it */
			tb.hash[h] = p
			p = hnext
		}
	}
	if newsize < tb.size {
		assert(tb.hash[newsize] == nil && tb.hash[tb.size-1] == nil)
		var temp []*TString = make([]*TString, newsize)
		copy(temp, tb.hash[:newsize])
		tb.hash = temp
	}
	tb.size = newsize
}
func luaS_init(L *lua_State) {
	var g *global_State = L.l_G
	luaS_resize(L, MINSTRTABSIZE) /* initial size of string table */
	/* pre-create memory-error message */
	g.memerrmsg = luaS_newliteral(L, MEMERRMSG)
	luaC_fix(L, obj2gco(g.memerrmsg)) /* it should never be collected */
	for i := 0; i < STRCACHE_N; i++ { /* fill cache with valid strings */
		for j := 0; j < STRCACHE_M; j++ {
			g.strcache[i][j] = g.memerrmsg
		}
	}
}
func luaS_new(L *lua_State, str string) *TString {
	var ret *TString = luaC_newobj(L, LUA_TSTRING).(*TString)
	ret.data = str
	return ret
}

func createstrobj(L *lua_State, l size_t, tag int, h uint) *TString {
	var ts *TString
	var o GCObject
	//var totalsize size_t = sizelstring(l) /* total size of TString object */
	o = luaC_newobj(L, tag /*, totalsize*/)
	ts = gco2ts(o)
	ts.hash = h
	ts.extra = 0
	//getstr(ts)[l] = '\0';  /* ending 0 */
	return ts
}

/*
** checks whether short string exists and reuses it or creates a new one
 */
func internshrstr(L *lua_State, str []byte, l size_t) *TString {
	var ts *TString
	var g *global_State = L.l_G
	var h uint = luaS_hash(str, l, g.seed)
	var list **TString = &g.strt.hash[lmod(h, g.strt.size)]
	assert(str != nil)
	for ts = *list; ts != nil; ts = ts.u.hnext {
		if l == size_t(ts.shrlen) && string(str) == getstr(ts) {
			if isdead(g, ts) {
				changewhite(ts)
			}
			return ts
		}
	}
	if g.strt.nuse >= g.strt.size && g.strt.size <= math.MaxInt32/2 {
		luaS_resize(L, g.strt.size*2)
		list = &g.strt.hash[lmod(h, g.strt.size)] /* recompute with new size */
	}
	ts = createstrobj(L, l, LUA_TSHRSTR, h)
	ts.shrlen = lu_byte(l)
	ts.data = string(str)
	ts.u.hnext = *list
	*list = ts
	g.strt.nuse++
	return ts
}

func luaS_newlstr(L *lua_State, str []byte, l size_t) *TString {
	if l <= LUAI_MAXSHORTLEN {
		return internshrstr(L, str, l)
	} else {
		var ts *TString
		if l >= MAX_SIZE {

		}
		ts = luaC_newobj(L, LUA_TSTRING).(*TString)
		ts.data = string(str[:l])
		return ts
	}
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
