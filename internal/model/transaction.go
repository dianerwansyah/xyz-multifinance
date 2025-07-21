package model

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	ID                uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	ContractNumber    string         `gorm:"uniqueIndex;not null" json:"contract_number"`
	CustomerID        uint           `gorm:"index;not null" json:"customer_id"`
	Tenor             int            `gorm:"not null" json:"tenor"`
	InstallmentAmount int64          `gorm:"column:installment_amount;not null" json:"installment_amount"`
	OTR               int64          `json:"otr"`
	AdminFee          int64          `gorm:"column:admin_fee" json:"admin_fee"`
	InterestAmount    int64          `gorm:"column:interest_amount" json:"interest_amount"`
	AssetName         string         `json:"asset_name"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}
