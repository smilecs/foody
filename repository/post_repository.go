package repository

import (
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/smilecs/foody/config"
	"github.com/smilecs/foody/schema"
)

type PostRepository struct {
	Database config.Database
}

type PostWithMedia struct {
	schema.Post
	MediaURL string `db:"media_url"`
}

func (r *PostRepository) CreatePost(post schema.Post, mediaID uuid.UUID, mediaURL string) error {
	query := `
		INSERT INTO post (post_id, author_id, media_id, media_url, title, body, tags)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id;
	`

	var userID int
	err := r.Database.QueryRowx(query, post.Id, post.AuthorId, mediaID, mediaURL, post.Title, post.Body, post.Tags).Scan(&userID)

	if err != nil {
		log.Println("error creating post: ", err)
		return err
	}

	return nil
}

func (r *PostRepository) GetPostByUserID(id uuid.UUID) (*PostWithMedia, error) {
	var post PostWithMedia
	err := r.Database.QueryRowx("SELECT * FROM post WHERE author_id = $1", id).StructScan(&post)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *PostRepository) GetPosts(limit, offset int) ([]PostWithMedia, error) {
	var posts []PostWithMedia

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
		var post PostWithMedia
		if err := rows.StructScan(&post); err != nil {
			log.Printf("error scanning posts: %v\n", err)
		}
		posts = append(posts, post)
	}
	return posts, nil
}
