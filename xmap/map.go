package xmap

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/codingeasygo/util/attrvalid"

	"github.com/codingeasygo/util/converter"
)

// BaseValuable is interface which can be store value
type BaseValuable interface {
	ValueVal(path ...string) (v interface{}, err error)
	//SetValue will set value to path
	SetValue(path string, val interface{}) (err error)
	//Delete will delete value on path
	Delete(path string) (err error)
	//Clear will clear all key on map
	Clear() (err error)
	//Length will return value count
	Length() (l int)
	//Exist will return true if key having
	Exist(path ...string) bool
}

// RawMapable will get the raw map
type RawMapable interface {
	RawMap() map[string]interface{}
}

// Valuable is interface which can be store and convert value
type Valuable interface {
	attrvalid.ValueGetter
	attrvalid.Validable
	BaseValuable
	//Raw will return raw base valuable
	Raw() BaseValuable
	//Int will convert path value to int
	Int(path ...string) int
	//IntDef will convert path value to int, return default when not exist or convert fail
	IntDef(def int, path ...string) (v int)
	//IntVal will convert path value to int, return error when not exist or convert fail
	IntVal(path ...string) (int, error)
	//Int64 will convert path value to int64
	Int64(path ...string) int64
	//Int64Def will convert path value to int64, return default when not exist or convert fail
	Int64Def(def int64, path ...string) (v int64)
	//Int64Val will convert path value to int64, return error when not exist or convert fail
	Int64Val(path ...string) (int64, error)
	//Uint64 will convert path value to uint64
	Uint64(path ...string) uint64
	//Uint64Def will convert path value to uint64, return default when not exist or convert fail
	Uint64Def(def uint64, path ...string) (v uint64)
	//Uint64Val will convert path value to uint64, return error when not exist or convert fail
	Uint64Val(path ...string) (uint64, error)
	//Float64 will convert path value to float64
	Float64(path ...string) float64
	//Float64Def will convert path value to float64, return default when not exist or convert fail
	Float64Def(def float64, path ...string) (v float64)
	//Float64Val will convert path value to float64, return error when not exist or convert fail
	Float64Val(path ...string) (float64, error)
	//Str will convert path value to string
	Str(path ...string) string
	//StrDef will convert path value to string, return default when not exist or convert fail
	StrDef(def string, path ...string) (v string)
	//StrVal will convert path value to string, return error when not exist or convert fail
	StrVal(path ...string) (string, error)
	//Map will convert path value to map
	Map(path ...string) M
	//MapDef will convert path value to map, return default when not exist or convert fail
	MapDef(def M, path ...string) (v M)
	//MapVal will convert path value to map, return error when not exist or convert fail
	MapVal(path ...string) (M, error)
	//ArrayDef will convert path value to array, return default when not exist or convert fail
	ArrayDef(def []interface{}, path ...string) []interface{}
	//ArrayVal will convert path value to interface{} array, return error when not exist or convert fail
	ArrayVal(path ...string) ([]interface{}, error)
	//ArrayMapDef will convert path value to map array, return default when not exist or convert fail
	ArrayMapDef(def []M, path ...string) []M
	//ArrayMapVal will convert path value to map array, return error when not exist or convert fail
	ArrayMapVal(path ...string) ([]M, error)
	//ArrayStrDef will convert path value to string array, return default when not exist or convert fail
	ArrayStrDef(def []string, path ...string) []string
	//ArrayStrVal will convert path value to string array, return error when not exist or convert fail
	ArrayStrVal(path ...string) ([]string, error)
	//ArrayIntDef will convert path value to int array, return default when not exist or convert fail
	ArrayIntDef(def []int, path ...string) []int
	//ArrayIntVal will convert path value to int array, return error when not exist or convert fail
	ArrayIntVal(path ...string) ([]int, error)
	//ArrayInt64Def will convert path value to int64 array, return default when not exist or convert fail
	ArrayInt64Def(def []int64, path ...string) []int64
	//ArrayInt64Val will convert path value to int64 array, return error when not exist or convert fail
	ArrayInt64Val(path ...string) ([]int64, error)
	//ArrayUint64Def will convert path value to uint64 array, return default when not exist or convert fail
	ArrayUint64Def(def []uint64, path ...string) []uint64
	//ArrayUint64Val will convert path value to uint64 array, return error when not exist or convert fail
	ArrayUint64Val(path ...string) ([]uint64, error)
	//ArrayFloat64Def will convert path value to float64 array, return default when not exist or convert fail
	ArrayFloat64Def(def []float64, path ...string) []float64
	//ArrayFloat64Val will convert path value to float64 array, return error when not exist or convert fail
	ArrayFloat64Val(path ...string) ([]float64, error)
	//Value will convert path value to interface{}
	Value(path ...string) interface{}
	//ReplaceAll will replace all value by ${path}
	ReplaceAll(in string, usingEnv, keepEmpty bool) (out string)
}

