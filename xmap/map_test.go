package xmap

import (
	"fmt"
	"sort"
	"testing"
)

func TestMap(t *testing.T) {
	m := New()
	testMap(t, m)
	m2 := NewSafe()
	testMap(t, m2)
	m3 := M{}
	testMap(t, m3)
}

func testMap(t *testing.T, m Valuable) {
	m.SetValue("abc", "123")
	m.SetValue("int", 123)
	m.SetValue("x", M{"a": 123})
	m.SetValue("/ary", []interface{}{1, 2, 3})
	m.SetValue("/ary2", "1,2,3")
	m.SetValue("/arm", []interface{}{M{}, M{}, M{}})
	assert := func(v bool) {
		if !v {
			panic("error")
		}
	}
	assert(m.Exist("abc"))
	//
	assert(nil != m.Raw())
	assert(nil != m.Map("x"))
	assert(nil != m.MapDef(nil, "x"))
	assert(nil == m.MapDef(nil, "not"))
	//
	assert(m.Str("abc") == "123")
	assert(m.StrDef("", "abc") == "123")
	assert(m.StrDef("", "not") == "")
	//
	assert(m.Int("int") == 123)
	assert(m.Int64("int") == 123)
	assert(m.Uint64("int") == 123)
	assert(m.Float64("int") == 123)
	assert(m.Int("x/a") == 123)
	assert(m.Int64("x/a") == 123)
	assert(m.Uint64("x/a") == 123)
	assert(m.Float64("x/a") == 123)
	//
	assert(m.IntDef(0, "int") == 123)
	assert(m.Int64Def(0, "int") == 123)
	assert(m.Uint64Def(0, "int") == 123)
	assert(m.Float64Def(0, "int") == 123)
	//
	assert(m.IntDef(0, "not") == 0)
	assert(m.Int64Def(0, "not") == 0)
	assert(m.Uint64Def(0, "not") == 0)
	assert(m.Float64Def(0, "not") == 0)
	//
	assert(len(m.ArrayDef(nil, "/ary")) == 3)
	assert(len(m.ArrayMapDef(nil, "/arm")) == 3)
	assert(len(m.ArrayIntDef(nil, "/ary")) == 3)
	assert(len(m.ArrayInt64Def(nil, "/ary")) == 3)
	assert(len(m.ArrayUint64Def(nil, "/ary")) == 3)
	assert(len(m.ArrayFloat64Def(nil, "/ary")) == 3)
	assert(len(m.ArrayStrDef(nil, "/ary")) == 3)
	//
	assert(len(m.ArrayDef(nil, "/not")) == 0)
	assert(len(m.ArrayMapDef(nil, "/not")) == 0)
	assert(len(m.ArrayIntDef(nil, "/not")) == 0)
	assert(len(m.ArrayInt64Def(nil, "/not")) == 0)
	assert(len(m.ArrayUint64Def(nil, "/v")) == 0)
	assert(len(m.ArrayFloat64Def(nil, "/not")) == 0)
	assert(len(m.ArrayStrDef(nil, "/not")) == 0)
	//
	//
	if v, err := m.StrVal("int"); true {
		assert(v == "123" && err == nil)
	}
	if v, err := m.IntVal("int"); true {
		assert(v == 123 && err == nil)
	}
	if v, err := m.Int64Val("int"); true {
		assert(v == 123 && err == nil)
	}
	if v, err := m.Uint64Val("int"); true {
		assert(v == 123 && err == nil)
	}
	if v, err := m.Float64Val("int"); true {
		assert(v == 123 && err == nil)
	}
	if v, err := m.MapVal("x"); true {
		assert(v != nil && err == nil)
	}
	if v, err := m.ValueVal("x"); true {
		assert(v != nil && err == nil)
	}
	//
	if v, err := m.ArrayVal("ary"); true {
		assert(len(v) == 3 && err == nil)
	}
	if v, err := m.ArrayStrVal("ary"); true {
		assert(len(v) == 3 && err == nil)
	}
	if v, err := m.ArrayIntVal("ary"); true {
		assert(len(v) == 3 && err == nil)
	}
	if v, err := m.ArrayInt64Val("ary"); true {
		assert(len(v) == 3 && err == nil)
	}
	if v, err := m.ArrayUint64Val("ary"); true {
		assert(len(v) == 3 && err == nil)
	}
	if v, err := m.ArrayFloat64Val("ary"); true {
		assert(len(v) == 3 && err == nil)
	}
	if v, err := m.ArrayMapVal("arm"); true {
		assert(len(v) == 3 && err == nil)
	}
	//
	//test remove
	m.SetValue("having", "123")
	assert(m.Value("having") != nil)
	m.Delete("having")
	assert(m.Value("having") == nil)
	//
	m.SetValue("having", "123")
	assert(m.Length() > 0)
	m.Clear()
	assert(m.Length() == 0)
}

