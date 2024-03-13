package models

import (
	"fmt"
	"log"
	"mqtt/database"
	"time"

	"gorm.io/gorm"
)

type Content0202 struct {
	ID uint64 `json:"id" gorm:"primaryKey;autoIncrement"`

	OnlineNo  string  `json:"onlineNo"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	Speed     int     `json:"speed"`
	Dir       int     `json:"dir"`
	GPSTime   string  `json:"gpsTime"`

	CanItem string `json:"canItem"`

	ObjectVersion uint           `json:"-" gorm:"default:1"`
	CreatedByID   uint           `json:"-"`
	CreatedAt     *time.Time     `json:"created_at,omitempty" gorm:"autoCreateTime"`
	UpdatedByID   uint           `json:"-"`
	UpdatedAt     *time.Time     `json:"-" gorm:"autoUpdateTime"`
	DeletedByID   uint           `json:"-" gorm:"column:deleted"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

func (u *Content0202) TableName() string {
	return "app_backend.bs_can_details"
}

func (u *Content0202) Migrate() error {
	db := database.DB
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	if err := db.AutoMigrate(&Content0202{}); err != nil {
		log.Println("Error migrating Content0202:", err)
		return err
	}
	return nil
}

func (u *Content0202) Drop() error {
	db := database.DB
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	if err := db.Exec("DROP TABLE " + u.TableName()).Error; err != nil {
		log.Println("Error dropping table:", err)
		return err
	}

	return nil
}
