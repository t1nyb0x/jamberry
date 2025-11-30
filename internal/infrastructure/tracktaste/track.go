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

// similarTrackResponse は類似トラック情報を表します
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
