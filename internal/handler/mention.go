package handler

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// HandleMessageCreate はメッセージ作成イベントを処理します
func (h *Handler) HandleMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Bot自身のメッセージは無視
	if m.Author.ID == s.State.User.ID {
		return
	}

	// メンションされているか確認
	isMentioned := false
	for _, mention := range m.Mentions {
		if mention.ID == s.State.User.ID {
			isMentioned = true
			break
		}
	}

	if !isMentioned {
		return
	}

	slog.Debug("bot mentioned",
		"guild_id", m.GuildID,
		"channel_id", m.ChannelID,
		"user_id", m.Author.ID,
		"content", m.Content,
	)

	// ヘルプを提案するメッセージを送信
	response := fmt.Sprintf(
		"こんにちは！ 🍇\n"+
			"jamberry の使い方は </help:%s> で確認できます。",
		getHelpCommandID(s),
	)

	if _, err := s.ChannelMessageSendReply(m.ChannelID, response, m.Reference()); err != nil {
		slog.Error("failed to send mention response", "error", err)
	}
}

// getHelpCommandID はhelpコマンドのIDを取得します
func getHelpCommandID(s *discordgo.Session) string {
	commands, err := s.ApplicationCommands(s.State.User.ID, "")
	if err != nil {
		slog.Warn("failed to get application commands", "error", err)
		return "help"
	}

	for _, cmd := range commands {
		if strings.EqualFold(cmd.Name, "help") {
			return cmd.ID
		}
	}

	return "help"
}
