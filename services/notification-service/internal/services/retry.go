package services

import (
	"log"
	"math"
	"time"

	"notification-service/internal/models"
	"notification-service/internal/repository"
)

const (
	MaxRetries     = 3
	BaseRetryDelay = 5 * time.Second
)

type RetryService struct {
	repo         *repository.NotificationRepository
	emailService *EmailService
	smsService   *SMSService
	pushService  *PushService
}

func NewRetryService(
	repo *repository.NotificationRepository,
	email *EmailService,
	sms *SMSService,
	push *PushService,
) *RetryService {
	return &RetryService{
		repo:         repo,
		emailService: email,
		smsService:   sms,
		pushService:  push,
	}
}

// RetryFailedNotifications attempts to resend failed notifications
func (s *RetryService) RetryFailedNotifications() {
	log.Println("Starting retry job for failed notifications...")

	// In a real implementation, you'd query failed notifications from DB
	// This is a simplified example
}

// RetryWithBackoff attempts to send a notification with exponential backoff
func (s *RetryService) RetryWithBackoff(notification *models.Notification) error {
	var lastErr error

	for attempt := 0; attempt < MaxRetries; attempt++ {
		// Calculate delay with exponential backoff
		delay := time.Duration(math.Pow(2, float64(attempt))) * BaseRetryDelay

		if attempt > 0 {
			log.Printf("Retry attempt %d for notification %s, waiting %v",
				attempt+1, notification.ID, delay)
			time.Sleep(delay)
		}

		// Attempt to send
		var err error
		switch notification.Type {
		case models.NotificationTypeEmail:
			err = s.emailService.Send(notification)
		case models.NotificationTypeSMS:
			err = s.smsService.Send(notification)
		case models.NotificationTypePush:
			err = s.pushService.Send(notification)
		}

		if err == nil {
			// Success!
			s.repo.UpdateStatus(notification.ID, models.StatusSent, "")
			log.Printf("Successfully sent notification %s on attempt %d",
				notification.ID, attempt+1)
			return nil
		}

		lastErr = err
		s.repo.IncrementRetryCount(notification.ID)
		log.Printf("Attempt %d failed for notification %s: %v",
			attempt+1, notification.ID, err)
	}

	// All retries exhausted
	s.repo.UpdateStatus(notification.ID, models.StatusFailed, lastErr.Error())
	return lastErr
}

// IsRetryable checks if a notification should be retried
func (s *RetryService) IsRetryable(notification *models.Notification) bool {
	return notification.RetryCount < MaxRetries &&
		notification.Status == models.StatusFailed
}
