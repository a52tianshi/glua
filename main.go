package main

import (
	//"fmt"
	"io"
	//	"io/ioutil"
	"os"
	//	"time"
	"flag"

	"github.com/golang/glog"
)

const (
	LUA_PROMPT   = "> "
	LUA_PROMPT2  = ">> "
	LUA_PROGNAME = "lua"
	LUA_MAXINPUT = 512
)

func lua_stdin_is_tty() int { return 1 }

/*
** lua_readline defines how to show a prompt and then read a line from
** the standard input.
** lua_saveline defines how to "save" a read line in a "history".
** lua_freeline defines how to free a line read by lua_readline.
 */

func lua_readline(L *lua_State, b []byte, p string) int {
	io.WriteString(os.Stdout, p)
	os.Stdout.Sync() /* show prompt */
	n, _ := os.Stdin.Read(b)
	glog.Infoln(b[:n], n)
	return n
}
func lua_saveline(L *lua_State, line string) {

}

var globalL *lua_State = nil
var progname string = LUA_PROGNAME

func print_usage(badoption string) {
	lua_writestringerror("%s: ", progname)
	if badoption[1] == 'e' || badoption[1] == 'l' {
		lua_writestringerror("'%s' needs argument\n", badoption)
	} else {
		lua_writestringerror("unrecognized option '%s'\n", badoption)
	}
	lua_writestringerror(
		"usage: %s [options] [script [args]]\n"+
			"Available options are:\n"+
			"  -e stat  execute string 'stat'\n"+
			"  -i       enter interactive mode after executing 'script'\n"+
			"  -l name  require library 'name'\n"+
			"  -v       show version information\n"+
			"  -E       ignore environment variables\n"+
			"  --       stop handling options\n"+
			"  -        stop handling options and execute stdin\n",
		progname)
}

func l_message(pname string, msg string) {
	if pname != "" {
		lua_writestringerror("%s: ", pname)
	}
	lua_writestringerror("%s\n", msg)
}

func report(L *lua_State, status int) int {
	if status != LUA_OK {
		var msg string = lua_tostring(L, -1)
		l_message(progname, msg)
		lua_pop(L, 1) /* remove message */
	}
	return status
}

func print_version() {
	lua_writestring(LUA_COPYRIGHT)
	lua_writeline()
}

/*
** Create the 'arg' table, which stores all arguments from the
** command line ('argv'). It should be aligned so that, at index 0,
** it has 'argv[script]', which is the script name. The arguments
** to the script (everything after 'script') go to positive indices;
** other arguments (before the script name) go to negative indices.
** If there is no script name, assume interpreter's name as base.
 */
func createargtable(L *lua_State, argv []string, argc int, script int) {
	var i, narg int
	//  if (script == argc) script = 0;  /* no script name? */
	//  narg = argc - (script + 1);  /* number of positive indices */
	lua_createtable(L, narg, script+1)
	for i = 0; i < argc; i++ {
		//    lua_pushstring(L, argv[i]);
		//    lua_rawseti(L, -2, i - script);
	}
	//  lua_setglobal(L, "arg");
}

/*
** Returns the string to be used as a prompt by the interpreter.
 */
func get_prompt(L *lua_State, firstline int) string {
	var p string
	lua_getglobal(L, ITE_string(firstline != 0, "_PROMPT", "_PROMPT2"))
	p = lua_tostring(L, -1)
	if p == "" {
		p = ITE_string(firstline != 0, LUA_PROMPT, LUA_PROMPT2)
	}
	return p
}

/*
** Prompt the user, read a line, and push it into the Lua stack.
 */
func pushline(L *lua_State, firstline int) int {
	var buffer [LUA_MAXINPUT]byte
	var b []byte = buffer[:]
	var l size_t
	var prmt string = get_prompt(L, firstline)
	var readstatus = lua_readline(L, b, prmt)
	if readstatus == 0 {
		return 0
	}
	lua_pop(L, 1)
	l = strlen(b)
	if l > 0 && b[l-1] == '\n' {
		l--
		b[l] = 0 /* remove it */
	}
	if firstline != 0 && b[0] == '=' {
		lua_pushfstring(L, "return %s", string(b[1:readstatus]))
	} else {
		lua_pushlstring(L, b, l)
	}
	return 1
}

/*
** Try to compile line on the stack as 'return <line>;'; on return, stack
** has either compiled chunk or original line (if compilation failed).
 */
func addreturn(L *lua_State) int {
	var line string = lua_tostring(L, -1)
	var retline string = lua_pushfstring(L, "return %s;", line)
	glog.Infoln("cqadd", line)
	glog.Infoln("cqadd", retline)
	var status int = luaL_loadbuffer(L, retline, size_t(len(retline)), "=stdin")
	if status == LUA_OK {
		lua_remove(L, -2) /* remove modified line */
		if line != "" {
			lua_saveline(L, line)
		}
	} else {
		lua_pop(L, 2)
	}
	return status
}

