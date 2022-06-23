package attrscan

import (
	"reflect"
	"strings"
	"testing"

	"github.com/codingeasygo/util/converter"
)

type IsNilArray []string

type IsZeroArray []string

func (s IsNilArray) IsNil() bool {
	return s == nil
}

func (s IsZeroArray) IsZero() bool {
	return len(s) < 1
}

type Simple struct {
	_  string      `json:"-"`
	A0 int64       `json:"a0"`
	A1 *int64      `json:"a1"`
	AX int64       `json:"ax"`
	B0 string      `json:"b0"`
	B1 *string     `json:"b1"`
	BX string      `json:"bx"`
	C0 float64     `json:"c0"`
	C1 *float64    `json:"c1"`
	CX float64     `json:"cx" filter:"#all"`
	D0 IsNilArray  `json:"d0"`
	D1 IsZeroArray `json:"d1"`
	XX string      `json:"-"`
}

func TestFilterField(t *testing.T) {
	var fields []string
	simple := &Simple{
		A0: 100,
		A1: converter.Int64Ptr(100),
		B0: "100",
		B1: converter.StringPtr("100"),
		C0: 100,
		C1: converter.Float64Ptr(100),
	}
	{
		if !CheckValue(reflect.ValueOf(nil), true, true) {
			t.Error("error")
			return
		}
		var intPtr *int
		if CheckValue(reflect.ValueOf(intPtr), false, false) {
			t.Error("error")
			return
		}
	}
	{
		fields = nil
		FilterFieldCall("test", simple, "", func(fieldName, fieldFunc string, field reflect.StructField, value interface{}) {
			fields = append(fields, fieldName)
		})
		if strings.Join(fields, ",") != "a0,a1,b0,b1,c0,c1,cx" {
			t.Errorf("%v", fields)
			return
		}
	}
	{
		fields = nil
		FilterFieldCall("test", simple, "#all", func(fieldName, fieldFunc string, field reflect.StructField, value interface{}) {
			fields = append(fields, fieldName)
		})
		if strings.Join(fields, ",") != "a0,a1,ax,b0,b1,bx,c0,c1,cx,d0,d1" {
			t.Errorf("%v", fields)
			return
		}
	}
	{
		fields = nil
		FilterFieldCall("test", simple, "a0,a1,ax#all", func(fieldName, fieldFunc string, field reflect.StructField, value interface{}) {
			fields = append(fields, fieldName)
		})
		if strings.Join(fields, ",") != "a0,a1,ax" {
			t.Errorf("%v", fields)
			return
		}
	}
	{
		fields = nil
		FilterFieldCall("test", simple, "^a0,a1,ax#all", func(fieldName, fieldFunc string, field reflect.StructField, value interface{}) {
			fields = append(fields, fieldName)
		})
		if strings.Join(fields, ",") != "b0,b1,bx,c0,c1,cx,d0,d1" {
			t.Errorf("%v", fields)
			return
		}
	}
	{
		fields = nil
		FilterFieldCall("test", simple, "count(a0),count(a1),ax#all", func(fieldName, fieldFunc string, field reflect.StructField, value interface{}) {
			fields = append(fields, fieldName)
		})
		if strings.Join(fields, ",") != "a0,a1,ax" {
			t.Errorf("%v", fields)
			return
		}
	}
}
