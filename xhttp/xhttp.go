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
	//ContentTypeForm is web form content type
	ContentTypeForm = "application/x-www-form-urlencoded"
	//ContentTypeJSON is json content type
	ContentTypeJSON = "application/json;charset=utf-8"
	//ContentTypeXML is xml content type
	ContentTypeXML = "application/xml;charset=utf-8"
)

var insecureSkipVerify = false
var enableCookie = false

func init() {
	initClient()
}

func initClient() {
	DefaultTransport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecureSkipVerify},
	}
	DefaultClient = &http.Client{
		Transport: DefaultTransport,
	}
	if enableCookie {
		DefaultClient.Jar, _ = cookiejar.New(nil)
	}
}

//DefaultTransport is default http transport
var DefaultTransport *http.Transport

//DefaultClient is default http client
var DefaultClient *http.Client

//Shared is share client
var Shared *Client = &Client{Raw: defaultRaw}

//DisableInsecureVerify will disable verify insecure
func DisableInsecureVerify() {
	insecureSkipVerify = true
	initClient()
}

//EnableInsecureVerify will enable verify insecure
func EnableInsecureVerify() {
	insecureSkipVerify = false
	initClient()
}

//EnableCookie will enable cookie
func EnableCookie() {
	enableCookie = true
	initClient()
}

//DisableCookie will disable cookie
func DisableCookie() {
	enableCookie = false
	initClient()
}

//ClearCookie will clear cookie
func ClearCookie() {
	if enableCookie {
		DefaultClient.Jar, _ = cookiejar.New(nil)
	} else {
		DefaultClient.Jar = nil
	}
}

func defaultRaw(method, uri string, header xmap.M, body io.Reader) (req *http.Request, res *http.Response, err error) {
	c := &RawClient{C: DefaultClient}
	req, res, err = c.RawRequest(method, uri, header, body)
	return
}

//RawClient is http raw request impl
type RawClient struct {
	C *http.Client
}

//RawRequest will do raw request
func (r *RawClient) RawRequest(method, uri string, header xmap.M, body io.Reader) (req *http.Request, res *http.Response, err error) {
	req, err = http.NewRequest(method, uri, body)
	if err != nil {
		return
	}
	for k, v := range header {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}
	res, err = r.C.Do(req)
	return
}

//GetBytes will get bytes from remote
func GetBytes(format string, args ...interface{}) (data []byte, err error) {
	return Shared.GetBytes(format, args...)
}

//GetHeaderBytes will get bytes from remote
func GetHeaderBytes(header xmap.M, format string, args ...interface{}) (data []byte, res *http.Response, err error) {
	return Shared.GetHeaderBytes(header, format, args...)
}

//GetText will get text from remote
func GetText(format string, args ...interface{}) (data string, err error) {
	return Shared.GetText(format, args...)
}

//GetHeaderText will get text from remote
func GetHeaderText(header xmap.M, format string, args ...interface{}) (data string, res *http.Response, err error) {
	return Shared.GetHeaderText(header, format, args...)
}

//GetMap will get map from remote
func GetMap(format string, args ...interface{}) (data xmap.M, err error) {
	return Shared.GetMap(format, args...)
}

//GetHeaderMap will get map from remote
func GetHeaderMap(header xmap.M, format string, args ...interface{}) (data xmap.M, res *http.Response, err error) {
	return Shared.GetHeaderMap(header, format, args...)
}

//GetJSON will get json from remote
func GetJSON(result interface{}, format string, args ...interface{}) (err error) {
	return Shared.GetJSON(result, format, args...)
}

//GetHeaderMap will get map from remote
func GetHeaderJSON(result interface{}, header xmap.M, format string, args ...interface{}) (res *http.Response, err error) {
	return Shared.GetHeaderJSON(result, header, format, args...)
}

//PostBytes will get bytes from remote
func PostBytes(body io.Reader, format string, args ...interface{}) (data []byte, err error) {
	return Shared.PostBytes(body, format, args...)
}

