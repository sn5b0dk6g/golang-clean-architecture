package router

import (
	"go-rest-api/adapter/api/action"
	"go-rest-api/adapter/logger"
	"go-rest-api/adapter/presenter"
	"go-rest-api/adapter/repository"
	"go-rest-api/domain"
	"go-rest-api/usecase"

	"github.com/labstack/echo/v4"
)

func BuildCreateUserAction(e echoServer) echo.HandlerFunc {
	return func(c echo.Context) error {
		var db domain.UserRepository
		if e.useDB {
			db = repository.NewUserSQL(e.dbSQL)
		} else {
			db = repository.NewUserRedis(e.dbNoSQL, e.redisExp)
		}
		requestID := c.Response().Header().Get(echo.HeaderXRequestID)
		rLog := e.log.WithFields(logger.Fields{"id": requestID})
		var (
			uc = usecase.NewCreateUserInteractor(
				db,
				presenter.NewCreateUserPresenter(),
			)
			act = action.NewCreateUserAction(uc, rLog, e.validator)
		)
		return act.Execute(c)
	}
}

func BuildAuthenticationUserAction(e echoServer) echo.HandlerFunc {
	return func(c echo.Context) error {
		requestID := c.Response().Header().Get(echo.HeaderXRequestID)
		rLog := e.log.WithFields(logger.Fields{"id": requestID})
		var (
			uc = usecase.NewAuthenticationUserInteractor(
				repository.NewUserSQL(e.dbSQL),
				repository.NewUserRedis(e.dbNoSQL, e.redisExp),
				presenter.NewAuthenticationUserPresenter(),
			)
			act = action.NewAuthenticationUserAction(uc, rLog, e.validator)
		)
		return act.Execute(c)
	}
}

func BuildLogoutUserAction(e echoServer) echo.HandlerFunc {
	return func(c echo.Context) error {
		requestID := c.Response().Header().Get(echo.HeaderXRequestID)
		rLog := e.log.WithFields(logger.Fields{"id": requestID})
		var (
			uc = usecase.NewLogoutUserInteractor(
				repository.NewUserRedis(e.dbNoSQL, e.redisExp),
			)
			act = action.NewLogoutUserAction(uc, rLog, e.validator)
		)
		return act.Execute(c)
	}
}
