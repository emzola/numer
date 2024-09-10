package rabbitmq

import (
	"context"
	"encoding/json"
	"log"

	"github.com/emzola/numer/activity-service/internal/models"
	"github.com/emzola/numer/activity-service/internal/repository"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
	repo    *repository.ActivityRepository
}

// Initialize RabbitMQ connection
func NewConsumer(url, queueName string, repo *repository.ActivityRepository) (*Consumer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &Consumer{
		conn:    conn,
		channel: ch,
		queue:   queueName,
		repo:    repo,
	}, nil
}

// Start consuming messages from the queue
func (c *Consumer) ConsumeMessages() {
	msgs, err := c.channel.Consume(
		c.queue, // queue
		"",      // consumer
		true,    // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
	if err != nil {
		log.Fatalf("failed to register consumer: %v", err)
	}

	go func() {
		for msg := range msgs {
			log.Printf("received a message: %s", msg.Body)

			// Deserialize the JSON to a Go map
			var message map[string]interface{}
			err := json.Unmarshal(msg.Body, &message)
			if err != nil {
				log.Fatalf("failed to deserialize JSON: %s", err)
			}

			// Extract data from the message map
			invoiceID := int64(message["invoice_id"].(float64)) // Convert float64 to int64
			userID := int64(message["user_id"].(float64))       // Convert float64 to int64
			action := message["action"].(string)
			description := message["description"].(string)

			activity := &models.Activity{
				InvoiceID:   invoiceID,
				UserID:      userID,
				Action:      action,
				Description: description,
			}

			// Store activity log in the database
			c.repo.LogActivity(context.TODO(), activity)
		}
	}()
}

func (c *Consumer) Close() {
	c.channel.Close()
	c.conn.Close()
}
