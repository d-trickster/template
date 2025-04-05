package logging

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
)

type TerminalHandler struct {
	w     io.Writer
	level slog.Level
	attrs []slog.Attr
}

func NewTerminalHandler(w io.Writer, level slog.Level) *TerminalHandler {
	return &TerminalHandler{
		w:     w,
		level: level,
	}
}

func (h *TerminalHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *TerminalHandler) Handle(_ context.Context, r slog.Record) error {
	timeStr := r.Time.Format("2006-01-02 15:04:05.000")

	levelStr := strings.ToUpper(r.Level.String())
	switch levelStr {
	case "INFO":
		levelStr = fmt.Sprintf("\033[32m%s\033[0m", levelStr)
	case "DEBUG":
	case "WARN":
		levelStr = fmt.Sprintf("\033[33m%s\033[0m", levelStr)
	case "ERROR":
		levelStr = fmt.Sprintf("\033[31m%s\033[0m", levelStr)
	}

	msg := fmt.Sprintf("\033[90m%s\033[0m [%s] %s", timeStr, levelStr, r.Message)

	attrs := make([]string, 0, len(h.attrs)+r.NumAttrs())
	for i := range h.attrs {
		attrs = append(attrs, fmt.Sprintf("%s=%v", h.attrs[i].Key, h.attrs[i].Value.Any()))
	}
	r.Attrs(func(a slog.Attr) bool {
		if a.Equal(slog.Attr{}) {
			return true
		}
		attrs = append(attrs, fmt.Sprintf("%s=%v", a.Key, a.Value.Any()))
		return true
	})

	if len(attrs) > 0 {
		msg += fmt.Sprintf(" \033[90m%s\033[0m", strings.Join(attrs, " "))
	}

	_, err := fmt.Fprintln(h.w, msg)
	return err
}

func (h *TerminalHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handler := NewTerminalHandler(h.w, h.level)
	handler.attrs = append(handler.attrs, attrs...)
	return handler
}

func (h *TerminalHandler) WithGroup(name string) slog.Handler {
	return h
}
