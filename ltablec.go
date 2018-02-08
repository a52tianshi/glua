package main

import (
	"math"

	"github.com/golang/glog"
)

const (
	MAXABITS = 7
	MAXASIZE = 1 << MAXABITS
	MAXHBITS = MAXABITS - 1
)

var dummynode *Node = &dummynode_
var dummynode_ Node //go自带初始化

//func init() {
//	//dummynode_.i_val.value_ = nil      /* value */
//	dummynode_.i_val.tt_ = LUA_TNIL    /* value */
//	//dummynode_.i_key.nk.value_ = nil   /* key */
//	dummynode_.i_key.nk.tt_ = LUA_TNIL /* key */
//	dummynode_.i_key.nk.next = 0       /* key */
//	dummynode_.i_key.tvk = nil         /* key */
//	dummynode_.i_key.tvk = LUA_TNIL    /* key */
//}
func hashpow2(t *Table, n uint) *Node {
	return gnode(t, lmod(n, int(sizenode(t))))
}
func hashint(t *Table, i lua_Integer) *Node {
	return hashpow2(t, uint(i))
}

/*
** for some types, it is better to avoid modulus by power of 2, as
** they tend to have many 2 factors.
 */
func hashmod(t *Table, n uint) *Node {
	return gnode(t, int(size_t(n)%((sizenode(t)-1)|1)))
}
func hashpointer(t *Table, p GCObject) *Node {
	return hashmod(t, point2uint(p))
}

//返回桶位置
func mainposition(t *Table, key *TValue) *Node {
	glog.Infoln("main ", ttype(key))
	switch ttype(key) {
	case LUA_TNUMINT:
		return hashint(t, ivalue(key))
	//    case LUA_TNUMFLT:
	//      return hashmod(t, l_hashfloat(fltvalue(key)));
	//    case LUA_TSHRSTR:
	//      return hashstr(t, tsvalue(key));
	//    case LUA_TLNGSTR:
	//      return hashpow2(t, luaS_hashlongstr(tsvalue(key)));
	//    case LUA_TBOOLEAN:
	//      return hashboolean(t, bvalue(key));
	//    case LUA_TLIGHTUSERDATA:
	//      return hashpointer(t, pvalue(key));
	//    case LUA_TLCF:
	//      return hashpointer(t, fvalue(key));
	default:
		assert(!ttisdeadkey(key))
		return hashpointer(t, gcvalue(key))
	}
}
func arrayindex(key *TValue) uint {
	if ttisinteger(key) {
		var k lua_Integer = ivalue(key)
		if k > 0 && k <= MAXASIZE {
			return uint(k)
		}
	}
	return 0
}
func computesizes(nums *[MAXABITS + 1]uint, pna *uint) uint {
	var optimal uint
	var twotoi uint = 1 /* 2^i */
	var a uint          /* number of elements smaller than 2^i */
	var na uint         /* number of elements to go to array part */
	for i := 0; *pna > twotoi/2; i++ {
		if nums[i] > 0 {
			a += nums[i]
			if a > twotoi/2 {
				optimal = twotoi
				na = a
			}
		}

		twotoi *= 2
	}
	assert((optimal == 0 || (optimal/2) < na) && na <= optimal)
	*pna = na
	return optimal
}
func countint(key *TValue, nums *[MAXABITS + 1]uint) uint {
	var k uint = arrayindex(key)
	if k != 0 { /* is 'key' an appropriate array index? */
		nums[luaO_ceillog2(k)]++ /* count as such */
		return 1
	}
	return 0
}
func numusearray(t *Table, nums *[MAXABITS + 1]uint) uint {
	var ttlg uint = 1
	var ause uint
	var i uint = 1
	for lg := 0; lg <= MAXABITS; {
		var lc uint
		var lim uint = ttlg
		if lim > t.sizearray {
			lim = t.sizearray
			if i > lim {
				break
			}
		}
		/* count elements in range (2^(lg - 1), 2^lg] */
		for ; i <= lim; i++ {
			if !ttisnil(&t.array[i-1]) {
				lc++
			}
		}
		nums[lg] += lc
		ause += lc

		lg++
		ttlg *= 2
	}
	return ause
}
func numusehash(t *Table, nums *[MAXABITS + 1]uint, pna *uint) int {
	var totaluse int = 0 /* total number of elements */
	var ause uint = 0    /* elements added to 'nums' (can go to array part) */
	var i int = int(sizenode(t)) - 1
	for ; i > 0; i-- {
		var n *Node = &t.node[i]
		if !ttisnil(gval(n)) {
			ause += countint(gkey(n), nums)
			totaluse++
		}
	}
	*pna += ause
	return totaluse
}
func setarrayvector(L *lua_State, t *Table, size uint) {
	var i uint
	t.array = make([]TValue, size)
	//luaM_reallocvector(L, t->array, t->sizearray, size, TValue);
	for i = t.sizearray; i < size; i++ {
		setnilvalue(&t.array[i])
	}
	t.sizearray = size
}

