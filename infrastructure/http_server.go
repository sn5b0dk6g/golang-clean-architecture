package infrastructure

import (
	"go-rest-api/adapter/logger"
	"go-rest-api/adapter/validator"
	"go-rest-api/infrastructure/database"
	"go-rest-api/infrastructure/log"
	"go-rest-api/infrastructure/router"
	"go-rest-api/infrastructure/validation"
	"strconv"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Server interface {
	Listen()
}

type config struct {
	appName       string
	logger        logger.Logger
	validator     validator.Validator
	dbSQL         *gorm.DB
	dbNoSQL       *redis.Client
	webServerPort router.Port
	webServer     Server
}

func NewConfig() *config {
	return &config{}
}

func (c *config) Name(name string) *config {
	c.appName = name
	return c
}

func (c *config) Logger(instance int) *config {
	log, err := log.NewLoggerFactory(instance)
	if err != nil {
		log.Fatalln(err)
	}

	c.logger = log
	c.logger.Infof("Successfully configured log")
	return c
}

func (c *config) DbSQL() *config {
	db, err := database.OpenPostgres()
	if err != nil {
		c.logger.Fatalln(err, "Could not make a connection to the database")
	}

	c.logger.Infof("Successfully connected to the SQL database")

	c.dbSQL = db
	return c
}

func (c *config) DbNoSQL() *config {
	db, err := database.OpenRedis()
	if err != nil {
		c.logger.Fatalln(err, "Could not make a connection to the database")
	}

	c.logger.Infof("Successfully connected to the NoSQL database")

	c.dbNoSQL = db
	return c
}

func (c *config) Validator() *config {
	v, err := validation.NewOzzo()
	if err != nil {
		c.logger.Fatalln(err)
	}

	c.logger.Infof("Successfully configured validator")

	c.validator = v
	return c
}

func (c *config) WebServerPort(port string) *config {
	p, err := strconv.ParseInt(port, 10, 64)
	if err != nil {
		c.logger.Fatalln(err)
	}

	c.webServerPort = router.Port(p)
	return c
}

func (c *config) WebServer() *config {
	s := router.NewEchoServer(c.logger, c.dbSQL, c.dbNoSQL, true, c.webServerPort, c.validator)

	c.logger.Infof("Successfully configured router server")

	c.webServer = s
	return c
}

func (c *config) Start() {
	c.webServer.Listen()
}
