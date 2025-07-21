package model

import (
	"time"

	"gorm.io/gorm"
)

type Limit struct {
	ID         uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	CustomerID uint           `gorm:"index;not null" json:"customer_id"`
	Tenor      int            `gorm:"column:tenor_month;not null" json:"tenor_month"`
	Limit      int64          `gorm:"column:limit_amount;not null" json:"limit_amount"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
