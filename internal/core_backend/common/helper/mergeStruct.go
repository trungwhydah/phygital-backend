package helper

import "reflect"

// `origin` is the original struct that have fields will be updated by the `updated` struct
// `origin` MUST be passed in as the pointer to the original struct
// `updated` MUST be passed as the updated struct (not pointer)
// All field names in `updated` MUST be in `origin`
// All field data type in `updated` must be pointer to the orignal type
func MergeStructsField(origin interface{}, updated interface{}) error {
	valUpdated := reflect.ValueOf(updated)
	valOrigin := reflect.ValueOf(origin).Elem()
	typA := valUpdated.Type()

	for i := 0; i < valUpdated.NumField(); i++ {
		fieldUpdated := valUpdated.Field(i)
		fieldName := typA.Field(i).Name

		// This is a pointer address not the actual value
		fieldValue := fieldUpdated.Interface()

		if fieldValue != nil && !reflect.ValueOf(fieldValue).IsNil() {
			// Dereference the pointer to get the real value
			realValue := reflect.Indirect(fieldUpdated).Interface()

			// Assign corresponding field in struct B
			fieldOrigin := valOrigin.FieldByName(fieldName)
			if fieldOrigin.IsValid() && fieldOrigin.Type().AssignableTo(reflect.TypeOf(realValue)) {
				fieldOrigin.Set(reflect.ValueOf(realValue))
			}
		}
	}

	return nil
}
