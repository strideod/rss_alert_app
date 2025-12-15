package log

import (
	"log/slog"
	"os"
)

// var (
// 	infoLogger *log.Logger
// 	warningLogger *log.Logger
// 	errorLogger *log.Logger
// )

// func Init() {
// 	infoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
// 	warningLogger = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
// 	errorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
// }

// func Info(msg string) {
// 	infoLogger.Output(2, msg)
// }

// func Warning(msg string) {
// 	warningLogger.Output(2, msg)
// }

// func Error(msg string) {
// 	errorLogger.Output(2, msg)
// }

var Logger *slog.Logger
var Debug bool

func Init() {
	// Default json handler
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level: slog.LevelInfo,
	})

	Logger = slog.New(handler)
}

func SetDebug() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	})

	Logger = slog.New(handler)
}

func LogInfo(msg string, keysAndValues ...any) {
	Logger.Info(msg, keysAndValues...)
}

func LogDebug(msg string, keysAndValues ...any) {
	Logger.Debug(msg, keysAndValues...)
}

func LogWarning(msg string, keysAndValues ...any) {
	Logger.Warn(msg, keysAndValues...)
}

func LogError(msg string, keysAndValues ...any) {
	Logger.Error(msg, keysAndValues...)
}