package presenter

import "go-rest-api/usecase"

type authenticationUserPresenter struct{}

func NewAuthenticationUserPresenter() usecase.AuthenticationUserPresenter {
	return authenticationUserPresenter{}
}

func (p authenticationUserPresenter) Output(jwtTokenString string) usecase.AuthenticationUserOutput {
	return usecase.AuthenticationUserOutput{
		JwtTokenString: jwtTokenString,
	}
}
