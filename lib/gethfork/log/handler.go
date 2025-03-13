package log

import (
	"context"
	"io"
	"log/slog"
	"math"
	"sync"
)

const levelMaxVerbosity slog.Level = math.MinInt

type discardHandler struct{}

// DiscardHandler returns a no-op handler
func DiscardHandler() slog.Handler {
	return &discardHandler{}
}

func (h *discardHandler) Handle(_ context.Context, r slog.Record) error {
	return nil
}

func (h *discardHandler) Enabled(_ context.Context, level slog.Level) bool {
	return false
}

func (h *discardHandler) WithGroup(name string) slog.Handler {
	panic("not implemented")
}

func (h *discardHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &discardHandler{}
}

type TerminalHandler struct {
	mu       sync.Mutex
	wr       io.Writer
	lvl      slog.Level
	useColor bool
	attrs    []slog.Attr
	// fieldPadding is a map with maximum field value lengths seen until now
	// to allow padding log contexts in a bit smarter way.
	fieldPadding map[string]int

	buf []byte
}

// NewTerminalHandler returns a handler which formats log records at all levels optimized for human readability on
// a terminal with color-coded level output and terser human friendly timestamp.
// This format should only be used for interactive programs or while developing.
//
//	[LEVEL] [TIME] MESSAGE key=value key=value ...
//
// Example:
//
//	[DBUG] [May 16 20:58:45] remove route ns=haproxy addr=127.0.0.1:50002
func NewTerminalHandler(wr io.Writer, useColor bool) *TerminalHandler {
	return NewTerminalHandlerWithLevel(wr, levelMaxVerbosity, useColor)
}

// NewTerminalHandlerWithLevel returns the same handler as NewTerminalHandler but only outputs
// records which are less than or equal to the specified verbosity level.
func NewTerminalHandlerWithLevel(wr io.Writer, lvl slog.Level, useColor bool) *TerminalHandler {
	return &TerminalHandler{
		wr:           wr,
		lvl:          lvl,
		useColor:     useColor,
		fieldPadding: make(map[string]int),
	}
}

func (h *TerminalHandler) Handle(_ context.Context, r slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	buf := h.format(h.buf, r, h.useColor)
	h.wr.Write(buf)
	h.buf = buf[:0]
	return nil
}

func (h *TerminalHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.lvl
}

func (h *TerminalHandler) WithGroup(name string) slog.Handler {
	panic("not implemented")
}

func (h *TerminalHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &TerminalHandler{
		wr:           h.wr,
		lvl:          h.lvl,
		useColor:     h.useColor,
		attrs:        append(h.attrs, attrs...),
		fieldPadding: make(map[string]int),
	}
}

// ResetFieldPadding zeroes the field-padding for all attribute pairs.
func (h *TerminalHandler) ResetFieldPadding() {
	h.mu.Lock()
	h.fieldPadding = make(map[string]int)
	h.mu.Unlock()
}
