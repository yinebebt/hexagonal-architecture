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
} //consider this a DB where pair of PW and username stored, the credential struct is an instance/trial for this key/value pair.

func (service *loginService) Login(username string, password string) bool {
	return service.authorizedUsername == username &&
		service.authorizedPassword == password
}
