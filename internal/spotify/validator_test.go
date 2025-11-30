package spotify

import (
	"testing"
)

func TestValidateInput_EmptyInput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected EntityType
	}{
		{"empty string", "", EntityTrack},
		{"whitespace only", "   ", EntityTrack},
		{"tab only", "\t", EntityTrack},
		{"newline only", "\n", EntityTrack},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateInput(tt.input, tt.expected)
			if result.Valid {
				t.Errorf("expected invalid, got valid")
			}
			if result.Error != "❌ Spotify の URL / ID として認識できませんでした。" {
				t.Errorf("unexpected error message: %s", result.Error)
			}
		})
	}
}

func TestValidateInput_ValidURL(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedType EntityType
		wantValid    bool
		wantURL      string
		wantID       string
	}{
		{
			name:         "track URL",
			input:        "https://open.spotify.com/track/4iV5W9uYEdYUVa79Axb7Rh",
			expectedType: EntityTrack,
			wantValid:    true,
			wantURL:      "https://open.spotify.com/track/4iV5W9uYEdYUVa79Axb7Rh",
			wantID:       "4iV5W9uYEdYUVa79Axb7Rh",
		},
		{
			name:         "artist URL",
			input:        "https://open.spotify.com/artist/0OdUWJ0sBjDrqHygGUXeCF",
			expectedType: EntityArtist,
			wantValid:    true,
			wantURL:      "https://open.spotify.com/artist/0OdUWJ0sBjDrqHygGUXeCF",
			wantID:       "0OdUWJ0sBjDrqHygGUXeCF",
		},
		{
			name:         "album URL",
			input:        "https://open.spotify.com/album/4aawyAB9vmqN3uQ7FjRGTy",
			expectedType: EntityAlbum,
			wantValid:    true,
			wantURL:      "https://open.spotify.com/album/4aawyAB9vmqN3uQ7FjRGTy",
			wantID:       "4aawyAB9vmqN3uQ7FjRGTy",
		},
		{
			name:         "URL with query params",
			input:        "https://open.spotify.com/track/4iV5W9uYEdYUVa79Axb7Rh?si=abc123",
			expectedType: EntityTrack,
			wantValid:    true,
			wantURL:      "https://open.spotify.com/track/4iV5W9uYEdYUVa79Axb7Rh",
			wantID:       "4iV5W9uYEdYUVa79Axb7Rh",
		},
		{
			name:         "URL with trailing slash",
			input:        "https://open.spotify.com/track/4iV5W9uYEdYUVa79Axb7Rh/",
			expectedType: EntityTrack,
			wantValid:    true, // Trailing slash is handled by TrimSpace
			wantURL:      "https://open.spotify.com/track/4iV5W9uYEdYUVa79Axb7Rh",
			wantID:       "4iV5W9uYEdYUVa79Axb7Rh",
		},
		{
			name:         "URL with whitespace",
			input:        "  https://open.spotify.com/track/4iV5W9uYEdYUVa79Axb7Rh  ",
			expectedType: EntityTrack,
			wantValid:    true,
			wantURL:      "https://open.spotify.com/track/4iV5W9uYEdYUVa79Axb7Rh",
			wantID:       "4iV5W9uYEdYUVa79Axb7Rh",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateInput(tt.input, tt.expectedType)
			if result.Valid != tt.wantValid {
				t.Errorf("expected valid=%v, got valid=%v, error=%s", tt.wantValid, result.Valid, result.Error)
			}
			if tt.wantValid {
				if result.URL != tt.wantURL {
					t.Errorf("expected URL=%s, got URL=%s", tt.wantURL, result.URL)
				}
				if result.ID != tt.wantID {
					t.Errorf("expected ID=%s, got ID=%s", tt.wantID, result.ID)
				}
				if result.EntityType != tt.expectedType {
					t.Errorf("expected EntityType=%s, got EntityType=%s", tt.expectedType, result.EntityType)
				}
			}
		})
	}
}

