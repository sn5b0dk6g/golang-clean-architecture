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

type FindAllTaskAction struct {
	log            logger.Logger
	uc             usecase.FindAllTaskUsecase
	validator      validator.Validator
	logKey, logMsg string
}

func NewFindAllTaskAction(
	log logger.Logger,
	uc usecase.FindAllTaskUsecase,
	v validator.Validator,
) FindAllTaskAction {
	return FindAllTaskAction{
		uc:        uc,
		log:       log,
		validator: v,
		logKey:    "find_all_task",
		logMsg:    "find all task",
	}
}

func (t FindAllTaskAction) Execute(c echo.Context) error {
	userId := utility.GetUserIdByToken(c)
	input := usecase.FindAllTaskInput{UserID: userId}

	output, err := t.uc.Execute(input)
	if err != nil {
		logging.NewError(t.log, err, t.logKey, http.StatusInternalServerError).Log(t.logMsg)
		return response.NewError(err, http.StatusInternalServerError).SendJSON(c)
	}

	logging.NewInfo(t.log, t.logKey, http.StatusOK).Log(t.logMsg)

	return c.JSON(http.StatusOK, output)
}
