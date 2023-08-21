package logger

import (
	"time"
)

const (
	LogDefault   = 0
	LogShortFile = 1 << iota
	LogLongFile
	LogFuncName
	LogUTC
	LogDate
	LogTime
	LogMicroSeconds
	LogLevelColored
	LogWithTrace
	LogPrefix
)
const (
	LogNoneLevel  = 0
	LogPanicLevel = iota * 10
	LogFatalLevel
	LogErrorLevel
	LogWarnLevel
	LogInfoLevel
	LogDebugLevel

	LogPanicPrefix = "PANIC"
	LogFatalPrefix = "FATAL"
	LogErrorPrefix = "ERROR"
	LogWarnPrefix  = "WARN "
	LogInfoPrefix  = "INFO "
	LogDebugPrefix = "DEBUG"

	// LogFatalExitCode -- число, отдаваемое ОС при завершении программы из ...Fatal() вызовов.
	LogFatalExitCode = 500

	// CtxTraceId -- упрощенная версия идентификатора сквозного лога
	CtxTraceId = "trace_id"

	// BaseDefFlags местный стандарт формата вывода:
	BaseDefFlags = LogUTC | LogMicroSeconds | LogShortFile

	// @see New():
	// Нумерация дополнительных параметров конструктору логера в требуемом ему порядке:

	// BaselogJsonHandlerId порядковый номер доп. параметра конструкторам логгеров - внешний преобразователь в json
	// можно задать конструктору свой, если не задан применит BaseLogMarshal()
	BaselogJsonHandlerId = 0
	// CtxlogDefersId Номер в параметрах конструктору для предопределенного набора []func() -- завершателей контекстов
	CtxlogDefersId = 1

	// можно дополнять при необходимости..
)

// LogJsonHandler -- обработчик преобразования сообщения в json м.б. внешним
type LogJsonHandler func(buf *[]byte, depth int, now time.Time, level, message string) error

// Loggable -- Тот, кто умеет выводить сообщения разного уровня в логгер:
type Loggable interface {
	// GetLevel -- получить наименьший разрешенный уровень сообщений
	GetLevel() int
	// Outlog -- вывести сообщение с заданной глубины вызовов, от времени, с таким уровнем, текстом и аргументами
	Outlog(depth int, now time.Time, level, message string)
}

// Levelable -- тот, кто умеет выводить разные сообщения с уровнем не ниже указанного
type Levelable interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Fatal(msg string, args ...any)
	Panic(msg string, args ...any)
}
