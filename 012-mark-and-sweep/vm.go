package main

import "errors"

const STACK_MAX = 256

type Object struct {
	Mark  int
	Next  *Object
	Value int
}

func newObject(v *VM, val int) *Object {
	// Full of objects, run GC
	if v.NumObjects == v.MaxObjects {
		v.GC()
	}

	// Create new object
	obj := &Object{Value: val}

	// Update Top of the linked list
	obj.Next = v.Top
	v.Top = obj

	// Mark the object as unreachable by default
	obj.Mark = 0

	// Update total objects of the VM
	v.NumObjects += 1
	return obj
}

// VM contains a stack and a linked list of objects
type VM struct {
	// init
	Stack      []*Object
	MaxObjects int

	// Update during Push/ Pop
	StackSize int

	// Update during NewObject/ Sweep
	Top        *Object
	NumObjects int
}

func NewVM() *VM {
	vm := &VM{Stack: make([]*Object, STACK_MAX), MaxObjects: 8}
	return vm
}

func (v *VM) PushInt(val int) error {
	// Init a new object, also append to the top of the linked list
	obj := newObject(v, val)

	// Push the object to the stack
	return v.push(obj)
}

// Simple stack implementation
func (v *VM) Pop() (*Object, error) {
	if v.StackSize < 0 {
		return nil, errors.New("stack Underflow")
	}
	obj := v.Stack[v.StackSize]
	v.StackSize -= 1
	return obj, nil
}

func (v *VM) push(val *Object) error {
	if v.StackSize > STACK_MAX {
		return errors.New("stack Overflow")
	}
	v.Stack[v.StackSize] = val
	v.StackSize += 1
	return nil
}
