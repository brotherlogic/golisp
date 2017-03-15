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
	{[]string{"(defun double (n) (* n 2))", "(defun quadruple (n) (double (double n)))", "(quadruple 5)"}, []string{"nil", "nil", "20"}},
	{[]string{"(defun test () (* 85 97))", "(test)"}, []string{"nil", "8245"}},
	{[]string{"(quote foo)"}, []string{"foo"}},
	{[]string{"(quote (hello world))"}, []string{"(hello world)"}},
	{[]string{"'foo"}, []string{"foo"}},
	{[]string{"''foo"}, []string{"(quote foo)"}},
	{[]string{"(list 'quote 'foo)"}, []string{"(quote foo)"}},
	{[]string{"(first ''foo)"}, []string{"quote"}},
	{[]string{"(rest ''foo)"}, []string{"(foo)"}},
	{[]string{"(length ''foo)"}, []string{"2"}},
	{[]string{"'(+ 2 2)"}, []string{"(+ 2 2)"}},
	{[]string{"(eval '(+ 2 2))"}, []string{"4"}},
	{[]string{"'''boing"}, []string{"(quote (quote boing))"}},
	{[]string{"(eval '''boing)"}, []string{"(quote boing)"}},
	{[]string{"(eval (eval '''boing))"}, []string{"boing"}},
	{[]string{"'(list '* 9 6)"}, []string{"(list (quote *) 9 6)"}},
	{[]string{"(eval '(list '* 9 6))"}, []string{"(* 9 6)"}},
	{[]string{"(eval (eval '(list '* 9 6)))"}, []string{"54"}},
	{[]string{"(apply #'+ '(2 3))"}, []string{"5"}},
	{[]string{"(apply #'equal '(12 17))"}, []string{"nil"}},
	{[]string{"(apply #'cons '(as (you like it)))"}, []string{"(as you like it)"}},
	{[]string{"(if (oddp 1) 'odd 'even)"}, []string{"odd"}},
	{[]string{"(if (oddp 2) 'odd 'even)"}, []string{"even"}},
	{[]string{"(if t 'test-was-true 'test-was-false)"}, []string{"test-was-true"}},
	{[]string{"(if nil 'test-was-true 'test-was-false)"}, []string{"test-was-false"}},
	{[]string{"(if (symbolp 'foo) (* 5 5) (+ 5 5))"}, []string{"25"}},
	{[]string{"(if (symbolp 1) (* 5 5) (+ 5 5))"}, []string{"10"}},
	{[]string{"(defun my-abs (x) (if (< x 0) (- x) x))", "(my-abs -5)", "(my-abs 5)"}, []string{"nil", "5", "5"}},
	{[]string{"(defun symbol-test (x) (if (symbolp x) (list 'yes x 'is 'a 'symbol) (list 'no x 'is 'not 'a 'symbol)))", "(symbol-test 'rutabaga)", "(symbol-test 12345)"}, []string{"nil", "(yes rutabaga is a symbol)", "(no 12345 is not a symbol)"}},
	{[]string{"(if t 'happy)"}, []string{"happy"}},
	{[]string{"(if nil 'happy)"}, []string{"nil"}},
	{[]string{"(defun compare (x y) (cond ((equal x y) 'numbers-are-the-same) ((< x y) 'first-is-smaller) ((> x y) 'first-is-bigger)))", "(compare 3 5)", "(compare 7 2)", "(compare 4 4)"}, []string{"nil", "first-is-smaller", "first-is-bigger", "numbers-are-the-same"}},
	{[]string{"(defun compare (x y) (cond ((< x y) 'first-is-smaller) ((> x y) 'first-is-bigger)))", "(compare 5 5)"}, []string{"nil", "nil"}},
	{[]string{"(defun where-is (x) (cond ((equal x 'paris) 'france) ((equal x 'london) 'england) ((equal x 'beijing) 'china) (t 'unknown)))", "(where-is 'london)", "(where-is 'beijing)", "(where-is 'hackensack)"}, []string{"nil", "england", "china", "unknown"}},
	{[]string{"(defun emphasize (x) (cond ((equal (first x) 'good) (cons 'great (rest x))) ((equal (first x) 'bad) (cons 'awful (rest x)))))", "(emphasize '(good mystery story))", "(emphasize '(mediocre mystery story))"}, []string{"nil", "(great mystery story)", "nil"}},
	{[]string{"(defun emphasize2 (x) (cond ((equal (first x) 'good) (cons 'great (rest x))) ((equal (first x) 'bad) (cons 'awful (rest x))) (t x)))", "(emphasize2 '(good day))", "(emphasize2 '(bad day))", "(emphasize2 '(long day))"}, []string{"nil", "(great day)", "(awful day)", "(long day)"}},
	{[]string{"(defun compute (op x y) (cond ((equal op 'sum-of) (+ x y)) ((equal op 'product-of) (* x y)) (t '(that does not compute))))", "(compute 'sum-of 3 7)", "(compute 'product-of 2 4)", "(compute 'zorch-of 3 1)"}, []string{"nil", "10", "8", "(that does not compute)"}},
	{[]string{"(defun double (n) (* n 2))", "(double 5)"}, []string{"nil", "10"}},
	{[]string{"(setf vowels '(a e i o u))", "(length vowels)", "(rest vowels)", "vowels", "(setf vowels '(a e i o u and sometimes y))", "(rest (rest vowels))"}, []string{"(a e i o u)", "5", "(e i o u)", "(a e i o u)", "(a e i o u and sometimes y)", "(i o u and sometimes y)"}},
	{[]string{"(setf long-list '(a b c d e f g h i))", "(setf head (first long-list))", "(setf tail (rest long-list))", "(cons head tail)", "(equal long-list (cons head tail))", "(list head tail)"}, []string{"(a b c d e f g h i)", "a", "(b c d e f g h i)", "(a b c d e f g h i)", "t", "(a (b c d e f g h i))"}},
	{[]string{"(defun poor-style (p) (setf p (+ p 5)) (list 'result 'is p))", "(poor-style 8)"}, []string{"nil", "(result is 13)"}},
	{[]string{"(defun average (x y) (let ((sum (+ x y))) (list x y 'average 'is (/ sum 2.0))))", "(average 3 7)"}, []string{"nil", "(3 7 average is 5.0)"}},
	{[]string{"(defun price-change (old new) (let* ((diff (- new old)) (proportion (/ diff old)) (percentage (* proportion 100.0))) (list 'widgets 'changed 'by percentage 'percent)))", "(price-change 1.25 1.35)"}, []string{"nil", "(widgets changed by 8.0 percent)"}},
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
	{[]string{"(defun intro (x y) (list x this is y))", "(intro 'stanley 'livingstone)"}, []bool{false, true}, []string{"", "Error! this unassigned variable"}},
	{[]string{"(defun double (n) (* n 2))", "n"}, []bool{false, true}, []string{"", "Error! n unassigned variable"}},
	{[]string{"(defun test () (* 85 97))", "(test 1)"}, []bool{false, true}, []string{"", "Error! Too many arguments"}},
	{[]string{"(defun test () (* 85 97))", "test"}, []bool{false, true}, []string{"", "Error! test unassigned variable"}},
	{[]string{"(eval (eval (eval '''boing)))"}, []bool{true}, []string{"Error! boing unassigned variable"}},
	{[]string{"(defun double (n) (* n 2))", "(double 5)", "n"}, []bool{false, false, true}, []string{"nil", "10", "Error! n unassigned variable"}},
	{[]string{"(defun poor-style (p) (setf p (+ p 5)) (list 'result 'is p))", "(poor-style 8)", "p"}, []bool{false, false, false}, []string{"", "", "Error! p unassigned variable"}},
}

