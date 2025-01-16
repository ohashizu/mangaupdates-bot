package commands

import (
	"github.com/disgoorg/disgo/discord"
)

func errorUserAlreadySetupEmbed() discord.Embed {
	embed := discord.NewEmbedBuilder().
		SetTitle("Error").
		SetDescription("You (your DMs) are already setup.").
		SetColor(0xff4f4f).
		Build()
	return embed
}
