package domain

// Track はトラック情報を表すドメインエンティティです
type Track struct {
	ID          string
	Name        string
	URL         string
	DurationMs  int
	TrackNumber int
	DiscNumber  int
	Explicit    bool
	Popularity  *int
	ISRC        *string
	Album       Album
	Artists     []Artist
}

// SimilarTrack は類似トラック情報を表します
type SimilarTrack struct {
	ID          string
	Name        string
	URL         string
	DurationMs  int
	TrackNumber int
	Explicit    bool
	Popularity  *int
	ISRC        *string
	UPC         *string
	Album       Album
}

// SearchResult は検索結果を表します
type SearchResult struct {
	Tracks []Track
	Total  int
}

// RecommendResult はレコメンド結果を表します
type RecommendResult struct {
	SourceTrack   Track
	SimilarTracks []SimilarTrack
}
