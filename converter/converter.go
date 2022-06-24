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

func valueConvert(v interface{}, targetType reflect.Type, defaultRet interface{}, parse func(string) (interface{}, error)) (result interface{}, err error) {
	targetValue := reflect.ValueOf(v)
	if !targetValue.IsValid() {
		return defaultRet, ErrNil
	}
	if v, ok := v.(time.Time); ok {
		targetValue = reflect.ValueOf(v.Local().UnixNano() / 1e6)
	}
	if targetValue.Kind() == reflect.String {
		result, err = parse(targetValue.String())
		if err != nil {
			return defaultRet, err
		}
		targetValue = reflect.ValueOf(result)
	}
	if targetValue.CanConvert(targetType) {
		return targetValue.Convert(targetType).Interface(), nil
	}
	if targetValue.Kind() == reflect.Ptr {
		targetValue = reflect.Indirect(targetValue)
		if !targetValue.IsValid() {
			return defaultRet, ErrNil
		}
		result, err = valueConvert(targetValue.Interface(), targetType, defaultRet, parse)
		return
	}
	return defaultRet, fmt.Errorf("incompactable kind(%v)", targetValue.Kind())
}

func Int(v interface{}) (val int) {
	val, _ = IntVal(v)
	return
}

var _intType = reflect.TypeOf(0)

func IntVal(v interface{}) (int, error) {
	ret, err := valueConvert(v, _intType, 0, func(s string) (interface{}, error) {
		return strconv.ParseInt(strings.TrimSpace(s), 10, 64)
	})
	if err != nil {
		return 0, err
	}
	return ret.(int), nil
}

func Int64(v interface{}) int64 {
	val, _ := Int64Val(v)
	return val
}

var _int64Type = reflect.TypeOf(int64(0))

func Int64Val(v interface{}) (int64, error) {
	ret, err := valueConvert(v, _int64Type, 0, func(s string) (interface{}, error) {
		return strconv.ParseInt(strings.TrimSpace(s), 10, 64)
	})
	if err != nil {
		return 0, err
	}
	return ret.(int64), nil
}

func Uint64(v interface{}) uint64 {
	val, _ := Uint64Val(v)
	return val
}

var _uint64Type = reflect.TypeOf(uint64(0))

func Uint64Val(v interface{}) (uint64, error) {
	ret, err := valueConvert(v, _uint64Type, 0, func(s string) (interface{}, error) {
		return strconv.ParseInt(strings.TrimSpace(s), 10, 64)
	})
	if err != nil {
		return 0, err
	}
	return ret.(uint64), nil
}

func Float64(v interface{}) float64 {
	val, _ := Float64Val(v)
	return val
}

var _float64Type = reflect.TypeOf(float64(0))

func Float64Val(v interface{}) (float64, error) {
	ret, err := valueConvert(v, _float64Type, 0, func(s string) (interface{}, error) {
		return strconv.ParseFloat(strings.TrimSpace(s), 64)
	})
	if err != nil {
		return 0, err
	}
	return ret.(float64), nil
}

func String(v interface{}) string {
	val, _ := StringVal(v)
	return val
}

func StringVal(v interface{}) (res string, err error) {
	if v == nil {
		return "", ErrNil
	}
	switch v := v.(type) {
	case string:
		return v, nil
	case *string:
		return *v, nil
	case []byte:
		return string(v), nil
	case *[]byte:
		return string(*v), nil
	default:
		return fmt.Sprintf("%v", v), nil
	}
}

//ArrayVal will convert value to array, if v is string will split it by comma, if v is slice will loop element to array, other will error
func ArrayVal(v interface{}) ([]interface{}, error) {
	return ArrayValAll(v, false)
}

//ArrayValAll will convert all value to array, if v is string will split it by comma, if v is slice will loop element to array, other will return []interface{}{v} when all is true
func ArrayValAll(v interface{}, all bool) ([]interface{}, error) {
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
	if sval, ok := v.(*string); ok {
		if sval == nil {
			return nil, ErrNil
		}
		vals := []interface{}{}
		for _, val := range strings.Split(*sval, ",") {
			vals = append(vals, val)
		}
		return vals, nil
	}
	vals := reflect.ValueOf(v)
	if vals.Kind() != reflect.Slice {
		if all {
			return []interface{}{v}, nil
		}
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
	if sval, ok := v.(*string); ok {
		if sval == nil {
			return nil, ErrNil
		}
		svals = strings.Split(*sval, ",")
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
	if sval, ok := v.(*string); ok {
		if sval == nil {
			return nil, ErrNil
		}
		for _, val := range strings.Split(*sval, ",") {
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
	if sval, ok := v.(*string); ok {
		if sval == nil {
			return nil, ErrNil
		}
		for _, val := range strings.Split(*sval, ",") {
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
	if sval, ok := v.(*string); ok {
		if sval == nil {
			return nil, ErrNil
		}
		for _, val := range strings.Split(*sval, ",") {
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
	if sval, ok := v.(*string); ok {
		if sval == nil {
			return nil, ErrNil
		}
		for _, val := range strings.Split(*sval, ",") {
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

func XML(v interface{}) string {
	data, err := xml.Marshal(v)
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

const (
	JoinPolicyNotSliceEmpty  = 1 << 0
	JoinPolicyNotSliceString = 1 << 1
	JoinPolicyNilSkip        = 1 << 5
	JoinPolicyNilString      = 1 << 6
	JoinPolicyDefault        = JoinPolicyNotSliceEmpty | JoinPolicyNilSkip
)

func JoinSafe(v interface{}, sep string, policy int) string {
	vtype := reflect.TypeOf(v)
	if vtype.Kind() != reflect.Slice {
		if policy&JoinPolicyNotSliceString == JoinPolicyNotSliceString {
			return fmt.Sprintf("%v", v)
		} else {
			return ""
		}
	}
	vval := reflect.ValueOf(v)
	if vval.Len() < 1 {
		return ""
	}
	stringVal := func(v reflect.Value) (string, bool) {
		if v.Kind() == reflect.Ptr && v.IsNil() {
			if policy&JoinPolicyNilString == JoinPolicyNilString {
				return fmt.Sprintf("%v", v), true
			} else {
				return "", false
			}
		}
		v = reflect.Indirect(v)
		return fmt.Sprintf("%v", v.Interface()), true
	}
	val := ""
	for i := 0; i < vval.Len(); i++ {
		s, ok := stringVal(vval.Index(i))
		if !ok {
			continue
		}
		if len(val) > 0 {
			val += fmt.Sprintf("%v%v", sep, s)
		} else {
			val = s
		}
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
