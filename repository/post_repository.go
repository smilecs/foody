package repository

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
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

func (r *PostRepository) GetPosts(limit, offset int) ([]schema.Post, error) {
	var posts []schema.Post

	query := `SELECT * FROM post ORDER BY created_at DESC LIMIT $1 OFFSET $2`

	rows, err := r.Database.Queryx(query, limit, offset)
	if err != nil {
		log.Printf("error retrieving posts: %v\n", err)
		return nil, err
	}
	defer func(rows *sqlx.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("error retrieving posts database: %v\n", err)
		}
	}(rows)

	for rows.Next() {
		var post schema.Post
		if err := rows.StructScan(&post); err != nil {
			log.Printf("error scanning posts: %v\n", err)
		}
		posts = append(posts, post)
	}
	return posts, nil
}
