package tracktaste

import "github.com/t1nyb0x/jamberry/internal/domain"

// trackResponse はtracktaste APIのトラックレスポンスを表します
type trackResponse struct {
	Album       albumResponse  `json:"album"`
	Artists     []artistBasic  `json:"artists"`
	DiscNumber  int            `json:"disc_number"`
	Popularity  *int           `json:"popularity"`
	ISRC        *string        `json:"isrc"`
	URL         string         `json:"url"`
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	TrackNumber int            `json:"track_number"`
	DurationMs  int            `json:"duration_ms"`
	Explicit    bool           `json:"explicit"`
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
	Album       albumResponse  `json:"album"`
	Artists     []artistBasic  `json:"artists"`
	DiscNumber  int            `json:"disc_number"`
	Popularity  *int           `json:"popularity"`
	ISRC        *string        `json:"isrc"`
	URL         string         `json:"url"`
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	TrackNumber int            `json:"track_number"`
	DurationMs  int            `json:"duration_ms"`
	Explicit    bool           `json:"explicit"`
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

// audioFeaturesResponse はAudio Features情報を表します
type audioFeaturesResponse struct {
	Tempo            float64 `json:"tempo"`
	Energy           float64 `json:"energy"`
	Danceability     float64 `json:"danceability"`
	Valence          float64 `json:"valence"`
	Acousticness     float64 `json:"acousticness"`
	Instrumentalness float64 `json:"instrumentalness"`
	Speechiness      float64 `json:"speechiness"`
	Liveness         float64 `json:"liveness"`
	Loudness         float64 `json:"loudness"`
	Key              int     `json:"key"`
	Mode             int     `json:"mode"`
	TimeSignature    int     `json:"time_signature"`
}

func (a *audioFeaturesResponse) toDomain(trackID string) *domain.AudioFeatures {
	return &domain.AudioFeatures{
		TrackID:          trackID,
		Tempo:            a.Tempo,
		Energy:           a.Energy,
		Danceability:     a.Danceability,
		Valence:          a.Valence,
		Acousticness:     a.Acousticness,
		Instrumentalness: a.Instrumentalness,
		Speechiness:      a.Speechiness,
		Liveness:         a.Liveness,
		Loudness:         a.Loudness,
		Key:              a.Key,
		Mode:             a.Mode,
		TimeSignature:    a.TimeSignature,
	}
}

// seedTrackResponse はシードトラック情報を表します
type seedTrackResponse struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	Artists       []artistBasic          `json:"artists"`
	AudioFeatures *audioFeaturesResponse `json:"audio_features"`
}

func (s *seedTrackResponse) toDomain() domain.Track {
	artists := make([]domain.Artist, len(s.Artists))
	for i, a := range s.Artists {
		artists[i] = a.toDomain()
	}

	return domain.Track{
		ID:      s.ID,
		Name:    s.Name,
		Artists: artists,
	}
}

// recommendTrackResponse はレコメンドトラック情報を表します
type recommendTrackResponse struct {
	Album           albumResponse          `json:"album"`
	Artists         []artistBasic          `json:"artists"`
	ISRC            *string                `json:"isrc"`
	URL             string                 `json:"url"`
	ID              string                 `json:"id"`
	Name            string                 `json:"name"`
	Popularity      *int                   `json:"popularity"`
	TrackNumber     int                    `json:"track_number"`
	DurationMs      int                    `json:"duration_ms"`
	Explicit        bool                   `json:"explicit"`
	SimilarityScore *float64               `json:"similarity_score"`
	MatchReasons    []string               `json:"match_reasons"`
	AudioFeatures   *audioFeaturesResponse `json:"audio_features"`
}

func (t *recommendTrackResponse) toDomain() domain.SimilarTrack {
	artists := make([]domain.Artist, len(t.Artists))
	for i, a := range t.Artists {
		artists[i] = a.toDomain()
	}

	track := domain.SimilarTrack{
		ID:              t.ID,
		Name:            t.Name,
		URL:             t.URL,
		DurationMs:      t.DurationMs,
		TrackNumber:     t.TrackNumber,
		Explicit:        t.Explicit,
		Popularity:      t.Popularity,
		ISRC:            t.ISRC,
		Album:           t.Album.toDomainBasic(),
		Artists:         artists,
		SimilarityScore: t.SimilarityScore,
		MatchReasons:    t.MatchReasons,
	}

	if t.AudioFeatures != nil {
		track.AudioFeatures = t.AudioFeatures.toDomain(t.ID)
	}

	return track
}

// recommendResponse はレコメンドAPIのレスポンス形式を表します
type recommendResponse struct {
	SeedTrack seedTrackResponse        `json:"seed_track"`
	Items     []recommendTrackResponse `json:"items"`
	Mode      string                   `json:"mode"`
}

func (r *recommendResponse) toDomain() *domain.RecommendResult {
	items := make([]domain.SimilarTrack, len(r.Items))
	for i, item := range r.Items {
		items[i] = item.toDomain()
	}

	result := &domain.RecommendResult{
		SeedTrack: r.SeedTrack.toDomain(),
		Items:     items,
		Mode:      domain.RecommendMode(r.Mode),
	}

	if r.SeedTrack.AudioFeatures != nil {
		result.SeedFeatures = r.SeedTrack.AudioFeatures.toDomain(r.SeedTrack.ID)
	}

	return result
}

// similarTrackResponse は類似トラック情報を表します（従来API互換）
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

// similarResponse は類似トラックのレスポンス形式を表します
type similarResponse struct {
	Items []similarTrackResponse `json:"items"`
}

// searchResponse は検索結果のレスポンス形式を表します
type searchResponse struct {
	Items []searchTrackResponse `json:"items"`
}
