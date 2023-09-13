package structtraversal

import (
	"reflect"
)

// TraverseObject for traverse through an Object tree and execute callback function on each node
// Only pointer nodes will be call by the callback.
func TraverseObject(obj any, callback func(args ...any)) {
	if reflect.TypeOf(obj).Kind() != reflect.Pointer {
		return
	}

	callback(obj)

	val := reflect.ValueOf(obj).Elem()
	typ := reflect.TypeOf(val.Interface())

	for i := 0; i < typ.NumField(); i++ {
		field := val.Field(i)
		if !field.CanInterface() {
			continue
		}

		if field.Kind() == reflect.Pointer {
			realType := field.Elem().Kind()

			if realType == reflect.Struct || realType == reflect.Map {
				TraverseObject(field.Elem().Interface(), callback)
			}
		} else if field.Kind() == reflect.Slice {
			TraverseSlice(field.Interface(), callback)
		}
	}
}
