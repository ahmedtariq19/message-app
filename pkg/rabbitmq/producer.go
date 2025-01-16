package rabbitmq

import (
	"fmt"
	"log"
	"message-app/conf"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer interface {
	Publish(queue string, message string) error
}
type producer struct {
	conn *amqp.Connection
}

func NewProducer(conf conf.RabbitMQ) (Producer, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://guest:guest@%s:%s/", conf.Host, conf.Port))
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	return &producer{conn: conn}, nil
}

func (p *producer) Publish(queue string, message string) error {
	ch, err := p.conn.Channel()
	if err != nil {
		return err
	}
	_, err = ch.QueueDeclare(queue, false, false, false, false, nil)
	if err != nil {
		return err
	}
	return ch.Publish("", queue, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(message),
	})
}