func setnodevector(L *lua_State, t *Table, size uint) {
	if size == 0 {
		t.node = []Node{*dummynode}
		t.lsizenode = 0
		t.lastfree = -1
	} else {
		var lsize int = luaO_ceillog2(size)
		if lsize > MAXHBITS {
			luaG_runerror(L, "table overflow")
		}
		size = uint(twoto(byte(lsize)))
		t.node = make([]Node, size)
		//for i := 0; i < int(size); i++ {//go默认就是0 不用赋值
		//var n *Node = gnode(t, i)
		//n.i_key.nk.next = 0
		//setnilvalue(n.i_key.nk)
		//setnilvalue(&n.i_val)
		//}
		t.lsizenode = byte(lsize)
		t.lastfree = int(size)
	}
}

func luaH_resize(L *lua_State, t *Table, nasize, nhsize uint) {
	var i uint
	var j int
	var oldasize uint
	var oldhsize int = allocsizenode(t)
	var nold []Node = t.node
	if nasize > oldasize {
		setarrayvector(L, t, nasize)
	}
	setnodevector(L, t, nhsize)
	if nasize < oldasize {
		t.sizearray = nasize
		for i = nasize; i < oldasize; i++ {
			if !ttisnil(&t.array[i]) {
				luaH_setint(L, t, lua_Integer(i+1), &t.array[i])
			}
		}
		//luaM_reallocvector(L, t.array, oldasize, nasize, TValue)
		t.array = make([]TValue, nasize)
	}

	for j = oldhsize - 1; j >= 0; j-- {
		var old *Node = &nold[j]
		if !ttisnil(gval(old)) {
			setobjt2t(L, luaH_set(L, t, gkey(old)), gval(old))
		}
	}

	if oldhsize > 0 {
		//luaM_freearray(L, nold, size_t(oldhsize))
	}
}

/*
** nums[i] = number of keys 'k' where 2^(i - 1) < k <= 2^i
 */
func rehash(L *lua_State, t *Table, ek *TValue) {
	var asize uint /* 优化后的size */
	var na uint    /* number of keys in the array part */
	var nums [MAXABITS + 1]uint
	na = numusearray(t, &nums)            /* count keys in array part */
	var totaluse int = int(na)            /* all those keys are integer keys */
	totaluse += numusehash(t, &nums, &na) /* count keys in hash part */
	/* count extra key */
	na += countint(ek, &nums)
	totaluse++
	/* compute new size for array part */
	asize = computesizes(&nums, &na)
	/* resize the table to new computed sizes */
	luaH_resize(L, t, asize, uint(totaluse)-na)
}

/*
** }=============================================================
 */

func luaH_new(L *lua_State) *Table {
	var o GCObject = luaC_newobj(L, LUA_TTABLE /*, sizeof(Table)*/)
	var t *Table = gco2t(o)
	t.metatable = nil
	t.flags = math.MaxUint8
	t.array = nil
	t.sizearray = 0
	setnodevector(L, t, 0)
	return t
}
func getfreepos(t *Table) *Node {
	if !isdummy(t) {
		for t.lastfree > 0 {
			t.lastfree--
			if ttisnil(gkey(&t.node[t.lastfree])) {
				return &t.node[t.lastfree]
			}
		}
	}
	return nil /* 找不到空的桶 */
}

/*
** inserts a new key into a hash table; first, check whether key's main
** position is free. If not, check whether colliding node is in its main
** position or not: if it is not, move colliding node to an empty place and
** put new key in its main position; otherwise (colliding node is in its main
** position), new key goes to an empty position.
 */

