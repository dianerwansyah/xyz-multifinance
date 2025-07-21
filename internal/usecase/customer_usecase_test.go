package usecase_test

import (
	"testing"

	"xyz-multifinance/internal/model"
	"xyz-multifinance/internal/usecase"
)

type mockCustomerRepo struct {
	FindByNIKFunc func(nik string) (*model.Customer, error)
	FindByIDFunc  func(id uint) (*model.Customer, error)
	CreateFunc    func(customer *model.Customer) error
	UpdateFunc    func(nik string, fields map[string]interface{}) error
	DeleteFunc    func(id uint) error
}

func (m *mockCustomerRepo) FindByNIK(nik string) (*model.Customer, error) {
	if m.FindByNIKFunc != nil {
		return m.FindByNIKFunc(nik)
	}
	return nil, nil
}

func (m *mockCustomerRepo) FindByID(id uint) (*model.Customer, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}

func (m *mockCustomerRepo) Create(customer *model.Customer) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(customer)
	}
	return nil
}

func (m *mockCustomerRepo) Update(nik string, fields map[string]interface{}) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(nik, fields)
	}
	return nil
}

func (m *mockCustomerRepo) Delete(id uint) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}

type mockUserRepo struct {
	FindByUsernameFunc func(username string) (*model.User, error)
	FindByIDFunc       func(id uint) (*model.User, error)
	CreateFunc         func(user *model.User) error
}

func (m *mockUserRepo) FindByUsername(username string) (*model.User, error) {
	if m.FindByUsernameFunc != nil {
		return m.FindByUsernameFunc(username)
	}
	return nil, nil
}

func (m *mockUserRepo) FindByID(id uint) (*model.User, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}

func (m *mockUserRepo) Create(user *model.User) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(user)
	}
	return nil
}

func TestCreateCustomer_Success(t *testing.T) {
	mockUser := &model.User{ID: 1, Username: "Admin"}

	uc := usecase.NewCustomerUsecase(&mockCustomerRepo{
		FindByNIKFunc: func(nik string) (*model.Customer, error) {
			return nil, nil
		},
		CreateFunc: func(c *model.Customer) error {
			return nil
		},
	}, &mockUserRepo{
		FindByIDFunc: func(id uint) (*model.User, error) {
			return mockUser, nil
		},
	})

	customer := &model.Customer{
		FullName:   "Agustiansyah",
		NIK:        "1234567890",
		UserID:     1,
		PlaceBirth: "Jakarta",
	}

	err := uc.CreateCustomer(customer)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

func TestCreateCustomer_NIKExists(t *testing.T) {
	existing := &model.Customer{NIK: "1234567890"}

	uc := usecase.NewCustomerUsecase(&mockCustomerRepo{
		FindByNIKFunc: func(nik string) (*model.Customer, error) {
			return existing, nil
		},
		CreateFunc: func(c *model.Customer) error {
			return nil
		},
	}, &mockUserRepo{
		FindByIDFunc: func(id uint) (*model.User, error) {
			return &model.User{ID: 1, Username: "Admin"}, nil
		},
	})

	customer := &model.Customer{
		NIK:    "1234567890",
		UserID: 1,
	}

	err := uc.CreateCustomer(customer)
	if err == nil || err.Error() != "customer with this NIK already exists" {
		t.Errorf("expected error about NIK, got %v", err)
	}
}
