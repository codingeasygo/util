package xmap

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/codingeasygo/util/attrvalid"

	"github.com/codingeasygo/util/converter"
)

//M is type define to map[string]interface{}
type M map[string]interface{}

func MapVal(v interface{}) (M, error) {
	if mv, ok := v.(M); ok {
		return mv, nil
	} else if mv, ok := v.(map[string]interface{}); ok {
		return M(mv), nil
	} else {
		return nil, fmt.Errorf("incompactable kind(%v)", reflect.ValueOf(v).Kind())
	}
}

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

func (m M) Int(path ...string) int {
	return converter.Int(m.Value(path...))
}

func (m M) IntDef(def int, path ...string) (v int) {
	v, err := m.IntVal(path...)
	if err != nil {
		v = def
	}
	return
}

func (m M) IntVal(path ...string) (int, error) {
	return converter.IntVal(m.Value(path...))
}

func (m M) Int64(path ...string) int64 {
	return converter.Int64(m.Value(path...))
}

func (m M) Int64Def(def int64, path ...string) (v int64) {
	v, err := m.Int64Val(path...)
	if err != nil {
		v = def
	}
	return
}

func (m M) Int64Val(path ...string) (int64, error) {
	return converter.Int64Val(m.Value(path...))
}

func (m M) Uint64(path ...string) uint64 {
	return converter.Uint64(m.Value(path...))
}

func (m M) Uint64Def(def uint64, path ...string) (v uint64) {
	v, err := m.Uint64Val(path...)
	if err != nil {
		v = def
	}
	return
}

func (m M) Uint64Val(path ...string) (uint64, error) {
	return converter.Uint64Val(m.Value(path...))
}

func (m M) Float64(path ...string) float64 {
	return converter.Float64(m.Value(path...))
}

func (m M) Float64Def(def float64, path ...string) (v float64) {
	v, err := m.Float64Val(path...)
	if err != nil {
		v = def
	}
	return
}

func (m M) Float64Val(path ...string) (float64, error) {
	return converter.Float64Val(m.Value(path...))
}

func (m M) Str(path ...string) string {
	return converter.String(m.Value(path...))
}

func (m M) StrDef(def string, path ...string) (v string) {
	v, err := m.StrVal(path...)
	if err != nil {
		v = def
	}
	return
}

func (m M) StrVal(path ...string) (string, error) {
	return converter.StringVal(m.Value(path...))
}

func (m M) Map(path ...string) M {
	v, _ := MapVal(m.Value(path...))
	return v
}

func (m M) MapDef(def M, path ...string) (v M) {
	v, err := m.MapVal(path...)
	if err != nil {
		v = def
	}
	return
}

func (m M) MapVal(path ...string) (M, error) {
	return MapVal(m.Value(path...))
}

func (m M) ArrayVal(path ...string) ([]interface{}, error) {
	return converter.ArrayVal(m.Value(path...))
}

func (m M) ArrayMapVal(path ...string) ([]M, error) {
	return ArrayMapVal(m.Value(path...))
}

func (m M) ArrayStrVal(path ...string) ([]string, error) {
	return converter.ArrayStringVal(m.Value(path...))
}

func (m M) ArrayIntVal(path ...string) ([]int, error) {
	return converter.ArrayIntVal(m.Value(path...))
}

func (m M) ArrayInt64Val(path ...string) ([]int64, error) {
	return converter.ArrayInt64Val(m.Value(path...))
}

func (m M) ArrayUint64Val(path ...string) ([]uint64, error) {
	return converter.ArrayUint64Val(m.Value(path...))
}

func (m M) ArrayFloat64Val(path ...string) ([]float64, error) {
	return converter.ArrayFloat64Val(m.Value(path...))
}

func (m M) Value(path ...string) interface{} {
	v, _ := m.ValueVal(path...)
	return v
}

func (m M) ValueVal(path ...string) (v interface{}, err error) {
	for _, p := range path {
		v, err = m.pathValue(p)
		if err == nil {
			break
		}
	}
	return
}

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

func (m M) Delete(path string) {
	m.setValue(path, nil)
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

//Get is implement for attrvalid.ValueGetter
func (m M) Get(path string) (interface{}, error) {
	return m.Value(strings.Split(path, "|")...), nil
}

//ValidFormat is implement for attrvalid.Validable
func (m M) ValidFormat(f string, args ...interface{}) error {
	return attrvalid.ValidAttrFormat(f, m, true, args...)
}

//Exist will return true if key having
func (m M) Exist(path ...string) bool {
	_, err := m.ValueVal(path...)
	return err == nil
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

// func Maps2Map(ms []Map, path ...string) Map {
// 	var res = Map{}
// 	for _, m := range ms {
// 		res[m.StrValP(path)] = m
// 	}
// 	return res
// }
