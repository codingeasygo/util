package xsql

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/codingeasygo/util/converter"
	"github.com/codingeasygo/util/xmap"
	"github.com/codingeasygo/util/xtime"
)

//Time is database value to parse data from database and parset time.Time to timestamp on json mashal
type Time time.Time

// TimeUnix will return time by timestamp
func TimeUnix(timestamp int64) Time {
	return Time(time.Unix(0, timestamp*1e6))
}

// TimeZero will return zero time
func TimeZero() Time {
	return Time(time.Unix(0, 0*1e6))
}

// TimeNow return current Time
func TimeNow() Time {
	return Time(time.Now())
}

// TimeStartOfToday return 00:00:00 of today
func TimeStartOfToday() Time {
	return Time(xtime.TimeStartOfToday())
}

// TimeStartOfWeek return 00:00:00 of week
func TimeStartOfWeek() Time {
	return Time(xtime.TimeStartOfWeek())
}

// TimeStartOfMonth return 00:00:00 of today
func TimeStartOfMonth() Time {
	return Time(xtime.TimeStartOfMonth())
}

//Timestamp return timestamp
func (t Time) Timestamp() int64 {
	return time.Time(t).Local().UnixNano() / 1e6
}

//MarshalJSON marshal time to string
func (t *Time) MarshalJSON() ([]byte, error) {
	raw := t.Timestamp()
	if raw < 0 {
		return []byte("0"), nil
	}
	stamp := fmt.Sprintf("%v", raw)
	return []byte(stamp), nil
}

//UnmarshalJSON unmarshal string to time
func (t *Time) UnmarshalJSON(bys []byte) (err error) {
	val := strings.TrimSpace(string(bys))
	if val == "null" {
		return
	}
	timestamp, err := strconv.ParseInt(val, 10, 64)
	if err == nil {
		*t = Time(time.Unix(0, timestamp*1e6))
	}
	return
}

//Scan is sql.Sanner
func (t *Time) Scan(src interface{}) (err error) {
	if src != nil {
		if timeSrc, ok := src.(time.Time); ok {
			*t = Time(timeSrc)
		}
	}
	return
}

func (t Time) AsTime() time.Time {
	return time.Time(t)
}

func (t Time) String() string {
	return time.Time(t).String()
}

//M is database value to parse json data to map value
type M map[string]interface{}

//RawMap will return raw map value
func (m M) RawMap() map[string]interface{} {
	return m
}

//Scan is sql.Sanner
func (m *M) Scan(src interface{}) (err error) {
	if src != nil {
		if jsonSrc, ok := src.(string); ok {
			err = json.Unmarshal([]byte(jsonSrc), m)
		} else {
			err = fmt.Errorf("the %v,%v is not string", reflect.TypeOf(src), src)
		}
	}
	return
}

//Value will parse to json value
func (m M) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	bys, err := json.Marshal(m)
	return string(bys), err
}

func (m M) AsMap() xmap.M { return xmap.M(m) }

//MArray is database value to parse json data to map value
type MArray []M

//Scan is sql.Sanner
func (m *MArray) Scan(src interface{}) (err error) {
	if src != nil {
		if jsonSrc, ok := src.(string); ok {
			err = json.Unmarshal([]byte(jsonSrc), m)
		} else {
			err = fmt.Errorf("the %v,%v is not string", reflect.TypeOf(src), src)
		}
	}
	return
}

//Value will parse to json value
func (m MArray) Value() (driver.Value, error) {
	if m == nil {
		return "[]", nil
	}
	bys, err := json.Marshal(m)
	return string(bys), err
}

