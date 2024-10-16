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

// InitRabbitMQ initializes the RabbitMQ configuration
func InitRabbitMQ() *RabbitMQConfig {
	// Fetch RabbitMQ connection details from environment variables
	username := os.Getenv("RABBITMQ_USER")
	password := os.Getenv("RABBITMQ_PASSWORD")
	host := os.Getenv("RABBITMQ_HOST")
	port := os.Getenv("RABBITMQ_PORT")

	// Construct the connection URL
	url := "amqp://" + username + ":" + password + "@" + host + ":" + port + "/"
	rabbitMQConfig, err := NewRabbitMQConfig(url)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	log.Println("RabbitMQ connection established successfully!")
	return rabbitMQConfig
}

func (r *RabbitMQConfig) DeclareQueue(queueName string, durable bool) {
	_, err := r.Channel.QueueDeclare(
		queueName,
		durable,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
	}
}
