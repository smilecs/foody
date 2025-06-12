-- Create users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    media_id UUID,
    date_of_birth DATE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create media table
CREATE TABLE media (
    id SERIAL PRIMARY KEY,
    media_id UUID NOT NULL UNIQUE,
    url VARCHAR(255) NOT NULL,
    media_type VARCHAR(50) NOT NULL CHECK (media_type IN ('video', 'image')),
    author_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (author_id) REFERENCES users(user_id) ON DELETE CASCADE
);

-- Add foreign key constraint to users table for media_id
ALTER TABLE users
ADD CONSTRAINT fk_users_media
FOREIGN KEY (media_id) REFERENCES media(media_id) ON DELETE SET NULL;

-- Create recipe table
CREATE TABLE recipe (
    id SERIAL PRIMARY KEY,
    recipe_id UUID NOT NULL UNIQUE,
    author_id UUID NOT NULL,
    media_id UUID,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    prep_time INTERVAL,
    cook_time INTERVAL,
    total_time INTERVAL,
    servings INT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (author_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (media_id) REFERENCES media(media_id) ON DELETE SET NULL
);

-- Create post table
CREATE TABLE post (
    id SERIAL PRIMARY KEY,
    post_id UUID NOT NULL UNIQUE,
    author_id UUID NOT NULL,
    media_id UUID,
    media_url VARCHAR(255),
    title VARCHAR(255) NOT NULL,
    body TEXT NOT NULL,
    tags TEXT[] DEFAULT '{}',
    recipe_id UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (author_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (media_id) REFERENCES media(media_id) ON DELETE SET NULL,
    FOREIGN KEY (recipe_id) REFERENCES recipe(recipe_id) ON DELETE SET NULL
);

-- Create recipe_ingredients table
CREATE TABLE recipe_ingredients (
    id SERIAL PRIMARY KEY,
    recipe_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    quantity DECIMAL(10,2) NOT NULL,
    unit VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (recipe_id) REFERENCES recipe(recipe_id) ON DELETE CASCADE
);

-- Create recipe_steps table
CREATE TABLE recipe_steps (
    id SERIAL PRIMARY KEY,
    recipe_id UUID NOT NULL,
    step_order INT NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (recipe_id) REFERENCES recipe(recipe_id) ON DELETE CASCADE
);

-- Create meal_plan table
CREATE TABLE meal_plan (
    id SERIAL PRIMARY KEY,
    meal_plan_id UUID NOT NULL UNIQUE,
    recipe_id UUID NOT NULL,
    author_id UUID NOT NULL,
    meal_type VARCHAR(50) NOT NULL CHECK (meal_type IN ('breakfast', 'lunch', 'dinner', 'snack')),
    date DATE NOT NULL,
    verified BOOLEAN DEFAULT FALSE,
    photo_id UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (recipe_id) REFERENCES recipe(recipe_id) ON DELETE CASCADE,
    FOREIGN KEY (author_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (photo_id) REFERENCES media(media_id) ON DELETE SET NULL
);

-- Create indexes for better query performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_media_author_id ON media(author_id);
CREATE INDEX idx_post_author_id ON post(author_id);
CREATE INDEX idx_post_recipe_id ON post(recipe_id);
CREATE INDEX idx_recipe_author_id ON recipe(author_id);
CREATE INDEX idx_recipe_ingredients_recipe_id ON recipe_ingredients(recipe_id);
CREATE INDEX idx_recipe_steps_recipe_id ON recipe_steps(recipe_id);
CREATE INDEX idx_meal_plan_author_id ON meal_plan(author_id);
CREATE INDEX idx_meal_plan_date ON meal_plan(date);
CREATE INDEX idx_meal_plan_recipe_id ON meal_plan(recipe_id); 