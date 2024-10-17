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
		return nil, err
	}
	log.Println("RabbitMQ connection established successfully!")
	return rabbitMQ, nil
}

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
		return err
	}
	log.Printf("Queue '%s' declared successfully.", queueName)
	return nil
}

func (r *RabbitMQConfig) PublishMessage(queueName string, message []byte) error {
	err := r.Channel.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		})
	if err != nil {
		log.Printf("Failed to publish message: %s", err)
		return err
	}

	log.Printf("Message published to queue %s: %s", queueName, message)
	return nil
}
