package cache

import (
	"reflect"
)

// Size returns how many bytes Write would generate to encode the value v, which
// must be a fixed-size value or a slice of fixed-size values, or a pointer to such data.
// If v is neither of these, Size returns -1.
func size(v interface{}) int {
	return dataSize(reflect.Indirect(reflect.ValueOf(v)))
}

// dataSize returns the number of bytes the actual data represented by v occupies in memory.
// For compound structures, it sums the sizes of the elements. Thus, for instance, for a slice
// it returns the length of the slice times the element size and does not count the memory
// occupied by the header. If the type of v is not acceptable, dataSize returns -1.
func dataSize(v reflect.Value) int {
	if v.Kind() == reflect.Slice {
		sum := 0
		for i := 0; i < v.Len(); i++ {
			s := sizeof(v.Index(i), v.Type().Elem())
			if s < 0 {
				return -1
			}
			sum += s
		}
		return sum
	}
	return sizeof(v, v.Type())
}

// sizeof returns the size >= 0 of variables for the given type or -1 if the type is not acceptable.
func sizeof(v reflect.Value, t reflect.Type) int {
	switch t.Kind() {
	case reflect.Array:
		if s := sizeof(v, t.Elem()); s >= 0 {
			return s * t.Len()
		}

	case reflect.Struct:
		sum := 0
		for i, n := 0, t.NumField(); i < n; i++ {
			s := dataSize(v.Field(i))
			if s < 0 {
				return -1
			}
			sum += s
		}
		return sum

	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128,
		reflect.Uintptr, reflect.Ptr:
		return int(t.Size())

	case reflect.String:
		return len(v.String())

	case reflect.Map:
		sum := 0
		for _, mapKey := range v.MapKeys() {
			keySize := dataSize(mapKey)
			valueSize := dataSize(v.MapIndex(mapKey))
			if keySize < 0 || valueSize < 0 {
				return -1
			}
			sum += keySize + valueSize
		}
		return sum
	}

	return -1
}
