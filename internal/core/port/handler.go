package port

type VideoHandler interface {
	Save(ctx interface{})
	FindAll(ctx interface{})
	ShowAll(ctx interface{})
	Update(ctx interface{})
	Delete(ctx interface{})
}

type LoginHandler interface {
	Login(ctx interface{})
}
