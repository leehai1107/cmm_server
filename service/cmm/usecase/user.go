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

type IUserUsecase interface {
	Login(ctx context.Context, req request.Login) (string, error)
	Register(ctx context.Context, req request.Register) error
}

type userUsecase struct {
	repo repository.IUserRepo
}

func NewUserUsecase(
	repo repository.IUserRepo,
) IUserUsecase {
	return &userUsecase{
		repo: repo,
	}
}

func (u *userUsecase) Login(ctx context.Context, req request.Login) (string, error) {
	user, err := u.repo.GetUserByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("invalid credentials") // avoid revealing that user doesn't exist
		}
		return "", err
	}

	// Compare password hash
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Generate and return token (placeholder)
	return "jwt_token_placeholder", nil
}

func (u *userUsecase) Register(ctx context.Context, req request.Register) error {
	logger.EnhanceWith(ctx).Info("Register usecase called")

	// Check if u is nil
	if u == nil {
		logger.EnhanceWith(ctx).Error("Usecase instance is nil")
		return errors.New("usecase instance is nil")
	}

	// Ensure userRepo is not nil
	if u.repo == nil {
		logger.EnhanceWith(ctx).Error("Repository is nil in Register usecase")
		return errors.New("repository is not initialized")
	}

	logger.EnhanceWith(ctx).Infow("Checking if email already exists", "email", req.Email)
	// Check if email already exists
	_, err := u.repo.GetUserByEmail(req.Email)
	if err == nil {
		return errors.New("email already registered")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
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
		Role:         entity.RoleAdmin, // Adjust role as needed
		CreatedAt:    time.Now(),
	}

	// Save user
	if err := u.repo.CreateUser(user); err != nil {
		return err
	}

	// Create wallet
	wallet := &entity.Wallet{
		UserID:  user.ID,
		Balance: 0,
	}

	return u.repo.CreateWallet(wallet)
}