/*
** Read a line and try to load (compile) it first as an expression (by
** adding "return " in front of it) and second as a statement. Return
** the final status of load/call with the resulting function (if any)
** in the top of the stack.
 */
func loadline(L *lua_State) int {
	var status int
	lua_settop(L, 0)
	if pushline(L, 1) == 0 {
		return -1 /* no input */
	}
	status = addreturn(L)
	if status != LUA_OK { /* 'return ...' did not work? */
		//status = multiline(L) /* try as command, maybe with continuation lines */
	}
	lua_remove(L, 1) /* remove line from the stack */
	glog.Infoln("cqtest %d", lua_gettop(L))
	assert(lua_gettop(L) == 1)
	return status
}

/*
** Prints (calling the Lua 'print' function) any values on the stack
 */
func l_print(L *lua_State) {
	var n int = lua_gettop(L)
	if n > 0 { /* any result to be printed? */
		luaL_checkstack(L, LUA_MINSTACK, "too many results to print")
		lua_getglobal(L, "print")
		lua_insert(L, 1)
		if lua_pcall(L, n, 0, 0) != LUA_OK {
			l_message(progname, lua_pushfstring(L, "error calling 'print' (%s)", lua_tostring(L, -1)))
		}
	}
}

/*
** Do the REPL: repeatedly read (load) a line, evaluate (call) it, and
** print any results.
 */
func doREPL(L *lua_State) {
	var status int
	var oldprogname string = progname
	progname = "" /* no 'progname' on errors in interactive mode */
	for status = loadline(L); status != -1; status = loadline(L) {
		if status == LUA_OK {
			//      status = docall(L, 0, LUA_MULTRET);
		}
		if status == LUA_OK {
			l_print(L)
		} else {
			report(L, status)
		}
	}
	lua_settop(L, 0) /* clear stack */
	lua_writeline()
	progname = oldprogname
}

/* bits of various argument indicators in 'args' */
const (
	has_error = 1  /* bad option */
	has_i     = 2  /* -i */
	has_v     = 4  /* -v */
	has_e     = 8  /* -e */
	has_E     = 16 /* -E */
)

//检查初始的启动命令参数
func collectargs(argv []string, first *int) int {
	var args int = 0
	var i int
	for i = 1; len(argv) > i; i++ {
		*first = i
		if argv[i][0] != '-' {
			return args
		}
		switch argv[i][1] {
		case '-':
			if argv[i][2] != 0 {
				return has_error
			}
			*first = i + 1
			return args
		case 0:
			return args
		case 'v':
			if len(argv[i]) > 2 {
				return has_error
			}
			args |= has_v
			break
		default:
			return has_error
		}
	}
	*first = i
	return args
}

/*
** Main body of stand-alone interpreter (to be called in protected mode).
** Reads the options and handles them all.
 */
func pmain(L *lua_State) int {
	var argc int = int(lua_tointeger(L, 1))
	var argv []string = lua_touserdata(L, 2).([]string)
	var script int
	var args int = collectargs(argv, &script)
	luaL_checkversion(L)
	if argv[0] != "" && argv[0][0] != 0 {
		progname = argv[0]
	}
	if args == has_error {
		print_usage(argv[script])
		return 0
	}
	if (args & has_v) != 0 {
		print_version()
	}
	luaL_openlibs(L)                      /* open standard libraries */
	createargtable(L, argv, argc, script) /* create table 'arg' */

	if (args & has_i) != 0 {
		doREPL(L)
	} else if (script == argc) && ((args & (has_e | has_v)) == 0) {
		if lua_stdin_is_tty() != 0 {
			print_version()
			doREPL(L)
		} else {

		}
	}
	return 1
}
func main() {
	flag.Set("logtostderr", "true")
	flag.Parse()
	var status, result int
	L := luaL_newstate()
	if L == nil {
		l_message(os.Args[0], "cannot create state: not enough memory")
		os.Exit(EXIT_FAILURE)
	}
	lua_pushcfunction(L, pmain)
	lua_pushinteger(L, lua_Integer(len(os.Args)))
	lua_pushlightuserdata(L, os.Args)
	status = lua_pcall(L, 2, 1, 0) /* do the call */
	result = lua_toboolean(L, -1)  /* get result */
	report(L, status)
	lua_close(L)
	os.Exit(ITE_int(result != 0 && status == LUA_OK, EXIT_SUCCESS, EXIT_FAILURE))
}
