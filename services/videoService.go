package services

import (
	"gitlab.com/Yinebeb-01/simpleAPI/entity"
)

type VideoService interface {
	Save(entity.Video) entity.Video
	FindAll() []entity.Video
}

type videoSer struct {
	videos []entity.Video
}

// New is a constructor to iniialize a viddeosService.
func New() VideoService {
	return &videoSer{}
}

// Save will add append video to Videos, and return the newly added video.
func (service *videoSer) Save(video entity.Video) entity.Video {
	service.videos = append(service.videos, video)
	return video
}

// find will return slice of all videos saved.
func (service *videoSer) FindAll() []entity.Video {
	return service.videos
}
