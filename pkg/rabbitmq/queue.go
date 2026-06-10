package rabbitmq

import amqp "github.com/rabbitmq/amqp091-go"

func (r *RabbitMQ) DeclareQueue(name string) (amqp.Queue, error) {

	return r.Channel.QueueDeclare(
		name,  // queue name
		true,  // durable
		false, // auto delete
		false, // exclusive
		false, // no wait
		nil,   // args
	)
}