type impl struct {
	BaseValuable
}

// Raw will return raw base valuable
func (i *impl) Raw() BaseValuable {
	return i.BaseValuable
}

// Int will convert path value to int
func (i *impl) Int(path ...string) int {
	return converter.Int(i.Value(path...))
}

// IntDef will convert path value to int, return default when not exist or convert fail
func (i *impl) IntDef(def int, path ...string) (v int) {
	v, err := i.IntVal(path...)
	if err != nil {
		v = def
	}
	return
}

// IntVal will convert path value to int, return error when not exist or convert fail
func (i *impl) IntVal(path ...string) (int, error) {
	return converter.IntVal(i.Value(path...))
}

// Int64 will convert path value to int64
func (i *impl) Int64(path ...string) int64 {
	return converter.Int64(i.Value(path...))
}

// Int64Def will convert path value to int64, return default when not exist or convert fail
func (i *impl) Int64Def(def int64, path ...string) (v int64) {
	v, err := i.Int64Val(path...)
	if err != nil {
		v = def
	}
	return
}

// Int64Val will convert path value to int64, return error when not exist or convert fail
func (i *impl) Int64Val(path ...string) (int64, error) {
	return converter.Int64Val(i.Value(path...))
}

// Uint64 will convert path value to uint64
func (i *impl) Uint64(path ...string) uint64 {
	return converter.Uint64(i.Value(path...))
}

// Uint64Def will convert path value to uint64, return default when not exist or convert fail
func (i *impl) Uint64Def(def uint64, path ...string) (v uint64) {
	v, err := i.Uint64Val(path...)
	if err != nil {
		v = def
	}
	return
}

// Uint64Val will convert path value to uint64, return error when not exist or convert fail
func (i *impl) Uint64Val(path ...string) (uint64, error) {
	return converter.Uint64Val(i.Value(path...))
}

// Float64 will convert path value to float64
func (i *impl) Float64(path ...string) float64 {
	return converter.Float64(i.Value(path...))
}

// Float64Def will convert path value to float64, return default when not exist or convert fail
func (i *impl) Float64Def(def float64, path ...string) (v float64) {
	v, err := i.Float64Val(path...)
	if err != nil {
		v = def
	}
	return
}

// Float64Val will convert path value to float64, return error when not exist or convert fail
func (i *impl) Float64Val(path ...string) (float64, error) {
	return converter.Float64Val(i.Value(path...))
}

// Str will convert path value to string
func (i *impl) Str(path ...string) string {
	return converter.String(i.Value(path...))
}

// StrDef will convert path value to string, return default when not exist or convert fail
func (i *impl) StrDef(def string, path ...string) (v string) {
	v, err := i.StrVal(path...)
	if err != nil {
		v = def
	}
	return
}

// StrVal will convert path value to string, return error when not exist or convert fail
func (i *impl) StrVal(path ...string) (string, error) {
	return converter.StringVal(i.Value(path...))
}

// Map will convert path value to map
func (i *impl) Map(path ...string) M {
	v, _ := MapVal(i.Value(path...))
	return v
}

// MapDef will convert path value to map, return default when not exist or convert fail
func (i *impl) MapDef(def M, path ...string) (v M) {
	v, err := i.MapVal(path...)
	if err != nil {
		v = def
	}
	return
}

// MapVal will convert path value to map, return error when not exist or convert fail
func (i *impl) MapVal(path ...string) (M, error) {
	return MapVal(i.Value(path...))
}

// ArrayDef will convert path value to interface{} array, return default when not exist or convert fail
func (i *impl) ArrayDef(def []interface{}, path ...string) []interface{} {
	vals, err := i.ArrayVal(path...)
	if err != nil {
		vals = def
	}
	return vals
}

