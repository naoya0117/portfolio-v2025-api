package resolvers

// THIS CODE WILL BE UPDATED WITH SCHEMA CHANGES. PREVIOUS IMPLEMENTATION FOR SCHEMA CHANGES WILL BE KEPT IN THE COMMENT SECTION. IMPLEMENTATION FOR UNCHANGED SCHEMA WILL BE KEPT.

import (
	"context"
	"fmt"
	"time"

	"github.com/naoya0117/portfolio-v2025-api/internal/generated"
	"github.com/naoya0117/portfolio-v2025-api/internal/models"
	"github.com/naoya0117/portfolio-v2025-api/internal/database"
)

type Resolver struct{ DB *database.DB }

// BlogPost field resolvers
func (r *blogPostResolver) CreatedAt(ctx context.Context, obj *models.BlogPost) (string, error) {
	return obj.CreatedAt.Format(time.RFC3339), nil
}

// UpdatedAt is the resolver for the updatedAt field.
func (r *blogPostResolver) UpdatedAt(ctx context.Context, obj *models.BlogPost) (string, error) {
	return obj.UpdatedAt.Format(time.RFC3339), nil
}

// Monologue field resolvers
func (r *monologueResolver) CreatedAt(ctx context.Context, obj *models.Monologue) (string, error) {
	return obj.CreatedAt.Format(time.RFC3339), nil
}

// UpdatedAt is the resolver for the updatedAt field.
func (r *monologueResolver) UpdatedAt(ctx context.Context, obj *models.Monologue) (string, error) {
	return obj.UpdatedAt.Format(time.RFC3339), nil
}

// Mutation resolvers
func (r *mutationResolver) LikeMonologue(ctx context.Context, id string) (*models.LikeResponse, error) {
	if r.DB == nil {
		return nil, fmt.Errorf("database connection not available")
	}
	return r.DB.LikeMonologue(id)
}

// LikeBlogPost is the resolver for the likeBlogPost field.
func (r *mutationResolver) LikeBlogPost(ctx context.Context, id string) (*models.LikeResponse, error) {
	if r.DB == nil {
		return nil, fmt.Errorf("database connection not available")
	}
	return r.DB.LikeBlogPost(id)
}

// GenerateURLPreview is the resolver for the generateUrlPreview field.
func (r *mutationResolver) GenerateURLPreview(ctx context.Context, url string) (*models.URLPreview, error) {
	// For now, create a simple preview
	// In a real implementation, this would fetch the URL and extract metadata
	return &models.URLPreview{
		Title:       "Generated Preview",
		Description: stringPtr("This is a generated preview for " + url),
		ImageURL:    stringPtr("https://via.placeholder.com/400x200"),
		SiteName:    stringPtr("Example Site"),
		URL:         url,
		Favicon:     stringPtr("https://via.placeholder.com/32x32"),
		CreatedAt:   time.Now(),
	}, nil
}

// CreateBlogPost is the resolver for the createBlogPost field.
func (r *mutationResolver) CreateBlogPost(ctx context.Context, input models.CreateBlogPostInput) (*models.BlogPost, error) {
	if r.DB == nil {
		return nil, fmt.Errorf("database connection not available")
	}
	return r.DB.CreateBlogPost(input)
}

// UpdateBlogPost is the resolver for the updateBlogPost field.
func (r *mutationResolver) UpdateBlogPost(ctx context.Context, id string, input models.UpdateBlogPostInput) (*models.BlogPost, error) {
	if r.DB == nil {
		return nil, fmt.Errorf("database connection not available")
	}
	return r.DB.UpdateBlogPost(id, input)
}

// DeleteBlogPost is the resolver for the deleteBlogPost field.
func (r *mutationResolver) DeleteBlogPost(ctx context.Context, id string) (bool, error) {
	if r.DB == nil {
		return false, fmt.Errorf("database connection not available")
	}
	return r.DB.DeleteBlogPost(id)
}

// PublishBlogPost is the resolver for the publishBlogPost field.
func (r *mutationResolver) PublishBlogPost(ctx context.Context, id string) (*models.BlogPost, error) {
	if r.DB == nil {
		return nil, fmt.Errorf("database connection not available")
	}
	return r.DB.PublishBlogPost(id)
}