//PostTypeBytes will get bytes from remote
func PostTypeBytes(contentType string, body io.Reader, format string, args ...interface{}) (data []byte, err error) {
	return Shared.PostTypeBytes(contentType, body, format, args...)
}

//PostHeaderBytes will get bytes from remote
func PostHeaderBytes(header xmap.M, body io.Reader, format string, args ...interface{}) (data []byte, res *http.Response, err error) {
	return Shared.PostHeaderBytes(header, body, format, args...)
}

//PostText will get text from remote
func PostText(body io.Reader, format string, args ...interface{}) (data string, err error) {
	return Shared.PostText(body, format, args...)
}

//PostTypeText will get text from remote
func PostTypeText(contentType string, body io.Reader, format string, args ...interface{}) (data string, err error) {
	return Shared.PostTypeText(contentType, body, format, args...)
}

//PostHeaderText will get text from remote
func PostHeaderText(header xmap.M, body io.Reader, format string, args ...interface{}) (data string, res *http.Response, err error) {
	return Shared.PostHeaderText(header, body, format, args...)
}

//PostMap will get map from remote
func PostMap(body io.Reader, format string, args ...interface{}) (data xmap.M, err error) {
	return Shared.PostMap(body, format, args...)
}

//PostTypeMap will get map from remote
func PostTypeMap(contentType string, body io.Reader, format string, args ...interface{}) (data xmap.M, err error) {
	return Shared.PostTypeMap(contentType, body, format, args...)
}

//PostHeaderMap will get map from remote
func PostHeaderMap(header xmap.M, body io.Reader, format string, args ...interface{}) (data xmap.M, res *http.Response, err error) {
	return Shared.PostHeaderMap(header, body, format, args...)
}

//PostJSONMap will get map from remote
func PostJSONMap(body interface{}, format string, args ...interface{}) (data xmap.M, err error) {
	return Shared.PostJSONMap(body, format, args...)
}

//PostMap will get map from remote
func PostJSON(result interface{}, body io.Reader, format string, args ...interface{}) (err error) {
	return Shared.PostJSON(result, body, format, args...)
}

//PostTypeMap will get map from remote
func PostTypeJSON(result interface{}, contentType string, body io.Reader, format string, args ...interface{}) (err error) {
	return Shared.PostTypeJSON(result, contentType, body, format, args...)
}

//PostHeaderMap will get map from remote
func PostHeaderJSON(result interface{}, header xmap.M, body io.Reader, format string, args ...interface{}) (res *http.Response, err error) {
	return Shared.PostHeaderJSON(result, header, body, format, args...)
}

//PostJSONMap will get map from remote
func PostJSONJSON(result interface{}, body interface{}, format string, args ...interface{}) (err error) {
	return Shared.PostJSONJSON(result, body, format, args...)
}

//MethodBytes will do http request, read reponse and parse to bytes
func MethodBytes(method string, header xmap.M, body io.Reader, format string, args ...interface{}) (data []byte, res *http.Response, err error) {
	return Shared.MethodBytes(method, header, body, format, args...)
}

//MethodBytes will do http request, read reponse and parse to string
func MethodText(method string, header xmap.M, body io.Reader, format string, args ...interface{}) (data string, res *http.Response, err error) {
	return Shared.MethodText(method, header, body, format, args...)
}

//MethodBytes will do http request, read reponse and parse to map
func MethodMap(method string, header xmap.M, body io.Reader, format string, args ...interface{}) (data xmap.M, res *http.Response, err error) {
	return Shared.MethodMap(method, header, body, format, args...)
}

//PostXMLText will get text from remote
func PostXMLText(v interface{}, format string, args ...interface{}) (data string, err error) {
	return Shared.PostXMLText(v, format, args...)
}

