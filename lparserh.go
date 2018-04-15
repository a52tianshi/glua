package main

type expkind byte

const (
	VVOID  expkind = iota
	VLOCAL         /* local variable; info = local register */
)

type expdesc struct {
	k expkind
	u struct {
		ival lua_Integer /* for VKINT */
		nval lua_Number  /* for VKFLT */
		info int         /* for generic use */
		ind  struct {    /* for indexed variables (VINDEXED) */
			idx int16 /* index (R/K) */
			t   byte  /* table (register or upvalue) */
			vt  byte  /* whether 't' is register (VLOCAL) or upvalue (VUPVAL) */
		}
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
	nactvar byte
}

/* list of labels or gotos */
type Labellist struct {
	arr  []Labeldesc
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
	bl      *BlockCnt
	pc      int
	jpc     int
	nactvar byte
	nups    byte
	freereg byte
}
