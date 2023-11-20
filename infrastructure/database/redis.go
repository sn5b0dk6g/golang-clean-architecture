package database

import (
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func OpenRedis() (*redis.Client, error) {
	c := newConfigRedis()
	url := fmt.Sprintf("%s:%s", c.host, c.port)

	db, err := strconv.Atoi(c.database)
	if err != nil {
		return nil, err
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: c.password,
		DB:       db,
	})

	return rdb, nil
}
