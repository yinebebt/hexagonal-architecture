package service

import (
	"fmt"

	"github.com/yinebebt/hexagonal-architecture/internal/core/entity"
	"github.com/yinebebt/hexagonal-architecture/internal/core/port"
)

type VideoService interface {
	Save(entity.Video) (entity.Video, error)
	Update(entity.Video) (entity.Video, error)
	Delete(id int64) error
	FindAll() ([]entity.Video, error)
}

type video struct {
	videoRepository port.VideoRepository
}

// New is a constructor to initialize a videoService.
func New(vidRepo port.VideoRepository) VideoService {
	return &video{videoRepository: vidRepo}
}

// Save will add a video to the repository and return the newly added video.
func (v *video) Save(video entity.Video) (entity.Video, error) {
	if err := v.validateVideo(video); err != nil {
		return entity.Video{}, err
	}
	if err := v.videoRepository.Save(&video); err != nil {
		return entity.Video{}, err
	}
	return video, nil
}

func (v *video) Update(video entity.Video) (entity.Video, error) {
	if video.ID == 0 {
		return entity.Video{}, fmt.Errorf("video ID is required for update")
	}
	if err := v.validateVideo(video); err != nil {
		return entity.Video{}, err
	}
	err := v.videoRepository.Update(&video)
	if err != nil {
		return entity.Video{}, err
	}
	return video, nil
}

func (v *video) Delete(id int64) error {
	return v.videoRepository.Delete(id)
}

func (v *video) FindAll() ([]entity.Video, error) {
	return v.videoRepository.FindAll()
}

// validateVideo performs basic validation on video entity
func (v *video) validateVideo(video entity.Video) error {
	if video.Title == "" {
		return fmt.Errorf("title is required")
	}
	if video.URL == "" {
		return fmt.Errorf("URL is required")
	}
	return nil
}