// UnpublishBlogPost is the resolver for the unpublishBlogPost field.
func (r *mutationResolver) UnpublishBlogPost(ctx context.Context, id string) (*models.BlogPost, error) {
	if r.DB == nil {
		return nil, fmt.Errorf("database connection not available")
	}
	return r.DB.UnpublishBlogPost(id)
}

// CreateCodeCategory is the resolver for the createCodeCategory field.
func (r *mutationResolver) CreateCodeCategory(ctx context.Context, input models.CreateCodeCategoryInput) (*models.CodeCategory, error) {
	if r.DB == nil {
		return nil, fmt.Errorf("database connection not available")
	}
	return r.DB.CreateCodeCategory(input)
}

// UpdateCodeCategory is the resolver for the updateCodeCategory field.
func (r *mutationResolver) UpdateCodeCategory(ctx context.Context, id string, input models.UpdateCodeCategoryInput) (*models.CodeCategory, error) {
	if r.DB == nil {
		return nil, fmt.Errorf("database connection not available")
	}
	return r.DB.UpdateCodeCategory(id, input)
}

// DeleteCodeCategory is the resolver for the deleteCodeCategory field.
func (r *mutationResolver) DeleteCodeCategory(ctx context.Context, id string) (bool, error) {
	if r.DB == nil {
		return false, fmt.Errorf("database connection not available")
	}
	return r.DB.DeleteCodeCategory(id)
}

// CreateMonologue is the resolver for the createMonologue field.
func (r *mutationResolver) CreateMonologue(ctx context.Context, input models.CreateMonologueInput) (*models.Monologue, error) {
	fmt.Printf("[RESOLVER] CreateMonologue called with input: %+v\n", input)
	fmt.Printf("[RESOLVER] ContentType: %s\n", input.ContentType)
	if r.DB == nil {
		return nil, fmt.Errorf("database connection not available")
	}
	result, err := r.DB.CreateMonologue(input)
	if err != nil {
		fmt.Printf("[RESOLVER] CreateMonologue error: %v\n", err)
		return nil, err
	}
	fmt.Printf("[RESOLVER] CreateMonologue success: %+v\n", result)
	return result, nil
}

// UpdateMonologue is the resolver for the updateMonologue field.
func (r *mutationResolver) UpdateMonologue(ctx context.Context, id string, input models.UpdateMonologueInput) (*models.Monologue, error) {
	if r.DB == nil {
		return nil, fmt.Errorf("database connection not available")
	}
	return r.DB.UpdateMonologue(id, input)
}

// DeleteMonologue is the resolver for the deleteMonologue field.
func (r *mutationResolver) DeleteMonologue(ctx context.Context, id string) (bool, error) {
	if r.DB == nil {
		return false, fmt.Errorf("database connection not available")
	}
	return r.DB.DeleteMonologue(id)
}

// PublishMonologue is the resolver for the publishMonologue field.
func (r *mutationResolver) PublishMonologue(ctx context.Context, id string) (*models.Monologue, error) {
	if r.DB == nil {
		return nil, fmt.Errorf("database connection not available")
	}
	return r.DB.PublishMonologue(id)
}

// UnpublishMonologue is the resolver for the unpublishMonologue field.
func (r *mutationResolver) UnpublishMonologue(ctx context.Context, id string) (*models.Monologue, error) {
	if r.DB == nil {
		return nil, fmt.Errorf("database connection not available")
	}
	return r.DB.UnpublishMonologue(id)
}

// Query resolvers
func (r *queryResolver) Profile(ctx context.Context) (*models.Profile, error) {
	if r.DB == nil {
		return nil, fmt.Errorf("database connection not available")
	}
	return r.DB.GetDefaultProfile()
}

// Skills is the resolver for the skills field.
func (r *queryResolver) Skills(ctx context.Context) ([]*models.Skill, error) {
	if r.DB == nil {
		return nil, fmt.Errorf("database connection not available")
	}
	return r.DB.GetSkills()
}

