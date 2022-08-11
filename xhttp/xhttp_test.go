package xhttp

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/codingeasygo/util/xjson"
	"github.com/codingeasygo/util/xmap"
)

func init() {
	EnableInsecureVerify()
	DisableInsecureVerify()
	DisableCookie()
	ClearCookie()
	EnableCookie()
	ClearCookie()
	NewRawClient(nil)
	NewClient(&http.Client{})
}

func TestGet(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var format = r.URL.Query().Get("format")
		switch format {
		case "json":
			xjson.WriteJSON(w, map[string]interface{}{
				"abc": 1,
			})
		case "text":
			fmt.Fprintf(w, "1")
		case "header":
			fmt.Fprintf(w, "%v", r.Header.Get("abc"))
		default:
			w.WriteHeader(500)
			fmt.Fprintf(w, "error")
		}
	}))
	//
	bval, err := GetBytes("%v/?format=text", ts.URL)
	if err != nil || string(bval) != "1" {
		t.Error(err)
		return
	}
	bval, _, err = GetHeaderBytes(nil, "%v/?format=text", ts.URL)
	if err != nil || string(bval) != "1" {
		t.Error(err)
		return
	}
	//
	sval, err := GetText("%v/?format=text", ts.URL)
	if err != nil || sval != "1" {
		t.Error(err)
		return
	}
	sval, _, err = GetHeaderText(xmap.M{"abc": 1}, "%v/?format=text", ts.URL)
	if err != nil || sval != "1" {
		t.Error(err)
		return
	}
	//
	mval, err := GetMap("%v/?format=json", ts.URL)
	if err != nil {
		t.Errorf("%v,%v", mval, err)
		return
	}
	mval, _, err = GetHeaderMap(nil, "%v/?format=json", ts.URL)
	if err != nil {
		t.Error(err)
		return
	}
	var ival int
	err = mval.ValidFormat(`abc,r|i,r:0;`, &ival)
	if err != nil || ival != 1 {
		t.Error(err)
		return
	}
	{ //get json
		var jval xmap.M
		jval = xmap.New()
		err := GetJSON(&jval, "%v/?format=json", ts.URL)
		if err != nil || len(jval) < 1 {
			t.Errorf("%v,%v", mval, err)
			return
		}
		jval = xmap.New()
		_, err = GetHeaderJSON(&jval, nil, "%v/?format=json", ts.URL)
		if err != nil || len(jval) < 1 {
			t.Error(err)
			return
		}

	}
	//
	//test error
	_, err = GetText("%v/?format=error", ts.URL)
	if err == nil {
		t.Error(err)
		return
	}
	_, err = GetText("%v/\x06?format=text", ts.URL)
	if err == nil {
		t.Error(err)
		return
	}
	_, err = GetText("%v/?format=text", "http://127.0.0.1:32")
	if err == nil {
		t.Error(err)
		return
	}
}

