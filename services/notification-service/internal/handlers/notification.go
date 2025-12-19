package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"notification-service/internal/models"
	"notification-service/internal/services"

	"github.com/google/uuid"
)

type NotificationHandler struct {
	emailService *services.EmailService
	smsService   *services.SMSService
	pushService  *services.PushService
}

func NewNotificationHandler(
	email *services.EmailService,
	sms *services.SMSService,
	push *services.PushService,
) *NotificationHandler {
	return &NotificationHandler{
		emailService: email,
		smsService:   sms,
		pushService:  push,
	}
}

func (h *NotificationHandler) HandleOrderEvent(data []byte) error {
	var event models.OrderEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return fmt.Errorf("failed to unmarshal order event: %w", err)
	}

	log.Printf("Processing order event: %s for order %s",
		event.EventType, event.OrderID)

	switch event.EventType {
	case "order_created":
		return h.sendOrderConfirmation(event)
	case "order_shipped":
		return h.sendShippingNotification(event)
	case "order_delivered":
		return h.sendDeliveryNotification(event)
	case "order_cancelled":
		return h.sendCancellationNotification(event)
	default:
		log.Printf("Unknown order event type: %s", event.EventType)
		return nil
	}
}

func (h *NotificationHandler) sendOrderConfirmation(event models.OrderEvent) error {
	notification := &models.Notification{
		ID:         uuid.New().String(),
		UserID:     event.UserID,
		Type:       models.NotificationTypeEmail,
		Channel:    "email",
		Recipient:  event.Email,
		Subject:    fmt.Sprintf("Order Confirmation #%s", event.OrderID),
		TemplateID: "order_confirmation",
		Metadata: map[string]any{
			"order_id": event.OrderID,
			"items":    event.Items,
			"total":    event.Total,
		},
		Status:    models.StatusPending,
		CreatedAt: time.Now(),
	}

	if err := h.emailService.Send(notification); err != nil {
		log.Printf("Failed to send order confirmation: %v", err)
		return err
	}

	log.Printf("✅ Order confirmation sent for order %s", event.OrderID)
	return nil
}

func (h *NotificationHandler) HandleUserEvent(data []byte) error {
	var event models.UserEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return fmt.Errorf("failed to unmarshal user event: %w", err)
	}

	log.Printf("Processing user event: %s for user %s",
		event.EventType, event.UserID)

	switch event.EventType {
	case "user_registered":
		return h.sendWelcomeEmail(event)
	case "password_reset_requested":
		return h.sendPasswordResetEmail(event)
	default:
		return nil
	}
}

func (h *NotificationHandler) sendWelcomeEmail(event models.UserEvent) error {
	notification := &models.Notification{
		ID:         uuid.New().String(),
		UserID:     event.UserID,
		Type:       models.NotificationTypeEmail,
		Channel:    "email",
		Recipient:  event.Email,
		Subject:    "Welcome to Our Store!",
		TemplateID: "welcome",
		Metadata: map[string]any{
			"name": event.Name,
		},
		Status:    models.StatusPending,
		CreatedAt: time.Now(),
	}

	return h.emailService.Send(notification)
}

func (h *NotificationHandler) HandleInventoryEvent(data []byte) error {
	var event models.InventoryEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return fmt.Errorf("failed to unmarshal inventory event: %w", err)
	}

	if event.EventType == "low_stock" {
		return h.sendLowStockAlert(event)
	}

	return nil
}

func (h *NotificationHandler) sendLowStockAlert(event models.InventoryEvent) error {
	// Send alert to admin/inventory team
	notification := &models.Notification{
		ID:         uuid.New().String(),
		Type:       models.NotificationTypeEmail,
		Channel:    "email",
		Recipient:  "inventory@ecommerce.com",
		Subject:    fmt.Sprintf("Low Stock Alert: %s", event.ProductName),
		TemplateID: "low_stock_alert",
		Metadata: map[string]any{
			"product_id":   event.ProductID,
			"product_name": event.ProductName,
			"quantity":     event.Quantity,
			"threshold":    event.Threshold,
		},
		Status:    models.StatusPending,
		CreatedAt: time.Now(),
	}

	return h.emailService.Send(notification)
}

func (h *NotificationHandler) sendShippingNotification(event models.OrderEvent) error {
	notification := &models.Notification{
		ID:         uuid.New().String(),
		UserID:     event.UserID,
		Type:       models.NotificationTypeEmail,
		Channel:    "email",
		Recipient:  event.Email,
		Subject:    fmt.Sprintf("Your Order #%s has Shipped!", event.OrderID),
		TemplateID: "order_shipped",
		Metadata: map[string]any{
			"order_id": event.OrderID,
			"items":    event.Items,
		},
		Status:    models.StatusPending,
		CreatedAt: time.Now(),
	}

	if err := h.emailService.Send(notification); err != nil {
		log.Printf("Failed to send shipping notification: %v", err)
		return err
	}
	return nil
}

func (h *NotificationHandler) sendDeliveryNotification(event models.OrderEvent) error {
	notification := &models.Notification{
		ID:         uuid.New().String(),
		UserID:     event.UserID,
		Type:       models.NotificationTypeEmail,
		Channel:    "email",
		Recipient:  event.Email,
		Subject:    fmt.Sprintf("Order #%s Delivered", event.OrderID),
		TemplateID: "order_delivered",
		Metadata: map[string]any{
			"order_id": event.OrderID,
		},
		Status:    models.StatusPending,
		CreatedAt: time.Now(),
	}

	if err := h.emailService.Send(notification); err != nil {
		log.Printf("Failed to send delivery notification: %v", err)
		return err
	}
	return nil
}

func (h *NotificationHandler) sendCancellationNotification(event models.OrderEvent) error {
	notification := &models.Notification{
		ID:         uuid.New().String(),
		UserID:     event.UserID,
		Type:       models.NotificationTypeEmail,
		Channel:    "email",
		Recipient:  event.Email,
		Subject:    fmt.Sprintf("Order #%s Cancelled", event.OrderID),
		TemplateID: "order_cancelled",
		Metadata: map[string]any{
			"order_id": event.OrderID,
		},
		Status:    models.StatusPending,
		CreatedAt: time.Now(),
	}

	if err := h.emailService.Send(notification); err != nil {
		log.Printf("Failed to send cancellation notification: %v", err)
		return err
	}
	return nil
}

func (h *NotificationHandler) sendPasswordResetEmail(event models.UserEvent) error {
	notification := &models.Notification{
		ID:         uuid.New().String(),
		UserID:     event.UserID,
		Type:       models.NotificationTypeEmail,
		Channel:    "email",
		Recipient:  event.Email,
		Subject:    "Password Reset Request",
		TemplateID: "password_reset",
		Metadata: map[string]any{
			"reset_link": event.Metadata["reset_link"],
		},
		Status:    models.StatusPending,
		CreatedAt: time.Now(),
	}

	if err := h.emailService.Send(notification); err != nil {
		log.Printf("Failed to send password reset email: %v", err)
		return err
	}
	return nil
}
