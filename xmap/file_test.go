package xmap

import "testing"

func TestReadWrite(t *testing.T) {
	data := M{"abc": 1}
	err := WriteJSON(data, "/tmp/test_xmap.json")
	if err != nil {
		t.Error(err)
		return
	}
	readData, err := ReadJSON("/tmp/test_xmap.json")
	if err != nil || readData.Int("abc") != data.Int("abc") {
		t.Error(err)
		return
	}
}
