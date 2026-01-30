package rest

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// Logger returns a Gin middleware handler for HTTP request logging
func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(params gin.LogFormatterParams) string {
		return fmt.Sprintf("%s [%s] |%s| %d %s `%s`",
			params.ClientIP,
			params.TimeStamp,
			params.Method,
			params.StatusCode,
			params.Latency,
			params.Path,
		)
	})
}
