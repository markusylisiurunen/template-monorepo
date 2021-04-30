package logger

type Logger interface {
	Debugf(template string, args ...interface{})
	Debugw(msg string, keysAndValues ...interface{})

	Infof(template string, args ...interface{})
	Infow(msg string, keysAndValues ...interface{})

	Warnf(template string, args ...interface{})
	Warnw(msg string, keysAndValues ...interface{})

	Errorf(template string, args ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
}
