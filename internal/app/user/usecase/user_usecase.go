package usecase

import (
	"github.com/google/uuid"
	"github.com/mqqff/savebite-be/internal/app/user/repository"
	"github.com/mqqff/savebite-be/internal/domain/dto"
)

type UserUsecaseItf interface {
	GetProfile(userID uuid.UUID) (dto.UserProfile, error)
}

type UserUsecase struct {
	userRepo repository.UserRepoItf
}

func NewUserUsecase(r repository.UserRepoItf) UserUsecaseItf {
	return &UserUsecase{
		userRepo: r,
	}
}

func (u *UserUsecase) GetProfile(userID uuid.UUID) (dto.UserProfile, error) {
	user, err := u.userRepo.Show(dto.UserParam{ID: userID})
	if err != nil {
		return dto.UserProfile{}, err
	}

	userProfile := dto.UserProfile{
		ID:    user.ID.String(),
		Email: user.Email,
		Name:  user.Name,
	}

	return userProfile, err
}
