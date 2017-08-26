package main

func luaK_ret(fs *FuncState, first int, nret int) {
	luaK_codeABC(fs, OP_RETURN, first, nret+1, 0)
}

/*
** Format and emit an 'iABC' instruction. (Assertions check consistency
** of parameters versus opcode.)
 */
func luaK_codeABC(fs *FuncState, o OpCode, a int, b int, c int) int {
	return 0
}