//PostFormText will get text from remote
func PostFormText(form xmap.M, format string, args ...interface{}) (data string, err error) {
	return Shared.PostFormText(form, format, args...)
}

//PostFormMap will get map from remote
func PostFormMap(form xmap.M, format string, args ...interface{}) (data xmap.M, err error) {
	return Shared.PostFormMap(form, format, args...)
}

//PostMultipartBytes will get bytes from remote
func PostMultipartBytes(header, fields xmap.M, format string, args ...interface{}) (data []byte, res *http.Response, err error) {
	return Shared.PostMultipartBytes(header, fields, format, args...)
}

//PostMultipartText will get bytes from remote
func PostMultipartText(header, fields xmap.M, format string, args ...interface{}) (data string, err error) {
	return Shared.PostMultipartText(header, fields, format, args...)
}

//PostMultipartMap will get map from remote
func PostMultipartMap(header, fields xmap.M, format string, args ...interface{}) (data xmap.M, err error) {
	return Shared.PostMultipartMap(header, fields, format, args...)
}

//UploadText will get text from remote
func UploadText(fields xmap.M, filekey, filename, format string, args ...interface{}) (text string, err error) {
	return Shared.UploadText(fields, filekey, filename, format, args...)
}

//UploadHeaderText will get text from remote
func UploadHeaderText(header xmap.M, fields xmap.M, filekey, filename, format string, args ...interface{}) (text string, res *http.Response, err error) {
	return Shared.UploadHeaderText(header, fields, filekey, filename, format, args...)
}

//UploadMap will get map from remote
func UploadMap(fields xmap.M, filekey, filename, format string, args ...interface{}) (data xmap.M, err error) {
	return Shared.UploadMap(fields, filekey, filename, format, args...)
}

//Download will download file to save path
func Download(saveto string, format string, args ...interface{}) (saved int64, err error) {
	return Shared.Download(saveto, format, args...)
}

//DownloadHeader will download file to save path
func DownloadHeader(saveto string, header xmap.M, format string, args ...interface{}) (saved int64, err error) {
	return Shared.DownloadHeader(saveto, header, format, args...)
}

//RawRequestF is raw request func define
type RawRequestF func(method, uri string, header xmap.M, body io.Reader) (req *http.Request, res *http.Response, err error)

//Client is http get client
type Client struct {
	Raw RawRequestF
}

//NewRawClient will return new client
func NewRawClient(raw RawRequestF) (client *Client) {
	client = &Client{Raw: raw}
	return
}

//NewClient will return new client by http.Client
func NewClient(raw *http.Client) (client *Client) {
	c := &RawClient{C: raw}
	client = &Client{
		Raw: c.RawRequest,
	}
	return
}

//GetBytes will do http request and read the bytes response
func (c *Client) GetBytes(format string, args ...interface{}) (data []byte, err error) {
	data, _, err = c.GetHeaderBytes(nil, format, args...)
	return
}

