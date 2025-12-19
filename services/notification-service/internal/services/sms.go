package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"notification-service/internal/config"
	"notification-service/internal/models"
)

type SMSService struct {
	config config.TwilioConfig
	client *http.Client
}

func NewSMSService(cfg config.TwilioConfig) *SMSService {
	return &SMSService{
		config: cfg,
		client: &http.Client{},
	}
}

func (s *SMSService) Send(notification *models.Notification) error {
	// Build Twilio API URL
	apiURL := fmt.Sprintf(
		"https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json",
		s.config.AccountSID,
	)

	// Prepare form data
	data := url.Values{}
	data.Set("To", notification.Recipient)
	data.Set("From", s.config.FromNumber)
	data.Set("Body", notification.Content)

	// Create request
	req, err := http.NewRequest("POST", apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers and auth
	req.SetBasicAuth(s.config.AccountSID, s.config.AuthToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send request
	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send SMS: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode >= 400 {
		var errResp struct {
			Message string `json:"message"`
			Code    int    `json:"code"`
		}
		json.NewDecoder(resp.Body).Decode(&errResp)
		return fmt.Errorf("twilio error: %s (code: %d)", errResp.Message, errResp.Code)
	}

	return nil
}

// SendBulk sends SMS to multiple recipients
func (s *SMSService) SendBulk(notifications []*models.Notification) []error {
	errors := make([]error, len(notifications))
	for i, n := range notifications {
		errors[i] = s.Send(n)
	}
	return errors
}

// ValidatePhoneNumber checks if a phone number is valid
func (s *SMSService) ValidatePhoneNumber(phone string) bool {
	// Basic E.164 format validation
	if len(phone) < 10 || len(phone) > 15 {
		return false
	}
	if phone[0] != '+' {
		return false
	}
	return true
}
