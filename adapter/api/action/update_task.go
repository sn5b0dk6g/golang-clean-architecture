package action

import (
	"go-rest-api/adapter/api/logging"
	"go-rest-api/adapter/api/response"
	"go-rest-api/adapter/logger"
	"go-rest-api/adapter/validator"
	"go-rest-api/usecase"
	"go-rest-api/utility"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type UpdateTaskAction struct {
	log            logger.Logger
	uc             usecase.UpdateTaskUsecase
	validator      validator.Validator
	logKey, logMsg string
}

func NewUpdateTaskAction(
	log logger.Logger,
	uc usecase.UpdateTaskUsecase,
	v validator.Validator,
) UpdateTaskAction {
	return UpdateTaskAction{
		uc:        uc,
		log:       log,
		validator: v,
		logKey:    "update_task",
		logMsg:    "update a task",
	}
}

func (t UpdateTaskAction) Execute(c echo.Context) error {
	var input usecase.UpdateTaskInput
	if err := c.Bind(&input); err != nil {
		logging.NewError(t.log, err, t.logKey, http.StatusBadRequest).Log(t.logMsg)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	taskId, err := strconv.Atoi(c.Param("taskId"))
	if err != nil {
		logging.NewError(t.log, err, t.logKey, http.StatusBadRequest).Log(t.logMsg)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	userId := utility.GetUserIdByToken(c)
	input.TaskID = uint(taskId)
	input.UserID = userId

	output, err := t.uc.Execute(input)
	if err != nil {
		logging.NewError(t.log, err, t.logKey, http.StatusInternalServerError).Log(t.logMsg)
		return response.NewError(err, http.StatusInternalServerError).SendJSON(c)
	}

	logging.NewInfo(t.log, t.logKey, http.StatusOK).Log(t.logMsg)

	return c.JSON(http.StatusOK, output)
}
