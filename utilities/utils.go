package utilities

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
)

// createAMQPConnection establishes a connection to RabbitMQ.
func CreateAMQPConnection() (*amqp.Connection, error) {
	amqpURL := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, fmt.Errorf("Error occurred while establishing connection: %s", err)
	}
	return conn, nil
}

// createAMQPChannel creates a channel for the RabbitMQ connection.
func CreateAMQPChannel(conn *amqp.Connection) (*amqp.Channel, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("Error occurred while establishing channel: %s", err)
	}
	return channel, nil
}

// declareQueue declares a queue on the RabbitMQ channel.
func DeclareQueue(channel *amqp.Channel, queueName string) error {
	_, err := channel.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("Error occurred while declaring queue: %s", err)
	}
	return nil
}

// publishMessage publishes a message to a RabbitMQ queue.
func ConsumeMessage(ctx context.Context, channel *amqp.Channel, queueName string) {
	message, err := channel.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		log.Error().Msgf("Error occurred while publishing message: %s", err)
	}
	var forever chan struct{}
	var cveInfo []CVEInformation
	go func() {
		for d := range message {
			err = json.Unmarshal([]byte(d.Body), &cveInfo)
			if err != nil {
				log.Error().Msgf("Error unmarshalling data")
			}
			log.Info().Msgf("Received a message: %s", cveInfo)
			fmt.Println()
		}
	}()

	log.Info().Msg(" [*] Waiting for messages.")
	<-forever
}

// CreateRMQConn establishes a connection to RabbitMQ and publishes a message to a queue.
func SubscribeMessage(ctx context.Context) error {
	conn, err := CreateAMQPConnection()
	if err != nil {
		return err
	}
	defer conn.Close()
	channel, err := CreateAMQPChannel(conn)
	if err != nil {
		return err
	}
	defer channel.Close()
	if err := DeclareQueue(channel, QueueName); err != nil {
		return err
	}
	ConsumeMessage(ctx, channel, QueueName)

	log.Info().Msg(fmt.Sprintf("Successfully established connection and published message to queue '%s'", QueueName))
	return nil
}
