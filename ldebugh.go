package main

import (
	"unsafe"
)

func pcRel(pc *Instruction, p *Proto) int {
	return int((uintptr(unsafe.Pointer(pc)) - uintptr(unsafe.Pointer(p.code))/8)) - 1
}

func getfuncline(f *Proto, pc int) int {
	if f.lineinfo != nil {
		return f.lineinfo[pc]
	} else {
		return -1
	}
}
