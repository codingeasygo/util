package converter

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var ErrNil = fmt.Errorf("nil value")

func Int(v interface{}) (val int) {
	val, _ = IntVal(v)
	return
}

func IntVal(v interface{}) (int, error) {
	if v == nil {
		return 0, ErrNil
	}
	k := reflect.TypeOf(v)
	if k.Name() == "Time" {
		t := v.(time.Time)
		return int(t.Local().UnixNano() / 1e6), nil
	}
	switch k.Kind() {
	case reflect.Int:
		return int(v.(int)), nil
	case reflect.Int8:
		return int(v.(int8)), nil
	case reflect.Int16:
		return int(v.(int16)), nil
	case reflect.Int32:
		return int(v.(int32)), nil
	case reflect.Int64:
		return int(v.(int64)), nil
	case reflect.Uint:
		return int(v.(uint)), nil
	case reflect.Uint8:
		return int(v.(uint8)), nil
	case reflect.Uint16:
		return int(v.(uint16)), nil
	case reflect.Uint32:
		return int(v.(uint32)), nil
	case reflect.Uint64:
		return int(v.(uint64)), nil
	case reflect.Float32:
		return int(v.(float32)), nil
	case reflect.Float64:
		return int(v.(float64)), nil
	case reflect.String:
		fv, err := strconv.ParseInt(strings.TrimSpace(v.(string)), 10, 64)
		return int(fv), err
	default:
		return 0, fmt.Errorf("incompactable kind(%v)", k.Kind())
	}
}

func Int64(v interface{}) int64 {
	val, _ := Int64Val(v)
	return val
}

func Int64Val(v interface{}) (int64, error) {
	if v == nil {
		return 0, ErrNil
	}
	k := reflect.TypeOf(v)
	if k.Name() == "Time" {
		t := v.(time.Time)
		return int64(t.Local().UnixNano() / 1e6), nil
	}
	switch k.Kind() {
	case reflect.Int:
		return int64(v.(int)), nil
	case reflect.Int8:
		return int64(v.(int8)), nil
	case reflect.Int16:
		return int64(v.(int16)), nil
	case reflect.Int32:
		return int64(v.(int32)), nil
	case reflect.Int64:
		return int64(v.(int64)), nil
	case reflect.Uint:
		return int64(v.(uint)), nil
	case reflect.Uint8:
		return int64(v.(uint8)), nil
	case reflect.Uint16:
		return int64(v.(uint16)), nil
	case reflect.Uint32:
		return int64(v.(uint32)), nil
	case reflect.Uint64:
		return int64(v.(uint64)), nil
	case reflect.Float32:
		return int64(v.(float32)), nil
	case reflect.Float64:
		return int64(v.(float64)), nil
	case reflect.String:
		fv, err := strconv.ParseInt(strings.TrimSpace(v.(string)), 10, 64)
		return int64(fv), err
	default:
		return 0, fmt.Errorf("incompactable kind(%v)", k.Kind())
	}
}

func Uint64(v interface{}) uint64 {
	val, _ := Uint64Val(v)
	return val
}

func Uint64Val(v interface{}) (uint64, error) {
	if v == nil {
		return 0, ErrNil
	}
	k := reflect.TypeOf(v)
	if k.Name() == "Time" {
		t := v.(time.Time)
		return uint64(t.Local().UnixNano() / 1e6), nil
	}
	switch k.Kind() {
	case reflect.Int:
		return uint64(v.(int)), nil
	case reflect.Int8:
		return uint64(v.(int8)), nil
	case reflect.Int16:
		return uint64(v.(int16)), nil
	case reflect.Int32:
		return uint64(v.(int32)), nil
	case reflect.Int64:
		return uint64(v.(int64)), nil
	case reflect.Uint:
		return uint64(v.(uint)), nil
	case reflect.Uint8:
		return uint64(v.(uint8)), nil
	case reflect.Uint16:
		return uint64(v.(uint16)), nil
	case reflect.Uint32:
		return uint64(v.(uint32)), nil
	case reflect.Uint64:
		return uint64(v.(uint64)), nil
	case reflect.Float32:
		return uint64(v.(float32)), nil
	case reflect.Float64:
		return uint64(v.(float64)), nil
	case reflect.String:
		fv, err := strconv.ParseInt(strings.TrimSpace(v.(string)), 10, 64)
		return uint64(fv), err
	default:
		return 0, fmt.Errorf("incompactable kind(%v)", k.Kind())
	}
}

func Float64(v interface{}) float64 {
	val, _ := Float64Val(v)
	return val
}

