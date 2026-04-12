package port

import (
	"context"

	"github.com/yinebebt/hexagonal-architecture/internal/core/entity"
)

// VideoService defines the business logic contract for video operations.
// This is the primary port that driving adapters (REST, GraphQL, gRPC) depend on.
type VideoService interface {
	Save(ctx context.Context, video entity.Video) (entity.Video, error)
	Update(ctx context.Context, video entity.Video) (entity.Video, error)
	Delete(ctx context.Context, id int64) error
	FindByID(ctx context.Context, id int64) (entity.Video, error)
	FindAll(ctx context.Context) ([]entity.Video, error)
}
