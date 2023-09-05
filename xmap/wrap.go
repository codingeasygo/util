package xmap

import (
	"fmt"
	"reflect"
	"sync"
)

// Wrap will wrap raw map to safe
func Wrap(v interface{}) (m M) {
	if v == nil {
		panic("v is nil")
	}
	if mval, ok := v.(M); ok {
		m = mval
	} else if mval, ok := v.(map[string]interface{}); ok {
		m = M(mval)
	} else if rm, ok := v.(RawMapable); ok {
		m = M(rm.RawMap())
	} else {
		panic(fmt.Sprintf("not supported type %v", reflect.TypeOf(v)))
	}
	return
}

// WrapArray will wrap base values to array
func WrapArray(v interface{}) (ms []M) {
	if v == nil {
		return nil
	}
	vals := reflect.ValueOf(v)
	if vals.Kind() != reflect.Slice {
		panic(fmt.Errorf("incompactable kind(%v)", vals.Kind()))
	}
	for i := 0; i < vals.Len(); i++ {
		if vals.Index(i).IsZero() {
			ms = append(ms, nil)
		} else {
			ms = append(ms, Wrap(vals.Index(i).Interface()))
		}
	}
	return
}

// Parse will parse val to map
func Parse(v interface{}) (m Valuable, err error) {
	if v == nil {
		err = fmt.Errorf("v is nil")
		return
	}
	m, err = MapVal(v)
	if err == nil {
		return
	}
	err = nil
	if val, ok := v.(Valuable); ok {
		m = val
	} else if base, ok := v.(BaseValuable); ok {
		m = &impl{BaseValuable: base}
	} else {
		err = fmt.Errorf("not supported type %v", reflect.TypeOf(v))
	}
	return
}

// ParseArray will parse value to map array
func ParseArray(v interface{}) (ms []Valuable, err error) {
	if v == nil {
		return
	}
	raws, err := ArrayMapVal(v)
	if err == nil {
		for _, raw := range raws {
			ms = append(ms, raw)
		}
		return
	}
	err = nil
	vals := reflect.ValueOf(v)
	if vals.Kind() != reflect.Slice {
		err = fmt.Errorf("incompactable kind(%v)", vals.Kind())
		return
	}
	var m Valuable
	for i := 0; i < vals.Len(); i++ {
		if vals.Index(i).IsZero() {
			m = nil
		} else {
			m, err = Parse(vals.Index(i).Interface())
		}
		if err != nil {
			break
		}
		ms = append(ms, m)
	}
	return
}

// WrapSafe will wrap raw map to safe
func WrapSafe(raw interface{}) (m *SafeM) {
	if raw == nil {
		panic("raw is nil")
	}
	var b BaseValuable
	if m, ok := raw.(M); ok {
		b = m
	} else if base, ok := raw.(BaseValuable); ok {
		b = base
	} else if mval, ok := raw.(map[string]interface{}); ok {
		b = M(mval)
	} else if rm, ok := raw.(RawMapable); ok {
		b = M(rm.RawMap())
	} else {
		panic("not supported type " + reflect.TypeOf(raw).Kind().String())
	}
	m = &SafeM{
		raw:    b,
		locker: sync.RWMutex{},
	}
	m.Valuable = &impl{BaseValuable: m}
	return
}

// ParseSafe will parse val to map
func ParseSafe(v interface{}) (m Valuable, err error) {
	if v == nil {
		err = fmt.Errorf("v is nil")
		return
	}
	m, err = MapVal(v)
	if err == nil {
		m = WrapSafe(m)
		return
	}
	err = nil
	if base, ok := v.(BaseValuable); ok {
		m = WrapSafe(base)
	} else {
		err = fmt.Errorf("not supported type %v", reflect.TypeOf(v))
	}
	return
}

// WrapSafeArray will wrap raw map to safe
func WrapSafeArray(v interface{}) (ms []Valuable) {
	if v == nil {
		return nil
	}
	vals := reflect.ValueOf(v)
	if vals.Kind() != reflect.Slice {
		panic(fmt.Errorf("incompactable kind(%v)", vals.Kind()))
	}
	for i := 0; i < vals.Len(); i++ {
		if vals.Index(i).IsZero() {
			ms = append(ms, nil)
		} else {
			ms = append(ms, WrapSafe(vals.Index(i).Interface()))
		}
	}
	return
}

// ParseSafeArray will parse value to map array
func ParseSafeArray(v interface{}) (ms []Valuable, err error) {
	if v == nil {
		return
	}
	raws, err := ArrayMapVal(v)
	if err == nil {
		for _, raw := range raws {
			ms = append(ms, WrapSafe(raw))
		}
		return
	}
	err = nil
	vals := reflect.ValueOf(v)
	if vals.Kind() != reflect.Slice {
		err = fmt.Errorf("incompactable kind(%v)", vals.Kind())
		return
	}
	var m Valuable
	for i := 0; i < vals.Len(); i++ {
		if vals.Index(i).IsZero() {
			m = nil
		} else {
			m, err = ParseSafe(vals.Index(i).Interface())
		}
		if err != nil {
			break
		}
		ms = append(ms, m)
	}
	return
}
