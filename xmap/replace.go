package xmap

import (
	"os"
	"regexp"
	"strings"

	"github.com/codingeasygo/util/converter"
)

//ReplaceAll will replace input string by ${xx}, which xx is in values,
//if usingEnv is true, xx will check use env when vals is not having xx,
//if usingEmpty is true, xx will check use empty string when vals is not having xx and env is not exist
func ReplaceAll(vals func(key string) interface{}, in string, usingEnv, usingEmpty bool) (out string) {
	reg := regexp.MustCompile(`\$\{[^\}]*\}`)
	var rval string
	out = reg.ReplaceAllStringFunc(in, func(m string) string {
		keys := strings.Split(strings.Trim(m, "${}\t "), ",")
		for _, key := range keys {
			v := vals(key)
			if v != nil {
				rval = converter.IndirectString(v)
			}
			if len(rval) < 1 && usingEnv {
				rval = os.Getenv(key)
			}
			if len(rval) > 0 {
				break
			}
		}
		if len(rval) > 0 {
			return rval
		}
		if usingEmpty {
			return ""
		}
		return m
	})
	return
}
