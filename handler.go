// Package slogzap provides integration between the standard library's log/slog
// package and Uber's zap logging library.
package slogzap

import (
	"context"
	"log/slog"
	"math"
	"runtime"

	"github.com/takumakei/slogzap/levelconv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Handler is an implementation of [slog.Handler] that uses a [zap.Logger] as its backend.
type Handler struct {
	zap   *zap.Logger
	lvl   slog.Level
	limit slog.Level
}

var _ slog.Handler = (*Handler)(nil)

// New creates a new [slog.Handler] using the provided [zap.Logger] and options.
// It returns a [slog.Handler] that can be used with [slog.New].
//
// [levelconv.ToZap] is used to convert the [slog.Level] to [zapcore.Level].
//
// You can specify a limit using the [WithLimit] option.
// The limit represents the limit of the conversion.
// If the result of converting slog.Level is more serious than the limit, that level will be treated as the limit.
//
// While zap panics at a certain level (e.g. PanicLevel), slog has no such
// assumption. You can prevent panics by limiting the severe log level.
func New(logger *zap.Logger, options ...Option) slog.Handler {
	h := &Handler{
		zap:   logger,
		lvl:   levelconv.ToSlog(logger.Level()),
		limit: slog.Level(math.MaxInt),
	}
	for _, o := range options {
		o.apply(h)
	}
	return h
}

// Enabled implements [slog.Handler.Enabled].
// It returns true if the handler is enabled for the given level,
// determined by the underlying [zap.Logger] and limit option.
func (h *Handler) Enabled(_ context.Context, lvl slog.Level) bool {
	return h.lvl <= lvl
}

// Handle implements [slog.Handler.Handle].
// It processes the slog.Record and writes it to the underlying [zap.Logger].
func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	if ce := h.zap.Check(levelconv.ToZap(min(h.limit, r.Level)), r.Message); ce != nil {
		ce.Time = r.Time
		if f := runtime.FuncForPC(r.PC); f != nil {
			file, line := f.FileLine(r.PC)
			ce.Caller = zapcore.NewEntryCaller(r.PC, file, line, true)
		}
		fields := make([]zapcore.Field, 0, r.NumAttrs())
		r.Attrs(func(attr slog.Attr) bool {
			fields = append(fields, toField(attr))
			return true
		})
		ce.Write(fields...)
	}
	return nil
}

// WithAttrs implements [slog.Handler.WithAttrs].
// It returns a new Handler with the given attributes added to the logger.
func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) > 0 {
		fields := make([]zapcore.Field, len(attrs))
		for i, attr := range attrs {
			fields[i] = toField(attr)
		}
		h = h.clone()
		h.zap = h.zap.With(fields...)
	}
	return h
}

// WithGroup implements [slog.Handler.WithGroup].
// It returns a new Handler with the given group added to the logger.
func (h *Handler) WithGroup(name string) slog.Handler {
	h = h.clone()
	h.zap = h.zap.With(zap.Namespace(name))
	return h
}

func (h *Handler) clone() *Handler {
	o := *h
	return &o
}

// toField converts a slog.Attr to a zapcore.Field.
func toField(attr slog.Attr) zapcore.Field {
	return zap.Any(attr.Key, attr.Value.Any())
}

// Option is an interface for applying options to a [Handler].
type Option interface {
	apply(*Handler)
}

// WithLimit returns an Option that sets the maximum conversion level for the [Handler].
//
// See [New].
func WithLimit(zlvl zapcore.Level) Option {
	return withLimit(zlvl)
}

type withLimit zapcore.Level

func (o withLimit) apply(h *Handler) {
	h.limit = levelconv.ToSlog(zapcore.Level(o))
}
