package xhttp

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/codingeasygo/util/xmap"
)

const (
	ContentTypeForm = "application/x-www-form-urlencoded"
	ContentTypeJSON = "application/json;charset=utf-8"
	ContentTypeXML  = "application/xml;charset=utf-8"
)

func init() {
	initClient(true)
}

func initClient(insecureSkipVerify bool) {
	DefaultTransport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecureSkipVerify},
	}
	jar, _ := cookiejar.New(nil)
	DefaultClient = &http.Client{
		Transport: DefaultTransport,
		Jar:       jar,
	}

}

var DefaultTransport *http.Transport

var DefaultClient *http.Client

func DisableInsecureVerify() {
	initClient(true)
}

func EnableInsecureVerify() {
	initClient(false)
}

//GetBytes will do http request and read the bytes response
func GetBytes(format string, args ...interface{}) (data []byte, err error) {
	data, _, err = GetHeaderBytes(nil, format, args...)
	return
}

//GetHeaderBytes will do http request and read the text response
func GetHeaderBytes(header xmap.M, format string, args ...interface{}) (data []byte, res *http.Response, err error) {
	remote := fmt.Sprintf(format, args...)
	req, err := http.NewRequest("GET", remote, nil)
	if err != nil {
		return
	}
	for k, v := range header {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}
	res, err = DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	data, err = ioutil.ReadAll(res.Body)
	if res.StatusCode != 200 {
		err = fmt.Errorf("status code is %v", res.StatusCode)
	}
	return
}

//GetText will do http request and read the text response
func GetText(format string, args ...interface{}) (data string, err error) {
	data, _, err = GetHeaderText(nil, format, args...)
	return
}

//GetHeaderText will do http request and read the text response
func GetHeaderText(header xmap.M, format string, args ...interface{}) (data string, res *http.Response, err error) {
	bys, res, err := GetHeaderBytes(header, format, args...)
	data = string(bys)
	return
}

//GetMap will do http request, read reponse and parse to map
func GetMap(format string, args ...interface{}) (data xmap.Valuable, err error) {
	data, _, err = GetHeaderMap(nil, format, args...)
	return
}

//GetHeaderMap will do http request, read reponse and parse to map
func GetHeaderMap(header xmap.M, format string, args ...interface{}) (data xmap.Valuable, res *http.Response, err error) {
	text, res, err := GetHeaderBytes(header, format, args...)
	if err == nil {
		data, err = xmap.Parse(text)
	}
	return
}

//PostBytes will do http request and read the bytes response
func PostBytes(body io.Reader, format string, args ...interface{}) (data []byte, err error) {
	data, _, err = PostHeaderBytes(nil, body, format, args...)
	return
}

//PostTypeBytes will do http request and read the bytes response
func PostTypeBytes(contentType string, body io.Reader, format string, args ...interface{}) (data []byte, err error) {
	data, _, err = PostHeaderBytes(xmap.M{"Content-Type": contentType}, body, format, args...)
	return
}

//PostHeaderBytes will do http request and read the text response
func PostHeaderBytes(header xmap.M, body io.Reader, format string, args ...interface{}) (data []byte, res *http.Response, err error) {
	remote := fmt.Sprintf(format, args...)
	req, err := http.NewRequest("POST", remote, body)
	if err != nil {
		return
	}
	for k, v := range header {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}
	res, err = DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	data, err = ioutil.ReadAll(res.Body)
	if res.StatusCode != 200 {
		err = fmt.Errorf("status code is %v", res.StatusCode)
	}
	return
}

//PostText will do http request and read the text response
func PostText(body io.Reader, format string, args ...interface{}) (data string, err error) {
	data, _, err = PostHeaderText(nil, body, format, args...)
	return
}

//PostTypeText will do http request and read the text response
func PostTypeText(contentType string, body io.Reader, format string, args ...interface{}) (data string, err error) {
	data, _, err = PostHeaderText(xmap.M{"Content-Type": contentType}, body, format, args...)
	return
}

//PostHeaderText will do http request and read the text response
func PostHeaderText(header xmap.M, body io.Reader, format string, args ...interface{}) (data string, res *http.Response, err error) {
	bys, res, err := PostHeaderBytes(header, body, format, args...)
	data = string(bys)
	return
}

