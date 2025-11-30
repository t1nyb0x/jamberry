package domain

// Artist はアーティストの基本情報を表すドメインエンティティです
type Artist struct {
	ID   string
	Name string
	URL  string
}

// ArtistDetail はアーティストの詳細情報を表します
type ArtistDetail struct {
	ID         string
	Name       string
	URL        string
	Followers  string
	Genres     []string
	Popularity *int
	Images     []Image
}
