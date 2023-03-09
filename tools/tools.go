package tools

import (
	"fmt"
	"os"
	"strconv"
)

type Errors []error

var glErrors Errors

// GetErrorList -- получение списка ошибок за пределами пакета
func GetErrorList() Errors { return glErrors }

// LookupEnv @return -- значение переменной окружения или дефолтное (пофиг: нет или с ошибкой)
func LookupEnv(name string, defVal interface{}) interface{} {
	if val, ok := os.LookupEnv(name); ok {
		return val
	}
	return defVal
}

// ToInt -- Преобразование результата к целому
func ToInt(val interface{}) int {
	switch v := val.(type) {
	case int:
		return v
	case string:
		if intVal, err := strconv.Atoi(v); err == nil {
			return intVal
		}
	}
	glErrors = append(glErrors, fmt.Errorf(
		"ToInt() ERROR! value %v is not INTEGER or not be converted to INT!", val,
	))

	return 0
}

// ToString -- Преобразование результата к строке
func ToString(val interface{}) string {
	if strVal, ok := val.(string); ok {
		return strVal
	}
	glErrors = append(glErrors, fmt.Errorf(
		"ToString() ERROR! value %v is not STRING!", val,
	))
	return ""
}
