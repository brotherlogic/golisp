package golisp

import (
	"errors"
	"fmt"
	"log"
)

var funcs = map[string]func(*List) (Primitive, error){
	"cons":         cons,
	"+":            plus,
	"equal":        equal,
	"if":           iff,
	"symbolp":      symbolp,
	"<":            lessthan,
	">":            greaterthan,
	"-":            subtract,
	"max":          max,
	"min":          min,
	"append":       appendf,
	"reverse":      reversef,
	"cdr":          cdr,
	"nthcdr":       nthcdr,
	"nth":          nth,
	"car":          car,
	"last":         last,
	"member":       member,
	"intersection": intersection,
}

func intersectionp(l1, l2 List, build *listNode) {

	if l1.start != nil && !l1.start.value.isNil() {
		head := l1.start.value
		if !memberp(head, l2).isNil() {
			build.value = head
			build.next = &listNode{}
			intersectionp(List{start: l1.start.next}, l2, build.next)
		} else {
			intersectionp(List{start: l1.start.next}, l2, build)
		}
	}
}

func intersection(l *List) (Primitive, error) {
	retv := &List{start: &listNode{}}
	intersectionp(List{start: l.start.value.(List).start}, List{start: l.start.next.value.(List).start}, retv.start)

	if retv.start.value == nil {
		return &Nil{}, nil
	}
	return retv, nil
}

func memberp(p Primitive, l List) Primitive {
	if l.start != nil && !l.start.value.isNil() {
		if p.str() == l.start.value.str() {
			return l
		}
		return memberp(p, List{start: l.start.next})
	}

	return Nil{}
}

func member(l *List) (Primitive, error) {
	nl := l.start.next.value.(List)
	return memberp(l.start.value, nl), nil
}

func last(l *List) (Primitive, error) {
	if l.start.value == nil || l.start.value.isNil() {
		return l.start.value, nil
	}

	if l.start.value.isList() {
		listy := l.start.value.(List)
		return doLast(&listy)
	}

	return nil, errors.New(l.start.value.str() + " is not a list")
}

func doLast(l *List) (Primitive, error) {
	if l.start.next == nil || l.start.next.value.isNil() || l.start.next.single {
		return l, nil
	}

	return doLast(&List{start: l.start.next})
}

func car(l *List) (Primitive, error) {
	return l.start.value, nil
}

func nth(l *List) (Primitive, error) {
	nc, _ := nthcdr(l)
	if nc == nil || nc.isNil() {
		return nc, nil
	}
	ncl := nc.(List)
	return car(&ncl)
}

func cdr(l *List) (Primitive, error) {
	if l.start.single {
		return nil, errors.New(l.str() + " is not a list")
	}
	return List{start: l.start.next}, nil
}

func nthcdr(l *List) (Primitive, error) {
	log.Printf("NTHCDR: %v, %v", l.str(), l.start.next.value)

	if l.start.next.value == nil || l.start.next.value.isNil() || l.start.next.value.(List).start == nil {
		return Nil{}, nil
	}

	log.Printf("VALUE = %v", l.start.value.str())
	val := l.start.value.(Integer).value
	if val > 0 {
		ll := l.start.next.value.(List)
		cdrv, err := cdr(&ll)
		if err != nil {
			return nil, err
		}
		log.Printf("BUILT %v", cdrv.str())
		return nthcdr(&List{start: &listNode{value: Integer{value: val - 1}, next: &listNode{value: cdrv}}})
	}

	return l.start.next.value, nil
}

func copy(l List) List {
	nList := List{}
	lNode := &listNode{}
	nList.start = lNode
	cNode := l.start
	for cNode != nil {
		lNode.value = cNode.value
		if cNode.next != nil {
			lNode.next = &listNode{}
			lNode = lNode.next
		}
		cNode = cNode.next
	}

	return nList
}

func reverseh(l *listNode) (*listNode, *listNode) {
	if l.next == nil {
		return l, l
	}

	start, end := reverseh(l.next)
	end.next = l
	l.next = nil
	return start, l
}

func reversef(params *List) (Primitive, error) {
	if !params.start.value.isList() {
		return nil, errors.New("Wrong type input")
	}

	mainList := copy(params.start.value.(List))
	r, _ := reverseh(mainList.start)
	return List{start: r}, nil
}

func appendf(params *List) (Primitive, error) {

	if params.start.value.isNil() {
		return params.start.next.value, nil
	}

	if params.start.next.value.isNil() {
		return List{start: params.start.value.(List).start}, nil
	}

	if !params.start.value.isList() {
		return nil, errors.New(params.start.value.str() + " is not a list")
	}

	first := params.start.value.(List)
	rList := List{start: first.start}
	curr := rList.start
	for curr.next != nil {
		curr = curr.next
	}

	if params.start.next.value.isList() {
		second := params.start.next.value.(List)
		curr.next = second.start
	} else {
		curr.next = &listNode{value: params.start.next.value, single: true}
	}

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
