package queue

import (
	"log"

	"notification-service/internal/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQ(cfg config.RabbitMQConfig) *RabbitMQ {
	conn, err := amqp.Dial(cfg.URL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open channel: %v", err)
	}

	return &RabbitMQ{
		conn:    conn,
		channel: ch,
	}
}

func (r *RabbitMQ) DeclareQueue(name string) error {
	_, err := r.channel.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	return err
}

func (r *RabbitMQ) Consume(queueName string, handler func([]byte) error) error {
	if err := r.DeclareQueue(queueName); err != nil {
		return err
	}

	msgs, err := r.channel.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			if err := handler(msg.Body); err != nil {
				log.Printf("Error processing message: %v", err)
				msg.Nack(false, true) // Requeue on failure
			} else {
				msg.Ack(false)
			}
		}
	}()

	log.Printf("📥 Started consuming from queue: %s", queueName)
	return nil
}

func (r *RabbitMQ) Publish(queueName string, body []byte) error {
	if err := r.DeclareQueue(queueName); err != nil {
		return err
	}

	return r.channel.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (r *RabbitMQ) Close() {
	r.channel.Close()
	r.conn.Close()
}
