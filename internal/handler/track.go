package handler

import (
	"context"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/t1nyb0x/jamberry/internal/presenter"
	"github.com/t1nyb0x/jamberry/internal/usecase"
)

// handleTrack はトラック情報取得コマンドを処理します
func (h *Handler) handleTrack(s *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption) {
	if len(options) == 0 {
		slog.Info("validation failed: empty input", "command", "jam track")
		h.responder.RespondEphemeral(s, i, "❌ URL を入力してください。")
		return
	}

	input := options[0].StringValue()

	// DeferReply
	if err := h.responder.DeferReply(s, i); err != nil {
		slog.Error("failed to defer reply", "command", "jam track", "error", err)
		return
	}

	ctx := context.Background()
	output, err := h.trackUseCase.GetTrack(ctx, usecase.TrackInput{Input: input})
	if err != nil {
		if usecase.IsValidationError(err) {
			h.responder.EditResponse(s, i, err.Error())
			return
		}
		h.responder.EditResponse(s, i, err.Error())
		return
	}

	// Embed構築・返信
	emb := presenter.BuildTrackEmbed(output.Track)
	h.responder.EditResponseEmbed(s, i, emb)
	slog.Info("command completed", "command", "jam track", "track_name", output.Track.Name, "track_id", output.Track.ID)
}
