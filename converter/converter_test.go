package converter

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"
)

type XTIME time.Time

func TestNumber(t *testing.T) {
	type MyInt int
	var m map[string]interface{}
	//test all type
	m = map[string]interface{}{}
	m["abc"] = "123"
	m["abc2"] = int(1)
	m["myint"] = MyInt(1)
	m["float32"] = float32(1)
	m["float64"] = float64(1)
	m["int"] = int(1)
	m["int8"] = int8(1)
	m["int16"] = int16(1)
	m["int32"] = int32(1)
	m["int64"] = int64(1)
	m["uint"] = uint(1)
	m["uint8"] = uint8(1)
	m["uint16"] = uint16(1)
	m["uint32"] = uint32(1)
	m["uint64"] = uint64(1)
	m["time"] = time.Now()
	m["time1"] = time.Time{}
	m["time2"] = XTIME(time.Now())
	m["time3"] = XTIME(time.Time{})
	for key, val := range m {
		intA, err := IntVal(val)
		if err != nil {
			t.Error(key, err)
			return
		}
		intB := Int(val)
		if intA != intB {
			t.Errorf("%v  %v  %v", key, intA, intB)
			return
		}
		int64A, err := Int64Val(val)
		if err != nil {
			t.Error(key)
			return
		}
		int64B := Int64(val)
		if int64A != int64B {
			t.Errorf("%v  %v  %v", key, int64A, int64B)
			return
		}
		uint64A, err := Uint64Val(val)
		if err != nil {
			t.Error(key)
			return
		}
		uint64B := Uint64(val)
		if uint64A != uint64B {
			t.Errorf("%v  %v  %v", key, uint64A, uint64B)
			return
		}
		float64A, err := Float64Val(val)
		if err != nil {
			t.Error(key)
			return
		}
		float64B := Float64(val)
		if float64A != float64B {
			t.Errorf("%v  %v  %v", key, float64A, float64B)
			return
		}
	}
	//test error
	m = map[string]interface{}{}
	m["nil"] = nil
	m["abd"] = "a123"
	m["str1"] = []byte("akkkk")
	for key, val := range m {
		_, err := IntVal(val)
		if err == nil {
			t.Error(key)
			return
		}
		_, err = Int64Val(val)
		if err == nil {
			t.Error(key)
			return
		}
		_, err = Uint64Val(val)
		if err == nil {
			t.Error(key)
			return
		}
		_, err = Float64Val(val)
		if err == nil {
			t.Error(key)
			return
		}
	}
}

func TestString(t *testing.T) {
	var m map[string]interface{}
	//test all type
	m = map[string]interface{}{}
	type myString string
	var sbyte = []byte("akkkk")
	var sstr = "a123"
	m["abd"] = "a123"
	m["my_str"] = myString("jdldklal")
	m["str1"] = []byte("akkkk")
	m["str2"] = &sbyte
	m["str3"] = &sstr
	m["other"] = 111
	m["arr1"] = []int{1}
	m["arr1"] = []interface{}{nil, nil}
	for key, val := range m {
		strA, err := StringVal(val)
		if err != nil {
			t.Error(key)
			return
		}
		strB := String(val)
		if strA != strB {
			t.Errorf("%v  %v  %v", key, strA, strB)
			return
		}
	}
	//test error
	m = map[string]interface{}{}
	m["nil"] = nil
	for key, val := range m {
		_, err := StringVal(val)
		if err == nil {
			t.Error(key)
			return
		}
	}
}

