package usecase

import (
	"go-rest-api/domain"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type (
	CreateUserUsecase interface {
		Execute(CreateUserInput) (CreateUserOutput, error)
	}

	CreateUserInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	CreateUserPresenter interface {
		Output(domain.User) CreateUserOutput
	}

	CreateUserOutput struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	}

	CreateUserInteractor struct {
		userRepo  domain.UserRepository
		presenter CreateUserPresenter
	}
)

func NewCreateUserInteractor(
	userRepo domain.UserRepository,
	presenter CreateUserPresenter,
) CreateUserUsecase {
	return CreateUserInteractor{
		userRepo:  userRepo,
		presenter: presenter,
	}
}

func (u CreateUserInteractor) Execute(input CreateUserInput) (CreateUserOutput, error) {
	password, err := domain.NewPassword(input.Password)
	if err != nil {
		return u.presenter.Output(domain.User{}), err
	}

	strings.Split(input.Email, " ")

	user := domain.User{
		ID:       *domain.NewUserID(),
		Email:    *domain.NewEmailAddress(input.Email),
		Password: *password,
	}

	if err := u.userRepo.Create(&user); err != nil {
		return u.presenter.Output(domain.User{}), err
	}

	return u.presenter.Output(user), nil
}

func (v CreateUserInput) Validate() error {
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
