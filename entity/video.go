package entity

type person struct {
	FirstName string `json:"firstname" binding:"required"`
	LastName  string `json:"lastname" binding:"required"`
	Age       int8   `json:"age" binding:"gte=10,lte=127"`
	Email     string `json:"email" validate:"required,email"`
}

type Video struct {
	Title       string `json:"title" binding:"min=3,max=10" validate:"is-cool"`
	Description string `json:"description" binding:"max=25"`
	URL         string `json:"url" binding:"required,url"`
	Director    person `json:"author" binding:"required"`
}
