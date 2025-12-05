package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/leehai1107/cmm_server/pkg/logger"
	"github.com/leehai1107/cmm_server/service/cmm/model/entity"
	"github.com/leehai1107/cmm_server/service/cmm/model/request"
	"github.com/leehai1107/cmm_server/service/cmm/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type IAdminUsecase interface {
	CreateAccount(ctx context.Context, req request.CreateAccount) error
}

type adminUsecase struct {
	repo repository.IUserRepo
}

func NewAdminUsecase(
	repo repository.IUserRepo,
) IAdminUsecase {
	return &adminUsecase{
		repo: repo,
	}
}

func (a *adminUsecase) CreateAccount(ctx context.Context, req request.CreateAccount) error {
	log := logger.EnhanceWith(ctx)

	// Check if email already exists
	_, err := a.repo.GetUserByEmail(req.Email)
	if err == nil {
		log.Errorw("Email already exists", "email", req.Email)
		return errors.New("email already registered")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Errorw("Failed to check if email exists", "error", err)
		return err // Return other DB errors
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Create user entity
	user := &entity.User{
		ID:           uuid.New(),
		FullName:     req.FirstName + " " + req.LastName,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         entity.RoleAdmin,
		CreatedAt:    time.Now(),
	}

	// Save user
	if err := a.repo.CreateUser(user); err != nil {
		log.Errorw("Failed to create user", "error", err)
		return err
	}

	// Create wallet
	wallet := &entity.Wallet{
		UserID:  user.ID,
		Balance: 0,
	}

	return a.repo.CreateWallet(wallet)
}
