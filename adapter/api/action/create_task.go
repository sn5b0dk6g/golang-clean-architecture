package action

import (
	"go-rest-api/adapter/api/logging"
	"go-rest-api/adapter/api/response"
	"go-rest-api/adapter/logger"
	"go-rest-api/adapter/validator"
	"go-rest-api/usecase"
	"go-rest-api/utility"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CreateTaskAction struct {
	log            logger.Logger
	uc             usecase.CreateTaskUsecase
	validator      validator.Validator
	logKey, logMsg string
}

func NewCreateTaskAction(
	log logger.Logger,
	uc usecase.CreateTaskUsecase,
	v validator.Validator,
) CreateTaskAction {
	return CreateTaskAction{
		uc:        uc,
		log:       log,
		validator: v,
		logKey:    "create_task",
		logMsg:    "create a new task",
	}
}

func (t CreateTaskAction) Execute(c echo.Context) error {
	var input usecase.CreateTaskInput
	if err := c.Bind(&input); err != nil {
		logging.NewError(t.log, err, t.logKey, http.StatusBadRequest).Log(t.logMsg)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	userId := utility.GetUserIdByToken(c)
	input.UserID = userId

	output, err := t.uc.Execute(input)
	if err != nil {
		logging.NewError(t.log, err, t.logKey, http.StatusInternalServerError).Log(t.logMsg)
		return response.NewError(err, http.StatusInternalServerError).SendJSON(c)
	}

	logging.NewInfo(t.log, t.logKey, http.StatusCreated).Log(t.logMsg)

	return c.JSON(http.StatusCreated, output)
}
