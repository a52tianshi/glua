package main

import (
	"math"
)

//基础指令集 格式!
type OpMode byte

const (
	iABC OpMode = iota
	iABx
	iAsBx
	iAx
)

//各种指令集的位数 与偏置(位置)
const (
	SIZE_C  = 9
	SIZE_B  = 9
	SIZE_Bx = (SIZE_C + SIZE_B)
	SIZE_A  = 8
	SIZE_Ax = (SIZE_C + SIZE_B + SIZE_A)

	SIZE_OP = 6

	POS_OP = 0
	POS_A  = (POS_OP + SIZE_OP)
	POS_C  = (POS_A + SIZE_A)
	POS_B  = (POS_C + SIZE_C)
	POS_Bx = POS_C
	POS_Ax = POS_A
)

//LUAI_BITSINT 是32位的 (int 是32位的)所以下面这样定义
const (
	MAXARG_Bx  = math.MaxInt32
	MAXARG_sBx = math.MaxInt32
)

/* creates a mask with 'n' 1 bits at position 'p' */
//  n = 4 p = 2
//   0000111100
func MASK1(size byte, p byte) Instruction {
	//	((~((~(Instruction)0)<<(n)))<<(p))
	return (1<<size - 1) << p
}
func MASK0(size byte, p byte) Instruction {
	return Instruction(math.MaxUint64 - MASK1(size, p))
}
func getarg(i Instruction, pos, size byte) int {
	return int((i >> pos) & MASK1(size, 0))
}
func GETARG_Bx(i Instruction) int {
	return getarg(i, POS_Bx, SIZE_Bx)
}
func GETARG_sBx(i Instruction) int {
	return GETARG_Bx(i) - MAXARG_sBx
}

type OpCode byte

const (
	OP_MOVE     OpCode = iota /*	A B	R(A) := R(B)					*/
	OP_LOADK                  /*	A Bx	R(A) := Kst(Bx)					*/
	OP_LOADKX                 /*	A 	R(A) := Kst(extra arg)				*/
	OP_LOADBOOL               /*	A B C	R(A) := (Bool)B; if (C) pc++			*/
	OP_LOADNIL                /*	A B	R(A), R(A+1), ..., R(A+B) := nil		*/
	OP_GETUPVAL               /*	A B	R(A) := UpValue[B]				*/

	OP_GETTABUP /*	A B C	R(A) := UpValue[B][RK(C)]			*/
	OP_GETTABLE /*	A B C	R(A) := R(B)[RK(C)]				*/

	OP_SETTABUP /*	A B C	UpValue[A][RK(B)] := RK(C)			*/
	OP_SETUPVAL /*	A B	UpValue[B] := R(A)				*/
	OP_SETTABLE /*	A B C	R(A)[RK(B)] := RK(C)				*/

	OP_NEWTABLE /*	A B C	R(A) := {} (size = B,C)				*/

	OP_SELF /*	A B C	R(A+1) := R(B); R(A) := R(B)[RK(C)]		*/

	OP_ADD  /*	A B C	R(A) := RK(B) + RK(C)				*/
	OP_SUB  /*	A B C	R(A) := RK(B) - RK(C)				*/
	OP_MUL  /*	A B C	R(A) := RK(B) * RK(C)				*/
	OP_MOD  /*	A B C	R(A) := RK(B) % RK(C)				*/
	OP_POW  /*	A B C	R(A) := RK(B) ^ RK(C)				*/
	OP_DIV  /*	A B C	R(A) := RK(B) / RK(C)				*/
	OP_IDIV /*	A B C	R(A) := RK(B) // RK(C)				*/
	OP_BAND /*	A B C	R(A) := RK(B) & RK(C)				*/
	OP_BOR  /*	A B C	R(A) := RK(B) | RK(C)				*/
	OP_BXOR /*	A B C	R(A) := RK(B) ~ RK(C)				*/
	OP_SHL  /*	A B C	R(A) := RK(B) << RK(C)				*/
	OP_SHR  /*	A B C	R(A) := RK(B) >> RK(C)				*/
	OP_UNM  /*	A B	R(A) := -R(B)					*/
	OP_BNOT /*	A B	R(A) := ~R(B)					*/
	OP_NOT  /*	A B	R(A) := not R(B)				*/
	OP_LEN  /*	A B	R(A) := length of R(B)				*/

	OP_CONCAT /*	A B C	R(A) := R(B).. ... ..R(C)			*/

	OP_JMP /*	A sBx	pc+=sBx; if (A) close all upvalues >= R(A - 1)	*/
	OP_EQ  /*	A B C	if ((RK(B) == RK(C)) ~= A) then pc++		*/
	OP_LT  /*	A B C	if ((RK(B) <  RK(C)) ~= A) then pc++		*/
	OP_LE  /*	A B C	if ((RK(B) <= RK(C)) ~= A) then pc++		*/

	OP_TEST    /*	A C	if not (R(A) <=> C) then pc++			*/
	OP_TESTSET /*	A B C	if (R(B) <=> C) then R(A) := R(B) else pc++	*/

	OP_CALL     /*	A B C	R(A), ... ,R(A+C-2) := R(A)(R(A+1), ... ,R(A+B-1)) */
	OP_TAILCALL /*	A B C	return R(A)(R(A+1), ... ,R(A+B-1))		*/
	OP_RETURN   /*	A B	return R(A), ... ,R(A+B-2)	(see note)	*/

	OP_FORLOOP /*	A sBx	R(A)+=R(A+2);
		if R(A) <?= R(A+1) then { pc+=sBx; R(A+3)=R(A) }*/
	OP_FORPREP /*	A sBx	R(A)-=R(A+2); pc+=sBx				*/

	OP_TFORCALL /*	A C	R(A+3), ... ,R(A+2+C) := R(A)(R(A+1), R(A+2));	*/
	OP_TFORLOOP /*	A sBx	if R(A+1) ~= nil then { R(A)=R(A+1); pc += sBx }*/

	OP_SETLIST /*	A B C	R(A)[(C-1)*FPF+i] := R(A+i), 1 <= i <= B	*/

	OP_CLOSURE /*	A Bx	R(A) := closure(KPROTO[Bx])			*/

	OP_VARARG /*	A B	R(A), R(A+1), ..., R(A+B-2) = vararg		*/

	OP_EXTRAARG /*	Ax	extra (larger) argument for previous opcode	*/
)
