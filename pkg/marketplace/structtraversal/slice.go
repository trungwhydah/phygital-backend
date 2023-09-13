package structtraversal

import (
	"reflect"
)

// TraverseSlice traverse through each slice element and try to Traverse through their structure.
func TraverseSlice(sl any, callback func(args ...any)) {
	val := reflect.ValueOf(sl)
	typ := reflect.TypeOf(sl)

	if typ.Kind() != reflect.Slice {
		return
	}

	for i := 0; i < val.Len(); i++ {
		elemReflection := val.Index(i)

		if !elemReflection.CanInterface() {
			continue
		}

		if elemReflection.Kind() == reflect.Pointer {
			realType := elemReflection.Elem().Kind()

			if realType == reflect.Struct || realType == reflect.Map {
				TraverseObject(elemReflection.Interface(), callback)
			}
		} else if elemReflection.Kind() == reflect.Slice {
			TraverseSlice(elemReflection.Interface(), callback)
		}
	}
}
