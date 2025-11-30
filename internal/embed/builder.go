package embed

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/t1nyb0x/jamberry/internal/tracktaste"
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
func GetLargestImage(images []tracktaste.Image) string {
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
func JoinArtistNames(artists []tracktaste.Artist) string {
	names := make([]string, len(artists))
	for i, a := range artists {
		names[i] = a.Name
	}
	return strings.Join(names, ", ")
}

// BuildTrackEmbed はトラック情報のEmbedを構築します
func BuildTrackEmbed(track *tracktaste.Track) *discordgo.MessageEmbed {
	title := "🎵 " + track.Name
	if track.Explicit {
		title += " 🔞"
	}

	embed := &discordgo.MessageEmbed{
		Title:       title,
		Description: JoinArtistNames(track.Artists),
		Color:       SpotifyGreen,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "アルバム",
				Value:  track.Album.Name,
				Inline: true,
			},
			{
				Name:   "再生時間",
				Value:  FormatDuration(track.DurationMs),
				Inline: true,
			},
			{
				Name:   "リリース日",
				Value:  track.Album.ReleaseDate,
				Inline: true,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "🔗 Spotify で開く",
		},
		URL: track.URL,
	}

	// 人気度（欠損時は省略）
	if track.Popularity != nil {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "人気度",
			Value:  fmt.Sprintf("%d", *track.Popularity),
			Inline: true,
		})
	}

	// アルバムアート
	if imgURL := GetLargestImage(track.Album.Images); imgURL != "" {
		embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
			URL: imgURL,
		}
	}

	return embed
}

// BuildArtistEmbed はアーティスト情報のEmbedを構築します
func BuildArtistEmbed(artist *tracktaste.ArtistFull) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Title: "🎤 " + artist.Name,
		Color: SpotifyGreen,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "フォロワー",
				Value:  artist.Followers,
				Inline: true,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "🔗 Spotify で開く",
		},
		URL: artist.URL,
	}

	// ジャンル（最大3件、欠損時は「なし」）
	genreValue := "なし"
	if len(artist.Genres) > 0 {
		genres := artist.Genres
		if len(genres) > 3 {
			genres = genres[:3]
		}
		genreValue = strings.Join(genres, ", ")
	}
	embed.Fields = append([]*discordgo.MessageEmbedField{
		{
			Name:   "ジャンル",
			Value:  genreValue,
			Inline: true,
		},
	}, embed.Fields...)

	// 人気度（欠損時は省略）
	if artist.Popularity != nil {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "人気度",
			Value:  fmt.Sprintf("%d", *artist.Popularity),
			Inline: true,
		})
	}

	// アーティスト画像
	if imgURL := GetLargestImage(artist.Images); imgURL != "" {
		embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
			URL: imgURL,
		}
	}

	return embed
}

// BuildAlbumEmbed はアルバム情報のEmbedを構築します
func BuildAlbumEmbed(album *tracktaste.AlbumFull) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Title:       "💿 " + album.Name,
		Description: JoinArtistNames(album.Artists),
		Color:       SpotifyGreen,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "リリース日",
				Value:  album.ReleaseDate,
				Inline: true,
			},
			{
				Name:   "トラック数",
				Value:  fmt.Sprintf("%d 曲", len(album.Tracks.Items)),
				Inline: true,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "🔗 Spotify で開く",
		},
		URL: album.URL,
	}

	// 人気度（欠損時は省略）
	if album.Popularity != nil {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "人気度",
			Value:  fmt.Sprintf("%d", *album.Popularity),
			Inline: true,
		})
	}

	// 収録曲（先頭5曲）
	if len(album.Tracks.Items) > 0 {
		tracks := album.Tracks.Items
		if len(tracks) > 5 {
			tracks = tracks[:5]
		}
		var trackList []string
		for i, t := range tracks {
			trackList = append(trackList, fmt.Sprintf("%d. %s", i+1, t.Name))
		}
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  "収録曲",
			Value: strings.Join(trackList, "\n"),
		})
	}

	// アルバムアート
	if imgURL := GetLargestImage(album.Images); imgURL != "" {
		embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
			URL: imgURL,
		}
	}

	return embed
}

// BuildRecommendEmbed はレコメンド結果のEmbedを構築します
func BuildRecommendEmbed(originalTrackName string, items []tracktaste.SimilarTrack, page, pageSize, total int) *discordgo.MessageEmbed {
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
func BuildSearchEmbed(query string, items []tracktaste.SearchTrack, page, pageSize, total int) *discordgo.MessageEmbed {
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
