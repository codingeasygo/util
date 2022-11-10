package attrvalid

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/codingeasygo/util/converter"
)

type TestInt int
type TestIntArray []int
type TestIntPtrArray []*int

func TestCompatibleType(t *testing.T) {
	var iVal int
	var iArray []int
	var iPtrArray []*int
	var i interface{}
	if !CompatibleType(reflect.TypeOf(iVal)) {
		t.Error("err")
		return
	}
	if !CompatibleType(reflect.TypeOf(iArray)) {
		t.Error("err")
		return
	}
	if !CompatibleType(reflect.TypeOf(iPtrArray)) {
		t.Error("err")
		return
	}

	if !CompatibleType(reflect.TypeOf(TestInt(iVal))) {
		t.Error("err")
		return
	}
	if !CompatibleType(reflect.TypeOf(TestIntArray(iArray))) {
		t.Error("err")
		return
	}
	if !CompatibleType(reflect.TypeOf(TestIntPtrArray(iPtrArray))) {
		t.Error("err")
		return
	}
	if CompatibleType(reflect.TypeOf(i)) {
		t.Error("err")
		return
	}
	if CompatibleType(reflect.TypeOf(nil)) {
		t.Error("err")
		return
	}
}

func TestM(t *testing.T) {
	var a int
	m := M(map[string]interface{}{
		"a": 1,
	})
	err := m.ValidFormat(`a,r|i,r:0`, &a)
	if err != nil || a != 1 {
		t.Error(err)
		return
	}

}

func TestMS(t *testing.T) {
	var a int
	m := MS(map[string]string{
		"a": "1",
	})
	err := m.ValidFormat(`a,r|i,r:0`, &a)
	if err != nil || a != 1 {
		t.Error(err)
		return
	}

}

