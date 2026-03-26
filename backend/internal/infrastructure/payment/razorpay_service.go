package payment

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/aswinsreeraj/evntx/internal/domain"
	"github.com/aswinsreeraj/evntx/internal/repository"
)

const razorpayCreateOrderURL = "https://api.razorpay.com/v1/orders"

type RazorpayService struct {
	keyID      string
	keySecret  string
	httpClient *http.Client
}

func NewRazorpayService() repository.RazorpayService {
	return &RazorpayService{
		keyID:     os.Getenv("RAZORPAY_KEY_ID"),
		keySecret: os.Getenv("RAZORPAY_KEY_SECRET"),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *RazorpayService) GetKeyID() string {
	return s.keyID
}

func (s *RazorpayService) CreateOrder(amount int64, receipt string) (*domain.RazorpayOrder, error) {
	if s.keyID == "" || s.keySecret == "" {
		return nil, fmt.Errorf("razorpay credentials are not configured")
	}

	payload := struct {
		Amount   int64  `json:"amount"`
		Currency string `json:"currency"`
		Receipt  string `json:"receipt"`
	}{
		Amount:   amount,
		Currency: "INR",
		Receipt:  receipt,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, razorpayCreateOrderURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(s.keyID, s.keySecret)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("razorpay order creation failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	var order domain.RazorpayOrder
	if err := json.Unmarshal(respBody, &order); err != nil {
		return nil, err
	}

	order.RawResponse = json.RawMessage(respBody)

	return &order, nil
}
