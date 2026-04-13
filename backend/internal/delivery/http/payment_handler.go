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

func (h *PaymentHandler) VerifyRazorpayPayment(c *gin.Context) {
	var req struct {
		RazorpayOrderID   string `json:"razorpay_order_id" binding:"required"`
		RazorpayPaymentID string `json:"razorpay_payment_id" binding:"required"`
		RazorpaySignature string `json:"razorpay_signature" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.AppError(c, apiErrors.ErrInvalidRequestBody)
		return
	}

	if err := h.paymentUsecase.VerifyPayment(c.Request.Context(), req.RazorpayOrderID, req.RazorpayPaymentID, req.RazorpaySignature); err != nil {
		if err == apiErrors.ErrBookingExpiredPaymentSuccess {
			c.JSON(http.StatusOK, response.APIResponse{
				Success: true,
				Message: "Payment captured after booking expiry; refund initiated to source",
				Data: map[string]interface{}{
					"is_late_payment":  true,
					"is_source_refund": true,
				},
			})
			return
		}
		response.AppError(c, err)
		return
	}

	response.Success(c, "Payment verified successfully", map[string]interface{}{
		"is_late_payment": false,
	})
}
