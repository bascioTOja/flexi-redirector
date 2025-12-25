package models

import "time"

type ShortURL struct {
	ID          uint      `gorm:"primaryKey"`
	CreatedByID uint      `gorm:"index;not null"`
	Name        string    `gorm:"size:255;not null"`
	Slug        string    `gorm:"size:64;uniqueIndex;not null"`
	LongURL     string    `gorm:"size:2048;not null"`
	Views       uint      `gorm:"default:0"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

func (ShortURL) TableName() string {
	return "urlshortener_shorturl"
}
