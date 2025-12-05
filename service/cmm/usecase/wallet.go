package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/leehai1107/cmm_server/pkg/logger"
	"github.com/leehai1107/cmm_server/service/cmm/model/entity"
	"github.com/leehai1107/cmm_server/service/cmm/model/request"
	"github.com/leehai1107/cmm_server/service/cmm/model/response"
	"github.com/leehai1107/cmm_server/service/cmm/repository"
	"gorm.io/gorm"
)

type IWalletUsecase interface {
	GetWallet(ctx context.Context, userID uuid.UUID) (*response.WalletResponse, error)
	CreateTopup(ctx context.Context, userID uuid.UUID, req request.TopupWallet) (*response.TopupResponse, error)
	ConfirmTopup(ctx context.Context, req request.ConfirmTopup) error
	GetTopupHistory(ctx context.Context, userID uuid.UUID) ([]response.TopupResponse, error)
	GetTransactionHistory(ctx context.Context, userID uuid.UUID) ([]response.TransactionResponse, error)
}

type walletUsecase struct {
	walletRepo      repository.IWalletRepo
	transactionRepo repository.ITransactionRepo
}

func NewWalletUsecase(
	walletRepo repository.IWalletRepo,
	transactionRepo repository.ITransactionRepo,
) IWalletUsecase {
	return &walletUsecase{
		walletRepo:      walletRepo,
		transactionRepo: transactionRepo,
	}
}

func (u *walletUsecase) GetWallet(ctx context.Context, userID uuid.UUID) (*response.WalletResponse, error) {
	wallet, err := u.walletRepo.GetWalletByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("wallet not found")
		}
		return nil, err
	}

	return &response.WalletResponse{
		UserID:  wallet.UserID,
		Balance: wallet.Balance,
	}, nil
}

func (u *walletUsecase) CreateTopup(ctx context.Context, userID uuid.UUID, req request.TopupWallet) (*response.TopupResponse, error) {
	logger.EnhanceWith(ctx).Info("CreateTopup usecase called")

	// Verify wallet exists
	_, err := u.walletRepo.GetWalletByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("wallet not found")
		}
		return nil, err
	}

	topup := &entity.Topup{
		ID:        uuid.New(),
		UserID:    userID,
		Amount:    req.Amount,
		Method:    req.Method,
		Status:    "pending",
		CreatedAt: time.Now(),
	}

	if err := u.walletRepo.CreateTopup(topup); err != nil {
		return nil, err
	}

	return &response.TopupResponse{
		ID:        topup.ID,
		UserID:    topup.UserID,
		Amount:    topup.Amount,
		Method:    topup.Method,
		Status:    topup.Status,
		CreatedAt: topup.CreatedAt,
	}, nil
}

func (u *walletUsecase) ConfirmTopup(ctx context.Context, req request.ConfirmTopup) error {
	log := logger.EnhanceWith(ctx)
	log.Info("ConfirmTopup usecase called")

	topup, err := u.walletRepo.GetTopupByID(req.TopupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("topup not found")
		}
		return err
	}

	if topup.Status != "pending" {
		return errors.New("topup is not pending")
	}

	// Update topup status
	if err := u.walletRepo.UpdateTopupStatus(req.TopupID, "completed"); err != nil {
		return err
	}

	// Add balance to wallet
	if err := u.walletRepo.UpdateBalance(topup.UserID, topup.Amount); err != nil {
		log.Errorw("Failed to update wallet balance", "error", err)
		return errors.New("failed to update wallet balance")
	}

	// Create transaction record
	transaction := &entity.Transaction{
		ID:           uuid.New(),
		UserID:       topup.UserID,
		ServiceID:    2, // 2 for topup
		ServiceRefID: topup.ID,
		Amount:       topup.Amount,
		PaidAt:       time.Now(),
		Status:       "completed",
	}
	if err := u.transactionRepo.CreateTransaction(transaction); err != nil {
		log.Errorw("Failed to create transaction", "error", err)
	}

	return nil
}

func (u *walletUsecase) GetTopupHistory(ctx context.Context, userID uuid.UUID) ([]response.TopupResponse, error) {
	topups, err := u.walletRepo.GetTopupsByUser(userID)
	if err != nil {
		return nil, err
	}

	var result []response.TopupResponse
	for _, topup := range topups {
		result = append(result, response.TopupResponse{
			ID:        topup.ID,
			UserID:    topup.UserID,
			Amount:    topup.Amount,
			Method:    topup.Method,
			Status:    topup.Status,
			CreatedAt: topup.CreatedAt,
		})
	}

	return result, nil
}

func (u *walletUsecase) GetTransactionHistory(ctx context.Context, userID uuid.UUID) ([]response.TransactionResponse, error) {
	transactions, err := u.transactionRepo.GetTransactionsByUser(userID)
	if err != nil {
		return nil, err
	}

	var result []response.TransactionResponse
	for _, transaction := range transactions {
		result = append(result, response.TransactionResponse{
			ID:              transaction.ID,
			UserID:          transaction.UserID,
			ServiceID:       transaction.ServiceID,
			ServiceRefID:    transaction.ServiceRefID,
			Amount:          transaction.Amount,
			PaymentMethodID: transaction.PaymentMethodID,
			PaidAt:          transaction.PaidAt,
			Status:          transaction.Status,
		})
	}

	return result, nil
}
