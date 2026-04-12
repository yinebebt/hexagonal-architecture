package entity

import "time"

type Person struct {
	ID        uint64 `json:"id"`
	FirstName string `json:"firstname" binding:"required"`
	LastName  string `json:"lastname" binding:"required"`
	Age       int8   `json:"age" binding:"gte=10,lte=127"`
	Email     string `json:"email" validate:"required,email"`
}

type Video struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title" binding:"min=3,max=100"`
	Description string    `json:"description" binding:"max=500"`
	URL         string    `json:"url" binding:"required,url" example:"https://google.com/xyz-video"`
	Director    Person    `json:"author" binding:"required"`
	PersonID    uint64    `json:"-"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}
