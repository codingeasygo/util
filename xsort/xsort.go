package xsort

import (
	"reflect"
	"sort"
)

type Lesser interface {
	Less(other interface{}) bool
}

type Interface struct {
	swapper func(x, y int)
	lesser  func(x, y int) bool
	value   reflect.Value
}

func NewInterface(value interface{}) (v *Interface) {
	v = &Interface{
		value:   reflect.ValueOf(value),
		swapper: reflect.Swapper(value),
	}
	return
}

func (i *Interface) Len() int {
	return i.value.Len()
}

func (i *Interface) Less(x, y int) bool {
	if i.lesser == nil {
		a := i.value.Index(x).Interface().(Lesser)
		b := i.value.Index(y).Interface()
		return a.Less(b)
	} else {
		return i.lesser(x, y)
	}
}

func (i *Interface) Swap(x, y int) {
	i.swapper(x, y)
}

func Sort(v interface{}) {
	sort.Sort(NewInterface(v))
}

func SortFunc(v interface{}, less func(x, y int) bool) {
	value := NewInterface(v)
	value.lesser = less
	sort.Sort(value)
}