// ArrayVal will convert path value to interface{} array, return error when not exist or convert fail
func (i *impl) ArrayVal(path ...string) ([]interface{}, error) {
	return converter.ArrayVal(i.Value(path...))
}

// ArrayMapDef will convert path value to interface{} array, return default when not exist or convert fail
func (i *impl) ArrayMapDef(def []M, path ...string) []M {
	vals, err := i.ArrayMapVal(path...)
	if err != nil || len(vals) < 1 {
		vals = def
	}
	return vals
}

// ArrayMapVal will convert path value to map array, return error when not exist or convert fail
func (i *impl) ArrayMapVal(path ...string) ([]M, error) {
	return ArrayMapVal(i.Value(path...))
}

// ArrayStrDef will convert path value to string array, return default when not exist or convert fail
func (i *impl) ArrayStrDef(def []string, path ...string) []string {
	vals, err := i.ArrayStrVal(path...)
	if err != nil {
		vals = def
	}
	return vals
}

// ArrayStrVal will convert path value to string array, return error when not exist or convert fail
func (i *impl) ArrayStrVal(path ...string) ([]string, error) {
	return converter.ArrayStringVal(i.Value(path...))
}

// ArrayIntDef will convert path value to string array, return default when not exist or convert fail
func (i *impl) ArrayIntDef(def []int, path ...string) []int {
	vals, err := i.ArrayIntVal(path...)
	if err != nil {
		vals = def
	}
	return vals
}

// ArrayIntVal will convert path value to int array, return error when not exist or convert fail
func (i *impl) ArrayIntVal(path ...string) ([]int, error) {
	return converter.ArrayIntVal(i.Value(path...))
}

// ArrayInt64Def will convert path value to int64 array, return default when not exist or convert fail
func (i *impl) ArrayInt64Def(def []int64, path ...string) []int64 {
	vals, err := i.ArrayInt64Val(path...)
	if err != nil {
		vals = def
	}
	return vals
}

// ArrayInt64Val will convert path value to int64 array, return error when not exist or convert fail
func (i *impl) ArrayInt64Val(path ...string) ([]int64, error) {
	return converter.ArrayInt64Val(i.Value(path...))
}

// ArrayUint64Def will convert path value to uint64 array, return default when not exist or convert fail
func (i *impl) ArrayUint64Def(def []uint64, path ...string) []uint64 {
	vals, err := i.ArrayUint64Val(path...)
	if err != nil {
		vals = def
	}
	return vals
}

// ArrayUint64Val will convert path value to uint64 array, return error when not exist or convert fail
func (i *impl) ArrayUint64Val(path ...string) ([]uint64, error) {
	return converter.ArrayUint64Val(i.Value(path...))
}

// ArrayFloat64Def will convert path value to float64 array, return default when not exist or convert fail
func (i *impl) ArrayFloat64Def(def []float64, path ...string) []float64 {
	vals, err := i.ArrayFloat64Val(path...)
	if err != nil {
		vals = def
	}
	return vals
}

// ArrayFloat64Val will convert path value to float64 array, return error when not exist or convert fail
func (i *impl) ArrayFloat64Val(path ...string) ([]float64, error) {
	return converter.ArrayFloat64Val(i.Value(path...))
}

// Value will convert path value to interface{}
func (i *impl) Value(path ...string) interface{} {
	v, _ := i.ValueVal(path...)
	return v
}

// Get is implement for attrvalid.ValueGetter
func (i *impl) Get(path string) (interface{}, error) {
	return i.Value(strings.Split(path, "|")...), nil
}

// ValidFormat is implement for attrvalid.Validable
func (i *impl) ValidFormat(f string, args ...interface{}) error {
	return attrvalid.ValidAttrFormat(f, i, true, args...)
}

// ReplaceAll will replace value by ${xx}, which xx is having
func (i *impl) ReplaceAll(in string, usingEnv, keepEmpty bool) (out string) {
	out = ReplaceAll(func(key string) interface{} {
		v, err := i.ValueVal(key)
		if err != nil {
			return nil
		}
		return v
	}, in, usingEnv, keepEmpty)
	return
}

