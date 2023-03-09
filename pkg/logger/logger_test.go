package logger

import (
	"bytes"
	"encoding/json"
	"github.com/Arhat109/logger/pkg/config"
	"github.com/Arhat109/logger/pkg/dto"
	"log"
	"path/filepath"
	"testing"
	"time"
)

func PrepareLogger(toJson bool) (*config.LogConfig, dto.Loggable) {
	cfg := &config.LogConfig{
		ToJson:  toJson,
		Flags:   log.Ldate | log.Lshortfile,
		Level:   dto.LogInfoLevel,
		TraceId: "test_tracing",
	}
	lgr, _ := New(cfg)

	return cfg, lgr
}

// Остальное тестировать нет особого смысла, т.к. всё работает через эту функцию
func TestStdLogger_Outlog(t *testing.T) {

	_, lgr := PrepareLogger(true)

	expect := JsonMessage{
		Level:    dto.LogInfoPrefix,
		DateTime: time.Now().Format("2006-01-02 03:04:05.000"),
		FuncName: "testing.tRunner",
		FileName: "testing.go",
		Message:  "test 1",
	}

	if l, ok := lgr.(*StdLogger); !ok {
		t.Errorf("expected *StdLogger, but created %T", l)
	} else {
		lWriter := bytes.NewBuffer(make([]byte, 0, 2048))

		// подменить буферизованным Writer и сравнить с образцом вывод функции
		l.Outlog(lWriter, 3, dto.LogInfoPrefix, "test %d", 1)

		var actual JsonMessage
		err := json.Unmarshal(lWriter.Bytes(), &actual)
		if err == nil {
			expect.LineNum = actual.LineNum // в разных режимах компиляции номер строки может не совпадать!
			actual.FileName = filepath.Base(actual.FileName)
			if actual != expect {
				t.Errorf("\nexpected \"%#v\"\n, has    \"%#v\"\n", expect, actual)
			}
		} else {
			t.Errorf("test error is:%#v", err)
		}
	}
}
