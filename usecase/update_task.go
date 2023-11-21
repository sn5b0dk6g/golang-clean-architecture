package usecase

import (
	"go-rest-api/domain"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	UpdateTaskInput struct {
		UserID string `json:"user_id"`
		TaskID uint   `json:"task_id"`
		Title  string `json:"title"`
	}

	UpdateTaskOutput struct {
		ID        uint   `json:"id"`
		Title     string `json:"title"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	UpdateTaskUsecase interface {
		Execute(UpdateTaskInput) (UpdateTaskOutput, error)
	}

	UpdateTaskPresenter interface {
		Output(domain.Task) UpdateTaskOutput
	}

	UpdateTaskInteractor struct {
		taskRepo  domain.TaskRepository
		presenter UpdateTaskPresenter
	}
)

func NewUpdateTaskInteractor(
	taskRepo domain.TaskRepository,
	presenter UpdateTaskPresenter,
) UpdateTaskUsecase {
	return UpdateTaskInteractor{
		taskRepo:  taskRepo,
		presenter: presenter,
	}
}

func (t UpdateTaskInteractor) Execute(input UpdateTaskInput) (UpdateTaskOutput, error) {
	task := domain.Task{
		Title: input.Title,
	}
	err := t.taskRepo.Update(&task, *domain.UserIDWithArg(input.UserID), input.TaskID)
	if err != nil {
		return t.presenter.Output(domain.Task{}), err
	}
	return t.presenter.Output(task), nil
}

func (v UpdateTaskInput) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(
			&v.Title,
			validation.Required.Error("title is required"),
			validation.RuneLength(1, 30).Error("limited max 30 char"),
		),
	)
}
