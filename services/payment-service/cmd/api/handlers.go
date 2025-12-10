package main

import (
	"errors"
	"net/http"
	"strconv"

	"payment/data"

	"github.com/go-chi/chi/v5"
)

type DtoCollectionResponse struct {
	Collection interface{} `json:"collection"`
}

func (app *Config) GetAllPayments(w http.ResponseWriter, r *http.Request) {
	payments, err := app.Models.Payment.GetAll()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := DtoCollectionResponse{
		Collection: payments,
	}

	app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) GetPayment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "paymentId")
	paymentID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, errors.New("invalid payment id"))
		return
	}

	payment, err := app.Models.Payment.GetOne(paymentID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, payment)
}

func (app *Config) CreatePayment(w http.ResponseWriter, r *http.Request) {
	var payment data.Payment
	err := app.readJSON(w, r, &payment)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	newPayment, err := app.Models.Payment.Insert(payment)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, newPayment)
}

func (app *Config) UpdatePayment(w http.ResponseWriter, r *http.Request) {
	var payment data.Payment
	err := app.readJSON(w, r, &payment)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	updatedPayment, err := app.Models.Payment.Update(payment)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, updatedPayment)
}

func (app *Config) DeletePayment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "paymentId")
	paymentID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, errors.New("invalid payment id"))
		return
	}

	err = app.Models.Payment.Delete(paymentID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, true)
}
