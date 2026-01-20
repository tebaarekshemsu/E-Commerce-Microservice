package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"payment-service/data"
	pb "payment-service/proto"
	"payment-service/queue"

	"google.golang.org/grpc"
)

func StartGRPCServer(models data.Models, publisher *queue.Publisher, port string) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	pb.RegisterPaymentServiceServer(s, &PaymentServer{models: models, publisher: publisher})

	log.Printf("gRPC server listening on %v", lis.Addr())
	return s.Serve(lis)
}

type PaymentServer struct {
	pb.UnimplementedPaymentServiceServer
	models    data.Models
	publisher *queue.Publisher
}

func NewPaymentServer(models data.Models, publisher *queue.Publisher) *PaymentServer {
	return &PaymentServer{
		models:    models,
		publisher: publisher,
	}
}

func (s *PaymentServer) ProcessPayment(ctx context.Context, req *pb.PaymentRequest) (*pb.PaymentResponse, error) {
	log.Printf("gRPC ProcessPayment called for order %s, amount %.2f %s",
		req.OrderId, req.Amount, req.Currency)

	// Convert order ID to int
	orderID, err := strconv.Atoi(req.OrderId)
	if err != nil {
		return &pb.PaymentResponse{
			Success: false,
			Message: fmt.Sprintf("Invalid order ID: %v", err),
		}, nil
	}

	// Create payment record
	payment := data.Payment{
		OrderID:       orderID,
		IsPayed:       false,
		PaymentStatus: "processing",
	}

	// Insert payment record
	createdPayment, err := s.models.Payment.Insert(payment)
	if err != nil {
		log.Printf("Error creating payment record: %v", err)
		return &pb.PaymentResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to create payment record: %v", err),
		}, nil
	}

	// Simulate payment processing
	time.Sleep(100 * time.Millisecond)

	// For demo purposes, randomly succeed/fail based on amount
	success := req.Amount > 0 && req.Amount < 10000 // Fail if amount is too high

	// Update payment status
	createdPayment.IsPayed = success
	if success {
		createdPayment.PaymentStatus = "completed"
	} else {
		createdPayment.PaymentStatus = "failed"
	}

	_, err = s.models.Payment.Update(*createdPayment)
	if err != nil {
		log.Printf("Error updating payment status: %v", err)
	}

	transactionID := fmt.Sprintf("txn_%d_%d", createdPayment.ID, time.Now().Unix())
	message := "Payment completed successfully"
	if !success {
		message = "Payment failed - amount too high or invalid"
	}

	// Publish payment event if publisher is available
	if s.publisher != nil {
		event := map[string]interface{}{
			"order_id":       req.OrderId,
			"amount":         req.Amount,
			"currency":       req.Currency,
			"payment_method": req.PaymentMethod,
			"success":        success,
			"transaction_id": transactionID,
			"message":        message,
			"timestamp":      time.Now().UTC(),
		}

		if err := s.publisher.Push("payments", event); err != nil {
			log.Printf("Failed to publish payment event: %v", err)
		}
	}

	response := &pb.PaymentResponse{
		Success:       success,
		TransactionId: transactionID,
		Message:       message,
	}

	log.Printf("Payment processed: success=%t, transaction_id=%s", success, transactionID)

	return response, nil
}
