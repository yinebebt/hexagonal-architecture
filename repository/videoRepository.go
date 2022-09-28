package repository

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"gitlab.com/Yinebeb-01/simpleAPI/entity"
)

type VideoReposiory interface {
	Save(entity.Video)
	Update(entity.Video)
	Delete(entity.Video)
	FindAll() []entity.Video
	Close()
}

type Database struct {
	connection *gorm.DB
}

func NewVideoRepository() VideoReposiory {
	db, err := gorm.Open("sqlite3", "test.db") //creating new DB, with sqlite3-type
	if err != nil {
		panic("Failed to create a DB.")
	}
	db.AutoMigrate(&entity.Video{}, &entity.Person{})
	return &Database{
		connection: db,
	}
}

func (db *Database) Save(video entity.Video) {
	db.connection.Create(&video) //insert value to Database
}

func (db *Database) Update(video entity.Video) {
	db.connection.Save(&video)
}

func (db *Database) FindAll() []entity.Video {
	var videos []entity.Video
	db.connection.Set(`gorm:"auto_preload"`, true).Find(&videos) //set to fetch person struct/table too via the foreignkey
	return videos
}

func (db *Database) Delete(video entity.Video) {
	db.connection.Delete(video)
}

func (db *Database) Close() {
	err := db.connection.Close()
	if err != nil {
		panic("Failed to close database")
	}

}
