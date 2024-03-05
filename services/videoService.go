package services

import (
	"gitlab.com/Yinebeb-01/simpleapi/entity"
	"gitlab.com/Yinebeb-01/simpleapi/repository"
)

type VideoService interface {
	Save(entity.Video) entity.Video
	Update(entity.Video)
	Delete(entity.Video)
	FindAll() []entity.Video
}

type videoSer struct {
	videos repository.VideoReposiory
}

// New is a constructor to initialize a videoService.
func New(vidRepo repository.VideoReposiory) VideoService {
	return &videoSer{videos: vidRepo}
}

// Save will add append video to Videos, and return the newly added video.
func (repoService *videoSer) Save(video entity.Video) entity.Video {
	repoService.videos.Save(video)
	return video
}

func (repoService *videoSer) Update(video entity.Video) {
	repoService.videos.Update(video)
}

func (repoService *videoSer) Delete(video entity.Video) {
	repoService.videos.Delete(video)
}

func (repoService *videoSer) FindAll() []entity.Video {
	return repoService.videos.FindAll()
}
