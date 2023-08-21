package tests

import (
	"bytes"
	"github.com/Arhat109/logger/pkg/logger"
	"log"
	"runtime"
	"testing"
)

var glbuf = [1024]byte{}
var bufWriter = bytes.NewBuffer(glbuf[:0])
var baseLgr = logger.BaseLogger{}

func Benchmark_Baselog(b *testing.B) {
	var err error
	baseLgr.Init(&logger.LogConfig{
		IsJson: false,
		Flags:  logger.LogDate,
		Level:  logger.LogDebugLevel,
	}, &err)
	baseLgr.Out = bufWriter
	runtime.GC()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bufWriter.Reset()
		baseLgr.Debug("this is a message")
	}
}

var stdLgr = log.New(bufWriter, "", log.Ldate)

func Benchmark_Stdlog(b *testing.B) {
	runtime.GC()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bufWriter.Reset()
		stdLgr.Printf("any info message")
	}
}

/*
type BufSyncer struct {
	buf *bytes.Buffer
}

func (bs *BufSyncer) Write(content []byte) (int, error) {
	bs.buf.Reset()
	bs.buf.Write(content)
	return len(content), nil
}
func (bs *BufSyncer) Sync() error { return nil }

var bufSyncer = BufSyncer{buf: bufWriter}

var zapLgr = zap.New(
	zapcore.NewCore(
		zapcore.NewConsoleEncoder(
			zapcore.EncoderConfig{
				TimeKey:  "T",
				LevelKey: "L",
				//				CallerKey:        "C",
				MessageKey:  "M",
				LineEnding:  zapcore.DefaultLineEnding,
				EncodeLevel: zapcore.LowercaseLevelEncoder,
				EncodeTime:  zapcore.RFC3339TimeEncoder,
				//				EncodeCaller:     zapcore.ShortCallerEncoder,
				ConsoleSeparator: ": ",
			},
		), &bufSyncer, zap.DebugLevel,
	))

func Benchmark_Zaplog(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		zapLgr.Info("any info message")
	}
}

var GlRes string

func Benchmark_Sptintf(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GlRes = fmt.Sprintf("any info message")
	}
}

func Benchmark_Caller(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, GlRes, _, _ = runtime.Caller(1)
	}
}
*/
