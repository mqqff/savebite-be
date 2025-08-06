package repository

import (
	"github.com/mqqff/savebite-be/internal/domain/dto"
	"github.com/mqqff/savebite-be/internal/domain/entity"
	"github.com/mqqff/savebite-be/pkg/log"
	"gorm.io/gorm"
)

type UserRepoItf interface {
	Show(param dto.UserParam) (entity.User, error)
	Create(user *entity.User) error
}

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepoItf {
	return &UserRepo{db}
}

func (r *UserRepo) Create(user *entity.User) error {
	err := r.db.Create(user).Error
	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[UserRepo][Create] failed to create user")
	}

	return err
}

func (r *UserRepo) Show(param dto.UserParam) (entity.User, error) {
	user := entity.User{}
	err := r.db.Find(&user, param).Error
	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[UserRepo][Create] failed to get user")
		return entity.User{}, err
	}

	return user, err
}