func TestArrayMap(t *testing.T) {
	var m map[string]interface{}
	//test all type
	m = map[string]interface{}{}
	m["arr1"] = []map[string]interface{}{{}, {}}
	m["arr2"] = []interface{}{M{}, M{}}
	m["nil"] = nil
	for key, val := range m {
		_, err := ArrayMapVal(val)
		if err != nil {
			t.Error(key)
			return
		}
	}
	//test error
	m = map[string]interface{}{}
	m["int"] = 1
	m["i1"] = []interface{}{"aaa"}
	m["i2"] = []*testing.T{nil}
	for key, val := range m {
		_, err := ArrayMapVal(val)
		if err == nil {
			t.Error(key)
			return
		}
	}
}

func TestPathValue(t *testing.T) {
	//data
	m1 := map[string]interface{}{
		"s":   "str",
		"i":   int64(16),
		"f":   float64(16),
		"ary": []interface{}{1, 3, 4},
	}
	m2 := map[string]interface{}{
		"a":   "abc",
		"m":   m1,
		"ary": []interface{}{"1", "3", "4"},
	}
	m3 := map[string]interface{}{
		"b":   "abcc",
		"m":   m2,
		"ary": []interface{}{m1, m2},
	}
	m4 := Wrap(M{
		"test": 1,
		"ms":   []interface{}{m1, m2, m3},
		"m3":   m3,
		"m4":   "{}",
		"m5":   []byte("{}"),
		"m6":   "[{}]",
		"m7":   []byte("[{}]"),
		"ary2": []int{1, 3, 4},
		"me":   map[string]string{"a": "b"},
	})
	var v interface{}
	var err error
	v, err = m4.ValueVal("/path")
	assertError(t, v, err)
	v, err = m4.ValueVal("/test")
	assertNotError(t, v, err)
	v, err = m4.ValueVal("/ms")
	assertNotError(t, v, err)
	v, err = m4.ValueVal("/m3")
	assertNotError(t, v, err)
	v, err = m4.MapVal("/m4")
	assertNotError(t, v, err)
	v, err = m4.MapVal("/m5")
	assertNotError(t, v, err)
	v, err = m4.ArrayMapVal("/m6")
	assertNotError(t, v, err)
	v, err = m4.ArrayMapVal("/m7")
	assertNotError(t, v, err)
	//
	v, err = m4.ValueVal("/m3/b")
	assertNotError(t, v, err)
	v, err = m4.ValueVal("/m3/b2")
	assertError(t, v, err)
	v, err = m4.ValueVal("/m3/ary")
	assertNotError(t, v, err)
	v, err = m4.ValueVal("/ms/1")
	assertNotError(t, v, err)
	v, err = m4.ValueVal("/ms/100")
	assertError(t, v, err)
	v, err = m4.ValueVal("/ms/a")
	assertError(t, v, err)
	v, err = m4.ValueVal("/ary2/100")
	assertError(t, v, err)
	v, err = m4.ValueVal("/ms/@len")
	assertNotError(t, v, err)
	v, err = m4.ValueVal("/ary2/@len")
	assertError(t, v, err)
	v, err = m4.ValueVal("/test/abc")
	assertError(t, v, err)
	v, err = m4.ValueVal("/me/a")
	assertError(t, v, err)
	v, err = m4.ValueVal("/mekkkk/a")
	assertError(t, v, err)
}

func TestSetValue(t *testing.T) {
	var v interface{}
	var err error
	m := Wrap(M{
		"eary":  []string{},
		"ary":   []interface{}{456},
		"emap":  map[string]string{},
		"map":   map[string]interface{}{},
		"ntype": "kkkk",
	})
	m.SetValue("/abc", M{"a": 1})
	v, err = m.ValueVal("/abc/a")
	assertNotError(t, v, err)
	err = m.SetValue("/ary/0", 123)
	assertNotError(t, nil, err)

	err = m.SetValue("/map/a", 123)
	assertNotError(t, nil, err)
	_, err = m.ValueVal("/map/a")
	assertNotError(t, nil, err)
	//
	//error
	err = m.SetValue("/abcd/abc", 123)
	assertError(t, nil, err)
	err = m.SetValue("/eary/1", 123)
	assertError(t, nil, err)
	err = m.SetValue("/ary/5", 123)
	assertError(t, nil, err)
	err = m.SetValue("/ary/a", 123)
	assertError(t, nil, err)
	err = m.SetValue("/emap/a", 123)
	assertError(t, nil, err)
	err = m.SetValue("/ntype/a", 123)
	assertError(t, nil, err)
	err = m.SetValue("", 123)
	assertError(t, nil, err)
	//
	mv, err := m.MapVal("/abc")
	assertNotError(t, mv, err)
	v, err = mv.ValueVal("/a")
	assertNotError(t, v, err)
	//
	b := &M{}
	err = b.setPathValue("", 123)
	assertError(t, nil, err)
}

func assertNotError(t *testing.T, v interface{}, err error) {
	if err != nil {
		panic(err)
	}
	fmt.Println(v)
}

func assertError(t *testing.T, v interface{}, err error) {
	fmt.Println(err)
	if err == nil {
		panic("not error")
	}
}

// func TestArray2(t *testing.T) {
// 	fmt.Println([]int{1, 3, 5}[:3])
// }

