package handler

import (
	"github.com/bwmarrin/discordgo"
)

// Responder はDiscordレスポンスを管理します
type Responder struct{}

// NewResponder は新しいResponderを作成します
func NewResponder() *Responder {
	return &Responder{}
}

// RespondEphemeral はEphemeralメッセージで応答します
func (r *Responder) RespondEphemeral(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

// DeferReply は遅延レスポンスを開始します
func (r *Responder) DeferReply(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
}

// EditResponse はDeferred応答を編集します
func (r *Responder) EditResponse(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	_, _ = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &content,
	})
}

// EditResponseEmbed はDeferred応答をEmbedで編集します
func (r *Responder) EditResponseEmbed(s *discordgo.Session, i *discordgo.InteractionCreate, emb *discordgo.MessageEmbed) {
	_, _ = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{emb},
	})
}

// EditResponseWithComponents はDeferred応答をEmbed+Componentsで編集します
func (r *Responder) EditResponseWithComponents(s *discordgo.Session, i *discordgo.InteractionCreate, emb *discordgo.MessageEmbed, components []discordgo.MessageComponent) (*discordgo.Message, error) {
	return s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds:     &[]*discordgo.MessageEmbed{emb},
		Components: &components,
	})
}

// UpdateMessage はメッセージを更新します
func (r *Responder) UpdateMessage(s *discordgo.Session, i *discordgo.InteractionCreate, emb *discordgo.MessageEmbed, components []discordgo.MessageComponent) {
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{emb},
			Components: components,
		},
	})
}

// RespondEphemeralWithEmbed はEmbedを含むEphemeralメッセージで応答します
func (r *Responder) RespondEphemeralWithEmbed(s *discordgo.Session, i *discordgo.InteractionCreate, emb *discordgo.MessageEmbed, components []discordgo.MessageComponent) {
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{emb},
			Components: components,
			Flags:      discordgo.MessageFlagsEphemeral,
		},
	})
}

// UpdateEphemeralMessage はEphemeralメッセージを更新します
func (r *Responder) UpdateEphemeralMessage(s *discordgo.Session, i *discordgo.InteractionCreate, emb *discordgo.MessageEmbed, components []discordgo.MessageComponent) {
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{emb},
			Components: components,
		},
	})
}
