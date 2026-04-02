package pubsub

import (
	"encoding/json"

	"github.com/rabbitmq/amqp091-go"
)

func SubscribeJSON[T any](
	conn *amqp091.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
	handler func(T),
) error {
	chann, _, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		return err
	}

	delivery, err := chann.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() error {
		for msg := range delivery {
			var unmarshelMsg T
			err := json.Unmarshal(msg.Body, &unmarshelMsg)
			if err != nil {
				return err
			}
			handler(unmarshelMsg)
			msg.Ack(false)
		}
		return nil
	}()

	return nil
}
