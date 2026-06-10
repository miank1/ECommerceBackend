package rabbitmq

import (
	"log"
)

func (r *RabbitMQ) Consume(queueName string) error {

	msgs, err := r.Channel.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto ack
		false,     // exclusive
		false,     // no local
		false,     // no wait
		nil,       // args
	)

	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {

			log.Printf(
				"📦 Message Received: %s",
				string(msg.Body),
			)

			msg.Ack(false)
		}
	}()

	return nil
}
