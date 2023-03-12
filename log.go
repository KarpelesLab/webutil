package webutil

import (
	"context"
	"fmt"
	"log"
)

type contextLogger struct {
	context.Context
	l *log.Logger
}

type contextLoggerValue int

const contextLoggerLogger contextLoggerValue = 1

func NewContextLogger(ctx context.Context, l *log.Logger) context.Context {
	return &contextLogger{ctx, l}
}

func (ctx *contextLogger) Value(key any) any {
	if key == contextLoggerLogger {
		return ctx.l
	}
	return ctx.Context.Value(key)
}

func CtxPrintf(ctx context.Context, format string, v ...any) {
	if l, ok := ctx.Value(contextLoggerLogger).(*log.Logger); ok {
		l.Output(2, fmt.Sprintf(format, v...))
	} else {
		log.Output(2, fmt.Sprintf(format, v...))
	}
}

func CtxSetPrefix(ctx context.Context, prefix string) {
	if l, ok := ctx.Value(contextLoggerLogger).(*log.Logger); ok {
		l.SetPrefix(prefix)
	}
}
