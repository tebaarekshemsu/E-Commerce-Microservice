package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"

	"product-service/data"
	pb "product-service/proto"

	"google.golang.org/grpc"
)

func StartGRPCServer(models data.Models, port string) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	pb.RegisterProductServiceServer(s, &ProductServer{models: models})

	log.Printf("gRPC server listening on %v", lis.Addr())
	return s.Serve(lis)
}

type ProductServer struct {
	pb.UnimplementedProductServiceServer
	models data.Models
}

func NewProductServer(models data.Models) *ProductServer {
	return &ProductServer{
		models: models,
	}
}

func (s *ProductServer) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.ProductResponse, error) {
	log.Printf("gRPC GetProduct called with ID: %s", req.Id)

	// Convert string ID to int
	id, err := strconv.Atoi(req.Id)
	if err != nil {
		return nil, fmt.Errorf("invalid product ID: %v", err)
	}

	// Get product from database
	product, err := s.models.Product.GetOne(id)
	if err != nil {
		log.Printf("Error getting product: %v", err)
		return nil, fmt.Errorf("product not found: %v", err)
	}

	// Convert to protobuf response
	response := &pb.ProductResponse{
		Id:          strconv.Itoa(product.ID),
		Name:        product.Title,
		Description: fmt.Sprintf("SKU: %s", product.SKU), // Using SKU as description for now
		Price:       float32(product.PriceUnit),
		Stock:       int32(product.Quantity),
	}

	return response, nil
}

func (s *ProductServer) CheckAvailability(ctx context.Context, req *pb.CheckAvailabilityRequest) (*pb.CheckAvailabilityResponse, error) {
	log.Printf("gRPC CheckAvailability called with ProductID: %s, Quantity: %d", req.ProductId, req.Quantity)

	// Convert string ID to int
	id, err := strconv.Atoi(req.ProductId)
	if err != nil {
		return nil, fmt.Errorf("invalid product ID: %v", err)
	}

	// Get product from database
	product, err := s.models.Product.GetOne(id)
	if err != nil {
		// If product not found, it's not available
		return &pb.CheckAvailabilityResponse{Available: false}, nil
	}

	// Check stock
	if int32(product.Quantity) >= req.Quantity {
		return &pb.CheckAvailabilityResponse{Available: true}, nil
	}

	return &pb.CheckAvailabilityResponse{Available: false}, nil
}
