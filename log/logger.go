package log

type Logger interface {
	Debug(message string)
	Info(message string)
	Warn(message string)
	Error(error error, message string)
}
