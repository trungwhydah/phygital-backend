package helper

import (
	"fmt"
	"reflect"
	"strings"
)

// StructKeyToString convert keys of a struct to string. Ignore struct and slice type
func StructKeyToString(tableName string, bs interface{}) string {
	var str string
	val := reflect.Indirect(reflect.ValueOf(bs))

	for i := 0; i < val.NumField(); i++ {
		if val.Field(i).Kind() != reflect.Struct && val.Field(i).Kind() != reflect.Slice && val.Type().Field(i).Tag.Get("json") != "" {
			str += fmt.Sprintf("%s.%s, ", tableName, val.Type().Field(i).Tag.Get("json"))
		}
	}

	return strings.TrimSuffix(str, ", ")
}

func CheckFieldExists(data interface{}, fieldName string) bool {
	value := reflect.ValueOf(data)
	field := value.FieldByName(fieldName)
	return field.IsValid()
}

// ParseFieldName
// tagField: json or bson
func ParseFieldName(f reflect.StructField, tagField string) (name string, ignore bool) {
	tag := f.Tag.Get(tagField)
	// If we did not set it, return empty string
	if tag == "" {
		return "", false
	}

	// In case we set json/bson:"-"
	if tag == "-" {
		return "", true
	}

	if i := strings.Index(tag, ","); i != -1 {
		if i == 0 {
			return f.Name, false
		} else {
			return tag[:i], false
		}
	}

	// Not happy case, temporarily return "json"/"bson" and false
	return tag, false
}
