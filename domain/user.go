package domain

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type (
	UserRepository interface {
		FindByEmail(EmailAddress) (*User, error)
		Create(*User) error
	}

	UserRedisRepository interface {
		SaveToken(string, string) error
		RemoveToken(string) error
	}

	User struct {
		gorm.Model
		ID       UserID       `json:"id" gorm:"embedded;not null;unique"`
		Email    EmailAddress `json:"email" gorm:"embedded;not null;unique;index"`
		Password Password     `json:"password" gorm:"embedded;not null"`
	}
)
