package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

func NewConnection() (*DB, error) {
	// Get database configuration from environment variables
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// If database environment variables are not set, return nil (use mock data)
	if host == "" || user == "" || password == "" || dbname == "" {
		return nil, nil
	}

	// Create connection string
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Open database connection
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{db}, nil
}

func (db *DB) Close() error {
	if db.DB != nil {
		return db.DB.Close()
	}
	return nil
}

// CreateTables creates the necessary tables if they don't exist
func (db *DB) MigrateTables() error {
	// Add like_count column to existing blog_posts table if it doesn't exist
	_, err := db.Exec(`
		ALTER TABLE blog_posts 
		ADD COLUMN IF NOT EXISTS like_count INTEGER DEFAULT 0
	`)
	if err != nil {
		return fmt.Errorf("failed to add like_count column to blog_posts: %w", err)
	}

	// Add like_count column to existing monologues table if it doesn't exist
	_, err = db.Exec(`
		ALTER TABLE monologues 
		ADD COLUMN IF NOT EXISTS like_count INTEGER DEFAULT 0
	`)
	if err != nil {
		return fmt.Errorf("failed to add like_count column to monologues: %w", err)
	}

	return nil
}

func (db *DB) CreateTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS profiles (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(255) NOT NULL,
			title VARCHAR(255),
			bio TEXT,
			avatar_url VARCHAR(500),
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		)`,
		
		`CREATE TABLE IF NOT EXISTS social_links (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			profile_id UUID REFERENCES profiles(id),
			platform VARCHAR(100) NOT NULL,
			url VARCHAR(500) NOT NULL,
			icon VARCHAR(100),
			created_at TIMESTAMP DEFAULT NOW()
		)`,
		
		`CREATE TABLE IF NOT EXISTS skills (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(255) NOT NULL,
			category VARCHAR(255) NOT NULL,
			level INTEGER NOT NULL CHECK (level >= 1 AND level <= 10),
			icon_url VARCHAR(500),
			display_order INTEGER DEFAULT 0,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		)`,
		
		`CREATE TABLE IF NOT EXISTS experiences (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			company VARCHAR(255) NOT NULL,
			position VARCHAR(255) NOT NULL,
			description TEXT,
			start_date VARCHAR(20) NOT NULL,
			end_date VARCHAR(20),
			is_current BOOLEAN DEFAULT FALSE,
			technologies TEXT[],
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		)`,
		
		`CREATE TABLE IF NOT EXISTS code_categories (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(100) NOT NULL,
			slug VARCHAR(100) NOT NULL UNIQUE,
			description TEXT,
			parent_id UUID REFERENCES code_categories(id),
			color VARCHAR(7),
			icon VARCHAR(10),
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		)`,
		
		`CREATE TABLE IF NOT EXISTS monologues (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			content TEXT NOT NULL,
			content_type VARCHAR(20) NOT NULL CHECK (content_type IN ('POST', 'CODE', 'IMAGE', 'URL_PREVIEW')),
			code_language VARCHAR(50),
			code_snippet TEXT,
			tags TEXT[],
			is_published BOOLEAN DEFAULT FALSE,
			published_at TIMESTAMP,
			url VARCHAR(2048),
			series VARCHAR(255),
			category VARCHAR(255),
			code_category_id UUID REFERENCES code_categories(id),
			difficulty VARCHAR(20) CHECK (difficulty IN ('BEGINNER', 'INTERMEDIATE', 'ADVANCED')),
			like_count INTEGER DEFAULT 0,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		)`,
		
		`CREATE TABLE IF NOT EXISTS url_previews (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			monologue_id UUID REFERENCES monologues(id),
			title VARCHAR(500) NOT NULL,
			description TEXT,
			image_url VARCHAR(2048),
			site_name VARCHAR(200),
			url VARCHAR(2048) NOT NULL,
			favicon VARCHAR(2048),
			created_at TIMESTAMP DEFAULT NOW()
		)`,
		
		`CREATE TABLE IF NOT EXISTS blog_posts (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			title VARCHAR(500) NOT NULL,
			slug VARCHAR(500) NOT NULL UNIQUE,
			excerpt TEXT,
			content TEXT NOT NULL,
			cover_image_url VARCHAR(2048),
			tags TEXT[],
			status VARCHAR(20) NOT NULL CHECK (status IN ('DRAFT', 'PUBLISHED', 'ARCHIVED')),
			seo_title VARCHAR(500),
			seo_description TEXT,
			published_at TIMESTAMP,
			like_count INTEGER DEFAULT 0,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		)`,
		
		`CREATE TABLE IF NOT EXISTS monologue_likes (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			monologue_id UUID NOT NULL REFERENCES monologues(id),
			user_ip VARCHAR(45),
			user_id UUID,
			created_at TIMESTAMP DEFAULT NOW(),
			UNIQUE(monologue_id, user_ip)
		)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
	}

	return nil
}