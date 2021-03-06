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
	isStr() bool
	isNil() bool
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

func (l List) isStr() bool {
	return false
}

func (l List) isNil() bool {
	return false
}

func (l List) floatVal() float64 {
	return -1.0
}

// String Basic Type
type String struct {
	value string
}

func (s String) str() string {
	return s.value
}

func (s String) isList() bool {
	return false
}

func (s String) isSymbol() bool {
	return false
}

func (s String) isInt() bool {
	return false
}

func (s String) isNil() bool {
	return false
}

func (s String) isStr() bool {
	return true
}

func (s String) floatVal() float64 {
	return -1.0
}

func (l List) len() int {
	count := 0
	log.Printf("Running for %v with %v and %v", l, l.str(), l.start)
	st := l.start
	if st == nil {
		return 0
	} else if st.value != nil && !st.value.isNil() {
		count++
	}

	log.Printf("NOW %v", st.next)
	for st.next != nil {
		log.Printf("VALUE = %v", st.value.str())
		if !st.value.isNil() {
			count++
		}
		st = st.next
	}
	return count
}

func (l List) str() string {
	rep := "("
	if l.start != nil && l.start.single {
		rep = ""
	}
	node := l.start
	first := true
	for node != nil && node.value != nil && !node.value.isNil() {
		if first {
			log.Printf("Trying to print %v", node)
			if node.value != nil {
				rep += node.value.str()
			}
			first = false
		} else {
			log.Printf("Trying to print next %v", node)
			if node.single {
				rep += " . " + node.value.str()
			} else {
				rep += " " + node.value.str()
			}
		}
		node = node.next
	}
	if l.start != nil && !l.start.single {
		rep += ")"
	}
	return rep
}

type listNode struct {
	value  Primitive
	next   *listNode
	single bool
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

func (i Integer) isStr() bool {
	return false
}

func (i Integer) isNil() bool {
	return false
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
	if len(val) > 7 {
		log.Printf("SHRINKING %v", val)
		val = strconv.FormatFloat(f.value, 'f', 1, 64)
	}
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

func (f Float) isNil() bool {
	return false
}

func (f Float) isStr() bool {
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

func (r Ratio) isNil() bool {
	return false
}

func (r Ratio) isStr() bool {
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

func (t Truth) isStr() bool {
	return false
}

func (t Truth) isNil() bool {
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

func (s Symbol) isNil() bool {
	return false
}

func (s Symbol) isStr() bool {
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

func (n Nil) isNil() bool {
	return true
}

func (n Nil) isStr() bool {
	return false
}

func (n Nil) floatVal() float64 {
	return -1.0
}

// DeQuote converts quotes into explicit quote functions
func DeQuote(s string) string {
	str := strings.Replace(s, "#'", "'", -1)
	firstQuote := strings.Index(str, "'")

	for firstQuote >= 0 {
		log.Printf("DEQUOTING %v from index %v", str, firstQuote)
		nextIndex := firstQuote + 1
		if str[nextIndex] == '(' {
			bracketCount := 1
			for bracketCount > 0 {
				nextIndex++
				if str[nextIndex] == ')' {
					bracketCount--
				} else if str[nextIndex] == '(' {
					bracketCount++
				}
			}
			str = str[0:firstQuote] + "( quote " + str[firstQuote+1:nextIndex+1] + " )" + str[nextIndex+1:]
		} else {
			for nextIndex < len(str) && str[nextIndex] != ' ' {
				nextIndex++
			}
			str = str[0:firstQuote] + "( quote " + str[firstQuote+1:nextIndex] + " )" + str[nextIndex:]
		}

		firstQuote = strings.Index(str, "'")
	}

	return str
}

// ParseSingle Parses out a primitive
func ParseSingle(str string) Primitive {
	match, _ := regexp.MatchString("^-?[0-9]+$", str)
	if match {
		val, _ := strconv.Atoi(str)
		return Integer{value: val}
	}

	//Check for floats
	match, _ = regexp.MatchString("^[0-9]+\\.[0-9]+", str)
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
func Parse(strin string) Primitive {
	log.Printf("PARSING %v", strin)
	listReg, _ := regexp.Compile("\\s+")

	//Dequote first
	str := DeQuote(strin)

	// Lists should always start with '('
	if str[0] != '(' && !strings.HasPrefix(str, "'(") {
		return ParseSingle(str)
	}

	parseString := strings.Replace(strings.Replace(DeQuote(str), "(", "( ", -1), ")", " )", -1)
	log.Printf("PARSE STRING %v", parseString)
	elems := listReg.Split(parseString, -1)
	log.Printf("ELEMS = %v", elems)

	stack := listStack{}
	current := &listNode{value: &Nil{}}
	stack = stack.Push(current)
	isSingle := false
	for _, val := range elems[1 : len(elems)-1] {
		if val == "(" || val == "'(" {
			newln := &listNode{value: &Nil{}}
			if current.value != nil && !current.value.isNil() {
				newnd := &listNode{}
				current.next = newnd
				current = newnd
			}
			current.value = List{start: newln}
			stack = stack.Push(current)
			current = newln
		} else if val == ")" {
			stack, current = stack.Pop()
		} else if val == "." {
			isSingle = true
		} else {
			if current.value != nil && !current.value.isNil() {
				log.Printf("Adding new node for %v", current.value.str())
				newnd := &listNode{}
				current.next = newnd
				current = newnd
			}
			if isSingle {
				current.single = true
			}
			current.value = ParseSingle(val)
		}
	}
	log.Printf("POPPING FROM %v", stack)
	stack, sNode := stack.Pop()
	log.Printf("PARSE RESULT = %v", List{start: sNode}.str())
	retVal := List{start: sNode}

	return retVal
}
