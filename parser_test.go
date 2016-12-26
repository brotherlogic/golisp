package golisp

import "testing"

func TestIntegerRep(t *testing.T) {
	v := &Integer{value: 5}
	if v.str() != "5" {
		t.Errorf("Integer str method fails: %v from %v", v.str(), v)
	}
}

func TestListRep(t *testing.T) {
	v := &List{start: &listNode{value: &Integer{value: 5}, next: &listNode{value: &Integer{value: 10}}}}
	if v.str() != "( 5 10 )" {
		t.Errorf("List str method fails: %v from %v", v.str(), v)
	}
}
