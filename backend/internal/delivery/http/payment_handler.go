package http

import (
	"net/http"

	"github.com/aswinsreeraj/evntx/internal/usecase"
	apiErrors "github.com/aswinsreeraj/evntx/pkg/errors"
	"github.com/aswinsreeraj/evntx/pkg/response"
	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentUsecase *usecase.PaymentUsecase
}

func NewPaymentHandler(paymentUsecase *usecase.PaymentUsecase) *PaymentHandler {
	return &PaymentHandler{
		paymentUsecase: paymentUsecase,
	}
}

func (h *PaymentHandler) CreateRazorpayOrder(c *gin.Context) {
	var req struct {
		BookingID string `json:"booking_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.AppError(c, apiErrors.ErrInvalidRequestBody)
		return
	}

	userID := c.GetString("user_id")
	paymentOrder, err := h.paymentUsecase.CreatePaymentOrder(c.Request.Context(), req.BookingID, userID)
	if err != nil {
		response.AppError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response.APIResponse{
		Success: true,
		Message: "Payment order created successfully",
		Data:    paymentOrder,
	})
}
