package usecase

import "go-rest-api/domain"

type (
	FindAllTaskInput struct {
		UserID string `json:"user_id"`
	}

	FindAllTaskOutput struct {
		ID        uint   `json:"id"`
		Title     string `json:"title"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	FindAllTaskUsecase interface {
		Execute(FindAllTaskInput) ([]FindAllTaskOutput, error)
	}

	FindAllTaskPresenter interface {
		Output([]domain.Task) []FindAllTaskOutput
	}

	FindAllTaskInteractor struct {
		taskRepo  domain.TaskRepository
		presenter FindAllTaskPresenter
	}
)

func NewFindAllTaskInteractor(
	taskRepo domain.TaskRepository,
	presenter FindAllTaskPresenter,
) FindAllTaskUsecase {
	return FindAllTaskInteractor{
		taskRepo:  taskRepo,
		presenter: presenter,
	}
}

func (t FindAllTaskInteractor) Execute(input FindAllTaskInput) ([]FindAllTaskOutput, error) {
	userId := domain.UserIDWithArg(input.UserID)
	tasks, err := t.taskRepo.FindAll(*userId)
	if err != nil {
		return t.presenter.Output([]domain.Task{}), err
	}
	return t.presenter.Output(*tasks), nil
}
