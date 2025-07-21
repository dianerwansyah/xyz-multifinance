package repository

import (
	"xyz-multifinance/internal/model"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(tx *gorm.DB, transaction *model.Transaction) error
	Update(tx *gorm.DB, id uint, transaction *model.Transaction) error
	Delete(id uint) error
	FindByID(id uint) (*model.Transaction, error)
	FindByCustomerID(customerID uint) ([]model.Transaction, error)
	FindAll() ([]model.Transaction, error)
	SumUsedAmount(customerID uint, tenor int) (int64, error)
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(tx *gorm.DB, transaction *model.Transaction) error {
	if err := tx.Create(transaction).Error; err != nil {
		return err
	}
	return nil
}

func (r *transactionRepository) Update(tx *gorm.DB, id uint, transaction *model.Transaction) error {
	if err := tx.First(&model.Transaction{}, id).Error; err != nil {
		return err
	}
	transaction.ID = id
	if err := tx.Save(transaction).Error; err != nil {
		return err
	}
	return nil
}

func (r *transactionRepository) Delete(id uint) error {
	if err := r.db.Delete(&model.Transaction{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *transactionRepository) FindByID(id uint) (*model.Transaction, error) {
	var transaction model.Transaction
	if err := r.db.First(&transaction, id).Error; err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *transactionRepository) FindByCustomerID(customerID uint) ([]model.Transaction, error) {
	var transactions []model.Transaction
	if err := r.db.Where("customer_id = ?", customerID).Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *transactionRepository) FindAll() ([]model.Transaction, error) {
	var transactions []model.Transaction
	if err := r.db.Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *transactionRepository) SumUsedAmount(customerID uint, tenor int) (int64, error) {
	var total int64
	err := r.db.Model(&model.Transaction{}).
		Where("customer_id = ? AND tenor = ? AND status IN ?", customerID, tenor, []string{"success", "ongoing"}).
		Select("COALESCE(SUM(installment_amount), 0)").
		Scan(&total).Error

	if err != nil {
		return 0, err
	}
	return total, nil
}
