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
			Name:        "jam",
			Description: "jamberry - Spotify 情報取得 Bot",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "track",
					Description: "トラックの詳細情報を取得します",
					Options: []*discordgo.ApplicationCommandOption{
						urlOption,
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "artist",
					Description: "アーティストの詳細情報を取得します",
					Options: []*discordgo.ApplicationCommandOption{
						urlOption,
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "album",
					Description: "アルバムの詳細情報を取得します",
					Options: []*discordgo.ApplicationCommandOption{
						urlOption,
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "recommend",
					Description: "トラックに基づくおすすめ楽曲を取得します",
					Options: []*discordgo.ApplicationCommandOption{
						urlOption,
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "mode",
							Description: "レコメンドモード（similar: 雰囲気重視, related: 関連性重視, balanced: バランス）",
							Required:    false,
							Choices: []*discordgo.ApplicationCommandOptionChoice{
								{
									Name:  "バランス（デフォルト）",
									Value: "balanced",
								},
								{
									Name:  "雰囲気重視",
									Value: "similar",
								},
								{
									Name:  "関連性重視",
									Value: "related",
								},
							},
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "search",
					Description: "トラックを検索します",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "query",
							Description: "検索キーワードを入力",
							Required:    true,
						},
					},
				},
			},
		},
		{
			Name:        "tracktaste",
			Description: "TrackTaste API のステータスを確認します",
		},
	}
}
