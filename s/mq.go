package util

import (
	"github.com/streadway/amqp"
)

func Connect() (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, err
	}

	return conn, ch, nil
}

func PublishMessage(ch *amqp.Channel, exchange, queueName, body string) error {
	err := ch.Publish(
		exchange,
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	return err
}
