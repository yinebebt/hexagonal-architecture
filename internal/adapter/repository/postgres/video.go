// Package postgres is postgres implementation of video reository, it uses pgx
package postgres

import (
	"context"
	"log"

	"github.com/Yinebeb-01/hexagonalarch/internal/core/entity"
	"github.com/Yinebeb-01/hexagonalarch/internal/core/port"
	"github.com/jackc/pgx/v4"
)

type Database struct {
	conn *pgx.Conn
}

func NewVideoRepository(dsn string) port.VideoRepository {
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
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
		log.Fatalf("Failed to create schema: %v\n", err)
	}

	return &Database{
		conn: conn,
	}
}

func (db *Database) Save(video entity.Video) error {
	_, err := db.conn.Exec(context.Background(), "INSERT INTO videos (title, description,url) VALUES ($1, $2,$3)",
		video.Title, video.Description, video.URL)
	return err
}

func (db *Database) Update(video entity.Video) error {
	_, err := db.conn.Exec(context.Background(), "UPDATE videos SET title = $1, description = $2 WHERE id = $3",
		video.Title, video.Description, video.ID)
	return err
}

func (db *Database) FindAll() ([]entity.Video, error) {
	rows, err := db.conn.Query(context.Background(), "SELECT id, title, description,url FROM videos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videos []entity.Video
	for rows.Next() {
		var video entity.Video
		if err := rows.Scan(&video.ID, &video.Title, &video.Description, &video.URL); err != nil {
			return nil, err
		}
		videos = append(videos, video)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return videos, nil
}

func (db *Database) Delete(video entity.Video) error {
	_, err := db.conn.Exec(context.Background(), "DELETE FROM videos WHERE id = $1", video.ID)
	return err
}

func (db *Database) Clean() error {
	_, err := db.conn.Exec(context.Background(), "drop table if exists videos")
	return err
}
