package logger

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"strings"
)

// Source describes the location of a line of source code.
type Source struct {
	// Function is the package path-qualified function name containing the
	// source line. If non-empty, this string uniquely identifies a single
	// function in the program. This may be the empty string if not known.
	Function string `json:"function"`
	// File and Line are the file name and line number (1-based) of the source
	// line. These may be the empty string and zero, respectively, if not known.
	File string `json:"file"`
	Line int    `json:"line"`
}

func (s *Source) String() string {
	return fmt.Sprintf("%s:%d %s", s.File, s.Line, s.Function)
}

// sourceHandler 自定义 source 处理器，可以获取更上层的调用信息
type sourceHandler struct {
	slog.Handler
	callerSkip int // 跳过的调用层数
}

func newSourceHandler(handler slog.Handler, callerSkip int) *sourceHandler {
	return &sourceHandler{
		Handler:    handler,
		callerSkip: callerSkip,
	}
}

func (h *sourceHandler) Handle(ctx context.Context, r slog.Record) error {
	// 获取指定深度的调用信息
	pc, file, line, ok := runtime.Caller(4 + h.callerSkip) // 4 是 slog 内部调用栈的偏移
	if ok {
		// 替换原有的 source 信息
		funcName := runtime.FuncForPC(pc).Name()

		// 简化函数名，只保留包名和函数名
		if idx := strings.LastIndex(funcName, "/"); idx != -1 {
			funcName = funcName[idx+1:]
		}

		source := Source{
			File:     file,
			Line:     line,
			Function: funcName,
		}

		r.AddAttrs(slog.Any("source", source))
	}

	return h.Handler.Handle(ctx, r)
}

func (h *sourceHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return newSourceHandler(h.Handler.WithAttrs(attrs), h.callerSkip)
}

func (h *sourceHandler) WithGroup(name string) slog.Handler {
	return newSourceHandler(h.Handler.WithGroup(name), h.callerSkip)
}
