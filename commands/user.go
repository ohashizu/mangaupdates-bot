package commands

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/jckli/mangaupdates-bot/mubot"
	"github.com/jckli/mangaupdates-bot/utils"
)

var userCommand = discord.SlashCommandCreate{
	Name:        "user",
	Description: "Modify user settings",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionSubCommand{
			Name:        "setup",
			Description: "Sets up your user (DMs) for manga updates",
		},
	},
}

func userSetupHandler(e *handler.CommandEvent, b *mubot.Bot) error {
	userId := int64(e.User().ID)

	exists, err := utils.DbUserCheckExists(b, userId)
	if err != nil {
		b.Logger.Error(fmt.Sprintf("Error getting user in userSetupHandler: %s", err.Error()))
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageCreateBuilder().SetEmbeds(utils.DcErrorTechnicalErrorEmbed()).Build(),
		)
	}
	if exists {
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageCreateBuilder().SetEmbeds(errorUserAlreadySetupEmbed()).Build(),
		)
	}

	err = utils.DbAddUser(b, *e.User().GlobalName, userId)
	if err != nil {
		b.Logger.Error(fmt.Sprintf("Error adding user in userSetupHandler: %s", err.Error()))
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageCreateBuilder().SetEmbeds(utils.DcErrorTechnicalErrorEmbed()).Build(),
		)
	}

	embed := discord.NewEmbedBuilder().
		SetTitle("User Setup").
		SetDescription("Great! You (DMs) are now setup for manga updates. You can add manga now using the `/manga add` command.").
		SetColor(0x3083e3).
		Build()
	return e.Respond(
		discord.InteractionResponseTypeCreateMessage,
		discord.NewMessageCreateBuilder().SetEmbeds(embed).Build(),
	)
}
