package dto

import (
	"context"
)

const (
	LogDebugPrefix = "DEBUG"
	LogDebugLevel  = 6
	LogInfoPrefix  = "INFO"
	LogInfoLevel   = 5
	LogWarnPrefix  = "WARN"
	LogWarnLevel   = 4
	LogErrorPrefix = "ERROR"
	LogErrorLevel  = 3
	LogFatalPrefix = "FATAL"
	LogFatalLevel  = 2
	LogPanicPrefix = "PANIC"
	LogPanicLevel  = 1

	LogNoneLevel = 0

	// LogTracingId -- типовое название сквозного ключа трассировки запроса, ходящего по нескольким сервисам
	LogTracingId = "trace_id"

	// StdlogDefFlags местный стандарт формата вывода: log.Ldate log.Lmicroseconds log.Lshortfile log.Lmsgprefix: message\n
	StdlogDefFlags = -1

	// @see NewInternalLogger():
	// Дополнительные параметры генератору логера в требуемом ему порядке:

	// StdlogJsonHandlerId порядковый номер доп. параметра генератору - внешний преобразователь в json
	StdlogJsonHandlerId = 0
	// StdlogFatalDeferId -- порядковый номер в аргументах генератору логгера для функции-завершателя контекста при FATAL
	StdlogFatalDeferId = 1

	// можно дополнять при необходимости..
)

// LogJsonHandler -- обработчик преобразования сообщения в json м.б. внешним
type LogJsonHandler func(prefix, message string, args ...interface{}) (jsonMsg []byte, err error)

// Options -- настройки логирования, выделено в отд. стр-ру
type Options struct {
	EnvLevel int
	// Flags -- копия флагов форматирования, засунутая в стд. логгер
	Flags   int
	ToJson  LogJsonHandler
	TraceId string
	// FatalDefers -- слайс завершателей контекстов сервиса при ошибках круче ERROR
	FatalDefers []func()
}

// Loggable -- Возможность вывода сообщений разного уровня в логгер:
type Loggable interface {
	// GetLogger -- получить сам логгер по указателю в структуре или откуда ишо
	GetLogger() any
	// GetTraceId -- получить имя для сквозного trace_id префикса
	GetTraceId() string

	Sync() error // defer() если требуется логгеру

	// Простые методы логирования без контекста и trace_id

	Debug(msg string, args ...any) // stdout
	Info(msg string, args ...any)  // stdout
	Warn(msg string, args ...any)  // stdout
	Error(msg string, args ...any) // вывод в stderr без остановки
	Fatal(msg string, args ...any) // вывод в stderr с завершением

	// с контекстом, содержащим сквозной ключ trace_id:

	DebugCtx(ctx context.Context, msg string, args ...any)
	InfoCtx(ctx context.Context, msg string, args ...any)
	WarnCtx(ctx context.Context, msg string, args ...any)
	ErrorCtx(ctx context.Context, msg string, args ...any)
	FatalCtx(ctx context.Context, msg string, args ...any)
	PanicCtx(ctx context.Context, msg string, args ...any)
}
