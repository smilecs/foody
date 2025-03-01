package repository

import (
	"github.com/google/uuid"
	"github.com/smilecs/foody/db"
	"github.com/smilecs/foody/schema"
	"log"
)

type PostRepository struct {
	Database db.Database
}

func (r *PostRepository) CreatePost(post schema.Post) error {
	query := `
		INSERT INTO post (post_id, author_id, media_id, title, body, tags)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id;
	`

	var userID int
	err := r.Database.QueryRowx(query, post.Id, post.AuthorId, post.MediaId, post.Title, post.Body, post.Tags).Scan(&userID)

	if err != nil {
		log.Println("error creating post: ", err)
		return err
	}

	return nil
}

func (r *PostRepository) GetPostByUserID(id uuid.UUID) (*schema.Post, error) {
	var post schema.Post
	err := r.Database.QueryRowx("SELECT * FROM post WHERE author_id = $1", id).StructScan(&post)
	if err != nil {
		return nil, err
	}
	return &post, nil
}
