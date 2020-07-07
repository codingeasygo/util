package xsql

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

//Time is database value to parse data from database and parset time.Time to timestamp on json mashal
type Time time.Time

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
	now := time.Now()
	return Time(time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()))
}

// TimeStartOfMonth return 00:00:00 of today
func TimeStartOfMonth() Time {
	now := time.Now()
	return Time(time.Date(now.Year(), now.Month(), 0, 0, 0, 0, 0, now.Location()))
}

//Timestamp return timestamp
func (t Time) Timestamp() int64 {
	return time.Time(t).Local().UnixNano() / 1e6
}

//MarshalJSON marshal time to string
func (t *Time) MarshalJSON() ([]byte, error) {
	raw := t.Timestamp()
	if raw < 0 {
		return []byte("null"), nil
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

//IntArray is database value to parse data to []int64 value
type IntArray []int

//Scan is sql.Sanner
func (i *IntArray) Scan(src interface{}) (err error) {
	if src != nil {
		if jsonSrc, ok := src.(string); ok {
			err = json.Unmarshal([]byte(jsonSrc), i)
			if err != nil {
				err = fmt.Errorf("Unmarshal fail with %v by :%v", err, jsonSrc)
			}
		} else {
			err = fmt.Errorf("the %v,%v is not string", reflect.TypeOf(src), src)
		}
	}
	return
}

//Value is driver.Valuer
func (i IntArray) Value() (driver.Value, error) {
	if i == nil {
		return nil, nil
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

//IntPtrArray is database value to parse data to []int64 value
type IntPtrArray []*int

//Scan is sql.Sanner
func (i *IntPtrArray) Scan(src interface{}) (err error) {
	if src != nil {
		if jsonSrc, ok := src.(string); ok {
			err = json.Unmarshal([]byte(jsonSrc), i)
			if err != nil {
				err = fmt.Errorf("Unmarshal fail with %v by :%v", err, jsonSrc)
			}
		} else {
			err = fmt.Errorf("the %v,%v is not string", reflect.TypeOf(src), src)
		}
	}
	return
}

//Value is driver.Valuer
func (i IntPtrArray) Value() (driver.Value, error) {
	if i == nil {
		return nil, nil
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

//Int64Array is database value to parse data to []int64 value
type Int64Array []int64

//Scan is sql.Sanner
func (i *Int64Array) Scan(src interface{}) (err error) {
	if src != nil {
		if jsonSrc, ok := src.(string); ok {
			err = json.Unmarshal([]byte(jsonSrc), i)
			if err != nil {
				err = fmt.Errorf("Unmarshal fail with %v by :%v", err, jsonSrc)
			}
		} else {
			err = fmt.Errorf("the %v,%v is not string", reflect.TypeOf(src), src)
		}
	}
	return
}

//Value is driver.Valuer
func (i Int64Array) Value() (driver.Value, error) {
	if i == nil {
		return nil, nil
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

//Int64PtrArray is database value to parse data to []int64 value
type Int64PtrArray []*int64

//Scan is sql.Sanner
func (i *Int64PtrArray) Scan(src interface{}) (err error) {
	if src != nil {
		if jsonSrc, ok := src.(string); ok {
			err = json.Unmarshal([]byte(jsonSrc), i)
			if err != nil {
				err = fmt.Errorf("Unmarshal fail with %v by :%v", err, jsonSrc)
			}
		} else {
			err = fmt.Errorf("the %v,%v is not string", reflect.TypeOf(src), src)
		}
	}
	return
}

//Value is driver.Valuer
func (i Int64PtrArray) Value() (driver.Value, error) {
	if i == nil {
		return nil, nil
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
			if *v0 == v1 {
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

//Float64Array is database value to parse data to []int64 value
type Float64Array []float64

//Scan is sql.Sanner
func (f *Float64Array) Scan(src interface{}) (err error) {
	if src != nil {
		if jsonSrc, ok := src.(string); ok {
			err = json.Unmarshal([]byte(jsonSrc), f)
			if err != nil {
				err = fmt.Errorf("Unmarshal fail with %v by :%v", err, jsonSrc)
			}
		} else {
			err = fmt.Errorf("the %v,%v is not string", reflect.TypeOf(src), src)
		}
	}
	return
}

//Value is driver.Valuer
func (f Float64Array) Value() (driver.Value, error) {
	if f == nil {
		return nil, nil
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

//Float64PtrArray is database value to parse data to []int64 value
type Float64PtrArray []*float64

//Scan is sql.Sanner
func (f *Float64PtrArray) Scan(src interface{}) (err error) {
	if src != nil {
		if jsonSrc, ok := src.(string); ok {
			err = json.Unmarshal([]byte(jsonSrc), f)
			if err != nil {
				err = fmt.Errorf("Unmarshal fail with %v by :%v", err, jsonSrc)
			}
		} else {
			err = fmt.Errorf("the %v,%v is not string", reflect.TypeOf(src), src)
		}
	}
	return
}

//Value is driver.Valuer
func (f Float64PtrArray) Value() (driver.Value, error) {
	if f == nil {
		return nil, nil
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
			if *v0 == v1 {
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

//M is database value to parse json data to map value
type M map[string]interface{}

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

//StringArray is database value to parse data to []string value
type StringArray []string

//Scan is sql.Sanner
func (s *StringArray) Scan(src interface{}) (err error) {
	if src != nil {
		if jsonSrc, ok := src.(string); ok {
			err = json.Unmarshal([]byte(jsonSrc), s)
		} else {
			err = fmt.Errorf("the %v,%v is not string", reflect.TypeOf(src), src)
		}
	}
	return
}

//Value will parse to json value
func (s StringArray) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
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
