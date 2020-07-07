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

//BaseValuable is interface which can be store value
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
}

//RawMapable will get the raw map
type RawMapable interface {
	RawMap() map[string]interface{}
}

//Valuable is interface which can be store and convert value
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
	//Exist will return true if key having
	Exist(path ...string) bool
}

type impl struct {
	BaseValuable
}

//Raw will return raw base valuable
func (i *impl) Raw() BaseValuable {
	return i.BaseValuable
}

//Int will convert path value to int
func (i *impl) Int(path ...string) int {
	return converter.Int(i.Value(path...))
}

//IntDef will convert path value to int, return default when not exist or convert fail
func (i *impl) IntDef(def int, path ...string) (v int) {
	v, err := i.IntVal(path...)
	if err != nil {
		v = def
	}
	return
}

//IntVal will convert path value to int, return error when not exist or convert fail
func (i *impl) IntVal(path ...string) (int, error) {
	return converter.IntVal(i.Value(path...))
}

//Int64 will convert path value to int64
func (i *impl) Int64(path ...string) int64 {
	return converter.Int64(i.Value(path...))
}

//Int64Def will convert path value to int64, return default when not exist or convert fail
func (i *impl) Int64Def(def int64, path ...string) (v int64) {
	v, err := i.Int64Val(path...)
	if err != nil {
		v = def
	}
	return
}

//Int64Val will convert path value to int64, return error when not exist or convert fail
func (i *impl) Int64Val(path ...string) (int64, error) {
	return converter.Int64Val(i.Value(path...))
}

//Uint64 will convert path value to uint64
func (i *impl) Uint64(path ...string) uint64 {
	return converter.Uint64(i.Value(path...))
}

//Uint64Def will convert path value to uint64, return default when not exist or convert fail
func (i *impl) Uint64Def(def uint64, path ...string) (v uint64) {
	v, err := i.Uint64Val(path...)
	if err != nil {
		v = def
	}
	return
}

//Uint64Val will convert path value to uint64, return error when not exist or convert fail
func (i *impl) Uint64Val(path ...string) (uint64, error) {
	return converter.Uint64Val(i.Value(path...))
}

//Float64 will convert path value to float64
func (i *impl) Float64(path ...string) float64 {
	return converter.Float64(i.Value(path...))
}

//Float64Def will convert path value to float64, return default when not exist or convert fail
func (i *impl) Float64Def(def float64, path ...string) (v float64) {
	v, err := i.Float64Val(path...)
	if err != nil {
		v = def
	}
	return
}

//Float64Val will convert path value to float64, return error when not exist or convert fail
func (i *impl) Float64Val(path ...string) (float64, error) {
	return converter.Float64Val(i.Value(path...))
}

//Str will convert path value to string
func (i *impl) Str(path ...string) string {
	return converter.String(i.Value(path...))
}

//StrDef will convert path value to string, return default when not exist or convert fail
func (i *impl) StrDef(def string, path ...string) (v string) {
	v, err := i.StrVal(path...)
	if err != nil {
		v = def
	}
	return
}

//StrVal will convert path value to string, return error when not exist or convert fail
func (i *impl) StrVal(path ...string) (string, error) {
	return converter.StringVal(i.Value(path...))
}

//Map will convert path value to map
func (i *impl) Map(path ...string) M {
	v, _ := MapVal(i.Value(path...))
	return v
}

//MapDef will convert path value to map, return default when not exist or convert fail
func (i *impl) MapDef(def M, path ...string) (v M) {
	v, err := i.MapVal(path...)
	if err != nil {
		v = def
	}
	return
}

//MapVal will convert path value to map, return error when not exist or convert fail
func (i *impl) MapVal(path ...string) (M, error) {
	return MapVal(i.Value(path...))
}

//ArrayDef will convert path value to interface{} array, return default when not exist or convert fail
func (i *impl) ArrayDef(def []interface{}, path ...string) []interface{} {
	vals, err := i.ArrayVal(path...)
	if err != nil {
		vals = def
	}
	return vals
}

//ArrayVal will convert path value to interface{} array, return error when not exist or convert fail
func (i *impl) ArrayVal(path ...string) ([]interface{}, error) {
	return converter.ArrayVal(i.Value(path...))
}

//ArrayMapDef will convert path value to interface{} array, return default when not exist or convert fail
func (i *impl) ArrayMapDef(def []M, path ...string) []M {
	vals, err := i.ArrayMapVal(path...)
	if err != nil || len(vals) < 1 {
		vals = def
	}
	return vals
}

