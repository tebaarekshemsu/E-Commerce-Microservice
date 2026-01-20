package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	paymentpb "broker/proto/payment"
	productpb "broker/proto/product"

	"google.golang.org/grpc"
)

const webPort = "80"

type Config struct {
	ProductClient productpb.ProductServiceClient
	PaymentClient paymentpb.PaymentServiceClient
}

func main() {
	log.Printf("Starting broker service on port %s\n", webPort)

	productAddr := os.Getenv("PRODUCT_SERVICE_GRPC_URL")
	if productAddr == "" {
		productAddr = "product-service:9001"
	}

	paymentAddr := os.Getenv("PAYMENT_SERVICE_GRPC_URL")
	if paymentAddr == "" {
		paymentAddr = "payment-service:9002"
	}

	productCtx, productCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer productCancel()

	productConn, err := grpc.DialContext(productCtx, productAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("failed to connect to product gRPC service: %v", err)
	}
	defer productConn.Close()

	paymentCtx, paymentCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer paymentCancel()

	paymentConn, err := grpc.DialContext(paymentCtx, paymentAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("failed to connect to payment gRPC service: %v", err)
	}
	defer paymentConn.Close()

	app := Config{
		ProductClient: productpb.NewProductServiceClient(productConn),
		PaymentClient: paymentpb.NewPaymentServiceClient(paymentConn),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}
