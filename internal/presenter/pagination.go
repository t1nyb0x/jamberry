package presenter

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/t1nyb0x/jamberry/internal/domain"
)

// BuildRecommendEmbed はレコメンド結果のEmbedを構築します
func BuildRecommendEmbed(originalTrackName string, items []domain.SimilarTrack, page, pageSize, total int) *discordgo.MessageEmbed {
	start := page * pageSize
	end := start + pageSize
	if end > len(items) {
		end = len(items)
	}
	displayItems := items[start:end]

	description := fmt.Sprintf("「%s」に基づくレコメンド (%d-%d / %d 件)", originalTrackName, start+1, end, total)

	var trackListParts []string
	for i, track := range displayItems {
		// アーティスト名をalbum.artistsから取得
		artistNames := make([]string, len(track.Album.Artists))
		for j, a := range track.Album.Artists {
			artistNames[j] = a.Name
		}
		artistStr := strings.Join(artistNames, ", ")

		trackListParts = append(trackListParts, fmt.Sprintf(
			"**%d. %s** - %s\n📀 %s\n🔗 [Spotify](%s)",
			start+i+1, track.Name, artistStr, track.Album.Name, track.URL,
		))
	}

	return &discordgo.MessageEmbed{
		Title:       "🎶 おすすめトラック",
		Description: description + "\n\n" + strings.Join(trackListParts, "\n\n"),
		Color:       SpotifyGreen,
	}
}

// BuildSearchEmbed は検索結果のEmbedを構築します
func BuildSearchEmbed(query string, items []domain.Track, page, pageSize, total int) *discordgo.MessageEmbed {
	start := page * pageSize
	end := start + pageSize
	if end > len(items) {
		end = len(items)
	}
	displayItems := items[start:end]

	description := fmt.Sprintf("「%s」の検索結果 (%d-%d / %d 件)", query, start+1, end, total)

	var trackListParts []string
	for i, track := range displayItems {
		artistStr := JoinArtistNames(track.Artists)
		trackListParts = append(trackListParts, fmt.Sprintf(
			"**%d. %s** - %s\n📀 %s\n🔗 [Spotify](%s)",
			start+i+1, track.Name, artistStr, track.Album.Name, track.URL,
		))
	}

	return &discordgo.MessageEmbed{
		Title:       "🔍 検索結果",
		Description: description + "\n\n" + strings.Join(trackListParts, "\n\n"),
		Color:       SpotifyGreen,
	}
}

// BuildPaginationButtons はページングボタンを構築します
func BuildPaginationButtons(messageID string, page, totalPages int) []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "◀ 前へ",
					Style:    discordgo.SecondaryButton,
					CustomID: fmt.Sprintf("page_prev:%s:%d", messageID, page),
					Disabled: page == 0,
				},
				discordgo.Button{
					Label:    "次へ ▶",
					Style:    discordgo.SecondaryButton,
					CustomID: fmt.Sprintf("page_next:%s:%d", messageID, page),
					Disabled: page >= totalPages-1,
				},
				discordgo.Button{
					Label:    "👁 自分も見る",
					Style:    discordgo.PrimaryButton,
					CustomID: fmt.Sprintf("view_own:%s", messageID),
				},
			},
		},
	}
}
