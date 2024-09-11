package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
}

// Initialize RabbitMQ connection
func NewPublisher(url, queueName string) (*Publisher, error) {
	var conn *amqp.Connection
	var err error

	for i := 0; i < 5; i++ { // Retry 5 times
		conn, err = amqp.Dial(url)
		if err == nil {
			break
		}
		log.Printf("failed to connect to RabbitMQ, retrying in 5 seconds... (%d/5)", i+1)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &Publisher{
		conn:    conn,
		channel: ch,
		queue:   queueName,
	}, nil
}

// Publish an event to RabbitMQ
func (p *Publisher) Publish(message any) error {
	// Serialize the message to JSON
	body, err := json.Marshal(message)
	if err != nil {
		log.Fatal("failed to serialize map to json")
	}

	err = p.channel.Publish(
		"",      // exchange
		p.queue, // routing key (queue name)
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		log.Fatal("failed to publish a message")
	}
	log.Printf("Published message: %s", message)
	return nil
}

// Close the connection when done
func (p *Publisher) Close() {
	p.channel.Close()
	p.conn.Close()
}
