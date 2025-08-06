package repository

import (
	"github.com/google/uuid"
	"github.com/mqqff/savebite-be/internal/domain/dto"
	"github.com/mqqff/savebite-be/internal/domain/entity"
	"github.com/mqqff/savebite-be/pkg/log"
	"gorm.io/gorm"
)

type AnalysisRepoItf interface {
	Create(analysis *entity.Analysis) error
	GetByUserID(userID uuid.UUID, pagination dto.PaginationRequest) ([]entity.Analysis, int64, error)
}

type AnalysisRepo struct {
	db *gorm.DB
}

func NewAnalysisRepo(db *gorm.DB) AnalysisRepoItf {
	return &AnalysisRepo{db}
}

func (r *AnalysisRepo) Create(analysis *entity.Analysis) error {
	err := r.db.Create(analysis).Error
	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[AnalysisRepo][Create] failed to create analysis")
	}

	return err
}

func (r *AnalysisRepo) GetByUserID(userID uuid.UUID, pagination dto.PaginationRequest) ([]entity.Analysis, int64, error) {
	var totalItems int64
	var analyses []entity.Analysis
	err := r.db.Model(&entity.Analysis{}).Preload("Ingredients").Order("created_at desc").Offset(pagination.Offset).Limit(pagination.Limit).Find(&analyses, entity.Analysis{
		UserID: userID,
	}).Count(&totalItems).Error
	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[AnalysisRepo][GetByUserID] failed to retrieve analyses")
		return nil, 0, err
	}

	return analyses, totalItems, nil
}
