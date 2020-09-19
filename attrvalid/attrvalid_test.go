package attrvalid

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/codingeasygo/util/converter"
)

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
	v, err := ValidAttrTemple("测试", "r|s", "l:~10", true)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	v, err = ValidAttrTemple("测试测试测试测试", "r|s", "l:~10", true)
	if err == nil {
		t.Error("not error")
		return
	}
	//
	v, err = ValidAttrTemple("男", "r|s", "o:男~女", true)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	v, err = ValidAttrTemple("男ks", "r|s", "o:男~女", true)
	if err == nil {
		t.Error("not error")
		return
	}
	//
	v, err = ValidAttrTemple("centny@gmail.com", "r|s", "p:^.*\\@.*$", true)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	v, err = ValidAttrTemple("ks", "r|s", "p:^.*\\@.*$", true)
	if err == nil {
		t.Error("not error")
		return
	}
	//
	v, err = ValidAttrTemple("8", "o|i", "r:5~10", true)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	v, err = ValidAttrTemple("8", "o|i", "r:5~", true)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	v, err = ValidAttrTemple("12", "o|i", "r:5~10", true)
	if err == nil {
		t.Error("not error")
		return
	}
	//
	v, err = ValidAttrTemple("8", "o|f", "r:5~10", true)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	v, err = ValidAttrTemple("8", "o|f", "r:5~", true)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	v, err = ValidAttrTemple("12", "o|f", "r:5~10", true)
	if err == nil {
		t.Error("not error")
		return
	}
	//
	v, err = ValidAttrTemple("测", "o|s", "l:~8", true)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	v, err = ValidAttrTemple("测", "o|s", "l:2~", true)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	v, err = ValidAttrTemple("测度测度测度测度测度", "o|s", "l:~8", true)
	if err == nil {
		t.Error("not error")
		return
	}
	//
	v, err = ValidAttrTemple("测", "o|s", "l:2~8", true)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	v, err = ValidAttrTemple("a", "o|s", "l:2~8", true)
	if err == nil {
		t.Error("not error")
		return
	}
	//
	v, err = ValidAttrTemple("centny@gmail.com", "o|s", "p:^.*\\@.*$", true)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	v, err = ValidAttrTemple("ks", "o|s", "p:^.*\\@.*$", true)
	if err == nil {
		t.Error("not error")
		return
	}
	//
	v, err = ValidAttrTemple("1", "o|i", "o:1~2~3~4~5", true)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	v, err = ValidAttrTemple("11", "o|i", "o:1~2~3~4~5", true)
	if err == nil {
		t.Error("not error")
		return
	}
	//
	v, err = ValidAttrTemple("1.1", "o|f", "o:1.1~2.2~3.3~4~5", true)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	v, err = ValidAttrTemple("11", "o|f", "o:1~2~3~4~5", true)
	if err == nil {
		t.Error("not error")
		return
	}
	//
	v, err = ValidAttrTemple("测", "o|s", "l:a", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrTemple("测", "o|s", "KK:a", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrTemple("centny@gmail.com", "o|s", "p:*,..", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrTemple("测", "o|i", "r:8~9", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrTemple("测", "o|f", "r:8~9", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrTemple("测", "o|f", "o:8~9", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrTemple("测", "o|n", "r:8~9", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrTemple("5", "r|i", "r:~1", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrTemple("5", "r|i", "r:a~10", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrTemple("5", "r|i", "r:1~a", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrTemple("5", "r|f", "r:~1", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrTemple("5", "r|f", "r:a~10", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrTemple("5", "r|f", "r:1~a", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrTemple("5", "r|i", "m:1~a", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrTemple("5", "r|i", "o:1~a", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrTemple("5", "r|f", "o:1~a", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrTemple("5", "r|f", "m:1~k", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrTemple("5", "r|i", "o", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrTemple("5", "r", "o:1~10", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrTemple("", "r|i", "o:1~10", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrTemple("", "o|i", "o:1~10", true)
	if err != nil {
		t.Error(err.Error())
		return
	}
	v, err = ValidAttrTemple("a", "o|s", "l:a~8", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrTemple("a", "o|s", "l:2~a", true)
	if err == nil {
		t.Error("not error")
		return
	}
}

func TestValidAttrTemple2(t *testing.T) {
	v, err := ValidAttrTemple("测试", "R|S", "L:~10", true)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	v, err = ValidAttrTemple("测试测试测试测试", "R|S", "L:~10", true)
	if err == nil {
		t.Error("not error")
		return
	}
	//
	v, err = ValidAttrTemple("男", "R|S", "O:男~女", true)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	v, err = ValidAttrTemple("男ks", "R|S", "O:男~女", true)
	if err == nil {
		t.Error("not error")
		return
	}
	//
	v, err = ValidAttrTemple("centny@gmail.com", "R|S", "P:^.*\\@.*$", true)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	v, err = ValidAttrTemple("ks", "R|S", "P:^.*\\@.*$", true)
	if err == nil {
		t.Error("not error")
		return
	}
	//
	v, err = ValidAttrTemple("8", "O|I", "R:5~10", true)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	v, err = ValidAttrTemple("8", "O|I", "R:5~", true)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	v, err = ValidAttrTemple("12", "O|I", "R:5~10", true)
	if err == nil {
		t.Error("not error")
		return
	}
	//
	v, err = ValidAttrTemple("8", "O|F", "R:5~10", true)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	v, err = ValidAttrTemple("8", "O|F", "R:5~", true)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	v, err = ValidAttrTemple("12", "O|F", "R:5~10", true)
	if err == nil {
		t.Error("not error")
		return
	}
	//
	v, err = ValidAttrTemple("测", "O|S", "L:~8", true)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	v, err = ValidAttrTemple("测", "O|S", "L:2~", true)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	v, err = ValidAttrTemple("测度测度测度测度测度", "O|S", "L:~8", true)
	if err == nil {
		t.Error("not error")
		return
	}
	//
	v, err = ValidAttrTemple("测", "O|S", "L:2~8", true)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	v, err = ValidAttrTemple("a", "O|S", "L:2~8", true)
	if err == nil {
		t.Error("not error")
		return
	}
	//
	v, err = ValidAttrTemple("centny@gmail.com", "O|S", "P:^.*\\@.*$", true)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	v, err = ValidAttrTemple("ks", "O|S", "P:^.*\\@.*$", true)
	if err == nil {
		t.Error("not error")
		return
	}
	//
	v, err = ValidAttrTemple("1", "O|I", "O:1~2~3~4~5", true)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	v, err = ValidAttrTemple("11", "O|I", "O:1~2~3~4~5", true)
	if err == nil {
		t.Error("not error")
		return
	}
	//
	v, err = ValidAttrTemple("1.1", "O|F", "O:1.1~2.2~3.3~4~5", true)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	v, err = ValidAttrTemple("11", "O|F", "O:1~2~3~4~5", true)
	if err == nil {
		t.Error("not error")
		return
	}
	//
	v, err = ValidAttrTemple("测", "O|S", "L:a", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrTemple("测", "O|S", "KK:a", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrTemple("centny@gmail.com", "O|S", "P:*,..", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrTemple("测", "O|I", "R:8~9", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrTemple("测", "O|F", "R:8~9", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrTemple("测", "O|F", "O:8~9", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrTemple("测", "o|n", "R:8~9", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrTemple("5", "R|I", "R:~1", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrTemple("5", "R|I", "R:a~10", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrTemple("5", "R|I", "R:1~a", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrTemple("5", "R|F", "R:~1", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrTemple("5", "R|F", "R:a~10", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrTemple("5", "R|F", "R:1~a", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrTemple("5", "R|I", "m:1~a", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrTemple("5", "R|I", "O:1~a", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrTemple("5", "R|F", "O:1~a", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrTemple("5", "R|F", "m:1~k", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrTemple("5", "R|I", "o", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrTemple("5", "R", "O:1~10", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrTemple("", "R|I", "O:1~10", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrTemple("", "O|I", "O:1~10", true)
	if err != nil {
		t.Error(err.Error())
		return
	}
	v, err = ValidAttrTemple("a", "O|S", "L:a~8", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrTemple("a", "O|S", "L:2~a", true)
	if err == nil {
		t.Error("not error")
		return
	}
}

func TestValidAttrFormat(t *testing.T) {
	mv := map[string]interface{}{}
	mv["a"] = "abc"
	mv["i"] = "10"
	mv["f"] = "10.3"
	mv["ef"] = "20.3"
	mv["len"] = "11111111"
	mv["ary"] = "1,2,3,4,5"
	mv["ary2"] = "1,2,3,,4,5"
	mv["ary3"] = []interface{}{1, 2, 3, 4, 5}
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
		`, M(mv), true, &a, &i, &k, &ks, &f,
		&iv1, &iv1ary, &iv2, &iv3, &iv4, &iv5,
		&iv6, &iv7, &iv8, &iv9, &iv10, &iv10ary,
		&iv11, &iv12, &snot, &svary,
		&iv1ary2, &iv10ary2, &iv10ary3)
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
		ary1,r|s,l:~5;
	`, getter, true, &iary)
	if err == nil {
		t.Error("nil")
		return
	}
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
	err = NewStruct(&valuePtr2).ValidFormat(`
		int,R|I,R:0;
		float,R|I,R:0;
		string,R|S,L:0;
	`, &intValue, &floatValue, &stringValue)
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

func TestValid(t *testing.T) {
	var err error
	var intValue int
	//
	err = Valid(`int,R|I,R:0`, M(map[string]interface{}{"int": 100}), &intValue)
	if err != nil || intValue != 100 {
		t.Error(err)
		return
	}
	req, _ := http.NewRequest("GET", "http://test/?int=100", nil)
	err = Valid(`int,R|I,R:0`, req, &intValue)
	if err != nil || intValue != 100 {
		t.Error(err)
		return
	}
	err = Valid(`int,R|I,R:0`, req.URL.Query(), &intValue)
	if err != nil || intValue != 100 {
		t.Error(err)
		return
	}
	err = Valid(`int,R|I,R:0`, map[string]string{"int": "100"}, &intValue)
	if err != nil || intValue != 100 {
		t.Error(err)
		return
	}
	err = Valid(`int,R|I,R:0`, map[string]interface{}{"int": "100"}, &intValue)
	if err != nil || intValue != 100 {
		t.Error(err)
		return
	}
	err = Valid(`int,R|I,R:0`, xxx(map[string]interface{}{"int": 100}), &intValue)
	if err != nil || intValue != 100 {
		t.Error(err)
		return
	}
	err = Valid(`int,R|I,R:0`, &testStruct{Int: 100}, &intValue)
	if err != nil || intValue != 100 {
		t.Error(err)
		return
	}
}

func TestCheck(t *testing.T) {
	var err error
	err = Valid(`int,R|I,R:0`, M(map[string]interface{}{"int": 100}))
	if err != nil {
		t.Error(err)
		return
	}
	err = Valid(`int,R|I,R:1000`, M(map[string]interface{}{"int": 100}))
	if err == nil {
		t.Error(err)
		return
	}
}
