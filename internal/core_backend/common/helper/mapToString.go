package helper

import "fmt"

func MapToString(inputMap map[string]string) (result string) {
	for k, v := range inputMap {
		result += fmt.Sprintf("%s=%s&", k, v)
	}
	result = result[:len(result)-1]

	return
}
