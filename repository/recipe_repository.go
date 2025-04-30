package repository

import (
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/smilecs/foody/config"
	"github.com/smilecs/foody/schema"
)

type RecipeRepository struct {
	Database config.Database
}

func NewRecipeRepository(db config.Database) *RecipeRepository {
	return &RecipeRepository{Database: db}
}

type RecipeWithMedia struct {
	schema.Recipe
	MediaURL string `db:"media_url"`
}

func (r *RecipeRepository) CreateRecipe(recipe schema.Recipe, mediaID uuid.UUID, mediaURL string) error {
	tx, err := r.Database.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Insert recipe
	query := `
		INSERT INTO recipe (recipe_id, author_id, media_id, title, description, prep_time, cook_time, total_time, servings)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id;
	`

	var recipeID int
	err = tx.QueryRowx(query,
		recipe.Id,
		recipe.AuthorId,
		mediaID,
		recipe.Title,
		recipe.Description,
		recipe.PrepTime,
		recipe.CookTime,
		recipe.TotalTime,
		recipe.Servings,
	).Scan(&recipeID)

	if err != nil {
		log.Println("error creating recipe: ", err)
		return err
	}

	// Insert ingredients
	if len(recipe.Ingredients) > 0 {
		ingredientsQuery := `
			INSERT INTO recipe_ingredients (recipe_id, name, quantity, unit)
			VALUES ($1, $2, $3, $4)
		`
		for _, ingredient := range recipe.Ingredients {
			_, err = tx.Exec(ingredientsQuery,
				recipe.Id,
				ingredient.Name,
				ingredient.Quantity,
				ingredient.Unit,
			)
			if err != nil {
				log.Println("error creating ingredient: ", err)
				return err
			}
		}
	}

	// Insert steps
	if len(recipe.Steps) > 0 {
		stepsQuery := `
			INSERT INTO recipe_steps (recipe_id, step_order, description)
			VALUES ($1, $2, $3)
		`
		for _, step := range recipe.Steps {
			_, err = tx.Exec(stepsQuery,
				recipe.Id,
				step.Order,
				step.Description,
			)
			if err != nil {
				log.Println("error creating step: ", err)
				return err
			}
		}
	}

	return tx.Commit()
}

func (r *RecipeRepository) GetRecipes(limit, offset int) ([]RecipeWithMedia, error) {
	var recipes []RecipeWithMedia

	query := `
		SELECT r.*, m.url as media_url
		FROM recipe r
		LEFT JOIN media m ON r.media_id = m.media_id
		ORDER BY r.created_at DESC LIMIT $1 OFFSET $2
	`

	rows, err := r.Database.Queryx(query, limit, offset)
	if err != nil {
		log.Printf("error retrieving recipes: %v\n", err)
		return nil, err
	}
	defer func(rows *sqlx.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("error closing rows: %v\n", err)
		}
	}(rows)

	for rows.Next() {
		var recipe RecipeWithMedia
		if err := rows.StructScan(&recipe); err != nil {
			log.Printf("error scanning recipe: %v\n", err)
			continue
		}

		// Get ingredients
		ingredientsQuery := `SELECT name, quantity, unit FROM recipe_ingredients WHERE recipe_id = $1`
		ingredientRows, err := r.Database.Queryx(ingredientsQuery, recipe.Id)
		if err != nil {
			log.Printf("error retrieving ingredients: %v\n", err)
			continue
		}
		defer ingredientRows.Close()

		for ingredientRows.Next() {
			var ingredient schema.Ingredient
			if err := ingredientRows.StructScan(&ingredient); err != nil {
				log.Printf("error scanning ingredient: %v\n", err)
				continue
			}
			recipe.Ingredients = append(recipe.Ingredients, ingredient)
		}

		// Get steps
		stepsQuery := `SELECT step_order, description FROM recipe_steps WHERE recipe_id = $1 ORDER BY step_order`
		stepRows, err := r.Database.Queryx(stepsQuery, recipe.Id)
		if err != nil {
			log.Printf("error retrieving steps: %v\n", err)
			continue
		}
		defer stepRows.Close()

		for stepRows.Next() {
			var step schema.Step
			if err := stepRows.StructScan(&step); err != nil {
				log.Printf("error scanning step: %v\n", err)
				continue
			}
			recipe.Steps = append(recipe.Steps, step)
		}

		recipes = append(recipes, recipe)
	}
	return recipes, nil
}

