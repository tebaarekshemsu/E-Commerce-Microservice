package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"notification-service/internal/config"
	"notification-service/internal/models"
)

type NotificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(cfg config.DatabaseConfig) (*NotificationRepository, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &NotificationRepository{db: db}, nil
}

func (r *NotificationRepository) Create(n *models.Notification) error {
	metadata, _ := json.Marshal(n.Metadata)

	query := `
        INSERT INTO notifications (
            id, user_id, type, channel, recipient, subject, 
            content, template_id, metadata, status, created_at, updated_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
    `

	_, err := r.db.Exec(query,
		n.ID, n.UserID, n.Type, n.Channel, n.Recipient, n.Subject,
		n.Content, n.TemplateID, metadata, n.Status, n.CreatedAt, n.UpdatedAt,
	)

	return err
}

func (r *NotificationRepository) GetByID(id string) (*models.Notification, error) {
	query := `
        SELECT id, user_id, type, channel, recipient, subject, content, 
               template_id, metadata, status, sent_at, delivered_at, 
               failed_at, error_msg, retry_count, created_at, updated_at
        FROM notifications WHERE id = $1
    `

	var n models.Notification
	var metadata []byte
	var sentAt, deliveredAt, failedAt sql.NullTime

	err := r.db.QueryRow(query, id).Scan(
		&n.ID, &n.UserID, &n.Type, &n.Channel, &n.Recipient, &n.Subject,
		&n.Content, &n.TemplateID, &metadata, &n.Status, &sentAt, &deliveredAt,
		&failedAt, &n.ErrorMsg, &n.RetryCount, &n.CreatedAt, &n.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	json.Unmarshal(metadata, &n.Metadata)
	if sentAt.Valid {
		n.SentAt = &sentAt.Time
	}
	if deliveredAt.Valid {
		n.DeliveredAt = &deliveredAt.Time
	}
	if failedAt.Valid {
		n.FailedAt = &failedAt.Time
	}

	return &n, nil
}

func (r *NotificationRepository) UpdateStatus(id string, status models.NotificationStatus, errorMsg string) error {
	now := time.Now()
	var query string

	switch status {
	case models.StatusSent:
		query = `UPDATE notifications SET status = $1, sent_at = $2, updated_at = $2 WHERE id = $3`
	case models.StatusDelivered:
		query = `UPDATE notifications SET status = $1, delivered_at = $2, updated_at = $2 WHERE id = $3`
	case models.StatusFailed:
		query = `UPDATE notifications SET status = $1, failed_at = $2, error_msg = $3, updated_at = $2 WHERE id = $4`
		_, err := r.db.Exec(query, status, now, errorMsg, id)
		return err
	default:
		query = `UPDATE notifications SET status = $1, updated_at = $2 WHERE id = $3`
	}

	_, err := r.db.Exec(query, status, now, id)
	return err
}

func (r *NotificationRepository) GetByUserID(userID string, limit, offset int) ([]*models.Notification, error) {
	query := `
        SELECT id, user_id, type, channel, recipient, subject, content,
               status, created_at
        FROM notifications 
        WHERE user_id = $1 
        ORDER BY created_at DESC 
        LIMIT $2 OFFSET $3
    `

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []*models.Notification
	for rows.Next() {
		var n models.Notification
		err := rows.Scan(
			&n.ID, &n.UserID, &n.Type, &n.Channel, &n.Recipient,
			&n.Subject, &n.Content, &n.Status, &n.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, &n)
	}

	return notifications, nil
}

func (r *NotificationRepository) IncrementRetryCount(id string) error {
	query := `UPDATE notifications SET retry_count = retry_count + 1, updated_at = $1 WHERE id = $2`
	_, err := r.db.Exec(query, time.Now(), id)
	return err
}

func (r *NotificationRepository) Close() error {
	return r.db.Close()
}
