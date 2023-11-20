package log

import (
	"errors"
	"go-rest-api/adapter/logger"
)

const (
	InstanceZapLogger int = iota
	InstanceSlogLogger
)

var (
	errInvalidLoggerInstance = errors.New("invalid log instance")
)

func NewLoggerFactory(instance int) (logger.Logger, error) {
	switch instance {
	case InstanceZapLogger:
		return NewZapLogger()
	case InstanceSlogLogger:
		return NewSlogLogger(), nil
	default:
		return nil, errInvalidLoggerInstance
	}
}
