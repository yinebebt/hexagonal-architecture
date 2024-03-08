package gorm

import (
	"fmt"
	"github.com/Yinebeb-01/hexagonalarch/internal/core/entity"
	"github.com/Yinebeb-01/hexagonalarch/internal/core/port"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	connection *gorm.DB
}

func NewVideoRepository() port.VideoRepository {
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
	//set to fetch person object too via the foreign key
	db.connection.Set(`gorm:"auto_preload"`, true).Find(&videos)
	return videos
}

// fixme: query via unique id
func (db *Database) Delete(video entity.Video) {
	db.connection.Delete(video, fmt.Sprintf("title='%v'", video.Title))
}
