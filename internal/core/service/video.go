package service

import (
	"context"
	"fmt"

	"github.com/yinebebt/hexagonal-architecture/internal/core/entity"
	"github.com/yinebebt/hexagonal-architecture/internal/core/port"
)

type video struct {
	videoRepository port.VideoRepository
}

// New is a constructor to initialize a videoService.
func New(vidRepo port.VideoRepository) port.VideoService {
	return &video{videoRepository: vidRepo}
}

// Save will add a video to the repository and return the newly added video.
func (v *video) Save(ctx context.Context, video entity.Video) (entity.Video, error) {
	if err := v.validateVideo(video); err != nil {
		return entity.Video{}, err
	}
	if err := v.videoRepository.Save(ctx, &video); err != nil {
		return entity.Video{}, err
	}
	return video, nil
}

func (v *video) Update(ctx context.Context, video entity.Video) (entity.Video, error) {
	if video.ID == 0 {
		return entity.Video{}, fmt.Errorf("video ID is required for update")
	}
	if err := v.validateVideo(video); err != nil {
		return entity.Video{}, err
	}
	err := v.videoRepository.Update(ctx, &video)
	if err != nil {
		return entity.Video{}, err
	}
	return video, nil
}

func (v *video) Delete(ctx context.Context, id int64) error {
	return v.videoRepository.Delete(ctx, id)
}

func (v *video) FindByID(ctx context.Context, id int64) (entity.Video, error) {
	return v.videoRepository.FindByID(ctx, id)
}

func (v *video) FindAll(ctx context.Context) ([]entity.Video, error) {
	return v.videoRepository.FindAll(ctx)
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
