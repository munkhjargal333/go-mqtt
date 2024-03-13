package models

import (
	"fmt"
	"log"
	"mqtt/database"
	"time"

	"gorm.io/gorm"
)

type Content0800 struct {
	ID        uint64  `json:"id" gorm:"primaryKey;autoIncrement"`
	OnlineNo  string  `json:"onlineNo"`
	RouteCode string  `json:"routeCode"`
	DriverId  string  `json:"driverId"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	GPSTime   string  `json:"gpsTime"`
	Type      int     `json:"type"`

	ObjectVersion uint           `json:"-" gorm:"default:1"`
	CreatedByID   uint           `json:"-"`
	CreatedAt     *time.Time     `json:"created_at,omitempty" gorm:"autoCreateTime"`
	UpdatedByID   uint           `json:"-"`
	UpdatedAt     *time.Time     `json:"-" gorm:"autoUpdateTime"`
	DeletedByID   uint           `json:"-" gorm:"column:deleted"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

func (u *Content0800) TableName() string {
	return "app_backend.bs_driver_attendance"
}

func (u *Content0800) Migrate() error {
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

func (u *Content0800) Drop() error {
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
