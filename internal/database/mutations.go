package database

import (
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/naoya0117/portfolio-v2025-api/internal/models"
)

// Blog Post mutations
func (db *DB) CreateBlogPost(input models.CreateBlogPostInput) (*models.BlogPost, error) {
	now := time.Now()
	status := models.BlogStatusDraft
	if input.Status != nil {
		status = *input.Status
	}

	var publishedAt *string
	if status == models.BlogStatusPublished {
		publishedAtTime := now.Format(time.RFC3339)
		publishedAt = &publishedAtTime
	}

	query := `
		INSERT INTO blog_posts (title, slug, excerpt, content, cover_image_url, tags,
							   status, seo_title, seo_description, published_at, like_count)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at
	`

	post := &models.BlogPost{
		Title:          input.Title,
		Slug:           input.Slug,
		Excerpt:        input.Excerpt,
		Content:        input.Content,
		CoverImageURL:  input.CoverImageURL,
		Tags:           input.Tags,
		Status:         status,
		SeoTitle:       input.SeoTitle,
		SeoDescription: input.SeoDescription,
		PublishedAt:    publishedAt,
		LikeCount:      intPtr(0),
	}

	err := db.QueryRow(
		query, post.Title, post.Slug, ptrToNullString(post.Excerpt),
		post.Content, ptrToNullString(post.CoverImageURL), pq.Array(post.Tags),
		post.Status, ptrToNullString(post.SeoTitle), ptrToNullString(post.SeoDescription),
		ptrToNullString(post.PublishedAt), post.LikeCount,
	).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create blog post: %w", err)
	}

	return post, nil
}

func (db *DB) UpdateBlogPost(id string, input models.UpdateBlogPostInput) (*models.BlogPost, error) {
	// Build dynamic update query
	setParts := []string{"updated_at = NOW()"}
	args := []interface{}{}
	argIndex := 1

	if input.Title != nil {
		setParts = append(setParts, fmt.Sprintf("title = $%d", argIndex))
		args = append(args, *input.Title)
		argIndex++
	}
	if input.Slug != nil {
		setParts = append(setParts, fmt.Sprintf("slug = $%d", argIndex))
		args = append(args, *input.Slug)
		argIndex++
	}
	if input.Excerpt != nil {
		setParts = append(setParts, fmt.Sprintf("excerpt = $%d", argIndex))
		args = append(args, ptrToNullString(input.Excerpt))
		argIndex++
	}
	if input.Content != nil {
		setParts = append(setParts, fmt.Sprintf("content = $%d", argIndex))
		args = append(args, *input.Content)
		argIndex++
	}
	if input.CoverImageURL != nil {
		setParts = append(setParts, fmt.Sprintf("cover_image_url = $%d", argIndex))
		args = append(args, ptrToNullString(input.CoverImageURL))
		argIndex++
	}
	if input.Tags != nil {
		setParts = append(setParts, fmt.Sprintf("tags = $%d", argIndex))
		args = append(args, pq.Array(input.Tags))
		argIndex++
	}
	if input.Status != nil {
		setParts = append(setParts, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, *input.Status)
		argIndex++
		
		// Set publishedAt if changing to published
		if *input.Status == models.BlogStatusPublished {
			setParts = append(setParts, fmt.Sprintf("published_at = COALESCE(published_at, $%d)", argIndex))
			args = append(args, time.Now().Format(time.RFC3339))
			argIndex++
		}
	}
	if input.SeoTitle != nil {
		setParts = append(setParts, fmt.Sprintf("seo_title = $%d", argIndex))
		args = append(args, ptrToNullString(input.SeoTitle))
		argIndex++
	}
	if input.SeoDescription != nil {
		setParts = append(setParts, fmt.Sprintf("seo_description = $%d", argIndex))
		args = append(args, ptrToNullString(input.SeoDescription))
		argIndex++
	}

	// Add WHERE clause
	query := fmt.Sprintf("UPDATE blog_posts SET %s WHERE id = $%d", 
		joinStrings(setParts, ", "), argIndex)
	args = append(args, id)

	_, err := db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update blog post: %w", err)
	}

	// Return updated post
	posts, err := db.queryBlogPosts("SELECT id, title, slug, excerpt, content, cover_image_url, tags, status, seo_title, seo_description, published_at, like_count, created_at, updated_at FROM blog_posts WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	if len(posts) == 0 {
		return nil, fmt.Errorf("blog post not found after update")
	}

	return posts[0], nil
}