func sqlScan(src, dst interface{}, strConvert func(str string) (xerr error)) (err error) {
	if src == nil {
		return
	}
	str, ok := src.(string)
	if !ok {
		err = fmt.Errorf("the %v,%v is not string", reflect.TypeOf(src), src)
		return
	}
	if len(str) < 1 || str == "null" {
		return
	}
	if strings.HasPrefix(str, "[") {
		err = json.Unmarshal([]byte(str), dst)
		if err != nil {
			err = fmt.Errorf("unmarshal fail with %v by :%v", err, str)
		}
		return
	}
	if strings.HasPrefix(str, ",") {
		str = strings.TrimSpace(str)
		str = strings.Trim(str, ",")
		if len(str) > 0 {
			err = strConvert(str)
		}
		return
	}
	err = fmt.Errorf("the %v,%v is not invalid format", reflect.TypeOf(src), src)
	return
}

//IntArray is database value to parse data to []int64 value
type IntArray []int

//Scan is sql.Sanner
func (i *IntArray) Scan(src interface{}) (err error) {
	err = sqlScan(src, i, func(str string) (xerr error) {
		*i, xerr = converter.ArrayIntVal(str)
		return
	})
	return
}

//Value is driver.Valuer
func (i IntArray) Value() (driver.Value, error) {
	if i == nil {
		return "[]", nil
	}
	bys, err := json.Marshal(i)
	return string(bys), err
}

func (i IntArray) Len() int {
	return len(i)
}
func (i IntArray) Less(a, b int) bool {
	return i[a] < i[b]
}
func (i IntArray) Swap(a, b int) {
	i[a], i[b] = i[b], i[a]
}

//HavingOne will check if array having one value in vals
func (i IntArray) HavingOne(vals ...int) bool {
	for _, v0 := range i {
		for _, v1 := range vals {
			if v0 == v1 {
				return true
			}
		}
	}
	return false
}

//Join will parset to database array
func (i IntArray) Join(sep string) (res string) {
	for i, v := range i {
		if i > 0 {
			res += sep
		}
		res += fmt.Sprintf("%v", v)
	}
	return
}

//DbArray will join value to database array
func (i IntArray) DbArray() (res string) {
	res = "{" + i.Join(",") + "}"
	return
}

//StrArray will join value to string array by comma
func (i IntArray) StrArray() (res string) {
	res = "," + i.Join(",") + ","
	return
}

//RemoveDuplicate will remove duplicate and copy item to new array
func (i IntArray) RemoveDuplicate() IntArray {
	var arr IntArray
	m := map[int]bool{}
	for _, v := range i {
		if m[v] {
			continue
		}
		m[v] = true
		arr = append(arr, v)
	}
	return arr
}

//AsPtrArray will convet normla to ptr
func (i IntArray) AsPtrArray() (vals IntPtrArray) {
	for _, v := range i {
		vals = append(vals, converter.IntPtr(v))
	}
	return
}

//IntPtrArray is database value to parse data to []int64 value
type IntPtrArray []*int

//Scan is sql.Sanner
func (i *IntPtrArray) Scan(src interface{}) (err error) {
	err = sqlScan(src, i, func(str string) (xerr error) {
		vals, xerr := converter.ArrayIntVal(str)
		if xerr == nil {
			*i = IntArray(vals).AsPtrArray()
		}
		return
	})
	return
}

//Value is driver.Valuer
func (i IntPtrArray) Value() (driver.Value, error) {
	if i == nil {
		return "[]", nil
	}
	bys, err := json.Marshal(i)
	return string(bys), err
}

func (i IntPtrArray) Len() int {
	return len(i)
}
func (i IntPtrArray) Less(a, b int) bool {
	return i[a] == nil || (i[b] != nil && *i[a] < *i[b])
}
func (i IntPtrArray) Swap(a, b int) {
	i[a], i[b] = i[b], i[a]
}

//HavingOne will check if array having one value in vals
func (i IntPtrArray) HavingOne(vals ...int) bool {
	for _, v0 := range i {
		for _, v1 := range vals {
			if v0 != nil && *v0 == v1 {
				return true
			}
		}
	}
	return false
}

