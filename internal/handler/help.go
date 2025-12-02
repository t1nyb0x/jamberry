package handler

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/t1nyb0x/jamberry/internal/version"
)

// handleHelp はヘルプコマンドを処理します
func (h *Handler) handleHelp(s *discordgo.Session, i *discordgo.InteractionCreate) {
	slog.Debug("handling help command")

	embed := &discordgo.MessageEmbed{
		Title:       "🍇 jamberry ヘルプ",
		Description: "Spotify の楽曲・アーティスト・アルバム情報を Discord で検索・共有できる Bot です。",
		Color:       0x1DB954, // Spotify green
		Fields: []*discordgo.MessageEmbedField{
			{
				Name: "🎵 `/jam track <url>`",
				Value: "指定した Spotify トラックの詳細情報を表示します。\n" +
					"• 曲名、アーティスト、アルバム、リリース日\n" +
					"• 再生時間、人気度\n" +
					"• Spotify / KKBOX へのリンク",
				Inline: false,
			},
			{
				Name: "👤 `/jam artist <url>`",
				Value: "指定した Spotify アーティストの詳細情報を表示します。\n" +
					"• アーティスト名、ジャンル\n" +
					"• フォロワー数、人気度\n" +
					"• 代表曲（トップトラック）",
				Inline: false,
			},
			{
				Name: "💿 `/jam album <url>`",
				Value: "指定した Spotify アルバムの詳細情報を表示します。\n" +
					"• アルバム名、アーティスト、リリース日\n" +
					"• 収録曲数、総再生時間\n" +
					"• 収録トラック一覧",
				Inline: false,
			},
			{
				Name: "✨ `/jam recommend <url> [mode]`",
				Value: "指定したトラックに基づくおすすめ楽曲を5件表示します。\n" +
					"• **バランス**: 雰囲気と関連性の両方を考慮（デフォルト）\n" +
					"• **雰囲気重視**: BPM や音圧など音楽的特徴が似た曲\n" +
					"• **関連性重視**: 同じアーティストやジャンルの関連曲\n\n" +
					"📊 **スコアについて**\n" +
					"• 0〜100 の数値で類似度を表します\n" +
					"• 雰囲気スコア: 音楽的特徴（BPM/音圧等）の一致度\n" +
					"• 関連スコア: アーティスト/ジャンルの関連度\n\n" +
					"🎯 **ボーナス倍率**\n" +
					"• **×2.5**: 同一アーティストの別曲\n" +
					"• **×1.3**: 同じグループ/ユニットのメンバー\n" +
					"• **×1.2**: コラボ経験あり / 同じ声優\n" +
					"• **×1.1**: 同じレーベル/プロデューサー\n" +
					"• **×0.5**: 無関係なジャンル（ペナルティ）",
				Inline: false,
			},
			{
				Name: "🔍 `/jam search <query>`",
				Value: "キーワードでトラックを検索します。\n" +
					"• 最大10件の検索結果を表示\n" +
					"• ページネーション対応\n" +
					"• 結果から詳細情報を確認可能",
				Inline: false,
			},
			{
				Name: "🩺 `/tracktaste`",
				Value: "バックエンド API（TrackTaste）のステータスを確認します。\n" +
					"• API のバージョン、稼働時間\n" +
					"• 各サービスの接続状況",
				Inline: false,
			},
			{
				Name: "❓ `/help`",
				Value: "このヘルプメッセージを表示します。",
				Inline: false,
			},
			{
				Name: "📝 対応する入力形式",
				Value: "• **Spotify URL**: `https://open.spotify.com/track/xxxxx`\n" +
					"• **Spotify URI**: `spotify:track:xxxxx`\n" +
					"• **Spotify ID**: `xxxxx`（22文字の英数字）",
				Inline: false,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("jamberry v%s", version.GetVersion()),
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	}); err != nil {
		slog.Error("failed to respond with help", "error", err)
	}
}
