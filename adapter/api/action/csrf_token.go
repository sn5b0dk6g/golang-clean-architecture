package action

import (
	"go-rest-api/adapter/logger"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CsrfTokenAction struct {
	log            logger.Logger
	logKey, logMsg string
}

func NewCsrfTokenAction(
	log logger.Logger,
) CsrfTokenAction {
	return CsrfTokenAction{
		log:    log,
		logKey: "get_csrf_token",
		logMsg: "return csrf token",
	}
}

func (t CsrfTokenAction) Execute(c echo.Context) error {
	token := c.Get("csrf").(string)
	return c.JSON(http.StatusOK, echo.Map{
		"csrf_token": token,
	})
}