func TestPost(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var format = r.URL.Query().Get("format")
		switch format {
		case "json":
			xjson.WriteJSON(w, map[string]interface{}{
				"abc": 1,
			})
		case "text":
			fmt.Fprintf(w, "1")
		case "header":
			fmt.Fprintf(w, "%v", r.Header.Get("abc"))
		case "form":
			r.ParseForm()
			xjson.WriteJSON(w, map[string]interface{}{
				"abc": r.PostForm.Get("abc"),
			})
		case "part":
			r.ParseMultipartForm(1024)
			xjson.WriteJSON(w, map[string]interface{}{
				"abc": r.MultipartForm.Value["abc"][0],
			})
		case "body":
			io.Copy(w, r.Body)
		default:
			w.WriteHeader(500)
			fmt.Fprintf(w, "error")
		}
	}))
	//
	bval, err := PostBytes(nil, "%v/?format=text", ts.URL)
	if err != nil || string(bval) != "1" {
		t.Error(err)
		return
	}
	bval, _, err = MethodBytes("POST", nil, nil, "%v/?format=text", ts.URL)
	if err != nil || string(bval) != "1" {
		t.Error(err)
		return
	}
	bval, err = PostTypeBytes(ContentTypeForm, nil, "%v/?format=text", ts.URL)
	if err != nil || string(bval) != "1" {
		t.Error(err)
		return
	}
	bval, _, err = PostHeaderBytes(xmap.M{"Content-Type": ContentTypeForm}, nil, "%v/?format=text", ts.URL)
	if err != nil || string(bval) != "1" {
		t.Error(err)
		return
	}
	//
	sval, err := PostText(nil, "%v/?format=text", ts.URL)
	if err != nil || sval != "1" {
		t.Error(err)
		return
	}
	sval, _, err = MethodText("POST", nil, nil, "%v/?format=text", ts.URL)
	if err != nil || sval != "1" {
		t.Error(err)
		return
	}
	sval, err = PostTypeText(ContentTypeForm, nil, "%v/?format=text", ts.URL)
	if err != nil || sval != "1" {
		t.Error(err)
		return
	}
	sval, _, err = PostHeaderText(xmap.M{"Content-Type": ContentTypeForm}, nil, "%v/?format=text", ts.URL)
	if err != nil || sval != "1" {
		t.Error(err)
		return
	}
	//
	var ival int
	mval, err := PostMap(nil, "%v?format=json", ts.URL)
	if err != nil {
		t.Error(err)
		return
	}
	err = mval.ValidFormat(`abc,r|i,r:0;`, &ival)
	if err != nil || ival != 1 {
		t.Error(err)
		return
	}
	mval, _, err = MethodMap("POST", nil, nil, "%v?format=json", ts.URL)
	if err != nil {
		t.Error(err)
		return
	}
	err = mval.ValidFormat(`abc,r|i,r:0;`, &ival)
	if err != nil || ival != 1 {
		t.Error(err)
		return
	}
	mval, _, err = PostHeaderMap(xmap.M{"Content-Type": ContentTypeForm}, nil, "%v?format=json", ts.URL)
	if err != nil {
		t.Error(err)
		return
	}
	err = mval.ValidFormat(`abc,r|i,r:0;`, &ival)
	if err != nil || ival != 1 {
		t.Error(err)
		return
	}
	mval, err = PostTypeMap(ContentTypeForm, nil, "%v?format=json", ts.URL)
	if err != nil {
		t.Error(err)
		return
	}
	err = mval.ValidFormat(`abc,r|i,r:0;`, &ival)
	if err != nil || ival != 1 {
		t.Error(err)
		return
	}
	mval, err = PostJSONMap(map[string]interface{}{"abc": "1"}, "%v?format=body", ts.URL)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = PostJSONMap(t.Fail, "%v?format=body", ts.URL)
	if err == nil {
		t.Error(err)
		return
	}
	err = mval.ValidFormat(`abc,r|i,r:0;`, &ival)
	if err != nil || ival != 1 {
		t.Error(err)
		return
	}
	sval, err = PostFormText(map[string]interface{}{"abc": "1"}, "%v?format=json", ts.URL)
	if err != nil || len(sval) < 1 {
		t.Error(err)
		return
	}
	mval, err = PostFormMap(map[string]interface{}{"abc": "1"}, "%v?format=json", ts.URL)
	if err != nil {
		t.Error(err)
		return
	}
	err = mval.ValidFormat(`abc,r|i,r:0;`, &ival)
	if err != nil || ival != 1 {
		t.Error(err)
		return
	}
	sval, err = PostXMLText(t, "%v?format=body", ts.URL)
	if err != nil {
		t.Errorf("%v,%v", sval, err)
		return
	}
	_, err = PostXMLText(t.Fail, "%v?format=body", ts.URL)
	if err == nil {
		t.Errorf("%v,%v", sval, err)
		return
	}
	bval, _, err = PostMultipartBytes(nil, xmap.M{"abc": "123"}, "%v?format=part", ts.URL)
	if err != nil || !strings.Contains(string(bval), "123") {
		t.Errorf("err:%v,text:%v", err, bval)
		return
	}
	sval, err = PostMultipartText(nil, xmap.M{"abc": "123"}, "%v?format=part", ts.URL)
	if err != nil || !strings.Contains(sval, "123") {
		t.Errorf("err:%v,text:%v", err, sval)
		return
	}
	mval, err = PostMultipartMap(xmap.M{"abc": "123"}, xmap.M{"abc": "123"}, "%v?format=part", ts.URL)
	if err != nil || mval.Str("abc") != "123" {
		t.Errorf("err:%v,text:%v", err, sval)
		return
	}
	{ //json result
		var jval xmap.M
		jval = xmap.New()
		err = PostJSON(&jval, bytes.NewBufferString("{}"), "%v?format=json", ts.URL)
		if err != nil || len(jval) < 1 {
			t.Error(err)
			return
		}
		jval = xmap.New()
		err = PostTypeJSON(&jval, ContentTypeJSON, bytes.NewBufferString("{}"), "%v?format=json", ts.URL)
		if err != nil || len(jval) < 1 {
			t.Error(err)
			return
		}
		jval = xmap.New()
		_, err = PostHeaderJSON(&jval, nil, bytes.NewBufferString("{}"), "%v?format=json", ts.URL)
		if err != nil || len(jval) < 1 {
			t.Error(err)
			return
		}
		jval = xmap.New()
		err = PostJSONJSON(&jval, xmap.New(), "%v?format=json", ts.URL)
		if err != nil || len(jval) < 1 {
			t.Error(err)
			return
		}
		err = PostJSONJSON(&jval, t.Error, "%v?format=json", ts.URL)
		if err == nil {
			t.Error(err)
			return
		}
	}
	//
	//test error
	_, err = PostText(nil, "%v/?format=error", ts.URL)
	if err == nil {
		t.Error(err)
		return
	}
	_, err = PostText(nil, "%v/\x06?format=text", ts.URL)
	if err == nil {
		t.Error(err)
		return
	}
	_, err = PostText(nil, "%v/?format=text", "http://127.0.0.1:32")
	if err == nil {
		t.Error(err)
		return
	}
	_, _, err = MethodText("POST", nil, nil, "%v/?format=error", ts.URL)
	if err != nil {
		t.Error(err)
		return
	}
	_, _, err = MethodText("POST", nil, nil, "%v/?format=error", "http://127.0.0.1:32")
	if err == nil {
		t.Error(err)
		return
	}
	//
	_, err = PostMultipartText(nil, xmap.M{"abc": "123"}, "%v?format=error", ts.URL)
	if err == nil {
		t.Errorf("err:%v,text:%v", err, sval)
		return
	}
	_, err = PostMultipartText(nil, xmap.M{"abc": "123"}, "%v?\x06format=text", ts.URL)
	if err == nil {
		t.Errorf("err:%v,text:%v", err, sval)
		return
	}
	_, err = PostMultipartText(nil, xmap.M{"abc": "123"}, "%v?format=text", "http://127.0.0.1:32")
	if err == nil {
		t.Errorf("err:%v,text:%v", err, sval)
		return
	}
}

