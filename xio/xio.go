package xio

import (
	"encoding/json"
	"io"
)

//WriteJSON will marshal value to string and write to io.Writer
func WriteJSON(w io.Writer, v interface{}) (n int, err error) {
	data, err := json.Marshal(v)
	if err == nil {
		n, err = w.Write(data)
	}
	return
}
