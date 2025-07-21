package repository

import (
	"xyz-multifinance/internal/model"

	"gorm.io/gorm"
)

type CustomerRepository interface {
	FindByNIK(nik string) (*model.Customer, error)
	FindByID(id uint) (*model.Customer, error)
	Create(customer *model.Customer) error
	Update(nik string, fields map[string]interface{}) error
	Delete(id uint) error
}

type customerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &customerRepository{db: db}
}

func (r *customerRepository) FindByNIK(nik string) (*model.Customer, error) {
	var customer model.Customer
	if err := r.db.Where("nik = ?", nik).First(&customer).Error; err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *customerRepository) FindByID(id uint) (*model.Customer, error) {
	var customer model.Customer
	if err := r.db.First(&customer, id).Error; err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *customerRepository) Create(customer *model.Customer) error {
	return r.db.Create(customer).Error
}

func (r *customerRepository) Update(nik string, fields map[string]interface{}) error {
	return r.db.Model(&model.Customer{}).Where("nik = ?", nik).Updates(fields).Error
}
func (r *customerRepository) Delete(id uint) error {
	return r.db.Delete(&model.Customer{}, id).Error
}
