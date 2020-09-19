package xprop

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/codingeasygo/util/converter"
)

func assert(v bool) {
	if !v {
		panic("error")
	}
}

func TestEnvReplace(t *testing.T) {
	f := NewConfig()
	f.Base = "."
	f.SetValue("a", "b111111")
	fmt.Println(f.EnvReplace("sss${a} ${abc} ${da} ${HOME} ${} ${CONF_DIR}"))
	f.Clear()
	if f.Length() != 0 {
		t.Error("error")
	}
}

func TestInit(t *testing.T) {
	f := NewConfig()
	err := f.LoadWait("not_found.properties", false)
	if err == nil {
		panic("init error")
	}
	err = f.Load("test_data.properties")
	if err != nil {
		t.Error(err.Error())
		return
	}
	if f.Str("abc_conf") != `1
2
3` {
		fmt.Println("->\n", f.Str("abc_conf"))
		t.Error("error")
		return
	}
	for key, val := range f.config {
		fmt.Println(key, ":", val)
	}
	fmt.Println(f.StrVal("inta"))
	fmt.Println(f.StrVal("nfound"))
	fmt.Println(f.IntVal("inta"))
	fmt.Println(f.IntVal("nfound"))
	fmt.Println(f.IntVal("a"))
	fmt.Println(f.Float64Val("floata"))
	fmt.Println(f.Float64Val("nfound"))
	fmt.Println(f.Float64Val("a"))
	fmt.Println(f.StrVal("abc_conf"))
	fmt.Printf("\n\n\n\n")
	f.Delete("nfound")
	f.Delete("a")
	f.Print()
	f.PrintSection("loc")
}

func TestOpenError(t *testing.T) {
	f := NewConfig()
	fmt.Println(exec.Command("touch", "/tmp/fcg").Run())
	fmt.Println(exec.Command("chmod", "000", "/tmp/fcg").Run())
	fi, e := os.Open("/tmp/fcg")
	fmt.Println(fi, e)
	err := f.Load("/tmp/fcg")
	if err == nil {
		panic("init error")
	}
	fmt.Println(exec.Command("rm", "-f", "/tmp/fcg").Run())
}

func TestValue(t *testing.T) {
	f := NewConfig()
	err := f.Load("test_data.properties?ukk=123")
	if err != nil {
		t.Error(err)
		return
	}
	//
	assert(0 != f.IntDef(0, "inta"))
	assert(0 != f.Int64Def(0, "inta"))
	assert(0 != f.Uint64Def(0, "inta"))
	assert(0 != f.Float64Def(0, "floata"))
	assert("0" != f.StrDef("0", "floata"))
	assert(nil != f.MapDef(nil, "json"))
	//
	assert(0 == f.IntDef(0, "notxxx"))
	assert(0 == f.Int64Def(0, "notxxx"))
	assert(0 == f.Uint64Def(0, "notxxx"))
	assert(0 == f.Float64Def(0, "notxxx"))
	assert("0" == f.StrDef("0", "notxxx"))
	assert(nil == f.MapDef(nil, "notxxx"))
	//
	assert(nil != f.ArrayIntDef(nil, "inta"))
	assert(nil != f.ArrayInt64Def(nil, "inta"))
	assert(nil != f.ArrayUint64Def(nil, "inta"))
	assert(nil != f.ArrayFloat64Def(nil, "floata"))
	assert(nil != f.ArrayStrDef(nil, "floata"))
	assert(nil != f.ArrayMapDef(nil, "json2"))
	//
	assert(nil == f.ArrayIntDef(nil, "notxxx"))
	assert(nil == f.ArrayInt64Def(nil, "notxxx"))
	assert(nil == f.ArrayUint64Def(nil, "notxxx"))
	assert(nil == f.ArrayFloat64Def(nil, "notxxx"))
	assert(nil == f.ArrayStrDef(nil, "notxxx"))
	assert(nil == f.ArrayMapDef(nil, "notxxx"))
	//
	assert(os.ModePerm == f.FileModeDef(0, "mode"))
	assert(os.ModePerm == f.FileModeDef(os.ModePerm, "modex"))
	assert(os.ModePerm == f.FileModeDef(os.ModePerm, "json2"))
}

