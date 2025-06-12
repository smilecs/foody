package repository

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/smilecs/foody/config"
	"github.com/smilecs/foody/schema"
)

type MealPlanRepository struct {
	Database config.Database
}

func NewMealPlanRepository(db config.Database) *MealPlanRepository {
	return &MealPlanRepository{Database: db}
}

type MealPlanWithMedia struct {
	schema.MealPlan
	MediaURL string `db:"media_url"`
}

func (r *MealPlanRepository) CreateMealPlan(mealPlan schema.MealPlan) error {
	query := `
		INSERT INTO meal_plan (meal_plan_id, recipe_id, author_id, meal_type, date, verified, photo_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id;
	`

	var mealPlanID int
	err := r.Database.QueryRowx(query,
		mealPlan.Id,
		mealPlan.RecipeId,
		mealPlan.AuthorId,
		mealPlan.MealType,
		mealPlan.Date,
		mealPlan.Verified,
		mealPlan.PhotoId,
		time.Now(),
		time.Now(),
	).Scan(&mealPlanID)

	if err != nil {
		log.Printf("error creating meal plan: %v\n", err)
		return err
	}

	return nil
}

func (r *MealPlanRepository) GetMealPlansByAuthorID(authorID uuid.UUID) ([]MealPlanWithMedia, error) {
	var mealPlans []MealPlanWithMedia

	query := `
		SELECT mp.*, m.url as media_url
		FROM meal_plan mp
		LEFT JOIN media m ON mp.photo_id = m.media_id
		WHERE mp.author_id = $1
		ORDER BY mp.date DESC
	`

	rows, err := r.Database.Queryx(query, authorID)
	if err != nil {
		log.Printf("error retrieving meal plans: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var mealPlan MealPlanWithMedia
		if err := rows.StructScan(&mealPlan); err != nil {
			log.Printf("error scanning meal plan: %v\n", err)
			continue
		}
		mealPlans = append(mealPlans, mealPlan)
	}

	return mealPlans, nil
}

func (r *MealPlanRepository) GetMealPlanByID(id uuid.UUID) (*MealPlanWithMedia, error) {
	var mealPlan MealPlanWithMedia

	query := `
		SELECT mp.*, m.url as media_url
		FROM meal_plan mp
		LEFT JOIN media m ON mp.photo_id = m.media_id
		WHERE mp.meal_plan_id = $1
	`

	err := r.Database.QueryRowx(query, id).StructScan(&mealPlan)
	if err != nil {
		return nil, err
	}

	return &mealPlan, nil
}

func (r *MealPlanRepository) UpdateMealPlan(mealPlan schema.MealPlan) error {
	query := `
		UPDATE meal_plan 
		SET recipe_id = $1, meal_type = $2, date = $3, verified = $4, photo_id = $5, updated_at = $6
		WHERE meal_plan_id = $7
	`

	_, err := r.Database.Exec(query,
		mealPlan.RecipeId,
		mealPlan.MealType,
		mealPlan.Date,
		mealPlan.Verified,
		mealPlan.PhotoId,
		time.Now(),
		mealPlan.Id,
	)

	if err != nil {
		log.Printf("error updating meal plan: %v\n", err)
		return err
	}

	return nil
}

func (r *MealPlanRepository) DeleteMealPlan(id uuid.UUID) error {
	query := `DELETE FROM meal_plan WHERE meal_plan_id = $1`
	_, err := r.Database.Exec(query, id)
	if err != nil {
		log.Printf("error deleting meal plan: %v\n", err)
		return err
	}
	return nil
}
