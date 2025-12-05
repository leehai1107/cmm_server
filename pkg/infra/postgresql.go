package infra

import (
	"fmt"
	"log"
	"time"

	"github.com/leehai1107/cmm_server/pkg/config"
	"github.com/leehai1107/cmm_server/pkg/constant"
	"github.com/leehai1107/cmm_server/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbSingleton *gorm.DB

func InitPostgresql() {
	// Check if database is already initialized
	if dbSingleton != nil {
		logger.Info("Database already initialized")
		return
	}

	logger.Info("Initializing PostgreSQL database connection")
	dbCfg := config.DBConfig()

	// Log database configuration (without password)
	logger.Infof("Database config: host=%s, port=%s, user=%s, dbname=%s",
		dbCfg.PgHost,
		dbCfg.PgPort,
		dbCfg.PgUser,
		dbCfg.PgDatabase)

	// Validate database configuration
	if dbCfg.PgHost == "" {
		err := fmt.Errorf("database host is empty")
		logger.Error(err.Error())
		panic(err)
	}
	if dbCfg.PgPort == "" {
		err := fmt.Errorf("database port is empty")
		logger.Error(err.Error())
		panic(err)
	}
	if dbCfg.PgUser == "" {
		err := fmt.Errorf("database user is empty")
		logger.Error(err.Error())
		panic(err)
	}
	if dbCfg.PgDatabase == "" {
		err := fmt.Errorf("database name is empty")
		logger.Error(err.Error())
		panic(err)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Ho_Chi_Minh",
		dbCfg.PgHost,
		dbCfg.PgUser,
		dbCfg.PgPassword,
		dbCfg.PgDatabase,
		dbCfg.PgPort)

	logger.Infof("Connecting to PostgreSQL with DSN: host=%s user=%s dbname=%s port=%s",
		dbCfg.PgHost, dbCfg.PgUser, dbCfg.PgDatabase, dbCfg.PgPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{FullSaveAssociations: true})
	if err != nil {
		logger.Errorf("Failed to connect to database: %s", err.Error())
		panic(fmt.Sprintf("Unable to instantiate database: %s", err.Error()))
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Errorf("Failed to get sql.DB from gorm.DB: %s", err.Error())
		panic(fmt.Sprintf("Unable to get sql.DB from gorm.DB: %s", err.Error()))
	}

	if dbCfg.PgPoolSize > 0 {
		sqlDB.SetMaxOpenConns(dbCfg.PgPoolSize)
		logger.Infof("Set max open connections to %d", dbCfg.PgPoolSize)
	}

	if dbCfg.PgIdleConnTimeout > 0 {
		sqlDB.SetConnMaxIdleTime(time.Duration(dbCfg.PgIdleConnTimeout) * time.Second)
		logger.Infof("Set connection max idle time to %d seconds", dbCfg.PgIdleConnTimeout)
	}

	if dbCfg.PgMaxConnAge > 0 {
		sqlDB.SetConnMaxLifetime(time.Duration(dbCfg.PgMaxConnAge) * time.Second)
		logger.Infof("Set connection max lifetime to %d seconds", dbCfg.PgMaxConnAge)
	}

	if config.ServerConfig().ENV != constant.ProductionEnv {
		db = db.Debug()
		logger.Info("Enabled GORM debug mode")
	}

	// Test the connection
	err = sqlDB.Ping()
	if err != nil {
		logger.Errorf("Failed to ping database: %s", err.Error())
		panic(fmt.Sprintf("Database connection test failed: %s", err.Error()))
	}

	// Check if required tables exist
	var tableCount int64
	result := db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public'").Count(&tableCount)
	if result.Error != nil {
		logger.Errorf("Failed to check tables: %s", result.Error.Error())
	} else {
		logger.Infof("Database has %d tables in public schema", tableCount)

		// Check for specific tables
		var userTableExists bool
		db.Raw("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'users')").Scan(&userTableExists)
		logger.Infof("Users table exists: %v", userTableExists)

		if !userTableExists {
			logger.Warn("Users table does not exist. This may cause issues with user-related operations.")
		}

		var walletTableExists bool
		db.Raw("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'wallets')").Scan(&walletTableExists)
		logger.Infof("Wallets table exists: %v", walletTableExists)

		if !walletTableExists {
			logger.Warn("Wallets table does not exist. This may cause issues with wallet-related operations.")
		}
	}

	logger.Info("Successfully connected to PostgreSQL database")
	dbSingleton = db

	// Run migrations if tables don't exist
	if tableCount == 0 {
		logger.Info("No tables found. Running migrations...")
		if err := RunMigrations(db); err != nil {
			logger.Errorf("Failed to run migrations: %s", err.Error())
			panic(fmt.Sprintf("Database migration failed: %s", err.Error()))
		}
		logger.Info("Migrations completed successfully")
	}
}

func ClosePostgresql() error {
	sqlDB, err := dbSingleton.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func GetDB() *gorm.DB {
	if dbSingleton == nil {
		log.Printf("Connection to Database Postgres is not setup")
		// Initialize the database if it's not already initialized
		InitPostgresql()

		// If it's still nil after initialization, log a more severe error
		if dbSingleton == nil {
			logger.Error("Failed to initialize database connection")
			// Panic here to prevent the application from starting with a nil DB
			panic("Failed to initialize database connection. Check your database configuration.")
		}
	}

	return dbSingleton
}

// BeginTransaction start an Transaction, require defer ReleaseTransaction instantly
func BeginTransaction() *gorm.DB {
	return dbSingleton.Begin()
}

func ReleaseTransaction(tx *gorm.DB) {
	tx.Commit()
}