// func TestNewMap(t *testing.T) {
// 	fmt.Println(NewMap("map.json"))
// 	fmt.Println(NewMap("map.jsn"))
// 	fmt.Println(NewMaps("maps.json"))
// 	fmt.Println(NewMaps("maps.jsn"))
// }

func TestValidFormat(t *testing.T) {
	m := M(map[string]interface{}{
		"ab1": 1,
		"ab2": "xxx",
		"map": map[string]interface{}{
			"x1": 100,
		},
	})
	var v1 int64
	var v2 string
	var v3 int
	var v4 int
	err := m.ValidFormat(`
		ab1,R|I,R:0;
		ab2,R|S,L:0;
		/map/x1,R|I,R:0;
		not|ab1,R|I,R:0;
		`, &v1, &v2, &v3, &v4)
	if v1 != 1 || v2 != "xxx" || v3 != 100 || v4 != 1 {
		t.Error("error")
		return
	}
	fmt.Println(v1, v2, v3)
	if err != nil {
		t.Error(err)
		return
	}
	i := &impl{BaseValuable: m}
	err = i.ValidFormat(`
		ab1,R|I,R:0;
		ab2,R|S,L:0;
		/map/x1,R|I,R:0;
		not|ab1,R|I,R:0;
	`, &v1, &v2, &v3, &v4)
	if err != nil || v1 != 1 || v2 != "xxx" || v3 != 100 || v4 != 1 {
		t.Error("error")
		return
	}
	i.Raw()
}

func TestSafeValidFormat(t *testing.T) {
	m := WrapSafe(M(map[string]interface{}{
		"ab1": 1,
		"ab2": "xxx",
		"map": map[string]interface{}{
			"x1": 100,
		},
	}))
	var v1 int64
	var v2 string
	var v3 int
	var v4 int
	err := m.ValidFormat(`
		ab1,R|I,R:0;
		ab2,R|S,L:0;
		/map/x1,R|I,R:0;
		not|ab1,R|I,R:0;
		`, &v1, &v2, &v3, &v4)
	if v1 != 1 || v2 != "xxx" || v3 != 100 || v4 != 1 {
		t.Error("error")
		return
	}
	fmt.Println(v1, v2, v3)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestMSorter(t *testing.T) {
	newall := func() []M {
		v := WrapArray([]M{
			{
				"s": "a",
				"i": 1,
				"f": 1.0,
			},
			{
				"s": "c",
				"i": 3,
				"f": 3.0,
			},
			{
				"s": "b",
				"i": 2,
				"f": 2.0,
			},
		})
		return v
	}
	sort.Sort(NewMSorter(newall(), 0, false, "i"))
	sort.Sort(NewMSorter(newall(), 0, true, "i"))
	sort.Sort(NewMSorter(newall(), 1, false, "f"))
	sort.Sort(NewMSorter(newall(), 1, true, "f"))
	sort.Sort(NewMSorter(newall(), 2, false, "s"))
	sort.Sort(NewMSorter(newall(), 2, true, "s"))
}

func TestValuableSorter(t *testing.T) {
	newall := func() []Valuable {
		v := WrapSafeArray([]M{
			{
				"s": "a",
				"i": 1,
				"f": 1.0,
			},
			{
				"s": "c",
				"i": 3,
				"f": 3.0,
			},
			{
				"s": "b",
				"i": 2,
				"f": 2.0,
			},
		})
		return v
	}
	sort.Sort(NewValuableSorter(newall(), 0, false, "i"))
	sort.Sort(NewValuableSorter(newall(), 0, true, "i"))
	sort.Sort(NewValuableSorter(newall(), 1, false, "f"))
	sort.Sort(NewValuableSorter(newall(), 1, true, "f"))
	sort.Sort(NewValuableSorter(newall(), 2, false, "s"))
	sort.Sort(NewValuableSorter(newall(), 2, true, "s"))
}

func TestReplaceAll(t *testing.T) {
	val := M{
		"x1": 100,
		"x2": M{
			"a": 100,
		},
	}
	if val.ReplaceAll("${x1}", true, true) != "100" {
		t.Error("error")
		return
	}
	val2, _ := Parse(baseMap{"a": 100})
	if val2.ReplaceAll("${a}", true, true) != "100" {
		t.Error("error")
		return
	}
	if val2.ReplaceAll("${b}", true, true) != "" {
		t.Error("error")
		return
	}
	val3, _ := ParseSafe(`{"a":100}`)
	if val3.ReplaceAll("${a}", true, true) != "100" {
		t.Error("error")
		return
	}
	if val3.ReplaceAll("${b}", true, true) != "" {
		t.Error("error")
		return
	}
}

func TestEscapeValue(t *testing.T) {
	val := M{
		"x1/1": 100,
		"x1/2": M{
			"a": 100,
		},
	}
	if val.IntDef(0, "/x1\\/1") != 100 {
		t.Error("error")
		return
	}
	if val.IntDef(0, "/x1\\/2/a") != 100 {
		t.Error("error")
		return
	}
}
