package xio

import (
	"io/ioutil"
	"testing"
)

func TestSimple(t *testing.T) {
	WriteJSON(ioutil.Discard, map[string]interface{}{})
}
