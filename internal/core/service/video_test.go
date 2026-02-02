package service

import (
	"flag"
	"testing"

	"github.com/yinebebt/hexagonal-architecture/internal/adapter/repository"
	"github.com/yinebebt/hexagonal-architecture/internal/core/entity"

	"github.com/stretchr/testify/assert"
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
	videoRepository, err := repository.NewVideoRepository(*dbType, *dsn)
	assert.NoError(t, err)
	videoSer := New(videoRepository)

	// Clean up before test
	defer func() {
		_ = videoRepository.Clean()
	}()

	// When: admin runs FindAll method
	videos, err := videoSer.FindAll()

	// Then: video should be null/empty
	assert.Nil(t, err)
	assert.Empty(t, videos, "should be empty when no videos exist")
}

// TestFindAll_WhenVideoExists tests the scenario when there is a video in the database
func TestFindAll_WhenVideoExists(t *testing.T) {
	videoRepository, err := repository.NewVideoRepository(*dbType, *dsn)
	assert.NoError(t, err)
	videoSer := New(videoRepository)

	// Clean up after test
	defer func() {
		_ = videoRepository.Clean()
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
	_, err = videoSer.Save(video)
	assert.Nil(t, err)

	// When: admin runs FindAll method
	videos, err := videoSer.FindAll()

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
	videoRepo, err := repository.NewVideoRepository(*dbType, *dsn)
	assert.NoError(t, err)
	defer func() {
		_ = videoRepo.Clean()
	}()

	srv := New(videoRepo)

	_, err = srv.Save(getVideo())
	assert.NoError(t, err)

	videos, err := srv.FindAll()
	assert.NoError(t, err)
	assert.NotEmpty(t, videos)

	firstVideo := videos[0]
	assert.Equal(t, TITLE, firstVideo.Title)
	assert.Equal(t, DESCRIPTION, firstVideo.Description)
	assert.Equal(t, URL, firstVideo.URL)

	err = videoRepo.Delete(firstVideo.ID)
	assert.NoError(t, err)
}

func TestUpdate(t *testing.T) {
	videoRepo, err := repository.NewVideoRepository(*dbType, *dsn)
	assert.NoError(t, err)
	defer func() {
		_ = videoRepo.Clean()
	}()

	srv := New(videoRepo)

	// Save a video first
	video := getVideo()
	savedVideo, err := srv.Save(video)
	assert.NoError(t, err)

	// Update the video
	savedVideo.Title = "Updated Title"
	savedVideo.Description = "Updated Description"
	updatedVideo, err := srv.Update(savedVideo)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Title", updatedVideo.Title)
	assert.Equal(t, "Updated Description", updatedVideo.Description)
}

func TestUpdate_InvalidID(t *testing.T) {
	videoRepo, err := repository.NewVideoRepository(*dbType, *dsn)
	assert.NoError(t, err)
	defer func() {
		_ = videoRepo.Clean()
	}()

	srv := New(videoRepo)

	// Try to update with ID = 0
	video := getVideo()
	video.ID = 0
	_, err = srv.Update(video)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ID is required")
}

func TestDelete(t *testing.T) {
	videoRepo, err := repository.NewVideoRepository(*dbType, *dsn)
	assert.NoError(t, err)
	defer func() {
		_ = videoRepo.Clean()
	}()

	srv := New(videoRepo)

	// Save a video first
	video := getVideo()
	savedVideo, err := srv.Save(video)
	assert.NoError(t, err)

	// Delete the video
	err = srv.Delete(savedVideo.ID)
	assert.NoError(t, err)

	// Verify it's deleted
	videos, err := srv.FindAll()
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
