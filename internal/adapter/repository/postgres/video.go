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
		CREATE TABLE IF NOT EXISTS people (
			id BIGSERIAL PRIMARY KEY,
			first_name VARCHAR(32) NOT NULL,
			last_name VARCHAR(32) NOT NULL,
			age INTEGER NOT NULL,
			email VARCHAR(256) NOT NULL
		);

		CREATE TABLE IF NOT EXISTS videos (
			id BIGSERIAL PRIMARY KEY,
			title VARCHAR(100),
			description VARCHAR(500),
			url VARCHAR(256) UNIQUE,
			person_id BIGINT REFERENCES people(id),
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

func (db *database) Save(ctx context.Context, video *entity.Video) error {
	// Save the Director (Person) if it exists and doesn't have an ID
	if video.Director.ID == 0 && (video.Director.FirstName != "" || video.Director.LastName != "") {
		err := db.conn.QueryRow(ctx,
			"INSERT INTO people (first_name, last_name, age, email) VALUES ($1, $2, $3, $4) RETURNING id",
			video.Director.FirstName, video.Director.LastName, video.Director.Age, video.Director.Email).Scan(&video.Director.ID)
		if err != nil {
			return fmt.Errorf("failed to save person: %w", err)
		}
		video.PersonID = video.Director.ID
	}

	err := db.conn.QueryRow(ctx,
		"INSERT INTO videos (title, description, url, person_id) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at",
		video.Title, video.Description, video.URL, video.PersonID).Scan(&video.ID, &video.CreatedAt, &video.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to save video: %w", err)
	}
	return nil
}

func (db *database) Update(ctx context.Context, video *entity.Video) error {
	// Update the Director (Person) if it exists
	if video.PersonID > 0 && (video.Director.FirstName != "" || video.Director.LastName != "") {
		ct, err := db.conn.Exec(ctx,
			"UPDATE people SET first_name = $1, last_name = $2, age = $3, email = $4 WHERE id = $5",
			video.Director.FirstName, video.Director.LastName, video.Director.Age, video.Director.Email, video.PersonID)
		if err != nil {
			return fmt.Errorf("failed to update person: %w", err)
		}
		if ct.RowsAffected() == 0 {
			return fmt.Errorf("person with id %d not found", video.PersonID)
		}
	}

	ct, err := db.conn.Exec(ctx,
		"UPDATE videos SET title = $1, description = $2, url = $3, person_id = $4, updated_at = now() WHERE id = $5",
		video.Title, video.Description, video.URL, video.PersonID, video.ID)
	if err != nil {
		return fmt.Errorf("failed to update video: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("video with id %d not found", video.ID)
	}
	return nil
}

func (db *database) FindByID(ctx context.Context, id int64) (entity.Video, error) {
	var video entity.Video
	err := db.conn.QueryRow(ctx, `
		SELECT v.id, v.title, v.description, v.url, COALESCE(v.person_id, 0), v.created_at, v.updated_at,
		       COALESCE(p.id, 0), COALESCE(p.first_name, ''), COALESCE(p.last_name, ''), COALESCE(p.age, 0), COALESCE(p.email, '')
		FROM videos v
		LEFT JOIN people p ON v.person_id = p.id
		WHERE v.id = $1`, id).Scan(
		&video.ID, &video.Title, &video.Description, &video.URL,
		&video.PersonID, &video.CreatedAt, &video.UpdatedAt,
		&video.Director.ID, &video.Director.FirstName, &video.Director.LastName,
		&video.Director.Age, &video.Director.Email)
	if err == pgx.ErrNoRows {
		return entity.Video{}, fmt.Errorf("video with id %d not found", id)
	}
	if err != nil {
		return entity.Video{}, fmt.Errorf("failed to query video: %w", err)
	}
	return video, nil
}

func (db *database) FindAll(ctx context.Context) ([]entity.Video, error) {
	rows, err := db.conn.Query(ctx, `
		SELECT v.id, v.title, v.description, v.url, COALESCE(v.person_id, 0), v.created_at, v.updated_at,
		       COALESCE(p.id, 0), COALESCE(p.first_name, ''), COALESCE(p.last_name, ''), COALESCE(p.age, 0), COALESCE(p.email, '')
		FROM videos v
		LEFT JOIN people p ON v.person_id = p.id
		ORDER BY v.id`)
	if err != nil {
		return nil, fmt.Errorf("failed to query videos: %w", err)
	}
	defer rows.Close()

	var videos []entity.Video
	for rows.Next() {
		var video entity.Video
		if err := rows.Scan(&video.ID, &video.Title, &video.Description, &video.URL,
			&video.PersonID, &video.CreatedAt, &video.UpdatedAt,
			&video.Director.ID, &video.Director.FirstName, &video.Director.LastName,
			&video.Director.Age, &video.Director.Email); err != nil {
			return nil, fmt.Errorf("failed to scan video: %w", err)
		}
		videos = append(videos, video)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating videos: %w", err)
	}
	return videos, nil
}

func (db *database) Delete(ctx context.Context, id int64) error {
	ct, err := db.conn.Exec(ctx, "DELETE FROM videos WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete video: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("video with id %d not found", id)
	}
	return nil
}

func (db *database) Clean(ctx context.Context) error {
	if _, err := db.conn.Exec(ctx, "DROP TABLE IF EXISTS videos"); err != nil {
		return fmt.Errorf("failed to drop videos table: %w", err)
	}
	if _, err := db.conn.Exec(ctx, "DROP TABLE IF EXISTS people"); err != nil {
		return fmt.Errorf("failed to drop people table: %w", err)
	}
	return nil
}

func (db *database) Close() error {
	return db.conn.Close(context.Background())
}
