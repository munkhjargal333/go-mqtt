package main

import (
	"context"
	"fmt"
	"log"
	"mqtt/database"
	"mqtt/mqtt"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
	"gorm.io/gorm"
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

	database.MustConnect()

	channel, err := mqtt.InitializeRabbitMQ()
	if err != nil {
		log.Fatal("Error initializing RabbitMQ:", err)
	}
	defer channel.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		cancel()
	}()

	consumeMessages(ctx, channel)
}

func consumeMessages(ctx context.Context, channel *amqp.Channel) {
	msgs, err := channel.Consume(
		"ict.source",
		"for gps",
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
	log.Println(string(msg.Body))
}
