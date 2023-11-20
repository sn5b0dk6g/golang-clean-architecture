package router

import (
	"fmt"
	"go-rest-api/adapter/api/action"
	"go-rest-api/adapter/logger"
	"go-rest-api/adapter/presenter"
	repositoryNoSQL "go-rest-api/adapter/repository/nosql"
	repositorySQL "go-rest-api/adapter/repository/sql"
	"go-rest-api/adapter/validator"
	"go-rest-api/domain"
	"go-rest-api/usecase"
	"go-rest-api/utility"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Port int64

type echoServer struct {
	router    *echo.Echo
	log       logger.Logger
	dbSQL     *gorm.DB
	dbNoSQL   *redis.Client
	useDB     bool
	port      Port
	redisExp  time.Duration
	validator validator.Validator
}

func NewEchoServer(
	log logger.Logger,
	dbSQL *gorm.DB,
	dbNoSQL *redis.Client,
	useDB bool,
	port Port,
	validator validator.Validator,
) *echoServer {
	e := echo.New()
	return &echoServer{
		router:    e,
		log:       log,
		dbSQL:     dbSQL,
		dbNoSQL:   dbNoSQL,
		useDB:     useDB,
		port:      port,
		redisExp:  utility.GetRedisExpiration(),
		validator: validator,
	}
}

func (e echoServer) Listen() {
	e.router.Use(middleware.RequestID())
	//e.router.Use(middleware.Logger())
	e.router.Use(middleware.Recover())
	e.router.Use(middleware.RequestLoggerWithConfig(e.getRequestLoggerWithConfig()))

	// ログレベルの変更
	if l, ok := e.router.Logger.(*log.Logger); ok {
		l.SetLevel(log.INFO)
	}

	e.setAppHandlers()
	e.log.WithFields(logger.Fields{"port": e.port}).Infof("Starting HTTP Server")
	//e.router.Logger.Infof("Starting HTTP Server Port: %v", e.port)
	e.log.Fatalln(e.router.Start(fmt.Sprintf(":%v", e.port)))
}

func (e echoServer) setAppHandlers() {
	v1 := e.router.Group("/v1")
	v1.POST("/signup", e.buildCreateUserAction())
	v1.POST("/login", e.buildAuthenticationUserAction())
	v1.POST("/logout", e.buildLogoutUserAction())
}

func (e echoServer) buildCreateUserAction() echo.HandlerFunc {
	return func(c echo.Context) error {
		var db domain.UserRepository
		if e.useDB {
			db = repositorySQL.NewUserSQL(e.dbSQL)
		} else {
			db = repositoryNoSQL.NewUserRedis(e.dbNoSQL, e.redisExp)
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

func (e echoServer) buildAuthenticationUserAction() echo.HandlerFunc {
	return func(c echo.Context) error {
		requestID := c.Response().Header().Get(echo.HeaderXRequestID)
		rLog := e.log.WithFields(logger.Fields{"id": requestID})
		var (
			uc = usecase.NewAuthenticationUserInteractor(
				repositorySQL.NewUserSQL(e.dbSQL),
				repositoryNoSQL.NewUserRedis(e.dbNoSQL, e.redisExp),
				presenter.NewAuthenticationUserPresenter(),
			)
			act = action.NewAuthenticationUserAction(uc, rLog, e.validator)
		)
		return act.Execute(c)
	}
}

func (e echoServer) buildLogoutUserAction() echo.HandlerFunc {
	return func(c echo.Context) error {
		requestID := c.Response().Header().Get(echo.HeaderXRequestID)
		rLog := e.log.WithFields(logger.Fields{"id": requestID})
		var (
			uc = usecase.NewLogoutUserInteractor(
				repositoryNoSQL.NewUserRedis(e.dbNoSQL, e.redisExp),
			)
			act = action.NewLogoutUserAction(uc, rLog, e.validator)
		)
		return act.Execute(c)
	}
}

func (e echoServer) getRequestLoggerWithConfig() middleware.RequestLoggerConfig {
	return middleware.RequestLoggerConfig{
		LogStatus:    true,
		LogURI:       true,
		LogError:     true,
		LogRequestID: true,
		LogRemoteIP:  true,
		LogHost:      true,
		LogMethod:    true,
		LogUserAgent: true,
		LogLatency:   true,
		HandleError:  true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			fields := logger.Fields{
				"id":            v.RequestID,
				"remote_ip":     v.RemoteIP,
				"host":          v.Host,
				"method":        v.Method,
				"uri":           v.URI,
				"user_agent":    v.UserAgent,
				"status":        v.Status,
				"latency":       v.Latency,
				"latency_human": v.Latency.String(),
			}
			keys := []string{
				"id",
				"remote_ip",
				"host",
				"method",
				"uri",
				"user_agent",
				"status",
				"latency",
				"latency_human",
			}
			if v.Error != nil {
				fields["err"] = v.Error.Error()
			}
			e.log.WithIndexFields(fields, keys...).Infof("REQUEST")
			return nil
		},
	}
}

// func (e echoServer) getRequestLogAttr(v middleware.RequestLoggerValues) []slog.Attr {
// 	return []slog.Attr{
// 		slog.String("id", v.RequestID),
// 		slog.String("remote_ip", v.RemoteIP),
// 		slog.String("host", v.Host),
// 		slog.String("method", v.Method),
// 		slog.String("uri", v.URI),
// 		slog.String("user_agent", v.UserAgent),
// 		slog.Int("status", v.Status),
// 		slog.String("latency_human", v.Latency.String()),
// 	}
// }
