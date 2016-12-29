package golisp

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

const ()

// Primitive is the base of all elements in source code
type Primitive interface {
	str() string
	isList() bool
}

// List Basic Type
type List struct {
	start *listNode
}

func (l List) isList() bool {
	return true
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

func (i Integer) str() string {
	return fmt.Sprintf("%v", i.value)
}

func (i Integer) isList() bool {
	return false
}

// Truth is a base value of bool type
type Truth struct {
	value bool
}

func (t Truth) str() string {
	if t.value {
		return "t"
	}
	return "nil"
}

func (t Truth) isList() bool {
	return false
}

// Operator is a basic operator type
type Operator struct {
	value string
}

func (o Operator) str() string {
	return o.value
}

func (o Operator) isList() bool {
	return false
}

// ParseSingle Parses out a primitive
func ParseSingle(str string) Primitive {
	match, _ := regexp.MatchString("^[0-9]+$", str)
	if match {
		val, _ := strconv.Atoi(str)
		return Integer{value: val}
	}

	//Otherwise, assume operator
	return Operator{value: str}
}

//Parse parses a string to a primitive
func Parse(str string) Primitive {
	log.Printf("PARSING %v", str)
	listReg, _ := regexp.Compile("\\s+")

	// Lists should always start with '('
	if str[0] != '(' {
		return ParseSingle(str)
	}

	parseString := strings.Replace(strings.Replace(str, "(", "( ", -1), ")", " )", -1)
	elems := listReg.Split(parseString, -1)
	log.Printf("ELEMS = %v", elems)
	var stack []*listNode
	stackPointer := 0
	current := &listNode{}
	stack = append(stack, current)
	for _, val := range elems[1 : len(elems)-1] {
		log.Printf("WORKING on %v", val)
		if val == "(" {
			newln := &listNode{}
			if current.value != nil {
				newnd := &listNode{}
				current.next = newnd
				current = newnd
			}
			current.value = List{start: newln}
			stack = append(stack, newln)
			current = newln
			stackPointer++
			log.Printf("FOUND %p and %p", &stack[stackPointer], &newln)
		} else if val == ")" {
			stackPointer--
			current = stack[stackPointer]
		} else {
			if current.value != nil {
				log.Printf("Adding node to %v -> %v", current, stack[stackPointer])
				newnd := &listNode{}
				current.next = newnd
				current = newnd
			}
			current.value = ParseSingle(val)
			log.Printf("Set up current: %v (%v from %v)", current, stack[stackPointer], val)
		}
	}

	log.Printf("STACK = %v (%v)", stack, stackPointer)
	for i, v := range stack {
		log.Printf("%v. %v %p", i, v, v)
	}
	return List{start: stack[0]}
}
