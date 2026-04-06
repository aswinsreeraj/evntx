package http

import (
	"strconv"
	"strings"

	"github.com/aswinsreeraj/evntx/internal/domain"
	"github.com/aswinsreeraj/evntx/internal/repository"
	"github.com/aswinsreeraj/evntx/internal/usecase"
	pkgErrors "github.com/aswinsreeraj/evntx/pkg/errors"
	"github.com/aswinsreeraj/evntx/pkg/response"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	eventUsecase        *usecase.EventUsecase
	userUsecase         *usecase.UserUsecase
	walletUsecase       *usecase.WalletUsecase
	platformWalletRepo  repository.PlatformWalletRepository
}

func NewAdminHandler(eventUsecase *usecase.EventUsecase, userUsecase *usecase.UserUsecase, walletUsecase *usecase.WalletUsecase, platformWalletRepo repository.PlatformWalletRepository) *AdminHandler {
	return &AdminHandler{
		eventUsecase:       eventUsecase,
		userUsecase:        userUsecase,
		walletUsecase:      walletUsecase,
		platformWalletRepo: platformWalletRepo,
	}
}

func (h *AdminHandler) AdminListEvents(c *gin.Context) {
	search := c.Query("search")
	status := c.Query("status")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page := 1
	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}

	limit := 10
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}

	events, total, err := h.eventUsecase.AdminSearchEvents(search, status, page, limit)
	if err != nil {
		response.AppError(c, pkgErrors.ErrInternalServerError)
		return
	}

	eventResponses := make([]gin.H, 0, len(events))
	for _, evt := range events {
		eventResponses = append(eventResponses, gin.H{
			"id":             evt.ID,
			"slug":           evt.Slug,
			"title":          evt.Title,
			"organizer_name": evt.OrganizerName,
			"start_time":     evt.StartTime,
			"date":           evt.StartTime,
			"tickets_sold":   evt.TicketsSold,
			"revenue":        evt.Revenue,
			"city":           evt.City,
			"status":         evt.Status,
		})
	}

	response.Success(c, "Events retrieved successfully", gin.H{
		"events": eventResponses,
		"pagination": gin.H{
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

func (h *AdminHandler) ApproveEventHandler(c *gin.Context) {
	adminID := c.GetString("user_id")
	eventID := c.Param("event_id")

	err := h.eventUsecase.ApproveEvent(c.Request.Context(), adminID, eventID)
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "EVT_004") {
			response.AppError(c, pkgErrors.New(403, pkgErrors.ForbiddenAction, "Insufficient admin privileges"))
			return
		} else if strings.Contains(errMsg, "EVT_006") {
			response.AppError(c, pkgErrors.New(409, pkgErrors.InvalidStateTransition, "Event cannot be approved in current state"))
			return
		}
		response.AppError(c, pkgErrors.ErrInternalServerError)
		return
	}

	response.Success(c, "Event approved successfully", gin.H{
		"event_id": eventID,
		"status":   "approved",
	})
}

func (h *AdminHandler) RejectEventHandler(c *gin.Context) {
	adminID := c.GetString("user_id")
	eventID := c.Param("event_id")

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.AppError(c, pkgErrors.ErrInvalidRequestBody)
		return
	}

	err := h.eventUsecase.RejectEvent(c.Request.Context(), adminID, eventID, req.Reason)
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "EVT_004") {
			response.AppError(c, pkgErrors.New(403, pkgErrors.ForbiddenAction, "Insufficient admin privileges"))
			return
		} else if strings.Contains(errMsg, "EVT_006") {
			response.AppError(c, pkgErrors.New(409, pkgErrors.InvalidStateTransition, "Event cannot be rejected in current state"))
			return
		}
		response.AppError(c, pkgErrors.ErrInternalServerError)
		return
	}

	response.Success(c, "Event rejected successfully", gin.H{
		"event_id": eventID,
		"status":   "rejected",
	})
}

func (h *AdminHandler) SuspendEventHandler(c *gin.Context) {
	adminID := c.GetString("user_id")
	eventID := c.Param("event_id")

	var req struct {
		Reason string `json:"reason"`
	}

	if err := c.ShouldBindJSON(&req); err != nil || req.Reason == "" {
		response.AppError(c, pkgErrors.New(400, pkgErrors.InvalidRequestBody, "Reason is required"))
		return
	}

	err := h.eventUsecase.SuspendLiveEvent(c.Request.Context(), adminID, eventID, req.Reason)
	if err != nil {
		response.AppError(c, err)
		return
	}

	response.Success(c, "Event suspended successfully", gin.H{
		"event_id": eventID,
		"status":   "suspended",
	})
}

