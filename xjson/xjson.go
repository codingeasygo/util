package xjson

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
)

//WriteJSON will marshal value to string and write to io.Writer
func WriteJSON(w io.Writer, v interface{}) (n int, err error) {
	data, err := json.Marshal(v)
	if err == nil {
		n, err = w.Write(data)
	}
	return
}

//WriteJSON will read all data from io.Reader and unmashl to value
func ReadJSON(r io.Reader, v interface{}) (n int, err error) {
	data, err := ioutil.ReadAll(r)
	if err == nil {
		err = json.Unmarshal(data, v)
	}
	return
}

//WriteJSONFile will marshal value to string and write to file
func WriteJSONFile(filename string, v interface{}) (err error) {
	file, err := os.OpenFile(filename, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err == nil {
		_, err = WriteJSON(file, v)
		file.Close()
	}
	return
}

//ReadSONFile will read all data from io.Reader and unmashl to value
func ReadSONFile(filename string, v interface{}) (err error) {
	file, err := os.Open(filename)
	if err == nil {
		_, err = ReadJSON(file, v)
		file.Close()
	}
	return
}
