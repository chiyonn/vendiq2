package consumer

import (
	"context"
	"encoding/json"
	"log"

	"github.com/chiyonn/vendiq2/pricer/internal/bot"
	"github.com/streadway/amqp"
)

func StartConsumer(b bot.PricerBot) {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ (consumer): %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to create channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"price_jobs",
		true, false, false, false, nil,
	)
	if err != nil {
		log.Fatalf("failed to declare queue: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		false, false, false, false, nil,
	)
	if err != nil {
		log.Fatalf("failed to register consumer: %v", err)
	}

	log.Println("RabbitMQ consumer started.")

	for msg := range msgs {
		log.Printf("received message: %s", msg.Body)

		var task struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal(msg.Body, &task); err != nil {
			log.Printf("invalid message format: %v", err)
			msg.Nack(false, false)
			continue
		}

		if task.Type == "adjust_prices" {
			log.Println("executing price adjustment via PricerBot...")

			err := b.Run(context.Background())
			if err != nil {
				log.Printf("PricerBot failed: %v", err)
				msg.Nack(false, true)
				continue
			}

			log.Println("PricerBot completed price adjustment.")
		}

		msg.Ack(false)
	}
}
