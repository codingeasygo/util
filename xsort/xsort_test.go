package xsort

import "testing"

type testobj struct {
	Val string
}

func (t *testobj) Less(other interface{}) bool {
	return t.Val < other.(*testobj).Val
}

func TestSort(t *testing.T) {
	a, b := &testobj{Val: "b"}, &testobj{Val: "a"}
	vals := []*testobj{a, b}
	Sort(vals)
}
