package models

import (
	"fmt"
	"io"
	"strconv"
	"time"
)

type Skill struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Category     string    `json:"category"`
	Level        int       `json:"level"`
	IconURL      *string   `json:"iconUrl"`
	DisplayOrder int       `json:"displayOrder"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type BlogPost struct {
	ID              string     `json:"id"`
	Title           string     `json:"title"`
	Slug            string     `json:"slug"`
	Excerpt         *string    `json:"excerpt"`
	Content         string     `json:"content"`
	CoverImageURL   *string    `json:"coverImageUrl"`
	Tags            []string   `json:"tags"`
	Status          BlogStatus `json:"status"`
	SeoTitle        *string    `json:"seoTitle"`
	SeoDescription  *string    `json:"seoDescription"`
	PublishedAt     *string    `json:"publishedAt"`
	LikeCount       *int       `json:"likeCount"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
}

type Monologue struct {
	ID               string          `json:"id"`
	Content          string          `json:"content"`
	ContentType      ContentType     `json:"contentType"`
	CodeLanguage     *string         `json:"codeLanguage"`
	CodeSnippet      *string         `json:"codeSnippet"`
	Tags             []string        `json:"tags"`
	IsPublished      bool            `json:"isPublished"`
	PublishedAt      *string         `json:"publishedAt"`
	CreatedAt        time.Time       `json:"createdAt"`
	UpdatedAt        time.Time       `json:"updatedAt"`
	URL              *string         `json:"url"`
	URLPreview       *URLPreview     `json:"urlPreview"`
	RelatedBlogPosts []string        `json:"relatedBlogPosts"`
	Series           *string         `json:"series"`
	Category         *string         `json:"category"`
	LikeCount        *int            `json:"likeCount"`
}


type URLPreview struct {
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	ImageURL    *string   `json:"imageUrl"`
	SiteName    *string   `json:"siteName"`
	URL         string    `json:"url"`
	Favicon     *string   `json:"favicon"`
	CreatedAt   time.Time `json:"createdAt"`
}

// Enums
type ContentType string

const (
	ContentTypePost       ContentType = "POST"
	ContentTypeCode       ContentType = "CODE"
	ContentTypeImage      ContentType = "IMAGE"
	ContentTypeURLPreview ContentType = "URL_PREVIEW"
	ContentTypeBlog       ContentType = "BLOG"
)


type BlogStatus string

const (
	BlogStatusDraft     BlogStatus = "DRAFT"
	BlogStatusPublished BlogStatus = "PUBLISHED"
	BlogStatusArchived  BlogStatus = "ARCHIVED"
)

type Profile struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Title       *string        `json:"title"`
	Bio         *string        `json:"bio"`
	AvatarURL   *string        `json:"avatarUrl"`
	SocialLinks []*SocialLink  `json:"socialLinks"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
}

type SocialLink struct {
	Platform string  `json:"platform"`
	URL      string  `json:"url"`
	Icon     *string `json:"icon"`
}

type SkillCategory struct {
	Category string   `json:"category"`
	Skills   []*Skill `json:"skills"`
}

type Experience struct {
	ID           string    `json:"id"`
	Company      string    `json:"company"`
	Position     string    `json:"position"`
	Description  *string   `json:"description"`
	StartDate    string    `json:"startDate"`
	EndDate      *string   `json:"endDate"`
	IsCurrent    bool      `json:"isCurrent"`
	Technologies []string  `json:"technologies"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type MonologuesResponse struct {
	Nodes       []*Monologue `json:"nodes"`
	TotalCount  int          `json:"totalCount"`
	HasNextPage bool         `json:"hasNextPage"`
}

type RelatedContent struct {
	ID          string      `json:"id"`
	Title       string      `json:"title"`
	Type        ContentType `json:"type"`
	Excerpt     *string     `json:"excerpt"`
	Tags        []string    `json:"tags"`
	PublishedAt string      `json:"publishedAt"`
	ReadTime    *int        `json:"readTime"`
}

type LikeResponse struct {
	ID        string `json:"id"`
	LikeCount int    `json:"likeCount"`
	IsLiked   bool   `json:"isLiked"`
}



type CreateBlogPostInput struct {
	Title          string      `json:"title"`
	Slug           string      `json:"slug"`
	Excerpt        *string     `json:"excerpt"`
	Content        string      `json:"content"`
	CoverImageURL  *string     `json:"coverImageUrl"`
	Tags           []string    `json:"tags"`
	Status         *BlogStatus `json:"status"`
	SeoTitle       *string     `json:"seoTitle"`
	SeoDescription *string     `json:"seoDescription"`
}

type UpdateBlogPostInput struct {
	Title          *string     `json:"title"`
	Slug           *string     `json:"slug"`
	Excerpt        *string     `json:"excerpt"`
	Content        *string     `json:"content"`
	CoverImageURL  *string     `json:"coverImageUrl"`
	Tags           []string    `json:"tags"`
	Status         *BlogStatus `json:"status"`
	SeoTitle       *string     `json:"seoTitle"`
	SeoDescription *string     `json:"seoDescription"`
}

type CreateMonologueInput struct {
	Content        string       `json:"content"`
	ContentType    ContentType  `json:"contentType"`
	CodeLanguage   *string      `json:"codeLanguage"`
	CodeSnippet    *string      `json:"codeSnippet"`
	Tags           []string     `json:"tags"`
	IsPublished    *bool        `json:"isPublished"`
	URL            *string      `json:"url"`
	Series         *string      `json:"series"`
	Category       *string      `json:"category"`
}

type UpdateMonologueInput struct {
	Content        *string      `json:"content"`
	ContentType    *ContentType `json:"contentType"`
	CodeLanguage   *string      `json:"codeLanguage"`
	CodeSnippet    *string      `json:"codeSnippet"`
	Tags           []string     `json:"tags"`
	IsPublished    *bool        `json:"isPublished"`
	URL            *string      `json:"url"`
	Series         *string      `json:"series"`
	Category       *string      `json:"category"`
}

// GraphQL Marshaler methods for enums

func (c ContentType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(string(c)))
}

func (c *ContentType) UnmarshalGQL(v interface{}) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("ContentType must be a string")
	}
	*c = ContentType(s)
	return nil
}


func (b BlogStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(string(b)))
}

func (b *BlogStatus) UnmarshalGQL(v interface{}) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("BlogStatus must be a string")
	}
	*b = BlogStatus(s)
	return nil
}