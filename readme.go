package main

//所有变量变为64位为基础

//G(L) 直接拆开不用宏

//assert 直接用  不用封装了lua_assert  lua_longassert

//#define val_(o)		((o)->value_)   直接拆开不用宏

//luaC_newobj 改写

//LUA_NOCVTN2S 不定义  不开启

//luaO_pushvfstring
//markobject
//luaY_parser
