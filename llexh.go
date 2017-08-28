package main

const (
	FIRST_RESERVED = 257
	LUA_ENV        = "_ENV"
)
const (
	TK_AND = FIRST_RESERVED + iota
	TK_BREAK
	TK_ELSE
	TK_ELSEIF
	TK_END
	TK_RETURN
	TK_UNTIL
	TK_EOS
	TK_NAME
)

type SemInfo struct {
	r  lua_Number
	i  lua_Integer
	ts *TString
}
type Token struct {
	token   int
	seminfo SemInfo
}
type LexState struct {
	current    int
	linenumber int
	lastline   int
	t          Token
	lookahead  Token
	fs         *FuncState
	L          *lua_State
	z          *ZIO
	buff       *Mbuffer
	h          *Table
	dyd        *Dyndata
	source     *TString /* current source name */
	envn       *TString /* environment variable name */
}
