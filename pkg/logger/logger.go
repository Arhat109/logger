package logger

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/Arhat109/logger/pkg/config"
	"github.com/Arhat109/logger/pkg/dto"
)

// StdLogger -- местная реализация интерфейса через стандартный логгер
type StdLogger struct {
	Log *log.Logger
	dto.Options
}

// GetLogger -- получить логгер
func (t *StdLogger) GetLogger() any     { return t.Log }
func (t *StdLogger) GetTraceId() string { return t.TraceId }

// CloseDefers -- завершение всех defer контекста сервиса перед вызовом os.Exit()
func (t *StdLogger) CloseDefers() {
	if t.FatalDefers != nil {
		// вызываем в обратном порядке!
		for i := len(t.FatalDefers) - 1; i >= 0; i-- {
			t.FatalDefers[i]()
		}
	}
}

// Outlog -- собственно форматилка лога и его вывод куда сказано.
func (t *StdLogger) Outlog(logTo io.Writer, depth int, prefix, msg string, args ...any) {
	var err error

	t.Log.SetOutput(logTo)
	if t.ToJson != nil { // JSON! Все формируем тут по частям:
		var mb []byte
		mb, err = t.ToJson(prefix, msg, args...)
		if err != nil {
			panic(err) // а некуда, иначе!
		}

		// фиговое решение, т.к. t.Log.mutex не доступен.. :(
		t.Log.SetFlags(0)
		err = t.Log.Output(depth, string(mb))
		t.Log.SetFlags(t.Flags)
	} else {
		var message string

		t.Log.SetPrefix(prefix + ": ")
		message = fmt.Sprintf(msg, args...)
		err = t.Log.Output(depth, message)
	}
	if err != nil {
		panic(err) // нечем логировать..
	}
}

// OutlogCtx -- добавляет t.TraceId из контекста в префикс форматилки
func (t *StdLogger) OutlogCtx(ctx context.Context, logTo io.Writer, prefix, msg string, args ...any) {
	traceTag := ctx.Value(t.TraceId)
	if traceTag != nil {
		tracingId, ok := traceTag.(string)
		if ok {
			prefix = t.TraceId + ": " + tracingId + ", " + prefix
		}
	}
	t.Outlog(logTo, 4, prefix, msg, args...)
}

func (t *StdLogger) Sync() error { return nil } // тут не надо. Заглушка.

// Простое логирование по уровням с добавлением доп. полей по настройкам

func (t *StdLogger) Debug(msg string, args ...any) {
	if t.EnvLevel >= dto.LogDebugLevel {
		t.Outlog(os.Stdout, 3, dto.LogDebugPrefix, msg, args...)
	}
}
func (t *StdLogger) Info(msg string, args ...any) {
	if t.EnvLevel >= dto.LogInfoLevel {
		t.Outlog(os.Stdout, 3, dto.LogInfoPrefix, msg, args...)
	}
}
func (t *StdLogger) Warn(msg string, args ...any) {
	if t.EnvLevel >= dto.LogWarnLevel {
		t.Outlog(os.Stdout, 3, dto.LogWarnPrefix, msg, args...)
	}
}
func (t *StdLogger) Error(msg string, args ...any) {
	if t.EnvLevel >= dto.LogErrorLevel {
		t.Outlog(os.Stderr, 3, dto.LogErrorPrefix, msg, args...)
	}
}
func (t *StdLogger) Fatal(msg string, args ...any) {
	if t.EnvLevel >= dto.LogFatalLevel {
		t.Outlog(os.Stderr, 3, dto.LogFatalPrefix, msg, args...)
		t.CloseDefers()
		os.Exit(500)
	}
}
func (t *StdLogger) Panic(msg string, args ...any) {
	if t.EnvLevel >= dto.LogPanicLevel {
		t.Outlog(os.Stderr, 3, dto.LogPanicPrefix, msg, args...)
		t.CloseDefers()
		panic(msg)
	}
}

// Дополнительно вытаскивают из контекста сквозной идент запроса для логирования

func (t *StdLogger) DebugCtx(ctx context.Context, msg string, args ...any) {
	if t.EnvLevel >= dto.LogDebugLevel {
		t.OutlogCtx(ctx, os.Stdout, dto.LogDebugPrefix, msg, args...)
	}
}
func (t *StdLogger) InfoCtx(ctx context.Context, msg string, args ...any) {
	if t.EnvLevel >= dto.LogInfoLevel {
		t.OutlogCtx(ctx, os.Stdout, dto.LogInfoPrefix, msg, args...)
	}
}
func (t *StdLogger) WarnCtx(ctx context.Context, msg string, args ...any) {
	if t.EnvLevel >= dto.LogWarnLevel {
		t.OutlogCtx(ctx, os.Stdout, dto.LogWarnPrefix, msg, args...)
	}
}
func (t *StdLogger) ErrorCtx(ctx context.Context, msg string, args ...any) {
	if t.EnvLevel >= dto.LogErrorLevel {
		t.OutlogCtx(ctx, os.Stdout, dto.LogErrorPrefix, msg, args...)
	}
}
func (t *StdLogger) FatalCtx(ctx context.Context, msg string, args ...any) {
	if t.EnvLevel >= dto.LogFatalLevel {
		t.OutlogCtx(ctx, os.Stdout, dto.LogFatalPrefix, msg, args...)
		t.CloseDefers()
		ctx.Done()
		os.Exit(500)
	}
}
func (t *StdLogger) PanicCtx(ctx context.Context, msg string, args ...any) {
	if t.EnvLevel >= dto.LogPanicLevel {
		t.OutlogCtx(ctx, os.Stdout, dto.LogPanicPrefix, msg, args...)
		t.CloseDefers()
		ctx.Done()
		panic(msg)
	}
}

// New -- генератор нового логгера
// @param args -- доп. параметры конфигуратора (если надо!):
//
//	.. [0] | [ StdlogJsonHandlerId ] dto.LogJsonHandler -- преобразователь в json
//	.. [1] | [ StdlogFatalDeferId ]  []*func()          -- список defer для закрытия контескта сервиса при FATAL
func New(cfg *config.LogConfig, args ...any) (dto.Loggable, error) {

	stdLevel := cfg.Level
	flags := cfg.Flags
	if cfg.Flags == dto.StdlogDefFlags {
		flags = log.Ldate | log.Lmicroseconds | log.Lshortfile | log.Lmsgprefix
	}

	var toJson dto.LogJsonHandler
	var ok = true
	if cfg.ToJson {
		if len(args) > dto.StdlogJsonHandlerId {
			toJson, ok = args[dto.StdlogJsonHandlerId].(dto.LogJsonHandler)
		}
		if toJson == nil || !ok {
			toJson = StdlogMarshal
		}
	}

	var fatalDefers []func()
	if len(args) > dto.StdlogFatalDeferId {
		fatalDefers, ok = args[dto.StdlogFatalDeferId].([]func())
	}

	lgr := &StdLogger{
		Log: log.New(os.Stdout, "", flags),
		Options: dto.Options{
			EnvLevel:    stdLevel,
			ToJson:      toJson,
			Flags:       flags,
			FatalDefers: fatalDefers,
			TraceId:     cfg.TraceId,
		},
	}
	lgr.Info("Logger installed for with %d level", stdLevel)
	return lgr, nil
}
