package payment

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aswinsreeraj/evntx/internal/domain"
	"github.com/aswinsreeraj/evntx/internal/repository"
	razorpay "github.com/razorpay/razorpay-go"
	"github.com/razorpay/razorpay-go/utils"
	"github.com/aswinsreeraj/evntx/pkg/logger"
)

type RazorpayService struct {
	keyID		string
	keySecret	string
	client		*razorpay.Client
}

func NewRazorpayService() repository.RazorpayService {
	keyID := os.Getenv("RAZORPAY_KEY_ID")
	keySecret := os.Getenv("RAZORPAY_KEY_SECRET")
	client := razorpay.NewClient(keyID, keySecret)

	return &RazorpayService{
		keyID:		keyID,
		keySecret:	keySecret,
		client:		client,
	}
}

func (s *RazorpayService) GetKeyID() string {
	return s.keyID
}

func (s *RazorpayService) CreateOrder(amount int64, receipt string) (*domain.RazorpayOrder, error) {
	if s.keyID == "" || s.keySecret == "" {
		logger.Log.Error().Msg("razorpay credentials are not configured in environment")
		return nil, fmt.Errorf("razorpay credentials are not configured")
	}

	orderData := map[string]interface{}{
		"amount":	amount,
		"currency":	"INR",
		"receipt":	receipt,
	}

	body, err := s.client.Order.Create(orderData, nil)
	if err != nil {
		logger.Log.Error().
			Err(err).
			Interface("order_data", orderData).
			Msg("razorpay order creation failed at SDK level")
		return nil, fmt.Errorf("razorpay order creation failed: %w", err)
	}

	respBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal razorpay response: %w", err)
	}

	var order domain.RazorpayOrder
	if err := json.Unmarshal(respBody, &order); err != nil {
		return nil, fmt.Errorf("failed to unmarshal razorpay response to domain: %w", err)
	}

	order.RawResponse = json.RawMessage(respBody)

	return &order, nil
}

func (s *RazorpayService) VerifySignature(orderID string, paymentID string, signature string) (bool, error) {
	if s.keySecret == "" {
		return false, fmt.Errorf("razorpay credentials are not configured")
	}

	params := map[string]interface{}{
		"razorpay_order_id":	orderID,
		"razorpay_payment_id":	paymentID,
	}

	return utils.VerifyPaymentSignature(params, signature, s.keySecret), nil
}

func (s *RazorpayService) FetchOrder(orderID string) (*domain.RazorpayOrder, error) {
	if s.keyID == "" || s.keySecret == "" {
		return nil, fmt.Errorf("razorpay credentials are not configured")
	}

	body, err := s.client.Order.Fetch(orderID, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch razorpay order: %w", err)
	}

	respBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal razorpay response: %w", err)
	}

	var order domain.RazorpayOrder
	if err := json.Unmarshal(respBody, &order); err != nil {
		return nil, fmt.Errorf("failed to unmarshal razorpay response to domain: %w", err)
	}

	return &order, nil
}

func (s *RazorpayService) RefundPayment(paymentID string, amount int64) error {
	if s.keyID == "" || s.keySecret == "" {
		return fmt.Errorf("razorpay credentials are not configured")
	}

	refundData := map[string]interface{}{
		"speed": "normal",
	}

	_, err := s.client.Payment.Refund(paymentID, int(amount), refundData, nil)
	if err != nil {
		logger.Log.Error().
			Err(err).
			Str("payment_id", paymentID).
			Int64("amount", amount).
			Msg("razorpay refund request failed")
		return fmt.Errorf("razorpay refund request failed: %w", err)
	}

	return nil
}
