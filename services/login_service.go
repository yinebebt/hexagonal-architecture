package services

type LoginService interface {
	Login(username string, password string) bool
}

type loginService struct {
	authorizedUsername string
	authorizedPassword string
}

func NewLoginService() LoginService {
	return &loginService{
		authorizedUsername: "yinebeb",
		authorizedPassword: "silenat",
	}
}

// Login consider this a DB where pair of PW and username stored.
func (service *loginService) Login(username string, password string) bool {
	return service.authorizedUsername == username &&
		service.authorizedPassword == password
}