func (db *DB) DeleteBlogPost(id string) (bool, error) {
	query := "DELETE FROM blog_posts WHERE id = $1"
	result, err := db.Exec(query, id)
	if err != nil {
		return false, fmt.Errorf("failed to delete blog post: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}

func (db *DB) PublishBlogPost(id string) (*models.BlogPost, error) {
	query := `
		UPDATE blog_posts 
		SET status = 'PUBLISHED', 
			published_at = COALESCE(published_at, $1),
			updated_at = NOW()
		WHERE id = $2
	`

	_, err := db.Exec(query, time.Now().Format(time.RFC3339), id)
	if err != nil {
		return nil, fmt.Errorf("failed to publish blog post: %w", err)
	}

	posts, err := db.queryBlogPosts("SELECT id, title, slug, excerpt, content, cover_image_url, tags, status, seo_title, seo_description, published_at, like_count, created_at, updated_at FROM blog_posts WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	if len(posts) == 0 {
		return nil, fmt.Errorf("blog post not found after publish")
	}

	return posts[0], nil
}

func (db *DB) UnpublishBlogPost(id string) (*models.BlogPost, error) {
	query := `
		UPDATE blog_posts 
		SET status = 'DRAFT', updated_at = NOW()
		WHERE id = $1
	`

	_, err := db.Exec(query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to unpublish blog post: %w", err)
	}

	posts, err := db.queryBlogPosts("SELECT id, title, slug, excerpt, content, cover_image_url, tags, status, seo_title, seo_description, published_at, like_count, created_at, updated_at FROM blog_posts WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	if len(posts) == 0 {
		return nil, fmt.Errorf("blog post not found after unpublish")
	}

	return posts[0], nil
}

// Monologue mutations
func (db *DB) CreateMonologue(input models.CreateMonologueInput) (*models.Monologue, error) {
	now := time.Now()
	isPublished := false
	if input.IsPublished != nil {
		isPublished = *input.IsPublished
	}

	var publishedAt *string
	if isPublished {
		publishedAtTime := now.Format(time.RFC3339)
		publishedAt = &publishedAtTime
	}

	query := `
		INSERT INTO monologues (content, content_type, code_language, code_snippet, tags,
							   is_published, published_at, url, series, category)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at
	`

	mono := &models.Monologue{
		Content:      input.Content,
		ContentType:  input.ContentType,
		CodeLanguage: input.CodeLanguage,
		CodeSnippet:  input.CodeSnippet,
		Tags:         input.Tags,
		IsPublished:  isPublished,
		PublishedAt:  publishedAt,
		URL:          input.URL,
		Series:       input.Series,
		Category:     input.Category,
		LikeCount:    intPtr(0),
	}


	err := db.QueryRow(
		query, mono.Content, mono.ContentType, ptrToNullString(mono.CodeLanguage),
		ptrToNullString(mono.CodeSnippet), pq.Array(mono.Tags), mono.IsPublished,
		ptrToNullString(mono.PublishedAt), ptrToNullString(mono.URL),
		ptrToNullString(mono.Series), ptrToNullString(mono.Category),
	).Scan(&mono.ID, &mono.CreatedAt, &mono.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create monologue: %w", err)
	}

	// Generate URL preview if URL is provided
	if mono.URL != nil {
		preview, err := db.CreateURLPreview(mono.ID, *mono.URL)
		if err == nil {
			mono.URLPreview = preview
		}
	}

	return mono, nil
}

func (db *DB) UpdateMonologue(id string, input models.UpdateMonologueInput) (*models.Monologue, error) {
	setParts := []string{"updated_at = NOW()"}
	args := []interface{}{}
	argIndex := 1

	urlChanged := false

	if input.Content != nil {
		setParts = append(setParts, fmt.Sprintf("content = $%d", argIndex))
		args = append(args, *input.Content)
		argIndex++
	}
	if input.ContentType != nil {
		setParts = append(setParts, fmt.Sprintf("content_type = $%d", argIndex))
		args = append(args, *input.ContentType)
		argIndex++
	}
	if input.CodeLanguage != nil {
		setParts = append(setParts, fmt.Sprintf("code_language = $%d", argIndex))
		args = append(args, ptrToNullString(input.CodeLanguage))
		argIndex++
	}
	if input.CodeSnippet != nil {
		setParts = append(setParts, fmt.Sprintf("code_snippet = $%d", argIndex))
		args = append(args, ptrToNullString(input.CodeSnippet))
		argIndex++
	}
	if input.Tags != nil {
		setParts = append(setParts, fmt.Sprintf("tags = $%d", argIndex))
		args = append(args, pq.Array(input.Tags))
		argIndex++
	}
	if input.IsPublished != nil {
		setParts = append(setParts, fmt.Sprintf("is_published = $%d", argIndex))
		args = append(args, *input.IsPublished)
		argIndex++
		
		if *input.IsPublished {
			setParts = append(setParts, fmt.Sprintf("published_at = COALESCE(published_at, $%d)", argIndex))
			args = append(args, time.Now().Format(time.RFC3339))
			argIndex++
		}
	}
	if input.URL != nil {
		setParts = append(setParts, fmt.Sprintf("url = $%d", argIndex))
		args = append(args, ptrToNullString(input.URL))
		argIndex++
		urlChanged = true
	}
	if input.Series != nil {
		setParts = append(setParts, fmt.Sprintf("series = $%d", argIndex))
		args = append(args, ptrToNullString(input.Series))
		argIndex++
	}
	if input.Category != nil {
		setParts = append(setParts, fmt.Sprintf("category = $%d", argIndex))
		args = append(args, ptrToNullString(input.Category))
		argIndex++
	}

	query := fmt.Sprintf("UPDATE monologues SET %s WHERE id = $%d", 
		joinStrings(setParts, ", "), argIndex)
	args = append(args, id)

	_, err := db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update monologue: %w", err)
	}

	// Handle URL preview regeneration
	if urlChanged && input.URL != nil {
		// Delete existing preview
		db.DeleteURLPreviewByMonologueID(id)
		
		if *input.URL != "" {
			// Create new preview
			db.CreateURLPreview(id, *input.URL)
		}
	}

	return db.GetMonologueByID(id)
}

func (db *DB) DeleteMonologue(id string) (bool, error) {
	// Delete related URL previews first
	db.DeleteURLPreviewByMonologueID(id)
	
	query := "DELETE FROM monologues WHERE id = $1"
	result, err := db.Exec(query, id)
	if err != nil {
		return false, fmt.Errorf("failed to delete monologue: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}

func (db *DB) PublishMonologue(id string) (*models.Monologue, error) {
	query := `
		UPDATE monologues 
		SET is_published = true, 
			published_at = $1,
			updated_at = NOW()
		WHERE id = $2
	`

	_, err := db.Exec(query, time.Now().Format(time.RFC3339), id)
	if err != nil {
		return nil, fmt.Errorf("failed to publish monologue: %w", err)
	}

	return db.GetMonologueByID(id)
}

func (db *DB) UnpublishMonologue(id string) (*models.Monologue, error) {
	query := `
		UPDATE monologues 
		SET is_published = false, 
			published_at = NULL,
			updated_at = NOW()
		WHERE id = $1
	`

	_, err := db.Exec(query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to unpublish monologue: %w", err)
	}

	return db.GetMonologueByID(id)
}

func (db *DB) LikeMonologue(id string) (*models.LikeResponse, error) {
	// For now, just increment like count
	query := `
		UPDATE monologues 
		SET like_count = COALESCE(like_count, 0) + 1,
			updated_at = NOW()
		WHERE id = $1
		RETURNING like_count
	`

	var likeCount int
	err := db.QueryRow(query, id).Scan(&likeCount)
	if err != nil {
		return nil, fmt.Errorf("failed to like monologue: %w", err)
	}

	return &models.LikeResponse{
		ID:        id,
		LikeCount: likeCount,
		IsLiked:   true,
	}, nil
}

func (db *DB) LikeBlogPost(id string) (*models.LikeResponse, error) {
	// Handle the case where the ID might have "blog-" prefix
	cleanID := id
	if len(id) > 5 && id[:5] == "blog-" {
		cleanID = id[5:]
	}
	
	// Increment like count for blog post
	query := `
		UPDATE blog_posts 
		SET like_count = COALESCE(like_count, 0) + 1,
			updated_at = NOW()
		WHERE id = $1
		RETURNING like_count
	`

	var likeCount int
	err := db.QueryRow(query, cleanID).Scan(&likeCount)
	if err != nil {
		return nil, fmt.Errorf("failed to like blog post: %w", err)
	}

	return &models.LikeResponse{
		ID:        cleanID,
		LikeCount: likeCount,
		IsLiked:   true,
	}, nil
}


// URL Preview methods
func (db *DB) CreateURLPreview(monologueID, url string) (*models.URLPreview, error) {
	// In a real implementation, this would fetch the URL and extract metadata
	// For now, create a mock preview
	preview := &models.URLPreview{
		Title:       "Generated Preview",
		Description: stringPtr("This is a generated preview for " + url),
		ImageURL:    stringPtr("https://via.placeholder.com/400x200"),
		SiteName:    stringPtr("Example Site"),
		URL:         url,
		Favicon:     stringPtr("https://via.placeholder.com/32x32"),
		CreatedAt:   time.Now(),
	}

	query := `
		INSERT INTO url_previews (monologue_id, title, description, image_url, site_name, url, favicon)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := db.Exec(
		query, monologueID, preview.Title, ptrToNullString(preview.Description),
		ptrToNullString(preview.ImageURL), ptrToNullString(preview.SiteName),
		preview.URL, ptrToNullString(preview.Favicon),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create URL preview: %w", err)
	}

	return preview, nil
}

func (db *DB) DeleteURLPreviewByMonologueID(monologueID string) error {
	query := "DELETE FROM url_previews WHERE monologue_id = $1"
	_, err := db.Exec(query, monologueID)
	return err
}

// Helper functions
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}
	
	result := strs[0]
	for _, str := range strs[1:] {
		result += sep + str
	}
	return result
}

func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}