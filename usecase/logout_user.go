package usecase

import (
	"go-rest-api/domain"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type (
	LogoutUserInput struct {
		Email string
	}

	LogoutUserUsecase interface {
		Execute(LogoutUserInput) error
	}

	LogoutUserInteractor struct {
		redis domain.UserRedisRepository
	}
)

func NewLogoutUserInteractor(redis domain.UserRedisRepository) LogoutUserUsecase {
	return LogoutUserInteractor{
		redis: redis,
	}
}

func (u LogoutUserInteractor) Execute(input LogoutUserInput) error {
	err := u.redis.RemoveToken(input.Email)
	if err != nil {
		return err
	}
	return nil
}

func (v LogoutUserInput) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(
			&v.Email,
			validation.Required.Error("email is required"),
			validation.RuneLength(1, 30).Error("limited max 30 char"),
			is.Email.Error("is not valid email format"),
		),
	)
}