func TestValidAttrTemple(t *testing.T) {
	v, err := ValidAttrTemple("测试", "r|s", "l:~10", true, nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	_, err = ValidAttrTemple("测试测试测试测试", "r|s", "l:~10", true, nil)
	if err == nil {
		t.Error("not error")
		return
	}
	//
	v, err = ValidAttrTemple("男", "r|s", "o:男~女", true, nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	_, err = ValidAttrTemple("男ks", "r|s", "o:男~女", true, nil)
	if err == nil {
		t.Error("not error")
		return
	}
	//
	v, err = ValidAttrTemple("test@gmail.com", "r|s", "p:^.*\\@.*$", true, nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	_, err = ValidAttrTemple("ks", "r|s", "p:^.*\\@.*$", true, nil)
	if err == nil {
		t.Error("not error")
		return
	}
	//
	v, err = ValidAttrTemple("8", "o|i", "r:5~10", true, nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	v, err = ValidAttrTemple("8", "o|i", "r:5~", true, nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	_, err = ValidAttrTemple("12", "o|i", "r:5~10", true, nil)
	if err == nil {
		t.Error("not error")
		return
	}
	//
	v, err = ValidAttrTemple("8", "o|f", "r:5~10", true, nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	v, err = ValidAttrTemple("8", "o|f", "r:5~", true, nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	_, err = ValidAttrTemple("12", "o|f", "r:5~10", true, nil)
	if err == nil {
		t.Error("not error")
		return
	}
	//
	v, err = ValidAttrTemple("测", "o|s", "l:~8", true, nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	_, err = ValidAttrTemple("测", "o|s", "l:2~", true, nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	_, err = ValidAttrTemple("测度测度测度测度测度", "o|s", "l:~8", true, nil)
	if err == nil {
		t.Error("not error")
		return
	}
	//
	v, err = ValidAttrTemple("测", "o|s", "l:2~8", true, nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	_, err = ValidAttrTemple("a", "o|s", "l:2~8", true, nil)
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = ValidAttrTemple("a", "o|s", "n:", true, nil)
	if err != nil {
		t.Error("not error")
		return
	}
	//
	v, err = ValidAttrTemple("test@gmail.com", "o|s", "p:^.*\\@.*$", true, nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	_, err = ValidAttrTemple("ks", "o|s", "p:^.*\\@.*$", true, nil)
	if err == nil {
		t.Error("not error")
		return
	}
	//
	v, err = ValidAttrTemple("1", "o|i", "o:1~2~3~4~5", true, nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	_, err = ValidAttrTemple("11", "o|i", "o:1~2~3~4~5", true, nil)
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = ValidAttrTemple("11", "o|i", "n:", true, nil)
	if err != nil {
		t.Error("not error")
		return
	}
	//
	v, err = ValidAttrTemple("1.1", "o|f", "o:1.1~2.2~3.3~4~5", true, nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	_, err = ValidAttrTemple("11", "o|f", "o:1~2~3~4~5", true, nil)
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = ValidAttrTemple("11", "o|f", "n:", true, nil)
	if err != nil {
		t.Error("not error")
		return
	}
	//
	_, err = ValidAttrTemple("测", "o|s", "l:a", true, nil)
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = ValidAttrTemple("测", "o|s", "KK:a", true, nil)
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = ValidAttrTemple("test@gmail.com", "o|s", "p:*,..", true, nil)
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = ValidAttrTemple("测", "o|i", "r:8~9", true, nil)
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = ValidAttrTemple("测", "o|f", "r:8~9", true, nil)
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = ValidAttrTemple("测", "o|f", "o:8~9", true, nil)
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = ValidAttrTemple("测", "o|n", "r:8~9", true, nil)
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = ValidAttrTemple("5", "r|i", "r:~1", true, nil)
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = ValidAttrTemple("5", "r|i", "r:a~10", true, nil)
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = ValidAttrTemple("5", "r|i", "r:1~a", true, nil)
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = ValidAttrTemple("5", "r|f", "r:~1", true, nil)
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = ValidAttrTemple("5", "r|f", "r:a~10", true, nil)
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = ValidAttrTemple("5", "r|f", "r:1~a", true, nil)
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = ValidAttrTemple("5", "r|i", "m:1~a", true, nil)
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = ValidAttrTemple("5", "r|i", "o:1~a", true, nil)
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = ValidAttrTemple("5", "r|f", "o:1~a", true, nil)
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = ValidAttrTemple("5", "r|f", "m:1~k", true, nil)
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = ValidAttrTemple("5", "r|i", "o", true, nil)
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = ValidAttrTemple("5", "r", "o:1~10", true, nil)
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = ValidAttrTemple("", "r|i", "o:1~10", true, nil)
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = ValidAttrTemple("", "o|i", "o:1~10", true, nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	_, err = ValidAttrTemple("a", "o|s", "l:a~8", true, nil)
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = ValidAttrTemple("a", "o|s", "l:2~a", true, nil)
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = ValidAttrTemple(nil, "o|s", "l:0", true, nil)
	if err != nil {
		t.Error("not error")
		return
	}
	_, err = ValidAttrTemple(nil, "r|s", "l:0", true, nil)
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = ValidAttrTemple(0, "o|i", "r:0", true, nil)
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	_, err = ValidAttrTemple(0, "o|f", "r:0", true, nil)
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	_, err = ValidAttrTemple("", "o|s", "r:0", true, nil)
	if err != nil {
		t.Errorf("%v", err)
		return
	}
}

type EnumIntTest int

func (e *EnumIntTest) EnumValid(v interface{}) (err error) {
	val, err := converter.IntVal(v)
	if err == nil && val != 1 && val != 2 {
		err = fmt.Errorf("only supported 1,2")
	}
	return
}

type EnumIntArrayTest []int

func (e *EnumIntArrayTest) EnumValid(v interface{}) (err error) {
	val, err := converter.IntVal(v)
	if err == nil && val != 1 && val != 2 {
		err = fmt.Errorf("only supported 1,2")
	}
	return
}

type EnumFloatTest float64

func (e *EnumFloatTest) EnumValid(v interface{}) (err error) {
	val, err := converter.Float64Val(v)
	if err == nil && val != 1 && val != 2 {
		err = fmt.Errorf("only supported 1,2")
	}
	return
}

type EnumFloatArrayTest []int

func (e *EnumFloatArrayTest) EnumValid(v interface{}) (err error) {
	val, err := converter.Float64Val(v)
	if err == nil && val != 1 && val != 2 {
		err = fmt.Errorf("only supported 1,2")
	}
	return
}

type EnumStringTest float64

func (e *EnumStringTest) EnumValid(v interface{}) (err error) {
	val, err := converter.StringVal(v)
	if err == nil && val != "1" && val != "2" {
		err = fmt.Errorf("only supported 1,2")
	}
	return
}

type EnumStringArrayTest []int

func (e *EnumStringArrayTest) EnumValid(v interface{}) (err error) {
	val, err := converter.StringVal(v)
	if err == nil && val != "1" && val != "2" {
		err = fmt.Errorf("only supported 1,2")
	}
	return
}

type DefineInt64Test int64

func TestValidAttrFormat(t *testing.T) {
	mv := map[string]interface{}{}
	mv["a"] = "abc"
	mv["i"] = "10"
	mv["i4"] = "1"
	mv["f"] = "10.3"
	mv["ef"] = "20.3"
	mv["len"] = "11111111"
	mv["ary"] = "1,2,3,4,5"
	mv["ary2"] = "1,2,3,,4,5"
	mv["ary3"] = []interface{}{1, 2, 3, 4, 5}
	mv["ary4"] = "1,2"
	var a string
	var i int64
	var k string
	var ks []string
	var f float64
	var iv1 int
	var iv1ary []int
	var iv2 int16
	var iv3 int32
	var iv4 int64
	var iv5 uint
	var iv6 uint16
	var iv7 uint32
	var iv8 uint64
	var iv9 float32
	var iv10 float64
	var iv10ary []float64
	var iv11 string
	var iv12 int64
	var snot string
	var iv1ary2 []int
	var svary []string
	var iv10ary2 []float64
	var iv10ary3 []float64
	var eint EnumIntTest
	err := ValidAttrFormat(`//abc
		a,r|s,l:~5;//abc
		i,r|i,r:1~20;
		i,o|i,r:1~20;//sfdsj
		i,o|i,r:1~20;//sfdsj
		f,r|f,r:1.5~20;
		i,r|i,r:0;
		i,r|i,r:0;
		i,r|i,r:0;
		i,r|i,r:0;
		i,r|i,r:0;
		i,r|i,r:0;
		i,r|i,r:0;
		i,r|i,r:0;
		i,r|i,r:0;
		i,r|i,r:0;
		i,r|i,r:0;
		i,r|i,r:0;
		i,r|i,r:0;
		i,r|s,l:0;
		not,o|s,l:0;
		ary,r|s,l:0;
		ary,r|i,r:0;
		ary,r|f,r:0;
		ary3,r|f,r:0;
		i,r|i,r:0;
		`, M(mv), true, &a, &i, &k, &ks, &f,
		&iv1, &iv1ary, &iv2, &iv3, &iv4, &iv5,
		&iv6, &iv7, &iv8, &iv9, &iv10, &iv10ary,
		&iv11, &iv12, &snot, &svary,
		&iv1ary2, &iv10ary2, &iv10ary3,
		&eint,
	)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(k, ks, len(iv1ary2), len(iv10ary2))
	if k != "10" || ks[0] != "10" || iv1 != 10 || iv1ary2[0] != 1 || iv10 != 10 || iv10ary2[0] != 1 {
		t.Error("error")
		return
	}
	fmt.Println(len(svary), len(iv1ary2), len(iv10ary2))
	if len(svary) != 5 || len(iv1ary2) != 5 || len(iv10ary2) != 5 || len(iv10ary3) != 5 {
		t.Error("error")
		return
	}
	fmt.Println(a, i, k, f)
	//string to int on ValidValue
	err = ValidAttrFormat(`
		i,r|s,l:0;
		i,r|s,l:0;
		i,r|s,l:0;
		i,r|s,l:0;
		i,r|s,l:0;
		i,r|s,l:0;
		i,r|s,l:0;
		i,r|s,l:0;
		i,r|s,l:0;
		i,r|s,l:0;
		i,r|s,l:0;
		i,r|s,l:0;
		i,r|s,l:0;
		i,r|s,l:0;
		`, M(mv), true,
		&iv1, &iv1ary, &iv2, &iv3, &iv4, &iv5,
		&iv6, &iv7, &iv8, &iv9, &iv10, &iv10ary,
		&iv11, &iv12,
	)
	if err != nil {
		t.Error(err.Error())
		return
	}
	//
	//test array
	svary, iv1ary2, iv10ary2 = nil, nil, nil
	err = ValidAttrFormat(`
		ary2,r|s,l:0;
		ary2,r|i,r:0;
		ary2,r|f,r:0;
		`, M(mv), true, &svary, &iv1ary2, &iv10ary2)
	if err == nil {
		t.Error("error")
		return
	}
	svary, iv1ary2, iv10ary2 = nil, nil, nil
	err = ValidAttrFormat(`
		ary2,o|s,l:0;
		ary2,o|i,r:0;
		ary2,o|f,r:0;
		`, M(mv), true, &svary, &iv1ary2, &iv10ary2)
	if err != nil {
		t.Error("error")
		return
	}
	if len(svary) != 5 || len(iv1ary2) != 5 || len(iv10ary2) != 5 {
		t.Error("error")
		return
	}
	//
	//test enum
	var enumInt EnumIntTest
	var enumIntArray EnumIntArrayTest
	var enumFloat EnumFloatTest
	var enumFloatArray EnumFloatArrayTest
	var enumString EnumFloatTest
	var enumStringArray EnumFloatArrayTest
	err = ValidAttrFormat(`
			i,r|i,e:;
		`, M(mv), true, &enumInt)
	if err == nil {
		t.Error("error")
		return
	}
	err = ValidAttrFormat(`
			ary,r|i,e:0;
		`, M(mv), true, &enumIntArray)
	if err == nil {
		t.Error("error")
		return
	}
	err = ValidAttrFormat(`
		i,r|i,e:;
	`, M(mv), true, &i)
	if err == nil {
		t.Error("error")
		return
	}
	err = ValidAttrFormat(`
		i,r|f,e:;
	`, M(mv), true, &f)
	if err == nil {
		t.Error("error")
		return
	}
	err = ValidAttrFormat(`
		i,r|s,e:;
	`, M(mv), true, &a)
	if err == nil {
		t.Error("error")
		return
	}
	err = ValidAttrFormat(`
		i4,r|i,e:0;
		ary4,r|i,e:0;
		i4,r|f,e:0;
		ary4,r|f,e:0;
		i4,r|s,e:0;
		ary4,r|s,e:0;
	`, M(mv), true, &enumInt, &enumIntArray, &enumFloat, &enumFloatArray, &enumString, &enumStringArray)
	if err != nil {
		t.Error(err)
		return
	}
	if enumInt != 1 || len(enumIntArray) != 2 || enumIntArray[0] != 1 || enumIntArray[1] != 2 || enumFloat != 1 || len(enumFloatArray) != 2 || enumString != 1 || len(enumStringArray) != 2 {
		t.Error("error")
		return
	}
	//
	err = ValidAttrFormat(`
		a,r|s l:~5;
		`, M(mv), true, &a)
	if err == nil {
		t.Error("not error")
		return
	}
	//
	err = ValidAttrFormat(`
		len,r|s,l:~5;
		`, M(mv), true, &a)
	if err == nil {
		t.Error("not error")
		return
	}
	//
	var ea float32
	err = ValidAttrFormat(`
		a,r|s,l:~5;
		`, M(mv), true, &ea)
	if err == nil {
		t.Error("not error")
		return
	}
	fmt.Println(err.Error())
	//
	err = ValidAttrFormat(``, M(mv), true, &a)
	if err == nil {
		t.Error("not error")
		return
	}
	fmt.Println(err.Error())
	//
	err = ValidAttrFormat(`
		len,r|s,l:~5;
		len,r|s,l:~5;
		`, M(mv), true, &a)
	if err == nil {
		t.Error("not error")
		return
	}
	fmt.Println(err.Error())
	err = ValidAttrFormat(`
		len,r|s,l:~5,this is error message;
		`, M(mv), true, &a)
	if err == nil {
		t.Error("not error")
		return
	}
	fmt.Println(err.Error())
}

func TestValidAttrFormatPointer(t *testing.T) {
	mv := map[string]interface{}{}
	mv["a"] = "abc"
	mv["i"] = "10"
	mv["f"] = "10.3"
	mv["ef"] = "20.3"
	mv["len"] = "11111111"
	mv["ary"] = "1,2,3,4,5"
	mv["ary2"] = "1,2,3,,4,5"
	var a *string
	var i *int64
	var k *string
	var ks []*string
	var f *float64
	var iv1 *int
	var iv1ary []*int
	var iv2 *int16
	var iv3 *int32
	var iv4 *int64
	var iv5 *uint
	var iv6 *uint16
	var iv7 *uint32
	var iv8 *uint64
	var iv9 *float32
	var iv10 *float64
	var iv10ary []*float64
	var iv11 *string
	var iv12 *int64
	var iv1ary2 []*int
	var arystr []*string
	var iv10ary2 []*float64
	err := ValidAttrFormat(`//abc
		a,r|s,l:~5;//abc
		i,r|i,r:1~20;
		i,o|i,r:1~20;//sfdsj
		i,o|i,r:1~20;//sfdsj
		f,r|f,r:1.5~20;
		i,r|i,r:0;
		i,r|i,r:0;
		i,r|i,r:0;
		i,r|i,r:0;
		i,r|i,r:0;
		i,r|i,r:0;
		i,r|i,r:0;
		i,r|i,r:0;
		i,r|i,r:0;
		i,r|i,r:0;
		i,r|i,r:0;
		i,r|i,r:0;
		i,r|i,r:0;
		i,r|s,l:0;
		ary,r|s,l:0;
		ary,r|i,r:0;
		ary,r|f,r:0;
		`, M(mv), true, &a, &i, &k, &ks, &f,
		&iv1, &iv1ary, &iv2, &iv3, &iv4, &iv5,
		&iv6, &iv7, &iv8, &iv9, &iv10, &iv10ary,
		&iv11, &iv12, &arystr, &iv1ary2, &iv10ary2)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(k, ks, len(iv1ary), len(iv10ary))
	if *k != "10" || *ks[0] != "10" || *iv1 != 10 || *iv1ary[0] != 10 || *iv10 != 10 || *iv10ary[0] != 10 {
		t.Error("error")
		return
	}
	fmt.Println(arystr, iv1ary, iv10ary)
	if len(arystr) != 5 || len(iv1ary2) != 5 || len(iv10ary2) != 5 {
		t.Errorf("error,%v,%v,%v", len(arystr), len(iv1ary), len(iv10ary))
		return
	}
	fmt.Println(a, i, k, f)
}

func TestValidAttrFormatError(t *testing.T) {
	getter := ValueGetterF(func(key string) (v interface{}, err error) {
		if key == "not" {
			err = fmt.Errorf("not")
		} else if key == "ary1" {
			v = 1
		} else if key == "ary2" {
			v = []interface{}{"xxx"}
		}
		return
	})
	var err error
	//
	var sval string
	err = ValidAttrFormat(`
		not,r|s,l:~5;
		`, getter, true, &sval)
	if err == nil {
		t.Error("nil")
		return
	}
	//
	var iary []int
	err = ValidAttrFormat(`
		ary2,r|s,l:~5;
	`, getter, true, &iary)
	if err == nil {
		t.Error("nil")
		return
	}
	err = ValidAttrFormat(`
		xxx,r|s,l:~5;
	`, getter, true, &iary)
	if err == nil {
		t.Error("nil")
		return
	}
}

func TestEscape(t *testing.T) {
	//
	var a string
	err := ValidAttrFormat(`
		len,r|s,P:[^%N]*%N.*$;
		`, ValueGetterF(
		func(key string) (interface{}, error) {
			return "abc,ddf", nil
		},
	), true, &a)
	if err != nil {
		t.Error(err.Error())
		return
	}
}

func TestValidWeb(t *testing.T) {
	var (
		a   int
		b   string
		err error
		req *http.Request
	)
	req = httptest.NewRequest("GET", "http://localhost/?a=1&b=xxx", nil)
	err = QueryValidFormat(req, `
		a,r|i,r:0;
		b,r|s,l:0;
	`, &a, &b)
	if err != nil || a != 1 || b != "xxx" {
		t.Error(err)
		return
	}
	req.ParseForm()
	err = FormValidFormat(req, `
		a,r|i,r:0;
		b,r|s,l:0;
	`, &a, &b)
	if err != nil || a != 1 || b != "xxx" {
		t.Error(err)
		return
	}
	//
	req = httptest.NewRequest("POST", "http://localhost", bytes.NewBufferString("a=1&b=xxx"))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.ParseForm()
	err = PostFormValidFormat(req, `
		a,r|i,r:0;
		b,r|s,l:0;
	`, &a, &b)
	if err != nil || a != 1 || b != "xxx" {
		t.Error(err)
		return
	}
	//
	req = httptest.NewRequest("GET", "http://localhost", nil)
	req.PostForm = url.Values{}
	req.PostForm.Set("a", "1")
	req.PostForm.Set("b", "xxx")
	err = RequestValidFormat(req, `
		a,r|i,r:0;
		b,r|s,l:0;
	`, &a, &b)
	if err != nil || a != 1 || b != "xxx" {
		t.Error(err)
		return
	}
}

func TestValidNil(t *testing.T) {
	m := M{
		"a": nil,
		"b": []interface{}{},
	}
	var aryptr1, aryptr2 []*int64
	err := m.ValidFormat(`
		a,o|i,r:0;
		b,o|i,r:0;
	`, &aryptr1, &aryptr2)
	if err != nil {
		t.Error(err)
		return
	}
	err = m.ValidFormat(`
		a,r|i,r:0;
		b,r|i,r:0;
	`, &aryptr1, &aryptr2)
	if err == nil {
		t.Error(err)
		return
	}
}

func TestValidNoArray2Array(t *testing.T) {
	var data = `{"number":1000,"string":"1000"}`
	var m = map[string]interface{}{}
	var err = json.Unmarshal([]byte(data), &m)
	if err != nil {
		t.Error(err)
		return
	}
	mval := M(m)
	{
		var nval0 []int
		var nval1 []*int
		var nval2 []int64
		var nval3 []*int64
		var nval4 []float64
		var nval5 []*float64
		err = ValidAttrFormat(`
			number,O|I,R:-1;
			number,O|I,R:-1;
			number,O|I,R:-1;
			number,O|I,R:-1;
			number,O|I,R:-1;
			number,O|I,R:-1;
		`, mval, true, &nval0, &nval1, &nval2, &nval3, &nval4, &nval5)
		if err != nil ||
			len(nval0) != 1 || nval0[0] != 1000 ||
			len(nval1) != 1 || *nval1[0] != 1000 ||
			len(nval2) != 1 || nval2[0] != 1000 ||
			len(nval3) != 1 || *nval3[0] != 1000 ||
			len(nval4) != 1 || nval4[0] != 1000 ||
			len(nval5) != 1 || *nval5[0] != 1000 {
			t.Error(err)
			return
		}
	}
	{
		var nval0 []int
		var nval1 []*int
		var nval2 []int64
		var nval3 []*int64
		var nval4 []float64
		var nval5 []*float64
		err = ValidAttrFormat(`
			string,O|I,R:-1;
			string,O|I,R:-1;
			string,O|I,R:-1;
			string,O|I,R:-1;
			string,O|I,R:-1;
			string,O|I,R:-1;
		`, mval, true, &nval0, &nval1, &nval2, &nval3, &nval4, &nval5)
		if err != nil ||
			len(nval0) != 1 || nval0[0] != 1000 ||
			len(nval1) != 1 || *nval1[0] != 1000 ||
			len(nval2) != 1 || nval2[0] != 1000 ||
			len(nval3) != 1 || *nval3[0] != 1000 ||
			len(nval4) != 1 || nval4[0] != 1000 ||
			len(nval5) != 1 || *nval5[0] != 1000 {
			t.Error(err)
			return
		}
	}
}

type testSubStruct struct {
	_      string                 `xxx:"not exported"`
	xint   int                    `xxx:"not exported"`
	Int    int                    `json:"int"`
	Float  float64                `json:"float"`
	String string                 `json:"string"`
	Raw    map[string]interface{} `json:"raw"`
	Map    M                      `json:"map"`
}

type testStruct struct {
	Int    int                    `json:"int"`
	Float  float64                `json:"float"`
	String string                 `json:"string"`
	Raw    map[string]interface{} `json:"raw"`
	Map    M                      `json:"map"`
	Sub1   testSubStruct          `json:"sub1"`
	Sub2   *testSubStruct         `json:"sub2"`
}

type testStructPtr struct {
	Int    *int     `json:"int"`
	Float  *float64 `json:"float"`
	String *string  `json:"string"`
}

func TestValidStruct(t *testing.T) {
	value := testStruct{
		Int:    100,
		Float:  200,
		String: "300",
		Raw:    map[string]interface{}{"abc": 400},
		Map:    M{"abc": 500},
		Sub1: testSubStruct{
			xint:   100,
			Int:    100,
			Float:  200,
			String: "300",
			Raw:    map[string]interface{}{"abc": 400},
			Map:    M{"abc": 500},
		},
		Sub2: &testSubStruct{
			Int:    100,
			Float:  200,
			String: "300",
			Raw:    map[string]interface{}{"abc": 400},
			Map:    M{"abc": 500},
		},
	}
	var err error
	var intValue int
	var floatValue float64
	var stringValue, abc1Value, abc2Value string
	//
	//test json tag
	err = ValidStructAttrFormat(`
		int,R|I,R:0;
		float,R|I,R:0;
		string,R|S,L:0;
		raw/abc,R|S,L:0;
		map/abc,R|S,L:0;
	`, &value, true, &intValue, &floatValue, &stringValue, &abc1Value, &abc2Value)
	if err != nil || intValue != 100 || floatValue != 200 || stringValue != "300" || abc1Value != "400" || abc2Value != "500" {
		t.Errorf("%v,%v,%v,%v,%v,%v", err, intValue, floatValue, stringValue, abc1Value, abc2Value)
		return
	}
	err = ValidStructAttrFormat(`
		sub1/int,R|I,R:0;
		sub1/float,R|I,R:0;
		sub1/string,R|S,L:0;
		sub1/raw/abc,R|S,L:0;
		sub1/map/abc,R|S,L:0;
	`, &value, true, &intValue, &floatValue, &stringValue, &abc1Value, &abc2Value)
	if err != nil || intValue != 100 || floatValue != 200 || stringValue != "300" || abc1Value != "400" || abc2Value != "500" {
		t.Errorf("%v,%v,%v,%v,%v,%v", err, intValue, floatValue, stringValue, abc1Value, abc2Value)
		return
	}
	err = ValidStructAttrFormat(`
		sub2/int,R|I,R:0;
		sub2/float,R|I,R:0;
		sub2/string,R|S,L:0;
		sub2/raw/abc,R|S,L:0;
		sub2/map/abc,R|S,L:0;
	`, &value, true, &intValue, &floatValue, &stringValue, &abc1Value, &abc2Value)
	if err != nil || intValue != 100 || floatValue != 200 || stringValue != "300" || abc1Value != "400" || abc2Value != "500" {
		t.Errorf("%v,%v,%v,%v,%v,%v", err, intValue, floatValue, stringValue, abc1Value, abc2Value)
		return
	}
	//
	//test field name
	err = ValidStructAttrFormat(`
		Int,R|I,R:0;
		Float,R|I,R:0;
		String,R|S,L:0;
		Raw/abc,R|S,L:0;
		Map/abc,R|S,L:0;
	`, &value, true, &intValue, &floatValue, &stringValue, &abc1Value, &abc2Value)
	if err != nil || intValue != 100 || floatValue != 200 || stringValue != "300" || abc1Value != "400" || abc2Value != "500" {
		t.Errorf("%v,%v,%v,%v,%v,%v", err, intValue, floatValue, stringValue, abc1Value, abc2Value)
		return
	}
	err = ValidStructAttrFormat(`
		Sub1/Int,R|I,R:0;
		Sub1/Float,R|I,R:0;
		Sub1/String,R|S,L:0;
		Sub1/Raw/abc,R|S,L:0;
		Sub1/Map/abc,R|S,L:0;
	`, &value, true, &intValue, &floatValue, &stringValue, &abc1Value, &abc2Value)
	if err != nil || intValue != 100 || floatValue != 200 || stringValue != "300" || abc1Value != "400" || abc2Value != "500" {
		t.Errorf("%v,%v,%v,%v,%v,%v", err, intValue, floatValue, stringValue, abc1Value, abc2Value)
		return
	}
	err = ValidStructAttrFormat(`
		Sub2/int,R|I,R:0;
		Sub2/float,R|I,R:0;
		Sub2/string,R|S,L:0;
		Sub2/raw/abc,R|S,L:0;
		Sub2/map/abc,R|S,L:0;
	`, &value, true, &intValue, &floatValue, &stringValue, &abc1Value, &abc2Value)
	if err != nil || intValue != 100 || floatValue != 200 || stringValue != "300" || abc1Value != "400" || abc2Value != "500" {
		t.Errorf("%v,%v,%v,%v,%v,%v", err, intValue, floatValue, stringValue, abc1Value, abc2Value)
		return
	}
	//
	//test new struct
	err = NewStruct(&value).ValidFormat(`
		int,R|I,R:0;
		float,R|I,R:0;
		string,R|S,L:0;
		raw/abc,R|S,L:0;
		map/abc,R|S,L:0;
	`, &intValue, &floatValue, &stringValue, &abc1Value, &abc2Value)
	if err != nil || intValue != 100 || floatValue != 200 || stringValue != "300" || abc1Value != "400" || abc2Value != "500" {
		t.Errorf("%v,%v,%v,%v,%v,%v", err, intValue, floatValue, stringValue, abc1Value, abc2Value)
		return
	}
	//
	//test struct ptr
	valuePtr := testStructPtr{
		Int:    converter.IntPtr(100),
		Float:  converter.Float64Ptr(200),
		String: converter.StringPtr("300"),
	}
	err = NewStruct(&valuePtr).ValidFormat(`
		int,R|I,R:0;
		float,R|I,R:0;
		string,R|S,L:0;
	`, &intValue, &floatValue, &stringValue)
	if err != nil || intValue != 100 || floatValue != 200 || stringValue != "300" {
		t.Errorf("%v,%v,%v,%v,%v,%v", err, intValue, floatValue, stringValue, abc1Value, abc2Value)
		return
	}
	valuePtr2 := testStructPtr{}
	err = ValidFormat(`string,R|S,L:0;`, &valuePtr2, &intValue)
	if err == nil {
		t.Errorf("%v", err)
		return
	}
	err = ValidFormat(`string,R|S,L:0;`, &valuePtr2)
	if err == nil {
		t.Errorf("%v", err)
		return
	}
	//
	//test error
	func() {
		defer func() {
			recover()
		}()
		NewStruct(1)
	}()
}

type xxx map[string]interface{}

func (x xxx) RawMap() map[string]interface{} {
	return x
}

func TestValidFormat(t *testing.T) {
	var err error
	var intValue int
	//
	err = ValidFormat(`int,R|I,R:0`, M(map[string]interface{}{"int": 100}), &intValue)
	if err != nil || intValue != 100 {
		t.Error(err)
		return
	}
	req, _ := http.NewRequest("GET", "http://test/?int=100", nil)
	err = ValidFormat(`int,R|I,R:0`, req, &intValue)
	if err != nil || intValue != 100 {
		t.Error(err)
		return
	}
	err = ValidFormat(`int,R|I,R:0`, req.URL.Query(), &intValue)
	if err != nil || intValue != 100 {
		t.Error(err)
		return
	}
	err = ValidFormat(`int,R|I,R:0`, map[string]string{"int": "100"}, &intValue)
	if err != nil || intValue != 100 {
		t.Error(err)
		return
	}
	err = ValidFormat(`int,R|I,R:0`, map[string]interface{}{"int": "100"}, &intValue)
	if err != nil || intValue != 100 {
		t.Error(err)
		return
	}
	err = ValidFormat(`int,R|I,R:0`, xxx(map[string]interface{}{"int": 100}), &intValue)
	if err != nil || intValue != 100 {
		t.Error(err)
		return
	}
	err = ValidFormat(`int,R|I,R:0`, &testStruct{Int: 100}, &intValue)
	if err != nil || intValue != 100 {
		t.Error(err)
		return
	}
}

func TestCheck(t *testing.T) {
	var err error
	err = ValidFormat(`int,R|I,R:0`, M(map[string]interface{}{"int": 100}))
	if err != nil {
		t.Error(err)
		return
	}
	err = ValidFormat(`int,R|I,R:1000`, M(map[string]interface{}{"int": 100}))
	if err == nil {
		t.Error(err)
		return
	}
}

type Simple struct {
	A0 int64   `json:"a0" valid:"a0,r|i,r:0"`
	A1 *int64  `json:"a1" valid:"a1,r|i,r:0;"`
	AX []int64 `json:"ax" valid:"ax,r|i,r:0;"`
	XX string  `json:"xx"`
}

type SimpleArgs struct {
	A0 int64   `json:"a0" valid:"a0,r|i,r:0"`
	A1 *int64  `json:"a1" valid:"a1,r|i,r:0;"`
	AX []int64 `json:"ax" valid:"ax,r|i,r:0;"`
	XX struct {
		B0 int64  `json:"b0" valid:"b0,r|i,r:0"`
		B1 *int64 `json:"b1" valid:"b1,r|i,r:0;"`
	} `json:"xx" valid:"inline"`
}

type Product struct {
	Title    string  `json:"title,omitempty" valid:"title,r|s,l:0;"`         /* the product title */
	TitleSub *string `json:"title_sub,omitempty" valid:"title_sub,o|s,l:0;"` /* the product sub title */
}

func TestValid(t *testing.T) {
	var err error
	errObject := struct {
		A0 int64   `json:"a0" valid:"a0,r|i,r:0"`
		A1 *int64  `json:"a1" valid:"a1,r|i,r:0;"`
		A2 *int64  `json:"a2" valid:"a2,r|i,r:0;"`
		A3 int64   `json:"a3" valid:"|i,r:0;"`
		AX []int64 `json:"ax" valid:"ax,r|i,r:0;"`
		AY []int64 `json:"ay" valid:"ay,r|i,r:0;"`
		XX string  `json:"xx"`
	}{
		A2: converter.Int64Ptr(0),
		AY: []int64{0, 1},
	}
	err = Valid(&errObject, "#all", "")
	if err == nil {
		t.Error(err)
		return
	}
	ok0Object := struct {
		A0 int64    `json:"a0" valid:"a0,r|i,r:0"`
		A1 *int64   `json:"a1" valid:"a1,r|i,r:0;"`
		AX []int64  `json:"ax" valid:"ax,r|i,r:0;"`
		AY []*int64 `json:"ay" valid:"ay,r|i,r:0;"`
		S0 string   `json:"s0" valid:"s0,r|s,l:0;"`
		S1 *string  `json:"s1" valid:"s1,r|s,l:0;"`
		XX string   `json:"xx"`
	}{
		A0: 100,
		A1: converter.Int64Ptr(100),
		AX: []int64{1, 2, 3},
		AY: []*int64{converter.Int64Ptr(1), converter.Int64Ptr(2), converter.Int64Ptr(3)},
		S0: "abc",
		S1: converter.StringPtr("abc"),
	}
	err = Valid(&ok0Object, "#all", "")
	if err != nil {
		t.Error(err)
		return
	}
	ok1Object := struct {
		A0 int64    `json:"a0" valid:"a0,r|i,r:0"`
		A1 *int64   `json:"a1" valid:"a1,r|i,r:0;"`
		AX []int64  `json:"ax" valid:"ax,r|i,r:0;"`
		AY []*int64 `json:"ay" valid:"ay,r|i,r:0;"`
		S0 string   `json:"s0" valid:"s0,r|s,l:0;"`
		S1 *string  `json:"s1" valid:"s1,r|s,l:0;"`
		XX string   `json:"xx"`
	}{
		A0: 100,
		AX: []int64{1, 2, 3},
	}
	err = Valid(&ok1Object, "#all", "a1,ay,s0,s1")
	if err != nil {
		t.Error(err)
		return
	}
	err = Valid(&ok1Object, "#all", "^a0,ax")
	if err != nil {
		t.Error(err)
		return
	}
	ok2Object := struct {
		A0     int64 `json:"a0" valid:"a0,r|i,r:0"`
		Simple `filter:"inline"`
	}{
		A0: 100,
		Simple: Simple{
			A0: 100,
			A1: converter.Int64Ptr(100),
			AX: []int64{100},
		},
	}
	err = Valid(&ok2Object, "#all", "")
	if err != nil {
		t.Error(err)
		return
	}
	ok2Object.A0 = 0
	err = Valid(&ok2Object, "#all", "a0")
	if err != nil {
		t.Error(err)
		return
	}
	ok2Object.A0 = 0
	err = Valid(&ok2Object, "#all", "")
	if err == nil {
		t.Error(err)
		return
	}
	ok2Object.A0 = -1
	err = Valid(&ok2Object, "#all", "a0")
	if err == nil {
		t.Error(err)
		return
	}
	ok3Object := struct {
		Title    string  `json:"title,omitempty" valid:"title,r|s,l:0;"`         /* the product title */
		TitleSub *string `json:"title_sub,omitempty" valid:"title_sub,r|s,l:0;"` /* the product sub title */
	}{
		Title: "abc",
	}
	err = Valid(&ok3Object, "#all", "title_sub")
	if err != nil {
		t.Error(err)
		return
	}
	ok4Object := struct {
		Title    string  `json:"title,omitempty" valid:"title,r|s,l:0;"`         /* the product title */
		TitleSub *string `json:"title_sub,omitempty" valid:"title_sub,r|s,l:0;"` /* the product sub title */
	}{
		Title:    "abc",
		TitleSub: converter.StringPtr(""),
	}
	err = Valid(&ok4Object, "#all", "title_sub")
	if err != nil {
		t.Error(err)
		return
	}
	ok5Object := struct {
		Title    string  `json:"title,omitempty" valid:"title,r|s,l:0;"`         /* the product title */
		TitleSub *string `json:"title_sub,omitempty" valid:"title_sub,o|s,l:0;"` /* the product sub title */
	}{
		Title:    "abc",
		TitleSub: converter.StringPtr(""),
	}
	err = Valid(&ok5Object, "#all", "")
	if err != nil {
		t.Error(err)
		return
	}
	ok6Object := struct {
		Title    string  `json:"title,omitempty" valid:"title,r|s,l:0;"`         /* the product title */
		TitleSub *string `json:"title_sub,omitempty" valid:"title_sub,o|s,l:0;"` /* the product sub title */
		Value    int     `json:"value" valid:"value,r|i,r:0"`
	}{
		Title:    "abc",
		TitleSub: converter.StringPtr(""),
	}
	err = Valid(&ok6Object, "title,value|title_sub#all", "")
	if err != nil {
		t.Error(err)
		return
	}
	err6Object := struct {
		Title    string  `json:"title,omitempty" valid:"title,r|s,l:0;"`         /* the product title */
		TitleSub *string `json:"title_sub,omitempty" valid:"title_sub,r|s,l:0;"` /* the product sub title */
		Value    int     `json:"value" valid:"value,r|i,r:0"`
	}{
		Title:    "abc",
		TitleSub: converter.StringPtr(""),
	}
	err = Valid(&err6Object, "title,value|title_sub#all", "")
	if err == nil {
		t.Error(err)
		return
	}
}

func TestValidArgs(t *testing.T) {
	{
		simple := &Simple{}
		var b0, c0 string
		formats, args := ValidArgs(simple, "#all", `b0,r|s,l:0`, &b0, `c0,r|s,l:0;`, &c0)
		if len(formats) < 1 || strings.Count(formats, ";") != 5 || len(args) != 5 {
			t.Errorf("%v,%v,%v", formats, strings.Count(formats, ";"), len(args))
			return
		}
		err := ValidAttrFormat(formats, ValueGetterF(func(key string) (interface{}, error) {
			return "1", nil
		}), true, args...)
		if err != nil {
			t.Error(err)
			return
		}
	}
	{
		simpleArgs := &SimpleArgs{}
		formats, args := ValidArgs(simpleArgs, "#all")
		if len(formats) < 1 || strings.Count(formats, ";") != 5 || len(args) != 5 {
			t.Errorf("%v,%v,%v", formats, strings.Count(formats, ";"), len(args))
			return
		}
	}

}

type SetterTestObject struct {
	A0 int64
}

func (s *SetterTestObject) Set(v interface{}) (err error) {
	s.A0 = v.(int64)
	return
}

type SetterTestArray []SetterTestObject

func (s *SetterTestArray) Set(v interface{}) (err error) {
	value, err := converter.Int64Val(v)
	if err == nil {
		*s = append(*s, SetterTestObject{A0: value})
	}
	return
}

func TestSetter(t *testing.T) {
	var err error
	//setter
	var setter0 SetterTestObject
	var setter1 SetterTestArray
	err = ValidSetValue(&setter0, int64(0))
	if err != nil {
		t.Error(err)
		return
	}
	err = ValidAttrFormat(`x,R|I,R:0;x,R|I,R:0;x,R|I,R:0;`, ValueGetterF(func(key string) (interface{}, error) {
		return "1", nil
	}), true, &setter0, &setter1, ValueSetterF(func(i interface{}) error {
		if v, err := converter.Int64Val(i); err != nil || v != 1 {
			return fmt.Errorf("error")
		}
		return nil
	}))
	if err != nil || setter0.A0 != 1 || len(setter1) != 1 || setter1[0].A0 != 1 {
		t.Error(err)
		return
	}
	//not supported
	simple := &Simple{}
	err = ValidAttrFormat(`x,R|I,R:0`, ValueGetterF(func(key string) (interface{}, error) {
		return "1", nil
	}), true, simple)
	if err == nil {
		t.Error(err)
		return
	}
	//set error
	err = ValidAttrFormat(`x,R|I,R:0;`, ValueGetterF(func(key string) (interface{}, error) {
		return "x", nil
	}), true, &setter0)
	if err == nil {
		t.Error(err)
		return
	}
	fmt.Println(err)
	//set error
	err = ValidAttrFormat(`x,R|I,R:0;`, ValueGetterF(func(key string) (interface{}, error) {
		return "x", nil
	}), true, &setter1)
	if err == nil {
		t.Error(err)
		return
	}
	fmt.Println(err)
	//set error
	err = ValidAttrFormat(`x,R|S,L:0;`, ValueGetterF(func(key string) (interface{}, error) {
		return "x", nil
	}), true, &setter1)
	if err == nil {
		t.Error(err)
		return
	}
	fmt.Println(err)
}

type ScannerTestObject struct {
	A0 int64
}

func (s *ScannerTestObject) Scan(v interface{}) (err error) {
	s.A0 = v.(int64)
	return
}

type ScannerTestArray []ScannerTestObject

func (s *ScannerTestArray) Scan(v interface{}) (err error) {
	value, err := converter.Int64Val(v)
	if err == nil {
		*s = append(*s, ScannerTestObject{A0: value})
	}
	return
}

func TestScanner(t *testing.T) {
	var err error
	//setter
	var setter0 ScannerTestObject
	var setter1 ScannerTestArray
	err = ValidSetValue(&setter0, int64(0))
	if err != nil {
		t.Error(err)
		return
	}
	err = ValidAttrFormat(`x,R|I,R:0;x,R|I,R:0;x,R|I,R:0;`, ValueGetterF(func(key string) (interface{}, error) {
		return "1", nil
	}), true, &setter0, &setter1, ValueSetterF(func(i interface{}) error {
		if v, err := converter.Int64Val(i); err != nil || v != 1 {
			return fmt.Errorf("error")
		}
		return nil
	}))
	if err != nil || setter0.A0 != 1 || len(setter1) != 1 || setter1[0].A0 != 1 {
		t.Error(err)
		return
	}
	//not supported
	simple := &Simple{}
	err = ValidAttrFormat(`x,R|I,R:0`, ValueGetterF(func(key string) (interface{}, error) {
		return "1", nil
	}), true, simple)
	if err == nil {
		t.Error(err)
		return
	}
	//set error
	err = ValidAttrFormat(`x,R|I,R:0;`, ValueGetterF(func(key string) (interface{}, error) {
		return "x", nil
	}), true, &setter0)
	if err == nil {
		t.Error(err)
		return
	}
	fmt.Println(err)
	//set error
	err = ValidAttrFormat(`x,R|I,R:0;`, ValueGetterF(func(key string) (interface{}, error) {
		return "x", nil
	}), true, &setter1)
	if err == nil {
		t.Error(err)
		return
	}
	fmt.Println(err)
	//set error
	err = ValidAttrFormat(`x,R|S,L:0;`, ValueGetterF(func(key string) (interface{}, error) {
		return "x", nil
	}), true, &setter1)
	if err == nil {
		t.Error(err)
		return
	}
	fmt.Println(err)
}

type StringDef string
type IntDef int64
type FloatDef float64
type StringDefArray []StringDef
type IntDefArray []IntDef
type FloatDefArray []FloatDef

func TestValidDef(t *testing.T) {
	getter := ValueGetterF(func(key string) (interface{}, error) { return "1", nil })
	//
	var strVal StringDef
	var intVal IntDef
	var floatVal FloatDef
	err := ValidAttrFormat(`
		v,r|s,l:0;
		v,r|i,r:0;
		v,r|f,r:0;
	`, getter, true, &strVal, &intVal, &floatVal)
	if err != nil || strVal != "1" || intVal != 1 || floatVal != 1 {
		t.Error(err)
		return
	}
	var strArr StringDefArray
	var intArr IntDefArray
	var floatArr FloatDefArray
	err = ValidAttrFormat(`
		v,r|s,l:0;
		v,r|i,r:0;
		v,r|f,r:0;
	`, getter, true, &strArr, &intArr, &floatArr)
	if err != nil || len(strArr) < 1 || len(intArr) < 1 || len(floatArr) < 1 || strArr[0] != "1" || intArr[0] != 1 || floatArr[0] != 1 {
		t.Error(err)
		return
	}
}
