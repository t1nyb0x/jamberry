package presenter

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/t1nyb0x/jamberry/internal/domain"
)

// getModeLabel はレコメンドモードの日本語ラベルを返します
func getModeLabel(mode domain.RecommendMode) string {
	switch mode {
	case domain.RecommendModeSimilar:
		return "雰囲気重視"
	case domain.RecommendModeRelated:
		return "関連性重視"
	case domain.RecommendModeBalanced:
		return "バランス"
	default:
		return "バランス"
	}
}

// BuildRecommendEmbed はレコメンド結果のEmbedを構築します
func BuildRecommendEmbed(originalTrackName string, items []domain.SimilarTrack, page, pageSize, total int, mode domain.RecommendMode) *discordgo.MessageEmbed {
	start := page * pageSize
	end := start + pageSize
	if end > len(items) {
		end = len(items)
	}
	displayItems := items[start:end]

	modeLabel := getModeLabel(mode)
	description := fmt.Sprintf("「%s」に基づくレコメンド\n**モード**: %s (%d-%d / %d 件)", originalTrackName, modeLabel, start+1, end, total)

	var trackListParts []string
	for i, track := range displayItems {
		// アーティスト名を取得
		var artistStr string
		if len(track.Artists) > 0 {
			artistNames := make([]string, len(track.Artists))
			for j, a := range track.Artists {
				artistNames[j] = a.Name
			}
			artistStr = strings.Join(artistNames, ", ")
		} else if len(track.Album.Artists) > 0 {
			// フォールバック: albumのartistsを使用
			artistNames := make([]string, len(track.Album.Artists))
			for j, a := range track.Album.Artists {
				artistNames[j] = a.Name
			}
			artistStr = strings.Join(artistNames, ", ")
		}

		// 番号と曲名（太字）
		trackInfo := fmt.Sprintf("**%d. %s**", start+i+1, track.Name)

		// アーティスト名
		trackInfo += fmt.Sprintf("🎤 %s", artistStr)

		// アルバム名（あれば）
		if track.Album.Name != "" {
			trackInfo += fmt.Sprintf("\n📀 %s", track.Album.Name)
		}

		// スコア表示（final_scoreを優先、なければsimilarity_scoreを使用）
		if track.FinalScore != nil {
			if track.GenreBonus != nil && *track.GenreBonus != 1.0 {
				trackInfo += fmt.Sprintf("\n✨ %.2f (×%.1f)", *track.FinalScore, *track.GenreBonus)
			} else {
				trackInfo += fmt.Sprintf("\n✨ %.2f", *track.FinalScore)
			}
		} else if track.SimilarityScore != nil {
			trackInfo += fmt.Sprintf("\n✨ %.0f%%", *track.SimilarityScore*100)
		}

		// Spotifyリンク
		spotifyURL := track.URL
		if spotifyURL == "" && track.ID != "" {
			// v2 APIではURLが含まれないので、IDからURLを構築
			spotifyURL = fmt.Sprintf("https://open.spotify.com/track/%s", track.ID)
		}
		if spotifyURL != "" {
			trackInfo += fmt.Sprintf("\n🔗 [Spotify](%s)", spotifyURL)
		}

		trackListParts = append(trackListParts, trackInfo)
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
