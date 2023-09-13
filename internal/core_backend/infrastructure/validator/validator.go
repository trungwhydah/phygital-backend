package validation

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/go-playground/validator/v10"
)

type customValidator struct {
	Validator *validator.Validate
}

type CustomValidator interface {
	Validate(i interface{}) (err error)
}

// NewCustomValidator create custom validator struct
func NewCustomValidator(v *validator.Validate) CustomValidator {
	_ = v.RegisterValidation(`rfe`, requireWhen)
	_ = v.RegisterValidation(`max_if_parent`, maxIfParent)

	return &customValidator{Validator: v}
}

// Validate validate a struct
func (cv *customValidator) Validate(i interface{}) (err error) {
	return cv.Validator.Struct(i)
}

// callback function when using rfe tag (required filed equal)
func requireWhen(fl validator.FieldLevel) bool {
	params := strings.Split(fl.Param(), "&")
	for _, v := range params {
		reverseFlag := false
		param := strings.Split(v, ":")
		paramField := param[0]
		paramValue := param[1]

		if strings.HasPrefix(paramValue, "!") {
			reverseFlag = true
			paramValue = paramValue[1:]
		}

		if paramField == `` {
			return true
		}

		// param field reflect.Value.
		var paramFieldValue reflect.Value

		if fl.Parent().Kind() == reflect.Ptr {
			paramFieldValue = fl.Parent().Elem().FieldByName(paramField)
		} else {
			paramFieldValue = fl.Parent().FieldByName(paramField)
		}

		if isEq(paramFieldValue, paramValue) == reverseFlag {
			return true
		}

	}
	return hasValue(fl)
}

func hasValue(fl validator.FieldLevel) bool {
	return requireCheckFieldKind(fl, "")
}

func requireCheckFieldKind(fl validator.FieldLevel, param string) bool {
	field := fl.Field()
	if len(param) > 0 {
		if fl.Parent().Kind() == reflect.Ptr {
			field = fl.Parent().Elem().FieldByName(param)
		} else {
			field = fl.Parent().FieldByName(param)
		}
	}
	switch field.Kind() {
	case reflect.Slice, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Chan, reflect.Func:
		return !field.IsNil()
	default:
		_, _, nullable := fl.ExtractType(field)
		if nullable && field.Interface() != nil {
			return true
		}
		return field.IsValid() && field.Interface() != reflect.Zero(field.Type()).Interface()
	}
}

// maxIfParent Check validate max length with condition when using max_if_parent tag
func maxIfParent(fl validator.FieldLevel) bool {
	param := strings.Split(fl.Param(), " if ")
	paramValue := strings.TrimSpace(param[0])
	paramField := strings.TrimSpace(param[1])

	paramCondition := strings.Split(paramField, ":")
	conditionField := strings.TrimSpace(paramCondition[0])
	conditionValues := strings.TrimSpace(paramCondition[1])

	paramFieldValue := fl.Top().FieldByName(conditionField)

	for _, cv := range strings.Split(conditionValues, " ") {
		if cv != "" && isEq(paramFieldValue, cv) {
			return checkMax(fl, paramValue)
		}
	}

	return true
}

// checkMax check max length of field
func checkMax(fl validator.FieldLevel, max string) bool {
	field := fl.Field().Interface().(string)
	maxInt, _ := strconv.Atoi(max)

	return utf8.RuneCountInString(field) <= maxInt
}

func isEq(field reflect.Value, value string) bool {
	switch field.Kind() {
	case reflect.String:
		return field.String() == value
	case reflect.Slice, reflect.Map, reflect.Array:
		p := asInt(value)
		return int64(field.Len()) == p
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p := asInt(value)
		return field.Int() == p
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p := asUint(value)
		return field.Uint() == p
	case reflect.Float32, reflect.Float64:
		p := asFloat(value)
		return field.Float() == p
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}

func asInt(param string) int64 {
	i, err := strconv.ParseInt(param, 0, 64)
	panicIf(err)

	return i
}

func asUint(param string) uint64 {
	i, err := strconv.ParseUint(param, 0, 64)
	panicIf(err)

	return i
}

func asFloat(param string) float64 {
	i, err := strconv.ParseFloat(param, 64)
	panicIf(err)

	return i
}

func panicIf(err error) {
	if err != nil {
		panic(err.Error())
	}
}
