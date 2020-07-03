package xhttp

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
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
)

func init() {
	initClient(true)
}

func initClient(insecureSkipVerify bool) {
	DefaultTransport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecureSkipVerify},
	}
	DefaultClient = &http.Client{
		Transport: DefaultTransport,
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
func GetHeaderBytes(header map[string]interface{}, format string, args ...interface{}) (data []byte, res *http.Response, err error) {
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
func GetHeaderText(header map[string]interface{}, format string, args ...interface{}) (data string, res *http.Response, err error) {
	bys, res, err := GetHeaderBytes(header, format, args...)
	data = string(bys)
	return
}

//GetMap will do http request, read reponse and parse to map
func GetMap(format string, args ...interface{}) (data xmap.M, err error) {
	data, _, err = GetHeaderMap(nil, format, args...)
	return
}

//GetHeaderMap will do http request, read reponse and parse to map
func GetHeaderMap(header map[string]interface{}, format string, args ...interface{}) (data xmap.M, res *http.Response, err error) {
	text, res, err := GetHeaderBytes(header, format, args...)
	if err == nil {
		err = json.Unmarshal(text, &data)
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
	data, _, err = PostHeaderBytes(map[string]interface{}{"Content-Type": contentType}, body, format, args...)
	return
}

//PostHeaderBytes will do http request and read the text response
func PostHeaderBytes(header map[string]interface{}, body io.Reader, format string, args ...interface{}) (data []byte, res *http.Response, err error) {
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
	data, _, err = PostHeaderText(map[string]interface{}{"Content-Type": contentType}, body, format, args...)
	return
}

//PostHeaderText will do http request and read the text response
func PostHeaderText(header map[string]interface{}, body io.Reader, format string, args ...interface{}) (data string, res *http.Response, err error) {
	bys, res, err := PostHeaderBytes(header, body, format, args...)
	data = string(bys)
	return
}

//PostMap will do http request, read reponse and parse to map
func PostMap(body io.Reader, format string, args ...interface{}) (data xmap.M, err error) {
	data, _, err = PostHeaderMap(nil, body, format, args...)
	return
}

//PostTypeMap will do http request, read reponse and parse to map
func PostTypeMap(contentType string, body io.Reader, format string, args ...interface{}) (data xmap.M, err error) {
	data, _, err = PostHeaderMap(map[string]interface{}{"Content-Type": contentType}, body, format, args...)
	return
}

//PostHeaderMap will do http request, read reponse and parse to map
func PostHeaderMap(header map[string]interface{}, body io.Reader, format string, args ...interface{}) (data xmap.M, res *http.Response, err error) {
	text, res, err := PostHeaderBytes(header, body, format, args...)
	if err == nil {
		err = json.Unmarshal(text, &data)
	}
	return
}

//PostJSONMap will do http request, read reponse and parse to map
func PostJSONMap(v interface{}, format string, args ...interface{}) (data xmap.M, err error) {
	bys, err := json.Marshal(v)
	if err == nil {
		data, _, err = PostHeaderMap(map[string]interface{}{"Content-Type": ContentTypeJSON}, bytes.NewBuffer(bys), format, args...)
	}
	return
}

//PostFormMap will do http request, read reponse and parse to map
func PostFormMap(form map[string]interface{}, format string, args ...interface{}) (data xmap.M, err error) {
	query := url.Values{}
	for k, v := range form {
		query.Set(k, fmt.Sprintf("%v", v))
	}
	buf := bytes.NewBufferString(query.Encode())
	data, _, err = PostHeaderMap(map[string]interface{}{"Content-Type": ContentTypeForm}, buf, format, args...)
	return
}

func CreateMultipartBody(fields map[string]interface{}) (string, *bytes.Buffer) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	for k, v := range fields {
		bodyWriter.WriteField(k, fmt.Sprintf("%v", v))
	}
	ctype := bodyWriter.FormDataContentType()
	bodyWriter.Close()
	return ctype, bodyBuf
}

type FileBodyTask struct {
	reader *io.PipeReader
	writer *io.PipeWriter
}

func CreateFileBodyTask(fields map[string]interface{}, filekey string, filename string) (*FileBodyTask, io.Reader, string) {
	task := &FileBodyTask{}
	reader, ctype := task.Start(fields, filekey, filename)
	return task, reader, ctype
}

func (f *FileBodyTask) Start(fields map[string]interface{}, filekey string, filename string) (io.Reader, string) {
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

func (f *FileBodyTask) run(bodyWriter *multipart.Writer, fields map[string]interface{}, filekey string, filename string) error {
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

func UploadText(filekey, filename, format string, args ...interface{}) (text string, err error) {
	text, _, err = UploadHeaderText(nil, nil, filekey, filename, format, args...)
	return
}

func UploadHeaderText(header map[string]interface{}, fields map[string]interface{}, filekey, filename, format string, args ...interface{}) (text string, res *http.Response, err error) {
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

func UploadMap(filekey, filename, format string, args ...interface{}) (data xmap.M, err error) {
	text, _, err := UploadHeaderText(nil, nil, filekey, filename, format, args...)
	if err == nil {
		err = json.Unmarshal([]byte(text), &data)
	}
	return
}

func Download(saveto string, format string, args ...interface{}) (saved int64, err error) {
	saved, err = DownloadHeader(saveto, nil, format, args...)
	return
}

func DownloadHeader(saveto string, header map[string]interface{}, format string, args ...interface{}) (saved int64, err error) {
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
