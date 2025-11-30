package presenter

import (
	"fmt"
	"strings"

	"github.com/t1nyb0x/jamberry/internal/domain"
)

const (
	// SpotifyGreen はSpotifyのブランドカラーです
	SpotifyGreen = 0x1DB954
)

// FormatDuration はミリ秒を M:SS 形式にフォーマットします
func FormatDuration(ms int) string {
	seconds := ms / 1000
	minutes := seconds / 60
	secs := seconds % 60
	return fmt.Sprintf("%d:%02d", minutes, secs)
}

// FormatNumber は数値をカンマ区切りにフォーマットします
func FormatNumber(n int) string {
	s := fmt.Sprintf("%d", n)
	if len(s) <= 3 {
		return s
	}

	var result []byte
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result = append(result, ',')
		}
		result = append(result, byte(c))
	}
	return string(result)
}

// GetLargestImage は最大サイズの画像URLを返します
func GetLargestImage(images []domain.Image) string {
	if len(images) == 0 {
		return ""
	}

	largest := images[0]
	for _, img := range images[1:] {
		if img.Width*img.Height > largest.Width*largest.Height {
			largest = img
		}
	}
	return largest.URL
}

// JoinArtistNames はアーティスト名をカンマ区切りで結合します
func JoinArtistNames(artists []domain.Artist) string {
	names := make([]string, len(artists))
	for i, a := range artists {
		names[i] = a.Name
	}
	return strings.Join(names, ", ")
}
