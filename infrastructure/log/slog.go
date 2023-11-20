package log

import (
	"go-rest-api/adapter/logger"
	"log/slog"
	"os"
)

type slogLogger struct {
	logger *slog.Logger
}

func NewSlogLogger() logger.Logger {
	slogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	return &slogLogger{logger: slogger}
}

func (l *slogLogger) Infof(format string, args ...interface{}) {
	l.logger.Info(format, args...)
}

func (l *slogLogger) Warnf(format string, args ...interface{}) {
	l.logger.Warn(format, args...)
}

func (l *slogLogger) Errorf(format string, args ...interface{}) {
	l.logger.Error(format, args...)
}

func (l *slogLogger) Fatalln(args ...interface{}) {
	l.logger.Error("", args...)
	os.Exit(1)
}

func (l *slogLogger) WithFields(fields logger.Fields) logger.Logger {
	var f = make([]interface{}, 0)
	for index, field := range fields {
		f = append(f, index)
		f = append(f, field)
	}

	log := l.logger.With(f...)
	return &slogLogger{logger: log}
}

func (l *slogLogger) WithIndexFields(fields logger.Fields, indexs ...string) logger.Logger {
	var f = make([]interface{}, 0)
	for _, key := range indexs {
		f = append(f, key)
		f = append(f, fields[key])
	}
	log := l.logger.With(f...)
	return &slogLogger{logger: log}
}

func (l *slogLogger) WithError(err error) logger.Logger {
	var log = l.logger.With("error", err.Error())
	return &slogLogger{logger: log}
}
