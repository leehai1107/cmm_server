package cors

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/leehai1107/cmm_server/pkg/config"
	"github.com/leehai1107/cmm_server/pkg/logger"
)

func CorsCfg(production bool) gin.HandlerFunc {
	var corsCfg gin.HandlerFunc
	msg := fmt.Sprintf("CORS config in production is %v ", production)

	if production {
		msg = msg + "--> load production config!"
		corsCfg = cors.New(
			cors.Config{
				AllowOrigins:     buildAllowOrigins(),
				AllowMethods:     buildAllowMethods(),
				AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
				ExposeHeaders:    []string{"Content-Length"},
				AllowCredentials: true,
				MaxAge:           12 * time.Hour,
			},
		)
	} else {
		msg = msg + "--> load default config!"
		corsCfg = cors.Default()
	}

	logger.Info(msg)

	// Add a wrapper to catch any panics in the CORS middleware
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logger.Errorf("PANIC in CORS middleware: %v", r)
				c.AbortWithStatusJSON(500, gin.H{"error": "Internal Server Error"})
			}
		}()

		// If corsCfg is nil, use a default handler that does nothing
		if corsCfg == nil {
			logger.Warn("CORS config is nil, using default pass-through handler")
			c.Next()
			return
		}

		corsCfg(c)
	}
}

// buildAllowOrigins extracts the values from CorsCfg struct and returns them as a slice of strings.
func buildAllowOrigins() []string {
	// Get the CorsCfg struct from the CorsConfig function
	corsConfig := config.CorsConfig()

	// Initialize a slice to store the origins
	allowOrigins := make([]string, 0)

	// Loop through the CorsCfg struct
	// and add the values into the allowOrigins slice
	allowOrigins = append(allowOrigins, corsConfig.Google)
	allowOrigins = append(allowOrigins, corsConfig.Facebook)

	return allowOrigins
}

func buildAllowMethods() []string {
	return []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
}
