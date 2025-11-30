package presenter

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/t1nyb0x/jamberry/internal/domain"
)

// BuildTrackEmbed はトラック情報のEmbedを構築します
func BuildTrackEmbed(track *domain.Track) *discordgo.MessageEmbed {
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
			{
				Name:   "リンク",
				Value:  fmt.Sprintf("[🔗 Spotify で開く](%s)", track.URL),
				Inline: false,
			},
		},
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
func BuildArtistEmbed(artist *domain.ArtistDetail) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Title: "🎤 " + artist.Name,
		Color: SpotifyGreen,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "フォロワー",
				Value:  artist.Followers,
				Inline: true,
			},
			{
				Name:   "リンク",
				Value:  fmt.Sprintf("[🔗 Spotify で開く](%s)", artist.URL),
				Inline: false,
			},
		},
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
func BuildAlbumEmbed(album *domain.AlbumDetail) *discordgo.MessageEmbed {
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
				Value:  fmt.Sprintf("%d 曲", len(album.Tracks)),
				Inline: true,
			},
			{
				Name:   "リンク",
				Value:  fmt.Sprintf("[🔗 Spotify で開く](%s)", album.URL),
				Inline: false,
			},
		},
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
	if len(album.Tracks) > 0 {
		tracks := album.Tracks
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
