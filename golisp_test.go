package golisp

import (
	"errors"
	"log"
	"strings"
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
	{[]string{"(setf x 57)", "(defun newvar (x) (list 'value 'of 'x 'is x))", "x", "(newvar 'whoopie)", "x"}, []string{"57", "nil", "57", "(value of x is whoopie)", "57"}},
	{[]string{"(setf a 100)", "(defun f (a) (list a (g ( + a 1))))", "(defun g (b) (list a b))", "(f 3)"}, []string{"100", "nil", "nil", "(3 (100 4))"}},
	{[]string{"(cons 'w '(x y z))"}, []string{"(w x y z)"}},
	{[]string{"(cons '(a b c) 'd)"}, []string{"((a b c) . d)"}},
	{[]string{"(append '(friends romans) '(and countrymen))"}, []string{"(friends romans and countrymen)"}},
	{[]string{"(append '(l m n o) '(p q r))"}, []string{"(l m n o p q r)"}},
	{[]string{"(append '(april showers) nil)"}, []string{"(april showers)"}},
	{[]string{"(append nil '(bring may flowers))"}, []string{"(bring may flowers)"}},
	{[]string{"(append nil nil)"}, []string{"nil"}},
	{[]string{"(append '((a 1) (b 2)) '((c 3) (d 4)))"}, []string{"((a 1) (b 2) (c 3) (d 4))"}},
	{[]string{"(append '(w x y) 'z)"}, []string{"(w x y . z)"}},
	{[]string{"(append '(a b c) '(d))"}, []string{"(a b c d)"}},
	{[]string{"(defun add-to-end (x e) (append x (list e)))", "(add-to-end '(a b c) 'd)"}, []string{"nil", "(a b c d)"}},
	{[]string{"(cons 'rice '(and beans))"}, []string{"(rice and beans)"}},
	{[]string{"(list 'rice '(and beans))"}, []string{"(rice (and beans))"}},
	{[]string{"(cons '(here today) '(gone tomorrow))"}, []string{"((here today) gone tomorrow)"}},
	{[]string{"(list '(here today) '(gone tomorrow))"}, []string{"((here today) (gone tomorrow))"}},
	{[]string{"(append '(here today) '(gone tomorrow))"}, []string{"(here today gone tomorrow)"}},
	{[]string{"(cons '(eat at) 'joes)"}, []string{"((eat at) . joes)"}},
	{[]string{"(list '(eat at) 'joes)"}, []string{"((eat at) joes)"}},
	{[]string{"(append '(eat at) 'joes)"}, []string{"(eat at . joes)"}},
	{[]string{"(reverse '(one two three four five))"}, []string{"(five four three two one)"}},
	{[]string{"(reverse '(l i v e))"}, []string{"(e v i l)"}},
	{[]string{"(reverse '((my oversight) (your blunder) (his negligence)))"}, []string{"((his negligence) (your blunder) (my oversight))"}},
	{[]string{"(setf vow '(to have and to hold))", "(reverse vow)", "vow"}, []string{"(to have and to hold)", "(hold to and have to)", "(to have and to hold)"}},
	{[]string{"(defun add-to-end (x y) (reverse (cons y (reverse x))))", "(add-to-end '(a b c) 'd)"}, []string{"nil", "(a b c d)"}},
	{[]string{"(nthcdr 0 '(a b c))"}, []string{"(a b c)"}},
	{[]string{"(nthcdr 1 '(a b c))"}, []string{"(b c)"}},
	{[]string{"(nthcdr 2 '(a b c))"}, []string{"(c)"}},
	{[]string{"(nthcdr 3 '(a b c))"}, []string{"nil"}},
	{[]string{"(nthcdr 4 '(a b c))"}, []string{"nil"}},
	{[]string{"(nthcdr 5 '(a b c))"}, []string{"nil"}},
	{[]string{"(nthcdr 2 '(a b c . d))"}, []string{"(c . d)"}},
	{[]string{"(nthcdr 3 '(a b c . d))"}, []string{"d"}},
	{[]string{"(nth 0 '(a b c))"}, []string{"a"}},
	{[]string{"(nth 1 '(a b c))"}, []string{"b"}},
	{[]string{"(nth 2 '(a b c))"}, []string{"c"}},
	{[]string{"(nth 3 '(a b c))"}, []string{"nil"}},
	{[]string{"(last '(all is forgiven))"}, []string{"(forgiven)"}},
	{[]string{"(last nil)"}, []string{"nil"}},
	{[]string{"(last '(a b c . d))"}, []string{"(c . d)"}},
	{[]string{"(setf ducks '(huey dewey louie))", "(member 'huey ducks)", "(member 'dewey ducks)", "(member 'louie ducks)", "(member 'mickey ducks)"}, []string{"(huey dewey louie)", "(huey dewey louie)", "(dewey louie)", "(louie)", "nil"}},
	{[]string{"(defun beforep (x y l) (member y (member x l)))", "(beforep 'not 'whom '(ask not for whom the bell tolls))", "(beforep 'thee 'tolls '(it tolls for thee))"}, []string{"nil", "(whom the bell tolls)", "nil"}},
	{[]string{"(intersection '(fred john mary) '(sue mary fred))"}, []string{"(fred mary)"}},
	{[]string{"(intersection '(a s d f g) '(v w s r a))"}, []string{"(a s)"}},
	{[]string{"(intersection '(foo bar baz) '(xam gorp bletch))"}, []string{"nil"}},
	{[]string{"(union '(finger hand arm) '(toe finger foot leg))"}, []string{"(finger hand arm toe foot leg)"}},
	{[]string{"(union '(fred john mary) '(sue mary fred))"}, []string{"(fred john mary sue)"}},
	{[]string{"(union '(a s d f g) '(v w s r a))"}, []string{"(a s d f g v w r)"}},
	{[]string{"(set-difference '(alpha bravo charlie delta) '(bravo charlie))"}, []string{"(alpha delta)"}},
	{[]string{"(set-difference '(alpha bravo charlie delta) '(echo alpha foxtrot))"}, []string{"(bravo charlie delta)"}},
	{[]string{"(set-difference '(alpha bravo) '(bravo alpha))"}, []string{"nil"}},
	{[]string{"(setf line1 '(all things in moderation))", "(setf line2 '(moderation in the defence of liberty is not virtue))", "(set-difference line1 line2)", "(set-difference line2 line1)"}, []string{"(all things in moderation)", "(moderation in the defence of liberty is not virtue)", "(all things)", "(the defence of liberty is not virtue)"}},
	{[]string{"(subsetp '(a i) '(a e i o u))"}, []string{"t"}},
	{[]string{"(subsetp '(a x) '(a e i o u))"}, []string{"nil"}},
	{[]string{"(defun titledp (name) (member (first name) '(mr ms miss mrs)))", "(titledp '(jane doe))", "(titledp '(ms jane doe))"}, []string{"nil", "nil", "(ms miss mrs)"}},
	{[]string{"(setf male-first-names '(john kim richard fred george))", "(setf female-first-names '(jane mary wanda barbara kim))", "(defun malep (name) (and (member name male-first-names) (not (member name female-first-names))))", "(defun femalep (name) (and (member name female-first-names) (not (member name male-first-names))))", "(malep 'richard)", "(malep 'barbara)", "(femalep 'barbara)", "(malep 'kim)"}, []string{"(john kim richard fred george)", "(jane mary wanda barbara kim)", "nil", "nil", "t", "nil", "t", "nil"}},
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
	{[]string{"(defun intro ('x 'y) (list x 'this 'is y))"}, []bool{true}, []string{"Error! Bad argument list"}},
	{[]string{"(defun intro ((x) (y)) (list x 'this 'is y))"}, []bool{true}, []string{"Error! Bad argument list"}},
	{[]string{"(defun intro (x y) (list (x) 'this 'is (y)))", "(intro 'stanley 'livingstone)"}, []bool{false, true}, []string{"", "Error in function intro: 'x' undefined function"}},
	{[]string{"(defun intro (x y) (list x this is y))", "(intro 'stanley 'livingstone)"}, []bool{false, true}, []string{"", "Error in function intro: this unassigned variable"}},
	{[]string{"(defun double (n) (* n 2))", "n"}, []bool{false, true}, []string{"", "Error! n unassigned variable"}},
	{[]string{"(defun test () (* 85 97))", "(test 1)"}, []bool{false, true}, []string{"", "Error! Too many arguments"}},
	{[]string{"(defun test () (* 85 97))", "test"}, []bool{false, true}, []string{"", "Error! test unassigned variable"}},
	{[]string{"(eval (eval (eval '''boing)))"}, []bool{true}, []string{"Error! boing unassigned variable"}},
	{[]string{"(defun double (n) (* n 2))", "(double 5)", "n"}, []bool{false, false, true}, []string{"nil", "10", "Error! n unassigned variable"}},
	{[]string{"(defun poor-style (p) (setf p (+ p 5)) (list 'result 'is p))", "(poor-style 8)", "p"}, []bool{false, false, false}, []string{"", "", "Error! p unassigned variable"}},
	{[]string{"(defun faulty-size-range (x y z) (let ((biggest (max x y z)) (smallest (min x y z)) (r (/ biggest smallest 1.0))) (list 'factor 'of r)))", "(faulty-size-range 35 87 4)"}, []bool{false, true}, []string{"", "Error in function faulty-size-range: biggest unassigned variable"}},
	{[]string{"(append 'a '(b c d))"}, []bool{true}, []string{"Error! a is not a list"}},
	{[]string{"(append 'rice '(and beans))"}, []bool{true}, []string{"Error! rice is not a list"}},
	{[]string{"(reverse 'live)"}, []bool{true}, []string{"Error! Wrong type input"}},
	{[]string{"(nthcdr 4 '(a b c . d))"}, []bool{true}, []string{"Error! d is not a list"}},
	{[]string{"(last 'nevermore)"}, []bool{true}, []string{"Error! nevermore is not a list"}},
}

func TestGolispBad(t *testing.T) {
	for _, test := range baddata {
		i := Init()
		for j := range test.expression {
			log.Printf("TESTING %v", test.expression[j])
			e := Parse(test.expression[j])
			p, err := i.Eval(e.(Primitive), make([]Variable, 0))
			if err != nil && !strings.HasPrefix(err.Error(), "Error") {
				err = errors.New("Error! " + err.Error())
			}
			log.Printf("TESTING %v with %v", p, err)
			if test.fail[j] && err == nil {
				t.Fatalf("Executing %v has not failed and it should have done: %v -> %v", e.str(), p, p.str())
			} else if !test.fail[j] && err != nil {
				t.Fatalf("Executing %v has failed and it shouldn't have done: %v -> %v leads to %v", e.str(), p, p, err)
			}
			if err != nil && test.message[j] != "" {
				if test.message[j] != err.Error() {
					t.Fatalf("Error messages don't match %v vs %v but %v", test.message[j], err.Error(), err)
				}
			}
		}
	}
}

func TestGolispGood(t *testing.T) {
	for _, test := range testdata {
		i := Init()
		for j := range test.expression {
			log.Printf("TESTING running test on %v", test.expression[j])
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
