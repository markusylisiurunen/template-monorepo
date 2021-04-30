package logger

import (
	"fmt"
)

type NullLoggerEntry struct {
	Level   string
	Message string
	Meta    []interface{}
}

type NullLogger struct {
	Logs []NullLoggerEntry
}

func (l *NullLogger) Debugf(template string, args ...interface{}) {
	l.log(NullLoggerEntry{
		Level:   "debug",
		Message: fmt.Sprintf(template, args...),
	})
}

func (l *NullLogger) Debugw(msg string, keysAndValues ...interface{}) {
	meta := []interface{}{}
	meta = append(meta, keysAndValues...)

	l.log(NullLoggerEntry{
		Level:   "",
		Message: msg,
		Meta:    meta,
	})
}

func (l *NullLogger) Infof(template string, args ...interface{}) {
	l.log(NullLoggerEntry{
		Level:   "info",
		Message: fmt.Sprintf(template, args...),
	})
}

func (l *NullLogger) Infow(msg string, keysAndValues ...interface{}) {
	meta := []interface{}{}
	meta = append(meta, keysAndValues...)

	l.log(NullLoggerEntry{
		Level:   "",
		Message: msg,
		Meta:    meta,
	})
}

func (l *NullLogger) Warnf(template string, args ...interface{}) {
	l.log(NullLoggerEntry{
		Level:   "warn",
		Message: fmt.Sprintf(template, args...),
	})
}

func (l *NullLogger) Warnw(msg string, keysAndValues ...interface{}) {
	meta := []interface{}{}
	meta = append(meta, keysAndValues...)

	l.log(NullLoggerEntry{
		Level:   "",
		Message: msg,
		Meta:    meta,
	})
}

func (l *NullLogger) Errorf(template string, args ...interface{}) {
	l.log(NullLoggerEntry{
		Level:   "error",
		Message: fmt.Sprintf(template, args...),
	})
}

func (l *NullLogger) Errorw(msg string, keysAndValues ...interface{}) {
	meta := []interface{}{}
	meta = append(meta, keysAndValues...)

	l.log(NullLoggerEntry{
		Level:   "",
		Message: msg,
		Meta:    meta,
	})
}

func (l *NullLogger) Clear() {
	l.Logs = []NullLoggerEntry{}
}

func (l *NullLogger) log(entry NullLoggerEntry) {
	l.Logs = append(l.Logs, entry)
}

func NewNullLogger() Logger {
	return &NullLogger{
		Logs: []NullLoggerEntry{},
	}
}
