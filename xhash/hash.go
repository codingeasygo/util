package xhash

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"io"
	"os"

	"github.com/codingeasygo/util/xio"
)

func MD5(data []byte) string {
	h := md5.New()
	h.Write(data)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func SHA1(data []byte) string {
	h := sha1.New()
	h.Write(data)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func ReaderHash(reader io.Reader, enableSHA1, enableMD5 bool) (shah, md5h string, filesize int64, err error) {
	h1 := sha1.New()
	writer := []io.Writer{}
	if enableSHA1 {
		writer = append(writer, h1)
	}
	h2 := md5.New()
	if enableMD5 {
		writer = append(writer, h2)
	}
	filesize, err = xio.CopyMulti(writer, reader)
	if enableSHA1 {
		shah = fmt.Sprintf("%x", h1.Sum(nil))
	}
	if enableMD5 {
		md5h = fmt.Sprintf("%x", h2.Sum(nil))
	}
	return
}

func FileHash(filename string, enableSHA1, enableMD5 bool) (shah, md5h string, filesize int64, err error) {
	f, err := os.Open(filename)
	if err == nil {
		shah, md5h, filesize, err = ReaderHash(f, enableSHA1, enableMD5)
		f.Close()
	}
	return
}
