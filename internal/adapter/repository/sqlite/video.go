// Package sqlite is sqlite implementation of video reository, it uses gorm
package sqlite

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

func NewVideoRepository(dsn string) port.VideoRepository {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to create a DB.")
	}
	err = db.AutoMigrate(&entity.Video{}, &entity.Person{})
	if err != nil {
		return &Database{}
	}
	return &Database{
		connection: db,
	}
}

func (db *Database) Save(video entity.Video) error {
	if err := db.connection.Create(&video).Error; err != nil {
		return fmt.Errorf("failed to save video: %v", err)
	}
	return nil
}

func (db *Database) Update(video entity.Video) error {
	if err := db.connection.Save(&video).Error; err != nil {
		return fmt.Errorf("failed to update video: %v", err)
	}
	return nil
}

func (db *Database) FindAll() ([]entity.Video, error) {
	var videos []entity.Video
	//set to fetch person object too via the foreign key
	if err := db.connection.Set(`gorm:"auto_preload"`, true).Find(&videos).Error; err != nil {
		return nil, fmt.Errorf("failed to get products: %v", err)
	}
	return videos, nil
}

// fixme: query via unique id
func (db *Database) Delete(video entity.Video) error {
	if err := db.connection.Delete(video, fmt.Sprintf("title='%v'", video.Title)).Error; err != nil {
		return fmt.Errorf("failed to delete video: %v", err)
	}
	return nil
}

func (db *Database) Clean() error {
	if err := db.connection.Exec("drop table if exists videos").Error; err != nil {
		return fmt.Errorf("faile to drop table: %v", err)
	}
	return nil
}
