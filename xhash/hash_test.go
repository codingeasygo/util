package xhash

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestFileHash(t *testing.T) {
	defer os.Remove("test.tmp")
	ioutil.WriteFile("test.tmp", []byte("abc"), os.ModePerm)
	shah1 := SHA1([]byte("abc"))
	md5h1 := MD5([]byte("abc"))
	shah2, md5h2, _, err := FileHash("test.tmp", true, true)
	if err != nil || shah1 != shah2 || md5h1 != md5h2 {
		t.Error("error")
		return
	}
}