// MapVal will conventer value to map
func MapVal(v interface{}) (M, error) {
	if mv, ok := v.(M); ok {
		return mv, nil
	} else if mv, ok := v.(map[string]interface{}); ok {
		return M(mv), nil
	} else if rm, ok := v.(RawMapable); ok {
		mv = M(rm.RawMap())
		return mv, nil
	} else if sv, ok := v.(string); ok {
		mv := M{}
		return mv, json.Unmarshal([]byte(sv), &mv)
	} else if sv, ok := v.([]byte); ok {
		mv := M{}
		return mv, json.Unmarshal(sv, &mv)
	} else {
		return nil, fmt.Errorf("incompactable kind(%v)", reflect.ValueOf(v).Kind())
	}
}

// ArrayMapVal will convert value to map array
func ArrayMapVal(v interface{}) (mvals []M, err error) {
	if v == nil {
		return nil, nil
	}
	var mval M
	if vals, ok := v.([]M); ok {
		return vals, nil
	} else if vals, ok := v.([]interface{}); ok {
		for _, v := range vals {
			mval, err = MapVal(v)
			if err != nil {
				return
			}
			mvals = append(mvals, mval)
		}
		return
	} else if vals, ok := v.([]map[string]interface{}); ok {
		for _, v := range vals {
			mvals = append(mvals, M(v))
		}
	}
	if sv, ok := v.(string); ok {
		err = json.Unmarshal([]byte(sv), &mvals)
		return
	}
	if sv, ok := v.([]byte); ok {
		err = json.Unmarshal(sv, &mvals)
		return
	}
	vals := reflect.ValueOf(v)
	if vals.Kind() != reflect.Slice {
		err = fmt.Errorf("incompactable kind(%v)", vals.Kind())
		return
	}
	for i := 0; i < vals.Len(); i++ {
		if vals.Index(i).IsZero() {
			mval, err = MapVal(nil)
		} else {
			mval, err = MapVal(vals.Index(i).Interface())
		}
		if err != nil {
			return
		}
		mvals = append(mvals, mval)
	}
	return
}

// M is type define to map[string]interface{}
type M map[string]interface{}

// New will create map Valuable
func New() (m M) {
	return M{}
}

// ValueVal will convert path value to value, return error when not exist or convert fail
func (m M) ValueVal(path ...string) (v interface{}, err error) {
	for _, p := range path {
		v, err = m.pathValue(p)
		if err == nil {
			break
		}
	}
	return
}

// SetValue will set value to path
func (m M) SetValue(path string, val interface{}) (err error) {
	if len(path) < 1 {
		err = fmt.Errorf("path is empty")
		return
	}
	path = strings.TrimPrefix(path, "/")
	keys := strings.Split(path, "/")
	if len(keys) == 1 {
		m.setValue(keys[0], val)
	} else {
		err = m.setPathValue(path, val)
	}
	return
}

// Delete will delete value on path
func (m M) Delete(path string) (err error) {
	return m.setPathValue(path, nil)
}

// Clear will clear all key on map
func (m M) Clear() (err error) {
	for key := range m {
		delete(m, key)
	}
	return
}

func (m M) setValue(path string, val interface{}) {
	if val == nil {
		delete(m, path)
	} else {
		m[path] = val
	}
}

func (m M) pathValue(path string) (interface{}, error) {
	if v, ok := m[path]; ok {
		return v, nil
	}
	path = strings.TrimPrefix(path, "/")
	keys := []string{}
	key := strings.Builder{}
	wrapped := false
	for _, c := range path {
		switch c {
		case '\\':
			wrapped = true
		case '/':
			if wrapped {
				key.WriteRune(c)
				wrapped = false
			} else {
				keys = append(keys, key.String())
				key = strings.Builder{}
			}
		default:
			key.WriteRune(c)
		}
	}
	if key.Len() > 0 {
		keys = append(keys, key.String())
		key = strings.Builder{}
	}
	return m.valP(keys...)
}

