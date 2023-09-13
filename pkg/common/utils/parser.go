package utils

import (
	"reflect"

	"backend-service/pkg/common/logger"
)

// ToSlice returns slices of specific of struct from list interface.
func ToSlice[E comparable](sources any, dest *[]E) bool {
	valueInterfaces, ok := sources.([]interface{})
	if !ok {
		logger.Errorw(
			"value must is array",
			"value type", reflect.TypeOf(sources),
			"value", sources,
		)

		return false
	}

	for _, item := range valueInterfaces {
		result, ok := item.(E)
		if !ok {
			logger.Errorw(
				"value must is float64",
				"value type", reflect.TypeOf(item),
				"value", item,
			)

			return false
		}

		*dest = append(*dest, result)
	}

	return true
}
