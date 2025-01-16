package utils

import (
	"github.com/disgoorg/disgo/discord"
	"html"
)

func DcErrorTechnicalErrorEmbed() discord.Embed {
	embed := discord.NewEmbedBuilder().
		SetTitle("Error").
		SetDescription("Something went wrong. Please ask for assistance in the support server.").
		SetColor(0xff4f4f).
		Build()
	return embed
}

func TruncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	if maxLength <= 3 {
		return s[:maxLength]
	}
	return s[:maxLength-3] + "..."
}

func ParseHTMLEntities(s string) string {
	return html.UnescapeString(s)
}
