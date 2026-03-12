package http

import (
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/aswinsreeraj/evntx/internal/usecase"
	apiErrors "github.com/aswinsreeraj/evntx/pkg/errors"
	apiResponse "github.com/aswinsreeraj/evntx/pkg/response"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUsecase *usecase.UserUsecase
}

func NewUserHandler(u *usecase.UserUsecase) *UserHandler {
	return &UserHandler{userUsecase: u}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetString("user_id")

	user, err := h.userUsecase.GetProfile(userID)
	if err != nil {
		apiResponse.Error(c, http.StatusNotFound, apiErrors.ResourceNotFound, "User not found")
		return
	}

	apiResponse.Success(c, "Profile retrieved successfully", gin.H{
		"id":            user.ID,
		"name":          user.Name,
		"email":         user.Email,
		"mobile":        user.Mobile,
		"dob":           user.Dob,
		"gender":        user.Gender,
		"profile_image": user.ProfileImage,
		"locations":     user.Locations,
	})
}

type updateProfileRequest struct {
	Name      string   `json:"name" binding:"required"`
	Mobile    string   `json:"mobile"`
	Dob       string   `json:"dob"`
	Gender    string   `json:"gender"`
	Locations []string `json:"locations"`
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetString("user_id")

	var req updateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiResponse.Error(c, http.StatusBadRequest, apiErrors.InvalidRequestBody, "Invalid request body")
		return
	}

	nameRegex := regexp.MustCompile(`^[a-zA-Z\s]+$`)
	if !nameRegex.MatchString(req.Name) {
		apiResponse.Error(c, http.StatusBadRequest, apiErrors.InvalidRequestBody, "Name can only contain alphabets and spaces")
		return
	}

	if err := h.userUsecase.UpdateProfile(userID, req.Name, req.Mobile, req.Dob, req.Gender, req.Locations); err != nil {
		apiResponse.Error(c, http.StatusInternalServerError, apiErrors.InternalServerError, "Failed to update profile")
		return
	}

	apiResponse.Success(c, "Profile updated successfully", nil)
}

func (h *UserHandler) AdminListUsers(c *gin.Context) {

	search := c.Query("search")
	status := c.Query("status")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "5")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	users, total, err := h.userUsecase.AdminSearchUsers(search, status, page, limit)
	if err != nil {
		apiResponse.Error(c, http.StatusInternalServerError, apiErrors.InternalServerError, "Failed to retrieve users")
		return
	}

	response := make([]gin.H, 0)
	for _, u := range users {
		response = append(response, gin.H{
			"id":         u.ID,
			"name":       u.Name,
			"email":      u.Email,
			"is_active":  u.IsActive,
			"created_at": u.CreatedAt,
		})
	}

	apiResponse.Success(c, "Users retrieved successfully", gin.H{
		"users": response,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

type updateStatusRequest struct {
	IsActive bool `json:"is_active"`
}

func (h *UserHandler) AdminUpdateUserStatus(c *gin.Context) {

	userID := c.Param("id")

	var req updateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiResponse.Error(c, http.StatusBadRequest, apiErrors.InvalidRequestBody, "Invalid request body")
		return
	}

	if err := h.userUsecase.AdminUpdateUserStatus(userID, req.IsActive); err != nil {
		apiResponse.Error(c, http.StatusInternalServerError, apiErrors.InternalServerError, "Failed to update user status")
		return
	}

	apiResponse.Success(c, "User status updated successfully", nil)
}

func (h *UserHandler) UploadProfileImage(c *gin.Context) {
	userID := c.GetString("user_id")

	file, err := c.FormFile("profile_image")
	if err != nil {
		apiResponse.Error(c, http.StatusBadRequest, apiErrors.InvalidRequestBody, "Image file is required")
		return
	}

	dirPath := "assets/images/" + userID
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		apiResponse.Error(c, http.StatusInternalServerError, apiErrors.InternalServerError, "Failed to create directory")
		return
	}

	filepath := dirPath + "/" + file.Filename
	imageURL := "/" + filepath

	if err := c.SaveUploadedFile(file, filepath); err != nil {
		apiResponse.Error(c, http.StatusInternalServerError, apiErrors.InternalServerError, "Failed to save image")
		return
	}

	if err := h.userUsecase.UploadProfileImage(userID, imageURL); err != nil {
		apiResponse.Error(c, http.StatusInternalServerError, apiErrors.InternalServerError, "Failed to update profile image")
		return
	}

	apiResponse.Success(c, "Profile image uploaded successfully", gin.H{
		"profile_image": imageURL,
	})
}
