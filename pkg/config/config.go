package config

import (
	"github.com/Arhat109/logger/pkg/dto"
	"github.com/Arhat109/logger/tools"
)

const (
	// EnvLoggerJson -- отдавать строкой или в JSON?
	EnvLoggerJson = "LOG_JSON"
	DefLoggerJson = "json"
	DefLoggerText = "txt"
	// EnvLoggerFlags -- формат вывода сообщений "работать молча"
	EnvLoggerFlags = "LOG_FLAGS"
	DefLoggerFlags = -1
	// EnvLoggerLevel -- уровень вывода сообщений логгером:
	EnvLoggerLevel = "LOG_LEVEL"
	DefLoggerLevel = "info"
	// EnvTraceId -- идентификатор сквозной трассировки, если не типовой
	EnvTraceId = "LOG_TRACE_ID"
	DefTraceId = dto.LogTracingId
)

type LogConfig struct {
	// ToJson внутренний лог в JSON или строками log.go?
	ToJson bool
	// Flags флаги отображения. @see log
	Flags int
	// Level std: 6 - debug, 5 - info, 4 - warn, 3 - error, 2 - fatal, 1 - panic
	Level int
	// TraceId идент сквозной трассировки. Может приходить в контексте.
	TraceId string
}

// New returns application config instance
func New() *LogConfig {
	cfg := LogConfig{
		ToJson:  tools.ToString(tools.LookupEnv(EnvLoggerJson, DefLoggerText)) == DefLoggerJson, // если не задано, то текстом
		Flags:   tools.ToInt(tools.LookupEnv(EnvLoggerFlags, DefLoggerFlags)),
		TraceId: tools.ToString(tools.LookupEnv(EnvTraceId, DefTraceId)),
	}
	strLevel := tools.ToString(tools.LookupEnv(EnvLoggerLevel, DefLoggerLevel))
	var level int
	switch strLevel {
	case "panic", dto.LogPanicPrefix:
		level = dto.LogPanicLevel
	case "fatal", dto.LogFatalPrefix:
		level = dto.LogFatalLevel
	case "error", dto.LogErrorPrefix:
		level = dto.LogErrorLevel
	case "warn", dto.LogWarnPrefix:
		level = dto.LogWarnLevel
	case "info", dto.LogInfoPrefix:
		level = dto.LogInfoLevel
	case "debug", dto.LogDebugPrefix:
		level = dto.LogDebugLevel
	}
	cfg.Level = level
	return &cfg
}
