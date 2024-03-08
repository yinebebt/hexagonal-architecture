package port

import "gitlab.com/Yinebeb-01/hexagonalarch/internal/core/entity"

type VideoRepository interface {
	Save(entity.Video)
	Update(entity.Video)
	Delete(entity.Video)
	FindAll() []entity.Video
}
