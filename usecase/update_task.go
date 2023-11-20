package usecase

import "go-rest-api/domain"

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
