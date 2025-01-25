package service

import (
	"github.com/Yinebeb-01/hexagonalarch/internal/core/entity"
	"github.com/Yinebeb-01/hexagonalarch/internal/core/port"
)

type VideoService interface {
	Save(entity.Video) (entity.Video, error)
	Update(entity.Video) entity.Video
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

// Save will add append video to Videos, and return the newly added video.
func (v *video) Save(video entity.Video) (entity.Video, error) {
	err := v.videoRepository.Save(video)
	return video, err
}

func (v *video) Update(video entity.Video) entity.Video {
	v.videoRepository.Update(video)
	return video
}

func (v *video) Delete(id int64) error {
	return v.videoRepository.Delete(id)
}

func (v *video) FindAll() ([]entity.Video, error) {
	return v.videoRepository.FindAll()
}
