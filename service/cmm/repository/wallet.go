package repository

import (
	"github.com/google/uuid"
	"github.com/leehai1107/cmm_server/pkg/logger"
	"github.com/leehai1107/cmm_server/service/cmm/model/entity"
	"gorm.io/gorm"
)

type IWalletRepo interface {
	GetWalletByUserID(userID uuid.UUID) (*entity.Wallet, error)
	CreateWallet(wallet *entity.Wallet) error
	UpdateBalance(userID uuid.UUID, amount float64) error
	DeductBalance(userID uuid.UUID, amount float64) error

	// Top-up methods
	CreateTopup(topup *entity.Topup) error
	GetTopupByID(id uuid.UUID) (*entity.Topup, error)
	GetTopupsByUser(userID uuid.UUID) ([]entity.Topup, error)
	UpdateTopupStatus(id uuid.UUID, status string) error
}

type walletRepo struct {
	db *gorm.DB
}

func NewWalletRepo(db *gorm.DB) IWalletRepo {
	return &walletRepo{
		db: db,
	}
}

func (r *walletRepo) GetWalletByUserID(userID uuid.UUID) (*entity.Wallet, error) {
	logger.Info("GetWalletByUserID repository method called")
	var wallet entity.Wallet
	err := r.db.Where("user_id = ?", userID).First(&wallet).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *walletRepo) CreateWallet(wallet *entity.Wallet) error {
	logger.Info("CreateWallet repository method called")
	return r.db.Create(wallet).Error
}

func (r *walletRepo) UpdateBalance(userID uuid.UUID, amount float64) error {
	logger.Info("UpdateBalance repository method called")
	return r.db.Model(&entity.Wallet{}).
		Where("user_id = ?", userID).
		Update("balance", gorm.Expr("balance + ?", amount)).Error
}

func (r *walletRepo) DeductBalance(userID uuid.UUID, amount float64) error {
	logger.Info("DeductBalance repository method called")
	return r.db.Model(&entity.Wallet{}).
		Where("user_id = ? AND balance >= ?", userID, amount).
		Update("balance", gorm.Expr("balance - ?", amount)).Error
}

func (r *walletRepo) CreateTopup(topup *entity.Topup) error {
	logger.Info("CreateTopup repository method called")
	return r.db.Create(topup).Error
}

func (r *walletRepo) GetTopupByID(id uuid.UUID) (*entity.Topup, error) {
	logger.Info("GetTopupByID repository method called")
	var topup entity.Topup
	err := r.db.Where("id = ?", id).First(&topup).Error
	if err != nil {
		return nil, err
	}
	return &topup, nil
}

func (r *walletRepo) GetTopupsByUser(userID uuid.UUID) ([]entity.Topup, error) {
	logger.Info("GetTopupsByUser repository method called")
	var topups []entity.Topup
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&topups).Error
	return topups, err
}

func (r *walletRepo) UpdateTopupStatus(id uuid.UUID, status string) error {
	logger.Info("UpdateTopupStatus repository method called")
	return r.db.Model(&entity.Topup{}).Where("id = ?", id).Update("status", status).Error
}
