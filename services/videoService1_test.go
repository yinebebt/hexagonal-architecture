package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/Yinebeb-01/simpleAPI/entity"
	"gitlab.com/Yinebeb-01/simpleAPI/repository"
)

const (
	TITLE       = "Video Title"
	DESCRIPTION = "Video Description"
	URL         = "https://youtu.be/JgW-i2QjgHQ"
)

func getVideo() entity.Video {
	return entity.Video{
		Title:       TITLE,
		Description: DESCRIPTION,
		URL:         URL,
	}
}

func TestFindAll(t *testing.T) {
	vidreo := repository.NewVideoRepository()
	service := New(vidreo)

	service.Save(getVideo())

	videos := service.FindAll()

	firstVideo := videos[0]
	assert.Nil(t, videos)
	assert.Equal(t, TITLE, firstVideo.Title)
	assert.Equal(t, DESCRIPTION, firstVideo.Description)
	assert.Equal(t, URL, firstVideo.URL)
}
