package xmap

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/codingeasygo/util/converter"
)

const ShouldIsNil ShouldAction = "IsNil"
const ShouldIsNoNil ShouldAction = "IsNoNil"
const ShouldIsZero ShouldAction = "IsZero"
const ShouldIsNoZero ShouldAction = "IsNoZero"
const ShouldIsEmpty ShouldAction = "IsEmpty"
const ShouldIsNoEmpty ShouldAction = "IsNoEmpty"
const ShouldIsInt ShouldAction = "IsInt"
const ShouldIsUint ShouldAction = "IsUint"
const ShouldIsFloat ShouldAction = "IsFloat"
const ShouldEQ ShouldAction = "EQ"
const ShouldGT ShouldAction = "GT"
const ShouldGTE ShouldAction = "GTE"
const ShouldLT ShouldAction = "LT"
const ShouldLTE ShouldAction = "LTE"

type ShouldAction string

func (s ShouldAction) isEmpty(val reflect.Value) bool {
	return (val.Kind() == reflect.Map || val.Kind() == reflect.Slice || val.Kind() == reflect.Array || val.Kind() == reflect.String) && val.Len() == 0
}

func (s ShouldAction) Compare(x, y interface{}) bool {
	a, b := reflect.ValueOf(x), reflect.ValueOf(y)
	if a.CanInt() && b.CanInt() {
		av := a.Convert(reflect.TypeOf(int64(0))).Interface().(int64)
		bv := b.Convert(reflect.TypeOf(int64(0))).Interface().(int64)
		if s == ShouldEQ && av != bv {
			return false
		}
		if s == ShouldGT && av <= bv {
			return false
		}
		if s == ShouldGTE && av < bv {
			return false
		}
		if s == ShouldLT && av >= bv {
			return false
		}
		if s == ShouldLTE && av > bv {
			return false
		}
		return true
	}
	if a.CanUint() && b.CanUint() {
		av := a.Convert(reflect.TypeOf(uint64(0))).Interface().(uint64)
		bv := b.Convert(reflect.TypeOf(uint64(0))).Interface().(uint64)
		if s == ShouldEQ && av != bv {
			return false
		}
		if s == ShouldGT && av <= bv {
			return false
		}
		if s == ShouldGTE && av < bv {
			return false
		}
		if s == ShouldLT && av >= bv {
			return false
		}
		if s == ShouldLTE && av > bv {
			return false
		}
		return true
	}
	if a.CanFloat() && b.CanFloat() {
		av := a.Convert(reflect.TypeOf(float64(0))).Interface().(float64)
		bv := b.Convert(reflect.TypeOf(float64(0))).Interface().(float64)
		if s == ShouldEQ && av != bv {
			return false
		}
		if s == ShouldGT && av <= bv {
			return false
		}
		if s == ShouldGTE && av < bv {
			return false
		}
		if s == ShouldLT && av >= bv {
			return false
		}
		if s == ShouldLTE && av > bv {
			return false
		}
		return true
	}
	return ValueEqual(x, y)
}

func (s ShouldAction) Check(v interface{}) bool {
	val := reflect.ValueOf(v)
	if !val.IsValid() {
		return s == ShouldIsNil || s == ShouldIsZero
	}
	if s == ShouldIsNil && (val.Kind() != reflect.Ptr || (val.Kind() == reflect.Ptr && !val.IsNil())) {
		return false
	}
	if s == ShouldIsNoNil && val.Kind() == reflect.Ptr && val.IsNil() {
		return false
	}
	if s == ShouldIsZero && !val.IsZero() {
		return false
	}
	if s == ShouldIsNoZero && val.IsZero() {
		return false
	}
	if s == ShouldIsEmpty && !s.isEmpty(val) {
		return false
	}
	if s == ShouldIsNoEmpty && s.isEmpty(val) {
		return false
	}
	if s == ShouldIsInt && !val.CanInt() {
		return false
	}
	if s == ShouldIsUint && !val.CanUint() {
		return false
	}
	if s == ShouldIsFloat && !val.CanFloat() {
		return false
	}
	return true
}

func (m M) Should(args ...interface{}) (err error) {
	n := len(args)
	for i := 0; i < n; {
		if i+1 >= n {
			err = fmt.Errorf("args[%v] action is not setted", i)
			break
		}
		key, ok := args[i].(string)
		if !ok {
			err = fmt.Errorf("args[%v] key is not string", i)
			break
		}
		val := m.Value(key)
		action, ok := args[i+1].(ShouldAction)
		if !ok {
			if !ValueEqual(val, args[i+1]) {
				err = fmt.Errorf("m.%v(%v,%v)!=args[%v](%v,%v)", key, reflect.TypeOf(val), val, i+1, reflect.TypeOf(args[i+1]), args[i+1])
				break
			}
			i += 2
			continue
		}
		if strings.HasPrefix(string(action), "Is") {
			if !action.Check(val) {
				err = fmt.Errorf("m.%v(%v,%v)!=%v", key, reflect.TypeOf(val), val, action)
				break
			}
			i += 2
			continue
		}
		if i+2 >= n {
			err = fmt.Errorf("args[%v] compare value is not setted", i)
			break
		}
		if !action.Compare(val, args[i+2]) {
			err = fmt.Errorf("m.%v(%v,%v) %v args[%v](%v,%v)", key, reflect.TypeOf(val), val, action, i+2, reflect.TypeOf(args[i+2]), args[i+2])
			break
		}
		i += 3
	}
	return
}

type Shoulder struct {
	Log        *log.Logger
	testerFail func()
	testerSkip func()
	shouldErr  bool
	shouldArgs []interface{}
	onlyLog    bool
}

func (s *Shoulder) Should(t *testing.T, args ...interface{}) *Shoulder {
	if t != nil {
		s.testerFail, s.testerSkip, s.shouldArgs = t.Fail, t.SkipNow, args
	}
	s.shouldArgs = append(s.shouldArgs, args...)
	return s
}

func (s *Shoulder) ShouldError(t *testing.T) *Shoulder {
	if t != nil {
		s.testerFail, s.testerSkip = t.Fail, t.SkipNow
	}
	s.shouldErr = true
	return s
}

func (s *Shoulder) OnlyLog(only bool) *Shoulder {
	s.onlyLog = only
	return s
}

func (s *Shoulder) callError(depth int, err error) {
	if s.testerFail == nil {
		panic(err)
	}
	if s.Log == nil {
		s.Log = log.New(os.Stderr, "    ", log.Llongfile)
	}
	s.Log.Output(depth, err.Error())
	if !s.onlyLog {
		s.testerFail()
		s.testerSkip()
	}
}

func (s *Shoulder) validError(depth int, res M, err error) bool {
	if err != nil {
		s.callError(depth+1, fmt.Errorf("%v, res is %v", err, converter.JSON(res)))
		return false
	}
	return true
}

func (s *Shoulder) validShould(depth int, res M, err error) bool {
	if len(s.shouldArgs) < 1 {
		return true
	}
	xerr := res.Should(s.shouldArgs...)
	if xerr != nil {
		s.callError(depth+1, fmt.Errorf("%v, res is %v", xerr, converter.JSON(res)))
		return false
	}
	return true
}

func (s *Shoulder) Valid(depth int, res M, err error) bool {
	if s.shouldErr {
		if err == nil {
			s.callError(depth, fmt.Errorf("err is nil, res is %v", converter.JSON(res)))
			return false
		}
	} else {
		if !s.validError(depth+1, res, err) {
			return false
		}
		if !s.validShould(depth+1, res, err) {
			return false
		}
	}
	return true
}
