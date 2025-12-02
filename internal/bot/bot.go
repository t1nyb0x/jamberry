package bot

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/t1nyb0x/jamberry/internal/version"
)

// Bot はDiscord Botを表します
type Bot struct {
	session  *discordgo.Session
	commands []*discordgo.ApplicationCommand
}

// New は新しいBotを作成します
func New(token string) (*Bot, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	session.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages

	return &Bot{
		session:  session,
		commands: Commands(),
	}, nil
}

// Session はDiscordセッションを返します
func (b *Bot) Session() *discordgo.Session {
	return b.session
}

// Start はBotを起動します
func (b *Bot) Start() error {
	if err := b.session.Open(); err != nil {
		return err
	}

	slog.Info("connected to discord")

	// ステータスにバージョン情報を表示
	status := "v" + version.GetVersion()
	if err := b.session.UpdateGameStatus(0, status); err != nil {
		slog.Warn("failed to update status", "error", err)
	} else {
		slog.Info("updated bot status", "status", status)
	}

	// スラッシュコマンドの登録
	for _, cmd := range b.commands {
		registered, err := b.session.ApplicationCommandCreate(b.session.State.User.ID, "", cmd)
		if err != nil {
			slog.Error("failed to register command", "command", cmd.Name, "error", err)
			continue
		}
		slog.Info("registered command", "command", registered.Name)
	}

	return nil
}

// Stop はBotを停止します
func (b *Bot) Stop() error {
	return b.session.Close()
}

// AddHandler はイベントハンドラーを追加します
func (b *Bot) AddHandler(handler interface{}) {
	b.session.AddHandler(handler)
}
