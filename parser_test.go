package golisp

import (
	"reflect"
	"testing"
)

func TestIntegerRep(t *testing.T) {
	v := &Integer{value: 5}
	if v.str() != "5" {
		t.Errorf("Integer str method fails: %v from %v", v.str(), v)
	}
}

func TestListRep(t *testing.T) {
	v := &List{start: &listNode{value: &Integer{value: 5}, next: &listNode{value: &Integer{value: 10}}}}
	if v.str() != "(5 10)" {
		t.Errorf("List str method fails: %v from %v", v.str(), v)
	}
}

func TestTruthRep(t *testing.T) {
	v := &Truth{value: false}
	if v.str() != "nil" {
		t.Errorf("Truth str method fails: %v from %v", v.str(), v)
	}
}

var islistdata = []struct {
	p Primitive
	r bool
}{
	{&Truth{}, false},
	{&Operator{}, false},
}

func TestListAssert(t *testing.T) {
	for _, test := range islistdata {
		if test.p.isList() != test.r {
			t.Errorf("Error in list assertion %v (%v)", test.p, test.p.isList())
		}
	}
}

var parsetestdata = []struct {
	expression string
	parsed     Primitive
}{
	{"5", Integer{value: 5}},
	{"(1 2 3)", List{start: &listNode{value: Integer{value: 1},
		next: &listNode{value: Integer{value: 2},
			next: &listNode{value: Integer{value: 3}}}}}},
	{"(+ 2 3)", List{start: &listNode{value: Operator{value: "+"},
		next: &listNode{value: Integer{value: 2},
			next: &listNode{value: Integer{value: 3}}}}}},
	{"(1 (2 3))", List{start: &listNode{value: Integer{value: 1},
		next: &listNode{value: List{start: &listNode{value: Integer{value: 2},
			next: &listNode{value: Integer{value: 3}}}}}}}},
	{"(oddp (+ 1 6))", List{start: &listNode{value: Operator{value: "oddp"},
		next: &listNode{value: List{start: &listNode{value: Operator{value: "+"},
			next: &listNode{value: Integer{value: 1},
				next: &listNode{value: Integer{value: 6}}}}}}}}},
}

func TestParse(t *testing.T) {
	for _, test := range parsetestdata {
		result := Parse(test.expression)
		if !reflect.DeepEqual(test.parsed, result) {
			t.Errorf("Parse fail %v->%v (should be %v)", test.expression, result.str(), test.parsed.str())
		}
	}
}
