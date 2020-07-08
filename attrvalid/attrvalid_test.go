package attrvalid

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
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
