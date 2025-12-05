package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/leehai1107/cmm_server/cmd/banner"
	"github.com/leehai1107/cmm_server/di/apifx"
	"github.com/leehai1107/cmm_server/di/dbfx"
	"github.com/leehai1107/cmm_server/pkg/config"
	"github.com/leehai1107/cmm_server/pkg/errors"
	"github.com/leehai1107/cmm_server/pkg/graceful"
	"github.com/leehai1107/cmm_server/pkg/infra"
	"github.com/leehai1107/cmm_server/pkg/logger"
	"github.com/leehai1107/cmm_server/pkg/middleware/cors"
	"github.com/leehai1107/cmm_server/pkg/recover"
	"github.com/leehai1107/cmm_server/pkg/swagger"
	"github.com/leehai1107/cmm_server/pkg/utils/ginbuilder"
	"github.com/leehai1107/cmm_server/pkg/utils/timeutils"
	"github.com/leehai1107/cmm_server/pkg/websocket"
	"github.com/leehai1107/cmm_server/service/cmm/delivery/http"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Command of Internal Service",
	Long:  "CLI used to manage internal apis, datas when users access.",
	Run: func(_ *cobra.Command, _ []string) {
		NewServer().Run()
	},
	Version: "1.0.0",
}

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Run() {
	app := fx.New(
		fx.Invoke(config.InitConfig),
		fx.Invoke(initLogger),
		fx.Invoke(errors.Initialize),
		fx.Invoke(timeutils.Init),
		fx.Invoke(infra.InitPostgresql),
		fx.Invoke(websocket.InitHub),
		//... add module here
		dbfx.Module,
		apifx.Module,
		fx.Provide(provideGinEngine),
		fx.Invoke(
			registerService,
			registerSwaggerHandler),
		fx.Invoke(startServer),
		fx.Invoke(banner.Print),
	)
	logger.Info("Server started!")
	app.Run()
}

func provideGinEngine() *gin.Engine {
	return ginbuilder.BaseBuilder().Build()
}

func registerService(
	g *gin.Engine,
	router http.Router,
) {
	// Ensure router is not nil
	if router == nil {
		logger.Error("Router is nil in registerService")
		panic("Router is nil. Cannot register routes.")
	}

	// Ensure gin engine is not nil
	if g == nil {
		logger.Error("Gin engine is nil in registerService")
		panic("Gin engine is nil. Cannot register routes.")
	}

	logger.Info("Registering routes...")
	internal := g.Group("/internal")
	internal.Use(
		recover.RPanic,
		cors.CorsCfg(config.ServerConfig().CorsProduction))

	// Log middleware registration
	logger.Info("Middleware registered for /internal routes")

	// Register routes
	router.Register(internal)

	// Log successful registration
	logger.Info("Routes registered successfully")

	// Log all registered routes for debugging
	for _, route := range g.Routes() {
		logger.Infof("Registered route: %s %s", route.Method, route.Path)
	}
}

func registerSwaggerHandler(g *gin.Engine) {
	swaggerAPI := g.Group("/internal/swagger")
	swag := swagger.NewSwagger()
	swaggerAPI.Use(swag.SwaggerHandler(config.ServerConfig().Production))
	swag.Register(swaggerAPI)
}

func startServer(lifecycle fx.Lifecycle, g *gin.Engine) {
	gracefulService := graceful.NewService(graceful.WithStopTimeout(time.Second), graceful.WithWaitTime(time.Second))

	gracefulService.Register(g)
	lifecycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				port := os.Getenv("PORT")
				if port == "" {
					port = fmt.Sprintf("%d", config.ServerConfig().HTTPPort)
				}
				logger.Info("run on port:", port)
				go gracefulService.StartServer(g, port)
				return nil
			},
			OnStop: func(context.Context) error {
				gracefulService.Close()
				infra.ClosePostgresql() // nolint
				return nil
			},
		},
	)
}
func initLogger() {
	logger.Initialize(config.ServerConfig().Logger)
}
