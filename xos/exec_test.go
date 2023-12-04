package xos

import "testing"

func TestRun(t *testing.T) {
	text, err := Run("echo %v", "abc")
	if err != nil || text != "abc\n" {
		t.Errorf("err:%v,text:%v", err, text)
		return
	}
}
