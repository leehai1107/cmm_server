package repository

import (
	"github.com/google/uuid"
	"github.com/leehai1107/cmm_server/pkg/logger"
	"github.com/leehai1107/cmm_server/service/cmm/model/entity"
	"gorm.io/gorm"
)

type ITransactionRepo interface {
	CreateTransaction(transaction *entity.Transaction) error
	GetTransactionByID(id uuid.UUID) (*entity.Transaction, error)
	GetTransactionsByUser(userID uuid.UUID) ([]entity.Transaction, error)
	GetTransactionsByService(serviceID int, serviceRefID uuid.UUID) ([]entity.Transaction, error)
	UpdateTransactionStatus(id uuid.UUID, status string) error
}

type transactionRepo struct {
	db *gorm.DB
}

func NewTransactionRepo(db *gorm.DB) ITransactionRepo {
	return &transactionRepo{
		db: db,
	}
}

func (r *transactionRepo) CreateTransaction(transaction *entity.Transaction) error {
	logger.Info("CreateTransaction repository method called")
	return r.db.Create(transaction).Error
}

func (r *transactionRepo) GetTransactionByID(id uuid.UUID) (*entity.Transaction, error) {
	logger.Info("GetTransactionByID repository method called")
	var transaction entity.Transaction
	err := r.db.Where("id = ?", id).First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *transactionRepo) GetTransactionsByUser(userID uuid.UUID) ([]entity.Transaction, error) {
	logger.Info("GetTransactionsByUser repository method called")
	var transactions []entity.Transaction
	err := r.db.Where("user_id = ?", userID).Order("paid_at DESC").Find(&transactions).Error
	return transactions, err
}

func (r *transactionRepo) GetTransactionsByService(serviceID int, serviceRefID uuid.UUID) ([]entity.Transaction, error) {
	logger.Info("GetTransactionsByService repository method called")
	var transactions []entity.Transaction
	err := r.db.Where("service_id = ? AND service_ref_id = ?", serviceID, serviceRefID).
		Order("paid_at DESC").Find(&transactions).Error
	return transactions, err
}

func (r *transactionRepo) UpdateTransactionStatus(id uuid.UUID, status string) error {
	logger.Info("UpdateTransactionStatus repository method called")
	return r.db.Model(&entity.Transaction{}).Where("id = ?", id).Update("status", status).Error
}
