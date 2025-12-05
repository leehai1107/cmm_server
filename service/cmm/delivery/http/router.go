package http

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/leehai1107/cmm_server/pkg/apiwrapper"
)

type Router interface {
	Register(routerGroup gin.IRouter)
}

type routerImpl struct {
	handler IHandler
}

func NewRouter(
	handler IHandler,
) Router {
	return &routerImpl{
		handler: handler,
	}
}

func (p *routerImpl) Register(r gin.IRouter) {

	//routes for apis
	api := r.Group("api/v1")
	{
		api.GET("/ping", func(c *gin.Context) {
			apiwrapper.SendSuccess(c, time.Now())
		})
	}

	// Admin routes
	adminApi := api.Group("admin")
	{
		adminApi.POST("/create-account", p.handler.CreateAccount)
	}

	// User routes
	userApi := api.Group("user")
	{
		userApi.POST("/login", p.handler.Login)
		userApi.POST("/register", p.handler.Register)
	}

	// Coffee Shop routes
	coffeeShopApi := api.Group("coffee-shop")
	{
		coffeeShopApi.GET("/all", p.handler.GetAllCoffeeShops)
		coffeeShopApi.GET("/:id", p.handler.GetCoffeeShop)
		coffeeShopApi.GET("/commission/:shop_id", p.handler.GetCommissionRate)

		// Protected routes (require authentication)
		coffeeShopApi.POST("/create", p.handler.CreateCoffeeShop)
		coffeeShopApi.GET("/my-shops", p.handler.GetMyCoffeeShops)
		coffeeShopApi.PUT("/update", p.handler.UpdateCoffeeShop)
		coffeeShopApi.DELETE("/:id", p.handler.DeleteCoffeeShop)
		coffeeShopApi.POST("/commission/set", p.handler.SetCommissionRate)
	}

	// Meeting Room routes
	meetingRoomApi := api.Group("meeting-room")
	{
		meetingRoomApi.GET("/:id", p.handler.GetMeetingRoom)
		meetingRoomApi.GET("/shop/:shop_id", p.handler.GetMeetingRoomsByCoffeeShop)
		meetingRoomApi.GET("/shop/:shop_id/available", p.handler.GetAvailableMeetingRooms)

		// Protected routes
		meetingRoomApi.POST("/create", p.handler.CreateMeetingRoom)
		meetingRoomApi.PUT("/update", p.handler.UpdateMeetingRoom)
		meetingRoomApi.DELETE("/:id", p.handler.DeleteMeetingRoom)
	}

	// Booking routes
	bookingApi := api.Group("booking")
	{
		bookingApi.POST("/create", p.handler.CreateBooking)
		bookingApi.GET("/:id", p.handler.GetBooking)
		bookingApi.GET("/my-bookings", p.handler.GetMyBookings)
		bookingApi.POST("/:id/cancel", p.handler.CancelBooking)
		bookingApi.GET("/room/:room_id", p.handler.GetRoomBookings)
	}

	// Wallet routes
	walletApi := api.Group("wallet")
	{
		walletApi.GET("", p.handler.GetWallet)
		walletApi.POST("/topup", p.handler.CreateTopup)
		walletApi.POST("/topup/confirm", p.handler.ConfirmTopup) // Admin only
		walletApi.GET("/topup/history", p.handler.GetTopupHistory)
		walletApi.GET("/transactions", p.handler.GetTransactionHistory)
	}

	// Voucher routes
	voucherApi := api.Group("voucher")
	{
		voucherApi.GET("/valid", p.handler.GetValidVouchers)
		voucherApi.GET("/:id", p.handler.GetVoucher)
		voucherApi.POST("/apply", p.handler.ApplyVoucher)

		// Admin routes
		voucherApi.POST("/create", p.handler.CreateVoucher)
		voucherApi.GET("/all", p.handler.GetAllVouchers)
		voucherApi.PUT("/update", p.handler.UpdateVoucher)
		voucherApi.DELETE("/:id", p.handler.DeleteVoucher)
	}

	// Post routes
	postApi := api.Group("post")
	{
		// Shop posts
		shopPostApi := postApi.Group("/shop")
		{
			shopPostApi.GET("/all", p.handler.GetAllShopPosts)
			shopPostApi.GET("/:id", p.handler.GetShopPost)
			shopPostApi.GET("/coffee-shop/:shop_id", p.handler.GetShopPostsByCoffeeShop)

			// Protected routes
			shopPostApi.POST("/create", p.handler.CreateShopPost)
			shopPostApi.PUT("/update", p.handler.UpdateShopPost)
			shopPostApi.DELETE("/:id", p.handler.DeleteShopPost)
		}

		// Internal posts
		internalPostApi := postApi.Group("/internal")
		{
			internalPostApi.GET("/:id", p.handler.GetInternalPost)
			internalPostApi.GET("/my-posts", p.handler.GetMyInternalPosts)

			// Admin routes
			internalPostApi.POST("/create", p.handler.CreateInternalPost)
			internalPostApi.GET("/all", p.handler.GetAllInternalPosts)
			internalPostApi.PUT("/update", p.handler.UpdateInternalPost)
			internalPostApi.DELETE("/:id", p.handler.DeleteInternalPost)
		}
	}

	// WebSocket chat route
	chatApi := api.Group("chat")
	{
		chatApi.GET("/ws/:roomId", func(c *gin.Context) {
			defer func() {
				if r := recover(); r != nil {
					c.AbortWithStatusJSON(500, gin.H{"error": "Internal Server Error"})
				}
			}()
			p.handler.ServeWS(c)
		})
	}

}