func TestValidateInput_InvalidURL(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedType EntityType
		wantError    string
	}{
		{
			name:         "non-spotify domain",
			input:        "https://example.com/track/4iV5W9uYEdYUVa79Axb7Rh",
			expectedType: EntityTrack,
			wantError:    "❌ Spotify の URL / ID として認識できませんでした。",
		},
		{
			name:         "youtube URL",
			input:        "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
			expectedType: EntityTrack,
			wantError:    "❌ Spotify の URL / ID として認識できませんでした。",
		},
		{
			name:         "spotify embed URL",
			input:        "https://embed.spotify.com/track/4iV5W9uYEdYUVa79Axb7Rh",
			expectedType: EntityTrack,
			wantError:    "❌ Spotify の URL / ID として認識できませんでした。",
		},
		{
			name:         "http instead of https",
			input:        "http://open.spotify.com/track/4iV5W9uYEdYUVa79Axb7Rh",
			expectedType: EntityTrack,
			wantError:    "", // http is also valid
		},
		{
			name:         "missing path",
			input:        "https://open.spotify.com/",
			expectedType: EntityTrack,
			wantError:    "❌ Spotify の URL / ID として認識できませんでした。",
		},
		{
			name:         "invalid ID format in URL",
			input:        "https://open.spotify.com/track/abc",
			expectedType: EntityTrack,
			wantError:    "❌ Spotify の URL / ID として認識できませんでした。",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateInput(tt.input, tt.expectedType)
			if tt.wantError != "" {
				if result.Valid {
					t.Errorf("expected invalid, got valid")
				}
				if result.Error != tt.wantError {
					t.Errorf("expected error=%s, got error=%s", tt.wantError, result.Error)
				}
			}
		})
	}
}

func TestValidateInput_ValidURI(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedType EntityType
		wantValid    bool
		wantURL      string
		wantID       string
	}{
		{
			name:         "track URI",
			input:        "spotify:track:4iV5W9uYEdYUVa79Axb7Rh",
			expectedType: EntityTrack,
			wantValid:    true,
			wantURL:      "https://open.spotify.com/track/4iV5W9uYEdYUVa79Axb7Rh",
			wantID:       "4iV5W9uYEdYUVa79Axb7Rh",
		},
		{
			name:         "artist URI",
			input:        "spotify:artist:0OdUWJ0sBjDrqHygGUXeCF",
			expectedType: EntityArtist,
			wantValid:    true,
			wantURL:      "https://open.spotify.com/artist/0OdUWJ0sBjDrqHygGUXeCF",
			wantID:       "0OdUWJ0sBjDrqHygGUXeCF",
		},
		{
			name:         "album URI",
			input:        "spotify:album:4aawyAB9vmqN3uQ7FjRGTy",
			expectedType: EntityAlbum,
			wantValid:    true,
			wantURL:      "https://open.spotify.com/album/4aawyAB9vmqN3uQ7FjRGTy",
			wantID:       "4aawyAB9vmqN3uQ7FjRGTy",
		},
		{
			name:         "URI with whitespace",
			input:        "  spotify:track:4iV5W9uYEdYUVa79Axb7Rh  ",
			expectedType: EntityTrack,
			wantValid:    true,
			wantURL:      "https://open.spotify.com/track/4iV5W9uYEdYUVa79Axb7Rh",
			wantID:       "4iV5W9uYEdYUVa79Axb7Rh",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateInput(tt.input, tt.expectedType)
			if result.Valid != tt.wantValid {
				t.Errorf("expected valid=%v, got valid=%v, error=%s", tt.wantValid, result.Valid, result.Error)
			}
			if tt.wantValid {
				if result.URL != tt.wantURL {
					t.Errorf("expected URL=%s, got URL=%s", tt.wantURL, result.URL)
				}
				if result.ID != tt.wantID {
					t.Errorf("expected ID=%s, got ID=%s", tt.wantID, result.ID)
				}
			}
		})
	}
}

func TestValidateInput_InvalidURI(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedType EntityType
		wantError    string
	}{
		{
			name:         "not spotify prefix",
			input:        "other:track:4iV5W9uYEdYUVa79Axb7Rh",
			expectedType: EntityTrack,
			wantError:    "❌ Spotify の URL / ID として認識できませんでした。",
		},
		{
			name:         "too few segments",
			input:        "spotify:track",
			expectedType: EntityTrack,
			wantError:    "❌ Spotify の URL / ID として認識できませんでした。",
		},
		{
			name:         "too many segments",
			input:        "spotify:track:4iV5W9uYEdYUVa79Axb7Rh:extra",
			expectedType: EntityTrack,
			wantError:    "❌ Spotify の URL / ID として認識できませんでした。",
		},
		{
			name:         "invalid ID in URI",
			input:        "spotify:track:abc",
			expectedType: EntityTrack,
			wantError:    "❌ Spotify の URL / ID として認識できませんでした。",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateInput(tt.input, tt.expectedType)
			if result.Valid {
				t.Errorf("expected invalid, got valid")
			}
			if result.Error != tt.wantError {
				t.Errorf("expected error=%s, got error=%s", tt.wantError, result.Error)
			}
		})
	}
}

