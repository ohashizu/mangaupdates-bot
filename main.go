package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/disgoorg/disgo/bot"
	"github.com/jckli/mangaupdates-bot/commands"
	"github.com/jckli/mangaupdates-bot/mubot"
	update_sending "github.com/jckli/mangaupdates-bot/updates"
	"github.com/jckli/mangaupdates-bot/utils"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	mu := mubot.New(os.Getenv("VERSION"))

	h := commands.CommandHandlers(mu)

	client := mu.Setup(
		h,
		bot.NewListenerFunc(mu.ReadyEvent),
	)
	mu.Client = client

	var err error
	muToken, err := utils.MuLogin()
	if err != nil {
		mu.Logger.Error(fmt.Sprintf("Failed to login to MangaUpdates: %s", err))
	}
	mu.MuToken = muToken.Context.SessionToken

	mongoClient, err := utils.DbConnect()
	if err != nil {
		mu.Logger.Error(fmt.Sprintf("Failed to connect to MongoDB: %s", err))
	}
	mu.MongoClient = mongoClient

	if mu.Config.DevMode {
		mu.Logger.Info(
			fmt.Sprintf(
				"Running in dev mode. Syncing commands to server ID: %s",
				mu.Config.DevServerID,
			),
		)
		_, err = client.Rest().
			SetGuildCommands(client.ApplicationID(), mu.Config.DevServerID, commands.CommandList)
	} else {
		mu.Logger.Info(
			"Running in global mode. Syncing commands globally.",
		)
		_, err = client.Rest().SetGlobalCommands(client.ApplicationID(), commands.CommandList)
	}
	if err != nil {
		mu.Logger.Error(fmt.Sprintf("Failed to sync commands: %s", err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.OpenGateway(ctx)
	if err != nil {
		mu.Logger.Error(fmt.Sprintf("Error while connecting: %s", err))
	}
	defer client.Close(context.TODO())

	update_sending.StartRssCheck(mu)

	mu.Logger.Info("Bot is now running.")

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
	shutdown(mu)
}

func shutdown(mu *mubot.Bot) {
	if err := utils.MuLogout(mu); err != nil {
		mu.Logger.Error(fmt.Sprintf("Failed to logout from MangaUpdates: %s", err))
	}
	if err := utils.DbDisconnect(mu); err != nil {
		mu.Logger.Error(fmt.Sprintf("Failed to disconnect from MongoDB: %s", err))
	}
	mu.Logger.Info("Shutting down...")
}
