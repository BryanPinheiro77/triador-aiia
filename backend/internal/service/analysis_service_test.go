package service

import (
	"strings"
	"testing"

	"github.com/BryanPinheiro77/triador-aiia/internal/dto"
)

func TestValidateLLMResponse_ShouldReturnNil_WhenResponseIsValid(t *testing.T) {
	service := &AnalysisService{}

	response := dto.LLMAnalysisResponse{
		CandidateName:   "Bryan Mendes",
		Skills:          []string{"Go", "Next.js", "PostgreSQL"},
		YearsExperience: 1,
		FitScore:        85,
		Summary:         "Boa aderência para a vaga.",
	}

	err := service.validateLLMResponse(response)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidateLLMResponse_ShouldReturnError_WhenFitScoreIsInvalid(t *testing.T) {
	service := &AnalysisService{}

	response := dto.LLMAnalysisResponse{
		CandidateName:   "Bryan Mendes",
		Skills:          []string{"Go"},
		YearsExperience: 1,
		FitScore:        150,
		Summary:         "Score inválido.",
	}

	err := service.validateLLMResponse(response)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "fit_score") {
		t.Fatalf("expected fit_score error, got %v", err)
	}
}

func TestValidateLLMResponse_ShouldReturnError_WhenSkillsAreMissing(t *testing.T) {
	service := &AnalysisService{}

	response := dto.LLMAnalysisResponse{
		CandidateName:   "Bryan Mendes",
		Skills:          []string{},
		YearsExperience: 1,
		FitScore:        80,
		Summary:         "Resumo válido.",
	}

	err := service.validateLLMResponse(response)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "skills") {
		t.Fatalf("expected skills error, got %v", err)
	}
}

func TestSanitizeLLMResponse_ShouldRemoveMarkdownCodeBlock(t *testing.T) {
	service := &AnalysisService{}

	rawResponse := "```json\n{\"candidate_name\":\"Bryan\",\"skills\":[\"Go\"],\"years_experience\":1,\"fit_score\":80,\"summary\":\"Resumo.\"}\n```"

	cleaned := service.sanitizeLLMResponse(rawResponse)

	if strings.Contains(cleaned, "```") {
		t.Fatalf("expected markdown code block to be removed, got %s", cleaned)
	}

	if !strings.Contains(cleaned, "\"candidate_name\":\"Bryan\"") {
		t.Fatalf("expected JSON content to be preserved, got %s", cleaned)
	}
}

func TestSanitizeLLMResponse_ShouldKeepPlainJSON(t *testing.T) {
	service := &AnalysisService{}

	rawResponse := "{\"candidate_name\":\"Bryan\",\"skills\":[\"Go\"],\"years_experience\":1,\"fit_score\":80,\"summary\":\"Resumo.\"}"

	cleaned := service.sanitizeLLMResponse(rawResponse)

	if cleaned != rawResponse {
		t.Fatalf("expected plain JSON to remain unchanged, got %s", cleaned)
	}
}