func TestLoad(t *testing.T) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	handlerc := 0
	proph := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerc++
		if handlerc > 1 {
			w.WriteHeader(200)
			fmt.Fprintf(w, "a=1")
		} else {
			w.WriteHeader(500)
			fmt.Fprintf(w, "testing")
		}
	})
	confh := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerc++
		if handlerc > 1 {
			w.WriteHeader(200)
			fmt.Fprintf(w, `
[a]
1
			`)
		} else {
			w.WriteHeader(500)
			fmt.Fprintf(w, "testing")
		}
	})
	var err error
	var config *Config
	var ts *httptest.Server
	//
	handlerc = 0
	ts = httptest.NewServer(confh)
	config, err = LoadConf(ts.URL + "/x.conf")
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(converter.JSON(config.config))
	assert(config.IntDef(0, "a") == 1)
	ts.Close()
	//
	handlerc = 0
	ts = httptest.NewServer(proph)
	config, err = LoadConf(ts.URL)
	if err != nil {
		t.Error(err.Error())
		return
	}
	assert(config.IntDef(0, "a") == 1)
	ts.Close()
	//
	handlerc = 0
	ts = httptest.NewTLSServer(proph)
	config, err = LoadConf(ts.URL)
	if err != nil {
		t.Error(err.Error())
		return
	}
	assert(config.IntDef(0, "a") == 1)
	ts.Close()
	//
	config, err = LoadConf("data:text/conf,[a]\n1")
	if err != nil {
		t.Error(err.Error())
		return
	}
	assert(config.IntDef(0, "a") == 1)
	//
	config, err = LoadConf("data:text/prop,a=1")
	if err != nil {
		t.Error(err.Error())
		return
	}
	assert(config.IntDef(0, "a") == 1)
	//
	os.Remove("/tmp/xprop_test.properties")
	go func() {
		time.Sleep(300 * time.Millisecond)
		ioutil.WriteFile("/tmp/xprop_test.properties", []byte("a=1"), os.ModePerm)
	}()
	config, err = LoadConf("/tmp/xprop_test.properties")
	if err != nil {
		t.Error(err.Error())
		return
	}
	assert(config.IntDef(0, "a") == 1)
	//
	config = NewConfig()
	config.LoadPropString("a=1")
	assert(config.IntDef(0, "a") == 1)
	//
	config = NewConfig()
	config.LoadPropReader("", bytes.NewBufferString("a=1"))
	assert(config.IntDef(0, "a") == 1)
	//
	config = NewConfig()
	ioutil.WriteFile("/tmp/xprop_test.properties", []byte("a=1"), os.ModePerm)
	config.LoadFile("/tmp/xprop_test.properties")
	assert(config.IntDef(0, "a") == 1)
	//
	handlerc = 1
	config = NewConfig()
	ts = httptest.NewServer(proph)
	err = config.LoadWeb(ts.URL)
	if err != nil {
		t.Error(err.Error())
		return
	}
	assert(config.IntDef(0, "a") == 1)
	ts.Close()
}

func TestSection(t *testing.T) {
	config, err := LoadConf("test_data.properties?ukk=123")
	if err != nil {
		t.Error(err)
		return
	}
	if config.StrDef("", "ukk") != "123" {
		t.Error("not right")
		return
	}
	if config.StrDef("", "/ukk") != "123" {
		t.Error("not right")
		return
	}
	if config.StrDef("", "abc/txabc") != "1" {
		t.Error("not right")
		return
	}
	if config.StrDef("", "/abc/txabc") != "1" {
		t.Error("not right")
		return
	}
	if config.StrDef("", "abd/dxabc") != "1" {
		t.Error("not right")
		return
	}
}

func TestError(t *testing.T) {
	config := NewConfig()
	assert(nil == config.exec("", "   #xxxx", false))
	assert(nil != config.LoadWait("data:text/prop,@l:xxxx", false))
	assert(nil != config.LoadWait("data:text/prop,@l:http://127.0.0.1:2332", false))
	if _, err := config.webGet("\x01"); err == nil {
		t.Error("error")
		return
	}
}

func TestMerge(t *testing.T) {
	var cfga, cfgb, cfgc = NewConfig(), NewConfig(), NewConfig()
	cfga.Masks = map[string]string{
		".*_DB_.*": ".*:[^@]*",
	}
	cfga.Load("test_a.properties")
	cfgb.Load("test_b.properties")
	if len(cfga.Seces) != 3 {
		t.Error("error")
		return
	}
	cfga.Merge(cfgb)
	if len(cfga.Seces) != 4 {
		t.Error("error")
		return
	}
	cfga.Merge(nil)
	cfgd := cfga.Clone()
	if len(cfgd.Seces) != 4 {
		t.Error("error")
		return
	}
	cfgc.MergeSection("a", cfga)
	if len(cfgc.Seces) != 1 {
		t.Error("error")
		return
	}
}

func TestRange(t *testing.T) {
	config, _ := LoadConf("test_a.properties")
	config.Range("a", func(key string, val interface{}) {
		fmt.Println(key, val)
	})
}

// func TestFileMode(t *testing.T) {
// 	config := NewConfig()
// 	config.SetValue("abc", "077")
// 	fmt.Println(cfg.FileModeV("abc", os.ModePerm))
// 	fmt.Println(cfg.FileModeV("xx", os.ModePerm))
// 	fmt.Println(cfg.FileModeV("abc2", os.ModePerm))
// }

func TestPrintMask(t *testing.T) {
	config := NewConfig()
	config.SetValue("loc/DB_DB_URL", "cny:sco@localhost")
	config.SetValue("xxx/DB_DB_URL", "cny:sco@localhost")
	config.Masks = map[string]string{
		".*_DB_.*": ".*:[^@]*",
	}
	config.Print()
	config.PrintSection("loc")
}