//GetHeaderBytes will do http request and read the text response
func (c *Client) GetHeaderBytes(header xmap.M, format string, args ...interface{}) (data []byte, res *http.Response, err error) {
	remote := fmt.Sprintf(format, args...)
	_, res, err = c.Raw("GET", remote, header, nil)
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
func (c *Client) GetText(format string, args ...interface{}) (data string, err error) {
	data, _, err = c.GetHeaderText(nil, format, args...)
	return
}

//GetHeaderText will do http request and read the text response
func (c *Client) GetHeaderText(header xmap.M, format string, args ...interface{}) (data string, res *http.Response, err error) {
	bys, res, err := c.GetHeaderBytes(header, format, args...)
	data = string(bys)
	return
}

//GetMap will do http request, read reponse and parse to map
func (c *Client) GetMap(format string, args ...interface{}) (data xmap.M, err error) {
	data, _, err = c.GetHeaderMap(nil, format, args...)
	return
}

//GetHeaderMap will do http request, read reponse and parse to map
func (c *Client) GetHeaderMap(header xmap.M, format string, args ...interface{}) (data xmap.M, res *http.Response, err error) {
	text, res, err := c.GetHeaderBytes(header, format, args...)
	if len(text) > 0 {
		if err == nil {
			data, err = xmap.MapVal(text)
		} else {
			data, _ = xmap.MapVal(text)
		}
	}
	return
}

//GetJSON will do http request, read reponse and parse to json
func (c *Client) GetJSON(value interface{}, format string, args ...interface{}) (err error) {
	_, err = c.GetHeaderJSON(value, nil, format, args...)
	return
}

//GetHeaderMap will do http request, read reponse and parse to map
func (c *Client) GetHeaderJSON(value interface{}, header xmap.M, format string, args ...interface{}) (res *http.Response, err error) {
	text, res, err := c.GetHeaderBytes(header, format, args...)
	if err == nil {
		err = json.Unmarshal(text, value)
	}
	return
}

//PostBytes will do http request and read the bytes response
func (c *Client) PostBytes(body io.Reader, format string, args ...interface{}) (data []byte, err error) {
	data, _, err = c.PostHeaderBytes(nil, body, format, args...)
	return
}

//PostTypeBytes will do http request and read the bytes response
func (c *Client) PostTypeBytes(contentType string, body io.Reader, format string, args ...interface{}) (data []byte, err error) {
	data, _, err = c.PostHeaderBytes(xmap.M{"Content-Type": contentType}, body, format, args...)
	return
}

//PostHeaderBytes will do http request and read the text response
func (c *Client) PostHeaderBytes(header xmap.M, body io.Reader, format string, args ...interface{}) (data []byte, res *http.Response, err error) {
	remote := fmt.Sprintf(format, args...)
	_, res, err = c.Raw("POST", remote, header, body)
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
func (c *Client) PostText(body io.Reader, format string, args ...interface{}) (data string, err error) {
	data, _, err = c.PostHeaderText(nil, body, format, args...)
	return
}

//PostTypeText will do http request and read the text response
func (c *Client) PostTypeText(contentType string, body io.Reader, format string, args ...interface{}) (data string, err error) {
	data, _, err = c.PostHeaderText(xmap.M{"Content-Type": contentType}, body, format, args...)
	return
}

//PostHeaderText will do http request and read the text response
func (c *Client) PostHeaderText(header xmap.M, body io.Reader, format string, args ...interface{}) (data string, res *http.Response, err error) {
	bys, res, err := c.PostHeaderBytes(header, body, format, args...)
	data = string(bys)
	return
}

//PostMap will do http request, read reponse and parse to map
func (c *Client) PostMap(body io.Reader, format string, args ...interface{}) (data xmap.M, err error) {
	data, _, err = c.PostHeaderMap(nil, body, format, args...)
	return
}

//PostTypeMap will do http request, read reponse and parse to map
func (c *Client) PostTypeMap(contentType string, body io.Reader, format string, args ...interface{}) (data xmap.M, err error) {
	data, _, err = c.PostHeaderMap(xmap.M{"Content-Type": contentType}, body, format, args...)
	return
}

//PostHeaderMap will do http request, read reponse and parse to map
func (c *Client) PostHeaderMap(header xmap.M, body io.Reader, format string, args ...interface{}) (data xmap.M, res *http.Response, err error) {
	text, res, err := c.PostHeaderBytes(header, body, format, args...)
	if len(text) > 0 {
		if err == nil {
			data, err = xmap.MapVal(text)
		} else {
			data, _ = xmap.MapVal(text)
		}
	}
	return
}

//PostJSONMap will do http request, read reponse and parse to map
func (c *Client) PostJSONMap(body interface{}, format string, args ...interface{}) (data xmap.M, err error) {
	bys, err := json.Marshal(body)
	if err != nil {
		return
	}
	data, _, err = c.PostHeaderMap(xmap.M{"Content-Type": ContentTypeJSON}, bytes.NewBuffer(bys), format, args...)
	return
}

//PostMap will do http request, read reJSONponse and parse by json
func (c *Client) PostJSON(result interface{}, body io.Reader, format string, args ...interface{}) (err error) {
	data, _, err := c.PostHeaderBytes(nil, body, format, args...)
	if err == nil {
		err = json.Unmarshal(data, result)
	}
	return
}

//PostTypeJSON will do http request, read reponse and parse by json
func (c *Client) PostTypeJSON(result interface{}, contentType string, body io.Reader, format string, args ...interface{}) (err error) {
	data, _, err := c.PostHeaderBytes(xmap.M{"Content-Type": contentType}, body, format, args...)
	if err == nil {
		err = json.Unmarshal(data, result)
	}
	return
}

//PostHeaderJSON will do http request, read reponse and parse by json
func (c *Client) PostHeaderJSON(result interface{}, header xmap.M, body io.Reader, format string, args ...interface{}) (res *http.Response, err error) {
	data, res, err := c.PostHeaderBytes(header, body, format, args...)
	if err == nil {
		err = json.Unmarshal(data, result)
	}
	return
}

//PostJSONJSON will do http request, read reponse and parse by json
func (c *Client) PostJSONJSON(result interface{}, body interface{}, format string, args ...interface{}) (err error) {
	bys, err := json.Marshal(body)
	if err != nil {
		return
	}
	data, _, err := c.PostHeaderBytes(xmap.M{"Content-Type": ContentTypeJSON}, bytes.NewBuffer(bys), format, args...)
	if err == nil {
		err = json.Unmarshal(data, result)
	}
	return
}

//MethodBytes will do http request, read reponse and parse to bytes
func (c *Client) MethodBytes(method string, header xmap.M, body io.Reader, format string, args ...interface{}) (data []byte, res *http.Response, err error) {
	remote := fmt.Sprintf(format, args...)
	_, res, err = c.Raw(method, remote, header, body)
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

//MethodBytes will do http request, read reponse and parse to string
func (c *Client) MethodText(method string, header xmap.M, body io.Reader, format string, args ...interface{}) (data string, res *http.Response, err error) {
	bytes, res, err := c.MethodBytes(method, header, body, format, args...)
	if len(bytes) > 0 {
		data = string(bytes)
	}
	return
}

//MethodBytes will do http request, read reponse and parse to map
func (c *Client) MethodMap(method string, header xmap.M, body io.Reader, format string, args ...interface{}) (data xmap.M, res *http.Response, err error) {
	text, res, err := c.MethodBytes(method, header, body, format, args...)
	if len(text) > 0 {
		if err == nil {
			data, err = xmap.MapVal(text)
		} else {
			data, _ = xmap.MapVal(text)
		}
	}
	return
}

//PostXMLText will do http request, read reponse and parse to map
func (c *Client) PostXMLText(v interface{}, format string, args ...interface{}) (data string, err error) {
	bys, err := xml.Marshal(v)
	if err == nil {
		data, _, err = c.PostHeaderText(xmap.M{"Content-Type": ContentTypeXML}, bytes.NewBuffer(bys), format, args...)
	}
	return
}

//PostFormText will do http request, read reponse and parse to map
func (c *Client) PostFormText(form xmap.M, format string, args ...interface{}) (data string, err error) {
	query := url.Values{}
	for k, v := range form {
		query.Set(k, fmt.Sprintf("%v", v))
	}
	buf := bytes.NewBufferString(query.Encode())
	data, _, err = c.PostHeaderText(xmap.M{"Content-Type": ContentTypeForm}, buf, format, args...)
	return
}

//PostFormMap will do http request, read reponse and parse to map
func (c *Client) PostFormMap(form xmap.M, format string, args ...interface{}) (data xmap.M, err error) {
	query := url.Values{}
	for k, v := range form {
		query.Set(k, fmt.Sprintf("%v", v))
	}
	buf := bytes.NewBufferString(query.Encode())
	data, _, err = c.PostHeaderMap(xmap.M{"Content-Type": ContentTypeForm}, buf, format, args...)
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
func (c *Client) PostMultipartBytes(header, fields xmap.M, format string, args ...interface{}) (data []byte, res *http.Response, err error) {
	bodyBuf, ctype := CreateMultipartBody(fields)
	remote := fmt.Sprintf(format, args...)
	if header == nil {
		header = xmap.M{}
	}
	header.SetValue("Content-Type", ctype)
	_, res, err = c.Raw("POST", remote, header, bodyBuf)
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
func (c *Client) PostMultipartText(header, fields xmap.M, format string, args ...interface{}) (data string, err error) {
	bytes, _, err := c.PostMultipartBytes(header, fields, format, args...)
	if len(bytes) > 0 {
		data = string(bytes)
	}
	return
}

//PostMultipartMap will do http request, read reponse and parse to map
func (c *Client) PostMultipartMap(header, fields xmap.M, format string, args ...interface{}) (data xmap.M, err error) {
	bytes, _, err := c.PostMultipartBytes(header, fields, format, args...)
	if len(bytes) > 0 {
		if err == nil {
			data, err = xmap.MapVal(bytes)
		} else {
			data, _ = xmap.MapVal(bytes)
		}
	}
	return
}

//FileBodyTask is upload task
type FileBodyTask struct {
	reader *io.PipeReader
	writer *io.PipeWriter
}

//CreateFileBodyTask will create the file upload task
func CreateFileBodyTask(fields xmap.M, filekey string, filename string) (*FileBodyTask, io.Reader, string) {
	task := &FileBodyTask{}
	reader, ctype := task.Start(fields, filekey, filename)
	return task, reader, ctype
}

//Start will start the file upload body
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

//Close the upload file body
func (f *FileBodyTask) Close() error {
	if f.reader != nil {
		f.reader.Close()
	}
	if f.writer != nil {
		f.writer.Close()
	}
	return nil
}

//UploadText upload file and get text response
func (c *Client) UploadText(fields xmap.M, filekey, filename, format string, args ...interface{}) (text string, err error) {
	text, _, err = c.UploadHeaderText(nil, fields, filekey, filename, format, args...)
	return
}

//UploadHeaderText upload file and get text response
func (c *Client) UploadHeaderText(header xmap.M, fields xmap.M, filekey, filename, format string, args ...interface{}) (text string, res *http.Response, err error) {
	var ctype string
	var bodyBuf io.Reader
	var task *FileBodyTask
	task, bodyBuf, ctype = CreateFileBodyTask(fields, filekey, filename)
	defer task.Close()
	remote := fmt.Sprintf(format, args...)
	if header == nil {
		header = xmap.M{}
	}
	header.SetValue("Content-Type", ctype)
	_, res, err = c.Raw("POST", remote, header, bodyBuf)
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

//UploadMap will upload file and get map response
func (c *Client) UploadMap(fields xmap.M, filekey, filename, format string, args ...interface{}) (data xmap.M, err error) {
	text, _, err := c.UploadHeaderText(nil, fields, filekey, filename, format, args...)
	if len(text) > 0 {
		if err == nil {
			data, err = xmap.MapVal(text)
		} else {
			data, _ = xmap.MapVal(text)
		}
	}
	return
}

//Download will download the file to save path
func (c *Client) Download(saveto string, format string, args ...interface{}) (saved int64, err error) {
	saved, err = c.DownloadHeader(saveto, nil, format, args...)
	return
}

//DownloadHeader will download the file to save path
func (c *Client) DownloadHeader(saveto string, header xmap.M, format string, args ...interface{}) (saved int64, err error) {
	req, res, err := c.Raw("GET", fmt.Sprintf(format, args...), header, nil)
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
