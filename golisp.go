package golisp

import (
	"errors"
	"fmt"
	"log"
	"reflect"
)

// Interpreter is our lisp interpreter
type Interpreter struct {
	ops  []Operation
	vars []*Variable
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

func length(l List) Integer {
	curr := l.start
	len := 0
	for curr != nil {
		/*if curr.value.isList() {
			len += length(curr.value.(List)).value
		} else { */
		len++
		curr = curr.next
	}

	return Integer{value: len}
}

func (i *Interpreter) buildList(l *listNode, vs []Variable) (*List, error) {
	log.Printf("BUILDING LIST: %v", List{start: l}.str())
	head, err := i.Eval(l.value, vs)
	if err != nil {
		return nil, err
	}
	headNode := &listNode{value: head}
	currNode := headNode
	currProc := l.next
	for currProc != nil {
		currNode.next = &listNode{}
		currNode = currNode.next
		val, err := i.Eval(currProc.value, vs)
		if err != nil {
			return nil, err
		}
		currNode.value = val
		currProc = currProc.next
	}
	return &List{start: headNode}, nil
}

// Init prepares our interpreter
func Init() *Interpreter {
	i := &Interpreter{}

	i.vars = append(i.vars, &Variable{name: "pi", value: Float{value: 3.14159}})
	i.vars = append(i.vars, &Variable{name: "t", value: Truth{value: true}})
	i.vars = append(i.vars, &Variable{name: "nil", value: Nil{}})

	return i
}

// Eval evaluates lisp expressions
func (i *Interpreter) Eval(p Primitive, vs []Variable) (Primitive, error) {
	log.Printf("EVAL: %v with %v", p.str(), i)

	if !p.isList() {
		if p.isSymbol() {
			s := p.(Symbol)
			log.Printf("Searching for variable for %v with %v", s, vs)

			// Process stack variables first
			for _, v := range vs {
				if v.name == s.value {
					log.Printf("RETURNING %v", v.value)
					return v.value, nil
				}
			}

			for _, v := range i.vars {
				if v.name == s.value {
					log.Printf("RETURNING %v", v.value)
					return v.value, nil
				}
			}
			return nil, fmt.Errorf("Error! %v unassigned variable", p.str())
		}

		return p, nil
	}

	l := p.(List)

	// All evluatable lists start with an symbol
	symbol, found := l.start.value.(Symbol)
	log.Printf("SYMBOLD %v and %v -> %v", reflect.TypeOf(l.start.value), symbol, found)
	if found {
		log.Printf("SYMBOL %v", symbol.value)
		if symbol.value == "oddp" {
			first, err := i.Eval(l.start.next.value, vs)
			log.Printf("ERROR HERE is %v", err)
			if first.isInt() {
				return Truth{value: first.(Integer).value%2 == 1}, nil
			}
			return nil, fmt.Errorf("Error! Wrong type input to oddp")
		} else if symbol.value == "*" {
			first, _ := i.Eval(l.start.next.value, vs)
			second, _ := i.Eval(l.start.next.next.value, vs)
			return Integer{value: first.(Integer).value * second.(Integer).value}, nil
		} else if symbol.value == "/" {
			first, err1 := i.Eval(l.start.next.value, vs)
			second, err2 := i.Eval(l.start.next.next.value, vs)
			log.Printf("DIVIDE: %v and %v", err1, err2)
			if first.isInt() && second.isInt() {
				return Ratio{numerator: first.(Integer).value, denominator: second.(Integer).value}, nil
			}
			return Float{value: first.floatVal() / second.floatVal()}, nil
		} else if symbol.value == "first" {
			first, err := i.Eval(l.start.next.value, vs)
			if err != nil {
				return nil, err
			}
			if first.isList() {
				v := first.(List).start.value
				return v, nil
			}
		} else if symbol.value == "defun" {
			op := Operation{name: l.start.next.value.str(),
				binding: l.start.next.next.value.(List),
				expr:    l.start.next.next.next.value.(List)}
			log.Printf("OP = %v", op)

			//Verify the argument list
			curr := op.binding.start
			for curr != nil {
				if !curr.value.isSymbol() && !curr.value.isNil() {
					return nil, errors.New("Bad argument list")
				}
				curr = curr.next
			}

			i.ops = append(i.ops, op)
			return Nil{}, nil
		} else if symbol.value == "list" {
			log.Printf("LIST: %v", l.str())
			list := List{start: &listNode{}}
			currHead := list.start
			toadd := l.start.next
			for toadd != nil {
				val, err := i.Eval(toadd.value, vs)
				if err != nil {
					return nil, err
				}
				currHead.value = val
				if toadd.next != nil {
					currHead.next = &listNode{}
				}
				currHead = currHead.next
				toadd = toadd.next
			}
			return list, nil
		} else if symbol.value == "quote" {
			head := l.start.next.value
			log.Printf("QUOTE RETURN: %v", head.str())
			return head, nil
		} else if symbol.value == "rest" {
			log.Printf("REST: %v", l.start.next.value.str())
			procList, _ := i.Eval(l.start.next.value, vs)
			nlist := List{start: procList.(List).start.next}
			return nlist, nil
		} else if symbol.value == "length" {
			procList, _ := i.Eval(l.start.next.value, vs)
			return length(procList.(List)), nil
		} else if symbol.value == "eval" {
			temp, _ := i.Eval(l.start.next.value, vs)
			evalRes, err := i.Eval(temp, vs)
			return evalRes, err
		} else if symbol.value == "apply" {
			fname, _ := i.Eval(l.start.next.value, vs)
			li, _ := i.Eval(l.start.next.next.value, vs)
			ln := li.(List)

			log.Printf("RUNNING apply: %v and %v", fname, ln)
			if f, ok := funcs[fname.str()]; ok {
				log.Printf("Applying function %v", fname.str())
				return f(&ln)
			}
		} else if symbol.value == "cond" {
			evalList := List{start: l.start.next}
			curr := evalList.start
			for curr != nil {
				t, _ := i.Eval(curr.value.(List).start.value, vs)
				log.Printf("EVALd %v -> %v", curr.value.(List).start.value.str(), t)
				if t != nil && t.(Truth).value {
					r, _ := i.Eval(curr.value.(List).start.next.value, vs)
					return r, nil
				}
				curr = curr.next
			}
			return Truth{value: false}, nil
		} else if symbol.value == "setf" {
			name := l.start.next.value.str()
			value, _ := i.Eval(l.start.next.next.value, vs)

			vf := false
			for _, v := range i.vars {
				if v.name == name {
					v.value = value
					vf = true
				}
			}

			if !vf {
				i.vars = append(i.vars, &Variable{name: name, value: value})
			}
			return value, nil
		}

		// Search through the function table
		log.Printf("Searching the function table")
		if f, ok := funcs[symbol.value]; ok {
			log.Printf("Applying function %v", symbol.value)
			flist, err := i.buildList(l.start.next, vs)
			if err != nil {
				return nil, err
			}

			return f(flist)
		}

		// If no operator is found, search through local ops
		log.Printf("Searching %v", i.ops)
		for _, op := range i.ops {
			if op.name == symbol.value {
				lvars := List{start: l.start.next}
				log.Printf("ARG LEN %v (%v) and %v (%v)", lvars.len(), lvars.str(), op.binding.len(), op.binding.str())
				if lvars.len() != op.binding.len() {
					log.Printf("Unable to run %v on %v - mismatch in var length", op, l)
					return nil, errors.New("Error! Too many arguments")
				}

				//Bind the variables
				vr := lvars.start
				val := op.binding.start
				for vr != nil {
					eval, _ := i.Eval(vr.value, vs)
					log.Printf("BINDING %v to %v", eval.str(), val.value.str())

					//Override first, then bind a new variable
					found = false
					for i, v := range vs {
						if v.name == val.value.str() {
							log.Printf("BEFORE %v %v", vs, v)
							vs[i].value = eval
							log.Printf("AFTER %v %v", vs, v)
							found = true
						}
					}
					if !found {
						vs = append(vs, Variable{name: val.value.str(), value: eval})
					}

					vr = vr.next
					val = val.next
				}

				log.Printf("BOUND %v to eval %v", vs, op.expr.str())
				res, err := i.Eval(op.expr, vs)
				return res, err
			}
		}
	}

	listv, found := l.start.value.(List)
	if found {
		l.start.value, _ = i.Eval(listv, vs)
		return i.Eval(l, vs)
	}

	log.Printf("Bottoming out")
	return nil, errors.New("Error! '" + l.start.value.str() + "' undefined function")
}