func Float64Val(v interface{}) (float64, error) {
	if v == nil {
		return 0, ErrNil
	}
	k := reflect.TypeOf(v)
	if k.Name() == "Time" {
		return float64(v.(time.Time).Local().UnixNano() / 1e6), nil
	}
	switch k.Kind() {
	case reflect.Int:
		return float64(v.(int)), nil
	case reflect.Int8:
		return float64(v.(int8)), nil
	case reflect.Int16:
		return float64(v.(int16)), nil
	case reflect.Int32:
		return float64(v.(int32)), nil
	case reflect.Int64:
		return float64(v.(int64)), nil
	case reflect.Uint:
		return float64(v.(uint)), nil
	case reflect.Uint8:
		return float64(v.(uint8)), nil
	case reflect.Uint16:
		return float64(v.(uint16)), nil
	case reflect.Uint32:
		return float64(v.(uint32)), nil
	case reflect.Uint64:
		return float64(v.(uint64)), nil
	case reflect.Float32:
		return float64(v.(float32)), nil
	case reflect.Float64:
		return float64(v.(float64)), nil
	case reflect.String:
		fv, err := strconv.ParseFloat(strings.TrimSpace(v.(string)), 10)
		return float64(fv), err
	default:
		return 0, fmt.Errorf("incompactable kind(%v)", k.Kind())
	}
}

func String(v interface{}) string {
	val, _ := StringVal(v)
	return val
}

func StringVal(v interface{}) (res string, err error) {
	if v == nil {
		return "", ErrNil
	}
	switch reflect.TypeOf(v).Kind() {
	case reflect.String:
		return v.(string), nil
	case reflect.Slice:
		if bys, ok := v.([]byte); ok {
			return string(bys), nil
		}
		fallthrough
	default:
		return fmt.Sprintf("%v", v), nil
	}
}

func ArrayVal(v interface{}) ([]interface{}, error) {
	if v == nil {
		return nil, ErrNil
	}
	if vals, ok := v.([]interface{}); ok {
		return vals, nil
	}
	if sval, ok := v.(string); ok {
		vals := []interface{}{}
		for _, val := range strings.Split(sval, ",") {
			vals = append(vals, val)
		}
		return vals, nil
	}
	vals := reflect.ValueOf(v)
	if vals.Kind() != reflect.Slice {
		return nil, fmt.Errorf("incompactable kind(%v)", vals.Kind())
	}
	var vs = []interface{}{}
	for i := 0; i < vals.Len(); i++ {
		vs = append(vs, vals.Index(i).Interface())
	}
	return vs, nil
}

func ArrayStringVal(v interface{}) (svals []string, err error) {
	if v == nil {
		return nil, ErrNil
	}
	var sval string
	if vals, ok := v.([]interface{}); ok {
		for _, v := range vals {
			sval, err = StringVal(v)
			if err != nil {
				return
			}
			svals = append(svals, sval)
		}
		return
	}
	if sval, ok := v.(string); ok {
		svals = strings.Split(sval, ",")
		return
	}
	vals := reflect.ValueOf(v)
	if vals.Kind() != reflect.Slice {
		err = fmt.Errorf("incompactable kind(%v)", vals.Kind())
		return
	}
	for i := 0; i < vals.Len(); i++ {
		if vals.Index(i).IsZero() {
			sval, err = StringVal(nil)
		} else {
			sval, err = StringVal(vals.Index(i).Interface())
		}
		if err != nil {
			return
		}
		svals = append(svals, sval)
	}
	return
}

func ArrayIntVal(v interface{}) (ivals []int, err error) {
	if v == nil {
		return nil, ErrNil
	}
	var ival int
	if vals, ok := v.([]interface{}); ok {
		for _, v := range vals {
			ival, err = IntVal(v)
			if err != nil {
				return
			}
			ivals = append(ivals, ival)
		}
		return
	}
	if sval, ok := v.(string); ok {
		for _, val := range strings.Split(sval, ",") {
			ival, err = IntVal(val)
			if err != nil {
				return
			}
			ivals = append(ivals, ival)
		}
		return
	}
	vals := reflect.ValueOf(v)
	if vals.Kind() != reflect.Slice {
		err = fmt.Errorf("incompactable kind(%v)", vals.Kind())
		return
	}
	for i := 0; i < vals.Len(); i++ {
		ival, err = IntVal(vals.Index(i).Interface())
		if err != nil {
			return
		}
		ivals = append(ivals, ival)
	}
	return
}

func ArrayInt64Val(v interface{}) (ivals []int64, err error) {
	if v == nil {
		return nil, ErrNil
	}
	var ival int64
	if vals, ok := v.([]interface{}); ok {
		for _, v := range vals {
			ival, err = Int64Val(v)
			if err != nil {
				return
			}
			ivals = append(ivals, ival)
		}
		return
	}
	if sval, ok := v.(string); ok {
		for _, val := range strings.Split(sval, ",") {
			ival, err = Int64Val(val)
			if err != nil {
				return
			}
			ivals = append(ivals, ival)
		}
		return
	}
	vals := reflect.ValueOf(v)
	if vals.Kind() != reflect.Slice {
		err = fmt.Errorf("incompactable kind(%v)", vals.Kind())
		return
	}
	for i := 0; i < vals.Len(); i++ {
		ival, err = Int64Val(vals.Index(i).Interface())
		if err != nil {
			return
		}
		ivals = append(ivals, ival)
	}
	return
}

