package controller

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/Yinebeb-01/hexagonalarch/internal/adapter/dto"
	"gitlab.com/Yinebeb-01/hexagonalarch/services"
)

type LoginController interface {
	Login(ctx *gin.Context) string
}

type loginController struct {
	loginService services.LoginService
	jWtService   services.JWTService
}

func NewLoginController(loginService services.LoginService,
	jWtService services.JWTService) LoginController {
	return &loginController{
		loginService: loginService,
		jWtService:   jWtService,
	}
}

func (controller *loginController) Login(ctx *gin.Context) string {
	var credentials dto.Credentials
	err := ctx.ShouldBind(&credentials)
	if err != nil {
		return ""
	}
	if controller.loginService.Login(credentials.Username, credentials.Password) {
		return controller.jWtService.GenerateToken(credentials.Username, true)
	}
	return ""
}
