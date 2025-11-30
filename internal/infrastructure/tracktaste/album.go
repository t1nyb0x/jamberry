package tracktaste

import "github.com/t1nyb0x/jamberry/internal/domain"

// albumResponse はアルバム情報を表します
type albumResponse struct {
	URL         string          `json:"url"`
	ID          string          `json:"id"`
	Images      []imageResponse `json:"images"`
	Name        string          `json:"name"`
	ReleaseDate string          `json:"release_date"`
	Artists     []artistBasic   `json:"artists"`
	Tracks      *albumTracks    `json:"tracks,omitempty"`
	Popularity  *int            `json:"popularity,omitempty"`
	UPC         string          `json:"upc,omitempty"`
	Genres      []string        `json:"genres,omitempty"`
}

func (a *albumResponse) toDomainBasic() domain.Album {
	images := make([]domain.Image, len(a.Images))
	for i, img := range a.Images {
		images[i] = img.toDomain()
	}

	artists := make([]domain.Artist, len(a.Artists))
	for i, ar := range a.Artists {
		artists[i] = ar.toDomain()
	}

	return domain.Album{
		ID:          a.ID,
		Name:        a.Name,
		URL:         a.URL,
		ReleaseDate: a.ReleaseDate,
		Images:      images,
		Artists:     artists,
	}
}

func (a *albumResponse) toDomain() *domain.AlbumDetail {
	images := make([]domain.Image, len(a.Images))
	for i, img := range a.Images {
		images[i] = img.toDomain()
	}

	artists := make([]domain.Artist, len(a.Artists))
	for i, ar := range a.Artists {
		artists[i] = ar.toDomain()
	}

	var tracks []domain.AlbumTrack
	if a.Tracks != nil {
		tracks = make([]domain.AlbumTrack, len(a.Tracks.Items))
		for i, t := range a.Tracks.Items {
			tracks[i] = t.toDomain()
		}
	}

	return &domain.AlbumDetail{
		ID:          a.ID,
		Name:        a.Name,
		URL:         a.URL,
		ReleaseDate: a.ReleaseDate,
		Images:      images,
		Artists:     artists,
		Tracks:      tracks,
		Popularity:  a.Popularity,
		UPC:         a.UPC,
		Genres:      a.Genres,
	}
}

// albumTracks はアルバム内のトラックリストを表します
type albumTracks struct {
	Items []albumTrackResponse `json:"items"`
}

// albumTrackResponse はアルバム内のトラック情報を表します
type albumTrackResponse struct {
	Artists     []artistBasic `json:"artists"`
	URL         string        `json:"url"`
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	TrackNumber int           `json:"track_number"`
}

func (t *albumTrackResponse) toDomain() domain.AlbumTrack {
	artists := make([]domain.Artist, len(t.Artists))
	for i, a := range t.Artists {
		artists[i] = a.toDomain()
	}

	return domain.AlbumTrack{
		ID:          t.ID,
		Name:        t.Name,
		URL:         t.URL,
		TrackNumber: t.TrackNumber,
		Artists:     artists,
	}
}
