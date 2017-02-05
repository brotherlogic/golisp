package golisp

import (
	"log"
	"testing"
)

var testdata = []struct {
	expression []string
	result     []string
}{
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
	{[]string{"(list 'james t 'kirk)"}, []string{"(james t kirk)"}},
	{[]string{"(defun riddle (x y) (list 'why 'is 'a x 'like 'a y))", "(riddle 'raven 'writing-desk)"}, []string{"nil", "(why is a raven like a writing-desk)"}},
	{[]string{"(first (list 1 2 3))"}, []string{"1"}},
	{[]string{"(first '(we hold these truths))"}, []string{"we"}},
	{[]string{"'(+ 1 2)"}, []string{"(+ 1 2)"}},
	{[]string{"(oddp (+ 1 2))"}, []string{"t"}},
	{[]string{"(list 'a 'b 'c)"}, []string{"(a b c)"}},
	{[]string{"(cons 'a '(b c))"}, []string{"(a b c)"}},
	{[]string{"(+ 10 (- 5 2))"}, []string{"13"}},
	{[]string{"(list 'buy '(* 27 34) 'bagels)"}, []string{"(buy (* 27 34) bagels)"}},
	{[]string{"(list 'buy (* 27 34) 'bagels)"}, []string{"(buy 918 bagels)"}},
	{[]string{"'(foo bar baz)"}, []string{"(foo bar baz)"}},
	{[]string{"(list 'foo 'bar 'baz)"}, []string{"(foo bar baz)"}},
	{[]string{"(cons 'foo '(bar baz))"}, []string{"(foo bar baz)"}},
	{[]string{"(list 33 'squared 'is (* 33 33))"}, []string{"(33 squared is 1089)"}},
	{[]string{"'(33 squared is (* 33 33))"}, []string{"(33 squared is (* 33 33))"}},
	{[]string{"(defun intro (x y) (list x 'this 'is y))", "(intro 'stanley 'livingstone)"}, []string{"nil", "(stanley this is livingstone)"}},
	{[]string{"(defun intro (x y) (list 'x 'this 'is 'y))", "(intro 'stanley 'livingstone)"}, []string{"nil", "(x this is y)"}},
}

var baddata = []struct {
	expression []string
	fail       []bool
	message    []string
}{
	{[]string{"(1 2 3)"}, []bool{true}, []string{""}},
	{[]string{"(defun average (x y) (/ (+ x y) 2.0))", "(average 6 8 7)"}, []bool{false, true}, []string{"", ""}},
	{[]string{"(equal kirk spock)"}, []bool{true}, []string{""}},
	{[]string{"(list kirk 1 2)"}, []bool{true}, []string{""}},
	{[]string{"(first (we hold these truths))"}, []bool{true}, []string{"Error! 'we' undefined function"}},
	{[]string{"(first 1 2 3 4)"}, []bool{true}, []string{""}},
	{[]string{"(oddp '(+ 1 2))"}, []bool{true}, []string{"Error! Wrong type input to oddp"}},
	{[]string{"(cons 'a (b c))"}, []bool{true}, []string{"Error! 'b' undefined function"}},
	{[]string{"(+ 10 '(- 5 2))"}, []bool{true}, []string{"Error! Wrong type input to +"}},
	{[]string{"(- 10 '(- 5 2))"}, []bool{true}, []string{"Error! Wrong type input to -"}},
	{[]string{"('foo 'bar 'baz)"}, []bool{true}, []string{"Error! 'foo' undefined function"}},
	{[]string{"(list foo bar baz)"}, []bool{true}, []string{"Error! foo unassigned variable"}},
	{[]string{"(foo bar baz)"}, []bool{true}, []string{"Error! 'foo' undefined function"}},
	{[]string{"(defun intro ('x 'y) (list x 'this 'is y))"}, []bool{true}, []string{"Bad argument list"}},
	{[]string{"(defun intro ((x) (y)) (list x 'this 'is y))"}, []bool{true}, []string{"Bad argument list"}},
	{[]string{"(defun intro (x y) (list (x) 'this 'is (y)))", "(intro 'stanley 'livingstone)"}, []bool{false, true}, []string{"", "Error! 'x' undefined function"}},
}

func TestGolispBad(t *testing.T) {
	for _, test := range baddata {
		i := Init()
		for j := range test.expression {
			log.Printf("TESTING %v", test.expression[j])
			e := Parse(test.expression[j])
			p, err := i.Eval(e.(Primitive))
			log.Printf("TESTING %v with %v", p, err)
			if test.fail[j] && err == nil {
				t.Errorf("Executing %v has not failed and it should have done: %v -> %v", e.str(), p, p.str())
			} else if !test.fail[j] && err != nil {
				t.Errorf("Executing %v has failed and it shouldn't have done: %v -> %v leads to %v", e.str(), p, p, err)
			}
			if err != nil && test.message[j] != "" {
				if test.message[j] != err.Error() {
					t.Errorf("Error messages don't match %v vs %v", test.message[j], err.Error())
				}
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
