package main

const MAXUPVAL = 255

type UpVal struct {
	v        *TValue
	refcount lu_mem
	u        struct {
		open struct {
			next    *UpVal
			touched int
		}
		value TValue
	}
}

func upisopen(up *UpVal) bool {
	return up.v != &up.u.value
}
