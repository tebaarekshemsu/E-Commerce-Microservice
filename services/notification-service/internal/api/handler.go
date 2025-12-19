package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"notification-service/internal/models"
	"notification-service/internal/repository"
	"notification-service/internal/services"

	"github.com/google/uuid"
)

type APIHandler struct {
	repo         *repository.NotificationRepository
	emailService *services.EmailService
	smsService   *services.SMSService
	pushService  *services.PushService
}

func NewAPIHandler(
	repo *repository.NotificationRepository,
	email *services.EmailService,
	sms *services.SMSService,
	push *services.PushService,
) *APIHandler {
	return &APIHandler{
		repo:         repo,
		emailService: email,
		smsService:   sms,
		pushService:  push,
	}
}

// SendNotificationRequest represents the API request body
type SendNotificationRequest struct {
	UserID    string         `json:"user_id"`
	Type      string         `json:"type"`      // email, sms, push
	Recipient string         `json:"recipient"` // email, phone, or device token
	Subject   string         `json:"subject,omitempty"`
	Content   string         `json:"content"`
	Template  string         `json:"template,omitempty"`
	Data      map[string]any `json:"data,omitempty"`
}

// SendNotification handles POST /api/notifications
func (h *APIHandler) SendNotification(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SendNotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	// Validate request
	if req.Type == "" || req.Recipient == "" || req.Content == "" {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing required fields"})
		return
	}

	// Create notification
	notification := &models.Notification{
		ID:         uuid.New().String(),
		UserID:     req.UserID,
		Type:       models.NotificationType(req.Type),
		Channel:    req.Type,
		Recipient:  req.Recipient,
		Subject:    req.Subject,
		Content:    req.Content,
		TemplateID: req.Template,
		Metadata:   req.Data,
		Status:     models.StatusPending,
	}

	// Save to database
	if err := h.repo.Create(notification); err != nil {
		log.Printf("Failed to save notification: %v", err)
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create notification"})
		return
	}

	// Send notification based on type
	var sendErr error
	switch notification.Type {
	case models.NotificationTypeEmail:
		sendErr = h.emailService.Send(notification)
	case models.NotificationTypeSMS:
		sendErr = h.smsService.Send(notification)
	case models.NotificationTypePush:
		sendErr = h.pushService.Send(notification)
	}

	if sendErr != nil {
		h.repo.UpdateStatus(notification.ID, models.StatusFailed, sendErr.Error())
		respondJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to send notification",
			"id":    notification.ID,
		})
		return
	}

	h.repo.UpdateStatus(notification.ID, models.StatusSent, "")

	respondJSON(w, http.StatusCreated, map[string]string{
		"id":      notification.ID,
		"status":  "sent",
		"message": "Notification sent successfully",
	})
}

// GetNotification handles GET /api/notifications/{id}
func (h *APIHandler) GetNotification(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Notification ID required"})
		return
	}

	notification, err := h.repo.GetByID(id)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}

	if notification == nil {
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "Notification not found"})
		return
	}

	respondJSON(w, http.StatusOK, notification)
}

// GetUserNotifications handles GET /api/users/{user_id}/notifications
func (h *APIHandler) GetUserNotifications(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "User ID required"})
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	notifications, err := h.repo.GetByUserID(userID, limit, offset)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"notifications": notifications,
		"limit":         limit,
		"offset":        offset,
	})
}

// HealthCheck handles GET /health
func (h *APIHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{
		"status":  "healthy",
		"service": "notification-service",
	})
}

func respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
