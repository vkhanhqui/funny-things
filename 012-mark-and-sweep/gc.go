package main

import (
	"fmt"
)

// Mark the Stack as reachable
func (v *VM) MarkAll() {
	for i := 0; i < v.StackSize; i++ {
		v.Stack[i].Mark = 1
	}
}

func (v *VM) Sweep() {
	object := v.Top
	prevObj := &Object{}

	// Traverse the linked list of objects
	for object != nil {
		if object.Mark != 1 { // Case 1: Unreachable object, Simply remove the object
			if prevObj.Next == nil { // Case: First object
				v.Top = object.Next
			} else {
				prevObj.Next = object.Next
			}
			v.NumObjects -= 1
		} else { // Case 2: Reachable object
			object.Mark = 0  // This object was reached, so unmark it (for the next GC)
			prevObj = object // Processed, set it as previous object
		}
		object = object.Next
	}
}

func (v *VM) GC() {
	numObjects := v.NumObjects
	v.MarkAll()
	v.Sweep()
	v.MaxObjects = v.NumObjects * 2 // Dynamically adjust the maximum
	fmt.Printf("Collected %d objects , %d remaining \n", numObjects-v.NumObjects, v.NumObjects)
}
