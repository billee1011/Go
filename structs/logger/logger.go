package logger

type Fields map[string]interface{}

type Logger interface {
	WithField(key string, value interface{}) Logger
	WithFields(fields Fields) Logger

	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
}
