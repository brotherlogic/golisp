package golisp

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
)

const ()

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
	first := true
	for node != nil {
		if first {
			rep += node.value.str()
			first = false
		} else {
			rep += " " + node.value.str()
		}
		node = node.next
	}
	rep += ")"
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

// Operator is a basic operator type
type Operator struct {
	value string
}

func (o Operator) str() string {
	return o.value
}

func (i Integer) str() string {
	return fmt.Sprintf("%v", i.value)
}

//Parse parses a string to a primitive
func Parse(str string) Primitive {
	listReg, _ := regexp.Compile("[\\s\\(\\)]+")

	match, _ := regexp.MatchString("^[0-9]+$", str)
	if match {
		val, _ := strconv.Atoi(str)
		return Integer{value: val}
	}

	// Lists start with (
	if str[0] == '(' {
		elems := listReg.Split(str, -1)
		log.Printf("FROM %v to :%v (%v)", str, elems, len(elems))
		start := listNode{}
		follow := &start
		first := true
		for _, val := range elems {
			if len(val) > 0 {
				if !first {
					follow.next = &listNode{}
					follow = follow.next
				}
				follow.value = Parse(val)
				first = false
			}
		}

		return List{start: &start}
	}

	//Otherwise, assume operator
	return Operator{value: str}
}
