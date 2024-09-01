// Package levelconv provides functions to convert between slog.Level and zapcore.Level.
//
// [slog.Level] defines four constants: LevelDebug, LevelInfo, LevelWarn, and LevelError.
// Essentially, it's an int where higher values represent more severe events.
// There's a difference of 4 between each constant value.
//
// [zapcore.Level] is strictly defined by seven constants: DebugLevel, InfoLevel, WarnLevel,
// ErrorLevel, DPanicLevel, PanicLevel, and FatalLevel.
//
// [slog.Level] and [zapcore.Level] have common definitions for Debug, Info, Warn, and Error.
// This package uses these as the basis for defining the conversion.
//
// ### slog.Level
//
// https://pkg.go.dev/log/slog@go1.23.0#Level
//
//	slog.LevelDebug = -4
//	slog.LevelInfo  =  0
//	slog.LevelWarn  =  4
//	slog.LevelError =  8
//
// ### zapcore.Level
//
// https://pkg.go.dev/go.uber.org/zap@v1.27.0/zapcore#Level
//
//	zapcore.DebugLevel  = -1
//	zapcore.InfoLevel   = 0
//	zapcore.WarnLevel   = 1
//	zapcore.ErrorLevel  = 2
//	zapcore.DPanicLevel = 3
//	zapcore.PanicLevel  = 4
//	zapcore.FatalLevel  = 5
package levelconv

import (
	"log/slog"

	"go.uber.org/zap/zapcore"
)

var _ slog.Level = slog.LevelDebug
var _ zapcore.Level = zapcore.DebugLevel

// ToSlog converts a [zapcore.Level] to a [slog.Level].
//
// It maps [zapcore.InfoLevel] to [slog.LevelInfo] as the baseline.
// Using this conversion as a reference, it maps the 7 [zapcore.Level] constant values
// to [slog.Level] values with a difference of 4 between each level.
//
//	zapcore.DebugLevel  to slog.LevelDebug
//	zapcore.InfoLevel   to slog.LevelInfo
//	zapcore.WarnLevel   to slog.LevelWarn
//	zapcore.ErrorLevel  to slog.LevelError
//	zapcore.DPanicLevel to slog.LevelError + 4
//	zapcore.PanicLevel  to slog.LevelError + 8
//	zapcore.FatalLevel  to slog.LevelError + 12
func ToSlog(zlvl zapcore.Level) slog.Level {
	return toSlog[zlvl]
}

var toSlog = map[zapcore.Level]slog.Level{
	zapcore.DebugLevel:  slog.LevelDebug,
	zapcore.InfoLevel:   slog.LevelInfo,
	zapcore.WarnLevel:   slog.LevelWarn,
	zapcore.ErrorLevel:  slog.LevelError,
	zapcore.DPanicLevel: slog.LevelError + 4,
	zapcore.PanicLevel:  slog.LevelError + 8,
	zapcore.FatalLevel:  slog.LevelError + 12,
}

// ToZap converts a [slog.Level] to a [zapcore.Level].
//
// First, it maps values greater than or equal to [slog.LevelInfo] and less than [slog.LevelWarn]
// to [zapcore.InfoLevel].
//
// Then, since the difference between the 4 [slog.Level] constants is 4, it maps [slog.Level]
// values to [zapcore.Level] constants for every difference of 4.
//
// Finally, as [slog.Level] upper and lower bounds are within the int range, it rounds all
// values to the upper and lower bounds of [zapcore.Level].
//
//	All slog.Level values less than slog.LevelInfo (0) to zapcore.DebugLevel
//	slog.LevelInfo (0) to less than slog.LevelWarn (4) to zapcore.InfoLevel
//	slog.LevelWarn (4) to less than slog.LevelError (8) to zapcore.WarnLevel
//	slog.LevelError (8) to less than slog.LevelError (8) + 4 to zapcore.ErrorLevel
//	slog.LevelError (8) + 4 to less than slog.LevelError (8) + 8 to zapcore.DPanicLevel
//	slog.LevelError (8) + 8 to less than slog.LevelError (8) + 12 to zapcore.PanicLevel
//	slog.LevelError (8) + 12 and above to zapcore.FatalLevel
func ToZap(slvl slog.Level) zapcore.Level {
	ilvl := int(slvl)
	for _, it := range toZap {
		if ilvl < it.ilvl {
			return it.zlvl
		}
	}
	return zapcore.FatalLevel
}

var toZap = []*struct {
	ilvl int
	zlvl zapcore.Level
}{
	{ilvl: 4 * 0, zlvl: zapcore.DebugLevel},
	{ilvl: 4 * 1, zlvl: zapcore.InfoLevel},
	{ilvl: 4 * 2, zlvl: zapcore.WarnLevel},
	{ilvl: 4 * 3, zlvl: zapcore.ErrorLevel},
	{ilvl: 4 * 4, zlvl: zapcore.DPanicLevel},
	{ilvl: 4 * 5, zlvl: zapcore.PanicLevel},
}
