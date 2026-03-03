package http

import (
	"net/http"

	"github.com/aswinsreeraj/evntx/internal/usecase"
	apiErrors "github.com/aswinsreeraj/evntx/pkg/errors"
	"github.com/aswinsreeraj/evntx/pkg/response"
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
		response.Error(c, http.StatusBadRequest, apiErrors.InvalidRequestBody, "Invalid email format")
		return
	}

	otp, err := h.authUsecase.RequestEmailOTP(req.Email)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, apiErrors.InvalidRequestBody, "Failed to generate OTP")
		return
	}

	response.Success(c, "OTP sent successfully", gin.H{
		"otp": otp, // remove later - only for testing
	})
}

type otpVerifyRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required,len=6"`
}

func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req otpVerifyRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apiErrors.InvalidRequestBody, "Invalid request body")
		return
	}

	access, refresh, err := h.authUsecase.VerifyEmailOTP(
		req.Email,
		req.OTP,
		c.Request.UserAgent(),
		c.ClientIP(),
	)

	if err != nil {
		response.Error(c, http.StatusUnauthorized, apiErrors.InvalidOTP, "Invalid or expired OTP")
		return
	}

	response.Success(c, "Login successful", gin.H{
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
		response.Error(c, http.StatusBadRequest, apiErrors.InvalidRequestBody, "Invalid request body")
		return
	}

	access, err := h.authUsecase.RefreshToken(req.RefreshToken)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, apiErrors.InvalidOTP, "Invalid refresh token")
		return
	}

	response.Success(c, "Token refreshed successfully", gin.H{
		"access_token": access,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req refreshRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apiErrors.InvalidRequestBody, "Invalid request body")
		return
	}

	if err := h.authUsecase.Logout(req.RefreshToken); err != nil {
		response.Error(c, http.StatusUnauthorized, apiErrors.InvalidOTP, "Failed to logout")
		return
	}

	response.Success(c, "Logged out successfully", nil)
}
