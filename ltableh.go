package main

func gnode(t *Table, i int) *Node {
	return &t.node[i]
}
func gval(n *Node) *TValue {
	return &n.i_val
}
func gnext(n *Node) int {
	return n.i_key.nk.next
}

/* 'const' to avoid wrong writings that can mess up field 'next' */
func gkey(n *Node) *TValue {
	return &n.i_key.tvk
}
func isdummy(t *Table) bool {
	return t.lastfree != -1
}
func allocsizenode(t *Table) int {
	if isdummy(t) {
		return 0
	} else {
		return int(sizenode(t))
	}
}

/* returns the key, given the value of a table entry */
func keyfromval(v *TValue) *TValue {
	return gkey(TValue_Node(v))
}
