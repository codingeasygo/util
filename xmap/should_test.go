package xmap

import (
	"testing"
)

func TestShould(t *testing.T) {
	var ptr string
	var ptr2 *string
	vals := M{
		"int":     123,
		"int2":    0,
		"uint":    uint(123),
		"uint2":   uint(0),
		"float":   float64(123),
		"float2":  float64(0),
		"string":  "123",
		"string2": "",
		"ptr":     &ptr,
		"ptr2":    ptr2,
		"invalid": nil,
		"slice":   []int{123},
		"slice2":  []int{},
		"map": map[string]int{
			"a": 123,
		},
		"map2": map[string]int{},
	}
	{ //check
		if err := vals.Should("int", 123); err != nil {
			t.Error(err)
			return
		}
		if err := vals.Should("uint", uint(123)); err != nil {
			t.Error(err)
			return
		}
		if err := vals.Should("float", float64(123)); err != nil {
			t.Error(err)
			return
		}
		if err := vals.Should("string", "123"); err != nil {
			t.Error(err)
			return
		}
		if err := vals.Should("int", ShouldIsNoZero); err != nil {
			t.Error(err)
			return
		}
		if err := vals.Should("ptr", ShouldIsNoNil); err != nil {
			t.Error(err)
			return
		}
		if err := vals.Should("slice", ShouldIsNoEmpty); err != nil {
			t.Error(err)
			return
		}
		if err := vals.Should("map", ShouldIsNoEmpty); err != nil {
			t.Error(err)
			return
		}
		if err := vals.Should("ptr2", ShouldIsNil); err != nil {
			t.Error(err)
			return
		}
		if err := vals.Should("invalid", ShouldIsNil); err != nil {
			t.Error(err)
			return
		}
		if err := vals.Should("invalid", ShouldIsZero); err != nil {
			t.Error(err)
			return
		}
		if err := vals.Should("slice2", ShouldIsEmpty); err != nil {
			t.Error(err)
			return
		}
		if err := vals.Should("map2", ShouldIsEmpty); err != nil {
			t.Error(err)
			return
		}
		//not
		if err := vals.Should("int2", 123); err == nil {
			t.Error(err)
			return
		}
		if err := vals.Should("uint2", uint(123)); err == nil {
			t.Error(err)
			return
		}
		if err := vals.Should("float2", float64(123)); err == nil {
			t.Error(err)
			return
		}
		if err := vals.Should("string2", "123"); err == nil {
			t.Error(err)
			return
		}
		if err := vals.Should("int", ShouldIsZero); err == nil {
			t.Error(err)
			return
		}
		if err := vals.Should("ptr", ShouldIsNil); err == nil {
			t.Error(err)
			return
		}
		if err := vals.Should("slice", ShouldIsEmpty); err == nil {
			t.Error(err)
			return
		}
		if err := vals.Should("map", ShouldIsEmpty); err == nil {
			t.Error(err)
			return
		}
		if err := vals.Should("int2", ShouldIsNoZero); err == nil {
			t.Error(err)
			return
		}
		if err := vals.Should("ptr2", ShouldIsNoNil); err == nil {
			t.Error(err)
			return
		}
		if err := vals.Should("invalid", ShouldIsNoNil); err == nil {
			t.Error(err)
			return
		}
		if err := vals.Should("slice2", ShouldIsNoEmpty); err == nil {
			t.Error(err)
			return
		}
		if err := vals.Should("map2", ShouldIsNoEmpty); err == nil {
			t.Error(err)
			return
		}
		if err := vals.Should("ptr", ShouldIsInt); err == nil {
			t.Error(err)
			return
		}
		if err := vals.Should("ptr", ShouldIsUint); err == nil {
			t.Error(err)
			return
		}
		if err := vals.Should("ptr", ShouldIsFloat); err == nil {
			t.Error(err)
			return
		}
	}
	{ //compare
		if err := vals.Should("int", ShouldEQ, 123); err != nil {
			t.Error(err)
			return
		}
		if err := vals.Should("uint", ShouldEQ, uint(123)); err != nil {
			t.Error(err)
			return
		}
		if err := vals.Should("float", ShouldEQ, float64(123)); err != nil {
			t.Error(err)
			return
		}
		if err := vals.Should("int", ShouldEQ, uint(123)); err != nil {
			t.Error(err)
			return
		}
		if err := vals.Should("uint", ShouldEQ, 123); err != nil {
			t.Error(err)
			return
		}
		if err := vals.Should("float", ShouldEQ, 123); err != nil {
			t.Error(err)
			return
		}
		if err := vals.Should("int", ShouldGT, 122); err != nil {
			t.Error(err)
			return
		}
		if err := vals.Should("uint", ShouldGT, uint(122)); err != nil {
			t.Error(err)
			return
		}
		if err := vals.Should("float", ShouldGT, float64(122)); err != nil {
			t.Error(err)
			return
		}
		if err := vals.Should("int", ShouldGTE, 123); err != nil {
			t.Error(err)
			return
		}
		if err := vals.Should("uint", ShouldGTE, uint(123)); err != nil {
			t.Error(err)
			return
		}
		if err := vals.Should("float", ShouldGTE, float64(123)); err != nil {
			t.Error(err)
			return
		}
		if err := vals.Should("int", ShouldLT, 124); err != nil {
			t.Error(err)
			return
		}
		if err := vals.Should("uint", ShouldLT, uint(124)); err != nil {
			t.Error(err)
			return
		}
		if err := vals.Should("float", ShouldLT, float64(124)); err != nil {
			t.Error(err)
			return
		}
		if err := vals.Should("int", ShouldLTE, 123); err != nil {
			t.Error(err)
			return
		}
		if err := vals.Should("uint", ShouldLTE, uint(123)); err != nil {
			t.Error(err)
			return
		}
		if err := vals.Should("float", ShouldLTE, float64(123)); err != nil {
			t.Error(err)
			return
		}
		//not
		if err := vals.Should("int", ShouldEQ, 124); err == nil {
			t.Error(err)
			return
		}
		if err := vals.Should("uint", ShouldEQ, uint(124)); err == nil {
			t.Error(err)
			return
		}
		if err := vals.Should("float", ShouldEQ, float64(124)); err == nil {
			t.Error(err)
			return
		}
		if err := vals.Should("int", ShouldGT, 124); err == nil {
			t.Error(err)
			return
		}
		if err := vals.Should("uint", ShouldGT, uint(124)); err == nil {
			t.Error(err)
			return
		}
		if err := vals.Should("float", ShouldGT, float64(124)); err == nil {
			t.Error(err)
			return
		}
		if err := vals.Should("int", ShouldGTE, 124); err == nil {
			t.Error(err)
			return
		}
		if err := vals.Should("uint", ShouldGTE, uint(124)); err == nil {
			t.Error(err)
			return
		}
		if err := vals.Should("float", ShouldGTE, float64(124)); err == nil {
			t.Error(err)
			return
		}
		if err := vals.Should("int", ShouldLT, 122); err == nil {
			t.Error(err)
			return
		}
		if err := vals.Should("uint", ShouldLT, uint(122)); err == nil {
			t.Error(err)
			return
		}
		if err := vals.Should("float", ShouldLT, float64(122)); err == nil {
			t.Error(err)
			return
		}
		if err := vals.Should("int", ShouldLTE, 122); err == nil {
			t.Error(err)
			return
		}
		if err := vals.Should("uint", ShouldLTE, uint(122)); err == nil {
			t.Error(err)
			return
		}
		if err := vals.Should("float", ShouldLTE, float64(122)); err == nil {
			t.Error(err)
			return
		}
	}
	{ //error
		if err := vals.Should("xx"); err == nil {
			t.Error(err)
			return
		}
		if err := vals.Should(1, 2); err == nil {
			t.Error(err)
			return
		}
		if err := vals.Should("xxx", ShouldLTE); err == nil {
			t.Error(err)
			return
		}
	}
}
