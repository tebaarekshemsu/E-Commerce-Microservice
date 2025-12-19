
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

-- Now create the notifications database and tables

CREATE DATABASE notifications;

\c notifications

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE notification_type AS ENUM ('email', 'sms', 'push');
CREATE TYPE notification_status AS ENUM ('pending', 'sent', 'failed', 'delivered');

CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id VARCHAR(255) NOT NULL,
    type notification_type NOT NULL,
    channel VARCHAR(50) NOT NULL,
    recipient VARCHAR(255) NOT NULL,
    subject VARCHAR(255),
    content TEXT NOT NULL,
    template_id VARCHAR(100),
    metadata JSONB DEFAULT '{}',
    status notification_status DEFAULT 'pending',
    sent_at TIMESTAMP,
    delivered_at TIMESTAMP,
    failed_at TIMESTAMP,
    error_msg TEXT,
    retry_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_status ON notifications(status);
CREATE INDEX idx_notifications_created_at ON notifications(created_at);
CREATE INDEX idx_notifications_type ON notifications(type);

-- Templates table
CREATE TABLE notification_templates (
    id VARCHAR(100) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type notification_type NOT NULL,
    subject VARCHAR(255),
    content TEXT NOT NULL,
    variables JSONB DEFAULT '[]',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- User preferences table
CREATE TABLE user_notification_preferences (
    user_id VARCHAR(255) PRIMARY KEY,
    email_enabled BOOLEAN DEFAULT true,
    sms_enabled BOOLEAN DEFAULT false,
    push_enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Now create the user_service database
CREATE DATABASE user_service_db;
GRANT ALL PRIVILEGES ON DATABASE user_service_db TO postgres;

CREATE DATABASE user_service_db;
