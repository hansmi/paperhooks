package client

import (
	"fmt"
)

type stdLogger interface {
	Print(...any)
}

type Logger interface {
	Errorf(format string, v ...any)
	Warnf(format string, v ...any)
	Debugf(format string, v ...any)
}

type discardLogger struct{}

var _ Logger = (*discardLogger)(nil)

func (*discardLogger) Errorf(format string, v ...any) {
}

func (*discardLogger) Warnf(format string, v ...any) {
}

func (*discardLogger) Debugf(format string, v ...any) {
}

type wrappedStdLogger struct {
	stdLogger
}

var _ Logger = (*wrappedStdLogger)(nil)

func (l *wrappedStdLogger) log(prefix string, format string, v []any) {
	l.stdLogger.Print(prefix, fmt.Sprintf(format, v...))
}

func (l *wrappedStdLogger) Errorf(format string, v ...any) {
	l.log("[E] ", format, v)
}

func (l *wrappedStdLogger) Warnf(format string, v ...any) {
	l.log("[W] ", format, v)
}

func (l *wrappedStdLogger) Debugf(format string, v ...any) {
	l.log("[D] ", format, v)
}

type prefixLogger struct {
	wrapped Logger
	prefix  string
}

var _ Logger = (*prefixLogger)(nil)

func (l *prefixLogger) wrap(fn func(string, ...any), format string, v []any) {
	fn("%s%s", l.prefix, fmt.Sprintf(format, v...))
}

func (l *prefixLogger) Errorf(format string, v ...any) {
	l.wrap(l.wrapped.Errorf, format, v)
}

func (l *prefixLogger) Warnf(format string, v ...any) {
	l.wrap(l.wrapped.Warnf, format, v)
}

func (l *prefixLogger) Debugf(format string, v ...any) {
	l.wrap(l.wrapped.Debugf, format, v)
}
