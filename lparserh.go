package main

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
