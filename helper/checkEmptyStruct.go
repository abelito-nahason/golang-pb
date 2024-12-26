package helper

import "reflect"

// checks if the given value is the zero value for its type
func isZeroValue(v reflect.Value) bool {
	return v.Interface() == reflect.Zero(v.Type()).Interface()
}

func CheckEmptyStruct(a interface{}) bool {
	v := reflect.ValueOf(a)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return true
	}

	for i := 0; i < v.NumField(); i++ {
		if isZeroValue(v.Field(i)) {
			return true
		}
	}
	return false
}
