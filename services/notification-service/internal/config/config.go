package config

import (
	"os"
	"strconv"
)

type Config struct {
	SMTP     SMTPConfig
	Twilio   TwilioConfig
	Firebase FirebaseConfig
	RabbitMQ RabbitMQConfig
	Database DatabaseConfig
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

type TwilioConfig struct {
	AccountSID string
	AuthToken  string
	FromNumber string
}

type FirebaseConfig struct {
	ProjectID      string
	CredentialFile string
}

type RabbitMQConfig struct {
	URL string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func Load() *Config {
	return &Config{
		SMTP: SMTPConfig{
			Host:     getEnv("SMTP_HOST", "smtp.gmail.com"),
			Port:     getEnvAsInt("SMTP_PORT", 587),
			Username: getEnv("SMTP_USERNAME", ""),
			Password: getEnv("SMTP_PASSWORD", ""),
			From:     getEnv("SMTP_FROM", "noreply@ecommerce.com"),
		},
		Twilio: TwilioConfig{
			AccountSID: getEnv("TWILIO_ACCOUNT_SID", ""),
			AuthToken:  getEnv("TWILIO_AUTH_TOKEN", ""),
			FromNumber: getEnv("TWILIO_FROM_NUMBER", ""),
		},
		Firebase: FirebaseConfig{
			ProjectID:      getEnv("FIREBASE_PROJECT_ID", ""),
			CredentialFile: getEnv("FIREBASE_CREDENTIAL_FILE", ""),
		},
		RabbitMQ: RabbitMQConfig{
			URL: getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "notifications"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
