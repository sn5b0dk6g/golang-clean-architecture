package domain

import (
	"database/sql/driver"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Password struct {
	password string `gorm:"type:varchar(255)"`
}

func NewPassword(password string) (*Password, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return &Password{}, err
	}

	return &Password{password: string(hash)}, nil
}

func (p Password) String() string {
	return p.password
}

func (p Password) CompareHashAndPassword(password []byte) error {
	return bcrypt.CompareHashAndPassword([]byte(p.password), password)
}

func (p Password) Equals(other Password) bool {
	return p.password == other.password
}

func (p Password) Value() (driver.Value, error) {
	return p.String(), nil
}

func (p *Password) Scan(src interface{}) error {
	switch value := src.(type) {
	case string:
		p.password = value
	case []byte:
		p.password = string(value)
	default:
		return errors.New("incompatible type for Password")
	}

	return nil
}