func (m M) valP(keys ...string) (interface{}, error) {
	count := len(keys)
	var tv interface{} = m
	for i := 0; i < count; i++ {
		if tv == nil {
			break
		}
		switch reflect.TypeOf(tv).Kind() {
		case reflect.Slice: //if array
			ary, ok := tv.([]interface{}) //check if valid array
			if !ok {
				return nil, fmt.Errorf(
					"invalid array(%v) in path(/%v),expected []interface{}",
					reflect.TypeOf(tv).String(), strings.Join(keys[:i+1], "/"),
				)
			}
			if keys[i] == "@len" { //check if have @len
				return len(ary), nil //return the array length
			}
			idx, err := strconv.Atoi(keys[i]) //get the target index.
			if err != nil {
				return nil, fmt.Errorf(
					"invalid array index(/%v)", strings.Join(keys[:i+1], "/"),
				)
			}
			if idx >= len(ary) || idx < 0 { //check index valid
				return nil, fmt.Errorf(
					"array out of index in path(/%v)", strings.Join(keys[:i+1], "/"),
				)
			}
			tv = ary[idx]
			continue
		case reflect.Map: //if map
			tm, _ := MapVal(tv) //check map covert
			if tm == nil {
				return nil, fmt.Errorf(
					"invalid map in path(/%v)", strings.Join(keys[:i], "/"),
				)
			}
			tv = tm[keys[i]]
			continue
		default: //unknow type
			return nil, fmt.Errorf(
				"invalid type(%v) in path(/%v)",
				reflect.TypeOf(tv).Kind(), strings.Join(keys[:i], "/"),
			)
		}
	}
	if tv == nil { //if valud not found
		return nil, fmt.Errorf(
			"value not found in path(/%v)", strings.Join(keys, "/"),
		)
	}
	return tv, nil
}

func (m M) setPathValue(path string, val interface{}) error {
	if len(path) < 1 {
		return errors.New("path is empty")
	}
	path = strings.TrimPrefix(path, "/")
	keys := strings.Split(path, "/")
	//
	i := len(keys) - 1
	pv, err := m.valP(keys[:i]...)
	if err != nil {
		return err
	}
	switch reflect.TypeOf(pv).Kind() {
	case reflect.Slice:
		ary, ok := pv.([]interface{}) //check if valid array
		if !ok {
			return fmt.Errorf(
				"invalid array(%v) in path(/%v),expected []interface{}",
				reflect.TypeOf(pv).String(), strings.Join(keys[:i+1], "/"),
			)
		}
		idx, err := strconv.Atoi(keys[i]) //get the target index.
		if err != nil {
			return fmt.Errorf(
				"invalid array index(/%v)", strings.Join(keys[:i+1], "/"),
			)
		}
		if idx >= len(ary) || idx < 0 { //check index valid
			return fmt.Errorf(
				"array out of index in path(/%v)", strings.Join(keys[:i+1], "/"),
			)
		}
		ary[idx] = val
	case reflect.Map:
		tm, _ := MapVal(pv) //check map covert
		if tm == nil {
			return fmt.Errorf(
				"invalid map in path(/%v)", strings.Join(keys[:i], "/"),
			)
		}
		tm.setValue(keys[i], val)
	default: //unknow type
		return fmt.Errorf(
			"not map type(%v) in path(/%v)",
			reflect.TypeOf(pv).Kind(), strings.Join(keys[:i], "/"),
		)
	}
	return nil
}

// Length will return value count
func (m M) Length() (l int) {
	l = len(m)
	return
}

// Exist will return true if key having
func (m M) Exist(path ...string) bool {
	_, err := m.ValueVal(path...)
	return err == nil
}

// Raw will return raw base valuable
func (m M) Raw() BaseValuable {
	return m
}

// Int will convert path value to int
func (m M) Int(path ...string) int {
	return converter.Int(m.Value(path...))
}

// IntDef will convert path value to int, return default when not exist or convert fail
func (m M) IntDef(def int, path ...string) (v int) {
	v, err := m.IntVal(path...)
	if err != nil {
		v = def
	}
	return
}

// IntVal will convert path value to int, return error when not exist or convert fail
func (m M) IntVal(path ...string) (int, error) {
	return converter.IntVal(m.Value(path...))
}

// Int64 will convert path value to int64
func (m M) Int64(path ...string) int64 {
	return converter.Int64(m.Value(path...))
}

// Int64Def will convert path value to int64, return default when not exist or convert fail
func (m M) Int64Def(def int64, path ...string) (v int64) {
	v, err := m.Int64Val(path...)
	if err != nil {
		v = def
	}
	return
}

