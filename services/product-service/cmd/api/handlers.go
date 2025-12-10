package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"product/data"
)

type DtoCollectionResponse struct {
	Collection interface{} `json:"collection"`
}

// Product Handlers

func (app *Config) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := app.Models.Product.GetAll()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := DtoCollectionResponse{
		Collection: products,
	}

	app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "productId")
	productID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, errors.New("invalid product id"))
		return
	}

	product, err := app.Models.Product.GetOne(productID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, product)
}

func (app *Config) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product data.Product
	err := app.readJSON(w, r, &product)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	newProduct, err := app.Models.Product.Insert(product)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, newProduct)
}

func (app *Config) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	var product data.Product
	err := app.readJSON(w, r, &product)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	updatedProduct, err := app.Models.Product.Update(product)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, updatedProduct)
}

func (app *Config) UpdateProductWithID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "productId")
	productID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, errors.New("invalid product id"))
		return
	}

	var product data.Product
	err = app.readJSON(w, r, &product)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	product.ID = productID
	updatedProduct, err := app.Models.Product.Update(product)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, updatedProduct)
}

func (app *Config) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "productId")
	productID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, errors.New("invalid product id"))
		return
	}

	err = app.Models.Product.Delete(productID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, true)
}

// Category Handlers

func (app *Config) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := app.Models.Category.GetAll()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := DtoCollectionResponse{
		Collection: categories,
	}

	app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) GetCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "categoryId")
	categoryID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, errors.New("invalid category id"))
		return
	}

	category, err := app.Models.Category.GetOne(categoryID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, category)
}

func (app *Config) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var category data.Category
	err := app.readJSON(w, r, &category)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	newCategory, err := app.Models.Category.Insert(category)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, newCategory)
}

func (app *Config) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	var category data.Category
	err := app.readJSON(w, r, &category)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	updatedCategory, err := app.Models.Category.Update(category)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, updatedCategory)
}

func (app *Config) UpdateCategoryWithID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "categoryId")
	categoryID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, errors.New("invalid category id"))
		return
	}

	var category data.Category
	err = app.readJSON(w, r, &category)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	category.ID = categoryID
	updatedCategory, err := app.Models.Category.Update(category)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, updatedCategory)
}

func (app *Config) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "categoryId")
	categoryID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, errors.New("invalid category id"))
		return
	}

	err = app.Models.Category.Delete(categoryID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, true)
}
