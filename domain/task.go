package domain

import (
	"gorm.io/gorm"
)

type (
	TaskRepository interface {
		FindAll(userId UserID) (*[]Task, error)
		FindByID(userId UserID, taskId uint) (*Task, error)
		Create(task *Task) error
		Update(task *Task, userId UserID, taskId uint) error
		Delete(userId UserID, taskId uint) error
	}

	Task struct {
		gorm.Model
		Title  string `json:"title" gorm:"not null"`
		User   User   `json:"user" gorm:"foreignKey:UserId; constraint:OnDelete:CASCADE"`
		UserId UserID `json:"user_id" gorm:"not null;embedded;index"`
	}
)