func TestValidateInput_ValidRawID(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedType EntityType
		wantValid    bool
		wantURL      string
		wantID       string
	}{
		{
			name:         "valid 22-char ID for track",
			input:        "4iV5W9uYEdYUVa79Axb7Rh",
			expectedType: EntityTrack,
			wantValid:    true,
			wantURL:      "https://open.spotify.com/track/4iV5W9uYEdYUVa79Axb7Rh",
			wantID:       "4iV5W9uYEdYUVa79Axb7Rh",
		},
		{
			name:         "valid 22-char ID for artist",
			input:        "0OdUWJ0sBjDrqHygGUXeCF",
			expectedType: EntityArtist,
			wantValid:    true,
			wantURL:      "https://open.spotify.com/artist/0OdUWJ0sBjDrqHygGUXeCF",
			wantID:       "0OdUWJ0sBjDrqHygGUXeCF",
		},
		{
			name:         "valid 22-char ID for album",
			input:        "4aawyAB9vmqN3uQ7FjRGTy",
			expectedType: EntityAlbum,
			wantValid:    true,
			wantURL:      "https://open.spotify.com/album/4aawyAB9vmqN3uQ7FjRGTy",
			wantID:       "4aawyAB9vmqN3uQ7FjRGTy",
		},
		{
			name:         "ID with whitespace",
			input:        "  4iV5W9uYEdYUVa79Axb7Rh  ",
			expectedType: EntityTrack,
			wantValid:    true,
			wantURL:      "https://open.spotify.com/track/4iV5W9uYEdYUVa79Axb7Rh",
			wantID:       "4iV5W9uYEdYUVa79Axb7Rh",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateInput(tt.input, tt.expectedType)
			if result.Valid != tt.wantValid {
				t.Errorf("expected valid=%v, got valid=%v, error=%s", tt.wantValid, result.Valid, result.Error)
			}
			if tt.wantValid {
				if result.URL != tt.wantURL {
					t.Errorf("expected URL=%s, got URL=%s", tt.wantURL, result.URL)
				}
				if result.ID != tt.wantID {
					t.Errorf("expected ID=%s, got ID=%s", tt.wantID, result.ID)
				}
				if result.EntityType != tt.expectedType {
					t.Errorf("expected EntityType=%s, got EntityType=%s", tt.expectedType, result.EntityType)
				}
			}
		})
	}
}

func TestValidateInput_InvalidRawID(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedType EntityType
		wantError    string
	}{
		{
			name:         "too short ID",
			input:        "abc",
			expectedType: EntityTrack,
			wantError:    "❌ Spotify の URL / ID として認識できませんでした。",
		},
		{
			name:         "too long ID",
			input:        "4iV5W9uYEdYUVa79Axb7Rh123",
			expectedType: EntityTrack,
			wantError:    "❌ Spotify の URL / ID として認識できませんでした。",
		},
		{
			name:         "special characters",
			input:        "4iV5W9uYEdYUVa79Axb7!@",
			expectedType: EntityTrack,
			wantError:    "❌ Spotify の URL / ID として認識できませんでした。",
		},
		{
			name:         "21 characters",
			input:        "4iV5W9uYEdYUVa79Axb7R",
			expectedType: EntityTrack,
			wantError:    "❌ Spotify の URL / ID として認識できませんでした。",
		},
		{
			name:         "23 characters",
			input:        "4iV5W9uYEdYUVa79Axb7Rha",
			expectedType: EntityTrack,
			wantError:    "❌ Spotify の URL / ID として認識できませんでした。",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateInput(tt.input, tt.expectedType)
			if result.Valid {
				t.Errorf("expected invalid, got valid")
			}
			if result.Error != tt.wantError {
				t.Errorf("expected error=%s, got error=%s", tt.wantError, result.Error)
			}
		})
	}
}

