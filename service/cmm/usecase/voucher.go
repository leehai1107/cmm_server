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

type IVoucherUsecase interface {
	CreateVoucher(ctx context.Context, req request.CreateVoucher) (*response.VoucherResponse, error)
	GetVoucher(ctx context.Context, voucherID uuid.UUID) (*response.VoucherResponse, error)
	GetAllVouchers(ctx context.Context) ([]response.VoucherResponse, error)
	GetValidVouchers(ctx context.Context) ([]response.VoucherResponse, error)
	UpdateVoucher(ctx context.Context, req request.UpdateVoucher) error
	DeleteVoucher(ctx context.Context, voucherID uuid.UUID) error
	ApplyVoucher(ctx context.Context, req request.ApplyVoucher) (*response.VoucherCalculation, error)
}

type voucherUsecase struct {
	voucherRepo repository.IVoucherRepo
}

func NewVoucherUsecase(voucherRepo repository.IVoucherRepo) IVoucherUsecase {
	return &voucherUsecase{
		voucherRepo: voucherRepo,
	}
}

func (u *voucherUsecase) CreateVoucher(ctx context.Context, req request.CreateVoucher) (*response.VoucherResponse, error) {
	logger.EnhanceWith(ctx).Info("CreateVoucher usecase called")

	// Check if voucher code already exists
	_, err := u.voucherRepo.GetVoucherByCode(req.Code)
	if err == nil {
		return nil, errors.New("voucher code already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Validate dates
	if req.ValidTo.Before(req.ValidFrom) {
		return nil, errors.New("valid_to must be after valid_from")
	}

	voucher := &entity.Voucher{
		ID:              uuid.New(),
		Code:            req.Code,
		DiscountPercent: req.DiscountPercent,
		MaxUses:         req.MaxUses,
		UsedCount:       0,
		ServiceID:       req.ServiceID,
		ValidFrom:       req.ValidFrom,
		ValidTo:         req.ValidTo,
	}

	if err := u.voucherRepo.CreateVoucher(voucher); err != nil {
		return nil, err
	}

	return &response.VoucherResponse{
		ID:              voucher.ID,
		Code:            voucher.Code,
		DiscountPercent: voucher.DiscountPercent,
		MaxUses:         voucher.MaxUses,
		UsedCount:       voucher.UsedCount,
		ServiceID:       voucher.ServiceID,
		ValidFrom:       voucher.ValidFrom,
		ValidTo:         voucher.ValidTo,
	}, nil
}

func (u *voucherUsecase) GetVoucher(ctx context.Context, voucherID uuid.UUID) (*response.VoucherResponse, error) {
	voucher, err := u.voucherRepo.GetVoucherByID(voucherID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("voucher not found")
		}
		return nil, err
	}

	return &response.VoucherResponse{
		ID:              voucher.ID,
		Code:            voucher.Code,
		DiscountPercent: voucher.DiscountPercent,
		MaxUses:         voucher.MaxUses,
		UsedCount:       voucher.UsedCount,
		ServiceID:       voucher.ServiceID,
		ValidFrom:       voucher.ValidFrom,
		ValidTo:         voucher.ValidTo,
	}, nil
}

func (u *voucherUsecase) GetAllVouchers(ctx context.Context) ([]response.VoucherResponse, error) {
	vouchers, err := u.voucherRepo.GetAllVouchers()
	if err != nil {
		return nil, err
	}

	var result []response.VoucherResponse
	for _, voucher := range vouchers {
		result = append(result, response.VoucherResponse{
			ID:              voucher.ID,
			Code:            voucher.Code,
			DiscountPercent: voucher.DiscountPercent,
			MaxUses:         voucher.MaxUses,
			UsedCount:       voucher.UsedCount,
			ServiceID:       voucher.ServiceID,
			ValidFrom:       voucher.ValidFrom,
			ValidTo:         voucher.ValidTo,
		})
	}

	return result, nil
}

func (u *voucherUsecase) GetValidVouchers(ctx context.Context) ([]response.VoucherResponse, error) {
	vouchers, err := u.voucherRepo.GetValidVouchers()
	if err != nil {
		return nil, err
	}

	var result []response.VoucherResponse
	for _, voucher := range vouchers {
		result = append(result, response.VoucherResponse{
			ID:              voucher.ID,
			Code:            voucher.Code,
			DiscountPercent: voucher.DiscountPercent,
			MaxUses:         voucher.MaxUses,
			UsedCount:       voucher.UsedCount,
			ServiceID:       voucher.ServiceID,
			ValidFrom:       voucher.ValidFrom,
			ValidTo:         voucher.ValidTo,
		})
	}

	return result, nil
}

func (u *voucherUsecase) UpdateVoucher(ctx context.Context, req request.UpdateVoucher) error {
	voucher, err := u.voucherRepo.GetVoucherByID(req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("voucher not found")
		}
		return err
	}

	if req.Code != "" {
		// Check if new code conflicts with existing vouchers
		existing, err := u.voucherRepo.GetVoucherByCode(req.Code)
		if err == nil && existing.ID != req.ID {
			return errors.New("voucher code already exists")
		}
		voucher.Code = req.Code
	}
	if req.DiscountPercent > 0 {
		voucher.DiscountPercent = req.DiscountPercent
	}
	if req.MaxUses >= 0 {
		voucher.MaxUses = req.MaxUses
	}
	if req.ValidFrom != nil {
		voucher.ValidFrom = *req.ValidFrom
	}
	if req.ValidTo != nil {
		voucher.ValidTo = *req.ValidTo
	}

	// Validate dates
	if voucher.ValidTo.Before(voucher.ValidFrom) {
		return errors.New("valid_to must be after valid_from")
	}

	return u.voucherRepo.UpdateVoucher(voucher)
}

func (u *voucherUsecase) DeleteVoucher(ctx context.Context, voucherID uuid.UUID) error {
	_, err := u.voucherRepo.GetVoucherByID(voucherID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("voucher not found")
		}
		return err
	}

	return u.voucherRepo.DeleteVoucher(voucherID)
}

func (u *voucherUsecase) ApplyVoucher(ctx context.Context, req request.ApplyVoucher) (*response.VoucherCalculation, error) {
	voucher, err := u.voucherRepo.GetVoucherByCode(req.VoucherCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid voucher code")
		}
		return nil, err
	}

	// Validate voucher
	now := time.Now()
	if now.Before(voucher.ValidFrom) || now.After(voucher.ValidTo) {
		return nil, errors.New("voucher is not valid at this time")
	}
	if voucher.MaxUses > 0 && voucher.UsedCount >= voucher.MaxUses {
		return nil, errors.New("voucher has reached maximum uses")
	}

	// Calculate discount
	discountAmount := req.Amount * float64(voucher.DiscountPercent) / 100
	finalAmount := req.Amount - discountAmount

	return &response.VoucherCalculation{
		OriginalAmount:  req.Amount,
		DiscountAmount:  discountAmount,
		FinalAmount:     finalAmount,
		DiscountPercent: voucher.DiscountPercent,
		VoucherCode:     voucher.Code,
	}, nil
}