//Join will parset to database array
func (i IntPtrArray) Join(sep string) (res string) {
	for i, v := range i {
		if i > 0 {
			res += sep
		}
		if v == nil {
			res += "nil"
		} else {
			res += fmt.Sprintf("%v", *v)
		}
	}
	return
}

//DbArray will join value to database array
func (i IntPtrArray) DbArray() (res string) {
	res = "{" + i.Join(",") + "}"
	return
}

//StrArray will join value to string array by comma
func (i IntPtrArray) StrArray() (res string) {
	res = "," + i.Join(",") + ","
	return
}

//RemoveDuplicate will remove duplicate and copy item to new array
func (i IntPtrArray) RemoveDuplicate() IntPtrArray {
	var arr IntPtrArray
	m := map[int]bool{}
	for _, v := range i {
		if v == nil || m[*v] {
			continue
		}
		m[*v] = true
		arr = append(arr, v)
	}
	return arr
}

//AsArray will convet ptr to normal, skip nil
func (i IntPtrArray) AsArray() (vals IntArray) {
	for _, v := range i {
		if v == nil {
			continue
		}
		vals = append(vals, *v)
	}
	return
}

//Int64Array is database value to parse data to []int64 value
type Int64Array []int64

//Scan is sql.Sanner
func (i *Int64Array) Scan(src interface{}) (err error) {
	err = sqlScan(src, i, func(str string) (xerr error) {
		*i, xerr = converter.ArrayInt64Val(str)
		return
	})
	return
}

//Value is driver.Valuer
func (i Int64Array) Value() (driver.Value, error) {
	if i == nil {
		return "[]", nil
	}
	bys, err := json.Marshal(i)
	return string(bys), err
}

func (i Int64Array) Len() int {
	return len(i)
}
func (i Int64Array) Less(a, b int) bool {
	return i[a] < i[b]
}
func (i Int64Array) Swap(a, b int) {
	i[a], i[b] = i[b], i[a]
}

//HavingOne will check if array having one value in vals
func (i Int64Array) HavingOne(vals ...int64) bool {
	for _, v0 := range i {
		for _, v1 := range vals {
			if v0 == v1 {
				return true
			}
		}
	}
	return false
}

//Join will parset to database array
func (i Int64Array) Join(sep string) (res string) {
	for i, v := range i {
		if i > 0 {
			res += sep
		}
		res += fmt.Sprintf("%v", v)
	}
	return
}

//DbArray will join value to database array
func (i Int64Array) DbArray() (res string) {
	res = "{" + i.Join(",") + "}"
	return
}

//StrArray will join value to string array by comma
func (i Int64Array) StrArray() (res string) {
	res = "," + i.Join(",") + ","
	return
}

//RemoveDuplicate will remove duplicate and copy item to new array
func (i Int64Array) RemoveDuplicate() Int64Array {
	var arr Int64Array
	m := map[int64]bool{}
	for _, v := range i {
		if m[v] {
			continue
		}
		m[v] = true
		arr = append(arr, v)
	}
	return arr
}

//AsPtrArray will convet normla to ptr
func (i Int64Array) AsPtrArray() (vals Int64PtrArray) {
	for _, v := range i {
		vals = append(vals, converter.Int64Ptr(v))
	}
	return
}

//Int64PtrArray is database value to parse data to []int64 value
type Int64PtrArray []*int64

//Scan is sql.Sanner
func (i *Int64PtrArray) Scan(src interface{}) (err error) {
	err = sqlScan(src, i, func(str string) (xerr error) {
		vals, xerr := converter.ArrayInt64Val(str)
		if xerr == nil {
			*i = Int64Array(vals).AsPtrArray()
		}
		return
	})
	return
}

//Value is driver.Valuer
func (i Int64PtrArray) Value() (driver.Value, error) {
	if i == nil {
		return "[]", nil
	}
	bys, err := json.Marshal(i)
	return string(bys), err
}

