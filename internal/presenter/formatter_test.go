package presenter

import (
	"testing"

	"github.com/t1nyb0x/jamberry/internal/domain"
)

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		ms       int
		expected string
	}{
		{"zero", 0, "0:00"},
		{"one second", 1000, "0:01"},
		{"59 seconds", 59000, "0:59"},
		{"one minute", 60000, "1:00"},
		{"one minute one second", 61000, "1:01"},
		{"3:45 (typical song)", 225000, "3:45"},
		{"10 minutes", 600000, "10:00"},
		{"10:30", 630000, "10:30"},
		{"over hour", 3661000, "61:01"},
		{"5:00", 300000, "5:00"},
		{"9:59", 599000, "9:59"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatDuration(tt.ms)
			if result != tt.expected {
				t.Errorf("FormatDuration(%d) = %s, want %s", tt.ms, result, tt.expected)
			}
		})
	}
}

func TestFormatNumber(t *testing.T) {
	tests := []struct {
		name     string
		n        int
		expected string
	}{
		{"zero", 0, "0"},
		{"single digit", 5, "5"},
		{"two digits", 42, "42"},
		{"three digits", 123, "123"},
		{"four digits", 1234, "1,234"},
		{"five digits", 12345, "12,345"},
		{"six digits", 123456, "123,456"},
		{"seven digits", 1234567, "1,234,567"},
		{"million", 1000000, "1,000,000"},
		{"large number", 123456789, "123,456,789"},
		{"typical followers", 5823914, "5,823,914"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatNumber(tt.n)
			if result != tt.expected {
				t.Errorf("FormatNumber(%d) = %s, want %s", tt.n, result, tt.expected)
			}
		})
	}
}

func TestGetLargestImage(t *testing.T) {
	tests := []struct {
		name     string
		images   []domain.Image
		expected string
	}{
		{
			name:     "empty images",
			images:   []domain.Image{},
			expected: "",
		},
		{
			name: "single image",
			images: []domain.Image{
				{URL: "http://example.com/img.jpg", Width: 640, Height: 640},
			},
			expected: "http://example.com/img.jpg",
		},
		{
			name: "multiple images - largest first",
			images: []domain.Image{
				{URL: "http://example.com/large.jpg", Width: 640, Height: 640},
				{URL: "http://example.com/medium.jpg", Width: 300, Height: 300},
				{URL: "http://example.com/small.jpg", Width: 64, Height: 64},
			},
			expected: "http://example.com/large.jpg",
		},
		{
			name: "multiple images - largest last",
			images: []domain.Image{
				{URL: "http://example.com/small.jpg", Width: 64, Height: 64},
				{URL: "http://example.com/medium.jpg", Width: 300, Height: 300},
				{URL: "http://example.com/large.jpg", Width: 640, Height: 640},
			},
			expected: "http://example.com/large.jpg",
		},
		{
			name: "multiple images - largest in middle",
			images: []domain.Image{
				{URL: "http://example.com/small.jpg", Width: 64, Height: 64},
				{URL: "http://example.com/large.jpg", Width: 640, Height: 640},
				{URL: "http://example.com/medium.jpg", Width: 300, Height: 300},
			},
			expected: "http://example.com/large.jpg",
		},
		{
			name: "different aspect ratios",
			images: []domain.Image{
				{URL: "http://example.com/wide.jpg", Width: 800, Height: 400},   // 320000
				{URL: "http://example.com/square.jpg", Width: 640, Height: 640}, // 409600
				{URL: "http://example.com/tall.jpg", Width: 400, Height: 800},   // 320000
			},
			expected: "http://example.com/square.jpg",
		},
		{
			name: "same size images",
			images: []domain.Image{
				{URL: "http://example.com/first.jpg", Width: 300, Height: 300},
				{URL: "http://example.com/second.jpg", Width: 300, Height: 300},
			},
			expected: "http://example.com/first.jpg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetLargestImage(tt.images)
			if result != tt.expected {
				t.Errorf("GetLargestImage() = %s, want %s", result, tt.expected)
			}
		})
	}
}

func TestJoinArtistNames(t *testing.T) {
	tests := []struct {
		name     string
		artists  []domain.Artist
		expected string
	}{
		{
			name:     "empty artists",
			artists:  []domain.Artist{},
			expected: "",
		},
		{
			name: "single artist",
			artists: []domain.Artist{
				{Name: "Artist A"},
			},
			expected: "Artist A",
		},
		{
			name: "two artists",
			artists: []domain.Artist{
				{Name: "Artist A"},
				{Name: "Artist B"},
			},
			expected: "Artist A, Artist B",
		},
		{
			name: "three artists",
			artists: []domain.Artist{
				{Name: "Artist A"},
				{Name: "Artist B"},
				{Name: "Artist C"},
			},
			expected: "Artist A, Artist B, Artist C",
		},
		{
			name: "artists with special characters",
			artists: []domain.Artist{
				{Name: "米津玄師"},
				{Name: "YOASOBI"},
			},
			expected: "米津玄師, YOASOBI",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := JoinArtistNames(tt.artists)
			if result != tt.expected {
				t.Errorf("JoinArtistNames() = %s, want %s", result, tt.expected)
			}
		})
	}
}
