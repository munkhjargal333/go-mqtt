package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mqtt/content"
	"mqtt/helper"
	"mqtt/models"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Message struct {
	ID      string      `json:"id"`
	Content interface{} `json:"content"`
}

var DB *gorm.DB

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	initialDB()

	channel, err := InitializeRabbitMQ()
	if err != nil {
		log.Fatal("Error initializing RabbitMQ:", err)
	}
	defer channel.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	consumeMessages(ctx, channel)
	waitForShutdown()
}

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

	queueName := "ict.business"
	exchangeName := "ict_edit_trip"
	routingKey := "ict_edited_trip"

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

func initialDB() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	dbHost := os.Getenv("db_host")
	dbPort := os.Getenv("db_port")
	dbUser := os.Getenv("db_user")
	dbPassword := os.Getenv("db_password")
	dbName := os.Getenv("db_name")

	dsn := fmt.Sprintf(
		`host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Ulaanbaatar`,
		dbHost,
		dbUser,
		dbPassword,
		dbName,
		dbPort,
	)

	db, err := gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
			NowFunc: func() time.Time {
				return time.Now()
			},
		},
	)

	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
	}

	DB = db
}

func consumeMessages(ctx context.Context, channel *amqp.Channel) {
	msgs, err := channel.Consume(
		"ict.business",
		"trip",
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
	if msg.Body == nil {
		log.Println("Message content is nil")
		return
	}
	//log.Println(string(msg.Body))
	var message Message
	err := json.Unmarshal(msg.Body, &message)
	if err != nil {
		log.Println("Error decoding message:", err)
		return
	}

	switch message.ID {
	case "0DE0":
		log.Println(message)
	default:
		log.Println("Unknown message ID:", message.ID)
	}
}

func processMessage2(msg amqp.Delivery) {
	if msg.Body == nil {
		log.Println("Message content is nil")
		return
	}

	var message Message
	err := json.Unmarshal(msg.Body, &message)
	if err != nil {
		log.Println("Error decoding message:", err)
		return
	}

	switch message.ID {
	case "0DE0":
		var content0DE0 content.Content0DE0
		err := json.Unmarshal([]byte(message.Content.(string)), &content0DE0)
		if err != nil {
			log.Println("Error decoding content for ID 0200:", err)
			return
		}
		data := models.Content0DE0{
			PlanDate: content0DE0.PlanDate,
			PlanList: helper.ArrayToString(content0DE0.PlanList),
			State:    content0DE0.State,
		}

		DB.Create(data)

		log.Println(message)
	default:
		log.Println("Unknown message ID:", message.ID)
	}
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

func waitForShutdown() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	log.Println("Shutting down...")
	os.Exit(0)
}
