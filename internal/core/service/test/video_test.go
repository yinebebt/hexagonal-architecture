package test

import (
	"context"
	"errors"
	"flag"
	"testing"

	"github.com/Yinebeb-01/hexagonalarch/internal/adapter/repository"
	"github.com/Yinebeb-01/hexagonalarch/internal/core/entity"
	"github.com/Yinebeb-01/hexagonalarch/internal/core/port"
	"github.com/Yinebeb-01/hexagonalarch/internal/core/service"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"
)

const (
	TITLE       = "Video Title"
	DESCRIPTION = "Video Description"
	URL         = "https://youtu.be/JgW-i2QjgHQ"
)

var (
	dbType = flag.String("dbtype", "sqlite", "Type of the database (sqliteor postgres)")
	dsn    = flag.String("dsn", "test.db", "Data source name for the database")
)

var (
	videoRepository port.VideoRepository
	videoSer        service.VideoService
	video           entity.Video
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
	videos, err := videoSer.FindAll()
	if err != nil {
		return err
	}
	videoRes := videos[0]
	if videoRes.Title == video.Title && videoRes.Description == video.Description {
		return nil
	} else {
		return errors.New("not video")
	}
}

func videoShouldBeNull() error {
	videos, err := videoSer.FindAll()
	if err != nil {
		return nil
	}
	if assert.Empty(t, videos, "should be empty") {
		return nil
	} else {
		return errors.New("not null")
	}
}

func InitializeScenario(ctx *godog.ScenarioContext) {

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		videoRepository = repository.NewVideoRepository(*dbType, *dsn)
		videoSer = service.New(videoRepository)
		video = entity.Video{}
		return ctx, nil
	})

	ctx.Step(`^Admin post some video$`, adminPostSomeVideo)
	ctx.Step(`^Admin post no video$`, adminPostNoVideo)
	ctx.Step(`^admin run FindAll method$`, adminRunFindAllMethod)
	ctx.Step(`^video should be video$`, videoShouldBeVideo)
	ctx.Step(`^video should be null$`, videoShouldBeNull)

	ctx.After(func(ctx context.Context, sc *godog.Scenario, _ error) (context.Context, error) {
		err := videoRepository.Clean()
		return ctx, err
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
	video := repository.NewVideoRepository(*dbType, *dsn)
	service := service.New(video)

	_, err := service.Save(getVideo())
	assert.Nil(t, err)

	videos, err := service.FindAll()

	firstVideo := videos[0]
	assert.NotNil(t, videos)
	assert.Equal(t, TITLE, firstVideo.Title)
	assert.Equal(t, DESCRIPTION, firstVideo.Description)
	assert.Equal(t, URL, firstVideo.URL)
	assert.Nil(t, err)

	video.Delete(firstVideo)
}

func getVideo() entity.Video {
	return entity.Video{
		Title:       TITLE,
		Description: DESCRIPTION,
		URL:         URL,
	}
}
