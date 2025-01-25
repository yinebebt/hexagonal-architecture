package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

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
