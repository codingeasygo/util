package xmap

import (
	"testing"
)

func TestReplace(t *testing.T) {
	vals := Wrap(map[string]interface{}{
		"int":    1,
		"float":  1.0,
		"string": "abc",
	})
	if vals.ReplaceAll(`${int}`, true, true) != "1" {
		t.Error("error")
		return
	}
	if vals.ReplaceAll(`${float}`, true, true) != "1" {
		t.Error("error")
		return
	}
	if vals.ReplaceAll(`${string}`, true, true) != "abc" {
		t.Error("error")
		return
	}
	if vals.ReplaceAll(`${string}-${int}`, true, true) != "abc-1" {
		t.Error("error")
		return
	}
	if vals.ReplaceAll(`${PWD}`, true, true) == "" {
		t.Error("error")
		return
	}
	if vals.ReplaceAll(`${abc}`, true, true) != "" {
		t.Error("error")
		return
	}
	if vals.ReplaceAll(`${abc}`, true, false) != "${abc}" {
		t.Error("error")
		return
	}
}
