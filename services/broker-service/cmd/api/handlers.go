package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"time"

	paymentpb "broker/proto/payment"
	productpb "broker/proto/product"

	"github.com/go-chi/chi/v5"
)

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Broker hit",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) ProductServiceProxy() http.Handler {
	target, _ := url.Parse("http://product-service")
	return httputil.NewSingleHostReverseProxy(target)
}

func (app *Config) PaymentServiceProxy() http.Handler {
	target, _ := url.Parse("http://payment-service")
	return httputil.NewSingleHostReverseProxy(target)
}
func (app *Config) UserServiceProxy() http.Handler {
	target, _ := url.Parse("http://user-service:8000")
	return httputil.NewSingleHostReverseProxy(target)
}
func (app *Config) OrderServiceProxy() http.Handler {
	target, _ := url.Parse("http://order-service")
	return httputil.NewSingleHostReverseProxy(target)
}

func (app *Config) NotificationServiceProxy() http.Handler {
	target, _ := url.Parse("http://notification-service")
	return httputil.NewSingleHostReverseProxy(target)
}

func (app *Config) GRPCHealth(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	status := map[string]string{
		"product": "unavailable",
		"payment": "unavailable",
	}

	if app.ProductClient != nil {
		if _, err := app.ProductClient.CheckAvailability(ctx, &productpb.CheckAvailabilityRequest{ProductId: "1", Quantity: 1}); err == nil {
			status["product"] = "ok"
		}
	}

	if app.PaymentClient != nil {
		if _, err := app.PaymentClient.ProcessPayment(ctx, &paymentpb.PaymentRequest{OrderId: "0", Amount: 0, Currency: "USD", PaymentMethod: "health"}); err == nil {
			status["payment"] = "ok"
		}
	}

	payload := jsonResponse{Error: false, Message: "grpc health", Data: status}
	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) GetProductGRPC(w http.ResponseWriter, r *http.Request) {
	if app.ProductClient == nil {
		_ = app.errorJSON(w, errors.New("product gRPC client unavailable"), http.StatusServiceUnavailable)
		return
	}

	productID := chi.URLParam(r, "id")
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	resp, err := app.ProductClient.GetProduct(ctx, &productpb.GetProductRequest{Id: productID})
	if err != nil {
		_ = app.errorJSON(w, err, http.StatusBadGateway)
		return
	}

	payload := jsonResponse{Error: false, Message: "product fetched", Data: resp}
	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) CheckProductAvailability(w http.ResponseWriter, r *http.Request) {
	if app.ProductClient == nil {
		_ = app.errorJSON(w, errors.New("product gRPC client unavailable"), http.StatusServiceUnavailable)
		return
	}

	productID := chi.URLParam(r, "id")
	qtyStr := r.URL.Query().Get("quantity")
	qty, err := strconv.Atoi(qtyStr)
	if err != nil || qty <= 0 {
		_ = app.errorJSON(w, errors.New("quantity must be a positive integer"), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	resp, err := app.ProductClient.CheckAvailability(ctx, &productpb.CheckAvailabilityRequest{
		ProductId: productID,
		Quantity:  int32(qty),
	})
	if err != nil {
		_ = app.errorJSON(w, err, http.StatusBadGateway)
		return
	}

	payload := jsonResponse{Error: false, Message: "availability checked", Data: resp}
	_ = app.writeJSON(w, http.StatusOK, payload)
}

type paymentRequestPayload struct {
	OrderID       string  `json:"order_id"`
	Amount        float32 `json:"amount"`
	Currency      string  `json:"currency"`
	PaymentMethod string  `json:"payment_method"`
}

func (app *Config) ProcessPaymentGRPC(w http.ResponseWriter, r *http.Request) {
	if app.PaymentClient == nil {
		_ = app.errorJSON(w, errors.New("payment gRPC client unavailable"), http.StatusServiceUnavailable)
		return
	}

	var req paymentRequestPayload
	if err := app.readJSON(w, r, &req); err != nil {
		_ = app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	if req.OrderID == "" {
		_ = app.errorJSON(w, errors.New("order_id is required"), http.StatusBadRequest)
		return
	}

	if req.Currency == "" {
		req.Currency = "USD"
	}
	if req.PaymentMethod == "" {
		req.PaymentMethod = "credit_card"
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := app.PaymentClient.ProcessPayment(ctx, &paymentpb.PaymentRequest{
		OrderId:       req.OrderID,
		Amount:        req.Amount,
		Currency:      req.Currency,
		PaymentMethod: req.PaymentMethod,
	})
	if err != nil {
		_ = app.errorJSON(w, err, http.StatusBadGateway)
		return
	}

	payload := jsonResponse{Error: false, Message: "payment processed", Data: resp}
	_ = app.writeJSON(w, http.StatusOK, payload)
}
