package data

import (
	"context"
	"database/sql"
	"time"
)

const dbTimeout = time.Second * 3

type Models struct {
	Product  ProductModel
	Category CategoryModel
}

func New(db *sql.DB) Models {
	return Models{
		Product:  ProductModel{DB: db},
		Category: CategoryModel{DB: db},
	}
}

type Product struct {
	ID        int       `json:"productId"`
	Title     string    `json:"productTitle"`
	ImageURL  string    `json:"imageUrl"`
	SKU       string    `json:"sku"`
	PriceUnit float64   `json:"priceUnit"`
	Quantity  int       `json:"quantity"`
	Category  *Category `json:"category"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type Category struct {
	ID             int       `json:"categoryId"`
	Title          string    `json:"categoryTitle"`
	ImageURL       string    `json:"imageUrl"`
	ParentCategory *Category `json:"parentCategory,omitempty"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
}

type ProductModel struct {
	DB *sql.DB
}

type CategoryModel struct {
	DB *sql.DB
}

// Product methods

func (m *ProductModel) GetAll() ([]*Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT p.product_id, p.product_title, p.image_url, p.sku, p.price_unit, p.quantity, p.created_at, p.updated_at,
		       c.category_id, c.category_title, c.image_url
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.category_id
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*Product

	for rows.Next() {
		var p Product
		var c Category
		var cID sql.NullInt32
		var cTitle sql.NullString
		var cImage sql.NullString

		err := rows.Scan(
			&p.ID,
			&p.Title,
			&p.ImageURL,
			&p.SKU,
			&p.PriceUnit,
			&p.Quantity,
			&p.CreatedAt,
			&p.UpdatedAt,
			&cID,
			&cTitle,
			&cImage,
		)
		if err != nil {
			return nil, err
		}

		if cID.Valid {
			c.ID = int(cID.Int32)
			c.Title = cTitle.String
			c.ImageURL = cImage.String
			p.Category = &c
		}

		products = append(products, &p)
	}

	return products, nil
}

func (m *ProductModel) GetOne(id int) (*Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT p.product_id, p.product_title, p.image_url, p.sku, p.price_unit, p.quantity, p.created_at, p.updated_at,
		       c.category_id, c.category_title, c.image_url
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.category_id
		WHERE p.product_id = $1
	`

	var p Product
	var c Category
	var cID sql.NullInt32
	var cTitle sql.NullString
	var cImage sql.NullString

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&p.ID,
		&p.Title,
		&p.ImageURL,
		&p.SKU,
		&p.PriceUnit,
		&p.Quantity,
		&p.CreatedAt,
		&p.UpdatedAt,
		&cID,
		&cTitle,
		&cImage,
	)

	if err != nil {
		return nil, err
	}

	if cID.Valid {
		c.ID = int(cID.Int32)
		c.Title = cTitle.String
		c.ImageURL = cImage.String
		p.Category = &c
	}

	return &p, nil
}

func (m *ProductModel) Insert(product Product) (*Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var categoryID *int
	if product.Category != nil {
		categoryID = &product.Category.ID
	}

	query := `
		INSERT INTO products (product_title, image_url, sku, price_unit, quantity, category_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING product_id
	`

	err := m.DB.QueryRowContext(ctx, query,
		product.Title,
		product.ImageURL,
		product.SKU,
		product.PriceUnit,
		product.Quantity,
		categoryID,
		time.Now(),
		time.Now(),
	).Scan(&product.ID)

	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (m *ProductModel) Update(product Product) (*Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var categoryID *int
	if product.Category != nil {
		categoryID = &product.Category.ID
	}

	query := `
		UPDATE products
		SET product_title = $1, image_url = $2, sku = $3, price_unit = $4, quantity = $5, category_id = $6, updated_at = $7
		WHERE product_id = $8
	`

	_, err := m.DB.ExecContext(ctx, query,
		product.Title,
		product.ImageURL,
		product.SKU,
		product.PriceUnit,
		product.Quantity,
		categoryID,
		time.Now(),
		product.ID,
	)

	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (m *ProductModel) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `DELETE FROM products WHERE product_id = $1`

	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}

// Category methods

func (m *CategoryModel) GetAll() ([]*Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT c.category_id, c.category_title, c.image_url, c.created_at, c.updated_at,
		       p.category_id, p.category_title, p.image_url
		FROM categories c
		LEFT JOIN categories p ON c.parent_category_id = p.category_id
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*Category

	for rows.Next() {
		var c Category
		var p Category
		var pID sql.NullInt32
		var pTitle sql.NullString
		var pImage sql.NullString

		err := rows.Scan(
			&c.ID,
			&c.Title,
			&c.ImageURL,
			&c.CreatedAt,
			&c.UpdatedAt,
			&pID,
			&pTitle,
			&pImage,
		)
		if err != nil {
			return nil, err
		}

		if pID.Valid {
			p.ID = int(pID.Int32)
			p.Title = pTitle.String
			p.ImageURL = pImage.String
			c.ParentCategory = &p
		}

		categories = append(categories, &c)
	}

	return categories, nil
}

func (m *CategoryModel) GetOne(id int) (*Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT c.category_id, c.category_title, c.image_url, c.created_at, c.updated_at,
		       p.category_id, p.category_title, p.image_url
		FROM categories c
		LEFT JOIN categories p ON c.parent_category_id = p.category_id
		WHERE c.category_id = $1
	`

	var c Category
	var p Category
	var pID sql.NullInt32
	var pTitle sql.NullString
	var pImage sql.NullString

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&c.ID,
		&c.Title,
		&c.ImageURL,
		&c.CreatedAt,
		&c.UpdatedAt,
		&pID,
		&pTitle,
		&pImage,
	)

	if err != nil {
		return nil, err
	}

	if pID.Valid {
		p.ID = int(pID.Int32)
		p.Title = pTitle.String
		p.ImageURL = pImage.String
		c.ParentCategory = &p
	}

	return &c, nil
}

func (m *CategoryModel) Insert(category Category) (*Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var parentID *int
	if category.ParentCategory != nil {
		parentID = &category.ParentCategory.ID
	}

	query := `
		INSERT INTO categories (category_title, image_url, parent_category_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING category_id
	`

	err := m.DB.QueryRowContext(ctx, query,
		category.Title,
		category.ImageURL,
		parentID,
		time.Now(),
		time.Now(),
	).Scan(&category.ID)

	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (m *CategoryModel) Update(category Category) (*Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var parentID *int
	if category.ParentCategory != nil {
		parentID = &category.ParentCategory.ID
	}

	query := `
		UPDATE categories
		SET category_title = $1, image_url = $2, parent_category_id = $3, updated_at = $4
		WHERE category_id = $5
	`

	_, err := m.DB.ExecContext(ctx, query,
		category.Title,
		category.ImageURL,
		parentID,
		time.Now(),
		category.ID,
	)

	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (m *CategoryModel) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `DELETE FROM categories WHERE category_id = $1`

	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}
