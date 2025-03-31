package log

type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelDanger
	LogLevelFatal
)

var level LogLevel = LogLevelDebug

func SetLogLevel(l LogLevel) {
	level = l
}
