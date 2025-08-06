package usecase

import (
	"errors"
	"github.com/google/uuid"
	"github.com/mqqff/savebite-be/internal/app/user/repository"
	"github.com/mqqff/savebite-be/internal/domain/dto"
	"github.com/mqqff/savebite-be/internal/domain/entity"
	"github.com/mqqff/savebite-be/internal/domain/env"
	"github.com/mqqff/savebite-be/pkg/jwt"
	"github.com/mqqff/savebite-be/pkg/oauth"
)

type AuthUsecaseItf interface {
	HandleRedirect(state string) (string, error)
	HandleCallback(data dto.GoogleCallbackRequest) (string, error)
}

type AuthUsecase struct {
	userRepo repository.UserRepoItf
	oauth    oauth.OAuthItf
	jwt      jwt.JWTIf
}

func NewAuthUsecase(r repository.UserRepoItf, o oauth.OAuthItf, j jwt.JWTIf) AuthUsecaseItf {
	return &AuthUsecase{
		userRepo: r,
		oauth:    o,
		jwt:      j,
	}
}

func (u *AuthUsecase) HandleRedirect(state string) (string, error) {
	return u.oauth.GenerateAuthLink(state), nil
}

func (u *AuthUsecase) HandleCallback(data dto.GoogleCallbackRequest) (string, error) {
	if data.State != env.AppEnv.OAuthState {
		return "", errors.New("invalid state")
	}

	token, err := u.oauth.ExchangeToken(data.Code)
	if err != nil {
		return "", err
	}

	userInfo, err := u.oauth.GetUserInfo(token)
	if err != nil {
		return "", err
	}

	user, err := u.userRepo.Show(dto.UserParam{Email: userInfo.Email})
	if err != nil {
		return "", err
	}

	if user.ID == uuid.Nil {
		user = entity.User{
			ID:    uuid.New(),
			Email: userInfo.Email,
			Name:  userInfo.Name,
		}

		err := u.userRepo.Create(&user)
		if err != nil {
			return "", err
		}
	}

	jwtToken, err := u.jwt.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}