func ArrayUint64Val(v interface{}) (ivals []uint64, err error) {
	if v == nil {
		return nil, ErrNil
	}
	var ival uint64
	if vals, ok := v.([]interface{}); ok {
		for _, v := range vals {
			ival, err = Uint64Val(v)
			if err != nil {
				return
			}
			ivals = append(ivals, ival)
		}
		return
	}
	if sval, ok := v.(string); ok {
		for _, val := range strings.Split(sval, ",") {
			ival, err = Uint64Val(val)
			if err != nil {
				return
			}
			ivals = append(ivals, ival)
		}
		return
	}
	vals := reflect.ValueOf(v)
	if vals.Kind() != reflect.Slice {
		err = fmt.Errorf("incompactable kind(%v)", vals.Kind())
		return
	}
	for i := 0; i < vals.Len(); i++ {
		ival, err = Uint64Val(vals.Index(i).Interface())
		if err != nil {
			return
		}
		ivals = append(ivals, ival)
	}
	return
}

func ArrayFloat64Val(v interface{}) (ivals []float64, err error) {
	if v == nil {
		return nil, ErrNil
	}
	var ival float64
	if vals, ok := v.([]interface{}); ok {
		for _, v := range vals {
			ival, err = Float64Val(v)
			if err != nil {
				return
			}
			ivals = append(ivals, ival)
		}
		return
	}
	if sval, ok := v.(string); ok {
		for _, val := range strings.Split(sval, ",") {
			ival, err = Float64Val(val)
			if err != nil {
				return
			}
			ivals = append(ivals, ival)
		}
		return
	}
	vals := reflect.ValueOf(v)
	if vals.Kind() != reflect.Slice {
		err = fmt.Errorf("incompactable kind(%v)", vals.Kind())
		return
	}
	for i := 0; i < vals.Len(); i++ {
		ival, err = Float64Val(vals.Index(i).Interface())
		if err != nil {
			return
		}
		ivals = append(ivals, ival)
	}
	return
}

//ArrayHaving will return true if the array element having one is in objs
func ArrayHaving(ary interface{}, objs ...interface{}) bool {
	switch reflect.TypeOf(ary).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(ary)
		for i := 0; i < s.Len(); i++ {
			for _, obj := range objs {
				if obj == s.Index(i).Interface() {
					return true
				}
			}
		}
		return false
	default:
		return false
	}
}

func JSON(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		return err.Error()
	}
	return string(data)
}

//UnmarshalJSON will read bytes from reader and unmarshal to object
func UnmarshalJSON(r io.Reader, v interface{}) (data []byte, err error) {
	data, err = ioutil.ReadAll(r)
	if err == nil || err == io.EOF {
		err = json.Unmarshal(data, v)
	}
	return
}

//UnmarshalXML will read bytes from reader and unmarshal to object
func UnmarshalXML(r io.Reader, v interface{}) (data []byte, err error) {
	data, err = ioutil.ReadAll(r)
	if err == nil || err == io.EOF {
		err = xml.Unmarshal(data, v)
	}
	return
}

// Int8Ptr -
func Int8Ptr(arg int8) *int8 {
	return &arg
}

// Uint8Ptr -
func Uint8Ptr(arg uint8) *uint8 {
	return &arg
}

// Int16Ptr -
func Int16Ptr(arg int16) *int16 {
	return &arg
}

// Uint16Ptr -
func Uint16Ptr(arg uint16) *uint16 {
	return &arg
}

// IntPtr -
func IntPtr(arg int) *int {
	return &arg
}

// UintPtr -
func UintPtr(arg uint) *uint {
	return &arg
}

// Int32Ptr -
func Int32Ptr(arg int32) *int32 {
	return &arg
}

// Uint32Ptr -
func Uint32Ptr(arg uint32) *uint32 {
	return &arg
}

// Int64Ptr -
func Int64Ptr(arg int64) *int64 {
	return &arg
}

// Uint64Ptr -
func Uint64Ptr(arg uint64) *uint64 {
	return &arg
}

// Float32Ptr -
func Float32Ptr(arg float32) *float32 {
	return &arg
}

// Float64Ptr -
func Float64Ptr(arg float64) *float64 {
	return &arg
}

// StringPtr -
func StringPtr(arg string) *string {
	return &arg
}

//Join all slice to string
func Join(v interface{}, sep string) string {
	vtype := reflect.TypeOf(v)
	if vtype.Kind() != reflect.Slice {
		panic("not slice")
	}
	vval := reflect.ValueOf(v)
	if vval.Len() < 1 {
		return ""
	}
	val := fmt.Sprintf("%v", reflect.Indirect(vval.Index(0)).Interface())
	for i := 1; i < vval.Len(); i++ {
		val += fmt.Sprintf("%v%v", sep, reflect.Indirect(vval.Index(i)).Interface())
	}
	return val
}

func IndirectString(val interface{}) string {
	if val == nil {
		return "nil"
	}
	rval := reflect.ValueOf(val)
	if rval.Kind() == reflect.Ptr && rval.IsNil() {
		return "nil"
	}
	return fmt.Sprintf("%v", reflect.Indirect(rval).Interface())
}
