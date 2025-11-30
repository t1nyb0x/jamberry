package bot

import (
	"github.com/bwmarrin/discordgo"
)

// Commands はスラッシュコマンドの定義を返します
func Commands() []*discordgo.ApplicationCommand {
	urlOption := &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "url",
		Description: "Spotify の URL, URI, または ID を入力",
		Required:    true,
	}

	return []*discordgo.ApplicationCommand{
		{
			Name:        "track",
			Description: "Spotifyトラックの詳細情報を取得します",
			Options: []*discordgo.ApplicationCommandOption{
				urlOption,
			},
		},
		{
			Name:        "artist",
			Description: "Spotifyアーティストの詳細情報を取得します",
			Options: []*discordgo.ApplicationCommandOption{
				urlOption,
			},
		},
		{
			Name:        "album",
			Description: "Spotifyアルバムの詳細情報を取得します",
			Options: []*discordgo.ApplicationCommandOption{
				urlOption,
			},
		},
		{
			Name:        "recommend",
			Description: "トラックに基づくおすすめ楽曲を取得します",
			Options: []*discordgo.ApplicationCommandOption{
				urlOption,
			},
		},
		{
			Name:        "search",
			Description: "Spotifyでトラックを検索します",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "query",
					Description: "検索キーワードを入力",
					Required:    true,
				},
			},
		},
	}
}
