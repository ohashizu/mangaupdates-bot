package mubot

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/snowflake/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type Config struct {
	Token       string
	DevMode     bool
	DevServerID snowflake.ID
}

type Bot struct {
	Client      bot.Client
	Logger      *slog.Logger
	MuToken     string
	MongoClient *mongo.Client
	Version     string
	Config      Config
}

func New(version string) *Bot {
	devServerID, _ := strconv.Atoi(os.Getenv("DEV_SERVER_ID"))

	logger := slog.Default()
	logger.Info("Starting bot version: " + version)

	return &Bot{
		Client:      nil,
		Logger:      logger,
		MuToken:     "",
		MongoClient: nil,
		Version:     version,
		Config: Config{
			Token:       os.Getenv("TOKEN"),
			DevMode:     os.Getenv("DEV_MODE") == "true",
			DevServerID: snowflake.ID(devServerID),
		},
	}
}

func (b *Bot) Setup(listeners ...bot.EventListener) bot.Client {
	client, err := disgo.New(
		b.Config.Token,
		bot.WithLogger(b.Logger),
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuilds,
			),
		),
		bot.WithCacheConfigOpts(
			cache.WithCaches(cache.FlagGuilds),
		),
		bot.WithEventListeners(listeners...),
	)
	if err != nil {
		b.Logger.Error(fmt.Sprintf("Error while building DisGo client: %s", err))
	}

	return client
}

func (b *Bot) ReadyEvent(_ *events.Ready) {
	err := b.Client.SetPresence(
		context.TODO(),
		gateway.WithPlayingActivity("/help"),
		gateway.WithOnlineStatus(discord.OnlineStatusOnline),
	)
	if err != nil {
		b.Logger.Error(fmt.Sprintf("Error while setting presence: %s", err))
	}

	b.Logger.Info("Bot presence set successfully.")
}
