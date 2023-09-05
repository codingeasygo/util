package xmap

import (
	"fmt"
	"testing"
)

type rawMap map[string]interface{}

func (r rawMap) RawMap() map[string]interface{} {
	return r
}

type baseMap map[string]interface{}

func (b baseMap) ValueVal(path ...string) (v interface{}, err error) {
	for _, p := range path {
		if val, ok := b[p]; ok {
			v = val
			return
		}
	}
	err = fmt.Errorf("not exists")
	return
}

// SetValue will set value to path
func (b baseMap) SetValue(path string, val interface{}) (err error) {
	return
}

// Delete will delete value on path
func (b baseMap) Delete(path string) (err error) {
	return
}

// Clear will clear all key on map
func (b baseMap) Clear() (err error) {
	return
}

// Length will return value count
func (b baseMap) Length() (l int) {
	return
}

// Exist will return true if key having
func (b baseMap) Exist(path ...string) bool {
	return false
}

func TestWrap(t *testing.T) {
	var err error
	var wrapVals1 = []interface{}{
		M{"a": 1},
		map[string]interface{}{"b": 2},
		rawMap{"c": 3},
	}
	var parseVals1 = append(wrapVals1, "{}", baseMap{"d": 4}, &impl{})
	var wrapSafeVals1 = append(wrapVals1, baseMap{"d": 4}, &impl{})
	var parseSafeVals1 = append(wrapSafeVals1, "{}")
	var errVals1 = []interface{}{nil, "xx"}
	for _, v := range wrapVals1 {
		Wrap(v)
	}
	for _, v := range parseVals1 {
		_, err = Parse(v)
		if err != nil {
			t.Errorf("err:%v,v:%v", err, v)
			return
		}
	}
	for _, v := range errVals1 {
		func() {
			defer func() {
				if recover() == nil {
					t.Error("error")
				}
			}()
			Wrap(v)
		}()
		_, err = Parse(v)
		if err == nil {
			t.Error("error")
		}
	}
	for _, v := range wrapSafeVals1 {
		WrapSafe(v)
	}
	for _, v := range parseSafeVals1 {
		_, err = ParseSafe(v)
		if err != nil {
			t.Errorf("err:%v,v:%v", err, v)
		}
	}
	for _, v := range errVals1 {
		func() {
			defer func() {
				if recover() == nil {
					t.Error("error")
				}
			}()
			WrapSafe(v)
		}()
		_, err = ParseSafe(v)
		if err == nil {
			t.Error("error")
		}
	}
}

func TestWrapArray(t *testing.T) {
	var err error
	var wrapVals1 = []interface{}{
		nil,
		[]interface{}{nil},
		[]M{{"a": 1}},
		[]map[string]interface{}{{"b": 2}},
		[]rawMap{{"c": 3}},
	}
	var parseVals1 = append(wrapVals1, "[{}]", []baseMap{{"d": 4}}, []*impl{{}})
	var wrapSafeVals1 = append(wrapVals1, []baseMap{{"d": 4}}, []*impl{{}})
	var parseSafeVals1 = append(wrapSafeVals1, "[{}]")
	var errVals1 = []interface{}{"xx", []string{"xx"}}
	for _, v := range wrapVals1 {
		WrapArray(v)
	}
	for _, v := range parseVals1 {
		_, err = ParseArray(v)
		if err != nil {
			t.Errorf("err:%v,v:%v", err, v)
			return
		}
	}
	for _, v := range errVals1 {
		func() {
			defer func() {
				if recover() == nil {
					t.Error("error")
				}
			}()
			WrapArray(v)
		}()
		_, err = ParseArray(v)
		if err == nil {
			t.Errorf("error:%v", v)
		}
	}
	for _, v := range wrapSafeVals1 {
		WrapSafeArray(v)
	}
	for _, v := range parseSafeVals1 {
		_, err = ParseSafeArray(v)
		if err != nil {
			t.Errorf("err:%v,v:%v", err, v)
		}
	}
	for _, v := range errVals1 {
		func() {
			defer func() {
				if recover() == nil {
					t.Error("error")
				}
			}()
			WrapSafeArray(v)
		}()
		_, err = ParseSafeArray(v)
		if err == nil {
			t.Error("error")
		}
	}
}
