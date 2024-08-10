package port

import "github.com/Yinebeb-01/hexagonalarch/internal/core/entity"

type VideoRepository interface {
	Save(entity.Video) error
	Update(entity.Video) error
	Delete(entity.Video) error
	FindAll() ([]entity.Video, error)
	Clean() error
}
