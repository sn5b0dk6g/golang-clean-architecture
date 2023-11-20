package database

import "gorm.io/gorm"

func newGormConfig() *gorm.Config {
	return &gorm.Config{
		TranslateError: true,
	}
}
