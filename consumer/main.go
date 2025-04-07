package main

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"homework-dontpanicw/consumer/domain"
	"homework-dontpanicw/consumer/image_processor"
	"log"
	"time"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {

	conn, err := amqp.Dial("amqp://guest:guest@broker:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	queue, err := ch.QueueDeclare(
		"photos",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var task domain.Task
			if err := json.Unmarshal(d.Body, &task); err != nil {
				log.Printf("Error decoding JSON: %v", err)
				continue
			}
			//чтобы файл успел сохраниться
			time.Sleep(2 * time.Second)
			err = image_processor.RedactPhoto(task)
			if err != nil {
				log.Printf("Error redacting image: %v", err)
			}

			//response := map[string]uuid.UUID{"PhotoId": task.PhotoId}
			//responseJson, _ := json.Marshal(response)
			response := task.PhotoId.String()

			err = ch.Publish("",
				d.ReplyTo,
				false,
				false,
				amqp.Publishing{
					ContentType:   "text/plain",
					Body:          []byte(response),
					CorrelationId: d.CorrelationId,
				},
			)
			if err != nil {
				log.Printf("Failed to publish a message: %s", err)
			}

			d.Ack(false)
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

}
