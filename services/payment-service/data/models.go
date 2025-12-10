package data

import (
	"context"
	"database/sql"
	"time"
)

const dbTimeout = time.Second * 3

type Models struct {
	Payment PaymentModel
}

func New(db *sql.DB) Models {
	return Models{
		Payment: PaymentModel{DB: db},
	}
}

type Payment struct {
	ID            int       `json:"paymentId"`
	OrderID       int       `json:"orderId"`
	IsPayed       bool      `json:"isPayed"`
	PaymentStatus string    `json:"paymentStatus"`
	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
}

type PaymentModel struct {
	DB *sql.DB
}

func (m *PaymentModel) GetAll() ([]*Payment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT payment_id, order_id, is_payed, payment_status, created_at, updated_at
		FROM payments
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []*Payment

	for rows.Next() {
		var p Payment
		err := rows.Scan(
			&p.ID,
			&p.OrderID,
			&p.IsPayed,
			&p.PaymentStatus,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		payments = append(payments, &p)
	}

	return payments, nil
}

func (m *PaymentModel) GetOne(id int) (*Payment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT payment_id, order_id, is_payed, payment_status, created_at, updated_at
		FROM payments
		WHERE payment_id = $1
	`

	var p Payment
	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&p.ID,
		&p.OrderID,
		&p.IsPayed,
		&p.PaymentStatus,
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (m *PaymentModel) Insert(payment Payment) (*Payment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		INSERT INTO payments (order_id, is_payed, payment_status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING payment_id
	`

	err := m.DB.QueryRowContext(ctx, query,
		payment.OrderID,
		payment.IsPayed,
		payment.PaymentStatus,
		time.Now(),
		time.Now(),
	).Scan(&payment.ID)

	if err != nil {
		return nil, err
	}

	return &payment, nil
}

func (m *PaymentModel) Update(payment Payment) (*Payment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		UPDATE payments
		SET order_id = $1, is_payed = $2, payment_status = $3, updated_at = $4
		WHERE payment_id = $5
	`

	_, err := m.DB.ExecContext(ctx, query,
		payment.OrderID,
		payment.IsPayed,
		payment.PaymentStatus,
		time.Now(),
		payment.ID,
	)

	if err != nil {
		return nil, err
	}

	return &payment, nil
}

func (m *PaymentModel) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `DELETE FROM payments WHERE payment_id = $1`

	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}
