package models

import (
	"time"
)

type NotificationType string

const (
	NotificationTypeEmail NotificationType = "email"
	NotificationTypeSMS   NotificationType = "sms"
	NotificationTypePush  NotificationType = "push"
)

type NotificationStatus string

const (
	StatusPending   NotificationStatus = "pending"
	StatusSent      NotificationStatus = "sent"
	StatusFailed    NotificationStatus = "failed"
	StatusDelivered NotificationStatus = "delivered"
)

type Notification struct {
	ID          string             `json:"id" db:"id"`
	UserID      string             `json:"user_id" db:"user_id"`
	Type        NotificationType   `json:"type" db:"type"`
	Channel     string             `json:"channel" db:"channel"`
	Recipient   string             `json:"recipient" db:"recipient"`
	Subject     string             `json:"subject,omitempty" db:"subject"`
	Content     string             `json:"content" db:"content"`
	TemplateID  string             `json:"template_id,omitempty" db:"template_id"`
	Metadata    map[string]any     `json:"metadata,omitempty" db:"metadata"`
	Status      NotificationStatus `json:"status" db:"status"`
	SentAt      *time.Time         `json:"sent_at,omitempty" db:"sent_at"`
	DeliveredAt *time.Time         `json:"delivered_at,omitempty" db:"delivered_at"`
	FailedAt    *time.Time         `json:"failed_at,omitempty" db:"failed_at"`
	ErrorMsg    string             `json:"error_msg,omitempty" db:"error_msg"`
	RetryCount  int                `json:"retry_count" db:"retry_count"`
	CreatedAt   time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" db:"updated_at"`
}

// E-commerce Event Models
type OrderEvent struct {
	EventType string    `json:"event_type"`
	OrderID   string    `json:"order_id"`
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone,omitempty"`
	Items     []Item    `json:"items"`
	Total     float64   `json:"total"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

type Item struct {
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
}

type UserEvent struct {
	EventType string    `json:"event_type"`
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone,omitempty"`
	Name      string    `json:"name"`
	Timestamp time.Time `json:"timestamp"`
}

type InventoryEvent struct {
	EventType   string    `json:"event_type"`
	ProductID   string    `json:"product_id"`
	ProductName string    `json:"product_name"`
	Quantity    int       `json:"quantity"`
	Threshold   int       `json:"threshold"`
	Timestamp   time.Time `json:"timestamp"`
}
