package database

import (
	"fmt"
	"log"
	"os"

	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func MustConnect() {
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
