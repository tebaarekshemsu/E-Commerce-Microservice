-- migrations/001_initial_schema.sql

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
    marketing_enabled BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);