package http

import (
	"log"
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

type otpVerifyRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required,len=6"`
}

func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req otpVerifyRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false})
		return
	}

	access, refresh, err := h.authUsecase.VerifyEmailOTP(
		req.Email,
		req.OTP,
		c.Request.UserAgent(),
		c.ClientIP(),
	)

	if err != nil {
		log.Println("VerifyEmailOTP error:", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":       true,
		"access_token":  access,
		"refresh_token": refresh,
	})
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req refreshRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false})
		return
	}

	access, err := h.authUsecase.RefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"access_token": access,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req refreshRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false})
		return
	}

	if err := h.authUsecase.Logout(req.RefreshToken); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
