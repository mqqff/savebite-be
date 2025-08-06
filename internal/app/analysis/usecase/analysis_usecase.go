package usecase

import (
	"bytes"
	"github.com/google/uuid"
	"github.com/mqqff/savebite-be/internal/app/analysis/repository"
	"github.com/mqqff/savebite-be/internal/domain/dto"
	"github.com/mqqff/savebite-be/internal/domain/entity"
	"github.com/mqqff/savebite-be/internal/infra/gemini"
	"github.com/mqqff/savebite-be/pkg/log"
	"github.com/mqqff/savebite-be/pkg/markdown"
	"github.com/mqqff/savebite-be/pkg/supabase"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
)

type AnalysisUsecaseItf interface {
	Analyze(file *multipart.FileHeader, userID uuid.UUID) (dto.AnalysisResponse, error)
	GetHistory(userID uuid.UUID, pagination dto.PaginationRequest) ([]dto.AnalysisResponse, dto.PaginationResponse, error)
}

type AnalysisUsecase struct {
	analysisRepo repository.AnalysisRepoItf
	supabase     supabase.SupabaseItf
	gemini       gemini.GeminiItf
	md           markdown.MarkdownItf
}

func NewAnalysisUsecase(r repository.AnalysisRepoItf, s supabase.SupabaseItf, g gemini.GeminiItf, m markdown.MarkdownItf) AnalysisUsecaseItf {
	return &AnalysisUsecase{
		analysisRepo: r,
		supabase:     s,
		gemini:       g,
		md:           m,
	}
}

func (u *AnalysisUsecase) Analyze(file *multipart.FileHeader, userID uuid.UUID) (dto.AnalysisResponse, error) {
	image, err := file.Open()
	defer image.Close()
	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[AnalysisUsecase][Analyze] failed to open image file")
		return dto.AnalysisResponse{}, err
	}

	imageData, err := io.ReadAll(image)
	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[AnalysisUsecase][Analyze] failed to read image data")
		return dto.AnalysisResponse{}, err
	}

	imageBody := bytes.NewReader(imageData)

	fileName := file.Filename
	ext := filepath.Ext(fileName)

	mimeType := file.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "image/jpeg"

		switch ext {
		case ".png":
			mimeType = "image/png"
		case ".jpg":
			mimeType = "image/jpeg"
		case ".gif":
			mimeType = "image/gif"
		case ".webp":
			mimeType = "image/webp"
		}
	}

	result, err := u.gemini.AnalyzeImage(mimeType, imageData)
	if err != nil {
		return dto.AnalysisResponse{}, err
	}

	result.Feedback = u.md.MdToHTML(result.Feedback)
	result.Feedback = strings.ReplaceAll(result.Feedback, "\n", "")
	result.Feedback = strings.ReplaceAll(result.Feedback, "\t", "")

	path := uuid.NewString() + ext
	bucket := "analysis-images"
	publicURL, err := u.supabase.UploadFile(bucket, path, mimeType, imageBody)
	if err != nil {
		return dto.AnalysisResponse{}, err
	}

	analysisID := uuid.New()

	var ingredients []entity.Ingredient

	for _, ingredient := range result.UsableIngredients {
		ingredients = append(ingredients, entity.Ingredient{
			AnalysisID: analysisID,
			Name:       ingredient,
			Status:     "usable",
		})
	}

	for _, ingredient := range result.UnusableIngredients {
		ingredients = append(ingredients, entity.Ingredient{
			AnalysisID: analysisID,
			Name:       ingredient,
			Status:     "unusable",
		})
	}

	analysis := &entity.Analysis{
		ID:          analysisID,
		UserID:      userID,
		ImageURL:    publicURL,
		Feedback:    result.Feedback,
		Ingredients: ingredients,
	}

	err = u.analysisRepo.Create(analysis)
	if err != nil {
		return dto.AnalysisResponse{}, nil
	}

	response := dto.AnalysisResponse{
		DetectedItems:       result.DetectedItems,
		UsableIngredients:   result.UsableIngredients,
		UnusableIngredients: result.UnusableIngredients,
		Feedback:            result.Feedback,
	}

	return response, nil
}

func (u *AnalysisUsecase) GetHistory(userID uuid.UUID, pagination dto.PaginationRequest) ([]dto.AnalysisResponse, dto.PaginationResponse, error) {
	if pagination.Limit < 1 {
		pagination.Limit = 10
	}

	if pagination.Page < 1 {
		pagination.Page = 1
	}

	pagination.Offset = (pagination.Page - 1) * pagination.Limit

	analysesDB, totalItems, err := u.analysisRepo.GetByUserID(userID, pagination)
	if err != nil {
		return nil, dto.PaginationResponse{}, err
	}

	var resp []dto.AnalysisResponse

	for _, analysis := range analysesDB {
		var detectedItems []string
		var usableIngredients []string
		var unusableIngredients []string

		for _, ingredient := range analysis.Ingredients {
			detectedItems = append(detectedItems, ingredient.Name)
			if ingredient.Status == "usable" {
				usableIngredients = append(usableIngredients, ingredient.Name)
			} else if ingredient.Status == "unusable" {
				unusableIngredients = append(unusableIngredients, ingredient.Name)
			}
		}

		resp = append(resp, dto.AnalysisResponse{
			DetectedItems:       detectedItems,
			UsableIngredients:   usableIngredients,
			UnusableIngredients: unusableIngredients,
			Feedback:            analysis.Feedback,
			CreatedAt:           analysis.CreatedAt.String(),
		})
	}

	totalPages := float64(totalItems) / float64(pagination.Limit)
	if totalPages > float64(int(totalPages)) {
		totalPages++
	}

	meta := dto.PaginationResponse{
		Page:       pagination.Page,
		Limit:      pagination.Limit,
		TotalItems: totalItems,
		TotalPages: int64(totalPages),
	}

	return resp, meta, nil
}
