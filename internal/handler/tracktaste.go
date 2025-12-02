package handler

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/t1nyb0x/jamberry/internal/version"
)

// handleTrackTaste はTrackTasteのヘルスチェックを行います
func (h *Handler) handleTrackTaste(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// 即時応答（thinking状態）
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	}); err != nil {
		slog.Error("failed to defer response", "error", err)
		return
	}

	// ヘルスチェック実行
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	health, err := h.ttClient.FetchHealth(ctx)
	if err != nil {
		slog.Warn("TrackTaste health check failed", "error", err)
		h.responder.EditResponse(s, i, "❌ TrackTaste API への接続に失敗しました。")
		return
	}

	// Embed作成
	embed := &discordgo.MessageEmbed{
		Title: "🎵 TrackTaste Status",
		Color: 0x1DB954, // Spotify green
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Status",
				Value:  fmt.Sprintf("```%s```", health.Status),
				Inline: true,
			},
			{
				Name:   "Version",
				Value:  fmt.Sprintf("```%s```", health.Version),
				Inline: true,
			},
			{
				Name:   "Uptime",
				Value:  fmt.Sprintf("```%s```", health.Uptime),
				Inline: true,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("jamberry v%s", version.GetVersion()),
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	// GitCommitがあれば追加
	if health.GitCommit != "" {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "Git Commit",
			Value:  fmt.Sprintf("```%s```", health.GitCommit),
			Inline: true,
		})
	}

	// BuildTimeがあれば追加
	if health.BuildTime != "" {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "Build Time",
			Value:  fmt.Sprintf("```%s```", health.BuildTime),
			Inline: true,
		})
	}

	// Servicesがあれば追加
	if len(health.Services) > 0 {
		servicesText := ""
		for name, status := range health.Services {
			icon := "✅"
			if status != "enabled" {
				icon = "❌"
			}
			servicesText += fmt.Sprintf("%s %s\n", icon, name)
		}
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "Services",
			Value:  servicesText,
			Inline: false,
		})
	}

	// 応答を編集
	if _, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{embed},
	}); err != nil {
		slog.Error("failed to edit response", "error", err)
	}
}
