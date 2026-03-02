package http

import (
	"net/http"

	"github.com/aswinsreeraj/evntx/internal/usecase"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUsecase *usecase.AuthUsecase
}

func NewAuthHandler(u *usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: u}
}

type otpRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func (h *AuthHandler) RequestOTP(c *gin.Context) {
	var req otpRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request body",
		})
		return
	}

	otp, err := h.authUsecase.RequestEmailOTP(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to generate OTP",
		})
		return
	}

	// temporary: return OTP for testing
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"otp":     otp,
	})
}
