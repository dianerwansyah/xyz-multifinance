package model

import (
	"time"

	"gorm.io/gorm"
)

type Customer struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      uint           `gorm:"column:user_id" json:"user_id"`
	NIK         string         `gorm:"uniqueIndex;size:16;not null" json:"nik"`
	FullName    string         `gorm:"not null" json:"full_name"`
	LegalName   string         `json:"legal_name"`
	PlaceBirth  string         `gorm:"column:birth_place" json:"place_of_birth"`
	DateBirth   time.Time      `gorm:"column:birth_date" json:"date_of_birth"`
	Salary      int64          `json:"salary"`
	KTPPhoto    string         `gorm:"column:photo_ktp" json:"ktp_photo"`
	SelfiePhoto string         `gorm:"column:photo_selfie" json:"selfie_photo"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
