package models

import (
	"fmt"
	"log"
	"mqtt/database"
	"time"

	"gorm.io/gorm"
)

type Content0230 struct {
	ID           uint64  `json:"id" gorm:"primaryKey;autoIncrement"`
	OnlineNo     string  `json:"onlineNo"`
	RouteCode    string  `json:"routeCode"`
	UpDown       int     `json:"upDown"`
	StationIndex int     `json:"stationIndex"`
	GetOn        int     `json:"getOn"`
	GetOff       int     `json:"getOff"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	Dir          int     `json:"dir"`
	GPSTime      string  `json:"gpsTime"`
	Speed        int     `json:"speed"`
	Type         int     `json:"type"`

	ObjectVersion uint           `json:"-" gorm:"default:1"`
	CreatedByID   uint           `json:"-"`
	CreatedAt     *time.Time     `json:"created_at,omitempty" gorm:"autoCreateTime"`
	UpdatedByID   uint           `json:"-"`
	UpdatedAt     *time.Time     `json:"-" gorm:"autoUpdateTime"`
	DeletedByID   uint           `json:"-" gorm:"column:deleted"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

func (u *Content0230) TableName() string {
	return "app_backend.bs_passenger_counter_details"
}

func (u *Content0230) Migrate() error {
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

func (u *Content0230) Drop() error {
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
