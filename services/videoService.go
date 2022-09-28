package services

import (
	"gitlab.com/Yinebeb-01/simpleAPI/entity"
	"gitlab.com/Yinebeb-01/simpleAPI/repository"
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

// New is a constructor to iniialize a viddeosService.
func New(vidrepo repository.VideoReposiory) VideoService {
	return &videoSer{videos: vidrepo}
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

// find will return slice of all videos saved.
func (repoService *videoSer) FindAll() []entity.Video {
	videos := repoService.videos.FindAll()
	return videos
}
