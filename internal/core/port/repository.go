package port

import "github.com/yinebebt/hexagonal-architecture/internal/core/entity"

type VideoRepository interface {
	Save(*entity.Video) error
	Update(*entity.Video) error
	Delete(id int64) error
	FindAll() ([]entity.Video, error)
	// Clean removes all data from the database. This method is intended for testing only.
	Clean() error
	// Close closes the database connection and releases any resources.
	Close() error
}
