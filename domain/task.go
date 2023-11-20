package domain

import (
	"gorm.io/gorm"
)

type (
	TaskRepository interface {
		FindAll(UserID) (*[]Task, error)
		FindByID(UserID, taskId uint) (*Task, error)
		Create(task *Task) error
		Update(task *Task, UserID, taskId uint) error
		Delete(UserID, taskId uint) error
	}

	Task struct {
		gorm.Model
		Title  string `json:"title" gorm:"not null"`
		User   User   `json:"user" gorm:"foreignKey:UserId; constraint:OnDelete:CASCADE"`
		UserId UserID `json:"user_id" gorm:"not null;embedded;index"`
	}
)
