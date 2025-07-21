package usecase

import (
	"errors"
	"xyz-multifinance/internal/model"
	"xyz-multifinance/internal/repository"
)

type LimitUsecase interface {
	CreateLimit(limit *model.Limit) error
	UpdateLimit(id uint, fields map[string]interface{}) error
	DeleteLimit(id uint) error
	GetLimitByID(id uint) (*LimitWithRemaining, error)
	GetLimitsByCustomer(customerID uint) ([]LimitWithRemaining, error)
}

type limitUsecase struct {
	limitRepo       repository.LimitRepository
	transactionRepo repository.TransactionRepository
}

func NewLimitUsecase(limitRepo repository.LimitRepository, transactionRepo repository.TransactionRepository) LimitUsecase {
	return &limitUsecase{
		limitRepo:       limitRepo,
		transactionRepo: transactionRepo,
	}
}

func (uc *limitUsecase) CreateLimit(limit *model.Limit) error {
	if limit.Tenor <= 0 {
		return errors.New("tenor must be greater than zero")
	}
	if limit.Limit <= 0 {
		return errors.New("limit must be greater than zero")
	}
	return uc.limitRepo.Create(limit)
}

func (uc *limitUsecase) UpdateLimit(id uint, updatedFields map[string]interface{}) error {
	_, err := uc.limitRepo.FindByID(id)
	if err != nil {
		return errors.New("limit not found")
	}

	return uc.limitRepo.Update(id, updatedFields)
}

func (uc *limitUsecase) DeleteLimit(id uint) error {
	_, err := uc.limitRepo.FindByID(id)
	if err != nil {
		return errors.New("limit not found")
	}

	return uc.limitRepo.Delete(id)
}

type LimitWithRemaining struct {
	model.Limit
	RemainingLimit int64 `json:"remaining_limit"`
}

func (uc *limitUsecase) GetLimitByID(id uint) (*LimitWithRemaining, error) {
	limit, err := uc.limitRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if limit == nil {
		return nil, errors.New("limit not found")
	}

	usedAmount, err := uc.transactionRepo.SumUsedAmount(limit.CustomerID, limit.Tenor)
	if err != nil {
		return nil, err
	}

	remaining := limit.Limit - usedAmount

	return &LimitWithRemaining{
		Limit:          *limit,
		RemainingLimit: remaining,
	}, nil
}

func (uc *limitUsecase) GetLimitsByCustomer(customerID uint) ([]LimitWithRemaining, error) {
	limits, err := uc.limitRepo.FindByCustomerID(customerID)
	if err != nil {
		return nil, err
	}

	var result []LimitWithRemaining
	for _, l := range limits {
		usedAmount, err := uc.transactionRepo.SumUsedAmount(customerID, l.Tenor)
		if err != nil {
			return nil, err
		}

		remaining := l.Limit - usedAmount

		result = append(result, LimitWithRemaining{
			Limit:          l,
			RemainingLimit: remaining,
		})
	}

	return result, nil
}