// Int64Val will convert path value to int64, return error when not exist or convert fail
func (m M) Int64Val(path ...string) (int64, error) {
	return converter.Int64Val(m.Value(path...))
}

// Uint64 will convert path value to uint64
func (m M) Uint64(path ...string) uint64 {
	return converter.Uint64(m.Value(path...))
}

// Uint64Def will convert path value to uint64, return default when not exist or convert fail
func (m M) Uint64Def(def uint64, path ...string) (v uint64) {
	v, err := m.Uint64Val(path...)
	if err != nil {
		v = def
	}
	return
}

// Uint64Val will convert path value to uint64, return error when not exist or convert fail
func (m M) Uint64Val(path ...string) (uint64, error) {
	return converter.Uint64Val(m.Value(path...))
}

// Float64 will convert path value to float64
func (m M) Float64(path ...string) float64 {
	return converter.Float64(m.Value(path...))
}

// Float64Def will convert path value to float64, return default when not exist or convert fail
func (m M) Float64Def(def float64, path ...string) (v float64) {
	v, err := m.Float64Val(path...)
	if err != nil {
		v = def
	}
	return
}

// Float64Val will convert path value to float64, return error when not exist or convert fail
func (m M) Float64Val(path ...string) (float64, error) {
	return converter.Float64Val(m.Value(path...))
}

// Str will convert path value to string
func (m M) Str(path ...string) string {
	return converter.String(m.Value(path...))
}

// StrDef will convert path value to string, return default when not exist or convert fail
func (m M) StrDef(def string, path ...string) (v string) {
	v, err := m.StrVal(path...)
	if err != nil {
		v = def
	}
	return
}

// StrVal will convert path value to string, return error when not exist or convert fail
func (m M) StrVal(path ...string) (string, error) {
	return converter.StringVal(m.Value(path...))
}

// Map will convert path value to map
func (m M) Map(path ...string) M {
	v, _ := MapVal(m.Value(path...))
	return v
}

// MapDef will convert path value to map, return default when not exist or convert fail
func (m M) MapDef(def M, path ...string) (v M) {
	v, err := m.MapVal(path...)
	if err != nil {
		v = def
	}
	return
}

// MapVal will convert path value to map, return error when not exist or convert fail
func (m M) MapVal(path ...string) (M, error) {
	return MapVal(m.Value(path...))
}

// ArrayDef will convert path value to interface{} array, return default when not exist or convert fail
func (m M) ArrayDef(def []interface{}, path ...string) []interface{} {
	vals, err := m.ArrayVal(path...)
	if err != nil {
		vals = def
	}
	return vals
}

// ArrayVal will convert path value to interface{} array, return error when not exist or convert fail
func (m M) ArrayVal(path ...string) ([]interface{}, error) {
	return converter.ArrayVal(m.Value(path...))
}

// ArrayMapDef will convert path value to interface{} array, return default when not exist or convert fail
func (m M) ArrayMapDef(def []M, path ...string) []M {
	vals, err := m.ArrayMapVal(path...)
	if err != nil || len(vals) < 1 {
		vals = def
	}
	return vals
}

// ArrayMapVal will convert path value to map array, return error when not exist or convert fail
func (m M) ArrayMapVal(path ...string) ([]M, error) {
	return ArrayMapVal(m.Value(path...))
}

// ArrayStrDef will convert path value to string array, return default when not exist or convert fail
func (m M) ArrayStrDef(def []string, path ...string) []string {
	vals, err := m.ArrayStrVal(path...)
	if err != nil {
		vals = def
	}
	return vals
}

// ArrayStrVal will convert path value to string array, return error when not exist or convert fail
func (m M) ArrayStrVal(path ...string) ([]string, error) {
	return converter.ArrayStringVal(m.Value(path...))
}

// ArrayIntDef will convert path value to string array, return default when not exist or convert fail
func (m M) ArrayIntDef(def []int, path ...string) []int {
	vals, err := m.ArrayIntVal(path...)
	if err != nil {
		vals = def
	}
	return vals
}

// ArrayIntVal will convert path value to int array, return error when not exist or convert fail
func (m M) ArrayIntVal(path ...string) ([]int, error) {
	return converter.ArrayIntVal(m.Value(path...))
}