func TestGolispBad(t *testing.T) {
	for _, test := range baddata {
		i := Init()
		for j := range test.expression {
			log.Printf("TESTING %v", test.expression[j])
			e := Parse(test.expression[j])
			p, err := i.Eval(e.(Primitive), make([]Variable, 0))
			log.Printf("TESTING %v with %v", p, err)
			if test.fail[j] && err == nil {
				t.Fatalf("Executing %v has not failed and it should have done: %v -> %v", e.str(), p, p.str())
			} else if !test.fail[j] && err != nil {
				t.Fatalf("Executing %v has failed and it shouldn't have done: %v -> %v leads to %v", e.str(), p, p, err)
			}
			if err != nil && test.message[j] != "" {
				if test.message[j] != err.Error() {
					t.Fatalf("Error messages don't match %v vs %v", test.message[j], err.Error())
				}
			}
		}
	}
}

func TestGolispGood(t *testing.T) {
	for _, test := range testdata {
		i := Init()
		for j := range test.expression {
			log.Printf("Running test on %v", test.expression[j])
			e := Parse(test.expression[j])
			r, err := i.Eval(e.(Primitive), make([]Variable, 0))
			if err != nil {
				t.Fatalf("Executing %v has failed for %v", e.str(), err)
			} else if r.str() != test.result[j] {
				t.Fatalf("%v did not lead to %v, it lead to %v with %v", test.expression[j], test.result[j], r.str(), test.expression)
			}
		}
	}
}
