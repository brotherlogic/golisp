package golisp

import "testing"

var testdata = []struct {
	expression string
	result     string
}{
	{"(1 2 3)", "(1 2 3)"},
	{"(+ 2 3)", "5"},
	{"(+ 1 6)", "7"},
	{"(oddp (+ 1 6))", "t"},
	{"(* 3 (+ 1 6))", "21"},
	{"(/ (* 2 11) (+ 1 6))", "22/7"},
	{"23", "23"},
	{"t", "t"},
	{"nil", "nil"},
}

func TestGolisp(t *testing.T) {
	for _, test := range testdata {
		e := Parse(test.expression)
		r := Eval(e.(Primitive))
		if r.str() != test.result {
			t.Errorf("%v did not lead to %v, it lead to %v", test.expression, test.result, r.str())
		}
	}
}
