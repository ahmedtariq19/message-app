package rabbitmq

import (
	"fmt"
	"log"
	"message-app/conf"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer interface {
	Consume(queueName string, msgCh chan string, done chan struct{}) error
}

type consumer struct {
	conn *amqp.Connection
}

func NewConsumer(conf conf.RabbitMQ) (Consumer, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://guest:guest@%s:%s/", conf.Host, conf.Port))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	return &consumer{conn: conn}, nil
}

func (c *consumer) Consume(queueName string, msgCh chan string, done chan struct{}) error {

	ch, err := c.conn.Channel()
	if err != nil {
		return nil
	}

	defer ch.Close()

	// Declare the queue if it doesn't exist
	_, err = ch.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	// Start consuming messages asynchronously
	msgs, err := ch.Consume(queueName, "", true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to start consuming messages: %w", err)
	}

	// Use a select statement to non-blocking consume messages
	for {
		select {
		case <-done:
			log.Println("Stopping writeToSocket due to connection close.")
			return nil
		case msg := <-msgs:
			msgCh <- string(msg.Body)
		default:
			continue
		}
	}
}
