package action

import (
	"go-rest-api/adapter/api/logging"
	"go-rest-api/adapter/api/response"
	"go-rest-api/adapter/logger"
	"go-rest-api/adapter/validator"
	"go-rest-api/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CreateUserAction struct {
	log            logger.Logger
	uc             usecase.CreateUserUsecase
	validator      validator.Validator
	logKey, logMsg string
}

func NewCreateUserAction(
	uc usecase.CreateUserUsecase,
	log logger.Logger,
	v validator.Validator,
) CreateUserAction {
	return CreateUserAction{
		uc:        uc,
		log:       log,
		validator: v,
		logKey:    "create_user",
		logMsg:    "creating a new user",
	}
}

func (u CreateUserAction) Execute(c echo.Context) error {
	var input usecase.CreateUserInput
	if err := c.Bind(&input); err != nil {
		logging.NewError(u.log, err, u.logKey, http.StatusBadRequest).Log(u.logMsg)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if errs := u.validateInput(input); errs != nil {
		return response.NewErrorMessage(errs, http.StatusInternalServerError).SendJSON(c)
	}

	output, err := u.uc.Execute(input)
	if err != nil {
		logging.NewError(u.log, err, u.logKey, http.StatusInternalServerError).Log(u.logMsg)
		return response.NewError(err, http.StatusInternalServerError).SendJSON(c)
	}

	logging.NewInfo(u.log, u.logKey, http.StatusCreated).Log(u.logMsg)

	return c.JSON(http.StatusCreated, output)
}

func (u CreateUserAction) validateInput(input usecase.CreateUserInput) []string {
	var errs []string
	if err := u.validator.Validate(input); err != nil {
		errs = u.validator.Messages()
	}
	return errs
}
