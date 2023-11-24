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

type AuthenticationUserAction struct {
	log            logger.Logger
	uc             usecase.AuthenticationUserUsecase
	validator      validator.Validator
	logKey, logMsg string
}

func NewAuthenticationUserAction(
	uc usecase.AuthenticationUserUsecase,
	log logger.Logger,
	v validator.Validator,
) AuthenticationUserAction {
	return AuthenticationUserAction{
		uc:        uc,
		log:       log,
		validator: v,
		logKey:    "authentication_user",
		logMsg:    "authentication a user",
	}
}

func (u AuthenticationUserAction) Execute(c echo.Context) error {
	var input usecase.AuthenticationUserInput
	if err := c.Bind(&input); err != nil {
		logging.NewError(u.log, response.ErrParameterInvalid, u.logKey, http.StatusBadRequest).Log(u.logMsg)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if errs := u.validateInput(input); errs != nil {
		logging.NewError(u.log, response.ErrInvalidInput, u.logKey, http.StatusInternalServerError).Log(u.logMsg)
		return response.NewErrorMessage(errs, http.StatusInternalServerError).SendJSON(c)
	}

	output, err := u.uc.Execute(input)
	if err != nil {
		logging.NewError(u.log, err, u.logKey, http.StatusInternalServerError).Log(u.logMsg)
		return response.NewError(err, http.StatusInternalServerError).SendJSON(c)
	}

	logging.NewInfo(u.log, u.logKey, http.StatusCreated).Log(u.logMsg)

	// cookieにトークンを追加する
	cookie := utility.CreateCookie(output.JwtTokenString, time.Now().Add(24*time.Hour))
	c.SetCookie(cookie)

	return c.NoContent(http.StatusOK)
}

func (u AuthenticationUserAction) validateInput(input usecase.AuthenticationUserInput) []string {
	var errs []string
	if err := u.validator.Validate(input); err != nil {
		errs = u.validator.Messages()
	}
	return errs
}
