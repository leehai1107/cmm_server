package http

import (
	"github.com/leehai1107/cmm_server/service/cmm/usecase"
)

// IHandler defines all handler interfaces
type IHandler interface {
	IUserHandler
	IAdminHandler
	ICoffeeShopHandler
	IChatHandler
	IBookingHandler
	IMeetingRoomHandler
	IWalletHandler
	IVoucherHandler
	IPostHandler
}

// Handler implements all handler interfaces
type Handler struct {
	userUsecase        usecase.IUserUsecase
	adminUsecase       usecase.IAdminUsecase
	bookingUsecase     usecase.IBookingUsecase
	coffeeShopUsecase  usecase.ICoffeeShopUsecase
	meetingRoomUsecase usecase.IMeetingRoomUsecase
	walletUsecase      usecase.IWalletUsecase
	voucherUsecase     usecase.IVoucherUsecase
	postUsecase        usecase.IPostUsecase
}

func NewHandler(
	userUsecase usecase.IUserUsecase,
	adminUsecase usecase.IAdminUsecase,
	bookingUsecase usecase.IBookingUsecase,
	coffeeShopUsecase usecase.ICoffeeShopUsecase,
	meetingRoomUsecase usecase.IMeetingRoomUsecase,
	walletUsecase usecase.IWalletUsecase,
	voucherUsecase usecase.IVoucherUsecase,
	postUsecase usecase.IPostUsecase,
) IHandler {
	return &Handler{
		userUsecase:        userUsecase,
		adminUsecase:       adminUsecase,
		bookingUsecase:     bookingUsecase,
		coffeeShopUsecase:  coffeeShopUsecase,
		meetingRoomUsecase: meetingRoomUsecase,
		walletUsecase:      walletUsecase,
		voucherUsecase:     voucherUsecase,
		postUsecase:        postUsecase,
	}
}
