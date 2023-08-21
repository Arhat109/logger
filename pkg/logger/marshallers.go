package logger

import (
	"encoding/json"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

// JsonMessage -- структура сообщения для маршалирования лога
type JsonMessage struct {
	Level    string `json:"Level"`
	DateTime string `json:"date_time"`
	FuncName string `json:"func_name"`
	FileName string `json:"file_name"`
	LineNum  int    `json:"line_num"`
	Message  string `json:"message"`
}

// BaseLogMarshal -- местный маршаллер в JSON лога с типовыми параметрами
func BaseLogMarshal(buf *[]byte, depth int, now time.Time, level, message string) error {
	frame, _ := GetCaller(depth + 1)

	mess := JsonMessage{
		Level:    level,
		DateTime: now.Format("2006-01-02 15:04:05.000"),
		FuncName: filepath.Base(frame.Function),
		FileName: frame.File,
		LineNum:  frame.Line,
		Message:  message,
	}
	var err error
	*buf, err = json.Marshal(mess)
	return err
}

// MapLogMarshal -- местный маршаллер в JSON лога с типовыми параметрами
func MapLogMarshal(buf *[]byte, depth int, now time.Time, level, message string, args ...any) ([]byte, error) {
	frame, _ := GetCaller(depth + 1)

	mess := map[string]string{
		"Level":    level,
		"DateTime": now.Format("2006-01-02 15:04:05.000"),
		"FuncName": filepath.Base(frame.Function),
		"FileName": frame.File,
		"LineNum":  strconv.Itoa(frame.Line),
		"Message":  message,
	}
	for i, arg := range args {
		mess[strconv.Itoa(i)] = ToString(arg)
	}
	return json.Marshal(mess)
}

// GetCallers -- отдает стек вызвавших логирование контекстов
func GetCallers(skip int) *runtime.Frames {
	rpc := make([]uintptr, 1)
	n := runtime.Callers(skip+1, rpc[:])
	if n < 1 {
		return nil
	}

	return runtime.CallersFrames(rpc)
}

// GetCaller -- отдает собственно контекст, вызвавший логирование
func GetCaller(skip int) (runtime.Frame, bool) {
	frames := GetCallers(skip + 1)

	return frames.Next()
}
