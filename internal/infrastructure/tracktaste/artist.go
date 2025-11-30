package tracktaste

import "github.com/t1nyb0x/jamberry/internal/domain"

// artistBasic はアーティストの基本情報を表します
type artistBasic struct {
	URL  string `json:"url"`
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (a *artistBasic) toDomain() domain.Artist {
	return domain.Artist{
		ID:   a.ID,
		Name: a.Name,
		URL:  a.URL,
	}
}

// artistResponse はアーティストの詳細情報を表します
type artistResponse struct {
	URL        string          `json:"url"`
	Followers  string          `json:"followers"`
	Genres     []string        `json:"genres"`
	ID         string          `json:"id"`
	Images     []imageResponse `json:"images"`
	Name       string          `json:"name"`
	Popularity *int            `json:"popularity"`
}

func (a *artistResponse) toDomain() *domain.ArtistDetail {
	images := make([]domain.Image, len(a.Images))
	for i, img := range a.Images {
		images[i] = img.toDomain()
	}

	return &domain.ArtistDetail{
		ID:         a.ID,
		Name:       a.Name,
		URL:        a.URL,
		Followers:  a.Followers,
		Genres:     a.Genres,
		Popularity: a.Popularity,
		Images:     images,
	}
}

// imageResponse は画像情報を表します
type imageResponse struct {
	URL    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

func (i *imageResponse) toDomain() domain.Image {
	return domain.Image{
		URL:    i.URL,
		Width:  i.Width,
		Height: i.Height,
	}
}
