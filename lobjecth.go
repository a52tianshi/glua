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
	tt     byte
	marked byte
}

type GCObject interface {
	Next() GCObject
	Tt() byte
	Marked() byte
	SetNext(GCObject)
	SetTt(byte)
	SetMarked(byte)
}

type Value struct { //union
	gc GCObject
	p  interface{}
	b  bool
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
func checktype(o *TValue, t int) bool {
	return ttnov(o) == t
}
func ttisfloat(o *TValue) bool {
	return checktag(o, LUA_TNUMFLT)
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
func ttisstring(o *TValue) bool {
	return checktype(o, LUA_TSTRING)
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
func ttisdeadkey(o *TValue) bool {
	return checktag((o), LUA_TDEADKEY)
}

/* Macros to access values */
func ivalue(o *TValue) lua_Integer {
	assert(ttisinteger(o))
	return o.value_.i
}
func fltvalue(o *TValue) lua_Number {
	assert(ttisfloat(o))
	return o.value_.n
}
func gcvalue(o *TValue) GCObject {
	assert(iscollectable(o))
	return o.value_.gc
}
func pvalue(o *TValue) interface{} {
	assert(ttislightuserdata(o))
	return o.value_.p
}
func tsvalue(o *TValue) *TString {
	assert(ttisstring(o))
	return gco2ts(o.value_.gc)
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
func bvalue(o *TValue) bool {
	assert(ttisboolean(o))
	return o.value_.b
}

/* a dead value may get the 'gc' field, but cannot access its contents */
func l_isfalse(o *TValue) bool {
	return ttisnil(o) || (ttisboolean(o) && bvalue(o) == false)
}

//判断是否是需要gc的类型
func iscollectable(o *TValue) bool { return (rttype(o) & BIT_ISCOLLECTABLE) != 0 }

/* Macros for internal tests */
func righttt(obj *TValue) bool {
	return ttype(obj) == int(gcvalue(obj).Tt())
}
func checkliveness(L *lua_State, obj *TValue) {
	assert(!iscollectable(obj) || (righttt(obj) && (L == nil || !isdead(L.l_G, gcvalue(obj)))))
}

//赋值TValue
func settt_(o *TValue, t int) {
	o.tt_ = t
}
func setfltvalue(obj *TValue, x lua_Number) {
	obj.value_.n = x
	settt_(obj, LUA_TNUMFLT)
}
func setivalue(obj *TValue, x lua_Integer) {
	obj.value_.i = x
	settt_(obj, LUA_TNUMINT)
}
func setnilvalue(obj *TValue) {
	settt_(obj, LUA_TNIL)
}
func setfvalue(obj *TValue, x lua_CFunction) {
	obj.value_.f = x
	settt_(obj, LUA_TLCF)
}
func setpvalue(obj *TValue, x interface{}) {
	obj.value_.p = x
	settt_(obj, LUA_TLIGHTUSERDATA)
}
func setbvalue(obj *TValue, x bool) {
	obj.value_.b = x
	settt_(obj, LUA_TBOOLEAN)
}
func setsvalue(L *lua_State, obj *TValue, x *TString) {
	var io *TValue = obj
	var x_ *TString = x
	io.value_.gc = obj2gco(x_)
	settt_(io, ctb(int(x_.Tt())))
	checkliveness(L, io)
}
func setclLvalue(L *lua_State, obj *TValue, x *LClosure) {
	var io *TValue = obj
	var x_ *LClosure = x
	io.value_.gc = obj2gco(x_)
	settt_(io, ctb(LUA_TLCL))
	checkliveness(L, io)
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

/*
** different types of assignments, according to destination
 */

/* from stack to (same) stack */
func setobjs2s(L *lua_State, obj1 *TValue, obj2 *TValue) {
	setobj(L, obj1, obj2)
}

/* to stack (not from same stack) */
func setobj2s(L *lua_State, obj1 *TValue, obj2 *TValue) {
	setobj(L, obj1, obj2)
}

func setsvalue2s(L *lua_State, idx int, ts *TString) {
	setsvalue(L, &L.stack[idx], ts)
}

func setobjt2t(L *lua_State, obj1 *TValue, obj2 *TValue) {
	setobj(L, obj1, obj2)
}

/* to table (define it as an expression to be used in macros) */
func setobj2t(L *lua_State, o1 *TValue, o2 *TValue) {
	*o1 = *o2
	checkliveness(L, (o1))
}

type StkId *TValue /* index to stack elements */

/*
** Header for string value; string bytes follow the end of this structure
** (aligned according to 'UTString'; see next).
 */

type TString struct {
	CommonHeader
	data   string
	extra  byte //cqtest
	shrlen byte
	hash   uint
	u      struct {
		lnglen size_t   //长字符串的长度
		hnext  *TString //指向hash链表
	}
}

/*
** Get the actual string (array of bytes) from a 'TString'.
** (Access to 'extra' ensures that value is really a 'TString'.)
 */
//不检验extra
func getstr(ts *TString) string {
	return ts.data
}
func svalue(o *TValue) string {
	return getstr(tsvalue(o))
}
func tsslen(s *TString) int {
	return len(s.data)
}
func vslen(o *TValue) size_t {
	return size_t(tsslen(tsvalue(o)))
}

/*
** Description of an upvalue for function prototypes
 */
type Upvaldesc struct {
	name    *TString /* upvalue name (for debug information) */
	instack byte     /* whether it is in stack (register) */
	idx     byte     /* index of upvalue (in stack or in outer function's list) */
}

/*
** 函数中的局部变量
** (used for debug information)
 */
type LocVar struct {
	varname *TString
	startpc int /* first point where variable is active */
	endpc   int /* first point where variable is dead */
}

/*
** Function Prototypes
 */

type Proto struct {
	CommonHeader
	numparams       byte
	is_vararg       byte
	maxstacksize    byte
	sizeupvalues    int
	sizek           int
	sizecode        int
	sizelineinfo    int
	linedefined     int
	lastlinedefined int
	k               *TValue
	code            *Instruction
	//...
	lineinfo []int
	//...
	upvalues []Upvaldesc /* upvalue information */
	source   *TString
}

/*
** Closures
 */

type ClosureHeader struct {
	CommonHeader
	nupvalues byte
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
	flags     byte //元方法标记
	lsizenode byte //取了对数
	sizearray uint
	array     []TValue
	node      []Node
	lastfree  int    //node的索引 -1表示空
	metatable *Table //元表
	gclist    GCObject
}

/*
** 'module' operation for hashing (size is always a power of 2)
 */

func lmod(s uint, size int) int {
	assert(size&(size-1) == 0)
	return int(s & uint(size-1))
}
func twoto(x byte) size_t {
	return size_t(1) << x
}
func sizenode(t *Table) size_t {
	return twoto(t.lsizenode)
}

/*
** (address of) a fixed nil value
 */
var luaO_nilobject *TValue = &luaO_nilobject_

var luaO_nilobject_ TValue