//ArrayMapVal will convert path value to map array, return error when not exist or convert fail
func (i *impl) ArrayMapVal(path ...string) ([]M, error) {
	return ArrayMapVal(i.Value(path...))
}

//ArrayStrDef will convert path value to string array, return default when not exist or convert fail
func (i *impl) ArrayStrDef(def []string, path ...string) []string {
	vals, err := i.ArrayStrVal(path...)
	if err != nil {
		vals = def
	}
	return vals
}

//ArrayStrVal will convert path value to string array, return error when not exist or convert fail
func (i *impl) ArrayStrVal(path ...string) ([]string, error) {
	return converter.ArrayStringVal(i.Value(path...))
}

//ArrayIntDef will convert path value to string array, return default when not exist or convert fail
func (i *impl) ArrayIntDef(def []int, path ...string) []int {
	vals, err := i.ArrayIntVal(path...)
	if err != nil {
		fmt.Println("\n\n--->", err)
		vals = def
	}
	return vals
}

//ArrayIntVal will convert path value to int array, return error when not exist or convert fail
func (i *impl) ArrayIntVal(path ...string) ([]int, error) {
	return converter.ArrayIntVal(i.Value(path...))
}

//ArrayInt64Def will convert path value to int64 array, return default when not exist or convert fail
func (i *impl) ArrayInt64Def(def []int64, path ...string) []int64 {
	vals, err := i.ArrayInt64Val(path...)
	if err != nil {
		vals = def
	}
	return vals
}

//ArrayInt64Val will convert path value to int64 array, return error when not exist or convert fail
func (i *impl) ArrayInt64Val(path ...string) ([]int64, error) {
	return converter.ArrayInt64Val(i.Value(path...))
}

//ArrayUint64Def will convert path value to uint64 array, return default when not exist or convert fail
func (i *impl) ArrayUint64Def(def []uint64, path ...string) []uint64 {
	vals, err := i.ArrayUint64Val(path...)
	if err != nil {
		vals = def
	}
	return vals
}

//ArrayUint64Val will convert path value to uint64 array, return error when not exist or convert fail
func (i *impl) ArrayUint64Val(path ...string) ([]uint64, error) {
	return converter.ArrayUint64Val(i.Value(path...))
}

//ArrayFloat64Def will convert path value to float64 array, return default when not exist or convert fail
func (i *impl) ArrayFloat64Def(def []float64, path ...string) []float64 {
	vals, err := i.ArrayFloat64Val(path...)
	if err != nil {
		vals = def
	}
	return vals
}

//ArrayFloat64Val will convert path value to float64 array, return error when not exist or convert fail
func (i *impl) ArrayFloat64Val(path ...string) ([]float64, error) {
	return converter.ArrayFloat64Val(i.Value(path...))
}

//Value will convert path value to interface{}
func (i *impl) Value(path ...string) interface{} {
	v, _ := i.ValueVal(path...)
	return v
}

//Get is implement for attrvalid.ValueGetter
func (i *impl) Get(path string) (interface{}, error) {
	return i.Value(strings.Split(path, "|")...), nil
}

//ValidFormat is implement for attrvalid.Validable
func (i *impl) ValidFormat(f string, args ...interface{}) error {
	return attrvalid.ValidAttrFormat(f, i, true, args...)
}

//Exist will return true if key having
func (i *impl) Exist(path ...string) bool {
	_, err := i.ValueVal(path...)
	return err == nil
}

//MapVal will conventer value to map
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

