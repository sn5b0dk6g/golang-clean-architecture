package usecase

import (
	"go-rest-api/domain"
	"os"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/golang-jwt/jwt/v4"
)

type (
	AuthenticationUserUsecase interface {
		Execute(AuthenticationUserInput) (AuthenticationUserOutput, error)
	}

	AuthenticationUserValidator interface {
		Validate(AuthenticationUserInput) error
	}

	AuthenticationUserInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	AuthenticationUserPresenter interface {
		Output(string) AuthenticationUserOutput
	}

	AuthenticationUserOutput struct {
		JwtTokenString string `json:"jwtTokenString"`
	}

	AuthenticationUserInteractor struct {
		userRepo  domain.UserRepository
		redis     domain.UserRedisRepository
		presenter AuthenticationUserPresenter
	}
)

func NewAuthenticationUserInteractor(
	userRepo domain.UserRepository,
	redis domain.UserRedisRepository,
	presenter AuthenticationUserPresenter,
) AuthenticationUserUsecase {
	return AuthenticationUserInteractor{
		userRepo:  userRepo,
		redis:     redis,
		presenter: presenter,
	}
}

func (u AuthenticationUserInteractor) Execute(input AuthenticationUserInput) (AuthenticationUserOutput, error) {
	user, err := u.userRepo.FindByEmail(*domain.NewEmailAddress(input.Email))
	if err != nil {
		return AuthenticationUserOutput{}, err
	}

	if err = user.Password.CompareHashAndPassword([]byte(input.Password)); err != nil {
		return AuthenticationUserOutput{}, err
	}

	// jwtトークンの生成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.String(),
		"exp":     time.Now().Add(12 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return AuthenticationUserOutput{}, err
	}

	if err = u.redis.SaveToken(tokenString, input.Email); err != nil {
		return AuthenticationUserOutput{}, err
	}

	return u.presenter.Output(tokenString), nil
}

func (v AuthenticationUserInput) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(
			&v.Email,
			validation.Required.Error("email is required"),
			validation.RuneLength(1, 30).Error("limited max 30 char"),
			is.Email.Error("is not valid email format"),
		),
		validation.Field(
			&v.Password,
			validation.Required.Error("password is required"),
			validation.RuneLength(6, 30).Error("limited min 6 max 30 char"),
		),
	)
}
