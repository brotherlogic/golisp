package golisp

import (
	"fmt"
	"log"
)

var (
	funcs = map[string]func(*List) (Primitive, error){
		"cons":  cons,
		"+":     plus,
		"equal": equal,
	}
)

func cons(params *List) (Primitive, error) {
	head := params.start.value
	rest := params.start.next.value.(List)
	return List{start: &listNode{value: head, next: rest.start}}, nil
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
