package test

import (
	"context"
	"errors"
	"gitlab.com/Yinebeb-01/hexagonalarch/internal/adapter/repository/gorm"
	"gitlab.com/Yinebeb-01/hexagonalarch/internal/core/entity"
	"gitlab.com/Yinebeb-01/hexagonalarch/internal/core/service"
	"testing"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"
)

const (
	TITLE       = "Video Title"
	DESCRIPTION = "Video Description"
	URL         = "https://youtu.be/JgW-i2QjgHQ"
)

var (
	videoRepository = gorm.NewVideoRepository()
	videoSer        = service.New(videoRepository)
	video           = entity.Video{}
	t               *testing.T
)

func adminPostNoVideo() error {
	video = entity.Video{}
	return nil
}

func adminPostSomeVideo() {
	video = entity.Video{
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
	videoSer.Save(video)
}

func adminRunFindAllMethod() error {
	return nil
}

func videoShouldBeVideo() error {
	videoRes := videoSer.FindAll()[0]
	if videoRes.Title == video.Title && videoRes.Description == video.Description {
		return nil
	} else {
		return errors.New("not video")
	}
}

func videoShouldBeNull() error {
	videos := videoSer.FindAll()
	if assert.Empty(t, videos, "should be empty") {
		return nil
	} else {
		return errors.New("not null")
	}
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		videoRepository = gorm.NewVideoRepository()
		videoSer = service.New(videoRepository)
		video = entity.Video{}

		return ctx, nil
	})

	ctx.Step(`^Admin post some video$`, adminPostSomeVideo)
	ctx.Step(`^Admin post no video$`, adminPostNoVideo)
	ctx.Step(`^admin run FindAll method$`, adminRunFindAllMethod)
	ctx.Step(`^video should be video$`, videoShouldBeVideo)
	ctx.Step(`^video should be null$`, videoShouldBeNull)

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		videoRepository.Delete(video)
		return ctx, nil
	})
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t, // Testing instance that will run subtests.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func TestFindAll(t *testing.T) {
	video := gorm.NewVideoRepository()
	service := service.New(video)

	service.Save(getVideo())

	videos := service.FindAll()

	firstVideo := videos[0]
	assert.NotNil(t, videos)
	assert.Equal(t, TITLE, firstVideo.Title)
	assert.Equal(t, DESCRIPTION, firstVideo.Description)
	assert.Equal(t, URL, firstVideo.URL)

	video.Delete(firstVideo)
}

func getVideo() entity.Video {
	return entity.Video{
		Title:       TITLE,
		Description: DESCRIPTION,
		URL:         URL,
	}
}
