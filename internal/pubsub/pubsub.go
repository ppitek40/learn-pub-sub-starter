package pubsub

import (
		amqp "github.com/rabbitmq/amqp091-go"
		"encoding/json"
		"encoding/gob"
		"context"
		"fmt"
		"bytes"
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	bytes, err := json.Marshal(val)
	if err != nil {
		return err
	}
	err = ch.PublishWithContext(context.Background(), exchange, key, false, false, amqp.Publishing{ ContentType: "application/json", Body: bytes })
	return err
}

func PublishGob[T any](ch *amqp.Channel, exchange, key string, val T) error {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	if err := enc.Encode(val); err != nil {
		return err
	}

	err := ch.PublishWithContext(context.Background(), exchange, key, false, false, amqp.Publishing{ ContentType: "application/gob", Body: buffer.Bytes()})
	return err
}

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error){
	channel, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, err
	}
	fmt.Println("Channel Created")
	fmt.Println(queueType == Transient)
	queue, err := channel.QueueDeclare(queueName, queueType == Durable, queueType == Transient, queueType == Transient, false, amqp.Table{ "x-dead-letter-exchange": "peril_dlx"})
	if err != nil {
		return nil, amqp.Queue{}, err
	}
	fmt.Println("Queue Created")
	err = channel.QueueBind(queueName, key, exchange, false, nil)
	if err != nil {
		return nil, amqp.Queue{}, err
	}
	fmt.Println("Queue Binded")
	return channel, queue, nil
}

func SubscribeJSON[T any](
    conn *amqp.Connection,
    exchange,
    queueName,
    key string,
    queueType SimpleQueueType, // an enum to represent "durable" or "transient"
    handler func(T) AckType,
) error {
	channel, _, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		return err
	}

	chDelivery, err := channel.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func () {
		for delivery := range chDelivery {
			var value T
			json.Unmarshal(delivery.Body, &value)
			ack := handler(value)
			switch ack {
			case Ack:
				delivery.Ack(false)
			case NackDiscard:
				delivery.Nack(false, false)
			case NackRequeue:
				delivery.Nack(false, true)
			}
		}
	}()
	return nil
}

func SubscribeGob[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
	handler func(T) AckType,
) error {
	channel, _, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		return err
	}

	if err := channel.Qos(10, 0, false); err != nil {
		return err
	}

	chDelivery, err := channel.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	unmarshaller := func(data []byte) (T, error) {
		var target T
		dec := gob.NewDecoder(bytes.NewBuffer(data))
		err := dec.Decode(&target)
		return target, err
	}

	go func() {
		defer channel.Close()
		for delivery := range chDelivery {
			value, err := unmarshaller(delivery.Body)
			if err != nil {
				fmt.Println(err)
				continue
			}
			switch ack := handler(value); ack {
			case Ack:
				delivery.Ack(false)
			case NackRequeue:
				delivery.Nack(false, true)
			case NackDiscard:
				delivery.Nack(false, false)
			default:
				fmt.Println("unknown ack")
			}
		}
	}()

	return nil
}

type AckType int
const (
	Ack AckType = iota
	NackRequeue
	NackDiscard
)

type SimpleQueueType int
const (
	Durable SimpleQueueType = iota
	Transient
)