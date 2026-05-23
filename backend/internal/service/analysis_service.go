package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/BryanPinheiro77/triador-aiia/internal/dto"
	"github.com/BryanPinheiro77/triador-aiia/internal/llm"
	"github.com/BryanPinheiro77/triador-aiia/internal/model"
	"github.com/BryanPinheiro77/triador-aiia/internal/repository"
)

type AnalysisService struct {
	repository *repository.AnalysisRepository
	llmClient  *llm.OpenAIClient
}

func NewAnalysisService(
	repository *repository.AnalysisRepository,
	llmClient *llm.OpenAIClient,
) *AnalysisService {
	return &AnalysisService{
		repository: repository,
		llmClient:  llmClient,
	}
}

func (s *AnalysisService) Create(
	ctx context.Context,
	request dto.AnalysisRequest,
) (*dto.AnalysisResponse, error) {
	prompt := s.buildPrompt(request)

	rawResponse, err := s.llmClient.Analyze(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrLLMProviderFailed, err)
	}

	var llmResponse dto.LLMAnalysisResponse

	cleanResponse := s.sanitizeLLMResponse(rawResponse)

	err = json.Unmarshal([]byte(cleanResponse), &llmResponse)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidLLMOutput, err)
	}

	err = s.validateLLMResponse(llmResponse)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidLLMOutput, err)
	}

	skillsJSON, err := json.Marshal(llmResponse.Skills)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidLLMOutput, err)
	}

	analysis := model.Analysis{
		CandidateName:   llmResponse.CandidateName,
		Skills:          string(skillsJSON),
		YearsExperience: llmResponse.YearsExperience,
		FitScore:        llmResponse.FitScore,
		Summary:         llmResponse.Summary,
	}

	err = s.repository.Save(&analysis)
	if err != nil {
		return nil, err
	}

	response := dto.AnalysisResponse{
		ID:              analysis.ID,
		CandidateName:   analysis.CandidateName,
		Skills:          llmResponse.Skills,
		YearsExperience: analysis.YearsExperience,
		FitScore:        analysis.FitScore,
		Summary:         analysis.Summary,
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

func (s *AnalysisService) buildPrompt(request dto.AnalysisRequest) string {
	return fmt.Sprintf(`
Analyze the resume according to the job description.

Return only valid JSON, without markdown, without explanations and without code block.

Expected JSON format:
{
  "candidate_name": "string",
  "skills": ["string"],
  "years_experience": 0,
  "fit_score": 0,
  "summary": "string"
}

Rules:
- fit_score must be between 0 and 100.
- years_experience must be an approximate integer based only on explicit or strongly implied professional experience from the resume. If experience is unclear, return 0.
- summary must be short, justify the score and always be written in Portuguese (Brazil).
- skills must include only technical skills found in the resume.
- Do not infer experience time only from education, academic projects or general knowledge.
- If the candidate name is not clear, use "Unknown candidate".

Resume:
%s

Job description:
%s
`, request.Resume, request.JobDescription)
}

func (s *AnalysisService) validateLLMResponse(response dto.LLMAnalysisResponse) error {
	if response.CandidateName == "" {
		return errors.New("LLM response missing candidate_name")
	}

	if len(response.Skills) == 0 {
		return errors.New("LLM response missing skills")
	}

	if response.YearsExperience < 0 {
		return errors.New("LLM response has invalid years_experience")
	}

	if response.FitScore < 0 || response.FitScore > 100 {
		return errors.New("LLM response has invalid fit_score")
	}

	if response.Summary == "" {
		return errors.New("LLM response missing summary")
	}

	return nil
}

func (s *AnalysisService) sanitizeLLMResponse(rawResponse string) string {
	cleaned := strings.TrimSpace(rawResponse)

	cleaned = strings.TrimPrefix(cleaned, "```json")
	cleaned = strings.TrimPrefix(cleaned, "```")
	cleaned = strings.TrimSuffix(cleaned, "```")

	return strings.TrimSpace(cleaned)
}
