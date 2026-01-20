package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"payment-service/data"
	grpcServer "payment-service/grpc"
	"payment-service/queue"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	amqp "github.com/rabbitmq/amqp091-go"
)

const webPort = "80"
const grpcPort = "9002"

type Config struct {
	Models data.Models
}

func main() {
	log.Println("Starting payment service")

	// Connect to DB
	dsn := os.Getenv("DSN")
	if dsn == "" {
		dsn = "host=localhost port=5432 user=postgres password=password dbname=payment_service sslmode=disable timezone=UTC connect_timeout=5"
	}

	conn, err := connectToDB(dsn)
	if err != nil {
		log.Panic(err)
	}

	app := Config{
		Models: data.New(conn),
	}

	// Connect to RabbitMQ
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@rabbitmq:5672/"
	}

	amqpConn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer amqpConn.Close()

	publisher := queue.NewPublisher(amqpConn)

	// Start gRPC server in a goroutine
	go func() {
		if err := grpcServer.StartGRPCServer(app.Models, publisher, grpcPort); err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	log.Printf("Starting payment service on HTTP port %s and gRPC port %s\n", webPort, grpcPort)
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB(dsn string) (*sql.DB, error) {
	counts := 0

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready...")
			counts++
		} else {
			log.Println("Connected to Postgres!")
			return connection, nil
		}

		if counts > 10 {
			log.Println(err)
			return nil, err
		}

		log.Println("Backing off for two seconds...")
		time.Sleep(2 * time.Second)
		continue
	}
}
