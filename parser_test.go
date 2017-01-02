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

var testData = []struct {
	p  Primitive
	l  bool
	i  bool
	s  bool
	fv float64
}{
	{&Truth{}, false, false, false, -1.0},
	{&Symbol{}, false, false, true, -1.0},
	{&Ratio{numerator: 1, denominator: 2}, false, false, false, 0.5},
	{&Nil{}, false, false, false, -1.0},
	{&List{}, true, false, false, -1.0},
}

func TestListAssert(t *testing.T) {
	for _, test := range testData {
		if test.p.isList() != test.l {
			t.Errorf("Error in list assertion %v (%v)", test.p, test.p.isList())
		}
		if test.p.isInt() != test.i {
			t.Errorf("Error in int assertion %v (%v)", test.i, test.p.isInt())
		}
		if test.p.isSymbol() != test.s {
			t.Errorf("Error in symbol assertion %v (%v)", test.s, test.p.isSymbol())
		}
		if test.p.floatVal() != test.fv {
			t.Errorf("Error in float val assertion %v (%v)", test.fv, test.p.floatVal())
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
	{"(+ 2 3)", List{start: &listNode{value: Symbol{value: "+"},
		next: &listNode{value: Integer{value: 2},
			next: &listNode{value: Integer{value: 3}}}}}},
	{"(1 (2 3))", List{start: &listNode{value: Integer{value: 1},
		next: &listNode{value: List{start: &listNode{value: Integer{value: 2},
			next: &listNode{value: Integer{value: 3}}}}}}}},
	{"(oddp (+ 1 6))", List{start: &listNode{value: Symbol{value: "oddp"},
		next: &listNode{value: List{start: &listNode{value: Symbol{value: "+"},
			next: &listNode{value: Integer{value: 1},
				next: &listNode{value: Integer{value: 6}}}}}}}}},
	{"(/ (* 2 11) (+ 1 6))", List{start: &listNode{value: Symbol{value: "/"},
		next: &listNode{value: List{start: &listNode{value: Symbol{value: "*"},
			next: &listNode{value: Integer{value: 2},
				next: &listNode{value: Integer{value: 11}}}}},
			next: &listNode{value: List{start: &listNode{value: Symbol{value: "+"},
				next: &listNode{value: Integer{value: 1},
					next: &listNode{value: Integer{value: 6}}}}}}}}}},
	{"2.0", Float{value: 2.0}},
	{"(defun average (x y) (/ (+ x y) 2.0))", List{start: &listNode{value: Symbol{value: "defun"},
		next: &listNode{value: Symbol{value: "average"},
			next: &listNode{value: List{start: &listNode{value: Symbol{value: "x"},
				next: &listNode{value: Symbol{value: "y"}}}},
				next: &listNode{value: List{start: &listNode{value: Symbol{value: "/"},
					next: &listNode{value: List{start: &listNode{value: Symbol{value: "+"},
						next: &listNode{value: Symbol{value: "x"},
							next: &listNode{value: Symbol{value: "y"}}}}},
						next: &listNode{value: Float{value: 2.0}}}}}}}}}}},
}

func TestParse(t *testing.T) {
	for _, test := range parsetestdata {
		result := Parse(test.expression)
		if !reflect.DeepEqual(test.parsed, result) {
			t.Errorf("Parse fail %v->%p (should be %p) or otherwise %v !-> %v", test.expression, result, test.parsed, result.str(), test.parsed.str())
		}
	}
}