// ArrayInt64Def will convert path value to int64 array, return default when not exist or convert fail
func (m M) ArrayInt64Def(def []int64, path ...string) []int64 {
	vals, err := m.ArrayInt64Val(path...)
	if err != nil {
		vals = def
	}
	return vals
}

// ArrayInt64Val will convert path value to int64 array, return error when not exist or convert fail
func (m M) ArrayInt64Val(path ...string) ([]int64, error) {
	return converter.ArrayInt64Val(m.Value(path...))
}

// ArrayUint64Def will convert path value to uint64 array, return default when not exist or convert fail
func (m M) ArrayUint64Def(def []uint64, path ...string) []uint64 {
	vals, err := m.ArrayUint64Val(path...)
	if err != nil {
		vals = def
	}
	return vals
}

// ArrayUint64Val will convert path value to uint64 array, return error when not exist or convert fail
func (m M) ArrayUint64Val(path ...string) ([]uint64, error) {
	return converter.ArrayUint64Val(m.Value(path...))
}

// ArrayFloat64Def will convert path value to float64 array, return default when not exist or convert fail
func (m M) ArrayFloat64Def(def []float64, path ...string) []float64 {
	vals, err := m.ArrayFloat64Val(path...)
	if err != nil {
		vals = def
	}
	return vals
}

// ArrayFloat64Val will convert path value to float64 array, return error when not exist or convert fail
func (m M) ArrayFloat64Val(path ...string) ([]float64, error) {
	return converter.ArrayFloat64Val(m.Value(path...))
}

// Value will convert path value to interface{}
func (m M) Value(path ...string) interface{} {
	v, _ := m.ValueVal(path...)
	return v
}

// Get is implement for attrvalid.ValueGetter
func (m M) Get(path string) (interface{}, error) {
	return m.Value(strings.Split(path, "|")...), nil
}

// ValidFormat is implement for attrvalid.Validable
func (m M) ValidFormat(f string, args ...interface{}) error {
	return attrvalid.ValidAttrFormat(f, m, true, args...)
}

// ReplaceAll will replace input string by ${xx}, which xx is in values,
// if usingEnv is true, xx will check use env when vals is not having xx,
// if usingEmpty is true, xx will check use empty string when vals is not having xx and env is not exist
func (m M) ReplaceAll(in string, usingEnv, usingEmpty bool) (out string) {
	out = ReplaceAll(func(key string) interface{} {
		v, err := m.ValueVal(key)
		if err != nil {
			return nil
		}
		return v
	}, in, usingEnv, usingEmpty)
	return
}

func ValueEqual(a, b interface{}) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	resultValue := reflect.ValueOf(a)
	if !resultValue.CanConvert(reflect.TypeOf(b)) {
		return false
	}
	targetValue := resultValue.Convert(reflect.TypeOf(b)).Interface()
	return reflect.DeepEqual(targetValue, b)
}

func (m M) ValueEqual(key string, value interface{}) bool {
	return ValueEqual(m.Value(key), value)
}

// SafeM is safe map
type SafeM struct {
	Valuable
	raw    BaseValuable
	locker sync.RWMutex
}

// NewSafe will return new safe map
func NewSafe() (m *SafeM) {
	m = &SafeM{
		raw:    M{},
		locker: sync.RWMutex{},
	}
	m.Valuable = &impl{BaseValuable: m}
	return
}

// Raw will return raw base
func (s *SafeM) Raw() BaseValuable {
	return s.raw
}

// ValueVal will convert path value to value, return error when not exist or convert fail
func (s *SafeM) ValueVal(path ...string) (v interface{}, err error) {
	s.locker.RLock()
	defer s.locker.RUnlock()
	return s.raw.ValueVal(path...)
}

// SetValue will set value to path
func (s *SafeM) SetValue(path string, val interface{}) (err error) {
	s.locker.Lock()
	defer s.locker.Unlock()
	return s.raw.SetValue(path, val)
}

// Delete will delete value on path
func (s *SafeM) Delete(path string) (err error) {
	s.locker.Lock()
	defer s.locker.Unlock()
	return s.raw.SetValue(path, nil)
}

// Clear will clear all key
func (s *SafeM) Clear() (err error) {
	s.locker.Lock()
	defer s.locker.Unlock()
	return s.raw.Clear()
}