func (i Int64PtrArray) Len() int {
	return len(i)
}
func (i Int64PtrArray) Less(a, b int) bool {
	return i[a] == nil || (i[b] != nil && *i[a] < *i[b])
}
func (i Int64PtrArray) Swap(a, b int) {
	i[a], i[b] = i[b], i[a]
}

//HavingOne will check if array having one value in vals
func (i Int64PtrArray) HavingOne(vals ...int64) bool {
	for _, v0 := range i {
		for _, v1 := range vals {
			if v0 != nil && *v0 == v1 {
				return true
			}
		}
	}
	return false
}

//Join will parset to database array
func (i Int64PtrArray) Join(sep string) (res string) {
	for i, v := range i {
		if i > 0 {
			res += sep
		}
		if v == nil {
			res += "nil"
		} else {
			res += fmt.Sprintf("%v", *v)
		}
	}
	return
}

//DbArray will join value to database array
func (i Int64PtrArray) DbArray() (res string) {
	res = "{" + i.Join(",") + "}"
	return
}

//StrArray will join value to string array by comma
func (i Int64PtrArray) StrArray() (res string) {
	res = "," + i.Join(",") + ","
	return
}

//RemoveDuplicate will remove duplicate and copy item to new array
func (i Int64PtrArray) RemoveDuplicate() Int64PtrArray {
	var arr Int64PtrArray
	m := map[int64]bool{}
	for _, v := range i {
		if v == nil || m[*v] {
			continue
		}
		m[*v] = true
		arr = append(arr, v)
	}
	return arr
}

//AsArray will convet ptr to normal, skip nil
func (i Int64PtrArray) AsArray() (vals Int64Array) {
	for _, v := range i {
		if v == nil {
			continue
		}
		vals = append(vals, *v)
	}
	return
}

//Float64Array is database value to parse data to []int64 value
type Float64Array []float64

//Scan is sql.Sanner
func (f *Float64Array) Scan(src interface{}) (err error) {
	err = sqlScan(src, f, func(str string) (xerr error) {
		*f, xerr = converter.ArrayFloat64Val(str)
		return
	})
	return
}

//Value is driver.Valuer
func (f Float64Array) Value() (driver.Value, error) {
	if f == nil {
		return "[]", nil
	}
	bys, err := json.Marshal(f)
	return string(bys), err
}

func (f Float64Array) Len() int {
	return len(f)
}
func (f Float64Array) Less(a, b int) bool {
	return f[a] < f[b]
}
func (f Float64Array) Swap(a, b int) {
	f[a], f[b] = f[b], f[a]
}

//HavingOne will check if array having one value in vals
func (f Float64Array) HavingOne(vals ...float64) bool {
	for _, v0 := range f {
		for _, v1 := range vals {
			if v0 == v1 {
				return true
			}
		}
	}
	return false
}

//Join will parset to database array
func (f Float64Array) Join(sep string) (res string) {
	for i, v := range f {
		if i > 0 {
			res += sep
		}
		res += fmt.Sprintf("%v", v)
	}
	return
}

//DbArray will join value to database array
func (f Float64Array) DbArray() (res string) {
	res = "{" + f.Join(",") + "}"
	return
}

//StrArray will join value to string array by comma
func (i Float64Array) StrArray() (res string) {
	res = "," + i.Join(",") + ","
	return
}

//RemoveDuplicate will remove duplicate and copy item to new array
func (f Float64Array) RemoveDuplicate() Float64Array {
	var arr Float64Array
	m := map[float64]bool{}
	for _, v := range f {
		if m[v] {
			continue
		}
		m[v] = true
		arr = append(arr, v)
	}
	return arr
}

//AsPtrArray will convet normla to ptr
func (f Float64Array) AsPtrArray() (vals Float64PtrArray) {
	for _, v := range f {
		vals = append(vals, converter.Float64Ptr(v))
	}
	return
}

//Float64PtrArray is database value to parse data to []int64 value
type Float64PtrArray []*float64

