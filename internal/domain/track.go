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

// AudioFeatures はトラックの音響特徴量を表します
type AudioFeatures struct {
	TrackID          string
	Tempo            float64
	Energy           float64
	Danceability     float64
	Valence          float64
	Acousticness     float64
	Instrumentalness float64
	Speechiness      float64
	Liveness         float64
	Loudness         float64
	Key              int
	Mode             int
	TimeSignature    int
}

// RecommendMode はレコメンドモードを表します
type RecommendMode string

const (
	RecommendModeSimilar  RecommendMode = "similar"  // 雰囲気重視
	RecommendModeRelated  RecommendMode = "related"  // 関連性重視
	RecommendModeBalanced RecommendMode = "balanced" // バランス（デフォルト）
)

// SimilarTrack は類似トラック情報を表します
type SimilarTrack struct {
	ID              string
	Name            string
	URL             string
	DurationMs      int
	TrackNumber     int
	Explicit        bool
	Popularity      *int
	ISRC            *string
	UPC             *string
	Album           Album
	Artists         []Artist
	SimilarityScore *float64
	MatchReasons    []string
	AudioFeatures   *AudioFeatures
}

// SearchResult は検索結果を表します
type SearchResult struct {
	Tracks []Track
	Total  int
}

// RecommendResult はレコメンド結果を表します
type RecommendResult struct {
	SeedTrack    Track
	SeedFeatures *AudioFeatures
	Items        []SimilarTrack
	Mode         RecommendMode
}