func TestUpload(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var format = r.URL.Query().Get("format")
		switch format {
		case "json":
			r.ParseMultipartForm(1024)
			f, _ := r.MultipartForm.File["abc"][0].Open()
			data, _ := ioutil.ReadAll(f)
			f.Close()
			xjson.WriteJSON(w, map[string]interface{}{
				"abc": string(data),
			})
		case "text":
			r.ParseMultipartForm(1024)
			f, _ := r.MultipartForm.File["abc"][0].Open()
			io.Copy(w, f)
			f.Close()
		default:
			w.WriteHeader(500)
			fmt.Fprintf(w, "error")
		}
	}))
	var ival int
	os.Remove("/tmp/xhttp_test.txt")
	ioutil.WriteFile("/tmp/xhttp_test.txt", []byte("1"), os.ModePerm)
	//
	sval, err := UploadText(nil, "abc", "/tmp/xhttp_test.txt", "%v?format=text", ts.URL)
	if err != nil || sval != "1" {
		t.Error(err)
		return
	}
	sval, _, err = UploadHeaderText(xmap.M{"h1": 1}, xmap.M{"f1": 1}, "abc", "/tmp/xhttp_test.txt", "%v?format=text", ts.URL)
	if err != nil || sval != "1" {
		t.Error(err)
		return
	}
	//
	mval, err := UploadMap(nil, "abc", "/tmp/xhttp_test.txt", "%v?format=json", ts.URL)
	if err != nil {
		t.Error(err)
		return
	}
	err = mval.ValidFormat(`abc,r|i,r:0;`, &ival)
	if err != nil || ival != 1 {
		t.Error(err)
		return
	}
	//
	//test error
	_, err = UploadText(nil, "abc", "/tmp/xhttp_test.txt", "%v/\x01?format=json", ts.URL)
	if err == nil {
		t.Error(err)
		return
	}
	_, err = UploadText(nil, "abc", "/tmp/xhttp_test.txt", "%v?format=error", ts.URL)
	if err == nil {
		t.Error(err)
		return
	}
	_, err = UploadText(nil, "abc", "/tmp/xxxxx", "%v?format=text", ts.URL)
	if err == nil {
		t.Error(err)
		return
	}
}

