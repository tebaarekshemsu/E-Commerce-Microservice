package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"notification-service/internal/config"
	"notification-service/internal/handlers"
	"notification-service/internal/queue"
	"notification-service/internal/services"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize services
	emailSvc := services.NewEmailService(cfg.SMTP)
	smsSvc := services.NewSMSService(cfg.Twilio)
	pushSvc := services.NewPushService(cfg.Firebase)

	// Initialize message queue
	messageQueue := queue.NewRabbitMQ(cfg.RabbitMQ)
	defer messageQueue.Close()

	// Create notification handler
	notificationHandler := handlers.NewNotificationHandler(
		emailSvc, smsSvc, pushSvc,
	)

	// Start consuming messages
	go messageQueue.Consume("orders", notificationHandler.HandleOrderEvent)
	go messageQueue.Consume("users", notificationHandler.HandleUserEvent)
	go messageQueue.Consume("inventory", notificationHandler.HandleInventoryEvent)

	log.Println("🚀 Notification Service started successfully")

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down notification service...")
}
