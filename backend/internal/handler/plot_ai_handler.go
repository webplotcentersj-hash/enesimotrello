package handler

import (
	"net/http"
	"strconv"
	"task-board/internal/service"

	"github.com/gin-gonic/gin"
)

type PlotAIHandler struct {
	plotAIService service.PlotAIService
}

func NewPlotAIHandler(plotAIService service.PlotAIService) *PlotAIHandler {
	return &PlotAIHandler{
		plotAIService: plotAIService,
	}
}

type ChatRequest struct {
	Message string `json:"message" binding:"required"`
}

func (h *PlotAIHandler) SendMessage(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.plotAIService.SendMessage(userID, req.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *PlotAIHandler) GetHistory(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50
	}

	history, err := h.plotAIService.GetHistory(userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"history": history})
}