// SkillsByCategory is the resolver for the skillsByCategory field.
func (r *queryResolver) SkillsByCategory(ctx context.Context) ([]*models.SkillCategory, error) {
	if r.DB == nil {
		return nil, fmt.Errorf("database connection not available")
	}

	skills, err := r.DB.GetSkills()
	if err != nil {
		return nil, err
	}

	categories := make(map[string][]*models.Skill)
	for _, skill := range skills {
		categories[skill.Category] = append(categories[skill.Category], skill)
	}

	result := make([]*models.SkillCategory, 0)
	for category, categorySkills := range categories {
		result = append(result, &models.SkillCategory{
			Category: category,
			Skills:   categorySkills,
		})
	}

	return result, nil
}

// Experiences is the resolver for the experiences field.
func (r *queryResolver) Experiences(ctx context.Context) ([]*models.Experience, error) {
	if r.DB == nil {
		return nil, fmt.Errorf("database connection not available")
	}
	return r.DB.GetExperiences()
}

// Monologue is the resolver for the monologue field.
func (r *queryResolver) Monologue(ctx context.Context, id string) (*models.Monologue, error) {
	if r.DB == nil {
		return nil, fmt.Errorf("database connection not available")
	}
	return r.DB.GetMonologueByID(id)
}

// Monologues is the resolver for the monologues field.
func (r *queryResolver) Monologues(ctx context.Context, limit *int, offset *int, categoryID *string, tags []string, difficulty *models.Difficulty) (*models.MonologuesResponse, error) {
	if r.DB == nil {
		return nil, fmt.Errorf("database connection not available")
	}

	monologues, err := r.DB.GetMonologues(limit, offset, categoryID, tags, difficulty)
	if err != nil {
		return nil, err
	}

	// Get total count for pagination
	allMonologues, err := r.DB.GetMonologues(nil, nil, categoryID, tags, difficulty)
	if err != nil {
		return nil, err
	}

	hasNextPage := false
	if limit != nil && offset != nil {
		hasNextPage = len(allMonologues) > (*offset + *limit)
	}

	return &models.MonologuesResponse{
		Nodes:       monologues,
		TotalCount:  len(allMonologues),
		HasNextPage: hasNextPage,
	}, nil
}

// BlogPost is the resolver for the blogPost field.
func (r *queryResolver) BlogPost(ctx context.Context, slug string) (*models.BlogPost, error) {
	if r.DB == nil {
		return nil, fmt.Errorf("database connection not available")
	}
	return r.DB.GetBlogPostBySlug(slug)
}

// BlogPosts is the resolver for the blogPosts field.
func (r *queryResolver) BlogPosts(ctx context.Context) ([]*models.BlogPost, error) {
	if r.DB == nil {
		return nil, fmt.Errorf("database connection not available")
	}
	return r.DB.GetBlogPosts()
}

// AdminBlogPosts is the resolver for the adminBlogPosts field.
func (r *queryResolver) AdminBlogPosts(ctx context.Context) ([]*models.BlogPost, error) {
	if r.DB == nil {
		return nil, fmt.Errorf("database connection not available")
	}
	return r.DB.GetAdminBlogPosts()
}

// AdminMonologues is the resolver for the adminMonologues field.
func (r *queryResolver) AdminMonologues(ctx context.Context) ([]*models.Monologue, error) {
	if r.DB == nil {
		return nil, fmt.Errorf("database connection not available")
	}
	return r.DB.GetAdminMonologues()
}

// CodeCategories is the resolver for the codeCategories field.
func (r *queryResolver) CodeCategories(ctx context.Context) ([]*models.CodeCategory, error) {
	if r.DB == nil {
		return nil, fmt.Errorf("database connection not available")
	}
	return r.DB.GetCodeCategories()
}

// CodeCategoriesHierarchy is the resolver for the codeCategoriesHierarchy field.
func (r *queryResolver) CodeCategoriesHierarchy(ctx context.Context) ([]*models.CodeCategory, error) {
	if r.DB == nil {
		return nil, fmt.Errorf("database connection not available")
	}
	return r.DB.GetCodeCategories()
}

