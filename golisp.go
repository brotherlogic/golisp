package golisp

import "log"

// Eval evaluates lisp expressions
func Eval(p Primitive) Primitive {
	log.Printf("EVAL: %v", p.str())

	if !p.isList() {
		return p
	}

	l := p.(List)
	// All evluatable lists start with an operator
	operator, found := l.start.value.(Operator)
	if found {
		if operator.value == "+" {
			first := Eval(l.start.next.value).(Integer)
			second := Eval(l.start.next.next.value).(Integer)
			return Integer{value: first.value + second.value}
		} else if operator.value == "oddp" {
			first := Eval(l.start.next.value).(Integer)
			return Truth{value: first.value%2 == 1}
		} else if operator.value == "*" {
			first := Eval(l.start.next.value).(Integer)
			second := Eval(l.start.next.next.value).(Integer)
			return Integer{value: first.value * second.value}
		} else if operator.value == "/" {
			first := Eval(l.start.next.value).(Integer)
			second := Eval(l.start.next.next.value).(Integer)
			return Ratio{numerator: first.value, denominator: second.value}
		} else if operator.value == "equal" {
			first := Eval(l.start.next.value).(Integer)
			second := Eval(l.start.next.next.value).(Integer)
			return Truth{value: first.value == second.value}
		}
	}

	return l
}
