package repository

import (
	"github.com/BryanPinheiro77/triador-aiia/internal/database"
	"github.com/BryanPinheiro77/triador-aiia/internal/model"
)

type AnalysisRepository struct {
}

func NewAnalysisRepository() *AnalysisRepository {
	return &AnalysisRepository{}
}

func (r *AnalysisRepository) Save(analysis *model.Analysis) error {
	return database.DB.Create(analysis).Error
}

func (r *AnalysisRepository) FindAll() ([]model.Analysis, error) {
	var analyses []model.Analysis

	err := database.DB.Order("created_at desc").Find(&analyses).Error

	return analyses, err
}