package middlewares

import "github.com/gin-gonic/gin"

// BasicAuth wrap gin.BasicAuth which returns a Basic HTTP Authorization middleware
func BasicAuth() gin.HandlerFunc {
	return gin.BasicAuth(gin.Accounts{
		"yinebeb": "silenat", //username: yinebeb ; password: silenat
	})
}
