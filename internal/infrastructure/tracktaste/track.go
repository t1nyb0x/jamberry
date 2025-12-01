package tracktaste

import "github.com/t1nyb0x/jamberry/internal/domain"

// trackResponse はtracktaste APIのトラックレスポンスを表します
type trackResponse struct {
	Album       albumResponse `json:"album"`
	Artists     []artistBasic `json:"artists"`
	DiscNumber  int           `json:"disc_number"`
	Popularity  *int          `json:"popularity"`
	ISRC        *string       `json:"isrc"`
	URL         string        `json:"url"`
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	TrackNumber int           `json:"track_number"`
	DurationMs  int           `json:"duration_ms"`
	Explicit    bool          `json:"explicit"`
}

func (t *trackResponse) toDomain() *domain.Track {
	artists := make([]domain.Artist, len(t.Artists))
	for i, a := range t.Artists {
		artists[i] = a.toDomain()
	}

	return &domain.Track{
		ID:          t.ID,
		Name:        t.Name,
		URL:         t.URL,
		DurationMs:  t.DurationMs,
		TrackNumber: t.TrackNumber,
		DiscNumber:  t.DiscNumber,
		Explicit:    t.Explicit,
		Popularity:  t.Popularity,
		ISRC:        t.ISRC,
		Album:       t.Album.toDomainBasic(),
		Artists:     artists,
	}
}

// searchTrackResponse は検索結果のトラック情報を表します
type searchTrackResponse struct {
	Album       albumResponse `json:"album"`
	Artists     []artistBasic `json:"artists"`
	DiscNumber  int           `json:"disc_number"`
	Popularity  *int          `json:"popularity"`
	ISRC        *string       `json:"isrc"`
	URL         string        `json:"url"`
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	TrackNumber int           `json:"track_number"`
	DurationMs  int           `json:"duration_ms"`
	Explicit    bool          `json:"explicit"`
}

func (t *searchTrackResponse) toDomain() domain.Track {
	artists := make([]domain.Artist, len(t.Artists))
	for i, a := range t.Artists {
		artists[i] = a.toDomain()
	}

	return domain.Track{
		ID:          t.ID,
		Name:        t.Name,
		URL:         t.URL,
		DurationMs:  t.DurationMs,
		TrackNumber: t.TrackNumber,
		DiscNumber:  t.DiscNumber,
		Explicit:    t.Explicit,
		Popularity:  t.Popularity,
		ISRC:        t.ISRC,
		Album:       t.Album.toDomainBasic(),
		Artists:     artists,
	}
}

// trackFeaturesResponse はTrackFeatures情報を表します（v2: Deezer + MusicBrainz）
type trackFeaturesResponse struct {
	BPM             float64  `json:"bpm"`
	DurationSeconds int      `json:"duration_seconds"`
	Gain            float64  `json:"gain"`
	Tags            []string `json:"tags"`
}

func (f *trackFeaturesResponse) toDomain() *domain.TrackFeatures {
	return &domain.TrackFeatures{
		BPM:             f.BPM,
		DurationSeconds: f.DurationSeconds,
		Gain:            f.Gain,
		Tags:            f.Tags,
	}
}

// seedTrackResponseV2 はシードトラック情報を表します（v2）
type seedTrackResponseV2 struct {
	ID      string         `json:"id"`
	Name    string         `json:"name"`
	Artists []artistBasic  `json:"artists"` // 実際のAPIは複数形
	Album   *albumResponse `json:"album"`
}

func (s *seedTrackResponseV2) toDomain() domain.Track {
	track := domain.Track{
		ID:   s.ID,
		Name: s.Name,
	}

	for _, artist := range s.Artists {
		track.Artists = append(track.Artists, artist.toDomain())
	}

	if s.Album != nil {
		track.Album = s.Album.toDomainBasic()
	}

	return track
}

// recommendTrackResponseV2 はレコメンドトラック情報を表します（v2）
// 実際のAPIレスポンスでは、トラック情報がフラットに返される
type recommendTrackResponseV2 struct {
	ID              string                 `json:"id"`
	Name            string                 `json:"name"`
	Artists         []artistBasic          `json:"artists"`
	SimilarityScore float64                `json:"similarity_score"`
	MatchReasons    []string               `json:"match_reasons"`
	Features        *trackFeaturesResponse `json:"features"`
}

func (t *recommendTrackResponseV2) toDomain() domain.SimilarTrack {
	track := domain.SimilarTrack{
		ID:              t.ID,
		Name:            t.Name,
		SimilarityScore: &t.SimilarityScore,
		MatchReasons:    t.MatchReasons,
	}

	for _, artist := range t.Artists {
		track.Artists = append(track.Artists, artist.toDomain())
	}

	if t.Features != nil {
		track.Features = t.Features.toDomain()
	}

	return track
}

// recommendResponseV2 はレコメンドAPIのレスポンス形式を表します（v2）
type recommendResponseV2 struct {
	SeedTrack    seedTrackResponseV2        `json:"seed_track"`
	SeedFeatures *trackFeaturesResponse     `json:"seed_features"`
	Items        []recommendTrackResponseV2 `json:"items"`
	Mode         string                     `json:"mode"`
}

func (r *recommendResponseV2) toDomain() *domain.RecommendResult {
	items := make([]domain.SimilarTrack, len(r.Items))
	for i, item := range r.Items {
		items[i] = item.toDomain()
	}

	result := &domain.RecommendResult{
		SeedTrack: r.SeedTrack.toDomain(),
		Items:     items,
		Mode:      domain.RecommendMode(r.Mode),
	}

	if r.SeedFeatures != nil {
		result.SeedFeatures = r.SeedFeatures.toDomain()
	}

	return result
}

// ============================================
// 以下は旧仕様（v1）のレスポンス構造体
// 注意: /v1/track/similar API用（従来API互換）
// ============================================

// similarTrackResponse は類似トラック情報を表します（旧仕様: v1 API）
type similarTrackResponse struct {
	Album       albumResponse `json:"album"`
	ISRC        *string       `json:"isrc"`
	UPC         *string       `json:"upc"`
	URL         string        `json:"url"`
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Popularity  *int          `json:"popularity"`
	TrackNumber int           `json:"track_number"`
	DurationMs  int           `json:"duration_ms"`
	Explicit    bool          `json:"explicit"`
}

func (t *similarTrackResponse) toDomain() domain.SimilarTrack {
	return domain.SimilarTrack{
		ID:          t.ID,
		Name:        t.Name,
		URL:         t.URL,
		DurationMs:  t.DurationMs,
		TrackNumber: t.TrackNumber,
		Explicit:    t.Explicit,
		Popularity:  t.Popularity,
		ISRC:        t.ISRC,
		UPC:         t.UPC,
		Album:       t.Album.toDomainBasic(),
	}
}

// similarResponse は類似トラックのレスポンス形式を表します（旧仕様: v1 API）
type similarResponse struct {
	Items []similarTrackResponse `json:"items"`
}

// searchResponse は検索結果のレスポンス形式を表します
type searchResponse struct {
	Items []searchTrackResponse `json:"items"`
}
