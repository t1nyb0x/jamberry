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

// AudioFeatures はトラックの音響特徴量を表します（旧仕様: Spotify Audio Features）
// 注意: Spotify Audio Features API は2024年11月に廃止されました
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

// TrackFeatures はトラックの特徴量を表します（v2: Deezer + MusicBrainz）
type TrackFeatures struct {
	BPM             float64  // Deezer: テンポ (0-250)
	DurationSeconds int      // Deezer: 曲の長さ（秒）
	Gain            float64  // Deezer: ReplayGain (dB)
	Tags            []string // MusicBrainz: ジャンル/スタイルタグ
}

// RecommendMode はレコメンドモードを表します
type RecommendMode string

const (
	RecommendModeSimilar  RecommendMode = "similar"  // 雰囲気重視（Deezer特徴量重視）
	RecommendModeRelated  RecommendMode = "related"  // 関連性重視（MusicBrainzタグ重視）
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
	GenreBonus      *float64       // v2: ジャンルボーナス倍率
	FinalScore      *float64       // v2: 最終スコア (similarity_score * genre_bonus)
	MatchReasons    []string
	AudioFeatures   *AudioFeatures // 旧仕様（v1）
	Features        *TrackFeatures // 新仕様（v2）
}

// SearchResult は検索結果を表します
type SearchResult struct {
	Tracks []Track
	Total  int
}

// RecommendResult はレコメンド結果を表します
type RecommendResult struct {
	SeedTrack    Track
	SeedFeatures *TrackFeatures // v2: Deezer + MusicBrainz features
	Items        []SimilarTrack
	Mode         RecommendMode
}
