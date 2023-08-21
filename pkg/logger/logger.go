package logger

import (
	"context"
	"fmt"
	"time"
)

type LogItem struct {
	logger *BaseLogger
	// ctx -- контекст, может содержать канал завершения для FATAL и PANIC
	ctx context.Context
	// traceVal -- идентификатор трейсинга
	traceVal string
	// buf -- указатель на слайс буфера формирования вывода
	buf *[]byte
	// depth -- глубина вызова в стеке, на которой надо искать file,function,line
	depth int
	// logTime -- время логирования.
	logTime time.Time
	// level -- уровень текущего сообщения
	level int
	// message -- само сообщение
	message string
	// prefixes -- стек префиксов сообщения
	prefixes []string
}

// Сеттеры установки данных:

func (lgr *BaseLogger) With(prefix string) *LogItem {
	buf := make([]byte, 0, bufSize)
	item := LogItem{
		logger:   lgr,
		buf:      &buf,
		depth:    2,
		logTime:  time.Now(),
		prefixes: make([]string, 0, 1),
	}
	item.AddPrefix(prefix)
	return &item
}
func (li *LogItem) SetCtx(ctx context.Context, traceId any) *LogItem {
	var ok bool
	if li.traceVal, ok = ctx.Value(traceId).(string); !ok {
		li.traceVal = ""
	}
	li.ctx = ctx
	return li
}
func (li *LogItem) SetBuf(buf *[]byte) *LogItem { li.buf = buf; return li }
func (li *LogItem) SetDepth(depth int) *LogItem { li.depth = depth; return li }
func (li *LogItem) SetLevel(level int) *LogItem { li.level = level; return li }
func (li *LogItem) AddPrefix(prefix string) *LogItem {
	li.prefixes = append(li.prefixes, prefix)
	return li
}
func (li *LogItem) DelPrefix() *LogItem {
	lenP := len(li.prefixes)
	if lenP > 0 {
		li.prefixes = li.prefixes[:lenP-1]
	}
	return li
}

// Геттеры. Дополнять по необходимости

func (li *LogItem) GetLevel() int { return li.logger.Level }

func (li *LogItem) Outlog(depth int, now time.Time, level, message string) {
	if li.logger.Level < li.level {
		return
	}
	li.logger.Outlog(depth+1, now, level, message)
}

func (li *LogItem) Debug(msg string, args ...any) {
	li.SetLevel(LogDebugLevel)
	li.Outlog(1, time.Now(), LogDebugPrefix, fmt.Sprintf(msg, args...))
}
func (li *LogItem) Info(msg string, args ...any)  {}
func (li *LogItem) Warn(msg string, args ...any)  {}
func (li *LogItem) Error(msg string, args ...any) {}
func (li *LogItem) Fatal(msg string, args ...any) {}
func (li *LogItem) Panic(msg string, args ...any) {}
