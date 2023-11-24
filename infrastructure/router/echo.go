package router

import (
	"fmt"
	"go-rest-api/adapter/api/action"
	"go-rest-api/adapter/logger"
	"go-rest-api/adapter/validator"
	"go-rest-api/utility"
	"net/http"
	"os"
	"time"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	e.router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000", os.Getenv("FE_URL")},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept,
			echo.HeaderAccessControlAllowHeaders, echo.HeaderXCSRFToken},
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE"},
		AllowCredentials: true,
	}))
	e.router.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		CookiePath:     "/",
		CookieDomain:   os.Getenv("API_DOMAIN"),
		CookieHTTPOnly: true,
		CookieSameSite: http.SameSiteDefaultMode,
	}))

	// ログレベルの変更
	// if l, ok := e.router.Logger.(*log.Logger); ok {
	// 	l.SetLevel(log.INFO)
	// }

	e.setAppHandlers()
	e.log.WithFields(logger.Fields{"port": e.port}).Infof("Starting HTTP Server")
	e.log.Fatalln(e.router.Start(fmt.Sprintf(":%v", e.port)))
}

func (e echoServer) setAppHandlers() {
	e.router.GET("/csrf", e.BuildCsrfTokenAction())
	v1 := e.router.Group("/v1")
	v1.POST("/signup", BuildCreateUserAction(e))
	v1.POST("/login", BuildAuthenticationUserAction(e))
	v1.POST("/logout", BuildLogoutUserAction(e))

	tasks := v1.Group("/tasks")
	tasks.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(os.Getenv("SECRET")), // 作成したときと同じSECRET
		TokenLookup: "cookie:token",              // どこに入っているか
	}))

	tasks.GET("", BuildFindAllTaskAction(e))
	tasks.GET("/:taskId", BuildFindIdTaskAction(e))
	tasks.POST("", BuildCreateTaskAction(e))
	tasks.POST("/:taskId", BuildUpdateTaskAction(e))
	tasks.DELETE("/:taskId", BuildDeleteTaskAction(e))
}

func (e echoServer) BuildCsrfTokenAction() echo.HandlerFunc {
	return func(c echo.Context) error {
		requestID := c.Response().Header().Get(echo.HeaderXRequestID)
		rLog := e.log.WithFields(logger.Fields{"id": requestID})
		act := action.NewCsrfTokenAction(rLog)
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
