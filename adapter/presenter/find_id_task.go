package presenter

import (
	"go-rest-api/domain"
	"go-rest-api/usecase"
	"time"
)

type findIdTaskPresenter struct{}

func NewFindIdTaskPresenter() usecase.FindIdTaskPresenter {
	return findIdTaskPresenter{}
}

func (t findIdTaskPresenter) Output(task domain.Task) usecase.FindIdTaskOutput {
	return usecase.FindIdTaskOutput{
		ID:        task.ID,
		Title:     task.Title,
		CreatedAt: task.CreatedAt.Format(time.DateTime),
		UpdatedAt: task.UpdatedAt.Format(time.DateTime),
	}
}
