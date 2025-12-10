-- The product_service database is created automatically by POSTGRES_DB in docker-compose
-- We just need to switch to it and create the tables

\c product_service

CREATE TABLE categories (
	category_id SERIAL PRIMARY KEY,
	parent_category_id INT,
	category_title VARCHAR(255),
	image_url VARCHAR(255),
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE products (
	product_id SERIAL PRIMARY KEY,
	category_id INT,
	product_title VARCHAR(255),
	image_url VARCHAR(255),
	sku VARCHAR(255),
	price_unit DECIMAL(10, 2),
	quantity INT,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT fk_category FOREIGN KEY (category_id) REFERENCES categories (category_id)
);

-- Now create the payment_service database and tables

CREATE DATABASE payment_service;
GRANT ALL PRIVILEGES ON DATABASE payment_service TO postgres;

\c payment_service

CREATE TABLE payments (
	payment_id SERIAL PRIMARY KEY,
	order_id INT,
	is_payed BOOLEAN,
	payment_status VARCHAR(255),
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
