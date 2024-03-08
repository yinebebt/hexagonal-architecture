package service

type LoginService interface {
	Login(username string, password string) bool
}

type login struct {
	authorizedUsername string
	authorizedPassword string
}

func NewLoginService() LoginService {
	return &login{
		authorizedUsername: "admin",
		authorizedPassword: "admin1234",
	}
}

// Login consider this a DB where pair of PW and username stored.
func (l *login) Login(username string, password string) bool {
	return l.authorizedUsername == username &&
		l.authorizedPassword == password
}
