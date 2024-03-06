package services

import (
	"gitlab.com/Yinebeb-01/hexagonalarch/entity"
	"gitlab.com/Yinebeb-01/hexagonalarch/repository"
)

type VideoService interface {
	Save(entity.Video) entity.Video
	Update(entity.Video)
	Delete(entity.Video)
	FindAll() []entity.Video
}

type videoService struct {
	videos repository.VideoRepository
}

// New is a constructor to initialize a videoService.
func New(vidRepo repository.VideoRepository) VideoService {
	return &videoService{videos: vidRepo}
}

// Save will add append video to Videos, and return the newly added video.
func (repoService *videoService) Save(video entity.Video) entity.Video {
	repoService.videos.Save(video)
	return video
}

func (repoService *videoService) Update(video entity.Video) {
	repoService.videos.Update(video)
}

func (repoService *videoService) Delete(video entity.Video) {
	repoService.videos.Delete(video)
}

func (repoService *videoService) FindAll() []entity.Video {
	return repoService.videos.FindAll()
}
