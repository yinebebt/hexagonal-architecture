package port

import "github.com/yinebebt/hexagonal-architecture/internal/core/entity"

type VideoRepository interface {
	Save(entity.Video) error
	Update(entity.Video) error
	Delete(id int64) error
	FindAll() ([]entity.Video, error)
	Clean() error
}
