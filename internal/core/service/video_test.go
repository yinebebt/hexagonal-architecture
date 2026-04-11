package service

import (
	"context"
	"flag"
	"testing"

	"github.com/yinebebt/hexagonal-architecture/internal/adapter/repository"
	"github.com/yinebebt/hexagonal-architecture/internal/core/entity"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	TITLE       = "Video Title"
	DESCRIPTION = "Video Description"
	URL         = "https://youtu.be/JgW-i2QjgHQ"
)

var (
	dbType = flag.String("dbtype", "sqlite", "Type of the database (sqlite or postgres)")
	dsn    = flag.String("dsn", "test.db", "Data source name for the database")
)

// TestFindAll_WhenNoVideo tests the scenario when there are no videos in the database
func TestFindAll_WhenNoVideo(t *testing.T) {
	ctx := context.Background()
	videoRepository, err := repository.NewVideoRepository(*dbType, *dsn)
	require.NoError(t, err)
	videoSer := New(videoRepository)

	// Clean up before test
	defer func() {
		_ = videoRepository.Clean(ctx)
	}()

	// When: admin runs FindAll method
	videos, err := videoSer.FindAll(ctx)

	// Then: video should be null/empty
	assert.Nil(t, err)
	assert.Empty(t, videos, "should be empty when no videos exist")
}

// TestFindAll_WhenVideoExists tests the scenario when there is a video in the database
func TestFindAll_WhenVideoExists(t *testing.T) {
	ctx := context.Background()
	videoRepository, err := repository.NewVideoRepository(*dbType, *dsn)
	require.NoError(t, err)
	videoSer := New(videoRepository)

	// Clean up after test
	defer func() {
		_ = videoRepository.Clean(ctx)
	}()

	// Given: Admin posts some video
	video := entity.Video{
		Title:       "cool video",
		Description: "video description",
		URL:         "https://www.yoe.com/embed/96np1mk",
		Director: entity.Person{
			FirstName: "Abel",
			LastName:  "Yisak",
			Age:       25,
			Email:     "abel@gmail.com",
		},
	}
	_, err = videoSer.Save(ctx, video)
	assert.Nil(t, err)

	// When: admin runs FindAll method
	videos, err := videoSer.FindAll(ctx)

	// Then: video should be returned
	assert.Nil(t, err)
	assert.NotEmpty(t, videos, "should not be empty when video exists")
	assert.GreaterOrEqual(t, len(videos), 1, "should have at least one video")

	videoRes := videos[0]
	assert.Equal(t, video.Title, videoRes.Title, "title should match")
	assert.Equal(t, video.Description, videoRes.Description, "description should match")
	assert.Equal(t, video.URL, videoRes.URL, "URL should match")
}

func TestFindAll(t *testing.T) {
	ctx := context.Background()
	videoRepo, err := repository.NewVideoRepository(*dbType, *dsn)
	require.NoError(t, err)
	defer func() {
		_ = videoRepo.Clean(ctx)
	}()

	srv := New(videoRepo)

	_, err = srv.Save(ctx, getVideo())
	assert.NoError(t, err)

	videos, err := srv.FindAll(ctx)
	assert.NoError(t, err)
	assert.NotEmpty(t, videos)

	firstVideo := videos[0]
	assert.Equal(t, TITLE, firstVideo.Title)
	assert.Equal(t, DESCRIPTION, firstVideo.Description)
	assert.Equal(t, URL, firstVideo.URL)

	err = videoRepo.Delete(ctx, firstVideo.ID)
	assert.NoError(t, err)
}

func TestFindByID(t *testing.T) {
	ctx := context.Background()
	videoRepo, err := repository.NewVideoRepository(*dbType, *dsn)
	require.NoError(t, err)
	defer func() {
		_ = videoRepo.Clean(ctx)
	}()

	srv := New(videoRepo)

	// Save a video first
	savedVideo, err := srv.Save(ctx, getVideo())
	require.NoError(t, err)

	// Find the video by ID
	found, err := srv.FindByID(ctx, savedVideo.ID)
	assert.NoError(t, err)
	assert.Equal(t, savedVideo.ID, found.ID)
	assert.Equal(t, TITLE, found.Title)
	assert.Equal(t, DESCRIPTION, found.Description)
	assert.Equal(t, URL, found.URL)

	// Try to find a non-existent video
	_, err = srv.FindByID(ctx, 99999)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestUpdate(t *testing.T) {
	ctx := context.Background()
	videoRepo, err := repository.NewVideoRepository(*dbType, *dsn)
	require.NoError(t, err)
	defer func() {
		_ = videoRepo.Clean(ctx)
	}()

	srv := New(videoRepo)

	// Save a video first
	video := getVideo()
	savedVideo, err := srv.Save(ctx, video)
	assert.NoError(t, err)

	// Update the video
	savedVideo.Title = "Updated Title"
	savedVideo.Description = "Updated Description"
	updatedVideo, err := srv.Update(ctx, savedVideo)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Title", updatedVideo.Title)
	assert.Equal(t, "Updated Description", updatedVideo.Description)
}

func TestUpdate_InvalidID(t *testing.T) {
	ctx := context.Background()
	videoRepo, err := repository.NewVideoRepository(*dbType, *dsn)
	require.NoError(t, err)
	defer func() {
		_ = videoRepo.Clean(ctx)
	}()

	srv := New(videoRepo)

	// Try to update with ID = 0
	video := getVideo()
	video.ID = 0
	_, err = srv.Update(ctx, video)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ID is required")
}

func TestDelete(t *testing.T) {
	ctx := context.Background()
	videoRepo, err := repository.NewVideoRepository(*dbType, *dsn)
	require.NoError(t, err)
	defer func() {
		_ = videoRepo.Clean(ctx)
	}()

	srv := New(videoRepo)

	// Save a video first
	video := getVideo()
	savedVideo, err := srv.Save(ctx, video)
	assert.NoError(t, err)

	// Delete the video
	err = srv.Delete(ctx, savedVideo.ID)
	assert.NoError(t, err)

	// Verify it's deleted
	videos, err := srv.FindAll(ctx)
	assert.NoError(t, err)
	assert.Empty(t, videos)
}

func getVideo() entity.Video {
	return entity.Video{
		Title:       TITLE,
		Description: DESCRIPTION,
		URL:         URL,
	}
}
