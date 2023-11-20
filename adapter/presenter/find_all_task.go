package presenter

import (
	"go-rest-api/domain"
	"go-rest-api/usecase"
	"time"
)

type findAllTaskPresenter struct{}

func NewFindAllTaskPresenter() usecase.FindAllTaskPresenter {
	return findAllTaskPresenter{}
}

func (t findAllTaskPresenter) Output(tasks []domain.Task) []usecase.FindAllTaskOutput {
	var o = make([]usecase.FindAllTaskOutput, 0)
	for _, task := range tasks {
		o = append(o, usecase.FindAllTaskOutput{
			ID:        task.ID,
			Title:     task.Title,
			CreatedAt: task.CreatedAt.Format(time.DateTime),
			UpdatedAt: task.UpdatedAt.Format(time.DateTime),
		})
	}
	return o
}