//PostMap will do http request, read reponse and parse to map
func PostMap(body io.Reader, format string, args ...interface{}) (data xmap.Valuable, err error) {
	data, _, err = PostHeaderMap(nil, body, format, args...)
	return
}

//PostTypeMap will do http request, read reponse and parse to map
func PostTypeMap(contentType string, body io.Reader, format string, args ...interface{}) (data xmap.Valuable, err error) {
	data, _, err = PostHeaderMap(xmap.M{"Content-Type": contentType}, body, format, args...)
	return
}

//PostHeaderMap will do http request, read reponse and parse to map
func PostHeaderMap(header xmap.M, body io.Reader, format string, args ...interface{}) (data xmap.Valuable, res *http.Response, err error) {
	text, res, err := PostHeaderBytes(header, body, format, args...)
	if err == nil {
		data, err = xmap.Parse(text)
	}
	return
}

//PostJSONMap will do http request, read reponse and parse to map
func PostJSONMap(v interface{}, format string, args ...interface{}) (data xmap.Valuable, err error) {
	bys, err := json.Marshal(v)
	if err == nil {
		data, _, err = PostHeaderMap(xmap.M{"Content-Type": ContentTypeJSON}, bytes.NewBuffer(bys), format, args...)
	}
	return
}

//PostXMLText will do http request, read reponse and parse to map
func PostXMLText(v interface{}, format string, args ...interface{}) (data string, err error) {
	bys, err := xml.Marshal(v)
	if err == nil {
		data, _, err = PostHeaderText(xmap.M{"Content-Type": ContentTypeXML}, bytes.NewBuffer(bys), format, args...)
	}
	return
}

//PostFormText will do http request, read reponse and parse to map
func PostFormText(form xmap.M, format string, args ...interface{}) (data string, err error) {
	query := url.Values{}
	for k, v := range form {
		query.Set(k, fmt.Sprintf("%v", v))
	}
	buf := bytes.NewBufferString(query.Encode())
	data, _, err = PostHeaderText(xmap.M{"Content-Type": ContentTypeForm}, buf, format, args...)
	return
}

//PostFormMap will do http request, read reponse and parse to map
func PostFormMap(form xmap.M, format string, args ...interface{}) (data xmap.Valuable, err error) {
	query := url.Values{}
	for k, v := range form {
		query.Set(k, fmt.Sprintf("%v", v))
	}
	buf := bytes.NewBufferString(query.Encode())
	data, _, err = PostHeaderMap(xmap.M{"Content-Type": ContentTypeForm}, buf, format, args...)
	return
}

//CreateMultipartBody will create multipart body
func CreateMultipartBody(fields xmap.M) (*bytes.Buffer, string) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	for k, v := range fields {
		bodyWriter.WriteField(k, fmt.Sprintf("%v", v))
	}
	ctype := bodyWriter.FormDataContentType()
	bodyWriter.Close()
	return bodyBuf, ctype
}

//PostMultipartBytes will do http request, read reponse and parse to map
func PostMultipartBytes(header, fields xmap.M, format string, args ...interface{}) (data []byte, res *http.Response, err error) {
	bodyBuf, ctype := CreateMultipartBody(fields)
	remote := fmt.Sprintf(format, args...)
	req, err := http.NewRequest("POST", remote, bodyBuf)
	if err != nil {
		return
	}
	for k, v := range header {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}
	req.Header.Set("Content-Type", ctype)
	res, err = DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	data, err = ioutil.ReadAll(res.Body)
	if res.StatusCode != 200 {
		err = fmt.Errorf("status code is %v", res.StatusCode)
	}
	return
}

//PostMultipartText will do http request, read reponse and parse to map
func PostMultipartText(header, fields xmap.M, format string, args ...interface{}) (data string, err error) {
	bys, _, err := PostMultipartBytes(header, fields, format, args...)
	if err == nil {
		data = string(bys)
	}
	return
}

//PostMultipartMap will do http request, read reponse and parse to map
func PostMultipartMap(header, fields xmap.M, format string, args ...interface{}) (data xmap.Valuable, err error) {
	bys, _, err := PostMultipartBytes(header, fields, format, args...)
	if err == nil {
		data, err = xmap.Parse(bys)
	}
	return
}

