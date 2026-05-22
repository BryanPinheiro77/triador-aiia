package service

import (
	"github.com/BryanPinheiro77/triador-aiia/internal/repository"
)

type AnalysisService struct {
	repository *repository.AnalysisRepository
}

func NewAnalysisService(
	repository *repository.AnalysisRepository,
) *AnalysisService {
	return &AnalysisService{
		repository: repository,
	}
}