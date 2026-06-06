package agent

import (
	"fmt"
	"reflect"
)

// deepEqualRobust handles comparison between potentially different numeric types
// and nested structures.
func deepEqualRobust(a, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	va := reflect.ValueOf(a)
	vb := reflect.ValueOf(b)

	// Handle numeric types by converting to float64 for comparison
	if isNumeric(va.Kind()) && isNumeric(vb.Kind()) {
		return va.Convert(reflect.TypeOf(float64(0))).Float() == vb.Convert(reflect.TypeOf(float64(0))).Float()
	}

	// Handle string vs numeric (common in XML transformations)
	if va.Kind() == reflect.String && isNumeric(vb.Kind()) {
		return va.String() == fmt.Sprintf("%v", b)
	}
	if isNumeric(va.Kind()) && vb.Kind() == reflect.String {
		return fmt.Sprintf("%v", a) == vb.String()
	}

	if va.Kind() != vb.Kind() {
		return false
	}

	switch va.Kind() {
	case reflect.Map:
		if va.Len() != vb.Len() {
			return false
		}
		for _, key := range va.MapKeys() {
			valA := va.MapIndex(key)
			valB := vb.MapIndex(key)
			if !valB.IsValid() || !deepEqualRobust(valA.Interface(), valB.Interface()) {
				return false
			}
		}
		return true
	case reflect.Slice, reflect.Array:
		if va.Len() != vb.Len() {
			return false
		}
		for i := 0; i < va.Len(); i++ {
			if !deepEqualRobust(va.Index(i).Interface(), vb.Index(i).Interface()) {
				return false
			}
		}
		return true
	default:
		return reflect.DeepEqual(a, b)
	}
}

func isNumeric(k reflect.Kind) bool {
	switch k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}
