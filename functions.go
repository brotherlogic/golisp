package golisp

import (
	"errors"
	"fmt"
	"log"
)

var funcs = map[string]func(*List) (Primitive, error){
	"cons":    cons,
	"+":       plus,
	"equal":   equal,
	"if":      iff,
	"symbolp": symbolp,
	"<":       lessthan,
	">":       greaterthan,
	"-":       subtract,
	"max":     max,
	"min":     min,
	"append":  appendf,
}

func appendf(params *List) (Primitive, error) {

	if params.start.value.isNil() {
		return params.start.next.value, nil
	}

	if params.start.next.value.isNil() {
		return List{start: params.start.value.(List).start}, nil
	}

	first := params.start.value.(List)
	second := params.start.next.value.(List)

	rList := List{start: first.start}
	curr := rList.start
	for curr.next != nil {
		curr = curr.next
	}
	curr.next = second.start
	return rList, nil
}

func max(params *List) (Primitive, error) {
	curr := params.start.next
	maxv := params.start.value.(Integer)
	for curr != nil {
		if curr.value.(Integer).value > maxv.value {
			maxv = curr.value.(Integer)
		}
		curr = curr.next
	}
	return maxv, nil
}

func min(params *List) (Primitive, error) {
	curr := params.start.next
	minv := params.start.value.(Integer)
	for curr != nil {
		if curr.value.(Integer).value < minv.value {
			minv = curr.value.(Integer)
		}
		curr = curr.next
	}
	return minv, nil
}

func subtract(params *List) (Primitive, error) {
	first := params.start.value
	if params.start.next != nil {
		second := params.start.next.value
		if first.isInt() && second.isInt() {
			return Integer{value: first.(Integer).value - second.(Integer).value}, nil
		} else if !second.isList() {
			return Float{value: first.floatVal() - second.floatVal()}, nil
		}

		return nil, errors.New("Error! Wrong type input to -")
	}
	return Integer{value: 0 - first.(Integer).value}, nil
}

func lessthan(params *List) (Primitive, error) {
	first := params.start.value
	second := params.start.next.value
	return Truth{value: first.floatVal() < second.floatVal()}, nil
}

func greaterthan(params *List) (Primitive, error) {
	first := params.start.value
	second := params.start.next.value
	return Truth{value: first.floatVal() > second.floatVal()}, nil
}

func symbolp(params *List) (Primitive, error) {
	head := params.start.value
	return Truth{value: head.isSymbol()}, nil
}

func cons(params *List) (Primitive, error) {
	head := params.start.value
	rest := params.start.next.value

	restList, ok := rest.(List)
	if ok {
		return List{start: &listNode{value: head, next: restList.start}}, nil
	}
	return List{start: &listNode{value: head, next: &listNode{value: rest, single: true}}}, nil
}

func plus(params *List) (Primitive, error) {
	first := params.start.value
	second := params.start.next.value
	if first.isInt() && second.isInt() {
		return Integer{value: first.(Integer).value + second.(Integer).value}, nil
	}
	return nil, fmt.Errorf("Error! Wrong type input to +")
}

func equal(params *List) (Primitive, error) {
	log.Printf("RUNNING EQUAL ON %v", params.str())
	first := params.start.value
	second := params.start.next.value
	if first.isInt() && second.isInt() {
		return Truth{value: first.(Integer).value == second.(Integer).value}, nil
	}
	return Truth{value: first.str() == second.str()}, nil
}

func iff(params *List) (Primitive, error) {
	if !params.start.value.isNil() && params.start.value.(Truth).value {
		return params.start.next.value, nil
	}
	if params.start.next.next != nil {
		return params.start.next.next.value, nil
	}
	return Truth{value: false}, nil
}
