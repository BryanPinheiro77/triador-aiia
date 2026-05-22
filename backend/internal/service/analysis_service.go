package service

import (
	"encoding/json"

	"github.com/BryanPinheiro77/triador-aiia/internal/dto"
	"github.com/BryanPinheiro77/triador-aiia/internal/model"
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

func (s *AnalysisService) Create(
	request dto.AnalysisRequest,
) (*dto.AnalysisResponse, error) {

	mockSkills := []string{
		"Go",
		"Next.js",
		"PostgreSQL",
	}

	skillsJSON, err := json.Marshal(mockSkills)
	if err != nil {
		return nil, err
	}

	analysis := model.Analysis{
		CandidateName: "Mock Candidate",
		Skills: string(skillsJSON),
		YearsExperience: 2,
		FitScore: 85,
		Summary: "Strong backend profile with good alignment.",
	}

	err = s.repository.Save(&analysis)
	if err != nil {
		return nil, err
	}

	response := dto.AnalysisResponse{
		ID: analysis.ID,
		CandidateName: analysis.CandidateName,
		Skills: mockSkills,
		YearsExperience: analysis.YearsExperience,
		FitScore: analysis.FitScore,
		Summary: analysis.Summary,
	}

	return &response, nil
}

func (s *AnalysisService) FindAll() ([]dto.AnalysisResponse, error) {
	analyses, err := s.repository.FindAll()
	if err != nil {
		return nil, err
	}

	responses := make([]dto.AnalysisResponse, 0, len(analyses))

	for _, analysis := range analyses {
		var skills []string

		err := json.Unmarshal([]byte(analysis.Skills), &skills)
		if err != nil {
			skills = []string{}
		}

		responses = append(responses, dto.AnalysisResponse{
			ID:              analysis.ID,
			CandidateName:   analysis.CandidateName,
			Skills:          skills,
			YearsExperience: analysis.YearsExperience,
			FitScore:        analysis.FitScore,
			Summary:         analysis.Summary,
		})
	}

	return responses, nil
}