package main

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

/* returns the key, given the value of a table entry */
func keyfromval(v *TValue) *TValue {
	return gkey(TValue_Node(v))
}
