package usecase

import "go-rest-api/domain"

type (
	CreateTaskInput struct {
		Title  string `json:"title"`
		UserID string `json:"user_id"`
	}

	CreateTaskOutput struct {
		ID        uint   `json:"id"`
		Title     string `json:"title"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	CreateTaskUsecase interface {
		Execute(CreateTaskInput) (CreateTaskOutput, error)
	}

	CreateTaskPresenter interface {
		Output(domain.Task) CreateTaskOutput
	}

	CreateTaskInteractor struct {
		taskRepo  domain.TaskRepository
		presenter CreateTaskPresenter
	}
)

func NewCreateTaskInteractor(
	taskRepo domain.TaskRepository,
	presenter CreateTaskPresenter,
) CreateTaskUsecase {
	return CreateTaskInteractor{
		taskRepo:  taskRepo,
		presenter: presenter,
	}
}

func (t CreateTaskInteractor) Execute(input CreateTaskInput) (CreateTaskOutput, error) {
	task := domain.Task{
		Title:  input.Title,
		UserId: *domain.UserIDWithArg(input.UserID),
	}
	if err := t.taskRepo.Create(&task); err != nil {
		return t.presenter.Output(domain.Task{}), err
	}
	return t.presenter.Output(task), nil
}
