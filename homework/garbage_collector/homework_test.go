package main

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

func Trace(stacks [][]uintptr) []uintptr {
	traced := make(map[uintptr]struct{})

	for i, stack := range stacks {
		for j := range stack {
			ptr := stacks[i][j]
			if ptr != 0 {
				traced[ptr] = struct{}{}
			}
		}
	}

	pointers := []uintptr{}

	var scan func(pointer uintptr)
	scan = func(pointer uintptr) {
		nextPointer := *(*uintptr)(unsafe.Pointer(pointer))
		if nextPointer != 0 {
			fmt.Printf("NEXT PTR %v\n", nextPointer)
			if _, ok := traced[nextPointer]; ok {
				return
			}
			pointers = append(pointers, nextPointer)
			traced[nextPointer] = struct{}{}
			scan(nextPointer)
		}
	}

	for i, stack := range stacks {
		for j := range stack {
			ptr := stacks[i][j]
			if ptr != 0 {
				fmt.Printf("SCANNED %v\n", ptr)
				pointers = append(pointers, ptr)
				scan(ptr)
			}
		}
	}

	return pointers
}

func TestTrace(t *testing.T) {
	var heapObjects = []int{
		0x00, 0x00, 0x00, 0x00, 0x00,
	}

	var heapPointer1 *int = &heapObjects[1]
	var heapPointer2 *int = &heapObjects[2]
	var heapPointer3 *int = nil
	var heapPointer4 **int = &heapPointer3

	var stacks = [][]uintptr{
		{
			uintptr(unsafe.Pointer(&heapPointer1)), 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, uintptr(unsafe.Pointer(&heapObjects[0])),
			0x00, 0x00, 0x00, 0x00,
		},
		{
			uintptr(unsafe.Pointer(&heapPointer2)), 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, uintptr(unsafe.Pointer(&heapObjects[1])),
			0x00, 0x00, 0x00, uintptr(unsafe.Pointer(&heapObjects[2])),
			uintptr(unsafe.Pointer(&heapPointer4)), 0x00, 0x00, 0x00,
		},
		{
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, uintptr(unsafe.Pointer(&heapObjects[3])),
		},
	}

	pointers := Trace(stacks)
	expectedPointers := []uintptr{
		uintptr(unsafe.Pointer(&heapPointer1)),
		uintptr(unsafe.Pointer(&heapObjects[0])),
		uintptr(unsafe.Pointer(&heapPointer2)),
		uintptr(unsafe.Pointer(&heapObjects[1])),
		uintptr(unsafe.Pointer(&heapObjects[2])),
		uintptr(unsafe.Pointer(&heapPointer4)),
		uintptr(unsafe.Pointer(&heapPointer3)),
		uintptr(unsafe.Pointer(&heapObjects[3])),
	}
	fmt.Printf("STACKS %v\n", stacks)
	fmt.Printf("EXPECTED %v\n", expectedPointers)
	fmt.Printf("ACTUAL %v\n", pointers)

	assert.True(t, reflect.DeepEqual(expectedPointers, pointers))
}
