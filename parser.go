package golisp

import "fmt"

// Primitive is the base of all elements in source code
type Primitive interface {
	str() string
}

// List Basic Type
type List struct {
	start *listNode
}

func (l List) str() string {
	rep := "("
	node := l.start
	for node != nil {
		rep += " " + node.value.str()
		node = node.next
	}
	rep += " )"
	return rep
}

type listNode struct {
	value Primitive
	next  *listNode
}

// Integer is a base value of int type
type Integer struct {
	value int
}

func (i Integer) str() string {
	return fmt.Sprintf("%v", i.value)
}
