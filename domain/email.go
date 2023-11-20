package domain

import (
	"database/sql/driver"
	"errors"
)

var (
	ErrInvalidEmailFormat = errors.New("invalid email format")
)

type EmailAddress struct {
	email string `gorm:"type:varchar(255)"`
}

func NewEmailAddress(email string) *EmailAddress {
	return &EmailAddress{email: email}
}

func (e EmailAddress) String() string {
	return e.email
}

func (e EmailAddress) Equals(other EmailAddress) bool {
	return e.email == other.email
}

func (e EmailAddress) Value() (driver.Value, error) {
	return e.String(), nil
}

func (e *EmailAddress) Scan(src interface{}) error {
	switch value := src.(type) {
	case string:
		e.email = value
	case []byte:
		e.email = string(value)
	default:
		return errors.New("incompatible type for EmailAddress")
	}

	return nil
}
