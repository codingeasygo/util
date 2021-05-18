package xjson

import (
	"os"
	"testing"
)

func TestSimple(t *testing.T) {
	of, err := os.CreateTemp(os.TempDir(), "xjson*.json")
	if err != nil {
		t.Error(err)
		return
	}
	of.Close()
	err = WriteJSONFile(of.Name(), map[string]interface{}{})
	if err != nil {
		t.Error(err)
		return
	}
	err = ReadSONFile(of.Name(), &map[string]interface{}{})
	if err != nil {
		t.Error(err)
		return
	}
}
