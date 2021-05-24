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
	a := i.value.Index(x).Interface().(Lesser)
	b := i.value.Index(y).Interface()
	return a.Less(b)
}

func (i *Interface) Swap(x, y int) {
	i.swapper(x, y)
}

func Sort(v interface{}) {
	sort.Sort(NewInterface(v))
}
