package main

//complete
const (
	ALPHABIT  = 0
	DIGITBIT  = 1
	PRINTBIT  = 2
	SPACEBIT  = 3
	XDIGITBIT = 4
)

func MASK(B byte) int {
	return 1 << B
}
func testprop(c int, p int) bool {
	return (int(luai_ctype_[(c)+1]) & (p)) != 0
}

func lislalpha(c int) bool { return testprop(c, MASK(ALPHABIT)) }
func lislalnum(c int) bool { return testprop(c, (MASK(ALPHABIT) | MASK(DIGITBIT))) }
func lisdigit(c int) bool  { return testprop(c, MASK(DIGITBIT)) }
func lisspace(c int) bool  { return testprop(c, MASK(SPACEBIT)) }
func lisprint(c int) bool  { return testprop(c, MASK(PRINTBIT)) }
func lisxdigit(c int) bool { return testprop(c, MASK(XDIGITBIT)) }

func ltolower(c int) int {
	return c | ('A' ^ 'a')
}
