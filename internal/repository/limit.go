package repository

import (
	"xyz-multifinance/internal/model"

	"gorm.io/gorm"
)

type LimitRepository interface {
	Create(limit *model.Limit) error
	Update(id uint, fields map[string]interface{}) error
	Delete(id uint) error
	FindByID(id uint) (*model.Limit, error)
	FindByCustomerID(customerID uint) ([]model.Limit, error)
}

type limitRepository struct {
	db *gorm.DB
}

func NewLimitRepository(db *gorm.DB) LimitRepository {
	return &limitRepository{db: db}
}

func (r *limitRepository) Create(limit *model.Limit) error {
	return r.db.Create(limit).Error
}

func (r *limitRepository) Update(id uint, fields map[string]interface{}) error {
	return r.db.Model(&model.Limit{}).Where("id = ?", id).Updates(fields).Error
}

func (r *limitRepository) Delete(id uint) error {
	return r.db.Delete(&model.Limit{}, id).Error
}

func (r *limitRepository) FindByID(id uint) (*model.Limit, error) {
	var limit model.Limit
	err := r.db.First(&limit, id).Error
	return &limit, err
}

func (r *limitRepository) FindByCustomerID(customerID uint) ([]model.Limit, error) {
	var limits []model.Limit
	err := r.db.Where("customer_id = ?", customerID).Find(&limits).Error
	return limits, err
}
