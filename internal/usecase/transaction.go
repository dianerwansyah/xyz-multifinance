package usecase

import (
	"errors"
	"xyz-multifinance/internal/model"
	"xyz-multifinance/internal/repository"

	"gorm.io/gorm"
)

type TransactionUsecase interface {
	CreateTransaction(tx *model.Transaction) error
	UpdateTransaction(id uint, tx *model.Transaction) error
	DeleteTransaction(id uint) error
	GetTransactionByID(id uint) (*model.Transaction, error)
	GetTransactionsByCustomer(customerID uint) ([]model.Transaction, error)
	GetAllTransactions() ([]model.Transaction, error)
}

type transactionUsecase struct {
	txRepo       repository.TransactionRepository
	limitRepo    repository.LimitRepository
	customerRepo repository.CustomerRepository
	db           *gorm.DB
}

func NewTransactionUsecase(
	txRepo repository.TransactionRepository,
	limitRepo repository.LimitRepository,
	customerRepo repository.CustomerRepository,
	db *gorm.DB,
) TransactionUsecase {
	return &transactionUsecase{
		txRepo:       txRepo,
		limitRepo:    limitRepo,
		customerRepo: customerRepo,
		db:           db,
	}
}

func (uc *transactionUsecase) CreateTransaction(tx *model.Transaction) error {
	_, err := uc.customerRepo.FindByID(tx.CustomerID)
	if err != nil {
		return errors.New("customer not found")
	}

	limits, err := uc.limitRepo.FindByCustomerID(tx.CustomerID)
	if err != nil {
		return errors.New("limit not found")
	}

	var limitAmount int64
	for _, l := range limits {
		if l.Tenor == tx.Tenor {
			limitAmount = l.Limit
			break
		}
	}
	if limitAmount == 0 {
		return errors.New("limit for tenor not found")
	}

	totalUsed, err := uc.txRepo.SumUsedAmount(tx.CustomerID, tx.Tenor)
	if err != nil {
		return errors.New("failed to calculate used limit")
	}

	if totalUsed+tx.InstallmentAmount > limitAmount {
		return errors.New("transaction amount exceeds limit")
	}

	err = uc.db.Transaction(func(txDB *gorm.DB) error {
		if err := uc.txRepo.Create(txDB, tx); err != nil {
			return err
		}
		return nil
	})

	return err
}

func (uc *transactionUsecase) UpdateTransaction(id uint, updatedTx *model.Transaction) error {
	existingTx, err := uc.txRepo.FindByID(id)
	if err != nil {
		return errors.New("transaction not found")
	}

	if existingTx.CustomerID != updatedTx.CustomerID {
		return errors.New("customer ID mismatch")
	}

	limits, err := uc.limitRepo.FindByCustomerID(updatedTx.CustomerID)
	if err != nil {
		return errors.New("limit not found")
	}

	var limitAmount int64
	for _, l := range limits {
		if l.Tenor == updatedTx.Tenor {
			limitAmount = l.Limit
			break
		}
	}
	if limitAmount == 0 {
		return errors.New("limit for tenor not found")
	}

	totalUsed, err := uc.txRepo.SumUsedAmount(updatedTx.CustomerID, updatedTx.Tenor)
	if err != nil {
		return errors.New("failed to calculate used limit")
	}

	newTotal := totalUsed - existingTx.InstallmentAmount + updatedTx.InstallmentAmount
	if newTotal > limitAmount {
		return errors.New("transaction amount exceeds limit")
	}

	err = uc.db.Transaction(func(txDB *gorm.DB) error {
		if err := uc.txRepo.Update(txDB, id, updatedTx); err != nil {
			return err
		}
		return nil
	})

	return err
}

func (uc *transactionUsecase) DeleteTransaction(id uint) error {
	_, err := uc.txRepo.FindByID(id)
	if err != nil {
		return errors.New("transaction not found")
	}

	return uc.txRepo.Delete(id)
}

func (uc *transactionUsecase) GetTransactionByID(id uint) (*model.Transaction, error) {
	return uc.txRepo.FindByID(id)
}

func (uc *transactionUsecase) GetTransactionsByCustomer(customerID uint) ([]model.Transaction, error) {
	return uc.txRepo.FindByCustomerID(customerID)
}

func (uc *transactionUsecase) GetAllTransactions() ([]model.Transaction, error) {
	return uc.txRepo.FindAll()
}
