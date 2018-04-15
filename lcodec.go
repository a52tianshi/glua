package main

const MAXREGS = 255 //lua 寄存器最大数量
/*
** 获取跳转指令集的地址 (用于连续跳转)
 */
func getjump(fs *FuncState, pc int) int {
	var offset int = GETARG_sBx(fs.f.code[pc])
	if offset == NO_JUMP {
		return NO_JUMP
	} else {
		return pc + 1 + offset
	}
}

/*
** Concatenate jump-list 'l2' into jump-list 'l1'
 */
func luaK_concat(fs *FuncState, l1 *int, l2 int) {
	if l2 == NO_JUMP {
		return
	} else if *l1 == NO_JUMP {
		*l1 = l2
	} else {
		var list int = *l1
		var next int = getjump(fs, list)
		for ; next != NO_JUMP; next = getjump(fs, list) {
			list = next
		}
		//fixjump(fs, list, l2) /* last element links to 'l2' */
	}

}

/*
** Create a jump instruction and return its position, so its destination
** can be fixed later (with 'fixjump'). If there are jumps to
** this position (kept in 'jpc'), link them all together so that
** 'patchlistaux' will fix all them directly to the final destination.
 */
func luaK_jump(fs *FuncState) int {
	var jpc int = fs.jpc /* save list of jumps to here */
	var j int
	fs.jpc = NO_JUMP /* no more jumps to here */
	j = luaK_codeAsBx(fs, OP_JMP, 0, NO_JUMP)
	luaK_concat(fs, &j, jpc) /* keep them on hold */
	return j
}
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
