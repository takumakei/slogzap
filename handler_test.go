package slogzap_test

import (
	"bytes"
	"context"
	"log/slog"
	"math"
	"testing"

	"github.com/takumakei/slogzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Example() {
	log := slog.New(slogzap.New(zap.NewExample()))
	log.Info("info")
	log.Log(context.TODO(), slog.LevelError, "error")
	log.Log(context.TODO(), slog.LevelError+4, "dpanic")
	// Output:
	// {"level":"info","msg":"info"}
	// {"level":"error","msg":"error"}
	// {"level":"dpanic","msg":"dpanic"}
}

func ExampleWithLimit() {
	log := slog.New(slogzap.New(zap.NewExample(), slogzap.WithLimit(zapcore.ErrorLevel)))
	log.Log(context.TODO(), slog.LevelError, "error")
	log.Log(context.TODO(), slog.LevelError+4, "dpanic") // slog.LevelError+4  => zapcore.DPanicLevel =(limit)=> zapcore.ErrorLevel
	log.Log(context.TODO(), slog.LevelError+8, "panic")  // slog.LevelError+8  => zapcore.PanicLevel  =(limit)=> zapcore.ErrorLevel
	log.Log(context.TODO(), slog.LevelError+12, "fatal") // slog.LevelError+12 => zapcore.FatalLevel  =(limit)=> zapcore.ErrorLevel
	log.Log(context.TODO(), slog.LevelInfo, "info")
	// Output:
	// {"level":"error","msg":"error"}
	// {"level":"error","msg":"dpanic"}
	// {"level":"error","msg":"panic"}
	// {"level":"error","msg":"fatal"}
	// {"level":"info","msg":"info"}
}

func TestHandler(t *testing.T) {
	example := func(log *slog.Logger) {
		log = log.With("a", 1)
		log = log.With("b", 2)
		log = log.WithGroup("A")
		log = log.With("c", 3)
		log = log.With("d", 4)
		log = log.WithGroup("B")
		log = log.With("e", 5)
		log.Debug("DEBUG", "f", 6)
		log.Info("INFO", "f", 6)
		log.Warn("WARN", "f", 6)
	}

	a := &bytes.Buffer{}
	example(
		slog.New(slog.NewJSONHandler(a, &slog.HandlerOptions{
			Level: slog.LevelWarn,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if len(groups) == 0 && a.Key == slog.TimeKey {
					return slog.Attr{}
				}
				return a
			},
		})),
	)

	b := &bytes.Buffer{}
	example(
		slog.New(slogzap.New(func() *zap.Logger {
			zc := zapcore.EncoderConfig{
				MessageKey:  "msg",
				LevelKey:    "level",
				EncodeLevel: zapcore.CapitalLevelEncoder,
			}
			core := zapcore.NewCore(zapcore.NewJSONEncoder(zc), zapcore.AddSync(b), zap.WarnLevel)
			return zap.New(core)
		}())),
	)

	want := a.String()
	got := b.String()
	t.Log(got)
	if want != got {
		t.Errorf("\nwant: %s\n got: %s", want, got)
	}
}

func TestGetLogger(t *testing.T) {
	want := zap.NewExample()
	got := slogzap.GetLogger(slog.New(slogzap.New(want)))
	if got != want {
		t.Error("should be equal")
	}
}

func TestHandlerLevel(t *testing.T) {
	h := slogzap.New(zap.NewExample()).(*slogzap.Handler)
	if h.Level() != slog.LevelDebug {
		t.Error("should be DEBUG, got ", h.Level())
	}
}

func TestHandlerLimit(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		h := slogzap.New(zap.NewExample()).(*slogzap.Handler)
		if h.Limit() != slog.Level(math.MaxInt) {
			t.Error("should be maximum, got ", h.Limit())
		}
	})
	t.Run("error", func(t *testing.T) {
		h := slogzap.New(zap.NewExample(), slogzap.WithLimit(zapcore.ErrorLevel)).(*slogzap.Handler)
		if h.Limit() != slog.LevelError {
			t.Error("should be ERROR, got ", h.Limit())
		}
	})
}
