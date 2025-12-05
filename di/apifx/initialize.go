package apifx

import (
	"github.com/leehai1107/cmm_server/service/cmm/delivery/http"
	"github.com/leehai1107/cmm_server/service/cmm/repository"
	"github.com/leehai1107/cmm_server/service/cmm/usecase"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

var Module = fx.Provide(
	provideRouter,
	provideHandler,

	// Repositories
	provideUserRepo,
	provideBookingRepo,
	provideCoffeeShopRepo,
	provideMeetingRoomRepo,
	provideWalletRepo,
	provideVoucherRepo,
	providePostRepo,
	provideTransactionRepo,

	// Usecases
	provideUserUsecase,
	provideAdminUsecase,
	provideBookingUsecase,
	provideCoffeeShopUsecase,
	provideMeetingRoomUsecase,
	provideWalletUsecase,
	provideVoucherUsecase,
	providePostUsecase,
)

func provideRouter(handler http.IHandler) http.Router {
	return http.NewRouter(handler)
}

func provideHandler(
	userUsecase usecase.IUserUsecase,
	adminUsecase usecase.IAdminUsecase,
	bookingUsecase usecase.IBookingUsecase,
	coffeeShopUsecase usecase.ICoffeeShopUsecase,
	meetingRoomUsecase usecase.IMeetingRoomUsecase,
	walletUsecase usecase.IWalletUsecase,
	voucherUsecase usecase.IVoucherUsecase,
	postUsecase usecase.IPostUsecase,
) http.IHandler {
	handler := http.NewHandler(
		userUsecase,
		adminUsecase,
		bookingUsecase,
		coffeeShopUsecase,
		meetingRoomUsecase,
		walletUsecase,
		voucherUsecase,
		postUsecase,
	)
	return handler
}

// Repository providers
func provideUserRepo(db *gorm.DB) repository.IUserRepo {
	return repository.NewUserRepo(db)
}

func provideBookingRepo(db *gorm.DB) repository.IBookingRepo {
	return repository.NewBookingRepo(db)
}

func provideCoffeeShopRepo(db *gorm.DB) repository.ICoffeeShopRepo {
	return repository.NewCoffeeShopRepo(db)
}

func provideMeetingRoomRepo(db *gorm.DB) repository.IMeetingRoomRepo {
	return repository.NewMeetingRoomRepo(db)
}

func provideWalletRepo(db *gorm.DB) repository.IWalletRepo {
	return repository.NewWalletRepo(db)
}

func provideVoucherRepo(db *gorm.DB) repository.IVoucherRepo {
	return repository.NewVoucherRepo(db)
}

func providePostRepo(db *gorm.DB) repository.IPostRepo {
	return repository.NewPostRepo(db)
}

func provideTransactionRepo(db *gorm.DB) repository.ITransactionRepo {
	return repository.NewTransactionRepo(db)
}

// Usecase providers
func provideUserUsecase(repo repository.IUserRepo) usecase.IUserUsecase {
	return usecase.NewUserUsecase(repo)
}

func provideAdminUsecase(repo repository.IUserRepo) usecase.IAdminUsecase {
	return usecase.NewAdminUsecase(repo)
}

func provideBookingUsecase(
	bookingRepo repository.IBookingRepo,
	meetingRoomRepo repository.IMeetingRoomRepo,
	walletRepo repository.IWalletRepo,
	voucherRepo repository.IVoucherRepo,
	transactionRepo repository.ITransactionRepo,
) usecase.IBookingUsecase {
	return usecase.NewBookingUsecase(bookingRepo, meetingRoomRepo, walletRepo, voucherRepo, transactionRepo)
}

func provideCoffeeShopUsecase(coffeeShopRepo repository.ICoffeeShopRepo) usecase.ICoffeeShopUsecase {
	return usecase.NewCoffeeShopUsecase(coffeeShopRepo)
}

func provideMeetingRoomUsecase(
	meetingRoomRepo repository.IMeetingRoomRepo,
	coffeeShopRepo repository.ICoffeeShopRepo,
) usecase.IMeetingRoomUsecase {
	return usecase.NewMeetingRoomUsecase(meetingRoomRepo, coffeeShopRepo)
}

func provideWalletUsecase(
	walletRepo repository.IWalletRepo,
	transactionRepo repository.ITransactionRepo,
) usecase.IWalletUsecase {
	return usecase.NewWalletUsecase(walletRepo, transactionRepo)
}

func provideVoucherUsecase(voucherRepo repository.IVoucherRepo) usecase.IVoucherUsecase {
	return usecase.NewVoucherUsecase(voucherRepo)
}

func providePostUsecase(
	postRepo repository.IPostRepo,
	coffeeShopRepo repository.ICoffeeShopRepo,
	userRepo repository.IUserRepo,
) usecase.IPostUsecase {
	return usecase.NewPostUsecase(postRepo, coffeeShopRepo, userRepo)
}
