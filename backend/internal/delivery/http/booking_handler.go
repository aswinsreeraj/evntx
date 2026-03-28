package http

import (
	"net/http"
	"strings"

	"github.com/aswinsreeraj/evntx/internal/domain"
	"github.com/aswinsreeraj/evntx/internal/usecase"
	"github.com/aswinsreeraj/evntx/pkg/response"
	"github.com/gin-gonic/gin"
)

type BookingHandler struct {
	bookingUsecase *usecase.BookingUsecase
	paymentUsecase *usecase.PaymentUsecase
}

func NewBookingHandler(bookingUsecase *usecase.BookingUsecase, paymentUsecase *usecase.PaymentUsecase) *BookingHandler {
	return &BookingHandler{
		bookingUsecase: bookingUsecase,
		paymentUsecase: paymentUsecase,
	}
}

func (h *BookingHandler) ReserveTickets(c *gin.Context) {
	userID := c.GetString("user_id")

	var req struct {
		EventID string                 `json:"event_id" binding:"required"`
		Tickets []domain.TicketRequest `json:"tickets" binding:"required,dive"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	booking, err := h.bookingUsecase.ReserveTickets(c.Request.Context(), userID, req.EventID, req.Tickets)
	if err != nil {
		errMsg := err.Error()

		if strings.Contains(errMsg, "EVT_009") {
			response.Error(c, http.StatusConflict, "EVT_009", "Tickets sold out")
			return
		} else if strings.Contains(errMsg, "EVT_012") {
			response.Error(c, http.StatusBadRequest, "EVT_012", "Event not live")
			return
		}

		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to reserve tickets")
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Booking reserved",
		"data": gin.H{
			"booking_id":   booking.ID,
			"expires_at":   booking.ExpiresAt,
			"total_amount": booking.TotalAmount,
		},
	})
}

func (h *BookingHandler) CancelBooking(c *gin.Context) {
	userID := c.GetString("user_id")
	bookingID := c.Param("booking_id")

	var req struct {
		Items []domain.TicketCancelRequest `json:"items" binding:"required,dive"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid cancellation request")
		return
	}

	err := h.bookingUsecase.CancelBooking(c.Request.Context(), bookingID, userID, req.Items)
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "not found") {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "Booking not found")
			return
		}
		if strings.Contains(errMsg, "cannot be cancelled") {
			response.Error(c, http.StatusBadRequest, "INVALID_STATE", "Booking cannot be cancelled")
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to cancel booking")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Booking cancelled successfully",
	})
}

func (h *BookingHandler) RefundBooking(c *gin.Context) {
	userID := c.GetString("user_id")
	bookingID := c.Param("booking_id")

	if err := h.paymentUsecase.RefundPaymentToWallet(c.Request.Context(), bookingID, userID); err != nil {
		response.AppError(c, err)
		return
	}

	response.Success(c, "Refund processed successfully", nil)
}
