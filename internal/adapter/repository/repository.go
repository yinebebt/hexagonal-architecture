package repository

import (
	"fmt"
	"log"

	"github.com/yinebebt/hexagonal-architecture/internal/adapter/repository/postgres"
	"github.com/yinebebt/hexagonal-architecture/internal/adapter/repository/sqlite"
	"github.com/yinebebt/hexagonal-architecture/internal/core/port"
)

func NewVideoRepository(dbType, dsn string) port.VideoRepository {
	switch dbType {
	case "sqlite":
		return sqlite.NewVideoRepository(dsn)
	case "postgres":
		return postgres.NewVideoRepository(dsn)
	default:
		log.Fatalf(fmt.Sprintf("Unsupported database type: %s", dbType))
		return nil
	}
}
