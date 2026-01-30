// Package postgres is postgres implementation of port.VideoRepository.
// It uses pgx postgres driver.

package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/yinebebt/hexagonal-architecture/internal/core/entity"
	"github.com/yinebebt/hexagonal-architecture/internal/core/port"
)

type database struct {
	conn *pgx.Conn
}

func NewVideoRepository(dsn string) (port.VideoRepository, error) {
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	_, err = conn.Exec(context.Background(), `
       CREATE TABLE IF NOT EXISTS videos (
            id SERIAL PRIMARY KEY,
            title VARCHAR(15),
            description VARCHAR(35),
            url VARCHAR(256) UNIQUE,
            person_id BIGINT,
            created_at TIMESTAMPTZ DEFAULT now(),
            updated_at TIMESTAMPTZ DEFAULT now()
        );
    `)
	if err != nil {
		conn.Close(context.Background())
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}

	return &database{
		conn: conn,
	}, nil
}

func (db *database) Save(video *entity.Video) error {
	err := db.conn.QueryRow(context.Background(),
		"INSERT INTO videos (title, description, url, person_id) VALUES ($1, $2, $3, $4) RETURNING id",
		video.Title, video.Description, video.URL, video.PersonID).Scan(&video.ID)
	if err != nil {
		return fmt.Errorf("failed to save video: %w", err)
	}
	return nil
}

func (db *database) Update(video *entity.Video) error {
	_, err := db.conn.Exec(context.Background(),
		"UPDATE videos SET title = $1, description = $2, url = $3, person_id = $4, updated_at = now() WHERE id = $5",
		video.Title, video.Description, video.URL, video.PersonID, video.ID)
	if err != nil {
		return fmt.Errorf("failed to update video: %w", err)
	}
	return nil
}

func (db *database) FindAll() ([]entity.Video, error) {
	rows, err := db.conn.Query(context.Background(),
		"SELECT id, title, description, url, person_id, created_at, updated_at FROM videos")
	if err != nil {
		return nil, fmt.Errorf("failed to query videos: %w", err)
	}
	defer rows.Close()

	var videos []entity.Video
	for rows.Next() {
		var video entity.Video
		if err := rows.Scan(&video.ID, &video.Title, &video.Description, &video.URL,
			&video.PersonID, &video.CreatedAt, &video.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan video: %w", err)
		}
		videos = append(videos, video)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating videos: %w", err)
	}
	return videos, nil
}

func (db *database) Delete(id int64) error {
	_, err := db.conn.Exec(context.Background(), "DELETE FROM videos WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete video: %w", err)
	}
	return nil
}

func (db *database) Clean() error {
	_, err := db.conn.Exec(context.Background(), "DROP TABLE IF EXISTS videos")
	if err != nil {
		return fmt.Errorf("failed to clean videos table: %w", err)
	}
	return nil
}

func (db *database) Close() error {
	return db.conn.Close(context.Background())
}
