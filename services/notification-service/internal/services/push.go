package services

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"

	"notification-service/internal/config"
	"notification-service/internal/models"
)

type PushService struct {
	config config.FirebaseConfig
	client *messaging.Client
	ctx    context.Context
}

func NewPushService(cfg config.FirebaseConfig) *PushService {
	ctx := context.Background()

	opt := option.WithCredentialsFile(cfg.CredentialFile)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Printf("Error initializing Firebase app: %v", err)
		return &PushService{config: cfg, ctx: ctx}
	}

	client, err := app.Messaging(ctx)
	if err != nil {
		log.Printf("Error getting Messaging client: %v", err)
		return &PushService{config: cfg, ctx: ctx}
	}

	return &PushService{
		config: cfg,
		client: client,
		ctx:    ctx,
	}
}

func (p *PushService) Send(notification *models.Notification) error {
	if p.client == nil {
		return fmt.Errorf("firebase client not initialized")
	}

	// Build FCM message
	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: notification.Subject,
			Body:  notification.Content,
		},
		Token: notification.Recipient, // Device token
		Data:  p.convertMetadata(notification.Metadata),
	}

	// Send message
	response, err := p.client.Send(p.ctx, message)
	if err != nil {
		return fmt.Errorf("failed to send push notification: %w", err)
	}

	log.Printf("Successfully sent push notification: %s", response)
	return nil
}

// SendToTopic sends a notification to all devices subscribed to a topic
func (p *PushService) SendToTopic(topic string, notification *models.Notification) error {
	if p.client == nil {
		return fmt.Errorf("firebase client not initialized")
	}

	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: notification.Subject,
			Body:  notification.Content,
		},
		Topic: topic,
		Data:  p.convertMetadata(notification.Metadata),
	}

	response, err := p.client.Send(p.ctx, message)
	if err != nil {
		return fmt.Errorf("failed to send topic notification: %w", err)
	}

	log.Printf("Successfully sent topic notification: %s", response)
	return nil
}

// SendMulticast sends to multiple device tokens
func (p *PushService) SendMulticast(tokens []string, notification *models.Notification) (*messaging.BatchResponse, error) {
	if p.client == nil {
		return nil, fmt.Errorf("firebase client not initialized")
	}

	message := &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title: notification.Subject,
			Body:  notification.Content,
		},
		Tokens: tokens,
		Data:   p.convertMetadata(notification.Metadata),
	}

	return p.client.SendEachForMulticast(p.ctx, message)
}

func (p *PushService) convertMetadata(metadata map[string]any) map[string]string {
	result := make(map[string]string)
	for k, v := range metadata {
		result[k] = fmt.Sprintf("%v", v)
	}
	return result
}
