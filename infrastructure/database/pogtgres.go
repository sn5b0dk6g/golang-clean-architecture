package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func OpenPostgres() (*gorm.DB, error) {
	c := newConfigPostgres()
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", c.user, c.password, c.host, c.port, c.database)
	db, err := gorm.Open(postgres.Open(url), newGormConfig())
	if err != nil {
		return nil, err
	}
	return db, nil
}

func ClosePostgres(db *gorm.DB) {
	sqlDB, _ := db.DB()
	if err := sqlDB.Close(); err != nil {
		log.Fatalln(err)
	}
}