func TestDownload(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var format = r.URL.Query().Get("format")
		switch format {
		case "filename":
			w.Header().Set("Content-Disposition", "binary;filename=\"xhttp_downt.txt\"")
			fmt.Fprintf(w, "1")
		case "text":
			fmt.Fprintf(w, "1")
		default:
			w.WriteHeader(500)
			fmt.Fprintf(w, "error")
		}
	}))
	//
	os.Remove("/tmp/xhttp_downt.txt")
	_, err := Download("/tmp/", "%v?format=filename", ts.URL)
	if err != nil {
		t.Error(err)
		return
	}
	bval, err := ioutil.ReadFile("/tmp/xhttp_downt.txt")
	if err != nil || string(bval) != "1" {
		t.Error(err)
		return
	}
	//
	os.Remove("/tmp/xhttp_downt2.txt")
	_, err = Download("/tmp/", "%v/xhttp_downt2.txt?format=text", ts.URL)
	if err != nil {
		t.Error(err)
		return
	}
	bval, err = ioutil.ReadFile("/tmp/xhttp_downt2.txt")
	if err != nil || string(bval) != "1" {
		t.Error(err)
		return
	}
	//
	os.Remove("/tmp/index.html")
	_, err = DownloadHeader("/tmp/", xmap.M{"x": 1}, "%v/?format=text", ts.URL)
	if err != nil {
		t.Error(err)
		return
	}
	bval, err = ioutil.ReadFile("/tmp/index.html")
	if err != nil || string(bval) != "1" {
		t.Error(err)
		return
	}
	//
	//test error
	_, err = Download("/tmp/", "%v\x01/?format=text", ts.URL)
	if err == nil {
		t.Error(err)
		return
	}
	_, err = Download("/tmp/", "%v/?format=error", ts.URL)
	if err == nil {
		t.Error(err)
		return
	}
	_, err = Download("/tmp/", "http://127.0.0.1:32")
	if err == nil {
		t.Error(err)
		return
	}
}

func TestCreateMultipartBody(t *testing.T) {
	CreateMultipartBody(map[string]interface{}{"A": 1})
}

// func TestReadAllStr(t *testing.T) {
// 	res, _ := readAllStr(nil)
// 	if len(res) > 0 {
// 		t.Error("not empty")
// 		return
// 	}
// 	r, _ := os.Open("name")
// 	res, _ = readAllStr(r)
// 	if len(res) > 0 {
// 		t.Error("not empty")
// 		return
// 	}
// }

