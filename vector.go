package main

import (
	"fmt"

	"github.com/ldsec/lattigo/ring"
)

// NewRandomVec generates a new vector of n variable uniformly distributed in Z(T)
func NewRandomVec(n, T uint64) []uint64 {
	var vector []uint64
	for i := uint64(0); i < n; i++ {
		vector = append(vector, ring.RandUniform(T, (2^37)-1))
	}
	return vector
}

// AddVec computes a+b modulo T with a,b vectors
func AddVec(a, b *[]uint64, T uint64) []uint64 {
	if len(*a) != len(*b) {
		fmt.Println("vectors should be of same size for the addition")
	}
	var vector []uint64
	for i := 0; i < len(*a); i++ {
		vector = append(vector, uint64(mod(int64((*a)[i]+(*b)[i]), int64(T))))
	}
	return vector
}

// SubVec computes a+b modulo T with a,b vectors
func SubVec(a, b *[]uint64, T uint64) []uint64 {
	if len(*a) != len(*b) {
		fmt.Println("vectors should be of same size for the substraction")
	}
	var vector []uint64
	for i := 0; i < len(*a); i++ {
		vector = append(vector, uint64(mod(int64((*a)[i]-(*b)[i]), int64(T))))
	}
	return vector
}

// MulVec computes a+b modulo T with a,b vectors
func MulVec(a, b *[]uint64, T uint64) []uint64 {
	if len(*a) != len(*b) {
		fmt.Println("vectors should be of same size for the multiplication")
	}
	var vector []uint64
	for i := 0; i < len(*a); i++ {
		vector = append(vector, uint64(mod(int64((*a)[i]*(*b)[i]), int64(T))))
	}
	return vector
}

// NegVec computes the vector -a modulo T
func NegVec(a *[]uint64, T uint64) []uint64 {
	var vector []uint64
	for i := 0; i < len(*a); i++ {
		vector = append(vector, uint64(mod(int64(-(*a)[i]), int64(T))))
	}
	return vector
}
