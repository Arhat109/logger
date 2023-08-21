package logger

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// BaseLogger -- местная реализация базового интерфейса.
// По мотивам пакета log с упрощениями и добавлением типовых потребностей для полноценного логирования:
// 1. Управление уровнем вывода и возможность вывода "в никуда".
// 2. Внутренний вызов fmt.Sprintf() для вставки в лог параметров на лету
// 3. Возможность вывода в формате json
// 4. Реентерабельность и возможность применения одного логера в нескольких горутинах. Внутри структуры ничего лишнего нет.
// 5. Закрытие файла вывода при логировании Fatal() и Panic()
// блокировка мьютексом только непосредственно вывода сообщения в поток(файл)
// Все поля структуры публичны, для полноценного внедрения по мере потребности программиста в развитии пакета.
type BaseLogger struct {
	Mu sync.Mutex
	// куда выводить сообщения (nil - в никуда)
	Out io.Writer
	// Level -- наибольший разрешенный уровень вывода сообщений
	Level int
	// Flags -- Что выводить в лог
	Flags int
	// ToJson -- маршаллер сообщений в json, если задан. Иначе - строка
	ToJson LogJsonHandler
}

const bufSize = 1024

var bufPool = sync.Pool{
	New: func() any { return new([bufSize]byte) },
}

func (baselog *BaseLogger) GetLevel() int { return baselog.Level }

// Init -- настройка логгера из структуры настроек @see ./config, возвращает себя (this)
// param Args -- доп. параметры конфигуратора (если надо!): тут можно задать маршаллер в json
// Ошибка открытия файла лога возвращается в параметре, исключительно для улучшения работы escape алгоритма.
func (baselog *BaseLogger) Init(cfg *LogConfig, retErr *error, args ...any) *BaseLogger {
	baselog.Level = cfg.Level
	baselog.Flags = cfg.Flags
	baselog.Mu = sync.Mutex{}

	switch cfg.Out {
	case "devnul":
		baselog.Out = nil
	case "stdout":
		baselog.Out = os.Stdout
	case "stderr":
		baselog.Out = os.Stderr
	default:
		if cfg.Out != "" {
			file, err := os.OpenFile(cfg.Out, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0x740)
			if err != nil {
				*retErr = fmt.Errorf("BaseLogger.Init() OpenFile has: %s", err.Error())
			}
			baselog.Out = file
		}
	}

	if cfg.IsJson {
		baselog.ToJson = BaseLogMarshal
		if len(args) > BaselogJsonHandlerId {
			toJson, ok := args[BaselogJsonHandlerId].(LogJsonHandler)
			if ok && toJson != nil {
				baselog.ToJson = toJson
			}
		}
	}
	if baselog.Level >= LogWarnLevel {
		baselog.Info("Logger installed for with %d Level", baselog.Level)
	}
	return baselog
}

// FormatString -- строчный режим вывода. Форматирует строку в заданном буфере
// порядок элементов в строке фиксирован:
// Level:{ yyyy-mm-dd hh:mm:ss.msec}{ file_name#line) | func_name#line }message\n
// {} -- опционально, если они есть. Склеиваются перед сообщением "как есть", разделять самостоятельно!
func (baselog *BaseLogger) FormatString(buf *[]byte, depth int, now time.Time, level, message string) {
	if color, ok := GlDefColors[level]; baselog.Flags&LogLevelColored != 0 && ok {
		FormatColored(buf, color, level)
		*buf = append(*buf, ':')
	} else {
		*buf = append(*buf, level...)
		*buf = append(*buf, ':')
	}

	if baselog.Flags&(LogDate|LogTime|LogMicroSeconds) != 0 {
		FormatTime(buf, now, baselog.Flags)
	}

	if baselog.Flags&(LogShortFile|LogLongFile) != 0 {
		FormatFileLine(buf, depth+1, baselog.Flags&LogShortFile != 0)
	} else if baselog.Flags&LogFuncName != 0 {
		FormatFuncLine(buf, depth+1)
	}

	*buf = append(*buf, ' ')
	*buf = append(*buf, message...)
	if len(message) == 0 || message[len(message)-1] != '\n' {
		*buf = append(*buf, '\n')
	}
}

// OutMessage -- вывод сообщения в поток логирования(файл) или в никуда
func (baselog *BaseLogger) OutMessage(content *[]byte) error {
	if baselog.Out != nil {
		baselog.Mu.Lock()
		_, err := baselog.Out.Write(*content)
		baselog.Mu.Unlock()
		return err
	}
	return nil
}

// Outlog -- собственно форматилка лога и его вывод куда сказано.
// Ориентировочно: level=6 символов, date=11, time=9, micro=4, long/short file=32/16, message <120> итого ~182символа
// Аллоцируем тут, для обеспечения реентерабельности в горутинах.
func (baselog *BaseLogger) Outlog(depth int, now time.Time, level, message string) {
	abuf := bufPool.Get().(*[bufSize]byte)
	buf := (*abuf)[:0]
	depth++
	if baselog.ToJson != nil { // JSON! Все формируем тут по частям:
		if err := baselog.ToJson(&buf, depth, now, level, message); err != nil {
			// преобразование в json не получилось, игнор ошибки т.к. далее не JSON:
			baselog.FormatString(&buf, depth, now, level, message)
		}
	} else {
		baselog.FormatString(&buf, depth, now, level, message)
	}

	if err := baselog.OutMessage(&buf); err != nil && baselog.Out != os.Stderr {
		if _, err := os.Stderr.Write(buf); err != nil {
			panic(err.Error())
		}
	}
	bufPool.Put(abuf)
}

// Простое логирование по уровням с добавлением доп. полей по настройкам

func (baselog *BaseLogger) Debug(msg string, args ...any) {
	if baselog.Level >= LogDebugLevel {
		baselog.Outlog(1, time.Now(), LogDebugPrefix, fmt.Sprintf(msg, args...))
	}
}
func (baselog *BaseLogger) Info(msg string, args ...any) {
	if baselog.Level >= LogInfoLevel {
		baselog.Outlog(1, time.Now(), LogInfoPrefix, fmt.Sprintf(msg, args...))
	}
}
func (baselog *BaseLogger) Warn(msg string, args ...any) {
	if baselog.Level >= LogWarnLevel {
		baselog.Outlog(1, time.Now(), LogWarnPrefix, fmt.Sprintf(msg, args...))
	}
}
func (baselog *BaseLogger) Error(msg string, args ...any) {
	if baselog.Level >= LogErrorLevel {
		baselog.Outlog(1, time.Now(), LogErrorPrefix, fmt.Sprintf(msg, args...))
	}
}
func (baselog *BaseLogger) Fatal(msg string, args ...any) {
	if baselog.Level >= LogFatalLevel {
		baselog.Outlog(1, time.Now(), LogFatalPrefix, fmt.Sprintf(msg, args...))
		os.Exit(1)
	}
}
func (baselog *BaseLogger) Panic(msg string, args ...any) {
	if baselog.Level >= LogPanicLevel {
		message := fmt.Sprintf(msg, args...)
		baselog.Outlog(1, time.Now(), LogPanicPrefix, message)
		panic(message)
	}
}

// NewBaseLogger -- генератор (в куче!) нового логгера из настроек @see ./config.go
// param Args -- доп. параметры конфигуратора (если надо!): тут можно задать маршаллер в json
func NewBaseLogger(cfg *LogConfig, args ...any) (Loggable, error) {
	var err error
	logger := &BaseLogger{}

	_ = logger.Init(cfg, &err, args...)

	return logger, err
}
