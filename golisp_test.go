package golisp

import "testing"

var testdata = []struct {
	expression []string
	result     []string
}{
	{[]string{"(1 2 3)"}, []string{"(1 2 3)"}},
	{[]string{"(+ 2 3)"}, []string{"5"}},
	{[]string{"(+ 1 6)"}, []string{"7"}},
	{[]string{"(oddp (+ 1 6))"}, []string{"t"}},
	{[]string{"(* 3 (+ 1 6))"}, []string{"21"}},
	{[]string{"(/ (* 2 11) (+ 1 6))"}, []string{"22/7"}},
	{[]string{"23"}, []string{"23"}},
	{[]string{"t"}, []string{"t"}},
	{[]string{"nil"}, []string{"nil"}},
	{[]string{"(equal (+ 7 5) (* 2 8))"}, []string{"nil"}},
	{[]string{"(/ (+ 6 8) 2.0)"}, []string{"7.0"}},
	{[]string{"(defun average (x y) (/ (+ x y) 2.0))", "(average 6 8)"}, []string{"nil", "7.0"}},
	{[]string{"(defun square (n) (* n n))", "(square 2)"}, []string{"nil", "4"}},
	{[]string{"(defun total-cost (quantity price handling-charge) (+ (* quantity price) handling-charge))", "(total-cost 2 3 4)"}, []string{"nil", "10"}},
	{[]string{"pi"}, []string{"3.14159"}},
	{[]string{"(equal 'kirk 'spock)"}, []string{"nil"}},
}

var baddata = []struct {
	expression []string
	fail       []bool
}{
	{[]string{"(defun average (x y) (/ (+ x y) 2.0))", "(average 6 8 7)"}, []bool{false, true}},
	{[]string{"(equal kirk spock)"}, []bool{true}},
}

func TestGolispBad(t *testing.T) {
	for _, test := range baddata {
		i := Init()
		for j := range test.expression {
			e := Parse(test.expression[j])
			p, err := i.Eval(e.(Primitive))
			if test.fail[j] && err == nil {
				t.Errorf("Executing %v has not failed and it should have done: %v", e.str(), p)
			}
		}
	}
}

func TestGolisp(t *testing.T) {
	for _, test := range testdata {
		i := Init()
		for j := range test.expression {
			e := Parse(test.expression[j])
			r, err := i.Eval(e.(Primitive))
			if err != nil {
				t.Errorf("Executing %v has failed for %v", e.str(), err)
			} else if r.str() != test.result[j] {
				t.Errorf("%v did not lead to %v, it lead to %v", test.expression[j], test.result[j], r.str())
			}
		}
	}
}
