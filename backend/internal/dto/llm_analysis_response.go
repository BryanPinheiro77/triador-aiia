package dto

type LLMAnalysisResponse struct {
	CandidateName   string   `json:"candidate_name"`
	Skills          []string `json:"skills"`
	YearsExperience int      `json:"years_experience"`
	FitScore        int      `json:"fit_score"`
	Summary         string   `json:"summary"`
}
