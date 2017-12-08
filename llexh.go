package main

const (
	FIRST_RESERVED = 257
	LUA_ENV        = "_ENV"
)
const (
	TK_AND = FIRST_RESERVED + iota
	TK_BREAK
	TK_DO
	TK_ELSE
	TK_ELSEIF
	TK_FALSE
	TK_FOR
	TK_FUNCTION
	TK_GOTO
	TK_IF
	TK_IN
	TK_LOCAL
	TK_NIL
	TK_NOT
	TK_OR
	TK_REPEAT
	TK_END
	TK_RETURN
	TK_THEN
	TK_TRUE
	TK_UNTIL
	TK_WHILE
	//////////////
	TK_IDIV
	TK_CONCAT
	TK_DOTS
	TK_EQ
	TK_GE
	TK_LE
	TK_NE
	TK_SHL
	TK_SHR
	TK_DBCOLON
	TK_EOS
	TK_FLT
	TK_INT
	TK_NAME
	TK_STRING
)

const NUM_RESERVED = int(TK_WHILE - FIRST_RESERVED + 1)

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
