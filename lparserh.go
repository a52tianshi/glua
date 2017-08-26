package main

type expkind byte

const (
	VVOID  expkind = iota
	VLOCAL         /* local variable; info = local register */
)

type expdesc struct {
	k expkind
	u struct {
		ival lua_Integer
	}
	t int /* patch list of 'exit when true' */
	f int /* patch list of 'exit when false' */
}

/* description of active local variable */
type Vardesc struct {
	idx int16
}

/* description of pending goto statements and label statements */
type Labeldesc struct {
	name    *TString
	pc      int
	line    int
	nactvar lu_byte
}

/* list of labels or gotos */
type Labellist struct {
	arr  *Labeldesc
	n    int
	size int
}

/* dynamic structures used by the parser */
type Dyndata struct {
	actvar struct {
		arr  *Vardesc
		n    int
		size int
	}
	gt    Labellist
	label Labellist
}

type FuncState struct {
	f       *Proto
	prev    *FuncState
	ls      *LexState
	vl      *BlockCnt
	pc      int
	nactvar lu_byte
	nups    lu_byte
	freereg lu_byte
}
