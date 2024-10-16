package config

import (
	"log"
	"os"

	"github.com/streadway/amqp"
)

type RabbitMQConfig struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

// NewRabbitMQConfig initializes a new RabbitMQConfig instance
func NewRabbitMQConfig(url string) (*RabbitMQConfig, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQConfig{
		Connection: conn,
		Channel:    channel,
	}, nil
}

// Close closes the RabbitMQ connection and channel
func (r *RabbitMQConfig) Close() {
	if r.Channel != nil {
		r.Channel.Close()
	}
	if r.Connection != nil {
		r.Connection.Close()
	}
}

// InitRabbitMQ initializes the RabbitMQ configuration and returns any error encountered
func InitRabbitMQ() (*RabbitMQConfig, error) {
	// Fetch RabbitMQ connection details from environment variables
	username := os.Getenv("RABBITMQ_USER")
	password := os.Getenv("RABBITMQ_PASSWORD")
	host := os.Getenv("RABBITMQ_HOST")
	port := os.Getenv("RABBITMQ_PORT")

	// Construct the connection URL
	url := "amqp://" + username + ":" + password + "@" + host + ":" + port + "/"
	rabbitMQ, err := NewRabbitMQConfig(url)
	if err != nil {
		return nil, err // Return error instead of logging and exiting
	}
	log.Println("RabbitMQ connection established successfully!")
	return rabbitMQ, nil
}

// DeclareQueue declares a queue with the specified properties and returns any error encountered
func (r *RabbitMQConfig) DeclareQueue(queueName string, durable bool) error {
	_, err := r.Channel.QueueDeclare(
		queueName,
		durable,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err // Return error instead of logging and exiting
	}
	log.Printf("Queue '%s' declared successfully.", queueName)
	return nil
}
