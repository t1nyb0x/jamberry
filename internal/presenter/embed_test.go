package presenter

import (
	"strings"
	"testing"

	"github.com/t1nyb0x/jamberry/internal/domain"
)

func TestBuildTrackEmbed(t *testing.T) {
	popularity := 85

	tests := []struct {
		name  string
		track *domain.Track
		check func(t *testing.T, embed *EmbedResult)
	}{
		{
			name: "basic track",
			track: &domain.Track{
				ID:         "trackId123",
				Name:       "Test Track",
				URL:        "https://open.spotify.com/track/trackId123",
				DurationMs: 225000,
				Explicit:   false,
				Popularity: &popularity,
				Album: domain.Album{
					Name:        "Test Album",
					ReleaseDate: "2024-01-15",
					Images: []domain.Image{
						{URL: "http://example.com/large.jpg", Width: 640, Height: 640},
					},
				},
				Artists: []domain.Artist{
					{Name: "Artist A"},
					{Name: "Artist B"},
				},
			},
			check: func(t *testing.T, embed *EmbedResult) {
				if embed.Title != "🎵 Test Track" {
					t.Errorf("Title = %s, want 🎵 Test Track", embed.Title)
				}
				if embed.Description != "Artist A, Artist B" {
					t.Errorf("Description = %s, want Artist A, Artist B", embed.Description)
				}
				if embed.URL != "https://open.spotify.com/track/trackId123" {
					t.Errorf("URL = %s, want https://open.spotify.com/track/trackId123", embed.URL)
				}
				if embed.Color != SpotifyGreen {
					t.Errorf("Color = %d, want %d", embed.Color, SpotifyGreen)
				}
				if embed.ThumbnailURL != "http://example.com/large.jpg" {
					t.Errorf("ThumbnailURL = %s, want http://example.com/large.jpg", embed.ThumbnailURL)
				}
			},
		},
		{
			name: "explicit track",
			track: &domain.Track{
				ID:         "trackId123",
				Name:       "Explicit Track",
				URL:        "https://open.spotify.com/track/trackId123",
				DurationMs: 180000,
				Explicit:   true,
				Album: domain.Album{
					Name:        "Test Album",
					ReleaseDate: "2024-01-15",
					Images:      []domain.Image{},
				},
				Artists: []domain.Artist{
					{Name: "Artist A"},
				},
			},
			check: func(t *testing.T, embed *EmbedResult) {
				if embed.Title != "🎵 Explicit Track 🔞" {
					t.Errorf("Title = %s, want 🎵 Explicit Track 🔞", embed.Title)
				}
			},
		},
		{
			name: "track without popularity",
			track: &domain.Track{
				ID:         "trackId123",
				Name:       "Test Track",
				URL:        "https://open.spotify.com/track/trackId123",
				DurationMs: 180000,
				Explicit:   false,
				Popularity: nil,
				Album: domain.Album{
					Name:        "Test Album",
					ReleaseDate: "2024",
					Images:      []domain.Image{},
				},
				Artists: []domain.Artist{
					{Name: "Artist A"},
				},
			},
			check: func(t *testing.T, embed *EmbedResult) {
				// 人気度フィールドが存在しないことを確認
				if embed.HasPopularity {
					t.Errorf("Expected no popularity field when Popularity is nil")
				}
			},
		},
		{
			name: "track without album art",
			track: &domain.Track{
				ID:         "trackId123",
				Name:       "Test Track",
				URL:        "https://open.spotify.com/track/trackId123",
				DurationMs: 180000,
				Explicit:   false,
				Popularity: &popularity,
				Album: domain.Album{
					Name:        "Test Album",
					ReleaseDate: "2024-01-15",
					Images:      []domain.Image{},
				},
				Artists: []domain.Artist{
					{Name: "Artist A"},
				},
			},
			check: func(t *testing.T, embed *EmbedResult) {
				if embed.ThumbnailURL != "" {
					t.Errorf("ThumbnailURL = %s, want empty", embed.ThumbnailURL)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			embed := BuildTrackEmbed(tt.track)
			result := &EmbedResult{
				Title:         embed.Title,
				Description:   embed.Description,
				URL:           embed.URL,
				Color:         embed.Color,
				HasPopularity: false,
			}
			if embed.Thumbnail != nil {
				result.ThumbnailURL = embed.Thumbnail.URL
			}
			for _, field := range embed.Fields {
				if field.Name == "人気度" {
					result.HasPopularity = true
				}
			}
			tt.check(t, result)
		})
	}
}

func TestBuildArtistEmbed(t *testing.T) {
	popularity := 92

	tests := []struct {
		name   string
		artist *domain.ArtistDetail
		check  func(t *testing.T, embed *EmbedResult)
	}{
		{
			name: "artist with all fields",
			artist: &domain.ArtistDetail{
				ID:         "artistId123",
				Name:       "Test Artist",
				URL:        "https://open.spotify.com/artist/artistId123",
				Followers:  "1,234,567",
				Genres:     []string{"Pop", "Rock", "Indie"},
				Popularity: &popularity,
				Images: []domain.Image{
					{URL: "http://example.com/artist.jpg", Width: 640, Height: 640},
				},
			},
			check: func(t *testing.T, embed *EmbedResult) {
				if embed.Title != "🎤 Test Artist" {
					t.Errorf("Title = %s, want 🎤 Test Artist", embed.Title)
				}
				if embed.URL != "https://open.spotify.com/artist/artistId123" {
					t.Errorf("URL = %s, want https://open.spotify.com/artist/artistId123", embed.URL)
				}
				if embed.ThumbnailURL != "http://example.com/artist.jpg" {
					t.Errorf("ThumbnailURL = %s, want http://example.com/artist.jpg", embed.ThumbnailURL)
				}
			},
		},
		{
			name: "artist without genres",
			artist: &domain.ArtistDetail{
				ID:         "artistId123",
				Name:       "Test Artist",
				URL:        "https://open.spotify.com/artist/artistId123",
				Followers:  "1,234,567",
				Genres:     []string{},
				Popularity: &popularity,
				Images:     []domain.Image{},
			},
			check: func(t *testing.T, embed *EmbedResult) {
				if !embed.HasGenreNone {
					t.Errorf("Expected genre to be 'なし' when no genres")
				}
			},
		},
		{
			name: "artist with more than 3 genres",
			artist: &domain.ArtistDetail{
				ID:         "artistId123",
				Name:       "Test Artist",
				URL:        "https://open.spotify.com/artist/artistId123",
				Followers:  "1,234,567",
				Genres:     []string{"Pop", "Rock", "Indie", "Jazz", "Blues"},
				Popularity: &popularity,
				Images:     []domain.Image{},
			},
			check: func(t *testing.T, embed *EmbedResult) {
				// ジャンルが最大3つであることを確認
				if embed.GenreCount > 3 {
					t.Errorf("GenreCount = %d, want <= 3", embed.GenreCount)
				}
			},
		},
		{
			name: "artist without popularity",
			artist: &domain.ArtistDetail{
				ID:         "artistId123",
				Name:       "Test Artist",
				URL:        "https://open.spotify.com/artist/artistId123",
				Followers:  "0",
				Genres:     []string{"Pop"},
				Popularity: nil,
				Images:     []domain.Image{},
			},
			check: func(t *testing.T, embed *EmbedResult) {
				if embed.HasPopularity {
					t.Errorf("Expected no popularity field when Popularity is nil")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			embed := BuildArtistEmbed(tt.artist)
			result := &EmbedResult{
				Title:         embed.Title,
				URL:           embed.URL,
				HasPopularity: false,
				HasGenreNone:  false,
				GenreCount:    0,
			}
			if embed.Thumbnail != nil {
				result.ThumbnailURL = embed.Thumbnail.URL
			}
			for _, field := range embed.Fields {
				if field.Name == "人気度" {
					result.HasPopularity = true
				}
				if field.Name == "ジャンル" {
					if field.Value == "なし" {
						result.HasGenreNone = true
					} else {
						result.GenreCount = strings.Count(field.Value, ",") + 1
					}
				}
			}
			tt.check(t, result)
		})
	}
}

func TestBuildAlbumEmbed(t *testing.T) {
	popularity := 75

	tests := []struct {
		name  string
		album *domain.AlbumDetail
		check func(t *testing.T, embed *EmbedResult)
	}{
		{
			name: "album with all fields",
			album: &domain.AlbumDetail{
				ID:          "albumId123",
				Name:        "Test Album",
				URL:         "https://open.spotify.com/album/albumId123",
				ReleaseDate: "2024-01-15",
				Popularity:  &popularity,
				Artists: []domain.Artist{
					{Name: "Artist A"},
					{Name: "Artist B"},
				},
				Tracks: []domain.AlbumTrack{
					{Name: "Track 1"},
					{Name: "Track 2"},
					{Name: "Track 3"},
					{Name: "Track 4"},
					{Name: "Track 5"},
					{Name: "Track 6"},
					{Name: "Track 7"},
				},
				Images: []domain.Image{
					{URL: "http://example.com/album.jpg", Width: 640, Height: 640},
				},
			},
			check: func(t *testing.T, embed *EmbedResult) {
				if embed.Title != "💿 Test Album" {
					t.Errorf("Title = %s, want 💿 Test Album", embed.Title)
				}
				if embed.Description != "Artist A, Artist B" {
					t.Errorf("Description = %s, want Artist A, Artist B", embed.Description)
				}
				if embed.ThumbnailURL != "http://example.com/album.jpg" {
					t.Errorf("ThumbnailURL = %s, want http://example.com/album.jpg", embed.ThumbnailURL)
				}
				// トラック数が5曲までに制限されていることを確認
				if embed.TrackCount > 5 {
					t.Errorf("TrackCount = %d, want <= 5 in display", embed.TrackCount)
				}
			},
		},
		{
			name: "album without popularity",
			album: &domain.AlbumDetail{
				ID:          "albumId123",
				Name:        "Test Album",
				URL:         "https://open.spotify.com/album/albumId123",
				ReleaseDate: "2024",
				Popularity:  nil,
				Artists: []domain.Artist{
					{Name: "Artist A"},
				},
				Tracks: []domain.AlbumTrack{
					{Name: "Track 1"},
				},
				Images: []domain.Image{},
			},
			check: func(t *testing.T, embed *EmbedResult) {
				if embed.HasPopularity {
					t.Errorf("Expected no popularity field when Popularity is nil")
				}
			},
		},
		{
			name: "album with no tracks",
			album: &domain.AlbumDetail{
				ID:          "albumId123",
				Name:        "Empty Album",
				URL:         "https://open.spotify.com/album/albumId123",
				ReleaseDate: "2024-01-15",
				Popularity:  &popularity,
				Artists: []domain.Artist{
					{Name: "Artist A"},
				},
				Tracks: []domain.AlbumTrack{},
				Images: []domain.Image{},
			},
			check: func(t *testing.T, embed *EmbedResult) {
				if embed.HasTracks {
					t.Errorf("Expected no tracks field when Tracks is empty")
				}
			},
		},
		{
			name: "album with exactly 5 tracks",
			album: &domain.AlbumDetail{
				ID:          "albumId123",
				Name:        "Test Album",
				URL:         "https://open.spotify.com/album/albumId123",
				ReleaseDate: "2024-01-15",
				Popularity:  &popularity,
				Artists: []domain.Artist{
					{Name: "Artist A"},
				},
				Tracks: []domain.AlbumTrack{
					{Name: "Track 1"},
					{Name: "Track 2"},
					{Name: "Track 3"},
					{Name: "Track 4"},
					{Name: "Track 5"},
				},
				Images: []domain.Image{},
			},
			check: func(t *testing.T, embed *EmbedResult) {
				if embed.TrackCount != 5 {
					t.Errorf("TrackCount = %d, want 5", embed.TrackCount)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			embed := BuildAlbumEmbed(tt.album)
			result := &EmbedResult{
				Title:         embed.Title,
				Description:   embed.Description,
				HasPopularity: false,
				HasTracks:     false,
				TrackCount:    0,
			}
			if embed.Thumbnail != nil {
				result.ThumbnailURL = embed.Thumbnail.URL
			}
			for _, field := range embed.Fields {
				if field.Name == "人気度" {
					result.HasPopularity = true
				}
				if field.Name == "収録曲" {
					result.HasTracks = true
					result.TrackCount = strings.Count(field.Value, "\n") + 1
				}
			}
			tt.check(t, result)
		})
	}
}

// EmbedResult はテスト用のEmbed結果を表します
type EmbedResult struct {
	Title         string
	Description   string
	URL           string
	Color         int
	ThumbnailURL  string
	HasPopularity bool
	HasGenreNone  bool
	GenreCount    int
	HasTracks     bool
	TrackCount    int
}