// func TestHTTP2(t *testing.T) {
// 	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.Write([]byte("{\"code\":1}"))
// 	}))
// 	ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.Write([]byte("{\"code:1}"))
// 	}))
// 	ts3 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 	}))
// 	res := HTTPGet2(ts.URL)
// 	fmt.Println(res)
// 	res = HTTPGet2(ts2.URL)
// 	fmt.Println(res)
// 	res = HTTPGet2(ts3.URL)
// 	fmt.Println(res)
// 	_, err := HPostF(ts.URL, map[string]string{"ma": "123"}, "abc", "")
// 	if err == nil {
// 		t.Error("not error")
// 		return
// 	}
// 	_, err = HPostF("hhh", map[string]string{"ma": "123"}, "abc", "test.txt")
// 	if err == nil {
// 		t.Error("not error")
// 		return
// 	}
// 	_, err = HPostF(ts.URL, map[string]string{"ma": "123"}, "abc", "test.txt")
// 	if err != nil {
// 		t.Error(err.Error())
// 		return
// 	}
// 	_, _, err = HTTPClient.HPostF_H(ts.URL, map[string]string{"ma": "123"}, map[string]string{"ma": "123"}, "abc", "test.txt")
// 	if err != nil {
// 		t.Error(err.Error())
// 		return
// 	}
// 	_, err = HPostF(ts.URL, map[string]string{"ma": "123"}, "abc", "/tmp")
// 	if err == nil {
// 		t.Error("not error")
// 		return
// 	}
// 	_, err = HPostF2(ts.URL, map[string]string{"ma": "123"}, "abc", "test.txt")
// 	if err != nil {
// 		t.Error(err.Error())
// 		return
// 	}
// 	_, _, err = HTTPClient.HGet_H(map[string]string{"ma": "123"}, "%s?abc=%s", ts.URL, "1111")
// 	if err != nil {
// 		t.Error(err.Error())
// 		return
// 	}
// 	HPostF2("kkk", map[string]string{"ma": "123"}, "abc", "test.txt")
// 	HPostF2("123%34%56://s", map[string]string{"ma": "123"}, "abc", "test.txt")
// 	HTTPPost(ts.URL, map[string]string{"ma": "123"})
// 	HTTPPost2(ts.URL, map[string]string{"ma": "123"})
// 	HTTPPost2("jhj", map[string]string{"ma": "123"})
// 	//
// 	HTTPClient.DLoad("/tmp/aa.log", map[string]string{"ma": "123"}, "%s", ts.URL)
// 	fmt.Println(HTTPClient.DLoad("/sg/aa.log", map[string]string{"ma": "123"}, "%s", ts.URL))
// }
// func TestHTTPErr(t *testing.T) {
// 	fmt.Println(HPostF2("123%45%6", map[string]string{"ma": "123"}, "abc", "test.txt"))
// 	fmt.Println(HGet("123%45%6"))
// 	fmt.Println(HPostN("123%45%6", "ABcc", nil))
// 	fmt.Println(DLoad("spath", "123%45%6"))
// }
// func TestHpp(t *testing.T) {
// 	HGet("123%45%67://s")
// 	HGet("kk")
// 	HGet2("kk")
// 	HPost("jjjj", nil)
// 	HPost2("kkk", nil)
// 	HGet2("kkk")
// }

// //
// type osize struct {
// }

// func (o *osize) Size() int64 {
// 	return 100
// }

// type ostat struct {
// 	F *os.File
// }

// func (o *ostat) Stat() (os.FileInfo, error) {
// 	return o.F.Stat()
// }
// func TestFormFSzie(t *testing.T) {
// 	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.Write([]byte("{\"code\":1}"))
// 		src, _, err := r.FormFile("abc")
// 		if err != nil {
// 			t.Error(err.Error())
// 			return
// 		}
// 		fsize := FormFSzie(src)
// 		if fsize < 1 {
// 			t.Error("not size")
// 		}
// 	}))
// 	_, err := HPostF(ts.URL, map[string]string{"ma": "123"}, "abc", "test.txt")
// 	if err != nil {
// 		t.Error(err.Error())
// 	}
// 	f, _ := os.Open("test.txt")
// 	defer f.Close()
// 	fsize := FormFSzie(f)
// 	if fsize < 1 {
// 		t.Error("not right")
// 	}
// 	fsize = FormFSzie(&osize{})
// 	if fsize < 1 {
// 		t.Error("not right")
// 	}
// }

// func TestMap2Query(t *testing.T) {
// 	mv := map[string]interface{}{}
// 	mv["abc"] = "123"
// 	mv["dd"] = "ee"
// 	fmt.Println(Map2Query(mv))
// }

