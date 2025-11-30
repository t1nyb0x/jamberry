package domain

import "context"

// TrackRepository はトラック情報を取得するリポジトリインターフェースです
type TrackRepository interface {
	// FetchTrack はトラック情報を取得します
	FetchTrack(ctx context.Context, spotifyURL string) (*Track, error)

	// FetchSimilar は類似トラックを取得します
	FetchSimilar(ctx context.Context, spotifyURL string) ([]SimilarTrack, error)

	// SearchTracks はトラックを検索します
	SearchTracks(ctx context.Context, query string) ([]Track, error)
}

// ArtistRepository はアーティスト情報を取得するリポジトリインターフェースです
type ArtistRepository interface {
	// FetchArtist はアーティスト情報を取得します
	FetchArtist(ctx context.Context, spotifyURL string) (*ArtistDetail, error)
}

// AlbumRepository はアルバム情報を取得するリポジトリインターフェースです
type AlbumRepository interface {
	// FetchAlbum はアルバム情報を取得します
	FetchAlbum(ctx context.Context, spotifyURL string) (*AlbumDetail, error)
}

// MusicRepository は音楽情報を取得する統合リポジトリインターフェースです
type MusicRepository interface {
	TrackRepository
	ArtistRepository
	AlbumRepository
}