//Scan is sql.Sanner
func (f *Float64PtrArray) Scan(src interface{}) (err error) {
	err = sqlScan(src, f, func(str string) (xerr error) {
		vals, xerr := converter.ArrayFloat64Val(str)
		if xerr == nil {
			*f = Float64Array(vals).AsPtrArray()
		}
		return
	})
	return
}

//Value is driver.Valuer
func (f Float64PtrArray) Value() (driver.Value, error) {
	if f == nil {
		return "[]", nil
	}
	bys, err := json.Marshal(f)
	return string(bys), err
}

func (f Float64PtrArray) Len() int {
	return len(f)
}
func (f Float64PtrArray) Less(a, b int) bool {
	return f[a] == nil || (f[b] != nil && *f[a] < *f[b])
}
func (f Float64PtrArray) Swap(a, b int) {
	f[a], f[b] = f[b], f[a]
}

//HavingOne will check if array having one value in vals
func (f Float64PtrArray) HavingOne(vals ...float64) bool {
	for _, v0 := range f {
		for _, v1 := range vals {
			if v0 != nil && *v0 == v1 {
				return true
			}
		}
	}
	return false
}

//Join will parset to database array
func (f Float64PtrArray) Join(sep string) (res string) {
	for i, v := range f {
		if i > 0 {
			res += sep
		}
		if v == nil {
			res += "nil"
		} else {
			res += fmt.Sprintf("%v", *v)
		}
	}
	return
}

//DbArray will join value to database array
func (f Float64PtrArray) DbArray() (res string) {
	res = "{" + f.Join(",") + "}"
	return
}

//StrArray will join value to string array by comma
func (i Float64PtrArray) StrArray() (res string) {
	res = "," + i.Join(",") + ","
	return
}

//RemoveDuplicate will remove duplicate and copy item to new array
func (f Float64PtrArray) RemoveDuplicate() Float64PtrArray {
	var arr Float64PtrArray
	m := map[float64]bool{}
	for _, v := range f {
		if v == nil || m[*v] {
			continue
		}
		m[*v] = true
		arr = append(arr, v)
	}
	return arr
}

//AsArray will convet ptr to normal, skip nil
func (f Float64PtrArray) AsArray() (vals Float64Array) {
	for _, v := range f {
		if v == nil {
			continue
		}
		vals = append(vals, *v)
	}
	return
}

//StringArray is database value to parse data to []string value
type StringArray []string

//Scan is sql.Sanner
func (s *StringArray) Scan(src interface{}) (err error) {
	err = sqlScan(src, s, func(str string) (xerr error) {
		*s, xerr = converter.ArrayStringVal(str)
		return
	})
	return
}

//Value will parse to json value
func (s StringArray) Value() (driver.Value, error) {
	if s == nil {
		return "[]", nil
	}
	bys, err := json.Marshal(s)
	return string(bys), err
}

func (s StringArray) Len() int {
	return len(s)
}
func (s StringArray) Less(a, b int) bool {
	return s[a] < s[b]
}
func (s StringArray) Swap(a, b int) {
	s[a], s[b] = s[b], s[a]
}

//HavingOne will check if array having one value in vals
func (s StringArray) HavingOne(vals ...string) bool {
	for _, v0 := range s {
		for _, v1 := range vals {
			if v0 == v1 {
				return true
			}
		}
	}
	return false
}

//Join will parset to database array
func (s StringArray) Join(sep string) (res string) {
	res = strings.Join(s, sep)
	return
}

//DbArray will join value to database array
func (s StringArray) DbArray() (res string) {
	res = "{" + s.Join(",") + "}"
	return
}

//StrArray will join value to string array by comma
func (i StringArray) StrArray() (res string) {
	res = "," + i.Join(",") + ","
	return
}

