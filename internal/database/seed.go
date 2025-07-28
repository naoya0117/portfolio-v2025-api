package database

import (
	"log"

	"github.com/lib/pq"
)

// SeedData populates the database with initial data
func (db *DB) SeedData() error {
	log.Println("Starting database seeding...")

	// Check if data already exists
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM profiles").Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		log.Println("Database already has data, skipping seeding")
		return nil
	}

	// Seed Profile
	profileID, err := db.seedProfile()
	if err != nil {
		return err
	}

	// Seed Skills
	if err := db.seedSkills(); err != nil {
		return err
	}

	// Seed Experiences
	if err := db.seedExperiences(); err != nil {
		return err
	}

	// Seed Code Categories
	categoryIDs, err := db.seedCodeCategories()
	if err != nil {
		return err
	}

	// Seed Blog Posts
	if err := db.seedBlogPosts(); err != nil {
		return err
	}

	// Seed Monologues
	if err := db.seedMonologues(categoryIDs); err != nil {
		return err
	}

	log.Printf("Database seeded successfully with profile ID: %s", profileID)
	return nil
}

func (db *DB) seedProfile() (string, error) {
	query := `
		INSERT INTO profiles (name, title, bio, avatar_url)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	var profileID string
	err := db.QueryRow(
		query,
		"山田太郎",
		"フルスタック開発者",
		"Next.js、React、TypeScript、Goを使った現代的なWebアプリケーション開発を専門としています。ユーザー体験とパフォーマンスを重視したソリューションの提供に情熱を注いでいます。",
		"/api/placeholder/150/150",
	).Scan(&profileID)

	if err != nil {
		return "", err
	}

	// Seed social links
	socialLinks := []struct {
		platform string
		url      string
		icon     string
	}{
		{"GitHub", "https://github.com", "github"},
		{"Twitter", "https://twitter.com", "twitter"},
		{"LinkedIn", "https://linkedin.com", "linkedin"},
	}

	for _, link := range socialLinks {
		_, err := db.Exec(
			"INSERT INTO social_links (profile_id, platform, url, icon) VALUES ($1, $2, $3, $4)",
			profileID, link.platform, link.url, link.icon,
		)
		if err != nil {
			return "", err
		}
	}

	return profileID, nil
}

func (db *DB) seedSkills() error {
	skills := []struct {
		name         string
		category     string
		level        int
		displayOrder int
	}{
		{"React", "Frontend", 9, 1},
		{"Next.js", "Frontend", 8, 2},
		{"TypeScript", "Language", 9, 3},
		{"JavaScript", "Language", 9, 4},
		{"Go", "Backend", 7, 5},
		{"Node.js", "Backend", 8, 6},
		{"PostgreSQL", "Database", 7, 7},
		{"GraphQL", "API", 8, 8},
		{"Tailwind CSS", "Styling", 9, 9},
		{"Docker", "DevOps", 7, 10},
		{"AWS", "Cloud", 6, 11},
		{"Git", "Tools", 9, 12},
	}

	for _, skill := range skills {
		_, err := db.Exec(
			"INSERT INTO skills (name, category, level, display_order) VALUES ($1, $2, $3, $4)",
			skill.name, skill.category, skill.level, skill.displayOrder,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *DB) seedExperiences() error {
	experiences := []struct {
		company      string
		position     string
		description  string
		startDate    string
		endDate      *string
		isCurrent    bool
		technologies []string
	}{
		{
			"テック株式会社",
			"シニア フルスタック開発者",
			"大規模Webアプリケーションの設計・開発・運用を担当。React、Next.js、Goでのマイクロサービスアーキテクチャによるシステム構築をリード。",
			"2022-04",
			nil,
			true,
			[]string{"React", "Next.js", "Go", "PostgreSQL", "Docker", "AWS"},
		},
		{
			"スタートアップ合同会社",
			"フロントエンド開発者",
			"SaaSプロダクトのフロントエンド開発を担当。Vue.js からReactへの移行プロジェクトをリード。UIコンポーネントライブラリの設計・構築。",
			"2020-06",
			seedStringPtr("2022-03"),
			false,
			[]string{"React", "Vue.js", "TypeScript", "Storybook", "Jest"},
		},
		{
			"システム開発会社",
			"ジュニア開発者",
			"受託開発プロジェクトにて、PHP、JavaScriptでのWebアプリケーション開発に従事。基礎的な開発スキルを習得。",
			"2018-04",
			seedStringPtr("2020-05"),
			false,
			[]string{"PHP", "JavaScript", "MySQL", "jQuery", "Bootstrap"},
		},
	}

	for _, exp := range experiences {
		_, err := db.Exec(`
			INSERT INTO experiences (company, position, description, start_date, end_date, is_current, technologies)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, exp.company, exp.position, exp.description, exp.startDate, exp.endDate, exp.isCurrent, pq.Array(exp.technologies))
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *DB) seedCodeCategories() (map[string]string, error) {
	categories := []struct {
		name        string
		slug        string
		description string
		color       string
		icon        string
	}{
		{"React Hooks", "react-hooks", "React Hooksの使い方とパターン", "#61dafb", "⚛️"},
		{"TypeScript Tips", "typescript-tips", "TypeScript活用のコツとベストプラクティス", "#3178c6", "🔷"},
		{"パフォーマンス最適化", "performance-optimization", "Webアプリのパフォーマンス改善テクニック", "#ff6b6b", "⚡"},
		{"API設計", "api-design", "REST APIとGraphQLの設計パターン", "#4ecdc4", "🔌"},
	}

	categoryIDs := make(map[string]string)

	for _, cat := range categories {
		var id string
		err := db.QueryRow(`
			INSERT INTO code_categories (name, slug, description, color, icon)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id
		`, cat.name, cat.slug, cat.description, cat.color, cat.icon).Scan(&id)

		if err != nil {
			return nil, err
		}

		categoryIDs[cat.slug] = id
	}

	return categoryIDs, nil
}