func TestValidateInput_EntityTypeMismatch(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedType EntityType
		wantError    string
	}{
		{
			name:         "track URL with artist expected",
			input:        "https://open.spotify.com/track/4iV5W9uYEdYUVa79Axb7Rh",
			expectedType: EntityArtist,
			wantError:    "❌ Spotify の ArtistURL を入力してください",
		},
		{
			name:         "track URL with album expected",
			input:        "https://open.spotify.com/track/4iV5W9uYEdYUVa79Axb7Rh",
			expectedType: EntityAlbum,
			wantError:    "❌ Spotify の AlbumURL を入力してください",
		},
		{
			name:         "artist URL with track expected",
			input:        "https://open.spotify.com/artist/0OdUWJ0sBjDrqHygGUXeCF",
			expectedType: EntityTrack,
			wantError:    "❌ Spotify の TrackURL を入力してください",
		},
		{
			name:         "artist URL with album expected",
			input:        "https://open.spotify.com/artist/0OdUWJ0sBjDrqHygGUXeCF",
			expectedType: EntityAlbum,
			wantError:    "❌ Spotify の AlbumURL を入力してください",
		},
		{
			name:         "album URL with track expected",
			input:        "https://open.spotify.com/album/4aawyAB9vmqN3uQ7FjRGTy",
			expectedType: EntityTrack,
			wantError:    "❌ Spotify の TrackURL を入力してください",
		},
		{
			name:         "album URL with artist expected",
			input:        "https://open.spotify.com/album/4aawyAB9vmqN3uQ7FjRGTy",
			expectedType: EntityArtist,
			wantError:    "❌ Spotify の ArtistURL を入力してください",
		},
		{
			name:         "track URI with artist expected",
			input:        "spotify:track:4iV5W9uYEdYUVa79Axb7Rh",
			expectedType: EntityArtist,
			wantError:    "❌ Spotify の ArtistURL を入力してください",
		},
		{
			name:         "artist URI with track expected",
			input:        "spotify:artist:0OdUWJ0sBjDrqHygGUXeCF",
			expectedType: EntityTrack,
			wantError:    "❌ Spotify の TrackURL を入力してください",
		},
		{
			name:         "album URI with artist expected",
			input:        "spotify:album:4aawyAB9vmqN3uQ7FjRGTy",
			expectedType: EntityArtist,
			wantError:    "❌ Spotify の ArtistURL を入力してください",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateInput(tt.input, tt.expectedType)
			if result.Valid {
				t.Errorf("expected invalid, got valid")
			}
			if result.Error != tt.wantError {
				t.Errorf("expected error=%s, got error=%s", tt.wantError, result.Error)
			}
		})
	}
}

func TestGetEntityMismatchError(t *testing.T) {
	tests := []struct {
		entityType EntityType
		expected   string
	}{
		{EntityTrack, "❌ Spotify の TrackURL を入力してください"},
		{EntityArtist, "❌ Spotify の ArtistURL を入力してください"},
		{EntityAlbum, "❌ Spotify の AlbumURL を入力してください"},
		{EntityUnknown, "❌ Spotify の URL / ID として認識できませんでした。"},
	}

	for _, tt := range tests {
		t.Run(string(tt.entityType), func(t *testing.T) {
			result := getEntityMismatchError(tt.entityType)
			if result != tt.expected {
				t.Errorf("expected=%s, got=%s", tt.expected, result)
			}
		})
	}
}

func TestSpotifyIDRegex(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"4iV5W9uYEdYUVa79Axb7Rh", true},
		{"0OdUWJ0sBjDrqHygGUXeCF", true},
		{"4aawyAB9vmqN3uQ7FjRGTy", true},
		{"ABC123abc456ABC123abc4", true},
		{"abc", false},
		{"4iV5W9uYEdYUVa79Axb7Rh1", false},
		{"4iV5W9uYEdYUVa79Axb7R", false},
		{"4iV5W9uYEdYUVa79Axb7!@", false},
		{"", false},
		{"4iV5W9uYEdYUVa79-xb7Rh", false},
		{"4iV5W9uYEdYUVa79_xb7Rh", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := SpotifyIDRegex.MatchString(tt.input)
			if result != tt.expected {
				t.Errorf("SpotifyIDRegex.MatchString(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}
