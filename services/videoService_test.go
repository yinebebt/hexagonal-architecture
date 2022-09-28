package services

import (
	"context"
	"errors"
	"testing"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"
	"gitlab.com/Yinebeb-01/simpleAPI/entity"
	"gitlab.com/Yinebeb-01/simpleAPI/repository"
)

var (
	videorepository repository.VideoReposiory = repository.NewVideoRepository()
	videoservice    VideoService              = New(videorepository)
	video           entity.Video              = entity.Video{}
	t               *testing.T
)

func adminPostNoVideo() error {
	video = entity.Video{} // null video
	return nil
}

func adminPostSomeVideo() {
	video = entity.Video{
		Title:       "yy cool",
		Description: "fy oto",
		URL:         "https://www.yoe.com/embed/96np1mk",
		Director: entity.Person{
			FirstName: "yina",
			LastName:  "tarku",
			Age:       45,
			Email:     "yinta5@gmail.co",
		},
	}
	videoservice.Save(video)
}

func adminRunFindAllMethod() error {
	return nil
}

func videoShouldBeVideo() error { //arg1 *godog.Table
	Video := entity.Video{}

	Video = videoservice.FindAll()[0]

	ist := assert.Equal(t, Video.Title, video.Title)
	isd := assert.Equal(t, Video.Description, video.Description)
	if ist && isd {
		return nil
	} else {
		return errors.New("not video")
	}
}

func videoShouldBeNull() error {
	video = videoservice.FindAll()[0]
	is := assert.Empty(t, video, "shoube empty")
	if is {
		return nil
	} else {
		return errors.New("not null")
	}
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		videorepository = repository.NewVideoRepository()
		videoservice = New(videorepository)
		video = entity.Video{}

		return ctx, nil
	})

	ctx.Step(`^Admin post some video$`, adminPostSomeVideo)
	ctx.Step(`^Admin post no video$`, adminPostNoVideo)
	ctx.Step(`^admin run FindAll method$`, adminRunFindAllMethod)
	ctx.Step(`^video should be video$`, videoShouldBeVideo)
	ctx.Step(`^video should be null$`, videoShouldBeNull)

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		videorepository.Close()
		return ctx, nil
	})
}
