package log

import (
	"fmt"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

type LogInfo map[string]interface{}

var logger zerolog.Logger

func init() {
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}

	fileName := fmt.Sprintf("./data/logs/app-%s.log", time.Now().Format("2006-01-02"))
	fileWriter := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    12,
		MaxAge:     28,
		MaxBackups: 1,
		LocalTime:  true,
		Compress:   true,
	}

	multi := zerolog.MultiLevelWriter(consoleWriter, fileWriter)

	logger = zerolog.New(multi).With().Timestamp().Logger()
}

func GetLogger() *zerolog.Logger {
	return &logger
}

func Debug(fields LogInfo, msg string) {
	var conv map[string]interface{} = fields

	event := logger.Debug()
	for key, value := range conv {
		event = event.Interface(key, value)
	}
	event.Msg(msg)
}

func Error(fields LogInfo, msg string) {
	var conv map[string]interface{} = fields

	event := logger.Error()
	for key, value := range conv {
		event = event.Interface(key, value)
	}
	event.Msg(msg)
}

func Fatal(fields LogInfo, msg string) {
	var conv map[string]interface{} = fields

	event := logger.Fatal()
	for key, value := range conv {
		event = event.Interface(key, value)
	}
	event.Msg(msg)
}

func Panic(fields LogInfo, msg string) {
	var conv map[string]interface{} = fields

	event := logger.Panic()
	for key, value := range conv {
		event = event.Interface(key, value)
	}
	event.Msg(msg)
}

func Info(fields LogInfo, msg string) {
	var conv map[string]interface{} = fields

	event := logger.Debug()
	for key, value := range conv {
		event = event.Interface(key, value)
	}
	event.Msg(msg)
}
