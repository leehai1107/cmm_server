package infra

import (
	"github.com/leehai1107/cmm_server/pkg/logger"
	"github.com/leehai1107/cmm_server/service/cmm/model/entity"
	"gorm.io/gorm"
)

// RunMigrations runs all database migrations
func RunMigrations(db *gorm.DB) error {
	logger.Info("Starting database migrations...")

	// List of all entities to migrate
	models := []interface{}{
		&entity.User{},
		&entity.Wallet{},
		&entity.Topup{},
		&entity.CoffeeShop{},
		&entity.CommissionRate{},
		&entity.MeetingRoom{},
		&entity.Booking{},
		&entity.Voucher{},
		&entity.Transaction{},
		&entity.ShopPost{},
		&entity.InternalPost{},
	}

	// Auto migrate all models
	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			logger.Errorf("Failed to migrate model %T: %v", model, err)
			return err
		}
		logger.Infof("Successfully migrated model: %T", model)
	}

	// Add foreign key constraints
	if err := addForeignKeys(db); err != nil {
		logger.Errorf("Failed to add foreign keys: %v", err)
		return err
	}

	logger.Info("Database migrations completed successfully")
	return nil
}

// addForeignKeys adds foreign key constraints to the database
func addForeignKeys(db *gorm.DB) error {
	logger.Info("Adding foreign key constraints...")

	// CoffeeShop foreign keys
	if err := db.Exec(`
		ALTER TABLE coffee_shops 
		DROP CONSTRAINT IF EXISTS fk_coffee_shops_owner;
	`).Error; err != nil {
		logger.Warnf("Could not drop constraint fk_coffee_shops_owner: %v", err)
	}

	if err := db.Exec(`
		ALTER TABLE coffee_shops 
		ADD CONSTRAINT fk_coffee_shops_owner 
		FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE;
	`).Error; err != nil {
		logger.Warnf("Could not add constraint fk_coffee_shops_owner: %v", err)
	}

	// CommissionRate foreign keys
	if err := db.Exec(`
		ALTER TABLE commission_rates 
		DROP CONSTRAINT IF EXISTS fk_commission_rates_coffee_shop;
	`).Error; err != nil {
		logger.Warnf("Could not drop constraint fk_commission_rates_coffee_shop: %v", err)
	}

	if err := db.Exec(`
		ALTER TABLE commission_rates 
		ADD CONSTRAINT fk_commission_rates_coffee_shop 
		FOREIGN KEY (coffee_shop_id) REFERENCES coffee_shops(id) ON DELETE CASCADE;
	`).Error; err != nil {
		logger.Warnf("Could not add constraint fk_commission_rates_coffee_shop: %v", err)
	}

	// MeetingRoom foreign keys
	if err := db.Exec(`
		ALTER TABLE meeting_rooms 
		DROP CONSTRAINT IF EXISTS fk_meeting_rooms_coffee_shop;
	`).Error; err != nil {
		logger.Warnf("Could not drop constraint fk_meeting_rooms_coffee_shop: %v", err)
	}

	if err := db.Exec(`
		ALTER TABLE meeting_rooms 
		ADD CONSTRAINT fk_meeting_rooms_coffee_shop 
		FOREIGN KEY (coffee_shop_id) REFERENCES coffee_shops(id) ON DELETE CASCADE;
	`).Error; err != nil {
		logger.Warnf("Could not add constraint fk_meeting_rooms_coffee_shop: %v", err)
	}

	// Booking foreign keys
	if err := db.Exec(`
		ALTER TABLE bookings 
		DROP CONSTRAINT IF EXISTS fk_bookings_customer;
	`).Error; err != nil {
		logger.Warnf("Could not drop constraint fk_bookings_customer: %v", err)
	}

	if err := db.Exec(`
		ALTER TABLE bookings 
		ADD CONSTRAINT fk_bookings_customer 
		FOREIGN KEY (customer_id) REFERENCES users(id) ON DELETE CASCADE;
	`).Error; err != nil {
		logger.Warnf("Could not add constraint fk_bookings_customer: %v", err)
	}

	if err := db.Exec(`
		ALTER TABLE bookings 
		DROP CONSTRAINT IF EXISTS fk_bookings_meeting_room;
	`).Error; err != nil {
		logger.Warnf("Could not drop constraint fk_bookings_meeting_room: %v", err)
	}

	if err := db.Exec(`
		ALTER TABLE bookings 
		ADD CONSTRAINT fk_bookings_meeting_room 
		FOREIGN KEY (meeting_room_id) REFERENCES meeting_rooms(id) ON DELETE CASCADE;
	`).Error; err != nil {
		logger.Warnf("Could not add constraint fk_bookings_meeting_room: %v", err)
	}

	// Wallet foreign keys
	if err := db.Exec(`
		ALTER TABLE wallets 
		DROP CONSTRAINT IF EXISTS fk_wallets_user;
	`).Error; err != nil {
		logger.Warnf("Could not drop constraint fk_wallets_user: %v", err)
	}

	if err := db.Exec(`
		ALTER TABLE wallets 
		ADD CONSTRAINT fk_wallets_user 
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
	`).Error; err != nil {
		logger.Warnf("Could not add constraint fk_wallets_user: %v", err)
	}

	// Topup foreign keys
	if err := db.Exec(`
		ALTER TABLE topups 
		DROP CONSTRAINT IF EXISTS fk_topups_user;
	`).Error; err != nil {
		logger.Warnf("Could not drop constraint fk_topups_user: %v", err)
	}

	if err := db.Exec(`
		ALTER TABLE topups 
		ADD CONSTRAINT fk_topups_user 
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
	`).Error; err != nil {
		logger.Warnf("Could not add constraint fk_topups_user: %v", err)
	}

	// Transaction foreign keys
	if err := db.Exec(`
		ALTER TABLE transactions 
		DROP CONSTRAINT IF EXISTS fk_transactions_user;
	`).Error; err != nil {
		logger.Warnf("Could not drop constraint fk_transactions_user: %v", err)
	}

	if err := db.Exec(`
		ALTER TABLE transactions 
		ADD CONSTRAINT fk_transactions_user 
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
	`).Error; err != nil {
		logger.Warnf("Could not add constraint fk_transactions_user: %v", err)
	}

	// ShopPost foreign keys
	if err := db.Exec(`
		ALTER TABLE shop_posts 
		DROP CONSTRAINT IF EXISTS fk_shop_posts_coffee_shop;
	`).Error; err != nil {
		logger.Warnf("Could not drop constraint fk_shop_posts_coffee_shop: %v", err)
	}

	if err := db.Exec(`
		ALTER TABLE shop_posts 
		ADD CONSTRAINT fk_shop_posts_coffee_shop 
		FOREIGN KEY (coffee_shop_id) REFERENCES coffee_shops(id) ON DELETE CASCADE;
	`).Error; err != nil {
		logger.Warnf("Could not add constraint fk_shop_posts_coffee_shop: %v", err)
	}

	// InternalPost foreign keys
	if err := db.Exec(`
		ALTER TABLE internal_posts 
		DROP CONSTRAINT IF EXISTS fk_internal_posts_created_by;
	`).Error; err != nil {
		logger.Warnf("Could not drop constraint fk_internal_posts_created_by: %v", err)
	}

	if err := db.Exec(`
		ALTER TABLE internal_posts 
		ADD CONSTRAINT fk_internal_posts_created_by 
		FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE;
	`).Error; err != nil {
		logger.Warnf("Could not add constraint fk_internal_posts_created_by: %v", err)
	}

	logger.Info("Foreign key constraints added successfully")
	return nil
}
