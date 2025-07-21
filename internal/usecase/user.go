package usecase

import (
	"errors"
	"time"

	"xyz-multifinance/internal/model"
	"xyz-multifinance/internal/repository"
	"xyz-multifinance/logger"
	"xyz-multifinance/pkg/jwtutil"

	"golang.org/x/crypto/bcrypt"
)

type UserUsecase interface {
	Login(username string, password string) (string, error)
	CreateUser(user *model.User) error
}

type userUsecase struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &userUsecase{userRepo: userRepo}
}

func (u *userUsecase) Login(username string, password string) (string, error) {
	user, err := u.userRepo.FindByUsername(username)
	if err != nil {
		return "", errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		logger.Log.Errorf("bcrypt compare error: %v", err)
		return "", errors.New("invalid credentials")
	}

	token, err := jwtutil.GenerateToken(user.ID, user.Role, time.Now().Add(24*time.Hour))
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}

func (uc *userUsecase) CreateUser(user *model.User) error {
	existingUser, _ := uc.userRepo.FindByUsername(user.Username)
	if existingUser != nil {
		return errors.New("username already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	user.Role = "user"
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	return uc.userRepo.Create(user)
}
