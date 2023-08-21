package logger

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

// ToInt -- Преобразование результата к целому @see config.go
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

// ToString -- Преобразование результата к строке @see config.go
func ToString(val interface{}) string {
	switch v := val.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case fmt.Stringer:
		return v.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}

// itoaBuf -- Cheap integer to fixed-width decimal ASCII. Give a negative width to avoid zero-padding.
// @author стырено log это лучше чем strconv.AppendInt()
func itoaBuf(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}
