package action

import (
	"go-rest-api/adapter/api/logging"
	"go-rest-api/adapter/api/response"
	"go-rest-api/adapter/logger"
	"go-rest-api/adapter/validator"
	"go-rest-api/usecase"
	"go-rest-api/utility"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type LogoutUserAction struct {
	log            logger.Logger
	uc             usecase.LogoutUserUsecase
	validator      validator.Validator
	logKey, logMsg string
}

func NewLogoutUserAction(
	uc usecase.LogoutUserUsecase,
	log logger.Logger,
	v validator.Validator,
) LogoutUserAction {
	return LogoutUserAction{
		uc:        uc,
		log:       log,
		validator: v,
		logKey:    "logout_user",
		logMsg:    "logout a user",
	}
}

func (u LogoutUserAction) Execute(c echo.Context) error {
	var input usecase.LogoutUserInput
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if errs := u.validateInput(input); errs != nil {
		return response.NewErrorMessage(errs, http.StatusInternalServerError).SendJSON(c)
	}

	err := u.uc.Execute(input)
	if err != nil {
		logging.NewError(u.log, err, u.logKey, http.StatusInternalServerError).Log(u.logMsg)
		return response.NewError(err, http.StatusInternalServerError).SendJSON(c)
	}

	cookie := utility.CreateCookie("", time.Now())
	c.SetCookie(cookie)

	logging.NewInfo(u.log, u.logKey, http.StatusOK).Log(u.logMsg)

	return c.NoContent(http.StatusOK)
}

func (u LogoutUserAction) validateInput(input usecase.LogoutUserInput) []string {
	var errs []string
	if err := u.validator.Validate(input); err != nil {
		errs = u.validator.Messages()
	}
	return errs
}
