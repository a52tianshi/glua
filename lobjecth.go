package main

const (
	LUA_TPROTO   = LUA_NUMTAGS
	LUA_TDEADKEY = LUA_NUMTAGS + 1

	LUA_TOTALTAGS = LUA_TPROTO + 2

	LUA_TLCL = LUA_TFUNCTION | (0 << 4)
	LUA_TLCF = LUA_TFUNCTION | (1 << 4) /* light C function */
	LUA_TCCL = LUA_TFUNCTION | (2 << 4)

	LUA_TSHRSTR = LUA_TSTRING | 0<<4
	LUA_TLNGSTR = LUA_TSTRING | 1<<4

	LUA_TNUMFLT = LUA_TNUMBER | 0<<4
	LUA_TNUMINT = LUA_TNUMBER | 1<<4

	BIT_ISCOLLECTABLE = 1 << 6
)

/* mark a tag as collectable */

func ctb(t int) int {
	return t | BIT_ISCOLLECTABLE
}

type CommonHeader struct {
	next   GCObject
	tt     lu_byte
	marked lu_byte
}

//type GCObject struct {
//	CommonHeader
//}
type GCObject interface {
	Next() GCObject
	Tt() lu_byte
	Marked() lu_byte
	SetNext(GCObject)
	SetTt(lu_byte)
	SetMarked(lu_byte)
}

type Value struct { //union
	gc GCObject
	p  interface{}
	b  int
	f  lua_CFunction
	i  lua_Integer
	n  lua_Number
}

type TValuefields struct {
	value_ Value
	tt_    int
}

type lua_TValue struct {
	TValuefields
}

//func val_(o *TValue) Value {
//	return o.value_
//}

type TValue lua_TValue

//func val_(o *TValue) int {
//	return o.value_
//}

/* raw type tag of a TValue */
func rttype(o *TValue) int {
	return o.tt_
}

/* tag with no variants (bits 0-3) */
func novariant(x int) int {
	return ((x) & 0x0F)
}

/* type tag of a TValue (bits 0-3 for tags + variant bits 4-5) */
func ttype(o *TValue) int {
	return (rttype(o) & 0x3F)
}

/* type tag of a TValue with no variants (bits 0-3) */
func ttnov(o *TValue) int {
	return novariant(rttype(o))
}

/* Macros to test type */
func checktag(o *TValue, t int) bool {
	return rttype(o) == t
}
func ttisinteger(o *TValue) bool {
	return checktag(o, LUA_TNUMINT)
}
func ttisnil(o *TValue) bool {
	return checktag(o, LUA_TNIL)
}
func ttisboolean(o *TValue) bool {
	return checktag(o, LUA_TBOOLEAN)
}
func ttislightuserdata(o *TValue) bool {
	return checktag(o, LUA_TLIGHTUSERDATA)
}
func ttistable(o *TValue) bool {
	return checktag(o, ctb(LUA_TTABLE))
}
func ttisLclosure(o *TValue) bool {
	return checktag(o, ctb(LUA_TLCL))
}
func ttislcf(o *TValue) bool {
	return checktag(o, LUA_TLCF)
}

/* Macros to access values */
func ivalue(o *TValue) lua_Integer {
	assert(ttisinteger(o))
	return o.value_.i
}
func gcvalue(o *TValue) GCObject {
	assert(iscollectable(o))
	return o.value_.gc
}
func pvalue(o *TValue) interface{} {
	assert(ttislightuserdata(o))
	return o.value_.p
}
func clLvalue(o *TValue) *LClosure {
	assert(ttisLclosure(o))
	return gco2lcl(o.value_.gc)
}
func fvalue(o *TValue) lua_CFunction {
	assert(ttislcf(o))
	return o.value_.f
}
func hvalue(o *TValue) *Table {
	assert(ttistable(o))
	return gco2t(o.value_.gc)
}
func bvalue(o *TValue) int {
	assert(ttisboolean(o))
	return o.value_.b
}

/* a dead value may get the 'gc' field, but cannot access its contents */
func l_isfalse(o *TValue) bool {
	return ttisnil(o) || (ttisboolean(o) && bvalue(o) == 0)
}

func iscollectable(o *TValue) bool {
	return (rttype(o) & BIT_ISCOLLECTABLE) != 0
}

/* Macros for internal tests */
func righttt(obj *TValue) bool {
	return ttype(obj) == int(gcvalue(obj).Tt())
}
func checkliveness(L *lua_State, obj *TValue) {
	assert(!iscollectable(obj) || (righttt(obj) && (L == nil || !isdead(L.l_G, gcvalue(obj)))))
}

/* Macros to set values */
func settt_(o *TValue, t int) {
	o.tt_ = t
}
func setivalue(obj *TValue, x lua_Integer) {
	var io *TValue = obj
	io.value_.i = x
	settt_(io, LUA_TNUMINT)
}
func setnilvalue(obj *TValue) {
	settt_(obj, LUA_TNIL)
}
func setfvalue(obj *TValue, x lua_CFunction) {
	var io *TValue = obj
	io.value_.f = x
	settt_(io, LUA_TLCF)
}
func setpvalue(obj *TValue, x interface{}) {
	var io *TValue = obj
	io.value_.p = x
	settt_(io, LUA_TLIGHTUSERDATA)
}
func sethvalue(L *lua_State, obj *TValue, x *Table) {
	var io *TValue = obj
	var x_ *Table = x
	io.value_.gc = obj2gco(x_)
	settt_(io, ctb(LUA_TTABLE))
	checkliveness(L, io)
}
func setobj(L *lua_State, obj1 *TValue, obj2 *TValue) {
	var io1 *TValue = obj1
	*io1 = *(obj2)
	checkliveness(L, io1)
}

type StkId *TValue /* index to stack elements */

/*
** Header for string value; string bytes follow the end of this structure
** (aligned according to 'UTString'; see next).
 */

type TString struct {
	CommonHeader
	data string
	//	extra  lu_byte
	//	shrlen lu_byte
	//	hash   uint
	//	u      struct {
	//		lnglen size_t
	//		hnext  *TString
	//	}
}

/*
** Get the actual string (array of bytes) from a 'TString'.
** (Access to 'extra' ensures that value is really a 'TString'.)
 */
//不检验extra
func getstr(ts *TString) string {
	return ts.data
}

/*
** Function Prototypes
 */

type Proto struct {
	CommonHeader
	numparams    lu_byte
	is_vararg    lu_byte
	maxstacksize lu_byte
	sizeupvalues int
	sizek        int
	sizecode     int
	sizelineinfo int
	//...
}

/*
** Closures
 */

type ClosureHeader struct {
	CommonHeader
	nupvalues lu_byte
	gclist    *GCObject
}

type LClosure struct {
	ClosureHeader
	p      *Proto
	upvals [1]*UpVal //[1];  /* list of upvalues */
}

/*
** Tables
 */
type TKey struct {
	nk struct {
		TValuefields
		next int
	}
	tvk TValue
}

type Node struct {
	i_val TValue
	i_key TKey
}

type Table struct {
	CommonHeader
	flags     lu_byte
	lsizenode lu_byte
	sizearray uint
	array     []TValue
	node      *Node
	lastfree  *Node
	metatable *Table
	//	gclist    *GCObject
}

/*
** 'module' operation for hashing (size is always a power of 2)
 */

func lmod(s uint, size int) int {
	assert(size&(size-1) == 0)
	return int(s & uint(size-1))
}

/*
** (address of) a fixed nil value
 */
var luaO_nilobject *TValue = &luaO_nilobject_

var luaO_nilobject_ TValue