func TestArray(t *testing.T) {
	var i0, i1 = 1, 2
	var s0 = "1,2,3"
	var m map[string]interface{}
	//test all type
	m = map[string]interface{}{}
	m["str1"] = []byte("123")
	m["str2"] = "1,2,3"
	m["str3"] = &s0
	m["arr1"] = []int{1, 1}
	m["arr2"] = []interface{}{1, 1}
	m["arr3"] = []*int{&i0, &i1}
	m["arr4"] = []interface{}{&i0, &i1}
	for key, val := range m {
		_, err := ArrayVal(val)
		if err != nil {
			t.Error(key, err)
			return
		}
		_, err = ArrayStringVal(val)
		if err != nil {
			t.Error(key, err)
			return
		}
		_, err = ArrayIntVal(val)
		if err != nil {
			t.Error(key, err)
			return
		}
		_, err = ArrayInt64Val(val)
		if err != nil {
			t.Error(key, err)
			return
		}
		_, err = ArrayUint64Val(val)
		if err != nil {
			t.Error(key, err)
			return
		}
		_, err = ArrayFloat64Val(val)
		if err != nil {
			t.Error(key, err)
			return
		}
	}
	//test error
	var snil *string
	var serr string = "xx"
	m = map[string]interface{}{}
	m["nil"] = nil
	m["int"] = 1
	m["str1"] = "xx"
	m["str2"] = snil
	m["str3"] = &serr
	m["i1"] = []interface{}{"aaa"}
	m["i2"] = []*testing.T{nil}
	// m["nil"] = []interface{}{"1", snil}
	for key, val := range m {
		_, err := ArrayIntVal(val)
		if err == nil {
			t.Error(key)
			return
		}
		_, err = ArrayInt64Val(val)
		if err == nil {
			t.Error(key)
			return
		}
		_, err = ArrayUint64Val(val)
		if err == nil {
			t.Error(key)
			return
		}
		_, err = ArrayFloat64Val(val)
		if err == nil {
			t.Error(key)
			return
		}
	}
	_, err := ArrayStringVal(nil)
	if err == nil {
		t.Error("xx")
		return
	}
	_, err = ArrayStringVal(snil)
	if err == nil {
		t.Error("xx")
		return
	}
	_, err = ArrayStringVal(t)
	if err == nil {
		t.Error("xx")
		return
	}
	_, err = ArrayStringVal([]interface{}{nil})
	if err == nil {
		t.Error("xx")
		return
	}
	_, err = ArrayStringVal([]*testing.T{nil})
	if err == nil {
		t.Error("xx")
		return
	}
	_, err = ArrayValAll(t, true)
	if err != nil {
		t.Error("xx")
		return
	}
	_, err = ArrayVal(nil)
	if err == nil {
		t.Error("xx")
		return
	}
	_, err = ArrayVal(snil)
	if err == nil {
		t.Error("xx")
		return
	}
	_, err = ArrayVal(t)
	if err == nil {
		t.Error("xx")
		return
	}
}

func TestArrayHaving(t *testing.T) {
	iary := []int{1, 2, 3, 4, 5, 6}
	if !ArrayHaving(iary, 2) {
		t.Error("value exis in array.")
		return
	}
	if ArrayHaving(iary, 8) {
		t.Error("value not exis in array.")
		return
	}
	//
	fary := []float32{1.0, 2.0, 3.0, 4.0, 5.0}
	if !ArrayHaving(fary, float32(1.0)) {
		t.Error("value exis in array.")
		return
	}
	if ArrayHaving(fary, float32(8.0)) {
		t.Error("value not exis in array.")
		return
	}
	//
	sary := []string{"a", "b", "c", "d", "e", "f"}
	if !ArrayHaving(sary, "c") {
		t.Error("value exis in array.")
		return
	}
	if ArrayHaving(sary, "g") {
		t.Error("value not exis in array.")
		return
	}
	ab := ""
	if ArrayHaving(ab, 8) {
		t.Error("value exis in array.")
		return
	}
}

func TestJSON(t *testing.T) {
	v1 := JSON("v")
	if v1 != "\"v\"" {
		t.Error(v1)
		return
	}
	v2 := JSON(TestJSON)
	if !strings.Contains(v2, "unsupported") {
		t.Error(v2)
		return
	}
}

func TestXML(t *testing.T) {
	v1 := XML("v")
	if v1 != "<string>v</string>" {
		t.Error(v1)
		return
	}
	v2 := XML(TestJSON)
	if !strings.Contains(v2, "unsupported") {
		t.Error(v2)
		return
	}
}