func (db *DB) seedBlogPosts() error {
	posts := []struct {
		title       string
		slug        string
		excerpt     string
		content     string
		tags        []string
		status      string
		publishedAt *string
	}{
		{
			"Next.js 15で変わったこと",
			"nextjs-15-changes",
			"Next.js 15の新機能と変更点について詳しく解説します。",
			"# Next.js 15で変わったこと\n\nNext.js 15がリリースされ、多くの新機能と改善が加えられました...",
			[]string{"Next.js", "React", "Web Development"},
			"PUBLISHED",
			seedStringPtr("2025-01-01T09:00:00Z"),
		},
		{
			"TypeScriptの型システムを理解する",
			"understanding-typescript-type-system",
			"TypeScriptの型システムの基礎から応用まで、実例とともに学びます。",
			"# TypeScriptの型システムを理解する\n\nTypeScriptの型システムは強力で、適切に使用することで...",
			[]string{"TypeScript", "Programming", "Type Safety"},
			"PUBLISHED",
			seedStringPtr("2024-12-15T14:00:00Z"),
		},
	}

	for _, post := range posts {
		_, err := db.Exec(`
			INSERT INTO blog_posts (title, slug, excerpt, content, tags, status, published_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, post.title, post.slug, post.excerpt, post.content, pq.Array(post.tags), post.status, post.publishedAt)

		if err != nil {
			return err
		}
	}

	return nil
}

func (db *DB) seedMonologues(categoryIDs map[string]string) error {
	monologues := []struct {
		content         string
		contentType     string
		codeLanguage    *string
		codeSnippet     *string
		tags            []string
		isPublished     bool
		publishedAt     *string
		url             *string
		category        string
		codeCategoryID  *string
		difficulty      *string
		likeCount       int
	}{
		{
			"React 19の新機能について調べていたところ、use()フックの存在を知りました。Promiseを直接扱えるようになるのは便利そうです。",
			"POST",
			nil,
			nil,
			[]string{"React", "JavaScript"},
			true,
			seedStringPtr("2025-01-15T10:00:00Z"),
			nil,
			"技術メモ",
			nil,
			nil,
			12,
		},
		{
			"カスタムフックでデータフェッチングを抽象化する方法",
			"CODE",
			seedStringPtr("typescript"),
			seedStringPtr(`const useAsyncData = <T>(asyncFn: () => Promise<T>) => {
  const [data, setData] = useState<T | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<Error | null>(null)

  useEffect(() => {
    asyncFn()
      .then(setData)
      .catch(setError)
      .finally(() => setLoading(false))
  }, [])

  return { data, loading, error }
}`),
			[]string{"React", "TypeScript", "Hooks"},
			true,
			seedStringPtr("2025-01-10T15:30:00Z"),
			nil,
			"",
			seedStringPtr(categoryIDs["react-hooks"]),
			seedStringPtr("INTERMEDIATE"),
			24,
		},
		{
			"便利なReact Hooksライブラリを見つけました。https://github.com/streamich/react-use には様々なカスタムフックが用意されていて、開発が効率的になりそうです。",
			"POST",
			nil,
			nil,
			[]string{"React", "Hooks", "ライブラリ"},
			true,
			seedStringPtr("2025-01-12T09:00:00Z"),
			seedStringPtr("https://github.com/streamich/react-use"),
			"ツール紹介",
			nil,
			nil,
			8,
		},
	}

	for _, mono := range monologues {
		var id string
		err := db.QueryRow(`
			INSERT INTO monologues (content, content_type, code_language, code_snippet, tags, is_published, 
								   published_at, url, category, code_category_id, difficulty, like_count)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
			RETURNING id
		`, mono.content, mono.contentType, mono.codeLanguage, mono.codeSnippet, pq.Array(mono.tags),
			mono.isPublished, mono.publishedAt, mono.url, mono.category, mono.codeCategoryID,
			mono.difficulty, mono.likeCount).Scan(&id)

		if err != nil {
			return err
		}

		// Create URL preview if URL exists
		if mono.url != nil {
			_, err = db.Exec(`
				INSERT INTO url_previews (monologue_id, title, description, image_url, site_name, url, favicon)
				VALUES ($1, $2, $3, $4, $5, $6, $7)
			`, id, "react-use", "Collection of essential React Hooks",
				"https://repository-images.githubusercontent.com/146641387/38ba6700-5db6-11ea-8af8-b5b0c92e5e2b",
				"GitHub", *mono.url, "https://github.githubassets.com/favicons/favicon.svg")

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func seedStringPtr(s string) *string {
	return &s
}