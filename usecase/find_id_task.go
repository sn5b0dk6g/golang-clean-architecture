package usecase

import "go-rest-api/domain"

type (
	FindIdTaskInput struct {
		UserID string `json:"user_id"`
		TaskID uint   `json:"task_id"`
	}

	FindIdTaskOutput struct {
		ID        uint   `json:"id"`
		Title     string `json:"title"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	FindIdTaskUsecase interface {
		Execute(FindIdTaskInput) (FindIdTaskOutput, error)
	}

	FindIdTaskPresenter interface {
		Output(domain.Task) FindIdTaskOutput
	}

	FindIdTaskInteractor struct {
		taskRepo  domain.TaskRepository
		presenter FindIdTaskPresenter
	}
)

func NewFindIdTaskInteractor(
	taskRepo domain.TaskRepository,
	presenter FindIdTaskPresenter,
) FindIdTaskUsecase {
	return FindIdTaskInteractor{
		taskRepo:  taskRepo,
		presenter: presenter,
	}
}

func (t FindIdTaskInteractor) Execute(input FindIdTaskInput) (FindIdTaskOutput, error) {
	task, err := t.taskRepo.FindByID(*domain.UserIDWithArg(input.UserID), input.TaskID)
	if err != nil {
		return t.presenter.Output(domain.Task{}), err
	}
	return t.presenter.Output(*task), nil
}