// RelatedContent is the resolver for the relatedContent field.
func (r *queryResolver) RelatedContent(ctx context.Context, monologueID string, limit *int) ([]*models.RelatedContent, error) {
	if r.DB == nil {
		return nil, fmt.Errorf("database connection not available")
	}

	// Get the current monologue
	currentMonologue, err := r.DB.GetMonologueByID(monologueID)
	if err != nil || currentMonologue == nil {
		return []*models.RelatedContent{}, nil
	}

	// Get all published monologues and blog posts
	allMonos, err := r.DB.GetMonologues(nil, nil, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	allPosts, err := r.DB.GetBlogPosts()
	if err != nil {
		return nil, err
	}

	result := make([]*models.RelatedContent, 0)
	maxLimit := 6
	if limit != nil && *limit < maxLimit {
		maxLimit = *limit
	}

	// Find related monologues by tag similarity
	for _, m := range allMonos {
		if m.ID == monologueID || !m.IsPublished || len(result) >= maxLimit {
			continue
		}

		// Check for tag overlap
		hasCommonTag := false
		for _, tag1 := range currentMonologue.Tags {
			for _, tag2 := range m.Tags {
				if tag1 == tag2 {
					hasCommonTag = true
					break
				}
			}
			if hasCommonTag {
				break
			}
		}

		if hasCommonTag {
			titleLen := min(50, len(m.Content))
			excerptLen := min(100, len(m.Content))
			result = append(result, &models.RelatedContent{
				ID:          m.ID,
				Title:       m.Content[:titleLen] + "...",
				Type:        models.ContentTypePost,
				Excerpt:     stringPtr(m.Content[:excerptLen] + "..."),
				Tags:        m.Tags,
				PublishedAt: *m.PublishedAt,
				ReadTime:    intPtr(len(m.Content)/200 + 1), // Estimated read time based on content length
			})
		}
	}

	// Add related blog posts
	for _, post := range allPosts {
		if len(result) >= maxLimit {
			break
		}

		// Check for tag overlap with blog posts
		hasCommonTag := false
		for _, tag1 := range currentMonologue.Tags {
			for _, tag2 := range post.Tags {
				if tag1 == tag2 {
					hasCommonTag = true
					break
				}
			}
			if hasCommonTag {
				break
			}
		}

		if hasCommonTag && post.Status == models.BlogStatusPublished {
			result = append(result, &models.RelatedContent{
				ID:          post.ID,
				Title:       post.Title,
				Type:        models.ContentTypePost,
				Excerpt:     post.Excerpt,
				Tags:        post.Tags,
				PublishedAt: *post.PublishedAt,
				ReadTime:    intPtr(len(post.Content)/200 + 1), // Estimated read time based on content length
			})
		}
	}

	return result, nil
}

// URL Preview resolver
func (r *urlPreviewResolver) CreatedAt(ctx context.Context, obj *models.URLPreview) (string, error) {
	return obj.CreatedAt.Format(time.RFC3339), nil
}

// BlogPost returns generated.BlogPostResolver implementation.
func (r *Resolver) BlogPost() generated.BlogPostResolver { return &blogPostResolver{r} }

// Monologue returns generated.MonologueResolver implementation.
func (r *Resolver) Monologue() generated.MonologueResolver { return &monologueResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// UrlPreview returns generated.UrlPreviewResolver implementation.
func (r *Resolver) UrlPreview() generated.UrlPreviewResolver { return &urlPreviewResolver{r} }

type blogPostResolver struct{ *Resolver }
type monologueResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type urlPreviewResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
/*
	type Resolver struct {
	DB *database.DB
}
func (r *blogPostResolver) LikeCount(ctx context.Context, obj *models.BlogPost) (*int, error) {
	if r.DB == nil {
		return intPtr(0), nil
	}
	count, err := r.DB.GetBlogPostLikeCount(obj.ID)
	if err != nil {
		return intPtr(0), nil
	}
	return intPtr(count), nil
}
func stringPtr(s string) *string {
	return &s
}
func intPtr(i int) *int {
	return &i
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
*/

func (r *blogPostResolver) LikeCount(ctx context.Context, obj *models.BlogPost) (*int, error) {
	if r.DB == nil {
		return intPtr(0), nil
	}
	count, err := r.DB.GetBlogPostLikeCount(obj.ID)
	if err  !=  nil {
		return intPtr(0), nil
	}
	return intPtr(count), nil
}

func intPtr(i int) *int { return &i }
func stringPtr(s string) *string { return &s }
