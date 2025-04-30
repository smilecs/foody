package schema

import (
	"time"

	"github.com/google/uuid"
)

type PostWithMedia struct {
	Id        uuid.UUID  `db:"post_id" json:"id"`
	AuthorId  uuid.UUID  `db:"author_id" json:"author_id"`
	MediaId   *uuid.UUID `db:"media_id" json:"media_id"`
	MediaURL  *string    `db:"media_url" json:"media_url"`
	Title     string     `db:"title" json:"title"`
	Body      string     `db:"body" json:"body"`
	Tags      []string   `db:"tags" json:"tags"`
	Recipe    *Recipe    `json:"recipe"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
}
