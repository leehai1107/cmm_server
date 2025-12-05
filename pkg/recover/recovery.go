package recover

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/leehai1107/cmm_server/pkg/logger"
)

func RPanic(c *gin.Context) {
	// Create a logger with the request context
	log := logger.EnhanceWith(c.Request.Context())

	// Log before executing the handler
	log.Infof("Executing handler for %v %v", c.Request.Method, c.Request.URL.EscapedPath())

	// Wrap the entire handler chain in a recover block
	defer func() {
		if err := recover(); err != nil {
			// Log detailed information about the request
			log.Errorf("PANIC RECOVERED: method %v, path %v, err %v",
				c.Request.Method,
				c.Request.URL.EscapedPath(),
				err,
			)

			// Log request headers and body for debugging
			log.Infof("Request headers: %v", c.Request.Header)

			// Log handler information
			log.Infof("Handler name: %v", c.HandlerName())

			// Log all available keys in the context
			keys := c.Keys
			if keys != nil {
				log.Infof("Context keys: %v", keys)
			} else {
				log.Info("No context keys available")
			}

			// Print the stack trace
			log.Error("Stack trace:")
			debug.PrintStack()

			// Return a 500 error to the client
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Internal server error",
			})
		}
	}()

	// Execute the next handler in the chain
	// This is wrapped in a separate defer/recover to catch any panics that might occur
	// during the execution of c.Next()
	func() {
		defer func() {
			if r := recover(); r != nil {
				// Re-panic to be caught by the outer recover
				panic(r)
			}
		}()
		c.Next()
	}()

	// Log after executing the handler
	log.Infof("Finished executing handler for %v %v", c.Request.Method, c.Request.URL.EscapedPath())
}
