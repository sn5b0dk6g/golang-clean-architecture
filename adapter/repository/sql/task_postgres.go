package repository

import (
	"errors"
	"go-rest-api/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrTaskNotExist = errors.New("task does not exist")
)

type TaskSQL struct {
	db *gorm.DB
}

func NewTaskSQL(db *gorm.DB) *TaskSQL {
	return &TaskSQL{db}
}

func (t TaskSQL) FindAll(userId domain.UserID) (*[]domain.Task, error) {
	tasks := []domain.Task{}
	err := t.db.Joins("User").Where("user_id=?", userId.String()).Order("created_at").Find(tasks).Error
	if err != nil {
		return &tasks, err
	}
	return &tasks, nil
}

func (t TaskSQL) FindByID(userId domain.UserID, taskId uint) (*domain.Task, error) {
	var task domain.Task
	err := t.db.Joins("User").Where("user_id=?", userId.String()).First(&task, taskId).Error
	if err != nil {
		return &task, err
	}
	return &task, nil
}

func (t TaskSQL) Create(task *domain.Task) error {
	if err := t.db.Create(task).Error; err != nil {
		return err
	}
	return nil
}

func (t TaskSQL) Update(task *domain.Task, userId domain.UserID, taskId uint) error {
	result := t.db.Model(task).Clauses(clause.Returning{}).Where("id=? AND user_id=?", taskId, userId).Update("title", task.Title)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected < 1 {
		return ErrTaskNotExist
	}
	return nil
}

func (t TaskSQL) Delete(userId domain.UserID, taskId uint) error {
	result := t.db.Where("id=? AND user_id=?", taskId, userId).Delete(&domain.Task{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected < 1 {
		return ErrTaskNotExist
	}
	return nil
}
