// Package sqlite is sqlite implementation of video repository using database/sql
package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
	"github.com/yinebebt/hexagonal-architecture/internal/core/entity"
	"github.com/yinebebt/hexagonal-architecture/internal/core/port"
)

type database struct {
	db *sql.DB
}

func NewVideoRepository(dsn string) (port.VideoRepository, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	repo := &database{db: db}

	if err := repo.createTables(); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return repo, nil
}

func (db *database) createTables() error {
	// Create people table first (referenced by videos)
	createPeopleTable := `
		CREATE TABLE IF NOT EXISTS people (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			first_name VARCHAR(32) NOT NULL,
			last_name VARCHAR(32) NOT NULL,
			age INTEGER NOT NULL,
			email VARCHAR(256) NOT NULL
		)
	`

	if _, err := db.db.Exec(createPeopleTable); err != nil {
		return fmt.Errorf("failed to create people table: %w", err)
	}

	// Create videos table
	createVideosTable := `
		CREATE TABLE IF NOT EXISTS videos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title VARCHAR(100) NOT NULL,
			description VARCHAR(500),
			url VARCHAR(256) NOT NULL UNIQUE,
			person_id INTEGER,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (person_id) REFERENCES people(id)
		)
	`

	if _, err := db.db.Exec(createVideosTable); err != nil {
		return fmt.Errorf("failed to create videos table: %w", err)
	}

	return nil
}

func (db *database) Save(ctx context.Context, video *entity.Video) error {
	// First, save the Director (Person) if it exists and doesn't have an ID
	if video.Director.ID == 0 && (video.Director.FirstName != "" || video.Director.LastName != "") {
		personQuery := `
			INSERT INTO people (first_name, last_name, age, email)
			VALUES (?, ?, ?, ?)
		`
		result, err := db.db.ExecContext(ctx, personQuery, video.Director.FirstName, video.Director.LastName, video.Director.Age, video.Director.Email)
		if err != nil {
			return fmt.Errorf("failed to save person: %w", err)
		}

		personID, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to get person ID: %w", err)
		}
		video.PersonID = uint64(personID)
		video.Director.ID = uint64(personID)
	}

	// Save the video
	query := `
		INSERT INTO videos (title, description, url, person_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	result, err := db.db.ExecContext(ctx, query, video.Title, video.Description, video.URL, video.PersonID, now, now)
	if err != nil {
		return fmt.Errorf("failed to save video: %w", err)
	}

	videoID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get video ID: %w", err)
	}

	video.ID = videoID
	video.CreatedAt = now
	video.UpdatedAt = now

	return nil
}

func (db *database) Update(ctx context.Context, video *entity.Video) error {
	// Update the Director (Person) if it exists
	if video.PersonID > 0 && (video.Director.FirstName != "" || video.Director.LastName != "") {
		personQuery := `
			UPDATE people
			SET first_name = ?, last_name = ?, age = ?, email = ?
			WHERE id = ?
		`
		_, err := db.db.ExecContext(ctx, personQuery, video.Director.FirstName, video.Director.LastName, video.Director.Age, video.Director.Email, video.PersonID)
		if err != nil {
			return fmt.Errorf("failed to update person: %w", err)
		}
	}

	// Update the video
	query := `
		UPDATE videos
		SET title = ?, description = ?, url = ?, person_id = ?, updated_at = ?
		WHERE id = ?
	`

	now := time.Now()
	result, err := db.db.ExecContext(ctx, query, video.Title, video.Description, video.URL, video.PersonID, now, video.ID)
	if err != nil {
		return fmt.Errorf("failed to update video: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("video with id %d not found", video.ID)
	}

	video.UpdatedAt = now

	return nil
}

func (db *database) FindByID(ctx context.Context, id int64) (entity.Video, error) {
	query := `
		SELECT v.id, v.title, v.description, v.url, COALESCE(v.person_id, 0), v.created_at, v.updated_at,
		       p.id, p.first_name, p.last_name, p.age, p.email
		FROM videos v
		LEFT JOIN people p ON v.person_id = p.id
		WHERE v.id = ?
	`

	var video entity.Video
	var personID sql.NullInt64
	var firstName, lastName, email sql.NullString
	var age sql.NullInt64

	err := db.db.QueryRowContext(ctx, query, id).Scan(
		&video.ID, &video.Title, &video.Description, &video.URL,
		&video.PersonID, &video.CreatedAt, &video.UpdatedAt,
		&personID, &firstName, &lastName, &age, &email,
	)
	if err == sql.ErrNoRows {
		return entity.Video{}, fmt.Errorf("video with id %d not found", id)
	}
	if err != nil {
		return entity.Video{}, fmt.Errorf("failed to query video: %w", err)
	}

	if personID.Valid {
		video.Director.ID = uint64(personID.Int64)
		if firstName.Valid {
			video.Director.FirstName = firstName.String
		}
		if lastName.Valid {
			video.Director.LastName = lastName.String
		}
		if age.Valid {
			video.Director.Age = int8(age.Int64)
		}
		if email.Valid {
			video.Director.Email = email.String
		}
	}

	return video, nil
}

func (db *database) FindAll(ctx context.Context) ([]entity.Video, error) {
	query := `
		SELECT v.id, v.title, v.description, v.url, COALESCE(v.person_id, 0), v.created_at, v.updated_at,
		       p.id, p.first_name, p.last_name, p.age, p.email
		FROM videos v
		LEFT JOIN people p ON v.person_id = p.id
		ORDER BY v.id
	`

	rows, err := db.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query videos: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var videos []entity.Video
	for rows.Next() {
		var video entity.Video
		var personID sql.NullInt64
		var firstName, lastName, email sql.NullString
		var age sql.NullInt64

		err := rows.Scan(
			&video.ID, &video.Title, &video.Description, &video.URL,
			&video.PersonID, &video.CreatedAt, &video.UpdatedAt,
			&personID, &firstName, &lastName, &age, &email,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan video: %w", err)
		}

		// Populate Director if person exists
		if personID.Valid {
			video.Director.ID = uint64(personID.Int64)
			if firstName.Valid {
				video.Director.FirstName = firstName.String
			}
			if lastName.Valid {
				video.Director.LastName = lastName.String
			}
			if age.Valid {
				video.Director.Age = int8(age.Int64)
			}
			if email.Valid {
				video.Director.Email = email.String
			}
		}

		videos = append(videos, video)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating videos: %w", err)
	}

	return videos, nil
}

func (db *database) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM videos WHERE id = ?`

	result, err := db.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete video: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("video with id %d not found", id)
	}

	return nil
}

func (db *database) Clean(ctx context.Context) error {
	if _, err := db.db.ExecContext(ctx, "DROP TABLE IF EXISTS videos"); err != nil {
		return fmt.Errorf("failed to drop videos table: %w", err)
	}

	if _, err := db.db.ExecContext(ctx, "DROP TABLE IF EXISTS people"); err != nil {
		return fmt.Errorf("failed to drop people table: %w", err)
	}

	return nil
}

func (db *database) Close() error {
	return db.db.Close()
}
