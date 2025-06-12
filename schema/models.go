package schema

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id       uuid.UUID `json:"user_id"`
	Name     string    `json:"name"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Media    Media     `json:"media"`
	DOB      string    `json:"date_of_birth"`
	Password string    `json:"-"`
}

type Post struct {
	Id       uuid.UUID `json:"post_id"`
	Title    string    `json:"title"`
	Body     string    `json:"body"`
	MediaId  uuid.UUID `json:"media_id"`
	AuthorId uuid.UUID `json:"author_id"`
	Tags     string    `json:"tags"`
	Recipe   *Recipe   `json:"recipe,omitempty"`
}

type Ingredient struct {
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
	Unit     string  `json:"unit"`
}

type Step struct {
	Order       int    `json:"order"`
	Description string `json:"description"`
}

type Recipe struct {
	Id          uuid.UUID      `json:"recipe_id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Ingredients []Ingredient   `json:"ingredients"`
	Steps       []Step         `json:"steps,omitempty"`
	PrepTime    *time.Duration `json:"prep_time,omitempty"`
	CookTime    *time.Duration `json:"cook_time,omitempty"`
	TotalTime   *time.Duration `json:"total_time,omitempty"`
	Servings    *int           `json:"servings,omitempty"`
	AuthorId    uuid.UUID      `json:"author_id"`
	MediaId     uuid.UUID      `json:"media_id,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type MealType string

const (
	Breakfast MealType = "breakfast"
	Lunch     MealType = "lunch"
	Dinner    MealType = "dinner"
	Snack     MealType = "snack"
)

type MealPlan struct {
	Id        uuid.UUID  `json:"meal_plan_id"`
	RecipeId  uuid.UUID  `json:"recipe_id"`
	AuthorId  uuid.UUID  `json:"author_id"`
	MealType  MealType   `json:"meal_type"`
	Date      time.Time  `json:"date"`
	Verified  bool       `json:"verified"`
	PhotoId   *uuid.UUID `json:"photo_id,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type MediaType string

const (
	Video MediaType = "video"
	Image MediaType = "image"
)

type Media struct {
	Id        uuid.UUID `json:"media_id"`
	URL       string    `json:"url"`
	MediaType MediaType `json:"media_type"`
	AuthorId  uuid.UUID `json:"author_id"`
}

type Interaction struct {
}
