package repository

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/smilecs/foody/config"
	"github.com/smilecs/foody/schema"
)

type PostRepository struct {
	Database config.Database
}

func NewPostRepository(db config.Database) *PostRepository {
	return &PostRepository{Database: db}
}

type PostWithMedia struct {
	schema.Post
	MediaURL  string     `db:"media_url"`
	RecipeID  *uuid.UUID `db:"recipe_id"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
}

func (r *PostRepository) CreatePost(post schema.Post, mediaID uuid.UUID, mediaURL string) error {
	query := `
		INSERT INTO post (post_id, author_id, media_id, media_url, title, body, tags, recipe_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id;
	`

	var recipeID *uuid.UUID
	if post.Recipe != nil {
		recipeID = &post.Recipe.Id
	}

	var userID int
	err := r.Database.QueryRowx(query, post.Id, post.AuthorId, mediaID, mediaURL, post.Title, post.Body, post.Tags, recipeID).Scan(&userID)

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

	query := `
		SELECT *
		FROM post
		ORDER BY created_at DESC LIMIT $1 OFFSET $2
	`

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
			continue
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *PostRepository) GetPostByID(id uuid.UUID) (*PostWithMedia, error) {
	var post PostWithMedia

	query := `
		SELECT *
		FROM post
		WHERE post_id = $1
	`

	err := r.Database.QueryRowx(query, id).StructScan(&post)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *PostRepository) UpdatePost(post schema.Post) error {
	query := `
		UPDATE post SET title = $1, body = $2, tags = $3 WHERE id = $4
	`
	_, err := r.Database.MustExec(query, post.Title, post.Body, post.Tags, post.Id)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostRepository) DeletePost(id uuid.UUID) error {
	query := `
		DELETE FROM post WHERE id = $1
	`
	_, err := r.Database.MustExec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostRepository) GetTotalPostsCount() (int, error) {
	var count int
	err := r.Database.QueryRowx("SELECT COUNT(*) FROM post").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
