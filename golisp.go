package golisp

import (
	"errors"
	"fmt"
	"log"
)

// Interpreter is our lisp interpreter
type Interpreter struct {
	ops  []Operation
	vars []Variable
}

//Variable holds a local variable value
type Variable struct {
	name  string
	value Primitive
}

// Operation defines a given operation that can run in code
type Operation struct {
	name    string
	binding List
	expr    List
}

// Init prepares our interpreter
func Init() *Interpreter {
	i := &Interpreter{}

	i.vars = append(i.vars, Variable{name: "pi", value: Float{value: 3.14159}})
	i.vars = append(i.vars, Variable{name: "t", value: Truth{value: true}})
	i.vars = append(i.vars, Variable{name: "nil", value: Nil{}})

	return i
}

// Eval evaluates lisp expressions
func (i *Interpreter) Eval(p Primitive) (Primitive, error) {
	log.Printf("EVAL: %v with %v", p.str(), i)

	if !p.isList() {
		if p.isSymbol() {
			s := p.(Symbol)
			log.Printf("Searching for variable for %v", s)
			for _, v := range i.vars {
				if v.name == s.value {
					log.Printf("RETURNING %v", v.value)
					return v.value, nil
				}
			}
			return nil, fmt.Errorf("Cannot convert %v into symbol", p)
		}

		return p, nil
	}

	l := p.(List)
	// All evluatable lists start with an symbol
	symbol, found := l.start.value.(Symbol)
	if found {
		log.Printf("SYMBOL %v", symbol.value)
		if symbol.value == "+" {
			first, _ := i.Eval(l.start.next.value)
			second, _ := i.Eval(l.start.next.next.value)
			return Integer{value: first.(Integer).value + second.(Integer).value}, nil
		} else if symbol.value == "oddp" {
			first, _ := i.Eval(l.start.next.value)
			return Truth{value: first.(Integer).value%2 == 1}, nil
		} else if symbol.value == "*" {
			first, _ := i.Eval(l.start.next.value)
			second, _ := i.Eval(l.start.next.next.value)
			return Integer{value: first.(Integer).value * second.(Integer).value}, nil
		} else if symbol.value == "/" {
			first, _ := i.Eval(l.start.next.value)
			second, _ := i.Eval(l.start.next.next.value)
			if first.isInt() && second.isInt() {
				return Ratio{numerator: first.(Integer).value, denominator: second.(Integer).value}, nil
			}
			return Float{value: first.floatVal() / second.floatVal()}, nil
		} else if symbol.value == "equal" {
			first, errf := i.Eval(l.start.next.value)
			second, errs := i.Eval(l.start.next.next.value)
			if errf != nil || errs != nil {
				return nil, fmt.Errorf("Error in eval %v or %v", errf, errs)
			}
			if first.isInt() && second.isInt() {
				return Truth{value: first.(Integer).value == second.(Integer).value}, nil
			}
			return Truth{value: first.(String).value == second.(String).value}, nil
		} else if symbol.value == "defun" {
			op := Operation{name: l.start.next.value.str(),
				binding: l.start.next.next.value.(List),
				expr:    l.start.next.next.next.value.(List)}
			log.Printf("OP = %v", op)
			i.ops = append(i.ops, op)
			return Nil{}, nil
		}

		// If no operator is found, search through local ops
		log.Printf("Searching %v", i.ops)
		for _, op := range i.ops {
			if op.name == symbol.value {
				lvars := List{start: l.start.next}
				if lvars.len() != op.binding.len() {
					log.Printf("Unable to run %v on %v - mismatch in var length", op, l)
					return nil, errors.New("Badly formed function expression")
				}

				//Bind the variables
				vr := lvars.start
				val := op.binding.start
				for vr != nil {
					i.vars = append(i.vars, Variable{name: val.value.str(), value: vr.value})

					vr = vr.next
					val = val.next
				}

				res, _ := i.Eval(op.expr)
				return res, nil
			}
		}
	}

	return l, nil
}