type FileBodyTask struct {
	reader *io.PipeReader
	writer *io.PipeWriter
}

func CreateFileBodyTask(fields xmap.M, filekey string, filename string) (*FileBodyTask, io.Reader, string) {
	task := &FileBodyTask{}
	reader, ctype := task.Start(fields, filekey, filename)
	return task, reader, ctype
}

func (f *FileBodyTask) Start(fields xmap.M, filekey string, filename string) (io.Reader, string) {
	f.reader, f.writer = io.Pipe()
	bodyWriter := multipart.NewWriter(f.writer)
	go func() {
		err := f.run(bodyWriter, fields, filekey, filename)
		bodyWriter.Close()
		if err == nil {
			f.writer.Close()
		} else {
			f.writer.CloseWithError(err)
		}
	}()
	return f.reader, bodyWriter.FormDataContentType()
}

func (f *FileBodyTask) run(bodyWriter *multipart.Writer, fields xmap.M, filekey string, filename string) error {
	for k, v := range fields {
		bodyWriter.WriteField(k, fmt.Sprintf("%v", v))
	}
	fileWriter, err := bodyWriter.CreateFormFile(filekey, filename)
	if err == nil {
		var file *os.File
		file, err = os.Open(filename)
		if err == nil {
			defer file.Close()
			_, err = io.Copy(fileWriter, file)
		}
	}
	return err
}

func (f *FileBodyTask) Close() error {
	if f.reader != nil {
		f.reader.Close()
	}
	if f.writer != nil {
		f.writer.Close()
	}
	return nil
}

func UploadText(fields xmap.M, filekey, filename, format string, args ...interface{}) (text string, err error) {
	text, _, err = UploadHeaderText(nil, fields, filekey, filename, format, args...)
	return
}

func UploadHeaderText(header xmap.M, fields xmap.M, filekey, filename, format string, args ...interface{}) (text string, res *http.Response, err error) {
	var ctype string
	var bodyBuf io.Reader
	var task *FileBodyTask
	task, bodyBuf, ctype = CreateFileBodyTask(fields, filekey, filename)
	defer task.Close()
	remote := fmt.Sprintf(format, args...)
	req, err := http.NewRequest("POST", remote, bodyBuf)
	if err != nil {
		return
	}
	for k, v := range header {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}
	req.Header.Set("Content-Type", ctype)
	res, err = DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	bys, err := ioutil.ReadAll(res.Body)
	text = string(bys)
	if res.StatusCode != 200 {
		err = fmt.Errorf("status code is %v", res.StatusCode)
	}
	return
}

func UploadMap(fields xmap.M, filekey, filename, format string, args ...interface{}) (data xmap.Valuable, err error) {
	text, _, err := UploadHeaderText(nil, fields, filekey, filename, format, args...)
	if err == nil {
		data, err = xmap.Parse(text)
	}
	return
}

func Download(saveto string, format string, args ...interface{}) (saved int64, err error) {
	saved, err = DownloadHeader(saveto, nil, format, args...)
	return
}

func DownloadHeader(saveto string, header xmap.M, format string, args ...interface{}) (saved int64, err error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(format, args...), nil)
	if err != nil {
		return
	}
	for k, v := range header {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}
	res, err := DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		err = fmt.Errorf("status code is %v", res.StatusCode)
		return
	}
	savepath := saveto
	if info, err := os.Stat(saveto); err == nil && info.IsDir() {
		var filename string
		disposition := res.Header.Get("Content-Disposition")
		parts := strings.SplitN(disposition, "filename", 2)
		if len(parts) == 2 {
			filename = strings.SplitN(parts[1], ";", 2)[0]
			filename = strings.TrimSpace(filename)
			filename = strings.Trim(filename, "=\"")
		}
		if len(filename) < 1 {
			_, filename = path.Split(req.URL.Path)
		}
		if len(filename) < 1 {
			filename = "index.html"
		}
		savepath = filepath.Join(saveto, filename)
	}
	file, err := os.OpenFile(savepath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err == nil {
		defer file.Close()
		saved, err = io.Copy(file, res.Body)
	}
	return
}
