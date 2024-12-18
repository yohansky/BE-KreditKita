package models

import (
	"time"

	"gorm.io/gorm"
)

type Consumer struct {
	Id           uint      `json:"id" gorm:"primaryKey"`
	NIK          uint      `json:"nik"`
	FullName     string    `json:"full_name"`
	LegalName    string    `json:"legal_name"`
	PlaceOfBirth string    `json:"place_of_birth"`
	DateOfBirth  string    `json:"date_of_birth"`
	Salary       float64   `json:"salary"`
	PhotoKTP     string    `json:"photo_ktp"`
	PhotoSelfie  string    `json:"photo_selfie"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
}

func (consumer *Consumer) Count(db *gorm.DB) int64 {
	var total int64

	db.Model(&Consumer{}).Count(&total)

	return total
}

func (consumer *Consumer) Take(db *gorm.DB, limit int, offset int) interface{} {
	var consumers []Consumer

	db.Offset(offset).Limit(limit).Find(&consumers)

	return consumers
}