// func TestAHttpPost(t *testing.T) {
// 	ts := httptest.NewServer(http.HandlerFunc(
// 		func(w http.ResponseWriter, r *http.Request) {
// 			r.ParseMultipartForm(10000000)
// 			// r.PostFormValue(key)
// 			fmt.Println(r.FormValue("ab"))
// 			fmt.Println(r.PostFormValue("ab"))
// 		}))
// 	HPostF2s(ts.URL, map[string]string{
// 		"ab": "233",
// 	}, "", "")
// }

// func HPostF2s(url string, fields map[string]string, fkey string, fp string) (string, error) {
// 	ctype, bodyBuf, err := CreateFormBody2(fields, fkey, fp)
// 	if err != nil {
// 		return "", err
// 	}
// 	res, err := http.Post(url, ctype, bodyBuf)
// 	if err != nil {
// 		return "", err
// 	}
// 	return readAllStr(res.Body)
// }

// func CreateFormBody2(fields map[string]string, fkey string, fp string) (string, *bytes.Buffer, error) {
// 	bodyBuf := &bytes.Buffer{}
// 	bodyWriter := multipart.NewWriter(bodyBuf)
// 	for k, v := range fields {
// 		bodyWriter.WriteField(k, v)
// 	}
// 	w, _ := bodyWriter.CreateFormField("kkk")
// 	w.Write([]byte("kkkkkkk"))
// 	if len(fkey) > 0 {
// 		fileWriter, err := bodyWriter.CreateFormFile(fkey, fp)
// 		if err != nil {
// 			return "", nil, err
// 		}
// 		fh, err := os.Open(fp)
// 		if err != nil {
// 			return "", nil, err
// 		}
// 		defer fh.Close()
// 		_, err = io.Copy(fileWriter, fh)
// 		if err != nil {
// 			return "", nil, err
// 		}
// 	}
// 	ctype := bodyWriter.FormDataContentType()
// 	bodyWriter.Close()
// 	return ctype, bodyBuf, nil
// }

// type ErrWriter struct {
// }

// func (e *ErrWriter) Write(p []byte) (n int, err error) {
// 	return 0, errors.New("test erro")
// }

// func TestCreateFileForm(t *testing.T) {
// 	bodyWriter := multipart.NewWriter(&ErrWriter{})
// 	err := CreateFileForm(bodyWriter, "sss", "sss")
// 	if err == nil {
// 		t.Error("not error")
// 	}
// 	fmt.Println(err.Error())
// }

// func TestJson2Ary(t *testing.T) {
// 	ary, err := Json2Ary(`
// 		[1,2,"ss"]
// 		`)
// 	if err != nil {
// 		t.Error(err.Error())
// 		return
// 	}
// 	fmt.Println(ary)
// 	_, err = Json2Ary(`
// 		[1,2,ss"]
// 		`)
// 	if err == nil {
// 		t.Error("not error")
// 		return
// 	}
// }

// type ErrReader struct {
// }

// func (e *ErrReader) Read(p []byte) (n int, err error) {
// 	return 0, errors.New("error")
// }
// func TestPostN(t *testing.T) {
// 	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.Write([]byte("Ok"))
// 	}))
// 	_, data, err := HPostN(ts.URL, "text/plain", bytes.NewBuffer([]byte("WWW")))
// 	if err != nil {
// 		t.Error(err.Error())
// 		return
// 	}
// 	fmt.Println(data)
// 	fmt.Println(HPostN("kkk://sssss", "text/plain", bytes.NewBuffer([]byte("WWW"))))
// 	fmt.Println(HPostN("http:///kkkfjdfsfsd", "text/plain", bytes.NewBuffer([]byte("WWW"))))
// 	// fmt.Println(HPostN("http://www.baidu.com", "text/plain", &ErrReader{}))
// }
// func TestPostN2(t *testing.T) {
// 	_, err := http.NewRequest("POT", "123%45%6://www.ss.com?", nil)
// 	fmt.Println(err)
// }

// func TestHttps(t *testing.T) {
// 	fmt.Println(HGet("https://qnear.com"))
// }
