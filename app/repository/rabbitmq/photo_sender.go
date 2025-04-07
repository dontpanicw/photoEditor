package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"homework-dontpanicw/app/domain"
	"homework-dontpanicw/app/usecases"
	"log"
	"time"
)

type RabbitMQSender struct {
	connection    *amqp.Connection
	channel       *amqp.Channel
	queueName     string
	responseQueue string
	responseChan  <-chan amqp.Delivery
}

func NewRabbitMQSender(amqpURL, queueName string, responseQueue string) (*RabbitMQSender, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %s", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %s", err)
	}

	_, err = ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Errorf("failed to declare a queue: %s", err)
	}

	_, err = ch.QueueDeclare(
		responseQueue,
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		return nil, fmt.Errorf("failed to declare response queue: %s", err)
	}

	return &RabbitMQSender{
		connection:    conn,
		channel:       ch,
		queueName:     queueName,
		responseQueue: responseQueue,
	}, nil
}

func (r *RabbitMQSender) SendTask(task domain.Task) (string, error) {
	body, err := json.Marshal(task)
	if err != nil {
		return "", err
	}

	correlationId := fmt.Sprintf("%d", time.Now().UnixNano())

	err = r.channel.Publish(
		"",
		r.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType:   "application/json",
			Body:          body,
			CorrelationId: correlationId,
			ReplyTo:       r.responseQueue,
		})
	if err != nil {
		return "", err
	}

	return correlationId, nil
}

func (r *RabbitMQSender) ListenForResponses(usecase usecases.Task) {
	msgs, err := r.channel.Consume(
		r.responseQueue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %s", err)
	}
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			//логика с response
			photoId, err := uuid.Parse(string(d.Body))
			if err != nil {
				log.Fatalf("Ошибка при конвертации: %v", err)
			}
			err = usecase.DoingTask(context.Background(), photoId)
			if err != nil {
				log.Printf("Doing task error: %v: ", err)
			}
			log.Printf("Photo is ready!: %s", d.Body)

		}
	}()
}
