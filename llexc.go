package main

func next(ls *LexState) {
	ls.current = zgetc(ls.z)
}
func currIsNewline(ls *LexState) bool {
	return (ls.current == '\n' || ls.current == '\r')
}
func luaX_token2str(ls *LexState, token int) string {
	return ""
	//	  if (token < FIRST_RESERVED) {  /* single-byte symbols? */
	//    lua_assert(token == cast_uchar(token));
	//    return luaO_pushfstring(ls->L, "'%c'", token);
	//  }
	//  else {
	//    const char *s = luaX_tokens[token - FIRST_RESERVED];
	//    if (token < TK_EOS)  /* fixed format (symbols and reserved words)? */
	//      return luaO_pushfstring(ls->L, "'%s'", s);
	//    else  /* names, strings, and numerals */
	//      return s;
	//  }
}
func lexerror(ls *LexState, msg string, token int) {
	//  msg = luaG_addinfo(ls->L, msg, ls->source, ls->linenumber);
	//  if (token)
	//    luaO_pushfstring(ls->L, "%s near %s", msg, txtToken(ls, token));
	//  luaD_throw(ls->L, LUA_ERRSYNTAX);
}
func luaX_syntaxerror(ls *LexState, msg string) {
	lexerror(ls, msg, ls.t.token)
}
func inclinenumber(ls *LexState) {
	var old int = ls.current
	assert(currIsNewline(ls))
	next(ls) /* skip '\n' or '\r' */
	if currIsNewline(ls) && ls.current != old {
		next(ls) /* skip '\n\r' or '\r\n' */
	}
	//  if (++ls->linenumber >= MAX_INT)
	//    lexerror(ls, "chunk has too many lines", 0);
}
func luaX_setinput(L *lua_State, ls *LexState, z *ZIO, source *TString, firstchar int) {
	ls.t.token = 0
	ls.L = L
	ls.current = firstchar
	ls.lookahead.token = TK_EOS /* no look-ahead token */
	ls.z = z
	ls.fs = nil
	ls.linenumber = 1
	ls.lastline = 1
	ls.source = source
	ls.envn = luaS_newliteral(L, LUA_ENV)           /* get env name */
	luaZ_resizebuffer(ls.L, ls.buff, LUA_MINBUFFER) /* initialize buffer */
}
func llex(ls *LexState, seminfo *SemInfo) int {
	luaZ_resetbuffer(ls.buff)
	for {
		switch ls.current {
		case '\n', '\r': /* line breaks */
			inclinenumber(ls)
		}
	}
	return 0
}
func luaX_next(ls *LexState) {
	ls.lastline = ls.linenumber
	if ls.lookahead.token != TK_EOS { /* is there a look-ahead token? */
		ls.t = ls.lookahead         /* use this one */
		ls.lookahead.token = TK_EOS /* and discharge it */
	} else {
		ls.t.token = llex(ls, &ls.t.seminfo) /* read next token */
	}
}
