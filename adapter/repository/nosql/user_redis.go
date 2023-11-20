package repository

import (
	"context"
	"encoding/json"
	"go-rest-api/domain"
	"time"

	"github.com/redis/go-redis/v9"
)

type UserRedis struct {
	db        *redis.Client
	exp       time.Duration
	keyPrefix string
}

func NewUserRedis(db *redis.Client, exp time.Duration) *UserRedis {
	return &UserRedis{
		db:        db,
		exp:       exp,
		keyPrefix: "user:",
	}
}

func findBy(db *redis.Client, key string) (*domain.User, error) {
	result, err := db.Get(context.Background(), key).Result()
	if err != nil {
		return &domain.User{}, err
	}

	var user domain.User
	err = json.Unmarshal([]byte(result), &user)
	if err != nil {
		return &domain.User{}, err
	}
	return &user, nil
}

func (u UserRedis) FindByEmail(email domain.EmailAddress) (*domain.User, error) {
	return findBy(u.db, u.keyPrefix+email.String())
}

func (u UserRedis) Create(user *domain.User) error {
	bson, err := json.Marshal(user)
	if err != nil {
		return err
	}

	err = u.db.Set(context.Background(), u.keyPrefix+user.Email.String(), bson, u.exp).Err()
	if err != nil {
		return err
	}

	return nil
}

func (u UserRedis) SaveToken(token string, email string) error {
	err := u.db.Set(context.Background(), email, token, u.exp).Err()
	if err != nil {
		return err
	}
	return nil
}

func (u UserRedis) RemoveToken(email string) error {
	err := u.db.Del(context.Background(), email).Err()
	if err != nil {
		return err
	}
	return nil
}
