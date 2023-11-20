package usecase

import "go-rest-api/domain"

type (
	DeleteTaskInput struct {
		UserID string `json:"user_id"`
		TaskID uint   `json:"task_id"`
	}

	DeleteTaskUsecase interface {
		Execute(DeleteTaskInput) error
	}

	DeleteTaskInteractor struct {
		taskRepo domain.TaskRepository
	}
)

func NewDeleteTaskInteractor(
	taskRepo domain.TaskRepository,
) DeleteTaskUsecase {
	return DeleteTaskInteractor{
		taskRepo: taskRepo,
	}
}

func (t DeleteTaskInteractor) Execute(input DeleteTaskInput) error {
	if err := t.taskRepo.Delete(*domain.UserIDWithArg(input.UserID), input.TaskID); err != nil {
		return err
	}
	return nil
}
