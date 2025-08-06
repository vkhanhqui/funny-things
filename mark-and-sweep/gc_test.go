package main

import (
	"fmt"
	"testing"
)

func TestObjectsOnStackArePreserved(t *testing.T) {
	vm := NewVM()

	err := vm.PushInt(1)
	if err != nil {
		panic(err)
	}
	err = vm.PushInt(2)
	if err != nil {
		panic(err)
	}

	vm.GC()
	if vm.NumObjects == 0 {
		fmt.Println("Should have collected objects")
	}
}

func TestUnreachedObjectsAreCollected(t *testing.T) {
	vm := NewVM()

	err := vm.PushInt(1)
	if err != nil {
		panic(err)
	}
	err = vm.PushInt(2)
	if err != nil {
		panic(err)
	}
	err = vm.PushInt(3)
	if err != nil {
		panic(err)
	}

	_, err = vm.Pop()
	if err != nil {
		panic(err)
	}
	_, err = vm.Pop()
	if err != nil {
		panic(err)
	}

	// 3 and 2 are already popped, so they are unreachable when calling MarkAll (in the GC)
	vm.GC()
	if vm.NumObjects != 1 {
		fmt.Println("Should have collected some objects")
	}
}
