package main

import (
	"os"
	"runtime/debug"
	"unsafe"
)

const (
	//cstdlib
	EXIT_SUCCESS = 0
	EXIT_FAILURE = 1
)

//strchr
func strchr(a string, c byte) int {
	temp := []byte(a)
	for i, k := range temp {
		if k == c {
			return i
		}
	}
	return -1
}

//strlen
func strlen(a []byte) size_t {
	for i, v := range a {
		if v == 0 {
			return size_t(i)
		}
	}
	return 0
}

//memcmp
func memcmp(a, b []byte, l size_t) bool {
	return string(a[:l]) == string(b[:l])
}

//三元运算符
func ITE_int(b bool, a int, c int) int {
	if b {
		return a
	}
	return c
}
func ITE_string(b bool, a string, c string) string {
	if b {
		return a
	}
	return c
}

//assert
func assert(i interface{}) {
	switch i.(type) {
	case bool:
		if i.(bool) == false {
			debug.PrintStack()
			os.Exit(1)
		}
	}
}
func abort() {
	os.Exit(1)
}

//func sizeof(i interface{}) {
//	return unsafe.Sizeof(i.(type))
//}

func GetStackByOpPtr(L *lua_State, add int, ptr *TValue) *TValue {
	v0 := unsafe.Pointer(&L.stack[0])
	vp := unsafe.Pointer(ptr)
	size := unsafe.Sizeof(L.stack[0])
	align := (uintptr(vp) - uintptr(v0)) / uintptr(size)
	return &L.stack[int(align)+add]
}

func GetNodeByOpPtr(n *Node, nx int) *Node {
	v0 := unsafe.Pointer(n)
	size := unsafe.Sizeof(n)
	ptr := uintptr(v0) + uintptr(size)*uintptr(nx)
	return (*Node)(unsafe.Pointer(ptr))
}
func TValue_Node(v *TValue) *Node {
	return (*Node)(unsafe.Pointer(v))
}

//type GCObject interface {
//	Next() GCObject
//	Tt() lu_byte
//	Marked() lu_byte
//	SetNext(GCObject)
//	SetTt(lu_byte)
//	SetMarked(lu_byte)
//}

func (self *Table) Next() GCObject {
	return self.next
}
func (self *Table) Tt() byte {
	return self.tt
}
func (self *Table) Marked() byte {
	return self.marked
}
func (self *Table) SetNext(a GCObject) {
	self.next = a
}
func (self *Table) SetTt(a byte) {
	self.tt = a
}
func (self *Table) SetMarked(a byte) {
	self.marked = a
}
func (self *TString) Next() GCObject {
	return self.next
}
func (self *TString) Tt() byte {
	return self.tt
}
func (self *TString) Marked() byte {
	return self.marked
}
func (self *TString) SetNext(a GCObject) {
	self.next = a
}
func (self *TString) SetTt(a byte) {
	self.tt = a
}
func (self *TString) SetMarked(a byte) {
	self.marked = a
}
func (self *LClosure) Next() GCObject {
	return self.next
}
func (self *LClosure) Tt() byte {
	return self.tt
}
func (self *LClosure) Marked() byte {
	return self.marked
}
func (self *LClosure) SetNext(a GCObject) {
	self.next = a
}
func (self *LClosure) SetTt(a byte) {
	self.tt = a
}
func (self *LClosure) SetMarked(a byte) {
	self.marked = a
}
func (self *Proto) Next() GCObject {
	return self.next
}
func (self *Proto) Tt() byte {
	return self.tt
}
func (self *Proto) Marked() byte {
	return self.marked
}
func (self *Proto) SetNext(a GCObject) {
	self.next = a
}
func (self *Proto) SetTt(a byte) {
	self.tt = a
}
func (self *Proto) SetMarked(a byte) {
	self.marked = a
}
