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

type DeleteTaskAction struct {
	log            logger.Logger
	uc             usecase.DeleteTaskUsecase
	validator      validator.Validator
	logKey, logMsg string
}

func NewDeleteTaskAction(
	log logger.Logger,
	uc usecase.DeleteTaskUsecase,
	v validator.Validator,
) DeleteTaskAction {
	return DeleteTaskAction{
		uc:        uc,
		log:       log,
		validator: v,
		logKey:    "delete_task",
		logMsg:    "delete a task",
	}
}

func (t DeleteTaskAction) Execute(c echo.Context) error {
	taskId, err := strconv.Atoi(c.Param("taskId"))
	if err != nil {
		logging.NewError(t.log, err, t.logKey, http.StatusBadRequest).Log(t.logMsg)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	userId := utility.GetUserIdByToken(c)
	input := usecase.DeleteTaskInput{
		UserID: userId,
		TaskID: uint(taskId),
	}

	err = t.uc.Execute(input)
	if err != nil {
		logging.NewError(t.log, err, t.logKey, http.StatusInternalServerError).Log(t.logMsg)
		return response.NewError(err, http.StatusInternalServerError).SendJSON(c)
	}

	logging.NewInfo(t.log, t.logKey, http.StatusNoContent).Log(t.logMsg)

	return c.NoContent(http.StatusNoContent)
}
