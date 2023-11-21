package router

import (
	"go-rest-api/adapter/api/action"
	"go-rest-api/adapter/logger"
	"go-rest-api/adapter/presenter"
	"go-rest-api/adapter/repository"
	"go-rest-api/usecase"

	"github.com/labstack/echo/v4"
)

func BuildFindAllTaskAction(e echoServer) echo.HandlerFunc {
	return func(c echo.Context) error {
		requestID := c.Response().Header().Get(echo.HeaderXRequestID)
		rLog := e.log.WithFields(logger.Fields{"id": requestID})
		var (
			uc = usecase.NewFindAllTaskInteractor(
				repository.NewTaskSQL(e.dbSQL),
				presenter.NewFindAllTaskPresenter(),
			)
			act = action.NewFindAllTaskAction(rLog, uc, e.validator)
		)
		return act.Execute(c)
	}
}

func BuildFindIdTaskAction(e echoServer) echo.HandlerFunc {
	return func(c echo.Context) error {
		requestID := c.Response().Header().Get(echo.HeaderXRequestID)
		rLog := e.log.WithFields(logger.Fields{"id": requestID})
		var (
			uc = usecase.NewFindIdTaskInteractor(
				repository.NewTaskSQL(e.dbSQL),
				presenter.NewFindIdTaskPresenter(),
			)
			act = action.NewFindIdTaskAction(rLog, uc, e.validator)
		)
		return act.Execute(c)
	}
}

func BuildCreateTaskAction(e echoServer) echo.HandlerFunc {
	return func(c echo.Context) error {
		requestID := c.Response().Header().Get(echo.HeaderXRequestID)
		rLog := e.log.WithFields(logger.Fields{"id": requestID})
		var (
			uc = usecase.NewCreateTaskInteractor(
				repository.NewTaskSQL(e.dbSQL),
				presenter.NewCreateTaskPresenter(),
			)
			act = action.NewCreateTaskAction(rLog, uc, e.validator)
		)
		return act.Execute(c)
	}
}

func BuildUpdateTaskAction(e echoServer) echo.HandlerFunc {
	return func(c echo.Context) error {
		requestID := c.Response().Header().Get(echo.HeaderXRequestID)
		rLog := e.log.WithFields(logger.Fields{"id": requestID})
		var (
			uc = usecase.NewUpdateTaskInteractor(
				repository.NewTaskSQL(e.dbSQL),
				presenter.NewUpdateTaskPresenter(),
			)
			act = action.NewUpdateTaskAction(rLog, uc, e.validator)
		)
		return act.Execute(c)
	}
}

func BuildDeleteTaskAction(e echoServer) echo.HandlerFunc {
	return func(c echo.Context) error {
		requestID := c.Response().Header().Get(echo.HeaderXRequestID)
		rLog := e.log.WithFields(logger.Fields{"id": requestID})
		var (
			uc = usecase.NewDeleteTaskInteractor(
				repository.NewTaskSQL(e.dbSQL),
			)
			act = action.NewDeleteTaskAction(rLog, uc, e.validator)
		)
		return act.Execute(c)
	}
}
