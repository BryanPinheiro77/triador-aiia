package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/BryanPinheiro77/triador-aiia/internal/dto"
	"github.com/BryanPinheiro77/triador-aiia/internal/service"
)

type AnalysisHandler struct {
	service *service.AnalysisService
}

func NewAnalysisHandler(
	service *service.AnalysisService,
) *AnalysisHandler {
	return &AnalysisHandler{
		service: service,
	}
}

func (h *AnalysisHandler) Create(c *gin.Context) {
	var request dto.AnalysisRequest

	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	response, err := h.service.Create(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *AnalysisHandler) FindAll(c *gin.Context) {
	responses, err := h.service.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, responses)
}