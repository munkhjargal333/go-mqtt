package mqtt

import (
	"fmt"
	"os"

	"github.com/streadway/amqp"
)

func InitializeRabbitMQ() (*amqp.Channel, error) {
	// Initialize RabbitMQ connection
	mqHost := os.Getenv("mq_host")
	mqPort := os.Getenv("mq_port")
	mqUser := os.Getenv("mq_user")
	mqPassword := os.Getenv("mq_password")

	amqpURI := fmt.Sprintf("amqp://%s:%s@%s:%s/", mqUser, mqPassword, mqHost, mqPort)

	// Create RabbitMQ channel
	channel, err := CreateChannel(amqpURI)
	if err != nil {
		return nil, err
	}

	queueName := "ict.source"
	exchangeName := "ict_edit_source"
	routingKey := "ict_edited_source"

	// Create exchange if it doesn't exist
	err = CreateExchange(channel, exchangeName, "direct")
	if err != nil {
		return nil, err
	}

	// Create queue if it doesn't exist
	_, err = CreateQueue(channel, queueName)
	if err != nil {
		return nil, err
	}

	// Bind queue to exchange with routing key
	err = channel.QueueBind(
		queueName,
		routingKey,
		exchangeName,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return channel, nil
}

func CreateQueue(ch *amqp.Channel, queueName string) (amqp.Queue, error) {
	que, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return amqp.Queue{}, err
	}

	return que, nil
}

func CreateExchange(ch *amqp.Channel, exchangeName string, exchangeType string) error {
	err := ch.ExchangeDeclare(
		exchangeName,
		exchangeType,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

func CreateChannel(amqpURI string) (*amqp.Channel, error) {
	conn, err := amqp.Dial(amqpURI)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return ch, nil
}
