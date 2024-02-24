package logs

import (
	"fmt"
	"io"
	"log"
	"strings"
)

const (
	Info = iota
	Warning
	Error
	lastLogLevelLevel
)

type LogLevel int

type LogfFn = func(level LogLevel, format string, args ...interface{})

type Logger LogfFn

func (l Logger) I(format string, args ...interface{}) {
	l(Info, format, args...)
}

func (l Logger) W(format string, args ...interface{}) {
	l(Warning, format, args...)
}

func (l Logger) E(format string, args ...interface{}) {
	l(Error, format, args...)
}

func (l Logger) WithPrefix(prefix string) Logger {
	prefix += ": "
	return func(level LogLevel, format string, args ...interface{}) {
		builder := strings.Builder{}
		_, _ = fmt.Fprintf(&builder, prefix)
		_, _ = fmt.Fprintf(&builder, format, args...)
		l(level, "%s", builder.String())
	}
}

func (l Logger) WithCallerFunc() Logger {
	return l.WithPrefix(callerFuncName())
}

func DefaultLogger() Logger {
	return defaultLogf
}

func NoLogger() Logger {
	return noLogf
}

func noLogf(LogLevel, string, ...interface{}) {
}

func defaultLogf(level LogLevel, format string, args ...interface{}) {
	defaultLogFns[level](format, args...)
}

var defaultLogFns = []func(string, ...interface{}){
	logInfof,
	logWarningf,
	logErrorf,
}

func logInfof(format string, args ...interface{}) {
	_, _ = io.WriteString(log.Writer(), "[I]: ")
	log.Printf(format, args...)
}

func logWarningf(format string, args ...interface{}) {
	_, _ = io.WriteString(log.Writer(), "[W]: ")
	log.Printf(format, args...)
}

func logErrorf(format string, args ...interface{}) {
	_, _ = io.WriteString(log.Writer(), "[E]: ")
	log.Printf(format, args...)
}