type xmlObj struct {
}

func TestUnmarshal(t *testing.T) {
	var err error
	_, err = UnmarshalJSON(bytes.NewBufferString("{}"), &map[string]interface{}{})
	if err != nil {
		t.Error(err)
		return
	}
	_, err = UnmarshalXML(bytes.NewBufferString("<xml></xml>"), &xmlObj{})
	if err != nil {
		t.Error(err)
		return
	}
}

func TestPtr(t *testing.T) {
	Int8Ptr(1)
	Uint8Ptr(1)
	Int16Ptr(1)
	Uint16Ptr(1)
	Int32Ptr(1)
	Uint32Ptr(1)
	IntPtr(1)
	UintPtr(1)
	Int64Ptr(1)
	Uint64Ptr(1)
	Float32Ptr(1)
	Float64Ptr(1)
	StringPtr("ss")
}

func TestJoin(t *testing.T) {
	if Join([]int{}, ",") != "" {
		t.Error("error")
		return
	}
	if Join([]int{1, 2, 3}, ",") != "1,2,3" {
		t.Error("error")
		return
	}
	if Join([]*int{IntPtr(1), IntPtr(2), IntPtr(3)}, ",") != "1,2,3" {
		t.Error("error")
		return
	}
	if Join([]*string{StringPtr("1"), StringPtr("2"), StringPtr("3")}, ",") != "1,2,3" {
		t.Error("error")
		return
	}
	func() {
		defer func() {
			recover()
		}()
		Join("xx", ",")
	}()
}

func TestJoinSafe(t *testing.T) {
	if JoinSafe([]int{}, ",", JoinPolicyDefault) != "" {
		t.Error("error")
		return
	}
	if JoinSafe([]int{1, 2, 3}, ",", JoinPolicyDefault) != "1,2,3" {
		t.Error("error")
		return
	}
	if JoinSafe([]*int{IntPtr(1), IntPtr(2), IntPtr(3)}, ",", JoinPolicyDefault) != "1,2,3" {
		t.Error("error")
		return
	}
	if JoinSafe([]*int{IntPtr(1), IntPtr(2), nil, IntPtr(3)}, ",", JoinPolicyDefault) != "1,2,3" {
		t.Error("error")
		return
	}
	if JoinSafe([]*int{IntPtr(1), IntPtr(2), nil, IntPtr(3)}, ",", JoinPolicyDefault) != "1,2,3" {
		t.Error("error")
		return
	}
	if JoinSafe([]*int{IntPtr(1), IntPtr(2), nil, IntPtr(3)}, ",", JoinPolicyNilString) != "1,2,<nil>,3" {
		t.Error(JoinSafe([]*int{IntPtr(1), IntPtr(2), nil, IntPtr(3)}, ",", JoinPolicyNilString))
		return
	}
	if JoinSafe([]*string{StringPtr("1"), StringPtr("2"), StringPtr("3")}, ",", JoinPolicyDefault) != "1,2,3" {
		t.Error("error")
		return
	}
	if JoinSafe([]*string{StringPtr("1"), StringPtr("2"), nil, StringPtr("3")}, ",", JoinPolicyDefault) != "1,2,3" {
		t.Error("error")
		return
	}
	if JoinSafe([]*string{StringPtr("1"), StringPtr("2"), nil, StringPtr("3")}, ",", JoinPolicyNilString) != "1,2,<nil>,3" {
		t.Error("error")
		return
	}
	if JoinSafe("xx", ",", JoinPolicyDefault) != "" {
		t.Error("error")
		return
	}
	if JoinSafe("xx", ",", JoinPolicyNotSliceString) != "xx" {
		t.Error("error")
		return
	}
}

func TestIndirectString(t *testing.T) {
	var x *int
	fmt.Println(IndirectString(x))
	fmt.Println(IndirectString(nil))
	fmt.Println(IndirectString(1))
}
