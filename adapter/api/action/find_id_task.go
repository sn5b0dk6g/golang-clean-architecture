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

type FindIdTaskAction struct {
	log            logger.Logger
	uc             usecase.FindIdTaskUsecase
	validator      validator.Validator
	logKey, logMsg string
}

func NewFindIdTaskAction(
	log logger.Logger,
	uc usecase.FindIdTaskUsecase,
	v validator.Validator,
) FindIdTaskAction {
	return FindIdTaskAction{
		uc:        uc,
		log:       log,
		validator: v,
		logKey:    "find_id_task",
		logMsg:    "find a task by id",
	}
}

func (t FindIdTaskAction) Execute(c echo.Context) error {
	taskId, err := strconv.Atoi(c.Param("taskId"))
	if err != nil {
		logging.NewError(t.log, err, t.logKey, http.StatusBadRequest).Log(t.logMsg)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	userId := utility.GetUserIdByToken(c)
	input := usecase.FindIdTaskInput{
		UserID: userId,
		TaskID: uint(taskId),
	}

	output, err := t.uc.Execute(input)
	if err != nil {
		logging.NewError(t.log, err, t.logKey, http.StatusInternalServerError).Log(t.logMsg)
		return response.NewError(err, http.StatusInternalServerError).SendJSON(c)
	}

	logging.NewInfo(t.log, t.logKey, http.StatusOK).Log(t.logMsg)

	return c.JSON(http.StatusOK, output)
}
