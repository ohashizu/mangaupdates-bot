package commands

import (
	"github.com/disgoorg/disgo/discord"
)

func errorServerAlreadySetupEmbed() discord.Embed {
	embed := discord.NewEmbedBuilder().
		SetTitle("Error").
		SetDescription("This server is already setup.").
		SetColor(0xff4f4f).
		Build()
	return embed
}

func errorNoPermissionsEmbed() discord.Embed {
	embed := discord.NewEmbedBuilder().
		SetTitle("Error").
		SetDescription("You do not have permission to run this command.").
		SetColor(0xff4f4f).
		Build()
	return embed
}
