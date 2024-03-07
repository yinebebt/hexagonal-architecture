package entity

import "time"

type Person struct {
	ID        uint64 `gorm:"primary_key;auto_increment" json:"id"`
	FirstName string `json:"firstname" binding:"required" gorm:"type:varchar(32)"`
	LastName  string `json:"lastname" binding:"required" gorm:"type:varchar(32)"`
	Age       int8   `json:"age" binding:"gte=10,lte=127"`
	Email     string `json:"email" validate:"required,email" gorm:"type:varchar(256)"`
}

type Video struct {
	ID          uint64    `gorm:"primary_key;auto_increment" json:"idd"`
	Title       string    `json:"title" binding:"min=3,max=10" validate:"is-cool" gorm:"type:varchar(10)"`
	Description string    `json:"description" binding:"max=25" gorm:"type:varchar(25)"`
	URL         string    `json:"url" binding:"required,url" gorm:"type:varchar(256); UNIQUE"`
	Director    Person    `json:"author" binding:"required" gorm:"foreignkey:PersonID"`
	PersonID    uint64    `json:"-"`
	CreatedAt   time.Time `json:"-" gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `json:"-" gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