//就是插入一个哈希表
func luaH_newkey(L *lua_State, t *Table, key *TValue) *TValue {
	var mp *Node
	var aux TValue
	glog.Infoln(key, key.TValuefields.tt_)
	if ttisnil(key) {
		luaG_runerror(L, "table index is nil")
	} else if ttisfloat(key) {
		var k lua_Integer
		if luaV_tointeger(key, &k, 0) != 0 { /* does index fit in an integer? */
			setivalue(&aux, k)
			key = &aux /* insert it as an integer */
		} else if luai_numisnan(fltvalue(key)) {
			luaG_runerror(L, "table index is NaN")
		}
	}
	mp = mainposition(t, key)
	if !ttisnil(gval(mp)) || isdummy(t) { /* main position is taken? */
		var othern *Node
		var f *Node = getfreepos(t) /* get a free place */
		if f == nil {               /* cannot find a free place? */
			rehash(L, t, key) /* grow table */
			/* whatever called 'newkey' takes care of TM cache */
			return luaH_set(L, t, key) /* insert key into grown table */
		}
		assert(!isdummy(t))
		othern = mainposition(t, gkey(mp))
		if othern != mp { /* is colliding node out of its main position? */
			//      /* yes; move colliding node into free position */
			//      while (othern + gnext(othern) != mp)  /* find previous */
			//        othern += gnext(othern);
			//      gnext(othern) = cast_int(f - othern);  /* rechain to point to 'f' */
			//      *f = *mp;  /* copy colliding node into free pos. (mp->next also goes) */
			//      if (gnext(mp) != 0) {
			//        gnext(f) += cast_int(mp - f);  /* correct 'next' */
			//        gnext(mp) = 0;  /* now 'mp' is free */
			//      }
			//      setnilvalue(gval(mp));
		} else { /* colliding node is in its own main position */
			//      /* new node will go into free position */
			//      if (gnext(mp) != 0)
			//        gnext(f) = cast_int((mp + gnext(mp)) - f);  /* chain new position */
			//      else lua_assert(gnext(f) == 0);
			//      gnext(mp) = cast_int(f - mp);
			//      mp = f;
		}
	}
	//  setnodekey(L, &mp->i_key, key);
	//  luaC_barrierback(L, t, key);
	assert(ttisnil(gval(mp)))
	return gval(mp)
}

/*
** search function for integers
 */
func luaH_getint(t *Table, key lua_Integer) *TValue {
	/* (1 <= key && key <= t->sizearray) */
	if uint(key)-1 < t.sizearray {
		return &t.array[key-1]
	} else {
		return luaO_nilobject
	}
}

/*
** "Generic" get version. (Not that generic: not valid for integers,
** which may be in array part, nor for floats with integral values.)
 */
func getgeneric(t *Table, key *TValue) *TValue {
	var n *Node = mainposition(t, key)
	for { /* check whether 'key' is somewhere in the chain */
		if luaV_rawequalobj(gkey(n), key) {
			return gval(n) /* that's it */
		} else {
			var nx int = gnext(n)
			if nx == 0 {
				return luaO_nilobject /* not found */
			}
			n = GetNodeByOpPtr(n, nx)
		}
	}
}

/*
** main search function
 */
func luaH_get(t *Table, key *TValue) *TValue {
	switch ttype(key) {
	//    case LUA_TSHRSTR: return luaH_getshortstr(t, tsvalue(key));
	//    case LUA_TNUMINT: return luaH_getint(t, ivalue(key));
	//    case LUA_TNIL: return luaO_nilobject;
	//    case LUA_TNUMFLT: {
	//      lua_Integer k;
	//      if (luaV_tointeger(key, &k, 0)) /* index is int? */
	//        return luaH_getint(t, k);  /* use specialized version */
	//      /* else... */
	//    }  /* FALLTHROUGH */
	default:
		return getgeneric(t, key)
	}
}

/*
** beware: when using this function you probably need to check a GC
** barrier and invalidate the TM cache.
 */
func luaH_set(L *lua_State, t *Table, key *TValue) *TValue {
	var p *TValue = luaH_get(t, key)
	if p != luaO_nilobject {
		return p
	} else {
		return luaH_newkey(L, t, key)
	}
}

func luaH_setint(L *lua_State, t *Table, key lua_Integer, value *TValue) {
	var p *TValue = luaH_getint(t, key)
	var cell *TValue
	if p != luaO_nilobject {
		cell = p
	} else {
		var k TValue
		setivalue(&k, key)
		cell = luaH_newkey(L, t, &k)
	}
	setobj2t(L, cell, value)
}
