package dto

type AnalysisRequest struct {
	Resume         string `json:"resume" binding:"required"`
	JobDescription string `json:"job_description" binding:"required"`
}
