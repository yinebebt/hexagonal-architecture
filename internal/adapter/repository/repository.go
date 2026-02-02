package repository

import (
	"fmt"

	"github.com/yinebebt/hexagonal-architecture/internal/adapter/repository/postgres"
	"github.com/yinebebt/hexagonal-architecture/internal/adapter/repository/sqlite"
	"github.com/yinebebt/hexagonal-architecture/internal/core/port"
)

func NewVideoRepository(dbType, dsn string) (port.VideoRepository, error) {
	switch dbType {
	case "sqlite":
		return sqlite.NewVideoRepository(dsn)
	case "postgres":
		return postgres.NewVideoRepository(dsn)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}
}
