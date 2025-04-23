package repository

import (
	"github.com/google/uuid"
	"github.com/smilecs/foody/config"
	"github.com/smilecs/foody/schema"
)

type MediaRepository struct {
	Database config.Database
}

func (r *MediaRepository) CreateMedia(media schema.Media) (uuid.UUID, error) {
	query := `
		INSERT INTO media (media_id, url, media_type, author_id)
		VALUES ($1, $2, $3, $4)
		RETURNING media_id;
	`

	var mediaID uuid.UUID
	err := r.Database.QueryRowx(query, media.Id, media.URL, media.MediaType, media.AuthorId).Scan(&mediaID)
	if err != nil {
		return uuid.Nil, err
	}

	return mediaID, nil
}

func (r *MediaRepository) GetMediaByID(id uuid.UUID) (*schema.Media, error) {
	var media schema.Media
	err := r.Database.QueryRowx("SELECT * FROM media WHERE media_id = $1", id).StructScan(&media)
	if err != nil {
		return nil, err
	}
	return &media, nil
}

func (r *MediaRepository) GetMediaByAuthorID(authorID uuid.UUID) ([]schema.Media, error) {
	var mediaList []schema.Media
	rows, err := r.Database.Queryx("SELECT * FROM media WHERE author_id = $1", authorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var media schema.Media
		if err := rows.StructScan(&media); err != nil {
			return nil, err
		}
		mediaList = append(mediaList, media)
	}

	return mediaList, nil
}
