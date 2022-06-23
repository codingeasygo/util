package attrscan

import (
	"reflect"
	"strings"
)

type NilChecker interface {
	IsNil() bool
}

type ZeroChecker interface {
	IsZero() bool
}

type NameConv func(on, name string, field reflect.StructField) string

type Scanner struct {
	Tag      string
	NameConv NameConv
}

var Default = &Scanner{
	Tag: "json",
	NameConv: func(on, name string, field reflect.StructField) string {
		return name
	},
}

func CheckValue(val reflect.Value, incNil, incZero bool) bool {
	return Default.CheckValue(val, incNil, incZero)
}

func (s *Scanner) CheckValue(val reflect.Value, incNil, incZero bool) bool {
	if !val.IsValid() {
		return incNil
	}
	v := val.Interface()
	if c, ok := v.(NilChecker); ok {
		return !c.IsNil() || incNil
	}
	if c, ok := v.(ZeroChecker); ok {
		return !c.IsZero() || incZero
	}
	kind := val.Kind()
	if kind == reflect.Ptr && val.IsNil() && !incNil {
		return false
	}
	if kind == reflect.Ptr && !val.IsNil() {
		val = reflect.Indirect(val)
	}
	if (!val.IsValid() || val.IsZero()) && !incZero {
		return false
	}
	return true
}

func FilterFieldCall(on string, v interface{}, filter string, call func(fieldName, fieldFunc string, field reflect.StructField, value interface{})) {
	Default.FilterFieldCall(on, v, filter, call)
}

func (s *Scanner) FilterFieldCall(on string, v interface{}, filter string, call func(fieldName, fieldFunc string, field reflect.StructField, value interface{})) {
	reflectValue := reflect.Indirect(reflect.ValueOf(v))
	reflectType := reflectValue.Type()
	var fieldAll = map[string]string{}
	var isExc = false
	var incNil, incZero bool
	if len(filter) > 0 {
		filter = strings.TrimSpace(filter)
		parts := strings.SplitN(filter, "#", 2)
		isExc = strings.HasPrefix(parts[0], "^")
		if len(parts[0]) > 0 {
			for _, fieldItem := range strings.Split(strings.TrimPrefix(parts[0], "^"), ",") {
				fieldParts := strings.SplitN(strings.Trim(strings.TrimSpace(fieldItem), ")"), "(", 2)
				if len(fieldParts) > 1 {
					fieldAll[fieldParts[1]] = fieldParts[0]
				} else {
					fieldAll[fieldParts[0]] = ""
				}
			}
		}
		if len(parts) > 1 && len(parts[1]) > 0 {
			incNil = strings.Contains(","+parts[1]+",", ",nil,") || strings.Contains(","+parts[1]+",", ",all,")
			incZero = strings.Contains(","+parts[1]+",", ",zero,") || strings.Contains(","+parts[1]+",", ",all,")
		}
	}
	numField := reflectType.NumField()
	for i := 0; i < numField; i++ {
		fieldValue := reflectValue.Field(i)
		fieldType := reflectType.Field(i)
		fieldName := strings.SplitN(fieldType.Tag.Get(s.Tag), ",", 2)[0]
		if len(fieldName) < 1 || fieldName == "-" {
			continue
		}
		if _, ok := fieldAll[fieldName]; (isExc && ok) || (!isExc && len(fieldAll) > 0 && !ok) {
			continue
		}
		filter := strings.TrimPrefix(fieldType.Tag.Get("filter"), "#")
		fieldIncNil, fieldIncZero := incNil, incZero
		if len(filter) > 0 {
			fieldIncNil = strings.Contains(","+filter+",", ",nil,") || strings.Contains(","+filter+",", ",all,")
			fieldIncZero = strings.Contains(","+filter+",", ",zero,") || strings.Contains(","+filter+",", ",all,")
		}
		if !s.checkValue(fieldValue, fieldIncNil, fieldIncZero) {
			continue
		}
		fieldName = s.NameConv(on, fieldName, fieldType)
		call(fieldName, fieldAll[fieldName], fieldType, fieldValue.Addr().Interface())
	}
}
