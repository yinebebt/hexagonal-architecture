package repository

import (
	"fmt"
	"gitlab.com/Yinebeb-01/simpleapi/entity"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type VideoReposiory interface {
	Save(entity.Video)
	Update(entity.Video)
	Delete(entity.Video)
	FindAll() []entity.Video
}

type Database struct {
	connection *gorm.DB
}

func NewVideoRepository() VideoReposiory {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to create a DB.")
	}
	_ = db.AutoMigrate(&entity.Video{}, &entity.Person{})
	if err != nil {
		return &Database{}
	}
	return &Database{
		connection: db,
	}
}

func (db *Database) Save(video entity.Video) {
	db.connection.Create(&video)
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
	db.connection.Delete(video, fmt.Sprintf("title='%v'", video.Title))
}
