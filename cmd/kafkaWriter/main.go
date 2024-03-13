package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mqtt/content"
	"mqtt/helper"
	"mqtt/mqtt"

	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
	"github.com/streadway/amqp"
	"gorm.io/gorm"
)

type Message struct {
	ID      string      `json:"id"`
	Content interface{} `json:"content"`
}

var DB *gorm.DB

const (
	kafkaBrokers = "43.231.115.108:9092"
	topic        = "mq-gps"
)

var kafkaWriter *kafka.Writer

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	channel, err := mqtt.InitializeRabbitMQ()
	if err != nil {
		log.Fatal("Error initializing RabbitMQ:", err)
	}
	defer channel.Close()

	if err := initKafkaWriter(); err != nil {
		log.Fatal("Error initializing Kafka writer:", err)
	}
	defer kafkaWriter.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	consumeMessages(ctx, channel)
	helper.WaitForShutdown()
}

func consumeMessages(ctx context.Context, channel *amqp.Channel) {
	msgs, err := channel.Consume(
		"ict.source",
		"kafka",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("Error consuming messages from queue:", err)
	}

	log.Println("Waiting for messages...")

	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down consumer...")
			return
		case msg := <-msgs:
			processMessage(msg)
		}
	}
}

func processMessage(msg amqp.Delivery) {
	newMessage, stateData := helper.RemoveState(string(msg.Body))

	var message Message
	err := json.Unmarshal([]byte(newMessage), &message)
	if err != nil {
		log.Println("Error decoding message:", err)
		return
	}

	switch content := message.Content.(type) {
	case string:

	case map[string]interface{}:
		contentBytes, err := json.Marshal(content)
		if err != nil {
			log.Println("Error converting content to string:", err)
			return
		}
		message.Content = string(contentBytes)
	default:
		log.Println("Unsupported content type")
		return
	}

	switch message.ID {

	case "0200":
		var content0200 content.Content0200
		err := json.Unmarshal([]byte(message.Content.(string)), &content0200)
		if err != nil {
			log.Println("Error decoding content for ID 0200:", err)
			return
		}
		content0200.State, err = helper.ParseStateData(stateData)
		if err != nil {
			log.Println("Error parsing state data:", err)
			return
		}

		contentJSON, err := json.Marshal(content0200)
		if err != nil {
			log.Println("Error encoding content0200 to JSON:", err)
			return
		}

		err = kafkaWriter.WriteMessages(context.Background(), kafka.Message{

			Key:   []byte("0200"), // Assuming ID is unique and can be used as the key
			Value: contentJSON,
		})
		if err != nil {
			log.Println("Error writing content0200 to Kafka:", err)
			return
		}

		log.Println("Written message with ID 0200 to Kafka")

	default:
		log.Println("Unknown message ID:", message.ID)
	}
}

func initKafkaWriter() error {
	kafkaWriter = kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{kafkaBrokers},
		Topic:   "gps-log1",
	})

	return nil
}
