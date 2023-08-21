package logger

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
	DefTraceId = CtxTraceId
)

// LogConfig -- настройки для инициализации логгеров. Одна структура на все типы.
// Если чего-то у логгера нет, то можно проигнорировать данный параметр тут.
type LogConfig struct {
	// "" по умолчанию в stderr, иначе полный путь к файлу лога или "stdout"|"devnul"
	Out string
	// IsJson формировать лог в JSON (true) или строками (false)?
	IsJson bool
	// Flags флаги отображения. ==0 вывод только сообщения через Sprintf
	Flags int
	// Level std: 60 - debug, 50 - info, 40 - warn, 30 - error, 20 - fatal, 10 - panic, <10 не выводим ничего.
	// позволяет в обертках применить расширение уровней на свое усмотрение..
	Level int
	// TraceId идент сквозной трассировки. Может приходить в контексте для CtxLogger
	TraceId string
}

// Init -- формирование настроек. Возвращает this
func (cfg *LogConfig) Init() *LogConfig {
	cfg.IsJson = ToString(LookupEnv(EnvLoggerJson, DefLoggerText)) == DefLoggerJson
	cfg.Flags = ToInt(LookupEnv(EnvLoggerFlags, DefLoggerFlags))
	cfg.TraceId = ToString(LookupEnv(EnvTraceId, DefTraceId))

	strLevel := ToString(LookupEnv(EnvLoggerLevel, DefLoggerLevel))
	var level int
	switch strLevel {
	case "panic", LogPanicPrefix:
		level = LogPanicLevel
	case "fatal", LogFatalPrefix:
		level = LogFatalLevel
	case "error", LogErrorPrefix:
		level = LogErrorLevel
	case "warn", LogWarnPrefix:
		level = LogWarnLevel
	case "info", LogInfoPrefix:
		level = LogInfoLevel
	case "debug", LogDebugPrefix:
		level = LogDebugLevel
	}
	cfg.Level = level

	return cfg
}

// NewLogConfig returns application config instance
func NewLogConfig() *LogConfig {
	cfg := LogConfig{}
	return (&cfg).Init()
}
