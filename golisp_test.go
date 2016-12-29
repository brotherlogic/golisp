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
}

func TestGolisp(t *testing.T) {
	for _, test := range testdata {
		e := Parse(test.expression)
		r := Eval(e.(List))
		if r.str() != test.result {
			t.Errorf("%v did not lead to %v, it lead to %v", test.expression, test.result, r.str())
		}
	}
}
