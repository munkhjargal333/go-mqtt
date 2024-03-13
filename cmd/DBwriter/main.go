package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mqtt/database"
	"mqtt/helper"
	"mqtt/models"
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

	consumeMessages(ctx, channel)
	helper.WaitForShutdown()
}

func consumeMessages(ctx context.Context, channel *amqp.Channel) {
	msgs, err := channel.Consume(
		"ict.source",
		"db",
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

var batch0200 []models.Content0200 // Accumulate messages for bulk insert
var batch0202 []models.Content0202
var batch0230 []models.Content0230
var batch0800 []models.Content0800

var batchSize = 1

func processMessage(msg amqp.Delivery) {
	newMessage, _ := helper.RemoveState(string(msg.Body))

	var message Message
	err := json.Unmarshal([]byte(newMessage), &message)
	if err != nil {
		log.Println("Error decoding message:", err)
		return
	}

	switch content := message.Content.(type) {
	case string:
		// Content is already a string, continue processing
	case map[string]interface{}:
		// Convert map to JSON string
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
	/*
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

			data := models.Content0200{
				Alarm:     helper.ArrayToString(content0200.Alarm),
				Dir:       content0200.Dir,
				GPSTime:   content0200.GPSTime,
				Latitude:  content0200.Latitude,
				Longitude: content0200.Longitude,
				Mileage:   content0200.Mileage,
				OnlineNo:  content0200.OnlineNo,
				RouteCode: content0200.RouteCode,
				Speed:     content0200.Speed,
				StaIndex:  content0200.StaIndex,
				State:     helper.MapToString(content0200.State),
				UpDown:    content0200.UpDown,
			}
			// Add data to batchMessages
			batch0200 = append(batch0200, data)
	*/

	case "0230":
		var content0230 models.Content0230

		contentStr, ok := message.Content.(string)
		if !ok {
			log.Println("Error: Content is not a string")
			return
		}

		err := json.Unmarshal([]byte(contentStr), &content0230)
		if err != nil {
			log.Println("Error decoding content for ID 0230:", err)
			return
		}

		batch0230 = append(batch0230, content0230)

	case "0202":
		var Content0202 models.Content0202

		contentStr, ok := message.Content.(string)
		if !ok {
			log.Println("Error: Content is not a string")
			return
		}

		err := json.Unmarshal([]byte(contentStr), &Content0202)
		if err != nil {
			log.Println("Error decoding content for ID 0202:", err)
		}
		batch0202 = append(batch0202, Content0202)

	case "0800":
		var content0800 models.Content0800

		contentStr, ok := message.Content.(string)
		if !ok {
			log.Println("Error: Content is not a string")
			return
		}

		err := json.Unmarshal([]byte(contentStr), &content0800)
		if err != nil {
			log.Println("error decodin content for ID 0800:", err)
		}
		batch0800 = append((batch0800), content0800)

	default:
		log.Println("Unknown message ID:", message.ID)
		return
	}

	if len(batch0230) >= batchSize {
		if err := bulk0230(batch0230); err != nil {
			log.Printf("Error performing bulk insert: 0230 %v\n", err)
		}
		batch0230 = nil
	}

	if len(batch0200) >= batchSize {
		if err := bulk0200(batch0200); err != nil {
			log.Println("Error performing bulk insert: 0200 %v\n", err)
		}
		batch0200 = nil
	}

	if len(batch0202) >= batchSize {
		if err := bulk0202(batch0202); err != nil {
			log.Println("Error performing bulk insert: 0202 %v\n", err)
		}
		batch0202 = nil
	}

	if len(batch0800) >= batchSize {
		if err := bulk0800(batch0800); err != nil {
			log.Println("Error performing bulk insert: 0800 %v\n", err)
		}
		batch0800 = nil
	}
}

func bulk0200(messages []models.Content0200) error {
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(messages).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func bulk0230(messages []models.Content0230) error {
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(messages).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func bulk0202(messages []models.Content0202) error {
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(messages).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func bulk0800(messages []models.Content0800) error {
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(messages).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func waitForShutdown() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	log.Println("Shutting down...")
	os.Exit(0)
}