//ArrayMapVal will convert value to map array
func ArrayMapVal(v interface{}) (mvals []M, err error) {
	if v == nil {
		return nil, nil
	}
	var mval M
	if vals, ok := v.([]interface{}); ok {
		for _, v := range vals {
			mval, err = MapVal(v)
			if err != nil {
				return
			}
			mvals = append(mvals, mval)
		}
		return
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

//M is type define to map[string]interface{}
type M map[string]interface{}

//Wrap map to Valuable
func (m M) Wrap() Valuable {
	return Wrap(m)
}

//ValueVal will convert path value to value, return error when not exist or convert fail
func (m M) ValueVal(path ...string) (v interface{}, err error) {
	for _, p := range path {
		v, err = m.pathValue(p)
		if err == nil {
			break
		}
	}
	return
}

//SetValue will set value to path
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

//Delete will delete value on path
func (m M) Delete(path string) (err error) {
	return m.setPathValue(path, nil)
}

//Clear will clear all key on map
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
	keys := strings.Split(path, "/")
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

//Length will return value count
func (m M) Length() (l int) {
	l = len(m)
	return
}

//New will create map Valuable
func New() (m Valuable) {
	return &impl{BaseValuable: M{}}
}

//Wrap will wrap raw map to safe
func Wrap(raw interface{}) (m Valuable) {
	if raw == nil {
		panic("raw is nil")
	}
	if base, ok := raw.(BaseValuable); ok {
		m = &impl{BaseValuable: base}
	} else if mval, ok := raw.(map[string]interface{}); ok {
		m = &impl{BaseValuable: M(mval)}
	} else if rm, ok := raw.(RawMapable); ok {
		m = &impl{BaseValuable: M(rm.RawMap())}
	} else {
		panic("not supported type " + reflect.TypeOf(raw).Kind().String())
	}
	return
}

//WrapArray will wrap base values to array
func WrapArray(raws ...interface{}) (ms []Valuable) {
	for _, raw := range raws {
		ms = append(ms, Wrap(raw))
	}
	return
}

//Parse will parse val to map
func Parse(v interface{}) (m Valuable, err error) {
	raw, err := MapVal(v)
	if err == nil {
		m = Wrap(raw)
	}
	return
}

//ParseArray will parse value to map array
func ParseArray(v interface{}) (ms []Valuable, err error) {
	raws, err := ArrayMapVal(v)
	if err == nil {
		for _, raw := range raws {
			ms = append(ms, Wrap(raw))
		}
	}
	return
}

//SafeM is safe map
type SafeM struct {
	Valuable
	raw    BaseValuable
	locker sync.RWMutex
}

//NewSafe will return new safe map
func NewSafe() (m *SafeM) {
	m = &SafeM{
		raw:    M{},
		locker: sync.RWMutex{},
	}
	m.Valuable = &impl{BaseValuable: m}
	return
}

//WrapSafe will wrap raw map to safe
func WrapSafe(raw interface{}) (m *SafeM) {
	if raw == nil {
		panic("raw is nil")
	}
	var b BaseValuable
	if base, ok := raw.(BaseValuable); ok {
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

//WrapSafeArray will wrap raw map to safe
func WrapSafeArray(raws ...interface{}) (ms []Valuable) {
	for _, raw := range raws {
		ms = append(ms, WrapSafe(raw))
	}
	return
}

//ParseSafe will parse val to map
func ParseSafe(v interface{}) (m Valuable, err error) {
	raw, err := MapVal(v)
	if err == nil {
		m = WrapSafe(raw)
	}
	return
}

//ParseSafeArray will parse value to map array
func ParseSafeArray(v interface{}) (ms []Valuable, err error) {
	raws, err := ArrayMapVal(v)
	if err == nil {
		for _, raw := range raws {
			ms = append(ms, WrapSafe(raw))
		}
	}
	return
}

//Raw will return raw base
func (s *SafeM) Raw() BaseValuable {
	return s.raw
}

//ValueVal will convert path value to value, return error when not exist or convert fail
func (s *SafeM) ValueVal(path ...string) (v interface{}, err error) {
	s.locker.RLock()
	defer s.locker.RUnlock()
	return s.raw.ValueVal(path...)
}

//SetValue will set value to path
func (s *SafeM) SetValue(path string, val interface{}) (err error) {
	s.locker.Lock()
	defer s.locker.Unlock()
	return s.raw.SetValue(path, val)
}

//Delete will delete value on path
func (s *SafeM) Delete(path string) (err error) {
	s.locker.Lock()
	defer s.locker.Unlock()
	return s.raw.SetValue(path, nil)
}

//Clear will clear all key
func (s *SafeM) Clear() (err error) {
	s.locker.Lock()
	defer s.locker.Unlock()
	return s.raw.Clear()
}

//Length will return value count
func (s *SafeM) Length() (l int) {
	s.locker.RLock()
	defer s.locker.RUnlock()
	return s.raw.Length()
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
	All  []Valuable
	Path []string
	Type int //0 is int,1 is float,2 is string
	Desc bool
}

func NewMSorter(all []Valuable, vtype int, desc bool, path ...string) *MSorter {
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

// func Maps2Map(ms []Map, path ...string) Map {
// 	var res = Map{}
// 	for _, m := range ms {
// 		res[m.StrValP(path)] = m
// 	}
// 	return res
// }
