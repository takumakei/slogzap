package levelconv_test

import (
	"log/slog"
	"testing"

	"github.com/takumakei/slogzap/levelconv"
	"go.uber.org/zap/zapcore"
)

func TestToSlog(t *testing.T) {
	tests := []struct {
		zlvl zapcore.Level
		want slog.Level
	}{
		{zlvl: zapcore.DebugLevel, want: slog.LevelDebug},
		{zlvl: zapcore.InfoLevel, want: slog.LevelInfo},
		{zlvl: zapcore.WarnLevel, want: slog.LevelWarn},
		{zlvl: zapcore.ErrorLevel, want: slog.LevelError},
		{zlvl: zapcore.DPanicLevel, want: slog.LevelError + 4},
		{zlvl: zapcore.PanicLevel, want: slog.LevelError + 8},
		{zlvl: zapcore.FatalLevel, want: slog.LevelError + 12},
	}
	for _, tt := range tests {
		got := levelconv.ToSlog(tt.zlvl)
		if got != tt.want {
			t.Errorf("ToSlog(%v) = %v, want %v", tt.zlvl, got, tt.want)
		}
	}
}

func TestToZap(t *testing.T) {
	tests := []struct {
		slvl slog.Level
		want zapcore.Level
	}{
		{slog.LevelDebug - 3, zapcore.DebugLevel},
		{slog.LevelDebug - 2, zapcore.DebugLevel},
		{slog.LevelDebug - 1, zapcore.DebugLevel},
		{slog.LevelDebug, zapcore.DebugLevel},
		{slog.LevelDebug + 1, zapcore.DebugLevel},
		{slog.LevelDebug + 2, zapcore.DebugLevel},
		{slog.LevelDebug + 3, zapcore.DebugLevel},

		{slog.LevelInfo, zapcore.InfoLevel},
		{slog.LevelInfo + 1, zapcore.InfoLevel},
		{slog.LevelInfo + 2, zapcore.InfoLevel},
		{slog.LevelInfo + 3, zapcore.InfoLevel},

		{slog.LevelWarn, zapcore.WarnLevel},
		{slog.LevelWarn + 1, zapcore.WarnLevel},
		{slog.LevelWarn + 2, zapcore.WarnLevel},
		{slog.LevelWarn + 3, zapcore.WarnLevel},

		{slog.LevelError, zapcore.ErrorLevel},
		{slog.LevelError + 1, zapcore.ErrorLevel},
		{slog.LevelError + 2, zapcore.ErrorLevel},
		{slog.LevelError + 3, zapcore.ErrorLevel},

		{slog.LevelError + 4, zapcore.DPanicLevel},
		{slog.LevelError + 4 + 1, zapcore.DPanicLevel},
		{slog.LevelError + 4 + 2, zapcore.DPanicLevel},
		{slog.LevelError + 4 + 3, zapcore.DPanicLevel},

		{slog.LevelError + 8, zapcore.PanicLevel},
		{slog.LevelError + 8 + 1, zapcore.PanicLevel},
		{slog.LevelError + 8 + 2, zapcore.PanicLevel},
		{slog.LevelError + 8 + 3, zapcore.PanicLevel},

		{slog.LevelError + 8 + 4, zapcore.FatalLevel},
		{slog.LevelError + 8 + 5, zapcore.FatalLevel},
		{slog.LevelError + 8 + 6, zapcore.FatalLevel},
		{slog.LevelError + 8 + 7, zapcore.FatalLevel},
		{slog.LevelError + 8 + 8, zapcore.FatalLevel},
	}
	for _, tt := range tests {
		got := levelconv.ToZap(tt.slvl)
		if got != tt.want {
			t.Errorf("ToZap(%v) = %v, want %v", tt.slvl, got, tt.want)
		}
	}
}