// Length will return value count
func (s *SafeM) Length() (l int) {
	s.locker.RLock()
	defer s.locker.RUnlock()
	return s.raw.Length()
}

// Exist will return true if key having
func (s *SafeM) Exist(path ...string) bool {
	s.locker.RLock()
	defer s.locker.RUnlock()
	return s.raw.Exist(path...)
}

// ReplaceAll will replace input string by ${xx}, which xx is in values,
// if usingEnv is true, xx will check use env when vals is not having xx,
// if usingEmpty is true, xx will check use empty string when vals is not having xx and env is not exist
func (s *SafeM) ReplaceAll(in string, usingEnv, usingEmpty bool) (out string) {
	out = ReplaceAll(func(key string) interface{} {
		v, err := s.ValueVal(key)
		if err != nil {
			return nil
		}
		return v
	}, in, usingEnv, usingEmpty)
	return
}

// func (m Map) ToS(dest interface{}) {
// 	M2S(m, dest)
// }

// func NewMap(f string) (Map, error) {
// 	bys, err := ioutil.ReadFile(f)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var kvs Map = Map{}
// 	err = json.Unmarshal(bys, &kvs)
// 	return kvs, err
// }

// func NewMaps(f string) ([]Map, error) {
// 	bys, err := ioutil.ReadFile(f)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var kvs []Map = []Map{}
// 	err = json.Unmarshal(bys, &kvs)
// 	return kvs, err
// }

type MSorter struct {
	All  []M
	Path []string
	Type int //0 is int,1 is float,2 is string
	Desc bool
}

func NewMSorter(all []M, vtype int, desc bool, path ...string) *MSorter {
	return &MSorter{
		All:  all,
		Path: path,
		Type: vtype,
		Desc: desc,
	}
}
func (m *MSorter) Len() int {
	return len(m.All)
}

func (m *MSorter) Less(i, j int) bool {
	switch m.Type {
	case 0:
		if m.Desc {
			return m.All[i].Int(m.Path...) > m.All[j].Int(m.Path...)
		}
		return m.All[i].Int(m.Path...) < m.All[j].Int(m.Path...)
	case 1:
		if m.Desc {
			return m.All[i].Float64(m.Path...) > m.All[j].Float64(m.Path...)
		}
		return m.All[i].Float64(m.Path...) < m.All[j].Float64(m.Path...)
	default:
		if m.Desc {
			return m.All[i].Str(m.Path...) > m.All[j].Str(m.Path...)
		}
		return m.All[i].Str(m.Path...) < m.All[j].Str(m.Path...)
	}
}
func (m *MSorter) Swap(i, j int) {
	m.All[i], m.All[j] = m.All[j], m.All[i]
}

type ValuableSorter struct {
	All  []Valuable
	Path []string
	Type int //0 is int,1 is float,2 is string
	Desc bool
}

func NewValuableSorter(all []Valuable, vtype int, desc bool, path ...string) *ValuableSorter {
	return &ValuableSorter{
		All:  all,
		Path: path,
		Type: vtype,
		Desc: desc,
	}
}
func (v *ValuableSorter) Len() int {
	return len(v.All)
}

func (v *ValuableSorter) Less(i, j int) bool {
	switch v.Type {
	case 0:
		if v.Desc {
			return v.All[i].Int(v.Path...) > v.All[j].Int(v.Path...)
		}
		return v.All[i].Int(v.Path...) < v.All[j].Int(v.Path...)
	case 1:
		if v.Desc {
			return v.All[i].Float64(v.Path...) > v.All[j].Float64(v.Path...)
		}
		return v.All[i].Float64(v.Path...) < v.All[j].Float64(v.Path...)
	default:
		if v.Desc {
			return v.All[i].Str(v.Path...) > v.All[j].Str(v.Path...)
		}
		return v.All[i].Str(v.Path...) < v.All[j].Str(v.Path...)
	}
}
func (v *ValuableSorter) Swap(i, j int) {
	v.All[i], v.All[j] = v.All[j], v.All[i]
}

// func Maps2Map(ms []Map, path ...string) Map {
// 	var res = Map{}
// 	for _, m := range ms {
// 		res[m.StrValP(path)] = m
// 	}
// 	return res
// }
