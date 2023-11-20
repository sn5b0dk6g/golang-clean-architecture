package domain

import (
	"database/sql/driver"
	"errors"
)

type UserID struct {
	userID string `gorm:"type:varchar(255)"`
}

func NewUserID() *UserID {
	return &UserID{userID: NewUUID()}
}

func UserIDWithArg(id string) *UserID {
	return &UserID{id}
}

func (u UserID) String() string {
	return u.userID
}

func (u UserID) Equals(other UserID) bool {
	return u.userID == other.userID
}

func (u UserID) Value() (driver.Value, error) {
	return u.String(), nil
}

func (u *UserID) Scan(src interface{}) error {
	switch value := src.(type) {
	case string:
		u.userID = value
	case []byte:
		u.userID = string(value)
	default:
		return errors.New("incompatible type for UserID")
	}

	return nil
}
