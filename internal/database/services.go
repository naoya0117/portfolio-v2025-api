package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/naoya0117/portfolio-v2025-api/internal/models"
)

// Profile methods
func (db *DB) GetProfile(id string) (*models.Profile, error) {
	query := `
		SELECT id, name, title, bio, avatar_url, created_at, updated_at
		FROM profiles WHERE id = $1
	`
	
	profile := &models.Profile{}
	var title, bio, avatarURL sql.NullString
	
	err := db.QueryRow(query, id).Scan(
		&profile.ID, &profile.Name, &title, &bio, &avatarURL,
		&profile.CreatedAt, &profile.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	profile.Title = nullStringToPtr(title)
	profile.Bio = nullStringToPtr(bio)
	profile.AvatarURL = nullStringToPtr(avatarURL)
	
	// Load social links
	socialLinks, err := db.GetSocialLinks(profile.ID)
	if err != nil {
		return nil, err
	}
	profile.SocialLinks = socialLinks
	
	return profile, nil
}

func (db *DB) GetDefaultProfile() (*models.Profile, error) {
	query := `SELECT id FROM profiles ORDER BY created_at LIMIT 1`
	var id string
	err := db.QueryRow(query).Scan(&id)
	if err != nil {
		return nil, err
	}
	return db.GetProfile(id)
}

func (db *DB) GetSocialLinks(profileID string) ([]*models.SocialLink, error) {
	query := `SELECT platform, url, icon FROM social_links WHERE profile_id = $1`
	
	rows, err := db.Query(query, profileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var links []*models.SocialLink
	for rows.Next() {
		link := &models.SocialLink{}
		var icon sql.NullString
		
		err := rows.Scan(&link.Platform, &link.URL, &icon)
		if err != nil {
			return nil, err
		}
		
		link.Icon = nullStringToPtr(icon)
		links = append(links, link)
	}
	
	return links, nil
}

// Skills methods
func (db *DB) GetSkills() ([]*models.Skill, error) {
	query := `
		SELECT id, name, category, level, icon_url, display_order, created_at, updated_at
		FROM skills ORDER BY display_order, name
	`
	
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var skills []*models.Skill
	for rows.Next() {
		skill := &models.Skill{}
		var iconURL sql.NullString
		
		err := rows.Scan(
			&skill.ID, &skill.Name, &skill.Category, &skill.Level,
			&iconURL, &skill.DisplayOrder, &skill.CreatedAt, &skill.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		skill.IconURL = nullStringToPtr(iconURL)
		skills = append(skills, skill)
	}
	
	return skills, nil
}

// Experiences methods
func (db *DB) GetExperiences() ([]*models.Experience, error) {
	query := `
		SELECT id, company, position, description, start_date, end_date, 
			   is_current, technologies, created_at, updated_at
		FROM experiences ORDER BY is_current DESC, start_date DESC
	`
	
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var experiences []*models.Experience
	for rows.Next() {
		exp := &models.Experience{}
		var description, endDate sql.NullString
		
		err := rows.Scan(
			&exp.ID, &exp.Company, &exp.Position, &description, &exp.StartDate,
			&endDate, &exp.IsCurrent, pq.Array(&exp.Technologies),
			&exp.CreatedAt, &exp.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		exp.Description = nullStringToPtr(description)
		exp.EndDate = nullStringToPtr(endDate)
		experiences = append(experiences, exp)
	}
	
	return experiences, nil
}

// Blog Posts methods
func (db *DB) GetBlogPosts() ([]*models.BlogPost, error) {
	query := `
		SELECT id, title, slug, excerpt, content, cover_image_url, tags,
			   status, seo_title, seo_description, published_at, like_count, created_at, updated_at
		FROM blog_posts WHERE status = 'PUBLISHED' ORDER BY published_at DESC
	`
	
	return db.queryBlogPosts(query)
}

func (db *DB) GetAdminBlogPosts() ([]*models.BlogPost, error) {
	query := `
		SELECT id, title, slug, excerpt, content, cover_image_url, tags,
			   status, seo_title, seo_description, published_at, like_count, created_at, updated_at
		FROM blog_posts ORDER BY created_at DESC
	`
	
	return db.queryBlogPosts(query)
}

func (db *DB) GetBlogPostBySlug(slug string) (*models.BlogPost, error) {
	query := `
		SELECT id, title, slug, excerpt, content, cover_image_url, tags,
			   status, seo_title, seo_description, published_at, like_count, created_at, updated_at
		FROM blog_posts WHERE slug = $1 AND status = 'PUBLISHED'
	`
	
	posts, err := db.queryBlogPosts(query, slug)
	if err != nil {
		return nil, err
	}
	
	if len(posts) == 0 {
		return nil, nil
	}
	
	return posts[0], nil
}

func (db *DB) queryBlogPosts(query string, args ...interface{}) ([]*models.BlogPost, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var posts []*models.BlogPost
	for rows.Next() {
		post := &models.BlogPost{}
		var excerpt, coverImageURL, seoTitle, seoDescription, publishedAt sql.NullString
		var likeCount sql.NullInt64
		
		err := rows.Scan(
			&post.ID, &post.Title, &post.Slug, &excerpt, &post.Content,
			&coverImageURL, pq.Array(&post.Tags), &post.Status,
			&seoTitle, &seoDescription, &publishedAt, &likeCount,
			&post.CreatedAt, &post.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		post.Excerpt = nullStringToPtr(excerpt)
		post.CoverImageURL = nullStringToPtr(coverImageURL)
		post.SeoTitle = nullStringToPtr(seoTitle)
		post.SeoDescription = nullStringToPtr(seoDescription)
		post.PublishedAt = nullStringToPtr(publishedAt)
		
		if likeCount.Valid {
			count := int(likeCount.Int64)
			post.LikeCount = &count
		} else {
			post.LikeCount = intPtr(0)
		}
		
		posts = append(posts, post)
	}
	
	return posts, nil
}

func (db *DB) GetBlogPostLikeCount(blogPostID string) (int, error) {
	query := `SELECT like_count FROM blog_posts WHERE id = $1`
	var likeCount sql.NullInt64
	
	err := db.QueryRow(query, blogPostID).Scan(&likeCount)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	
	if likeCount.Valid {
		return int(likeCount.Int64), nil
	}
	return 0, nil
}

// Monologues methods
func (db *DB) GetMonologues(limit, offset *int, categoryID *string, tags []string, difficulty *models.Difficulty) ([]*models.Monologue, error) {
	query := `
		SELECT m.id, m.content, m.content_type, m.code_language, m.code_snippet,
			   m.tags, m.is_published, m.published_at, m.url, m.series, m.category,
			   m.difficulty, m.like_count, m.created_at, m.updated_at,
			   c.id, c.name, c.slug, c.description, c.parent_id, c.color, c.icon
		FROM monologues m
		LEFT JOIN code_categories c ON m.code_category_id = c.id
		WHERE m.is_published = true
	`
	args := []interface{}{}
	argIndex := 1
	
	if categoryID != nil {
		query += fmt.Sprintf(" AND m.code_category_id = $%d", argIndex)
		args = append(args, *categoryID)
		argIndex++
	}
	
	if difficulty != nil {
		query += fmt.Sprintf(" AND m.difficulty = $%d", argIndex)
		args = append(args, string(*difficulty))
		argIndex++
	}
	
	if len(tags) > 0 {
		query += fmt.Sprintf(" AND m.tags && $%d", argIndex)
		args = append(args, pq.Array(tags))
		argIndex++
	}
	
	query += " ORDER BY m.published_at DESC"
	
	if limit != nil {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, *limit)
		argIndex++
	}
	
	if offset != nil {
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, *offset)
	}
	
	return db.queryMonologues(query, args...)
}

func (db *DB) GetAdminMonologues() ([]*models.Monologue, error) {
	query := `
		SELECT m.id, m.content, m.content_type, m.code_language, m.code_snippet,
			   m.tags, m.is_published, m.published_at, m.url, m.series, m.category,
			   m.difficulty, m.like_count, m.created_at, m.updated_at,
			   c.id, c.name, c.slug, c.description, c.parent_id, c.color, c.icon
		FROM monologues m
		LEFT JOIN code_categories c ON m.code_category_id = c.id
		ORDER BY m.created_at DESC
	`
	
	return db.queryMonologues(query)
}

func (db *DB) GetMonologueByID(id string) (*models.Monologue, error) {
	query := `
		SELECT m.id, m.content, m.content_type, m.code_language, m.code_snippet,
			   m.tags, m.is_published, m.published_at, m.url, m.series, m.category,
			   m.difficulty, m.like_count, m.created_at, m.updated_at,
			   c.id, c.name, c.slug, c.description, c.parent_id, c.color, c.icon
		FROM monologues m
		LEFT JOIN code_categories c ON m.code_category_id = c.id
		WHERE m.id = $1
	`
	
	monologues, err := db.queryMonologues(query, id)
	if err != nil {
		return nil, err
	}
	
	if len(monologues) == 0 {
		return nil, nil
	}
	
	return monologues[0], nil
}

func (db *DB) queryMonologues(query string, args ...interface{}) ([]*models.Monologue, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var monologues []*models.Monologue
	for rows.Next() {
		mono := &models.Monologue{}
		var codeLanguage, codeSnippet, publishedAt, url, series, category, difficulty sql.NullString
		var likeCount sql.NullInt64
		
		// Code category fields
		var catID, catName, catSlug, catDesc, catParentID, catColor, catIcon sql.NullString
		
		err := rows.Scan(
			&mono.ID, &mono.Content, &mono.ContentType, &codeLanguage, &codeSnippet,
			pq.Array(&mono.Tags), &mono.IsPublished, &publishedAt, &url, &series, &category,
			&difficulty, &likeCount, &mono.CreatedAt, &mono.UpdatedAt,
			&catID, &catName, &catSlug, &catDesc, &catParentID, &catColor, &catIcon,
		)
		if err != nil {
			return nil, err
		}
		
		mono.CodeLanguage = nullStringToPtr(codeLanguage)
		mono.CodeSnippet = nullStringToPtr(codeSnippet)
		mono.PublishedAt = nullStringToPtr(publishedAt)
		mono.URL = nullStringToPtr(url)
		mono.Series = nullStringToPtr(series)
		mono.Category = nullStringToPtr(category)
		
		if difficulty.Valid {
			d := models.Difficulty(difficulty.String)
			mono.Difficulty = &d
		}
		
		if likeCount.Valid {
			count := int(likeCount.Int64)
			mono.LikeCount = &count
		}
		
		// Set code category if present
		if catID.Valid {
			mono.CodeCategory = &models.CodeCategory{
				ID:          catID.String,
				Name:        catName.String,
				Slug:        catSlug.String,
				Description: nullStringToPtr(catDesc),
				ParentID:    nullStringToPtr(catParentID),
				Color:       nullStringToPtr(catColor),
				Icon:        nullStringToPtr(catIcon),
			}
		}
		
		// Load URL preview if URL exists
		if mono.URL != nil {
			urlPreview, _ := db.GetURLPreviewByMonologueID(mono.ID)
			mono.URLPreview = urlPreview
		}
		
		monologues = append(monologues, mono)
	}
	
	return monologues, nil
}

// Code Categories methods
func (db *DB) GetCodeCategories() ([]*models.CodeCategory, error) {
	query := `
		SELECT id, name, slug, description, parent_id, color, icon, created_at, updated_at
		FROM code_categories ORDER BY name
	`
	
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var categories []*models.CodeCategory
	for rows.Next() {
		cat := &models.CodeCategory{}
		var description, parentID, color, icon sql.NullString
		var createdAt, updatedAt time.Time
		
		err := rows.Scan(
			&cat.ID, &cat.Name, &cat.Slug, &description, &parentID,
			&color, &icon, &createdAt, &updatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		cat.Description = nullStringToPtr(description)
		cat.ParentID = nullStringToPtr(parentID)
		cat.Color = nullStringToPtr(color)
		cat.Icon = nullStringToPtr(icon)
		cat.CreatedAt = createdAt
		cat.UpdatedAt = updatedAt
		
		categories = append(categories, cat)
	}
	
	return categories, nil
}

// URL Preview methods
func (db *DB) GetURLPreviewByMonologueID(monologueID string) (*models.URLPreview, error) {
	query := `
		SELECT title, description, image_url, site_name, url, favicon, created_at
		FROM url_previews WHERE monologue_id = $1
	`
	
	preview := &models.URLPreview{}
	var description, imageURL, siteName, favicon sql.NullString
	
	err := db.QueryRow(query, monologueID).Scan(
		&preview.Title, &description, &imageURL, &siteName,
		&preview.URL, &favicon, &preview.CreatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	preview.Description = nullStringToPtr(description)
	preview.ImageURL = nullStringToPtr(imageURL)
	preview.SiteName = nullStringToPtr(siteName)
	preview.Favicon = nullStringToPtr(favicon)
	
	return preview, nil
}

// Helper functions
func nullStringToPtr(ns sql.NullString) *string {
	if !ns.Valid {
		return nil
	}
	return &ns.String
}

func ptrToNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *s, Valid: true}
}

func ptrToNullInt(i *int) sql.NullInt64 {
	if i == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: int64(*i), Valid: true}
}

