package repository

import (
	"errors"

	"github.com/leehai1107/cmm_server/pkg/logger"
	"github.com/leehai1107/cmm_server/service/cmm/model/entity"
	"gorm.io/gorm"
)

type IUserRepo interface {
	GetUserByEmail(email string) (*entity.User, error)
	CreateUser(user *entity.User) error
	CreateWallet(wallet *entity.Wallet) error
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) IUserRepo {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) GetUserByEmail(email string) (*entity.User, error) {
	logger.Info("GetUserByEmail repository method called")

	// Check if r is nil
	if r == nil {
		logger.Error("Repository instance is nil")
		return nil, errors.New("repository instance is nil")
	}

	if r.db == nil {
		logger.Error("Database connection is nil in GetUserByEmail")
		return nil, errors.New("database connection is not initialized")
	}

	logger.Infof("Querying database for user with email: %s", email)
	var user entity.User
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			logger.Infof("User with email %s not found", email)
			return nil, gorm.ErrRecordNotFound
		}
		logger.Errorf("Database error in GetUserByEmail: %v", result.Error)
		return nil, result.Error
	}
	logger.Infof("User with email %s found", email)
	return &user, nil
}

func (r *userRepo) CreateUser(user *entity.User) error {
	logger.Info("CreateUser repository method called")

	// Check if r is nil
	if r == nil {
		logger.Error("Repository instance is nil in CreateUser")
		return errors.New("repository instance is nil")
	}

	if r.db == nil {
		logger.Error("Database connection is nil in CreateUser")
		return errors.New("database connection is not initialized")
	}

	if user == nil {
		logger.Error("User entity is nil in CreateUser")
		return errors.New("user cannot be nil")
	}

	logger.Infof("Creating user with email: %s", user.Email)
	err := r.db.Create(user).Error
	if err != nil {
		logger.Errorf("Error creating user: %v", err)
		return err
	}

	logger.Infof("User created successfully with ID: %s", user.ID)
	return nil
}

func (r *userRepo) CreateWallet(wallet *entity.Wallet) error {
	logger.Info("CreateWallet repository method called")

	// Check if r is nil
	if r == nil {
		logger.Error("Repository instance is nil in CreateWallet")
		return errors.New("repository instance is nil")
	}

	if r.db == nil {
		logger.Error("Database connection is nil in CreateWallet")
		return errors.New("database connection is not initialized")
	}

	if wallet == nil {
		logger.Error("Wallet entity is nil in CreateWallet")
		return errors.New("wallet cannot be nil")
	}

	logger.Infof("Creating wallet for user ID: %s", wallet.UserID)
	err := r.db.Create(wallet).Error
	if err != nil {
		logger.Errorf("Error creating wallet: %v", err)
		return err
	}

	logger.Infof("Wallet created successfully for user ID: %s", wallet.UserID)
	return nil
}
