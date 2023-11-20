package repository

import (
	"errors"
	"go-rest-api/domain"

	"gorm.io/gorm"
)

type UserSQL struct {
	db *gorm.DB
}

func NewUserSQL(db *gorm.DB) *UserSQL {
	return &UserSQL{db}
}

func findFirst(db *gorm.DB, query interface{}, args ...interface{}) (*domain.User, error) {
	var user domain.User
	if err := db.Where(query, args).First(&user).Error; err != nil {
		return &user, err
	}

	return &user, nil
}

func (u UserSQL) FindByID(id domain.UserID) (*domain.User, error) {
	return findFirst(u.db, "id=?", id)
}

func (u UserSQL) FindByEmail(email domain.EmailAddress) (*domain.User, error) {
	return findFirst(u.db, "email=?", email)
}

func (u UserSQL) Create(user *domain.User) error {
	if err := u.db.Create(user).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errors.New("一意制約違反です。")
		}
		return err
	}

	return nil
}
