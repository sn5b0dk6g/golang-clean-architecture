package presenter

import (
	"go-rest-api/domain"
	"go-rest-api/usecase"
	"time"
)

type createTaskPresenter struct{}

func NewCreateTaskPresenter() usecase.CreateTaskPresenter {
	return createTaskPresenter{}
}

func (t createTaskPresenter) Output(task domain.Task) usecase.CreateTaskOutput {
	return usecase.CreateTaskOutput{
		ID:        task.ID,
		Title:     task.Title,
		CreatedAt: task.CreatedAt.Format(time.DateTime),
		UpdatedAt: task.UpdatedAt.Format(time.DateTime),
	}
}
