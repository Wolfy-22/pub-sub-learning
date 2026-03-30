package pubsub

import "github.com/rabbitmq/amqp091-go"

type SimpleQueueType struct {
	Durable   bool
	Transient bool
}

func DeclareAndBind(
	conn *amqp091.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
) (*amqp091.Channel, amqp091.Queue, error) {
	chann, err := conn.Channel()
	if err != nil {
		return &amqp091.Channel{}, amqp091.Queue{}, err
	}

	queue, err := chann.QueueDeclare(queueName, queueType.Durable, queueType.Transient, queueType.Transient, false, nil)
	if err != nil {
		return &amqp091.Channel{}, amqp091.Queue{}, err
	}

	chann.QueueBind(queueName, key, exchange, false, nil)

	return chann, queue, nil
}
