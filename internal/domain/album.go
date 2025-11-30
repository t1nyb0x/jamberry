package domain

// Album はアルバムの基本情報を表すドメインエンティティです
type Album struct {
	ID          string
	Name        string
	URL         string
	ReleaseDate string
	Images      []Image
	Artists     []Artist
}

// AlbumDetail はアルバムの詳細情報を表します
type AlbumDetail struct {
	ID          string
	Name        string
	URL         string
	ReleaseDate string
	Images      []Image
	Artists     []Artist
	Tracks      []AlbumTrack
	Popularity  *int
	UPC         string
	Genres      []string
}

// AlbumTrack はアルバム内のトラック情報を表します
type AlbumTrack struct {
	ID          string
	Name        string
	URL         string
	TrackNumber int
	Artists     []Artist
}

// Image は画像情報を表します
type Image struct {
	URL    string
	Width  int
	Height int
}
