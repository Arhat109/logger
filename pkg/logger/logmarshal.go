package logger

import (
	"encoding/json"
	"fmt"
	"runtime"
	"time"
)

// JsonMessage -- структура сообщения для маршалирования лога
type JsonMessage struct {
	Level    string `json:"level"`
	DateTime string `json:"dateTime"`
	FuncName string `json:"funcName"`
	FileName string `json:"fileName"`
	LineNum  int    `json:"lineNum"`
	Message  string `json:"message"`
}

// StdlogMarshal -- местный маршаллер в JSON лога с типовыми параметрами
// !!! Поскольку функции логирования с контекстом совмещают trace_id с уровнем ошибки, вывод будет таким:
// { "level":"trace_id_unique: ERROR", ... } м.б. даже удобно
func StdlogMarshal(prefix, msg string, args ...interface{}) ([]byte, error) {
	dTime := time.Now()
	frame, _ := GetCaller(6)

	mess := JsonMessage{
		Level:    prefix,
		DateTime: dTime.Format("2006-01-02 03:04:05.000"),
		FuncName: frame.Function,
		FileName: frame.File,
		LineNum:  frame.Line,
		Message:  fmt.Sprintf(msg, args...),
	}
	return json.Marshal(mess)
}

// GetCallers -- отдает стек вызвавших логирование контекстов
func GetCallers(skip int) *runtime.Frames {
	rpc := make([]uintptr, 1)
	n := runtime.Callers(skip, rpc[:]) // ?!? не понял, но типа так..
	if n < 1 {
		return nil
	}

	return runtime.CallersFrames(rpc)
}

// GetCaller -- отдает собственно контекст, вызвавший логирование
func GetCaller(skip int) (runtime.Frame, bool) {
	frames := GetCallers(skip)

	return frames.Next()
}