func (h *AdminHandler) CompleteEventHandler(c *gin.Context) {
	adminID := c.GetString("user_id")
	eventID := c.Param("event_id")

	if err := h.eventUsecase.CompleteEvent(c.Request.Context(), adminID, eventID); err != nil {
		response.AppError(c, err)
		return
	}

	response.Success(c, "Event completed successfully", gin.H{
		"event_id": eventID,
		"status":   "completed",
	})
}

func (h *AdminHandler) SettleEventHandler(c *gin.Context) {
	eventID := c.Param("event_id")

	if err := h.eventUsecase.SettleEventEarnings(c.Request.Context(), eventID); err != nil {
		response.AppError(c, err)
		return
	}

	response.Success(c, "Settlement completed", gin.H{
		"event_id": eventID,
	})
}

func (h *AdminHandler) AdminGetEvent(c *gin.Context) {
	slug := c.Param("slug")

	event, details, personnels, tickets, err := h.eventUsecase.AdminGetEvent(slug)

	if err != nil {
		response.AppError(c, pkgErrors.ErrResourceNotFound)
		return
	}

	host := gin.H(nil)
	if evt, ok := event.(*domain.Event); ok && h.userUsecase != nil {
		user, organizerDetail, _, userErr := h.userUsecase.GetProfile(evt.OrganizerID)
		if userErr == nil && user != nil {
			hostName := user.Name
			if organizerDetail != nil && organizerDetail.OrganizationName != "" {
				hostName = organizerDetail.OrganizationName
			}
			host = gin.H{
				"name":   hostName,
				"role":   "Event Organizer",
				"avatar": user.ProfileImage,
			}
		}
	}

	response.Success(c, "Event fetched successfully", gin.H{
		"event":        event,
		"details":      details,
		"personnels":   personnels,
		"ticket_types": tickets,
		"host":         host,
	})
}
func (h *AdminHandler) GetPlatformWallet(c *gin.Context) {
	wallet, err := h.platformWalletRepo.GetPlatformWallet()
	if err != nil {
		response.AppError(c, pkgErrors.ErrInternalServerError)
		return
	}

	response.Success(c, "Platform wallet retrieved successfully", gin.H{
		"available_balance": wallet.AvailableBalance,
		"pending_balance":   wallet.PendingBalance,
		"total_credited":    wallet.TotalCredited,
		"total_debited":     wallet.TotalDebited,
		"updated_at":        wallet.UpdatedAt,
	})
}

func (h *AdminHandler) AdminGetPayouts(c *gin.Context) {
	status := c.Query("status")

	payouts, total, err := h.walletUsecase.AdminGetPayoutRequests(c.Request.Context(), status, 1, 50)
	if err != nil {
		response.AppError(c, pkgErrors.ErrInternalServerError)
		return
	}

	response.Success(c, "Payouts retrieved successfully", gin.H{
		"payouts": payouts,
		"total":   total,
	})
}

func (h *AdminHandler) AdminApprovePayout(c *gin.Context) {
	adminID := c.GetString("user_id")
	payoutID := c.Param("id")

	if err := h.walletUsecase.AdminApprovePayout(c.Request.Context(), adminID, payoutID); err != nil {
		response.AppError(c, err)
		return
	}

	response.Success(c, "Payout approved successfully", nil)
}

func (h *AdminHandler) AdminRejectPayout(c *gin.Context) {
	adminID := c.GetString("user_id")
	payoutID := c.Param("id")

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.AppError(c, pkgErrors.ErrInvalidRequestBody)
		return
	}

	if err := h.walletUsecase.AdminRejectPayout(c.Request.Context(), adminID, payoutID, req.Reason); err != nil {
		response.AppError(c, err)
		return
	}

	response.Success(c, "Payout rejected successfully", nil)
}

func (h *AdminHandler) AdminBulkApprovePayouts(c *gin.Context) {
	adminID := c.GetString("user_id")

	var req struct {
		PayoutIDs []string `json:"payout_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.AppError(c, pkgErrors.ErrInvalidRequestBody)
		return
	}

	if err := h.walletUsecase.AdminBulkApprovePayouts(c.Request.Context(), adminID, req.PayoutIDs); err != nil {
		response.AppError(c, err)
		return
	}

	response.Success(c, "Payouts approved successfully", nil)
}
