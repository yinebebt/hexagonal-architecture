package middleware

import (
	"github.com/Yinebeb-01/hexagonalarch/internal/core/service"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// AuthorizeJWT validates the token from the http request, returning a 401 if it's not valid
func AuthorizeJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		const BearerSchema = "Bearer "
		authHeader := c.GetHeader("Authorization")
		tokenString := authHeader[len(BearerSchema):]

		token, err := service.NewJWTService().ValidateToken(tokenString)

		if token.Valid {
			claims := token.Claims.(jwt.MapClaims)
			log.Println("Claims[Name]: ", claims["name"])
			log.Println("Claims[Admin]: ", claims["admin"])
			log.Println("Claims[Issuer]: ", claims["iss"])
			log.Println("Claims[IssuedAt]: ", claims["iat"])
			log.Println("Claims[ExpiresAt]: ", claims["exp"])
		} else {
			log.Println(err)
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}

// BasicAuth wrap gin.BasicAuth which returns a Basic HTTP Authorization middleware
func BasicAuth() gin.HandlerFunc {
	return gin.BasicAuth(gin.Accounts{
		"admin":     "admin1234",
		"test_user": "user1234",
	})
}