func (r *RecipeRepository) GetRecipeByID(id uuid.UUID) (*RecipeWithMedia, error) {
	var recipe RecipeWithMedia
	query := `
		SELECT r.*, m.url as media_url
		FROM recipe r
		LEFT JOIN media m ON r.media_id = m.media_id
		WHERE r.recipe_id = $1
	`
	err := r.Database.QueryRowx(query, id).StructScan(&recipe)
	if err != nil {
		return nil, err
	}

	// Get ingredients
	ingredientsQuery := `SELECT name, quantity, unit FROM recipe_ingredients WHERE recipe_id = $1`
	ingredientRows, err := r.Database.Queryx(ingredientsQuery, id)
	if err != nil {
		return nil, err
	}
	defer ingredientRows.Close()

	for ingredientRows.Next() {
		var ingredient schema.Ingredient
		if err := ingredientRows.StructScan(&ingredient); err != nil {
			log.Printf("error scanning ingredient: %v\n", err)
			continue
		}
		recipe.Ingredients = append(recipe.Ingredients, ingredient)
	}

	// Get steps
	stepsQuery := `SELECT step_order, description FROM recipe_steps WHERE recipe_id = $1 ORDER BY step_order`
	stepRows, err := r.Database.Queryx(stepsQuery, id)
	if err != nil {
		return nil, err
	}
	defer stepRows.Close()

	for stepRows.Next() {
		var step schema.Step
		if err := stepRows.StructScan(&step); err != nil {
			log.Printf("error scanning step: %v\n", err)
			continue
		}
		recipe.Steps = append(recipe.Steps, step)
	}

	return &recipe, nil
}

func (r *RecipeRepository) GetRecipesByAuthorID(authorID uuid.UUID) ([]RecipeWithMedia, error) {
	var recipes []RecipeWithMedia

	query := `
		SELECT r.*, m.url as media_url
		FROM recipe r
		LEFT JOIN media m ON r.media_id = m.media_id
		WHERE r.author_id = $1
		ORDER BY r.created_at DESC
	`

	rows, err := r.Database.Queryx(query, authorID)
	if err != nil {
		log.Printf("error retrieving recipes: %v\n", err)
		return nil, err
	}
	defer func(rows *sqlx.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("error closing rows: %v\n", err)
		}
	}(rows)

	for rows.Next() {
		var recipe RecipeWithMedia
		if err := rows.StructScan(&recipe); err != nil {
			log.Printf("error scanning recipe: %v\n", err)
			continue
		}
		recipes = append(recipes, recipe)
	}
	return recipes, nil
}

func (r *RecipeRepository) UpdateRecipe(recipe schema.Recipe) error {
	tx, err := r.Database.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Update recipe
	query := `
		UPDATE recipe 
		SET title = $1, description = $2, prep_time = $3, cook_time = $4, total_time = $5, servings = $6
		WHERE recipe_id = $7
	`
	_, err = tx.Exec(query,
		recipe.Title,
		recipe.Description,
		recipe.PrepTime,
		recipe.CookTime,
		recipe.TotalTime,
		recipe.Servings,
		recipe.Id,
	)
	if err != nil {
		return err
	}

	// Delete existing ingredients and steps
	_, err = tx.Exec("DELETE FROM recipe_ingredients WHERE recipe_id = $1", recipe.Id)
	if err != nil {
		return err
	}
	_, err = tx.Exec("DELETE FROM recipe_steps WHERE recipe_id = $1", recipe.Id)
	if err != nil {
		return err
	}

	// Insert new ingredients
	if len(recipe.Ingredients) > 0 {
		ingredientsQuery := `
			INSERT INTO recipe_ingredients (recipe_id, name, quantity, unit)
			VALUES ($1, $2, $3, $4)
		`
		for _, ingredient := range recipe.Ingredients {
			_, err = tx.Exec(ingredientsQuery,
				recipe.Id,
				ingredient.Name,
				ingredient.Quantity,
				ingredient.Unit,
			)
			if err != nil {
				return err
			}
		}
	}

	// Insert new steps
	if len(recipe.Steps) > 0 {
		stepsQuery := `
			INSERT INTO recipe_steps (recipe_id, step_order, description)
			VALUES ($1, $2, $3)
		`
		for _, step := range recipe.Steps {
			_, err = tx.Exec(stepsQuery,
				recipe.Id,
				step.Order,
				step.Description,
			)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (r *RecipeRepository) DeleteRecipe(id uuid.UUID) error {
	query := `
		DELETE FROM recipe WHERE recipe_id = $1
	`
	_, err := r.Database.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *RecipeRepository) GetTotalRecipesCount() (int, error) {
	var count int
	err := r.Database.QueryRowx("SELECT COUNT(*) FROM recipe").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
