package schema

import "github.com/google/uuid"

type User struct {
	Name string `json:"name"`
	Email string `json:"email"`
	Media Media `json:"media"`
	DOB string `json:"date_of_birth"`
}

type Post struct {
	Id       uuid.UUID `json:"post_id"`
	Title    string    `json:"title"`
	Body     string    `json:"body"`
	MediaId  uuid.UUID `json:"media_id"`
	AuthorId uuid.UUID `json:"author_id"`
}

type Recipe struct {
	Id       uuid.UUID `json:"recipe_id"`
	Title    string    `json:"title"`
	Body     string    `json:"body"`
	Tags     string    `json:"tags"`
	MediaId  uuid.UUID `json:"media_id"`
	AuthorId uuid.UUID `json:"author_id"`
}

type MediaType string

const (
    Video MediaType = "video"
    Image MediaType = "image"
)

type Media struct {
	Id uuid.UUID `json:"media_id"`
	URL string `json:"url"`
	MediaType MediaType `json:"media_type"`
	AuthorId uuid.UUID `json:"author_id"`
}

type Interaction struct {
}
