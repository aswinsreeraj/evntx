package http

import (
	"net/http"

	"github.com/aswinsreeraj/evntx/pkg/logger"

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

	isNewUser, err := h.authUsecase.RequestEmailOTP(req.Email)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, apiErrors.InvalidRequestBody, "Failed to generate OTP")
		return
	}

	response.Success(c, "OTP sent successfully", gin.H{
		"is_new_user": isNewUser,
	})
}

type otpVerifyRequest struct {
	Email string  `json:"email" binding:"required,email"`
	OTP   string  `json:"otp" binding:"required,len=6"`
	Name  *string `json:"name"`
}

func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req otpVerifyRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apiErrors.InvalidRequestBody, "Invalid request body")
		return
	}

	name := ""
	if req.Name != nil {
		name = *req.Name
	}

	user, roles, access, refresh, err := h.authUsecase.VerifyEmailOTP(
		req.Email,
		req.OTP,
		name,
		c.Request.UserAgent(),
		c.ClientIP(),
	)

	if err != nil {
		response.Error(c, http.StatusUnauthorized, apiErrors.InvalidOTP, "Invalid or expired OTP")
		return
	}

	roleStrings := make([]string, 0, len(roles))
	for _, r := range roles {
		roleStrings = append(roleStrings, string(r))
	}

	response.Success(c, "Login successful", gin.H{
		"access_token":  access,
		"refresh_token": refresh,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"roles": roleStrings,
		},
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

type googleLoginRequest struct {
	IDToken string `json:"id_token" binding:"required"`
}

func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	var req googleLoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apiErrors.InvalidRequestBody, "Invalid request body")
		return
	}

	access, refresh, err := h.authUsecase.GoogleLogin(
		req.IDToken,
		c.Request.UserAgent(),
		c.ClientIP(),
	)

	if err != nil {
		logger.Log.Error().Msgf("AuthUsecase.GoogleLogin failed: %v", err)
		response.Error(c, http.StatusUnauthorized, apiErrors.UnauthorizedAccess, "Invalid Google token")
		return
	}

	response.Success(c, "Login successful", gin.H{
		"access_token":  access,
		"refresh_token": refresh,
	})
}

type registerRequest struct {
	Email  string `json:"email" binding:"required,email"`
	OTP    string `json:"otp" binding:"required,len=6"`
	Name   string `json:"name" binding:"required"`
	Dob    string `json:"dob"`
	Gender string `json:"gender"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apiErrors.InvalidRequestBody, "Invalid request body")
		return
	}

	user, roles, access, refresh, err := h.authUsecase.Register(
		req.Email,
		req.OTP,
		req.Name,
		req.Dob,
		req.Gender,
		c.Request.UserAgent(),
		c.ClientIP(),
	)

	if err != nil {
		response.Error(c, http.StatusUnauthorized, apiErrors.InvalidOTP, "Invalid OTP or registration failed")
		return
	}

	roleStrings := make([]string, 0, len(roles))
	for _, r := range roles {
		roleStrings = append(roleStrings, string(r))
	}

	response.Success(c, "Registration successful", gin.H{
		"access_token":  access,
		"refresh_token": refresh,
		"user": gin.H{
			"id":     user.ID,
			"name":   user.Name,
			"dob":    user.Dob,
			"gender": user.Gender,
			"roles":  roleStrings,
		},
	})
}
