package models

import (
	"fmt"
	"log"
	"mqtt/database"
	"time"

	"gorm.io/gorm"
)

type Content0200 struct {
	ID        uint64  `json:"id" gorm:"primaryKey;autoIncrement"`
	Alarm     string  `json:"alarm"`
	Dir       int     `json:"dir"`
	GPSTime   string  `json:"gpsTime"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Mileage   int     `json:"mileage"`
	OnlineNo  string  `json:"onlineNo"`
	RouteCode string  `json:"routeCode"`
	Speed     int     `json:"speed"`
	StaIndex  int     `json:"staIndex"`
	State     string  `json:"state"`
	UpDown    int     `json:"upDown"`

	ObjectVersion uint           `json:"-" gorm:"default:1"`
	CreatedByID   uint           `json:"-"`
	CreatedAt     *time.Time     `json:"created_at,omitempty" gorm:"autoCreateTime"`
	UpdatedByID   uint           `json:"-"`
	UpdatedAt     *time.Time     `json:"-" gorm:"autoUpdateTime"`
	DeletedByID   uint           `json:"-" gorm:"column:deleted"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

func (u *Content0200) TableName() string {
	return "app_backend.mq_gps"
}

func (u *Content0200) Migrate() error {
	db := database.DB
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	if err := db.AutoMigrate(&Content0230{}); err != nil {
		log.Println("Error migrating Content0230:", err)
		return err
	}
	return nil
}

func (u *Content0200) Drop() error {
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
