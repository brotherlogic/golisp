package golisp

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

const (
	nilstr string = "nil"
)

// Primitive is the base of all elements in source code
type Primitive interface {
	str() string
	isList() bool
	isSymbol() bool
	isInt() bool
	floatVal() float64
}

// List Basic Type
type List struct {
	start *listNode
}

func (l List) isList() bool {
	return true
}

func (l List) isSymbol() bool {
	return false
}

func (l List) isInt() bool {
	return false
}

func (l List) floatVal() float64 {
	return -1.0
}

func (l List) len() int {
	count := 0
	st := l.start
	for st.next != nil {
		count++
		st = st.next
	}
	return count
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

func (i Integer) isSymbol() bool {
	return false
}

func (i Integer) isInt() bool {
	return true
}

func (i Integer) floatVal() float64 {
	return float64(i.value)
}

// Float is a base value of float type
type Float struct {
	value float64
}

func (f Float) str() string {
	val := strconv.FormatFloat(f.value, 'f', -1, 64)
	if strings.Contains(val, ".") {
		return val
	}
	return val + ".0"
}

func (f Float) isList() bool {
	return false
}

func (f Float) isSymbol() bool {
	return false
}

func (f Float) isInt() bool {
	return false
}

func (f Float) floatVal() float64 {
	return f.value
}

// Ratio is a base value of ratio type
type Ratio struct {
	numerator   int
	denominator int
}

func (r Ratio) str() string {
	return fmt.Sprintf("%v/%v", r.numerator, r.denominator)
}

func (r Ratio) isList() bool {
	return false
}

func (r Ratio) isSymbol() bool {
	return false
}

func (r Ratio) isInt() bool {
	return false
}

func (r Ratio) floatVal() float64 {
	return float64(r.numerator) / float64(r.denominator)
}

// Truth is a base value of bool type
type Truth struct {
	value bool
}

func (t Truth) str() string {
	if t.value {
		return "t"
	}
	return nilstr
}

func (t Truth) isList() bool {
	return false
}

func (t Truth) isSymbol() bool {
	return false
}

func (t Truth) isInt() bool {
	return false
}

func (t Truth) floatVal() float64 {
	return -1.0
}

// Symbol is a basic operator type
type Symbol struct {
	value string
}

func (s Symbol) str() string {
	return s.value
}

func (s Symbol) isList() bool {
	return false
}

func (s Symbol) isSymbol() bool {
	return true
}

func (s Symbol) isInt() bool {
	return false
}

func (s Symbol) floatVal() float64 {
	return -1.0
}

// Nil is a basic nil type
type Nil struct{}

func (n Nil) str() string {
	return nilstr
}

func (n Nil) isList() bool {
	return false
}

func (n Nil) isSymbol() bool {
	return false
}

func (n Nil) isInt() bool {
	return false
}

func (n Nil) floatVal() float64 {
	return -1.0
}

// ParseSingle Parses out a primitive
func ParseSingle(str string) Primitive {
	match, _ := regexp.MatchString("^[0-9]+$", str)
	if match {
		val, _ := strconv.Atoi(str)
		return Integer{value: val}
	}

	//Check for floats
	match, _ = regexp.MatchString("^[0-9]+\\.[0.9]+", str)
	if match {
		val, _ := strconv.ParseFloat(str, 64)
		return Float{value: val}
	}
	//Otherwise, assume operator
	return Symbol{value: str}
}

type listStack []*listNode

func (s listStack) Push(v *listNode) listStack {
	return append(s, v)
}

func (s listStack) Pop() (listStack, *listNode) {
	l := len(s)
	return s[:l-1], s[l-1]
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
	stack := listStack{}
	current := &listNode{}
	stack = stack.Push(current)
	for _, val := range elems[1 : len(elems)-1] {
		if val == "(" {
			newln := &listNode{}
			if current.value != nil {
				newnd := &listNode{}
				current.next = newnd
				current = newnd
			}
			current.value = List{start: newln}
			stack = stack.Push(current)
			current = newln
		} else if val == ")" {
			stack, current = stack.Pop()
		} else {
			if current.value != nil {
				newnd := &listNode{}
				current.next = newnd
				current = newnd
			}
			current.value = ParseSingle(val)
		}
	}
	stack, sNode := stack.Pop()
	log.Printf("PARSE RESULT = %v", List{start: sNode}.str())
	return List{start: sNode}
}
