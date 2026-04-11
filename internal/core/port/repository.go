package port

import (
	"context"

	"github.com/yinebebt/hexagonal-architecture/internal/core/entity"
)

type VideoRepository interface {
	Save(ctx context.Context, video *entity.Video) error
	Update(ctx context.Context, video *entity.Video) error
	Delete(ctx context.Context, id int64) error
	FindByID(ctx context.Context, id int64) (entity.Video, error)
	FindAll(ctx context.Context) ([]entity.Video, error)
	// Clean removes all data from the database. This method is intended for testing only.
	Clean(ctx context.Context) error
	// Close closes the database connection and releases any resources.
	Close() error
}
