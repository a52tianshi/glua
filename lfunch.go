package main

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
