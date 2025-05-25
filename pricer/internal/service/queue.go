package service

import (
	"fmt"
	"log"

	"github.com/chiyonn/vendiq2/pricer/internal/model"
	"github.com/go-resty/resty/v2"
	"github.com/streadway/amqp"
)

type QueueService struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

func NewQueueService() *QueueService {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to create RabbitMQ channel: %v", err)
	}

	q, err := ch.QueueDeclare(
		"price_jobs", // queue name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Fatalf("failed to declare queue: %v", err)
	}

	return &QueueService{
		conn:    conn,
		channel: ch,
		queue:   q,
	}
}

func (h *QueueService) AddQueue() error {
	body := `{"type":"adjust_prices"}`

	err := h.channel.Publish(
		"",             // exchange
		h.queue.Name,   // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		},
	)
	if err != nil {
		log.Printf("failed to publish message: %v", err)
		return err
	}

	log.Println("message published: adjust_prices")
	return nil
}

func (h *QueueService) GetAllQueues() ([]model.QueueInfo, error) {
	client := resty.New()

	resp, err := client.R().
		SetBasicAuth("guest", "guest").
		SetHeader("Accept", "application/json").
		SetResult(&[]model.QueueInfo{}).
		Get("http://rabbitmq:15672/api/queues")

	if err != nil {
		return nil, fmt.Errorf("failed to connect to management API: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("management API returned error: %s", resp.Status())
	}

	queues := *resp.Result().(*[]model.QueueInfo)
	log.Printf("successfully retrieved queue info (%d queues)", len(queues))

	return queues, nil
}
