package golisp

import "log"

// Eval evaluates lisp expressions
func Eval(l List) Primitive {
	log.Printf("EVAL: %v", l.str())

	operator, found := l.start.value.(Operator)
	if found {
		first := l.start.next.value.(Integer)
		second := l.start.next.next.value.(Integer)

		if operator.value == "+" {
			return Integer{value: first.value + second.value}
		}
	}

	return l
}
