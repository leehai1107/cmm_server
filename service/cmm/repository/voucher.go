package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/leehai1107/cmm_server/pkg/logger"
	"github.com/leehai1107/cmm_server/service/cmm/model/entity"
	"gorm.io/gorm"
)

type IVoucherRepo interface {
	CreateVoucher(voucher *entity.Voucher) error
	GetVoucherByID(id uuid.UUID) (*entity.Voucher, error)
	GetVoucherByCode(code string) (*entity.Voucher, error)
	GetAllVouchers() ([]entity.Voucher, error)
	GetValidVouchers() ([]entity.Voucher, error)
	UpdateVoucher(voucher *entity.Voucher) error
	IncrementUsedCount(id uuid.UUID) error
	DeleteVoucher(id uuid.UUID) error
}

type voucherRepo struct {
	db *gorm.DB
}

func NewVoucherRepo(db *gorm.DB) IVoucherRepo {
	return &voucherRepo{
		db: db,
	}
}

func (r *voucherRepo) CreateVoucher(voucher *entity.Voucher) error {
	logger.Info("CreateVoucher repository method called")
	return r.db.Create(voucher).Error
}

func (r *voucherRepo) GetVoucherByID(id uuid.UUID) (*entity.Voucher, error) {
	logger.Info("GetVoucherByID repository method called")
	var voucher entity.Voucher
	err := r.db.Where("id = ?", id).First(&voucher).Error
	if err != nil {
		return nil, err
	}
	return &voucher, nil
}

func (r *voucherRepo) GetVoucherByCode(code string) (*entity.Voucher, error) {
	logger.Info("GetVoucherByCode repository method called")
	var voucher entity.Voucher
	err := r.db.Where("code = ?", code).First(&voucher).Error
	if err != nil {
		return nil, err
	}
	return &voucher, nil
}

func (r *voucherRepo) GetAllVouchers() ([]entity.Voucher, error) {
	logger.Info("GetAllVouchers repository method called")
	var vouchers []entity.Voucher
	err := r.db.Find(&vouchers).Error
	return vouchers, err
}

func (r *voucherRepo) GetValidVouchers() ([]entity.Voucher, error) {
	logger.Info("GetValidVouchers repository method called")
	now := time.Now()
	var vouchers []entity.Voucher
	err := r.db.Where("valid_from <= ? AND valid_to >= ? AND (max_uses = 0 OR used_count < max_uses)",
		now, now).Find(&vouchers).Error
	return vouchers, err
}

func (r *voucherRepo) UpdateVoucher(voucher *entity.Voucher) error {
	logger.Info("UpdateVoucher repository method called")
	return r.db.Save(voucher).Error
}

func (r *voucherRepo) IncrementUsedCount(id uuid.UUID) error {
	logger.Info("IncrementUsedCount repository method called")
	return r.db.Model(&entity.Voucher{}).
		Where("id = ?", id).
		Update("used_count", gorm.Expr("used_count + 1")).Error
}

func (r *voucherRepo) DeleteVoucher(id uuid.UUID) error {
	logger.Info("DeleteVoucher repository method called")
	return r.db.Delete(&entity.Voucher{}, "id = ?", id).Error
}