//RemoveDuplicate will remove duplicate and copy item to new array
func (s StringArray) RemoveDuplicate(trim, empty bool) StringArray {
	var arr StringArray
	m := map[string]bool{}
	for _, v := range s {
		nv := v
		if trim {
			nv = strings.TrimSpace(v)
		}
		if empty && len(nv) < 1 {
			continue
		}
		if m[nv] {
			continue
		}
		m[nv] = true
		arr = append(arr, nv)
	}
	return arr
}

//RemoveEmpty will remove empty and copy item to new array
func (s StringArray) RemoveEmpty(trim bool) StringArray {
	var arr StringArray
	for _, v := range s {
		nv := v
		if trim {
			nv = strings.TrimSpace(v)
		}
		if len(nv) < 1 {
			continue
		}
		arr = append(arr, nv)
	}
	return arr
}

//AsPtrArray will convet normla to ptr
func (s StringArray) AsPtrArray() (vals StringPtrArray) {
	for _, v := range s {
		vals = append(vals, converter.StringPtr(v))
	}
	return
}

//StringPtrArray is database value to parse data to []string value
type StringPtrArray []*string

//Scan is sql.Sanner
func (s *StringPtrArray) Scan(src interface{}) (err error) {
	err = sqlScan(src, s, func(str string) (xerr error) {
		vals, xerr := converter.ArrayStringVal(str)
		if xerr == nil {
			*s = StringArray(vals).AsPtrArray()
		}
		return
	})
	return
}

//Value will parse to json value
func (s StringPtrArray) Value() (driver.Value, error) {
	if s == nil {
		return "[]", nil
	}
	bys, err := json.Marshal(s)
	return string(bys), err
}

func (s StringPtrArray) Len() int {
	return len(s)
}
func (s StringPtrArray) Less(a, b int) bool {
	return s[a] == nil || (s[b] != nil && *s[a] < *s[b])
}
func (s StringPtrArray) Swap(a, b int) {
	s[a], s[b] = s[b], s[a]
}

//HavingOne will check if array having one value in vals
func (s StringPtrArray) HavingOne(vals ...string) bool {
	for _, v0 := range s {
		for _, v1 := range vals {
			if v0 != nil && *v0 == v1 {
				return true
			}
		}
	}
	return false
}

//Join will parset to database array
func (s StringPtrArray) Join(sep string) (res string) {
	for i, v := range s {
		if i > 0 {
			res += sep
		}
		if v == nil {
			res += "nil"
		} else {
			res += *v
		}
	}
	return
}

//DbArray will join value to database array
func (s StringPtrArray) DbArray() (res string) {
	res = "{" + s.Join(",") + "}"
	return
}

//StrArray will join value to string array by comma
func (i StringPtrArray) StrArray() (res string) {
	res = "," + i.Join(",") + ","
	return
}

//RemoveDuplicate will remove duplicate and copy item to new array
func (s StringPtrArray) RemoveDuplicate(trim, empty bool) StringPtrArray {
	var arr StringPtrArray
	m := map[string]bool{}
	for _, v := range s {
		if v == nil {
			continue
		}
		nv := v
		if trim {
			n := strings.TrimSpace(*v)
			nv = &n
		}
		if empty && len(*nv) < 1 {
			continue
		}
		if m[*nv] {
			continue
		}
		m[*nv] = true
		arr = append(arr, nv)
	}
	return arr
}

//RemoveEmpty will remove empty and copy item to new array
func (s StringPtrArray) RemoveEmpty(trim bool) StringPtrArray {
	var arr StringPtrArray
	for _, v := range s {
		if v == nil {
			continue
		}
		nv := v
		if trim {
			n := strings.TrimSpace(*v)
			nv = &n
		}
		if len(*nv) < 1 {
			continue
		}
		arr = append(arr, nv)
	}
	return arr
}

//AsArray will convet ptr to normal, skip nil
func (s StringPtrArray) AsArray() (vals StringArray) {
	for _, v := range s {
		if v == nil {
			continue
		}
		vals = append(vals, *v)
	}
	return
}
