package delivery

import (
	"net/http"

	"github.com/brkss/dextrace/internal/domain"
	"github.com/brkss/dextrace/internal/usecase"
	"github.com/gin-gonic/gin"
)

type GlucoseHandler struct {
	glucoseUseCase *usecase.SibionicUseCase
	userID        string
	user          domain.User
}

func NewGlucoseHandler(glucoseUseCase *usecase.SibionicUseCase, userID string, user domain.User) *GlucoseHandler {
	return &GlucoseHandler{
		glucoseUseCase: glucoseUseCase,
		userID:        userID,
		user:          user,
	}
}

func (h *GlucoseHandler) GetGlucoseData(c *gin.Context) {
	userID := h.userID
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	data, err := h.glucoseUseCase.GetGlucoseData(h.user, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}