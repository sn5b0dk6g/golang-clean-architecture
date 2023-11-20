package presenter

import (
	"go-rest-api/domain"
	"go-rest-api/usecase"
)

type createUserPresenter struct{}

func NewCreateUserPresenter() usecase.CreateUserPresenter {
	return createUserPresenter{}
}

func (p createUserPresenter) Output(user domain.User) usecase.CreateUserOutput {
	return usecase.CreateUserOutput{
		ID:    user.ID.String(),
		Email: user.Email.String(),
	}
}
