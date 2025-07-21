package usecase

import (
	"errors"
	"fmt"
	"time"

	"xyz-multifinance/internal/model"
	"xyz-multifinance/internal/repository"
)

type CustomerUsecase interface {
	CreateCustomer(cust *model.Customer) error
	GetCustomerByNIK(nik string) (*model.Customer, error)
	UpdateCustomer(nik string, updatedFields map[string]interface{}) error
	DeleteCustomer(nik string) error
}

type customerUsecase struct {
	customerRepo repository.CustomerRepository
	userRepo     repository.UserRepository
}

func NewCustomerUsecase(repo repository.CustomerRepository, uRepo repository.UserRepository) CustomerUsecase {
	return &customerUsecase{
		customerRepo: repo,
		userRepo:     uRepo,
	}
}

func (uc *customerUsecase) CreateCustomer(customer *model.Customer) error {
	user, err := uc.userRepo.FindByID(customer.UserID)
	if err != nil {
		return fmt.Errorf("failed to check user: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user_id %d not found", customer.UserID)
	}

	existing, _ := uc.customerRepo.FindByNIK(customer.NIK)
	if existing != nil && existing.NIK != "" {
		return errors.New("customer with this NIK already exists")
	}

	customer.CreatedAt = time.Now()
	customer.UpdatedAt = time.Now()

	return uc.customerRepo.Create(customer)
}

func (uc *customerUsecase) GetCustomerByNIK(nik string) (*model.Customer, error) {
	return uc.customerRepo.FindByNIK(nik)
}

func (uc *customerUsecase) UpdateCustomer(nik string, updatedFields map[string]interface{}) error {
	_, err := uc.customerRepo.FindByNIK(nik)
	if err != nil {
		return err
	}

	updatedFields["updated_at"] = time.Now()

	return uc.customerRepo.Update(nik, updatedFields)
}

func (uc *customerUsecase) DeleteCustomer(nik string) error {
	customer, err := uc.customerRepo.FindByNIK(nik)
	if err != nil {
		return errors.New("customer not found")
	}

	return uc.customerRepo.Delete(customer.ID)
}
