package rabbitmq

import (
	"log"
	"message-app/conf"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	rabbitMQHost = "127.0.0.1"
	rabbitMQPort = "4672"
	testQueue    = "test-queue"
	testMessage  = "Hello, RabbitMQ!"
)

func TestRabbitMQProducerConsumer(t *testing.T) {
	config := conf.RabbitMQ{
		Host: rabbitMQHost,
		Port: rabbitMQPort,
	}
	producer, err := NewProducer(config)
	if err != nil {
		t.Fatalf("Failed to create producer: %v", err)
	}
	consumer, err := NewConsumer(config)
	if err != nil {
		t.Fatalf("Failed to create consumer: %v", err)
	}

	msgCh := make(chan string)
	done := make(chan struct{})

	go func() {
		if err := consumer.Consume(testQueue, msgCh, done); err != nil {
			log.Fatalf("Failed to consume messages: %v", err)
		}
	}()

	time.Sleep(1 * time.Second)

	err = producer.Publish(testQueue, testMessage)
	assert.NoError(t, err, "Failed to publish message")

	select {
	case msg := <-msgCh:
		assert.Equal(t, testMessage, msg, "Received message does not match expected")
	case <-time.After(3 * time.Second):
		t.Fatal("Timed out waiting for message")
	}

	close(done)
}
