package presenter

import (
	"go-rest-api/domain"
	"go-rest-api/usecase"
	"time"
)

type updateTaskpPresenter struct{}

func NewUpdateTaskPresenter() usecase.UpdateTaskPresenter {
	return updateTaskpPresenter{}
}

func (t updateTaskpPresenter) Output(task domain.Task) usecase.UpdateTaskOutput {
	return usecase.UpdateTaskOutput{
		ID:        task.ID,
		Title:     task.Title,
		CreatedAt: task.CreatedAt.Format(time.DateTime),
		UpdatedAt: task.UpdatedAt.Format(time.DateTime),
	}
}
