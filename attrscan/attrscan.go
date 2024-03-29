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
	if c, ok := v.(NilChecker); ok && c.IsNil() {
		return incNil
	}
	if c, ok := v.(ZeroChecker); ok && c.IsZero() {
		return incZero
	}
	if val.CanAddr() {
		v := val.Addr().Interface()
		// if c, ok := v.(NilChecker); ok && c.IsNil() {
		// 	return incNil
		// }
		if c, ok := v.(ZeroChecker); ok && c.IsZero() {
			return incZero
		}
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
	for _, f := range strings.Split(filter, "|") {
		s.filterFieldCall(on, v, f, call)
	}
}

func (s *Scanner) filterFieldCall(on string, v interface{}, filter string, call func(fieldName, fieldFunc string, field reflect.StructField, value interface{})) {
	reflectValue := reflect.Indirect(reflect.ValueOf(v))
	reflectType := reflectValue.Type()
	var fieldAll = map[string]string{}
	var isExc = false
	var incNil, incZero bool
	var alias string
	if len(filter) > 0 {
		filter = strings.TrimSpace(filter)
		parts := strings.SplitN(filter, ".", 2)
		if len(parts) > 1 {
			alias = parts[0] + "."
			filter = parts[1]
		}
		parts = strings.SplitN(filter, "#", 2)
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
		fieldFilter := strings.TrimSpace(strings.TrimPrefix(fieldType.Tag.Get("filter"), "#"))
		fieldIncNil, fieldIncZero, fieldInline := incNil, incZero, false
		if len(fieldFilter) > 0 {
			fieldIncNil = strings.Contains(","+fieldFilter+",", ",nil,") || strings.Contains(","+fieldFilter+",", ",all,")
			fieldIncZero = strings.Contains(","+fieldFilter+",", ",zero,") || strings.Contains(","+fieldFilter+",", ",all,")
			fieldInline = strings.Contains(","+fieldFilter+",", ",inline,")
		}
		if fieldInline {
			s.FilterFieldCall(on, fieldValue.Addr().Interface(), filter, call)
			continue
		}
		if len(fieldName) < 1 || fieldName == "-" {
			continue
		}
		if _, ok := fieldAll[fieldName]; (isExc && ok) || (!isExc && len(fieldAll) > 0 && !ok) {
			continue
		}

		if !s.CheckValue(fieldValue, fieldIncNil, fieldIncZero) {
			continue
		}
		fieldName = s.NameConv(on, fieldName, fieldType)
		call(alias+fieldName, fieldAll[fieldName], fieldType, fieldValue.Addr().Interface())
	}
}
