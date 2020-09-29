package xmap

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

//ReadJSON will read json data from file and parset to M
func ReadJSON(filename string) (data M, err error) {
	jsonData, err := ioutil.ReadFile(filename)
	if err == nil {
		err = json.Unmarshal(jsonData, &data)
	}
	return
}

//WriteJSON will marshal M to json and write to file
func WriteJSON(data M, filename string) (err error) {
	dir, _ := filepath.Split(filename)
	os.MkdirAll(dir, os.ModePerm)
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err == nil {
		ioutil.WriteFile(filename, jsonData, os.ModePerm)
	}
	return
}
