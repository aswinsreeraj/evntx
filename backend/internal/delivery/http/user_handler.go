package http

import (
	"os"
	"regexp"
	"strconv"

	"github.com/aswinsreeraj/evntx/internal/usecase"
	apiErrors "github.com/aswinsreeraj/evntx/pkg/errors"
	apiResponse "github.com/aswinsreeraj/evntx/pkg/response"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUsecase    *usecase.UserUsecase
	bookingUsecase *usecase.BookingUsecase
}

func NewUserHandler(userUsecase *usecase.UserUsecase, bookingUsecase *usecase.BookingUsecase) *UserHandler {
	return &UserHandler{
		userUsecase:    userUsecase,
		bookingUsecase: bookingUsecase,
	}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetString("user_id")

	user, orgDetail, roles, err := h.userUsecase.GetProfile(userID)
	if err != nil {
		apiResponse.AppError(c, apiErrors.ErrResourceNotFound)
		return
	}

	resp := gin.H{
		"id":            user.ID,
		"name":          user.Name,
		"email":         user.Email,
		"mobile":        user.Mobile,
		"dob":           user.Dob,
		"gender":        user.Gender,
		"profile_image": user.ProfileImage,
		"locations":     user.Locations,
		"roles":         roles,
	}

	if orgDetail != nil {
		resp["organization_name"] = orgDetail.OrganizationName
		resp["address"] = orgDetail.Address
	}

	apiResponse.Success(c, "Profile retrieved successfully", resp)
}

type updateProfileRequest struct {
	Name             string   `json:"name" binding:"required"`
	Mobile           string   `json:"mobile"`
	Dob              string   `json:"dob"`
	Gender           string   `json:"gender"`
	OrganizationName string   `json:"organization_name"`
	Address          string   `json:"address"`
	Locations        []string `json:"locations"`
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetString("user_id")

	var req updateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiResponse.AppError(c, apiErrors.ErrInvalidRequestBody)
		return
	}

	nameRegex := regexp.MustCompile(`^[a-zA-Z\s]+$`)
	if !nameRegex.MatchString(req.Name) {
		apiResponse.AppError(c, apiErrors.New(400, apiErrors.InvalidRequestBody, "Name can only contain alphabets and spaces"))
		return
	}

	if err := h.userUsecase.UpdateProfile(userID, req.Name, req.Mobile, req.Dob, req.Gender, req.OrganizationName, req.Address, req.Locations); err != nil {
		apiResponse.AppError(c, apiErrors.ErrInternalServerError)
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
		apiResponse.AppError(c, apiErrors.ErrInternalServerError)
		return
	}

	resp := make([]gin.H, 0, len(users))
	for _, u := range users {
		resp = append(resp, gin.H{
			"id":         u.ID,
			"name":       u.Name,
			"email":      u.Email,
			"is_active":  u.IsActive,
			"created_at": u.CreatedAt,
		})
	}

	apiResponse.Success(c, "Users retrieved successfully", gin.H{
		"users": resp,
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
		apiResponse.AppError(c, apiErrors.ErrInvalidRequestBody)
		return
	}

	if err := h.userUsecase.AdminUpdateUserStatus(userID, req.IsActive); err != nil {
		apiResponse.AppError(c, apiErrors.ErrInternalServerError)
		return
	}

	apiResponse.Success(c, "User status updated successfully", nil)
}

func (h *UserHandler) AdminListOrganizers(c *gin.Context) {
	search := c.Query("search")
	status := c.Query("status")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "5")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	orgs, total, err := h.userUsecase.AdminSearchOrganizers(search, status, page, limit)
	if err != nil {
		apiResponse.AppError(c, apiErrors.ErrInternalServerError)
		return
	}

	resp := make([]gin.H, 0, len(orgs))
	for _, u := range orgs {
		resp = append(resp, gin.H{
			"id":                      u.ID,
			"name":                    u.Name,
			"email":                   u.Email,
			"is_active":               u.IsActive,
			"total_bookings":          u.TotalBookings,
			"total_events":            u.TotalEvents,
			"wallet_balance":          u.WalletBalance,
			"total_revenue_generated": u.TotalRevenue,
		})
	}

	apiResponse.Success(c, "Organizers retrieved successfully", gin.H{
		"organizers": resp,
		"pagination": gin.H{
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

func (h *UserHandler) GetMyBookingsHandler(c *gin.Context) {
	userID := c.GetString("user_id")

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	status := c.Query("status")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	bookings, total, err := h.bookingUsecase.GetUserBookings(c.Request.Context(), userID, page, limit, status)
	if err != nil {
		apiResponse.AppError(c, apiErrors.ErrInternalServerError)
		return
	}

	responseBookings := make([]map[string]interface{}, 0, len(bookings))
	for _, b := range bookings {
		responseBookings = append(responseBookings, map[string]interface{}{
			"booking_id":       b.BookingID,
			"event_id":         b.EventID,
			"event_title":      b.EventTitle,
			"event_city":       b.EventCity,
			"event_start_time": b.EventStartTime,
			"status":           b.Status,
			"total_amount":     b.TotalAmount,
			"ticket_count":     b.TicketCount,
			"created_at":       b.CreatedAt,
		})
	}

	apiResponse.Success(c, "Bookings fetched successfully", gin.H{
		"bookings": responseBookings,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func (h *UserHandler) GetMyTicketsHandler(c *gin.Context) {
	userID := c.GetString("user_id")

	eventID := c.Query("event_id")
	bookingID := c.Query("booking_id")
	status := c.Query("status")

	tickets, err := h.bookingUsecase.GetUserTickets(c.Request.Context(), userID, eventID, bookingID, status)
	if err != nil {
		apiResponse.AppError(c, apiErrors.ErrInternalServerError)
		return
	}

	responseTickets := make([]map[string]interface{}, 0, len(tickets))
	for _, t := range tickets {
		responseTickets = append(responseTickets, map[string]interface{}{
			"ticket_id":     t.TicketID,
			"ticket_code":   t.TicketCode,
			"event_id":      t.EventID,
			"event_title":   t.EventTitle,
			"ticket_type":   t.TicketType,
			"status":        t.Status,
			"checked_in_at": t.CheckedInAt,
		})
	}

	apiResponse.Success(c, "Tickets fetched successfully", gin.H{
		"tickets": responseTickets,
	})
}

func (h *UserHandler) UploadProfileImage(c *gin.Context) {
	userID := c.GetString("user_id")

	file, err := c.FormFile("profile_image")
	if err != nil {
		apiResponse.AppError(c, apiErrors.New(400, apiErrors.InvalidRequestBody, "Image file is required"))
		return
	}

	dirPath := "assets/images/" + userID
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		apiResponse.AppError(c, apiErrors.ErrInternalServerError)
		return
	}

	filepath := dirPath + "/" + file.Filename
	imageURL := "/" + filepath

	if err := c.SaveUploadedFile(file, filepath); err != nil {
		apiResponse.AppError(c, apiErrors.New(500, apiErrors.InternalServerError, "Failed to save image"))
		return
	}

	if err := h.userUsecase.UploadProfileImage(userID, imageURL); err != nil {
		apiResponse.AppError(c, apiErrors.New(500, apiErrors.InternalServerError, "Failed to update profile image"))
		return
	}

	apiResponse.Success(c, "Profile image uploaded successfully", gin.H{
		"profile_image": imageURL,
	})
}
