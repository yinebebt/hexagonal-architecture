package service

import (
	"github.com/Yinebeb-01/hexagonalarch/internal/core/entity"
	"github.com/Yinebeb-01/hexagonalarch/internal/core/port"
)

type VideoService interface {
	Save(entity.Video) (entity.Video, error)
	Update(entity.Video) entity.Video
	Delete(entity.Video)
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

// fixme: return from result
func (v *video) Update(video entity.Video) entity.Video {
	v.videoRepository.Update(video)
	return video
}

func (v *video) Delete(video entity.Video) {
	v.videoRepository.Delete(video)
}

func (v *video) FindAll() ([]entity.Video, error) {
	return v.videoRepository.FindAll()
}
