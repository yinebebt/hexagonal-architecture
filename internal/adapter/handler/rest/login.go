package rest

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/Yinebeb-01/hexagonalarch/internal/adapter/dto"
	"gitlab.com/Yinebeb-01/hexagonalarch/internal/core/port"
	"gitlab.com/Yinebeb-01/hexagonalarch/internal/core/service"
	"net/http"
)

type login struct {
	loginService service.LoginService
	jWtService   service.JWTService
}

func InitLogin(loginService service.LoginService, jwtService service.JWTService) port.LoginHandler {
	return &login{
		loginService: loginService,
		jWtService:   jwtService,
	}
}

// Login
// @Summary      Login user
// @Description  login user description
// @Tags         User
// @Accept       json
// @Produce      json
// @Success      200  {object}  string
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /login [post]
func (l *login) Login(ctx *gin.Context) {
	var credentials dto.Credentials
	err := ctx.ShouldBind(&credentials)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid input"})
		return
	}
	if !l.loginService.Login(credentials.Username, credentials.Password) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "invalid credentials"})
		return
	}
	token := l.jWtService.GenerateToken(credentials.Username, true)

	if token == "" {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
